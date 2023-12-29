// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package testutil

import (
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

type ClientOption func(*fake.ClientBuilder)

func WithRuntimeObject(object client.Object) ClientOption {
	return func(b *fake.ClientBuilder) {
		b.WithRuntimeObjects(object)
		b.WithStatusSubresource(object)
	}
}

type SchemeOption func(*runtime.Scheme)

func WithGroupVersion(addToScheme func(*runtime.Scheme) error) SchemeOption {
	return func(s *runtime.Scheme) {
		if err := addToScheme(s); err != nil {
			panic("Fatal")
		}
	}
}

func SetupScheme(opts ...SchemeOption) *runtime.Scheme {
	s := runtime.NewScheme()
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func SetupClient(scheme *runtime.Scheme, opts ...ClientOption) client.Client {
	builder := fake.NewClientBuilder()
	builder.WithScheme(scheme)
	for _, opt := range opts {
		opt(builder)
	}
	cl := builder.Build()
	return cl
}
