package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:generate=true

type PackageUpdateSpec struct {
}

type PackageUpdateStatus struct {
}

// +kubebuilder:object:root=true

type PackageUpdate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PackageUpdateSpec   `json:"spec,omitempty"`
	Status PackageUpdateStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

type PackageUpdateList struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Items []PackageUpdate `json:"items"`
}
