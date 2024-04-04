// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package scheduler

import (
	"context"

	v1 "k8s.io/api/batch/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (s *Scheduler[T]) processJob(ctx context.Context, task Task[T]) error {
	// 1. get kubernetes client
	// 2.
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
	s.log.Info("new job initiated", "job", task)
	return nil
}
