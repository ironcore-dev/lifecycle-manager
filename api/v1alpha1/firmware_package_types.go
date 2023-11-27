// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:generate=true

// FirmwarePackageSpec contains definition of concrete firmware package specific
// for manufacturer and machine type.
type FirmwarePackageSpec struct {
	// MachineTypeRef contain reference to MachineType object.
	// +kubebuilder:validation:Required
	MachineTypeRef corev1.LocalObjectReference `json:"machineTypeRef"`

	// Name contains the name of the package.
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// Version contains the version of the package.
	// +kubebuilder:validation:Required
	Version string `json:"version"`

	// Source contains the URL where the package can be downloaded from.
	// +kubebuilder:validation:Optional
	Source string `json:"source"`

	// RebootRequired reflects whether the machine reboot required during package installation.
	// +kubebuilder:validation:Required
	RebootRequired bool `json:"rebootRequired"`
}

// +kubebuilder:object:generate=true

// FirmwarePackageStatus contains observed state of FirmwarePackage object.
type FirmwarePackageStatus struct{}

// +kubebuilder:object:root=true

// FirmwarePackage is the schema for FirmwarePackage API object.
type FirmwarePackage struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FirmwarePackageSpec   `json:"spec,omitempty"`
	Status FirmwarePackageStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// FirmwarePackageList contains a list of FirmwarePackage objects.
type FirmwarePackageList struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Items []FirmwarePackage `json:"items"`
}
