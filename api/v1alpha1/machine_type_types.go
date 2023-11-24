// Copyright 2023 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
