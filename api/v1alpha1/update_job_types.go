// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:generate=true

// UpdateJobSpec contains desired state of the update task.
type UpdateJobSpec struct {
	// MachineLifecycleRef contains reference to MachineLifecycle object which
	// updates will be installed to.
	// +kubebuilder:validation:Required
	MachineLifecycleRef corev1.LocalObjectReference `json:"machineLifecycleRef"`

	// FirmwarePackageRef contains reference to FirmwarePackage object which should be installed.
	// +kubebuilder:validation:Required
	FirmwarePackageRef corev1.LocalObjectReference `json:"firmwarePackageRef"`
}

// +kubebuilder:object:generate=true

// UpdateJobStatus contains observed state of the UpdateJob object.
type UpdateJobStatus struct {
	// Conditions reflects the object state change flow.
	// Possible condition types:
	// - Pending
	// - Scheduled
	// - InProgress
	// - Finished
	// +kubebuilder:validation:Optional
	Conditions []metav1.Condition `json:"conditions"`

	// State reflects the final state of the job.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=Success;Failure
	State UpdateJobState `json:"state"`

	// Message contains verbose message related to the current State.
	// +kubebuilder:validation:Optional
	Message string `json:"message"`
}

// +kubebuilder:object:root=true

// UpdateJob is the schema for UpdateJob API object.
type UpdateJob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   UpdateJobSpec   `json:"spec,omitempty"`
	Status UpdateJobStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// UpdateJobList contains a list of UpdateJob objects.
type UpdateJobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []UpdateJob `json:"items"`
}
