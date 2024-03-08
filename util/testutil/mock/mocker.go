// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package mock

import (
	lifecyclev1alpha1 "github.com/ironcore-dev/lifecycle-manager/api/lifecycle/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

type UnstructuredBuilder struct {
	inner unstructured.Unstructured
}

func NewUnstructuredBuilder() *UnstructuredBuilder {
	return &UnstructuredBuilder{inner: unstructured.Unstructured{}}
}

func (b *UnstructuredBuilder) WithName(name string) *UnstructuredBuilder {
	b.inner.SetName(name)
	return b
}

func (b *UnstructuredBuilder) WithNamespace(namespace string) *UnstructuredBuilder {
	b.inner.SetNamespace(namespace)
	return b
}

func (b *UnstructuredBuilder) WithDeletionTimestamp(timestamp *metav1.Time) *UnstructuredBuilder {
	b.inner.SetDeletionTimestamp(timestamp)
	return b
}

func (b *UnstructuredBuilder) WithLabels(labels map[string]string) *UnstructuredBuilder {
	b.inner.SetLabels(labels)
	return b
}

func (b *UnstructuredBuilder) WithFinalizers(finalizers []string) *UnstructuredBuilder {
	b.inner.SetFinalizers(finalizers)
	return b
}

func (b *UnstructuredBuilder) ToMachine() *MachineMockBuilder {
	var m lifecyclev1alpha1.Machine
	b.inner.SetAPIVersion("lifecycle.ironcore.dev/v1alpha1")
	b.inner.SetKind("Machine")
	_ = runtime.DefaultUnstructuredConverter.FromUnstructured(b.inner.Object, &m)
	return &MachineMockBuilder{inner: &m}
}
