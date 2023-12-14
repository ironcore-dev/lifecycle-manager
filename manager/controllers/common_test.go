// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

type schemeOption func(*testing.T, *runtime.Scheme)
type clientOption func(*fake.ClientBuilder)

func setupScheme(t *testing.T, scheme *runtime.Scheme, opts ...schemeOption) {
	t.Helper()

	for _, opt := range opts {
		opt(t, scheme)
	}
}

func withGroupVersion(b *scheme.Builder) schemeOption {
	return func(t *testing.T, s *runtime.Scheme) {
		if err := b.AddToScheme(s); err != nil {
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

func setupPrerequisites(t *testing.T, schemeOpts []schemeOption, clientOpts []clientOption) (client.Client, *runtime.Scheme) {
	t.Helper()

	s := runtime.NewScheme()
	setupScheme(t, s, schemeOpts...)
	c := setupClient(t, s, clientOpts...)
	return c, s
}

func newOnboardingReconciler(t *testing.T, schemeOpts []schemeOption, clientOpts []clientOption) *OnboardingReconciler {
	c, s := setupPrerequisites(t, schemeOpts, clientOpts)
	return &OnboardingReconciler{
		Client:        c,
		Scheme:        s,
		RequeuePeriod: time.Second * 5,
	}
}

func newUpdateTaskReconciler(t *testing.T, schemeOpts []schemeOption, clientOpts []clientOption) *UpdateTaskReconciler {
	c, s := setupPrerequisites(t, schemeOpts, clientOpts)
	return &UpdateTaskReconciler{
		Client:   c,
		Scheme:   s,
		Recorder: &record.FakeRecorder{},
	}
}

func newMachineTypeReconciler(t *testing.T, schemeOpts []schemeOption, clientOpts []clientOption) *MachineTypeReconciler {
	c, s := setupPrerequisites(t, schemeOpts, clientOpts)
	return &MachineTypeReconciler{
		Client:  c,
		Scheme:  s,
		Horizon: time.Minute * 30,
	}
}
