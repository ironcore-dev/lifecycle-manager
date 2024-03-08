// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package mock

import (
	oobv1alpha1 "github.com/ironcore-dev/oob/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
)

type OOBMockBuilder struct {
	inner *oobv1alpha1.OOB
}

func (b *UnstructuredBuilder) OOBFromUnstructured() *OOBMockBuilder {
	var o oobv1alpha1.OOB
	b.inner.SetAPIVersion("ironcore.dev/v1alpha1")
	b.inner.SetKind("OOB")
	_ = runtime.DefaultUnstructuredConverter.FromUnstructured(b.inner.Object, &o)
	return &OOBMockBuilder{inner: &o}
}

func (b *OOBMockBuilder) WithManufacturer(manufacturer string) *OOBMockBuilder {
	b.inner.Status.Manufacturer = manufacturer
	return b
}

func (b *OOBMockBuilder) WithSKU(sku string) *OOBMockBuilder {
	b.inner.Status.SKU = sku
	return b
}

func (b *OOBMockBuilder) Complete() *oobv1alpha1.OOB {
	return b.inner
}
