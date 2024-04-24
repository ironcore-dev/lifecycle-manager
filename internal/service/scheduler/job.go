// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package scheduler

import (
	"context"
	"fmt"

	v1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

const (
	lifecycleJobIdLabel   = "lifecycle.ironcore.dev/job-id"
	lifecycleJobTypeLabel = "lifecycle.ironcore.dev/job-type"
)

func (s *Scheduler[T]) processJob(ctx context.Context, task Task[T]) error {
	job, err := s.BatchV1().Jobs(s.namespace).Get(ctx, task.Key, metav1.GetOptions{})
	if err == nil {
		// check job state and act accordingly
		return s.processExistingJob(ctx, task, job)
	}
	if apierrors.IsNotFound(err) {
		// create job
		return s.createJob(ctx, task)
	}
	return err
}

func (s *Scheduler[T]) processExistingJob(ctx context.Context, task Task[T], job *v1.Job) error {
	s.log.Info("new job initiated", "job", task)
	return nil
}

func (s *Scheduler[T]) createJob(ctx context.Context, task Task[T]) error {
	namespace := task.Target.GetNamespace()
	job := &v1.Job{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: fmt.Sprintf("%s-", task.Key),
			Namespace:    namespace,
			Labels: map[string]string{
				lifecycleJobIdLabel:   task.Key,
				lifecycleJobTypeLabel: string(task.Type),
			},
		},
		Spec: v1.JobSpec{
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						lifecycleJobIdLabel:   task.Key,
						lifecycleJobTypeLabel: string(task.Type),
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "lifecycle-job",
							Image: "",
							Args:  []string{"--job-id", task.Key},
						},
					},
				},
			},
			TTLSecondsAfterFinished: ptr.To(int32(30)),
		},
	}
	if _, err := s.Clientset.BatchV1().Jobs(namespace).Create(ctx, job, metav1.CreateOptions{}); err != nil {
		return err
	}
	s.log.Info("new job initiated", "job", task)
	return nil
}
