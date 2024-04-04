// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package mock

import (
	"k8s.io/apimachinery/pkg/runtime"

	lifecyclev1alpha1 "github.com/ironcore-dev/lifecycle-manager/api/lifecycle/v1alpha1"
)

type MachineTypeMockBuilder struct {
	inner *lifecyclev1alpha1.MachineType
}

func (b *UnstructuredBuilder) MachineTypeFromUnstructured() *MachineTypeMockBuilder {
	var m lifecyclev1alpha1.MachineType
	b.inner.SetAPIVersion("lifecycle.ironcore.dev/v1alpha1")
	b.inner.SetKind("MachineType")
	_ = runtime.DefaultUnstructuredConverter.FromUnstructured(b.inner.Object, &m)
	return &MachineTypeMockBuilder{inner: &m}
}

func (b *MachineTypeMockBuilder) WithMachineGroups(groups []lifecyclev1alpha1.MachineGroup) *MachineTypeMockBuilder {
	b.inner.Spec.MachineGroups = groups
	return b
}

func (b *MachineTypeMockBuilder) Complete() *lifecyclev1alpha1.MachineType {
	return b.inner
}
