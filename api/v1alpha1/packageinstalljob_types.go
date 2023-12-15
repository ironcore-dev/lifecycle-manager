// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:generate=true

// PackageInstallJobSpec contains desired state of the update task.
type PackageInstallJobSpec struct {
	// MachineRef contains reference to Machine object which
	// updates will be installed to.
	// +kubebuilder:validation:Required
	MachineRef corev1.LocalObjectReference `json:"machineRef"`

	// FirmwarePackageRef contains reference to FirmwarePackage object which should be installed.
	// +kubebuilder:validation:Required
	FirmwarePackageRef corev1.LocalObjectReference `json:"firmwarePackageRef"`
}

// +kubebuilder:object:generate=true

// PackageInstallJobStatus contains observed state of the PackageInstallJob object.
type PackageInstallJobStatus struct {
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
	State PackageInstallJobState `json:"state"`

	// Message contains verbose message related to the current State.
	// +kubebuilder:validation:Optional
	Message string `json:"message"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=pij
// +kubebuilder:storageversion
// +kubebuilder:printcolumn:name="State",type=string,JSONPath=`.status.state`,description="Job state"
// +kubebuilder:printcolumn:name="Message",priority=1,type=string,JSONPath=`.status.message`,description="Job state message reports about any issues during processing"

// PackageInstallJob is the schema for PackageInstallJob API object.
type PackageInstallJob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PackageInstallJobSpec   `json:"spec,omitempty"`
	Status PackageInstallJobStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// PackageInstallJobList contains a list of PackageInstallJob objects.
type PackageInstallJobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PackageInstallJob `json:"items"`
}
