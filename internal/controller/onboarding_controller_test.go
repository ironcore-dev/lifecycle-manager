// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"context"
	"testing"
	"time"

	oobv1alpha1 "github.com/onmetal/oob-operator/api/v1alpha1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/ironcore-dev/lifecycle-manager/api/lifecycle/v1alpha1"
)

func TestOnboardingReconciler_Reconcile(t *testing.T) {
	t.Skipf("not all dependencies were migrated to ironcore")
	t.Parallel()

	schemeOpts := []schemeOption{withGroupVersion(oobv1alpha1.AddToScheme)}

	tests := map[string]struct {
		target        *oobv1alpha1.OOB
		request       types.NamespacedName
		expectRequeue bool
		expectError   bool
	}{
		"absent-oob": {
			target: &oobv1alpha1.OOB{
				TypeMeta: metav1.TypeMeta{
					Kind:       "OOB",
					APIVersion: "v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "absent-oob",
					Namespace: "metal",
				},
				Spec: oobv1alpha1.OOBSpec{
					LocatorLED: "On",
					Power:      "On",
					Reset:      "None",
					Filler:     pointer.Int64(0),
				},
				Status: oobv1alpha1.OOBStatus{
					Manufacturer: "Dell",
					SKU:          "R440X0001",
				},
			},
			request: types.NamespacedName{
				Namespace: "default",
				Name:      "sample",
			},
			expectRequeue: false,
			expectError:   false,
		},
		"empty-status": {
			target: &oobv1alpha1.OOB{
				TypeMeta: metav1.TypeMeta{
					Kind:       "OOB",
					APIVersion: "v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "empty-oob",
					Namespace: "default",
				},
				Spec: oobv1alpha1.OOBSpec{
					LocatorLED: "On",
					Power:      "On",
					Reset:      "None",
					Filler:     pointer.Int64(0),
				},
			},
			request: types.NamespacedName{
				Namespace: "default",
				Name:      "empty-oob",
			},
			expectRequeue: true,
			expectError:   false,
		},
		"deleting-oob": {
			target: &oobv1alpha1.OOB{
				TypeMeta: metav1.TypeMeta{
					Kind:       "OOB",
					APIVersion: "v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:              "deleting-oob",
					Namespace:         "default",
					DeletionTimestamp: &metav1.Time{Time: time.Now()},
					Finalizers:        []string{"fake-finalizer"},
				},
				Spec: oobv1alpha1.OOBSpec{
					LocatorLED: "On",
					Power:      "On",
					Reset:      "None",
					Filler:     pointer.Int64(0),
				},
			},
			request: types.NamespacedName{
				Namespace: "default",
				Name:      "deleting-oob",
			},
			expectRequeue: false,
			expectError:   false,
		},
	}

	for n, tc := range tests {
		name, testCase := n, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			clientOpts := []clientOption{withRuntimeObject(testCase.target)}
			r := newOnboardingReconciler(t, schemeOpts, clientOpts)
			resp, err := r.Reconcile(context.Background(), ctrl.Request{NamespacedName: testCase.request})
			if testCase.expectError {
				assert.Error(t, err)
			}
			if !testCase.expectError {
				assert.NoError(t, err)
			}
			if testCase.expectRequeue {
				assert.NotEmpty(t, resp.RequeueAfter)
			}
		})
	}
}

func TestOnboardingReconciler_Onboarding(t *testing.T) {
	t.Skipf("not all dependencies were migrated to ironcore")
	t.Parallel()

	schemeOpts := []schemeOption{
		withGroupVersion(oobv1alpha1.AddToScheme),
		withGroupVersion(v1alpha1.AddToScheme),
	}

	targetOOB := &oobv1alpha1.OOB{
		TypeMeta: metav1.TypeMeta{
			Kind:       "OOB",
			APIVersion: "v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sample-oob",
			Namespace: "metal",
		},
		Spec: oobv1alpha1.OOBSpec{
			LocatorLED: "On",
			Power:      "On",
			Reset:      "None",
			Filler:     pointer.Int64(0),
		},
		Status: oobv1alpha1.OOBStatus{
			Manufacturer: "Dell",
			SKU:          "R440X0001",
		},
	}
	targetMachineType := &v1alpha1.MachineType{
		TypeMeta: metav1.TypeMeta{
			Kind:       "MachineType",
			APIVersion: "v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "Dell-R440",
			Namespace: "metal",
		},
		Spec: v1alpha1.MachineTypeSpec{
			Manufacturer: "Dell",
			Type:         "R440",
			ScanPeriod:   metav1.Duration{Duration: time.Hour * 24},
		},
	}

	clientOpts := []clientOption{withRuntimeObject(targetOOB)}
	r := newOnboardingReconciler(t, schemeOpts, clientOpts)
	_, err := r.Reconcile(context.Background(), ctrl.Request{NamespacedName: client.ObjectKeyFromObject(targetOOB)})
	assert.NoError(t, err)

	actualMachineType := &v1alpha1.MachineType{}
	err = r.Get(context.Background(), types.NamespacedName{Namespace: "metal", Name: "Dell-R440"}, actualMachineType)
	assert.NoError(t, err)
	assert.Equal(t, client.ObjectKeyFromObject(targetMachineType), client.ObjectKeyFromObject(actualMachineType))
	actualMachine := &v1alpha1.Machine{}
	err = r.Get(context.Background(), types.NamespacedName{Namespace: "metal", Name: "sample-oob"}, actualMachine)
	assert.NoError(t, err)
	assert.Equal(t, client.ObjectKeyFromObject(targetOOB), client.ObjectKeyFromObject(actualMachine))
	assert.Equal(t, targetOOB.Name, actualMachine.Spec.OOBMachineRef.Name)
	assert.Equal(t, actualMachineType.Name, actualMachine.Spec.MachineTypeRef.Name)
}
