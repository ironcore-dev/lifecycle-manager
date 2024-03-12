// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package scheduler

import (
	"context"
	"log/slog"
	"sync"
	"time"

	commonv1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/common/v1alpha1"
	"github.com/jellydator/ttlcache/v3"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type LifecycleObject interface {
	client.Object
	LastScanTime() time.Time
}

type JobType string

var EvictionReason = map[ttlcache.EvictionReason]string{
	ttlcache.EvictionReasonDeleted:         "deleted",
	ttlcache.EvictionReasonCapacityReached: "capacity reached",
	ttlcache.EvictionReasonExpired:         "ttl expired",
}

type Option[T LifecycleObject] func(scheduler *Scheduler[T])

type Scheduler[T LifecycleObject] struct {
	*kubernetes.Clientset
	log          *slog.Logger
	workqueue    *RingBufQueue[T]
	activeJobs   *ttlcache.Cache[string, Task[T]]
	pendingTasks *FIFOQueue[T]
	period       time.Duration
	jobTTL       time.Duration
	workers      uint64

	done chan struct{}

	workersWaitGroup   sync.WaitGroup
	schedulerWaitGroup sync.WaitGroup
	mu                 sync.Mutex

	namespace string
}

// NewScheduler creates a new Scheduler instance with the given parameters.
// The Scheduler manages the scheduling of tasks and their execution by workers.
// It uses a workqueue and a FIFO queue to track tasks.
// The workerCount parameter specifies the number of workers, it reflects the max number of kubernetes Jobs which
// could be run in parallel.
// The activeJobTTL parameter is the time-to-live for active jobs in the cache.
// The logger parameter is the logger instance to use for logging.
// The Scheduler instance is returned.
func NewScheduler[T LifecycleObject](
	logger *slog.Logger,
	cfg *rest.Config,
	namespace string,
	opts ...Option[T],
) *Scheduler[T] {
	kubeClient := kubernetes.NewForConfigOrDie(cfg)
	scheduler := &Scheduler[T]{
		Clientset: kubeClient,
		log:       logger,
		namespace: namespace,
	}
	for _, opt := range opts {
		opt(scheduler)
	}
	return scheduler
}

func WithWorkerCount[T LifecycleObject](workers uint64) Option[T] {
	return func(scheduler *Scheduler[T]) {
		scheduler.workers = workers
		scheduler.workqueue = NewRingBufQueue[T](workers)
		scheduler.done = make(chan struct{}, workers)
	}
}

func WithActiveJobCache[T LifecycleObject](capacity uint64, ttl time.Duration) Option[T] {
	return func(scheduler *Scheduler[T]) {
		scheduler.activeJobs = ttlcache.New[string, Task[T]](
			ttlcache.WithCapacity[string, Task[T]](capacity),
			ttlcache.WithTTL[string, Task[T]](ttl),
			ttlcache.WithDisableTouchOnHit[string, Task[T]]())
	}
}

func WithQueueCapacity[T LifecycleObject](capacity uint64) Option[T] {
	return func(scheduler *Scheduler[T]) {
		scheduler.pendingTasks = NewFIFOQueue[T](capacity)
	}
}

func (s *Scheduler[T]) dropFinishedJob(
	_ context.Context,
	reason ttlcache.EvictionReason,
	item *ttlcache.Item[string, Task[T]],
) {
	s.log.Info("task evicted from active", "task", item.Key(), "reason", EvictionReason[reason])
	s.done <- struct{}{}
}

// Schedule checks if a Task is already enqueued in one of the queues and returns the appropriate RequestResult.
// If the Task is not enqueued in any queue, it attempts to enqueue it in the workqueue.
// If the enqueuing in the workqueue is successful, it returns RequestResult_REQUEST_RESULT_SCHEDULED.
// Otherwise, it attempts to enqueue the Task in the FIFO queue.
// If the enqueuing in the FIFO queue is successful, it returns RequestResult_REQUEST_RESULT_SCHEDULED.
// If the Task is already enqueued in any of the queues, it returns RequestResult_REQUEST_RESULT_SCHEDULED.
// If the enqueuing in both queues fails, it returns RequestResult_REQUEST_RESULT_FAILURE.
func (s *Scheduler[T]) Schedule(item Task[T]) commonv1alpha1.RequestResult {
	enqueued := s.activeJobs.Has(item.Key) || s.workqueue.Has(item.Key) || s.pendingTasks.Has(item.Key)
	if enqueued {
		return commonv1alpha1.RequestResult_REQUEST_RESULT_SCHEDULED
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// try to put item directly to workqueue
	if s.workqueue.TryEnqueue(item) {
		return commonv1alpha1.RequestResult_REQUEST_RESULT_SUCCESS
	}

	// try to push item to FIFO queue
	if s.pendingTasks.TryPush(item) {
		return commonv1alpha1.RequestResult_REQUEST_RESULT_SUCCESS
	}
	return commonv1alpha1.RequestResult_REQUEST_RESULT_FAILURE
}

// ForgetFinishedJob deletes a finished job from the active job tracker.
// It takes a key string as input and removes the corresponding job from the tracker.
// If the job is successfully deleted, it is considered forgotten.
// The method does not return any value.
func (s *Scheduler[T]) ForgetFinishedJob(key string) {
	s.activeJobs.Delete(key)
}

// Start starts the scheduler.
// It sets up a callback function for when a task is evicted from the active jobs cache.
// It starts the active jobs cache and initiates the worker goroutines.
// It then enters a loop where it listens for signals to stop the scheduler or process the queues.
// If it receives a stop signal, it closes the done channel and waits for all worker goroutines to finish.
// It then stops the active jobs cache and exits.
func (s *Scheduler[T]) Start(ctx context.Context) {
	s.activeJobs.OnEviction(s.dropFinishedJob)
	go s.activeJobs.Start()

	for range s.workers {
		s.workersWaitGroup.Add(1)
		go s.workerFunc(ctx)
	}
	s.log.Info("scheduler started")

	s.schedulerWaitGroup.Add(1)
	go s.schedulingLoop(ctx)
	s.schedulerWaitGroup.Wait()
}

func (s *Scheduler[T]) schedulingLoop(ctx context.Context) {
	for {
		select {
		case <-s.done:
			s.processQueues()
		case <-ctx.Done():
			close(s.done)
			s.log.Debug("workqueue", "len", s.workqueue.Len(), "queue", s.workqueue.Print())
			s.log.Debug("pending_tasks", "len", s.pendingTasks.Len(), "queue", s.pendingTasks.Print())
			s.log.Debug("active_jobs", "len", s.activeJobs.Len())
			s.workersWaitGroup.Wait()
			s.activeJobs.Stop()
			s.log.Info("scheduler stopped")
			s.schedulerWaitGroup.Done()
			return
		}
	}
}

func (s *Scheduler[T]) workerFunc(ctx context.Context) {
	for {
		select {
		case <-s.workqueue.Enqueued:
			if s.activeJobs.Len() == int(s.workers) {
				break
			}
			s.mu.Lock()
			task, _ := s.workqueue.Dequeue()
			s.mu.Unlock()
			s.activeJobs.Set(task.Key, task, ttlcache.DefaultTTL)
			if err := s.processJob(ctx, task); err != nil {
				s.log.Error("failed to process task", "error", err.Error())
				s.activeJobs.Delete(task.Key)
			}
		case <-ctx.Done():
			s.log.Debug("stop worker function")
			s.workersWaitGroup.Done()
			return
		}
	}
}

func (s *Scheduler[T]) processQueues() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.log.Debug("checking queues...")

	// if both pre-processing and working queues are empty then do nothing
	if s.workqueue.IsEmpty() && s.pendingTasks.IsEmpty() {
		s.log.Debug("no pending tasks to process")
		return
	}

	// if active jobs track is full then do nothing as workers could not process task
	if s.activeJobs.Len() == int(s.workers) {
		s.log.Debug("no free workers")
		return
	}

	// assuming there is free capacity in active job tracker
	// if workqueue is full then trigger as many workers as active job tracker capacity
	if s.workqueue.IsFull() {
		for range int(s.workers) - s.activeJobs.Len() {
			s.workqueue.Enqueued <- struct{}{}
		}
	}

	// move as many tasks from pre-processing queue to workqueue as possible
	for range s.workqueue.FreeCapacity() {
		if s.pendingTasks.IsEmpty() {
			break
		}
		s.enqueueTask()
	}

	// ensure workers will be triggered to pick ip job from working queue
	if !s.workqueue.IsEmpty() && s.activeJobs.Len() < int(s.workers) {
		s.workqueue.Enqueued <- struct{}{}
	}
}

func (s *Scheduler[T]) enqueueTask() {
	if item, ok := s.pendingTasks.Pop(); ok {
		s.log.Debug("task moved to workqueue", "task", item)
		s.workqueue.Enqueue(item)
	}
}
