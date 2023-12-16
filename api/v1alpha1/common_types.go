// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ScanResult string

const (
	ScanFailure ScanResult = "Failure"
	ScanSuccess ScanResult = "Success"
)

func (in ScanResult) IsSuccess() bool {
	return in == ScanSuccess
}

func (in ScanResult) IsFailure() bool {
	return in == ScanFailure
}

type PackageInstallJobState string

const (
	PackageInstallJobStateFailure PackageInstallJobState = "Failure"
	PackageInstallJobStateSuccess PackageInstallJobState = "Success"
)

func (in PackageInstallJobState) IsSuccess() bool {
	return in == PackageInstallJobStateSuccess
}

func (in PackageInstallJobState) IsFailure() bool {
	return in == PackageInstallJobStateFailure
}

type ScanState string

const (
	ScanScheduled ScanState = "Scheduled"
	ScanFinished  ScanState = "Finished"
)

func (in ScanState) IsScheduled() bool {
	return in == ScanScheduled
}

func (in ScanState) IsFinished() bool {
	return in == ScanFinished
}

const (
	ConditionTypePending    = "Pending"
	ConditionTypeScheduled  = "Scheduled"
	ConditionTypeInProgress = "InProgress"
	ConditionTypeFinished   = "Finished"

	ConditionStatusTrue  metav1.ConditionStatus = "True"
	ConditionStatusFalse metav1.ConditionStatus = "False"
)

// +kubebuilder:object:generate=true

type PackageInstallJobCondition struct {
	metav1.Condition `json:",inline"`
}

func (in PackageInstallJobCondition) IsPending() bool {
	return in.Type == ConditionTypePending
}

func (in PackageInstallJobCondition) IsScheduled() bool {
	return in.Type == ConditionTypeScheduled
}

func (in PackageInstallJobCondition) IsInProgress() bool {
	return in.Type == ConditionTypeInProgress
}

func (in PackageInstallJobCondition) IsFinished() bool {
	return in.Type == ConditionTypeFinished
}

func (in PackageInstallJobCondition) IsTrue() bool {
	return in.Status == ConditionStatusTrue
}

func (in PackageInstallJobCondition) IsFalse() bool {
	return in.Status == ConditionStatusFalse
}
