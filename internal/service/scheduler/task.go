// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package scheduler

const (
	ScanJob    JobType = "scan"
	InstallJob JobType = "install"
)

type Task[T LifecycleObject] struct {
	Key        string
	Type       JobType
	Target     T
	TargetType string
}

func NewTask[T LifecycleObject](key string, taskType JobType, target T, targetType string) Task[T] {
	return Task[T]{
		Key:        key,
		Type:       taskType,
		Target:     target,
		TargetType: targetType,
	}
}
