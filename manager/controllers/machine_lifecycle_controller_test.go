// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/ironcore-dev/lifecycle-manager/api/v1alpha1"
	lcmicommon "github.com/ironcore-dev/lifecycle-manager/lcmi/api/common/v1alpha1"
	lcmimachinelifecycle "github.com/ironcore-dev/lifecycle-manager/lcmi/api/machine_lifecycle/v1alpha1"
	alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/meta/v1alpha1"
	"github.com/ironcore-dev/lifecycle-manager/lcmi/fake"
	"github.com/ironcore-dev/lifecycle-manager/util/uuidutil"
)

func TestMachineLifecycleReconciler_Reconcile(t *testing.T) {
	t.Parallel()

	schemeOpts := []schemeOption{withGroupVersion(v1alpha1.SchemeBuilder)}
	expectedLastScanTime := time.Now()
	expectedInstalledPackages := []*alpha1.LocalObjectReference{
		{Name: "bios-v1.0.0"},
		{Name: "nic-v1.0.0"},
		{Name: "raid-v1.0.0"},
	}

	tests := map[string]struct {
		target       *v1alpha1.MachineLifecycle
		brokerClient *fake.MachineLifecycleClient
		request      ctrl.Request
		expectError  bool
	}{
		"absent-machine-type": {
			target: &v1alpha1.MachineLifecycle{
				TypeMeta: metav1.TypeMeta{
					Kind:       "MachineLifecycle",
					APIVersion: "v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sample-machine-lifecycle",
					Namespace: "metal",
				},
				Spec: v1alpha1.MachineLifecycleSpec{
					MachineTypeRef: corev1.LocalObjectReference{Name: "sample-machine-type"},
					OOBMachineRef:  corev1.LocalObjectReference{Name: "sample-oob"},
					ScanPeriod:     metav1.Duration{Duration: time.Hour * 24},
				},
				Status: v1alpha1.MachineLifecycleStatus{},
			},
			request:     ctrl.Request{NamespacedName: types.NamespacedName{Name: "absent-machine-lifecycle", Namespace: "metal"}},
			expectError: false,
		},
		"deleting-machine-type": {
			target: &v1alpha1.MachineLifecycle{
				TypeMeta: metav1.TypeMeta{
					Kind:       "MachineLifecycle",
					APIVersion: "v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:              "sample-machine-lifecycle",
					Namespace:         "metal",
					DeletionTimestamp: &metav1.Time{Time: time.Now()},
					Finalizers:        []string{"fake-finalizer"},
				},
				Spec: v1alpha1.MachineLifecycleSpec{
					MachineTypeRef: corev1.LocalObjectReference{Name: "sample-machine-type"},
					OOBMachineRef:  corev1.LocalObjectReference{Name: "sample-oob"},
					ScanPeriod:     metav1.Duration{Duration: time.Hour * 24},
				},
				Status: v1alpha1.MachineLifecycleStatus{},
			},
			request:     ctrl.Request{NamespacedName: types.NamespacedName{Name: "sample-machine-lifecycle", Namespace: "metal"}},
			expectError: false,
		},
		"no-scan-required": {
			target: &v1alpha1.MachineLifecycle{
				TypeMeta: metav1.TypeMeta{
					Kind:       "MachineLifecycle",
					APIVersion: "v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sample-machine-lifecycle",
					Namespace: "metal",
				},
				Spec: v1alpha1.MachineLifecycleSpec{
					MachineTypeRef: corev1.LocalObjectReference{Name: "sample-machine-type"},
					OOBMachineRef:  corev1.LocalObjectReference{Name: "sample-oob"},
					ScanPeriod:     metav1.Duration{Duration: time.Hour * 24},
				},
				Status: v1alpha1.MachineLifecycleStatus{
					LastScanTime:   metav1.Time{Time: time.Now().Add(-10 * time.Minute)},
					LastScanResult: v1alpha1.ScanSuccess,
				},
			},
			request:     ctrl.Request{NamespacedName: types.NamespacedName{Name: "sample-machine-lifecycle", Namespace: "metal"}},
			expectError: false,
		},
		"scan-required-not-found": {
			target: &v1alpha1.MachineLifecycle{
				TypeMeta: metav1.TypeMeta{
					Kind:       "MachineLifecycle",
					APIVersion: "v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sample-machine-lifecycle",
					Namespace: "metal",
				},
				Spec: v1alpha1.MachineLifecycleSpec{
					MachineTypeRef: corev1.LocalObjectReference{Name: "sample-machine-type"},
					OOBMachineRef:  corev1.LocalObjectReference{Name: "sample-oob"},
					ScanPeriod:     metav1.Duration{Duration: time.Hour * 24},
				},
				Status: v1alpha1.MachineLifecycleStatus{},
			},
			brokerClient: fake.NewMachineLifecycleClient(),
			request:      ctrl.Request{NamespacedName: types.NamespacedName{Name: "sample-machine-lifecycle", Namespace: "metal"}},
			expectError:  false,
		},
		"last-scan-failed": {
			target: &v1alpha1.MachineLifecycle{
				TypeMeta: metav1.TypeMeta{
					Kind:       "MachineLifecycle",
					APIVersion: "v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sample-machine-lifecycle",
					Namespace: "metal",
				},
				Spec: v1alpha1.MachineLifecycleSpec{
					MachineTypeRef: corev1.LocalObjectReference{Name: "sample-machine-type"},
					OOBMachineRef:  corev1.LocalObjectReference{Name: "sample-oob"},
					ScanPeriod:     metav1.Duration{Duration: time.Hour * 24},
				},
				Status: v1alpha1.MachineLifecycleStatus{
					LastScanTime:   metav1.Time{Time: time.Now().Add(-10 * time.Minute)},
					LastScanResult: v1alpha1.ScanFailure,
				},
			},
			brokerClient: fake.NewMachineLifecycleClientWithScans(map[string]*lcmimachinelifecycle.ScanResponse{
				uuidutil.UUIDFromObjectKey(types.NamespacedName{Name: "sample-machine-lifecycle", Namespace: "metal"}): {
					Status: &lcmimachinelifecycle.MachineLifecycleStatus{
						LastScanTime:      expectedLastScanTime.UnixNano(),
						LastScanResult:    lcmicommon.ScanResult_SCAN_RESULT_SUCCESS,
						InstalledPackages: expectedInstalledPackages,
					},
					State: lcmicommon.ScanState_SCAN_STATE_FINISHED,
				},
			}),
			request:     ctrl.Request{NamespacedName: types.NamespacedName{Name: "sample-machine-lifecycle", Namespace: "metal"}},
			expectError: false,
		},
		"last-scan-not-in-horizon": {
			target: &v1alpha1.MachineLifecycle{
				TypeMeta: metav1.TypeMeta{
					Kind:       "MachineLifecycle",
					APIVersion: "v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sample-machine-lifecycle",
					Namespace: "metal",
				},
				Spec: v1alpha1.MachineLifecycleSpec{
					MachineTypeRef: corev1.LocalObjectReference{Name: "sample-machine-type"},
					OOBMachineRef:  corev1.LocalObjectReference{Name: "sample-oob"},
					ScanPeriod:     metav1.Duration{Duration: time.Hour * 24},
				},
				Status: v1alpha1.MachineLifecycleStatus{
					LastScanTime:   metav1.Time{Time: time.Now().Add(-50 * time.Minute)},
					LastScanResult: v1alpha1.ScanSuccess,
				},
			},
			brokerClient: fake.NewMachineLifecycleClientWithScans(map[string]*lcmimachinelifecycle.ScanResponse{
				uuidutil.UUIDFromObjectKey(types.NamespacedName{Name: "sample-machine-lifecycle", Namespace: "metal"}): {
					Status: &lcmimachinelifecycle.MachineLifecycleStatus{
						LastScanTime:      expectedLastScanTime.UnixNano(),
						LastScanResult:    lcmicommon.ScanResult_SCAN_RESULT_SUCCESS,
						InstalledPackages: expectedInstalledPackages,
					},
					State: lcmicommon.ScanState_SCAN_STATE_FINISHED,
				},
			}),
			request:     ctrl.Request{NamespacedName: types.NamespacedName{Name: "sample-machine-lifecycle", Namespace: "metal"}},
			expectError: false,
		},
	}

	for n, tc := range tests {
		name, testCase := n, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			clientOpts := []clientOption{withRuntimeObject(testCase.target)}
			r := newMachineLifecycleReconciler(t, schemeOpts, clientOpts)
			r.MachineLifecycleBroker = testCase.brokerClient
			resp, err := r.Reconcile(context.Background(), testCase.request)
			if testCase.expectError {
				assert.Error(t, err)
				assert.Empty(t, resp)
				return
			}
			assert.NoError(t, err)
			assert.Empty(t, resp)
			switch name {
			case "last-scan-failed", "last-scan-not-in-horizon":
				actualMachineLifecycle := &v1alpha1.MachineLifecycle{}
				err = r.Get(context.Background(), client.ObjectKeyFromObject(testCase.target), actualMachineLifecycle)
				assert.NoError(t, err)
				assert.True(t, actualMachineLifecycle.Status.LastScanResult.IsSuccess())
				assert.Equal(t, expectedLastScanTime.Unix(), actualMachineLifecycle.Status.LastScanTime.Time.Unix())
				assert.Equal(t, len(expectedInstalledPackages), len(actualMachineLifecycle.Status.InstalledPackages))
			}
		})
	}
}
