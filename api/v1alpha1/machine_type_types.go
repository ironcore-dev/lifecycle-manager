// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:generate=true

// MachineTypeSpec contains definition of concrete machine type.
type MachineTypeSpec struct {
	// Manufacturer refers to manufacturer, e.g. Lenovo, Dell etc.
	// +kubebuilder:validation:Required
	Manufacturer string `json:"manufacturer"`

	// Type refers to machine type, e.g. 7z21 for Lenovo, R440 for Dell etc.
	// +kubebuilder:validation:Required
	Type string `json:"type"`

	// ScanPeriod defines the interval between scans.
	// +kubebuilder:validation:Required
	ScanPeriod metav1.Duration `json:"scanPeriod"`
}

// +kubebuilder:object:generate=true

// MachineTypeStatus contains observed state of machine type in part of available firmware updates.
type MachineTypeStatus struct {
	// LastScanTime reflects the timestamp when the last scan of available packages was done.
	// +kubebuilder:validation:Optional
	LastScanTime metav1.Time `json:"lastScanTime"`

	// LastScanResult reflects the result of the last scan.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=Success;Failure
	LastScanResult ScanResult `json:"lastScanResult"`
}

// +kubebuilder:object:root=true

// MachineType is the schema for MachineType API object.
type MachineType struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MachineTypeSpec   `json:"spec,omitempty"`
	Status MachineTypeStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MachineTypeList contains a list of MachineType objects.
type MachineTypeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []MachineType `json:"items"`
}
