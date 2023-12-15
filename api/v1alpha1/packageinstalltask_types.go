// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:generate=true

// PackageInstallTaskSpec contains desired configuration of update approval.
type PackageInstallTaskSpec struct {
	// MachineTypeRef contain reference to MachineType object.
	// +kubebuilder:validation:Required
	MachineTypeRef corev1.LocalObjectReference `json:"machineTypeRef"`

	// Packages contains a list of references to FirmwarePackage objects
	// approved for installation.
	// +kubebuilder:validation:Required
	Packages []corev1.LocalObjectReference `json:"packages"`

	// Machines contains a list of references to Machine objects
	// approved for packages upgrade.
	// +kubebuilder:validation:Required
	Machines []corev1.LocalObjectReference `json:"machines"`
}

// +kubebuilder:object:generate=true

// PackageInstallTaskStatus contains observed state of PackageInstallTask object.
type PackageInstallTaskStatus struct {
	// JobsTotal reflects the total number of PackageInstallJob objects created.
	// according to the PackageInstallTask definition
	// +kubebuilder:validation:Optional
	JobsTotal uint32 `json:"jobsTotal"`

	// JobsFailed reflects the number of failed jobs.
	// +kubebuilder:validation:Optional
	JobsFailed uint32 `json:"jobsFailed"`

	// JobsSuccessful reflects the number of successfully completed jobs.
	// +kubebuilder:validation:Optional
	JobsSuccessful uint32 `json:"jobsSuccessful"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=pit
// +kubebuilder:storageversion
// +kubebuilder:printcolumn:name="Jobs Total",type=integer,JSONPath=`.status.jobsTotal`,description="Total number of jobs created from current task"
// +kubebuilder:printcolumn:name="Jobs Failed",type=integer,JSONPath=`.status.jobsFailed`,description="Total number of failed jobs"
// +kubebuilder:printcolumn:name="Jobs Successful",type=integer,JSONPath=`.status.jobsSuccessful`,description="Total number of succeeded jobs"

// PackageInstallTask is the schema for PackageInstallTask API object.
type PackageInstallTask struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PackageInstallTaskSpec   `json:"spec,omitempty"`
	Status PackageInstallTaskStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// PackageInstallTaskList contains a list of PackageInstallTask objects.
type PackageInstallTaskList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PackageInstallTask `json:"items"`
}
