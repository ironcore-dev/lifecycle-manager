// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MachineSpec defines the desired state of Machine.
type MachineSpec struct {
	// MachineTypeRef contain reference to MachineType object.
	// +kubebuilder:validation:Required
	MachineTypeRef corev1.LocalObjectReference `json:"machineTypeRef" protobuf:"bytes,1,name=machine_type_ref"`

	// OOBMachineRef contains reference to OOB machine object.
	// +kubebuilder:validation:Required
	OOBMachineRef corev1.LocalObjectReference `json:"oobMachineRef" protobuf:"bytes,2,name=oob_machine_ref"`

	// ScanPeriod defines the interval between scans.
	// +kubebuilder:validation:Required
	ScanPeriod metav1.Duration `json:"scanPeriod" protobuf:"bytes,3,name=scan_period"`

	// Packages defines the list of package versions to install.
	// +kubebuilder:validation:Optional
	Packages []PackageVersion `json:"packages" protobuf:"bytes,4,rep,name=packages"`
}

// MachineStatus defines the observed state of Machine.
type MachineStatus struct {
	// LastScanTime reflects the timestamp when the last scan job for installed
	// firmware versions was performed.
	// +kubebuilder:validation:Optional
	LastScanTime metav1.Time `json:"lastScanTime" protobuf:"bytes,1,opt,name=last_scan_time"`

	// LastScanResult reflects either success or failure of the last scan job.
	// +kubebuilder:validation:Optional
	LastScanResult ScanResult `json:"lastScanResult" protobuf:"bytes,2,opt,name=last_scan_result,casttype=ScanResult"`

	// InstalledPackages reflects the versions of installed firmware packages.
	// +kubebuilder:validation:Optional
	InstalledPackages []PackageVersion `json:"installedPackages" protobuf:"bytes,3,rep,name=installed_packages"`

	// Message contains verbose message explaining current state
	// +kubebuilder:validation:Optional
	Message string `json:"message" protobuf:"bytes,4,opt,name=message"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient

// Machine is the Schema for the machines API.
type Machine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   MachineSpec   `json:"spec,omitempty" protobuf:"bytes,2,name=spec"`
	Status MachineStatus `json:"status,omitempty" protobuf:"bytes,3,name=status"`
}

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MachineList contains a list of Machine.
type MachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []Machine `json:"items" protobuf:"bytes,2,name=items"`
}
