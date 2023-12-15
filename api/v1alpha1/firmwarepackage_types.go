// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:generate=true

// FirmwarePackageSpec defines the desired state of FirmwarePackage.
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

// FirmwarePackageStatus defines the observed state of FirmwarePackage.
type FirmwarePackageStatus struct{}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=fi
// +kubebuilder:storageversion
// +kubebuilder:printcolumn:name="Package Name",type=string,JSONPath=`.spec.name`,description="Firmware package name"
// +kubebuilder:printcolumn:name="Version",type=string,JSONPath=`.spec.version`,description="Firmware package version"
// +kubebuilder:printcolumn:name="Reboot required",type=boolean,JSONPath=`.spec.rebootRequired`,description="Reflects whether reboot required"

// FirmwarePackage is the Schema for the firmwarepackages API.
type FirmwarePackage struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FirmwarePackageSpec   `json:"spec,omitempty"`
	Status FirmwarePackageStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// FirmwarePackageList contains a list of FirmwarePackage.
type FirmwarePackageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FirmwarePackage `json:"items"`
}
