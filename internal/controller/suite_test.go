// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/runtime"

	// +kubebuilder:scaffold:imports

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

func TestControllers(t *testing.T) {
	t.Parallel()
	RegisterFailHandler(Fail)
	RunSpecs(t, "Controller Suite")
}

func NewOnboardingReconciler(c client.Client, s *runtime.Scheme) *OnboardingReconciler {
	return &OnboardingReconciler{
		Client:        c,
		Scheme:        s,
		RequeuePeriod: time.Second * 5,
	}
}

func NewMachineTypeReconciler(c client.Client, s *runtime.Scheme) *MachineTypeReconciler {
	return &MachineTypeReconciler{
		Client: c,
		Scheme: s,
	}
}

func NewMachineReconciler(c client.Client, s *runtime.Scheme) *MachineReconciler {
	return &MachineReconciler{
		Client: c,
		Scheme: s,
	}
}
