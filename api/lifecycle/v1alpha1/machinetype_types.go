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
	Manufacturer string `json:"manufacturer" protobuf:"bytes,1,name=manufacturer"`

	// Type refers to machine type, e.g. 7z21 for Lenovo, R440 for Dell etc.
	// +kubebuilder:validation:Required
	Type string `json:"type" protobuf:"bytes,2,name=type"`

	// ScanPeriod defines the interval between scans.
	// +kubebuilder:validation:Required
	ScanPeriod metav1.Duration `json:"scanPeriod" protobuf:"bytes,3,name=scan_period"`

	// MachineGroups defines list of MachineGroup
	// +kubebuilder:validation:Optional
	MachineGroups []MachineGroup `json:"machineGroups" protobuf:"bytes,4,rep,name=machine_groups"`
}

// MachineGroup defines group of Machine objects filtered by label selector
// and a list of firmware packages versions which should be installed by default.
type MachineGroup struct {
	// MachineSelector defines native kubernetes label selector to apply to Machine objects.
	// +kubebuilder:validation:Required
	MachineSelector metav1.LabelSelector `json:"machineSelector" protobuf:"bytes,1,name=machine_selector"`

	// Packages defines default firmware package versions for the group of Machine objects.
	// +kubebuilder:validation:Required
	Packages []PackageVersion `json:"packages" protobuf:"bytes,2,rep,name=packages"`
}

// MachineTypeStatus defines the observed state of MachineType.
type MachineTypeStatus struct {
	// LastScanTime reflects the timestamp when the last scan of available packages was done.
	// +kubebuilder:validation:Optional
	LastScanTime metav1.Time `json:"lastScanTime" protobuf:"bytes,1,opt,name=last_scan_time"`

	// LastScanResult reflects the result of the last scan.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=Success;Failure
	LastScanResult ScanResult `json:"lastScanResult" protobuf:"bytes,2,opt,name=last_scan_result,casttype=ScanResult"`

	// AvailablePackages reflects the list of AvailablePackageVersion
	// +kubebuilder:validation:Optional
	AvailablePackages []AvailablePackageVersions `json:"availablePackages" protobuf:"bytes,3,rep,name=available_packages"`

	// Message contains verbose message explaining current state
	// +kubebuilder:validation:Optional
	Message string `json:"message" protobuf:"bytes,4,opt,name=message"`
}

// AvailablePackageVersions defines a number of versions for concrete firmware package.
type AvailablePackageVersions struct {
	// Name reflects the name of the firmware package
	// +kubebuilder:validation:Required
	Name string `json:"name" protobuf:"bytes,1,name=name"`

	// Versions reflects the list of discovered package versions available for installation.
	// +kubebuilder:validation:Required
	Versions []string `json:"versions" protobuf:"bytes,2,rep,name=versions"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient

// MachineType is the Schema for the machinetypes API.
type MachineType struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   MachineTypeSpec   `json:"spec,omitempty" protobuf:"bytes,2,name=spec"`
	Status MachineTypeStatus `json:"status,omitempty" protobuf:"bytes,3,name=status"`
}

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MachineTypeList contains a list of MachineType.
type MachineTypeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []MachineType `json:"items" protobuf:"bytes,2,name=items"`
}
