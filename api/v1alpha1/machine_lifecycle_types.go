package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:generate=true

// MachineLifecycleSpec contains desired configuration of machine lifecycle.
type MachineLifecycleSpec struct {
}

// MachineLifecycleStatus contains observed state of MachineLifecycle object.
type MachineLifecycleStatus struct {
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
