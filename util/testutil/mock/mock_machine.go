// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package mock

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	lifecyclev1alpha1 "github.com/ironcore-dev/lifecycle-manager/api/lifecycle/v1alpha1"
)

type MachineMockBuilder struct {
	inner *lifecyclev1alpha1.Machine
}

func (b *MachineMockBuilder) WithMachineTypeRef(name string) *MachineMockBuilder {
	b.inner.Spec.MachineTypeRef.Name = name
	return b
}

func (b *MachineMockBuilder) WithOOBMachineRef(name string) *MachineMockBuilder {
	b.inner.Spec.OOBMachineRef.Name = name
	return b
}

func (b *MachineMockBuilder) WithInstalledPackages(pkg []lifecyclev1alpha1.PackageVersion) *MachineMockBuilder {
	b.inner.Status.InstalledPackages = pkg
	return b
}

func (b *MachineMockBuilder) Complete() *lifecyclev1alpha1.Machine {
	return b.inner
}

type MachineOption func(*lifecyclev1alpha1.Machine)

func MachineWithLabels(lbl map[string]string) MachineOption {
	return func(o *lifecyclev1alpha1.Machine) {
		o.Labels = lbl
	}
}

func MachineWithMachineTypeRef(name string) MachineOption {
	return func(o *lifecyclev1alpha1.Machine) {
		o.Spec.MachineTypeRef = corev1.LocalObjectReference{Name: name}
	}
}

func MachineStatusWithInstalledPackages(pkg []lifecyclev1alpha1.PackageVersion) MachineOption {
	return func(o *lifecyclev1alpha1.Machine) {
		o.Status.InstalledPackages = pkg
	}
}

func NewMachineObject(name, namespace string, opts ...MachineOption) *lifecyclev1alpha1.Machine {
	o := &lifecyclev1alpha1.Machine{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Machine",
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
