// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

type schemeOption func(*testing.T, *runtime.Scheme)
type clientOption func(*fake.ClientBuilder)

func setupScheme(t *testing.T, scheme *runtime.Scheme, opts ...schemeOption) {
	t.Helper()

	for _, opt := range opts {
		opt(t, scheme)
	}
}

func withGroupVersion(addToScheme func(*runtime.Scheme) error) schemeOption {
	return func(t *testing.T, s *runtime.Scheme) {
		if err := addToScheme(s); err != nil {
			t.Fatal(err)
		}
	}
}

func setupClient(t *testing.T, scheme *runtime.Scheme, opts ...clientOption) client.Client {
	t.Helper()

	builder := fake.NewClientBuilder()
	builder.WithScheme(scheme)
	for _, opt := range opts {
		opt(builder)
	}
	cl := builder.Build()
	return cl
}

func withRuntimeObject(object client.Object) clientOption {
	return func(b *fake.ClientBuilder) {
		b.WithRuntimeObjects(object)
		b.WithStatusSubresource(object)
	}
}

func setupPrerequisites(t *testing.T, sOpts []schemeOption, cOpts []clientOption) (client.Client, *runtime.Scheme) {
	t.Helper()

	s := runtime.NewScheme()
	setupScheme(t, s, sOpts...)
	c := setupClient(t, s, cOpts...)
	return c, s
}

func newOnboardingReconciler(
	t *testing.T,
	sOpts []schemeOption,
	cOpts []clientOption,
) *OnboardingReconciler {
	t.Helper()

	c, s := setupPrerequisites(t, sOpts, cOpts)
	return &OnboardingReconciler{
		Client:        c,
		Scheme:        s,
		RequeuePeriod: time.Second * 5,
	}
}

func newMachineTypeReconciler(
	t *testing.T,
	sOpts []schemeOption,
	cOpts []clientOption,
) *MachineTypeReconciler {
	t.Helper()

	c, s := setupPrerequisites(t, sOpts, cOpts)
	return &MachineTypeReconciler{
		Client: c,
		Scheme: s,
	}
}

func newMachineReconciler(
	t *testing.T,
	sOpts []schemeOption,
	cOpts []clientOption,
) *MachineReconciler {
	t.Helper()

	c, s := setupPrerequisites(t, sOpts, cOpts)
	return &MachineReconciler{
		Client: c,
		Scheme: s,
	}
}
