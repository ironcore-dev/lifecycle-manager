// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/ironcore-dev/lifecycle-manager/api/v1alpha1"
	lcmicommon "github.com/ironcore-dev/lifecycle-manager/lcmi/api/common/v1alpha1"
	lcmimachinetype "github.com/ironcore-dev/lifecycle-manager/lcmi/api/machine_type/v1alpha1"
	"github.com/ironcore-dev/lifecycle-manager/lcmi/fake"
	"github.com/ironcore-dev/lifecycle-manager/util/uuidutil"
)

func TestMachineTypeReconciler_Reconcile(t *testing.T) {
	t.Parallel()

	schemeOpts := []schemeOption{withGroupVersion(v1alpha1.SchemeBuilder)}
	expectedLastScanTime := time.Now()

	tests := map[string]struct {
		target       *v1alpha1.MachineType
		brokerClient *fake.MachineTypeClient
		request      ctrl.Request
		expectError  bool
	}{
		"absent-machine-type": {
			target: &v1alpha1.MachineType{
				TypeMeta: metav1.TypeMeta{
					Kind:       "MachineType",
					APIVersion: "v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sample-machine-type",
					Namespace: "metal",
				},
				Spec: v1alpha1.MachineTypeSpec{
					Manufacturer: "Dell",
					Type:         "R440",
					ScanPeriod:   metav1.Duration{Duration: time.Hour * 24},
				},
				Status: v1alpha1.MachineTypeStatus{},
			},
			request:     ctrl.Request{NamespacedName: types.NamespacedName{Name: "absent-machine-type", Namespace: "metal"}},
			expectError: false,
		},
		"deleting-machine-type": {
			target: &v1alpha1.MachineType{
				TypeMeta: metav1.TypeMeta{
					Kind:       "MachineType",
					APIVersion: "v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:              "deleting-machine-type",
					Namespace:         "metal",
					DeletionTimestamp: &metav1.Time{Time: time.Now()},
					Finalizers:        []string{"fake-finalizer"},
				},
				Spec: v1alpha1.MachineTypeSpec{
					Manufacturer: "Dell",
					Type:         "R440",
					ScanPeriod:   metav1.Duration{Duration: time.Hour * 24},
				},
				Status: v1alpha1.MachineTypeStatus{},
			},
			request:     ctrl.Request{NamespacedName: types.NamespacedName{Name: "deleting-machine-type", Namespace: "metal"}},
			expectError: false,
		},
		"no-scan-required": {
			target: &v1alpha1.MachineType{
				TypeMeta: metav1.TypeMeta{
					Kind:       "MachineType",
					APIVersion: "v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sample-machine-type",
					Namespace: "metal",
				},
				Spec: v1alpha1.MachineTypeSpec{
					Manufacturer: "Dell",
					Type:         "R440",
					ScanPeriod:   metav1.Duration{Duration: time.Hour * 24},
				},
				Status: v1alpha1.MachineTypeStatus{
					LastScanTime:   metav1.Time{Time: time.Now().Add(-10 * time.Minute)},
					LastScanResult: v1alpha1.ScanSuccess,
				},
			},
			request:     ctrl.Request{NamespacedName: types.NamespacedName{Name: "sample-machine-type", Namespace: "metal"}},
			expectError: false,
		},
		"scan-required-not-found": {
			target: &v1alpha1.MachineType{
				TypeMeta: metav1.TypeMeta{
					Kind:       "MachineType",
					APIVersion: "v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sample-machine-type",
					Namespace: "metal",
				},
				Spec: v1alpha1.MachineTypeSpec{
					Manufacturer: "Dell",
					Type:         "R440",
					ScanPeriod:   metav1.Duration{Duration: time.Hour * 24},
				},
				Status: v1alpha1.MachineTypeStatus{},
			},
			brokerClient: fake.NewFakeMachineTypeClient(),
			request:      ctrl.Request{NamespacedName: types.NamespacedName{Name: "sample-machine-type", Namespace: "metal"}},
			expectError:  false,
		},
		"last-scan-failed": {
			target: &v1alpha1.MachineType{
				TypeMeta: metav1.TypeMeta{
					Kind:       "MachineType",
					APIVersion: "v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sample-machine-type",
					Namespace: "metal",
				},
				Spec: v1alpha1.MachineTypeSpec{
					Manufacturer: "Dell",
					Type:         "R440",
					ScanPeriod:   metav1.Duration{Duration: time.Hour * 24},
				},
				Status: v1alpha1.MachineTypeStatus{
					LastScanTime:   metav1.Time{Time: time.Now().Add(-10 * time.Minute)},
					LastScanResult: v1alpha1.ScanFailure,
				},
			},
			brokerClient: fake.NewFakeMachineTypeClientWithScans(map[string]*lcmimachinetype.ScanResponse{
				uuidutil.UUIDFromObjectKey(types.NamespacedName{Name: "sample-machine-type", Namespace: "metal"}): {
					Status: &lcmimachinetype.MachineTypeStatus{
						LastScanTime:   expectedLastScanTime.Unix(),
						LastScanResult: lcmicommon.ScanResult_SCAN_RESULT_SUCCESS,
					},
					State: lcmicommon.ScanState_SCAN_STATE_FINISHED,
				},
			}),
			request:     ctrl.Request{NamespacedName: types.NamespacedName{Name: "sample-machine-type", Namespace: "metal"}},
			expectError: false,
		},
		"last-scan-not-in-horizon": {
			target: &v1alpha1.MachineType{
				TypeMeta: metav1.TypeMeta{
					Kind:       "MachineType",
					APIVersion: "v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sample-machine-type",
					Namespace: "metal",
				},
				Spec: v1alpha1.MachineTypeSpec{
					Manufacturer: "Dell",
					Type:         "R440",
					ScanPeriod:   metav1.Duration{Duration: time.Hour * 24},
				},
				Status: v1alpha1.MachineTypeStatus{
					LastScanTime:   metav1.Time{Time: time.Now().Add(-50 * time.Minute)},
					LastScanResult: v1alpha1.ScanSuccess,
				},
			},
			brokerClient: fake.NewFakeMachineTypeClientWithScans(map[string]*lcmimachinetype.ScanResponse{
				uuidutil.UUIDFromObjectKey(types.NamespacedName{Name: "sample-machine-type", Namespace: "metal"}): {
					Status: &lcmimachinetype.MachineTypeStatus{
						LastScanTime:   time.Now().UnixNano(),
						LastScanResult: lcmicommon.ScanResult_SCAN_RESULT_SUCCESS,
					},
					State: lcmicommon.ScanState_SCAN_STATE_FINISHED,
				},
			}),
			request:     ctrl.Request{NamespacedName: types.NamespacedName{Name: "sample-machine-type", Namespace: "metal"}},
			expectError: false,
		},
	}

	for n, tc := range tests {
		name, testCase := n, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			clientOpts := []clientOption{withRuntimeObject(testCase.target)}
			r := newMachineTypeReconciler(t, schemeOpts, clientOpts)
			r.MachineTypeBroker = testCase.brokerClient
			resp, err := r.Reconcile(context.Background(), testCase.request)
			if testCase.expectError {
				assert.Error(t, err)
				assert.Empty(t, resp)
				return
			}
			assert.NoError(t, err)
			assert.Empty(t, resp)
			switch name {
			case "last-scan-failed":
				actualMachineType := &v1alpha1.MachineType{}
				err = r.Get(context.Background(), client.ObjectKeyFromObject(testCase.target), actualMachineType)
				assert.NoError(t, err)
				assert.True(t, actualMachineType.Status.LastScanResult.IsSuccess())
			case "last-scan-not-in-horizon":
				actualMachineType := &v1alpha1.MachineType{}
				err = r.Get(context.Background(), client.ObjectKeyFromObject(testCase.target), actualMachineType)
				assert.NoError(t, err)
				assert.Equal(t, expectedLastScanTime.Unix(), actualMachineType.Status.LastScanTime.Time.Unix())
			}
		})
	}
}
