// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:generate=true

// UpdateTaskSpec contains desired configuration of update approval.
type UpdateTaskSpec struct {
	// MachineTypeRef contain reference to MachineType object.
	// +kubebuilder:validation:Required
	MachineTypeRef corev1.LocalObjectReference `json:"machineTypeRef"`

	// Packages contains a list of references to FirmwarePackage objects
	// approved for installation.
	// +kubebuilder:validation:Required
	Packages []corev1.LocalObjectReference `json:"packages"`

	// Machines contains a list of references to MachineLifecycle objects
	// approved for packages upgrade.
	// +kubebuilder:validation:Required
	Machines []corev1.LocalObjectReference `json:"machines"`
}

// +kubebuilder:object:generate=true

// UpdateTaskStatus contains observed state of UpdateTask object.
type UpdateTaskStatus struct {
	// JobsTotal reflects the total number of UpdateJob objects created.
	// according to the UpdateTask definition
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

// UpdateTask is the schema for UpdateTask API object.
type UpdateTask struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   UpdateTaskSpec   `json:"spec,omitempty"`
	Status UpdateTaskStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// UpdateTaskList contains a list of UpdateTask objects.
type UpdateTaskList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []UpdateTask `json:"items"`
}
