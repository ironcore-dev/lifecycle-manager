package controllers

import (
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
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

func withRuntimeObject(object runtime.Object) clientOption {
	return func(b *fake.ClientBuilder) {
		b.WithRuntimeObjects(object)
	}
}

func setupReconciler(t *testing.T, schemeOpts []schemeOption, clientOpts []clientOption) *OnboardingReconciler {
	t.Helper()

	scheme := runtime.NewScheme()
	setupScheme(t, scheme, schemeOpts...)
	c := setupClient(t, scheme, clientOpts...)
	return &OnboardingReconciler{
		Client:        c,
		Scheme:        scheme,
		RequeuePeriod: time.Second * 5,
	}
}
