// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package mock

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	lifecyclev1alpha1 "github.com/ironcore-dev/lifecycle-manager/api/lifecycle/v1alpha1"
)

type MachineTypeOption func(*lifecyclev1alpha1.MachineType)

func MachineTypeWithDeletionTimestamp() MachineTypeOption {
	return func(o *lifecyclev1alpha1.MachineType) {
		o.DeletionTimestamp = &metav1.Time{Time: time.Now()}
	}
}

func MachineTypeWithFinalizer() MachineTypeOption {
	return func(o *lifecyclev1alpha1.MachineType) {
		o.Finalizers = []string{"test-suite-finalizer"}
	}
}

func MachineTypeWithGroup(group lifecyclev1alpha1.MachineGroup) MachineTypeOption {
	return func(o *lifecyclev1alpha1.MachineType) {
		o.Spec.MachineGroups = append(o.Spec.MachineGroups, group)
	}
}

func NewMachineTypeObject(name, namespace string, opts ...MachineTypeOption) *lifecyclev1alpha1.MachineType {
	o := &lifecyclev1alpha1.MachineType{
		TypeMeta: metav1.TypeMeta{
			Kind:       "MachineType",
			APIVersion: "v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}
