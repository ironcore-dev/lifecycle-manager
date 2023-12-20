// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MachineTypeSpec defines the desired state of MachineType
type MachineTypeSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of MachineType. Edit machinetype_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// MachineTypeStatus defines the observed state of MachineType
type MachineTypeStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// MachineType is the Schema for the machinetypes API
type MachineType struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MachineTypeSpec   `json:"spec,omitempty"`
	Status MachineTypeStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MachineTypeList contains a list of MachineType
type MachineTypeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MachineType `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MachineType{}, &MachineTypeList{})
}
