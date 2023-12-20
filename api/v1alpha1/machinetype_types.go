// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MachineTypeSpec defines the desired state of MachineType.
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

	// MachineGroups defines list of MachineGroup
	// +kubebuilder:validation:Optional
	MachineGroups []MachineGroup `json:"machineGroups"`
}

// MachineGroup defines group of Machine objects filtered by label selector
// and a list of firmware packages versions which should be installed by default.
type MachineGroup struct {
	// MachineSelector defines native kubernetes label selector to apply to Machine objects.
	// +kubebuilder:validation:Required
	MachineSelector metav1.LabelSelector `json:"machineSelector"`

	// Packages defines default firmware package versions for the group of Machine objects.
	// +kubebuilder:validation:Required
	Packages []PackageVersion `json:"packages"`
}

// MachineTypeStatus defines the observed state of MachineType.
type MachineTypeStatus struct {
	// LastScanTime reflects the timestamp when the last scan of available packages was done.
	// +kubebuilder:validation:Optional
	LastScanTime metav1.Time `json:"lastScanTime"`

	// LastScanResult reflects the result of the last scan.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=Success;Failure
	LastScanResult ScanResult `json:"lastScanResult"`

	// AvailablePackages reflects the list of AvailablePackageVersion
	// +kubebuilder:validation:Optional
	AvailablePackages []AvailablePackageVersions `json:"availablePackages"`
}

// AvailablePackageVersions defines a number of versions for concrete firmware package.
type AvailablePackageVersions struct {
	// Name reflects the name of the firmware package
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// Versions reflects the list of discovered package versions available for installation.
	// +kubebuilder:validation:Required
	Versions []string `json:"versions"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// MachineType is the Schema for the machinetypes API.
type MachineType struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MachineTypeSpec   `json:"spec,omitempty"`
	Status MachineTypeStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MachineTypeList contains a list of MachineType.
type MachineTypeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MachineType `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MachineType{}, &MachineTypeList{})
}
