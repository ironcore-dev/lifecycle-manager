// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:generate=true

// MachineLifecycleSpec contains desired configuration of machine lifecycle.
type MachineLifecycleSpec struct {
	// MachineTypeRef contain reference to MachineType object.
	// +kubebuilder:validation:Required
	MachineTypeRef corev1.LocalObjectReference `json:"machineTypeRef"`

	// OOBMachineRef contains reference to OOB machine object.
	// +kubebuilder:validation:Required
	OOBMachineRef corev1.LocalObjectReference `json:"oobMachineRef"`

	// ScanPeriod defines the interval between scans.
	// +kubebuilder:validation:Required
	ScanPeriod metav1.Duration `json:"scanPeriod"`
}

// +kubebuilder:object:generate=true

// MachineLifecycleStatus contains observed state of MachineLifecycle object.
type MachineLifecycleStatus struct {
	// LastScanTime reflects the timestamp when the last scan of available packages was done.
	// +kubebuilder:validation:Optional
	LastScanTime metav1.Time `json:"lastScanTime"`

	// LastScanResult reflects the state of the last scan.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=Success;Failure
	LastScanResult ScanResult `json:"lastScanResult"`

	// InstalledPackages contains the list of references to FirmwarePackage objects.
	// +kubebuilder:validation:Optional
	InstalledPackages []corev1.LocalObjectReference `json:"installedPackages"`
}

// +kubebuilder:object:root=true

// MachineLifecycle is the schema for MachineLifecycle API object.
type MachineLifecycle struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MachineLifecycleSpec   `json:"spec,omitempty"`
	Status MachineLifecycleStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MachineLifecycleList contains a list of MachineLifecycle objects.
type MachineLifecycleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []MachineLifecycle `json:"items"`
}
