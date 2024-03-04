// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package testutil

import (
	"time"

	oobv1alpha1 "github.com/ironcore-dev/oob/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type OOBObjectOption func(*oobv1alpha1.OOB)

func OOBWithDeletionTimestamp() OOBObjectOption {
	return func(o *oobv1alpha1.OOB) {
		o.DeletionTimestamp = &metav1.Time{Time: time.Now()}
	}
}

func OOBWithFinalizer() OOBObjectOption {
	return func(o *oobv1alpha1.OOB) {
		o.Finalizers = []string{"test-suite-finalizer"}
	}
}

func OOBWithStatus() OOBObjectOption {
	return func(o *oobv1alpha1.OOB) {
		o.Status.Manufacturer = "Sample"
		o.Status.SKU = "0000X0000"
	}
}

func NewOOBObject(name, namespace string, opts ...OOBObjectOption) *oobv1alpha1.OOB {
	o := &oobv1alpha1.OOB{
		TypeMeta: metav1.TypeMeta{
			Kind:       "OOB",
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
