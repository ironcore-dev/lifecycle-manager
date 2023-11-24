package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:generate=true

// MachineTypeSpec contains definition of concrete machine type.
type MachineTypeSpec struct {
}

// MachineTypeStatus contains observed state of machine type in part of available firmware updates.
type MachineTypeStatus struct {
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
