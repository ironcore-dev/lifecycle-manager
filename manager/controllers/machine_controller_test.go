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
	lcmimachine "github.com/ironcore-dev/lifecycle-manager/lcmi/api/machine/v1alpha1"
	lcmimeta "github.com/ironcore-dev/lifecycle-manager/lcmi/api/meta/v1alpha1"
	"github.com/ironcore-dev/lifecycle-manager/lcmi/fake"
	"github.com/ironcore-dev/lifecycle-manager/util/uuidutil"
)

func TestMachineReconciler_Reconcile(t *testing.T) {
	t.Parallel()

	schemeOpts := []schemeOption{withGroupVersion(v1alpha1.SchemeBuilder)}
	expectedLastScanTime := time.Now()
	expectedInstalledPackages := []*lcmimeta.LocalObjectReference{
		{Name: "bios-v1.0.0"},
		{Name: "nic-v1.0.0"},
		{Name: "raid-v1.0.0"},
	}

	tests := map[string]struct {
		target       *v1alpha1.Machine
		brokerClient *fake.MachineClient
		request      ctrl.Request
		expectError  bool
	}{
		"absent-machine-type": {
			target: &v1alpha1.Machine{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Machine",
					APIVersion: "v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sample-machine",
					Namespace: "metal",
				},
				Spec: v1alpha1.MachineSpec{
					MachineTypeRef: corev1.LocalObjectReference{Name: "sample-machine-type"},
					OOBMachineRef:  corev1.LocalObjectReference{Name: "sample-oob"},
					ScanPeriod:     metav1.Duration{Duration: time.Hour * 24},
				},
				Status: v1alpha1.MachineStatus{},
			},
			request:     ctrl.Request{NamespacedName: types.NamespacedName{Name: "absent-machine", Namespace: "metal"}},
			expectError: false,
		},
		"deleting-machine-type": {
			target: &v1alpha1.Machine{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Machine",
					APIVersion: "v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:              "sample-machine",
					Namespace:         "metal",
					DeletionTimestamp: &metav1.Time{Time: time.Now()},
					Finalizers:        []string{"fake-finalizer"},
				},
				Spec: v1alpha1.MachineSpec{
					MachineTypeRef: corev1.LocalObjectReference{Name: "sample-machine-type"},
					OOBMachineRef:  corev1.LocalObjectReference{Name: "sample-oob"},
					ScanPeriod:     metav1.Duration{Duration: time.Hour * 24},
				},
				Status: v1alpha1.MachineStatus{},
			},
			request:     ctrl.Request{NamespacedName: types.NamespacedName{Name: "sample-machine", Namespace: "metal"}},
			expectError: false,
		},
		"no-scan-required": {
			target: &v1alpha1.Machine{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Machine",
					APIVersion: "v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sample-machine",
					Namespace: "metal",
				},
				Spec: v1alpha1.MachineSpec{
					MachineTypeRef: corev1.LocalObjectReference{Name: "sample-machine-type"},
					OOBMachineRef:  corev1.LocalObjectReference{Name: "sample-oob"},
					ScanPeriod:     metav1.Duration{Duration: time.Hour * 24},
				},
				Status: v1alpha1.MachineStatus{
					LastScanTime:   metav1.Time{Time: time.Now().Add(-10 * time.Minute)},
					LastScanResult: v1alpha1.ScanSuccess,
				},
			},
			request:     ctrl.Request{NamespacedName: types.NamespacedName{Name: "sample-machine", Namespace: "metal"}},
			expectError: false,
		},
		"scan-required-not-found": {
			target: &v1alpha1.Machine{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Machine",
					APIVersion: "v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sample-machine",
					Namespace: "metal",
				},
				Spec: v1alpha1.MachineSpec{
					MachineTypeRef: corev1.LocalObjectReference{Name: "sample-machine-type"},
					OOBMachineRef:  corev1.LocalObjectReference{Name: "sample-oob"},
					ScanPeriod:     metav1.Duration{Duration: time.Hour * 24},
				},
				Status: v1alpha1.MachineStatus{},
			},
			brokerClient: fake.NewMachineClient(),
			request:      ctrl.Request{NamespacedName: types.NamespacedName{Name: "sample-machine", Namespace: "metal"}},
			expectError:  false,
		},
		"last-scan-failed": {
			target: &v1alpha1.Machine{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Machine",
					APIVersion: "v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sample-machine",
					Namespace: "metal",
				},
				Spec: v1alpha1.MachineSpec{
					MachineTypeRef: corev1.LocalObjectReference{Name: "sample-machine-type"},
					OOBMachineRef:  corev1.LocalObjectReference{Name: "sample-oob"},
					ScanPeriod:     metav1.Duration{Duration: time.Hour * 24},
				},
				Status: v1alpha1.MachineStatus{
					LastScanTime:   metav1.Time{Time: time.Now().Add(-10 * time.Minute)},
					LastScanResult: v1alpha1.ScanFailure,
				},
			},
			brokerClient: fake.NewMachineClientWithScans(map[string]*lcmimachine.ScanResponse{
				uuidutil.UUIDFromObjectKey(types.NamespacedName{Name: "sample-machine", Namespace: "metal"}): {
					Status: &lcmimachine.MachineStatus{
						LastScanTime:      expectedLastScanTime.UnixNano(),
						LastScanResult:    lcmicommon.ScanResult_SCAN_RESULT_SUCCESS,
						InstalledPackages: expectedInstalledPackages,
					},
					State: lcmicommon.ScanState_SCAN_STATE_FINISHED,
				},
			}),
			request:     ctrl.Request{NamespacedName: types.NamespacedName{Name: "sample-machine", Namespace: "metal"}},
			expectError: false,
		},
		"last-scan-not-in-horizon": {
			target: &v1alpha1.Machine{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Machine",
					APIVersion: "v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sample-machine",
					Namespace: "metal",
				},
				Spec: v1alpha1.MachineSpec{
					MachineTypeRef: corev1.LocalObjectReference{Name: "sample-machine-type"},
					OOBMachineRef:  corev1.LocalObjectReference{Name: "sample-oob"},
					ScanPeriod:     metav1.Duration{Duration: time.Hour * 24},
				},
				Status: v1alpha1.MachineStatus{
					LastScanTime:   metav1.Time{Time: time.Now().Add(-50 * time.Minute)},
					LastScanResult: v1alpha1.ScanSuccess,
				},
			},
			brokerClient: fake.NewMachineClientWithScans(map[string]*lcmimachine.ScanResponse{
				uuidutil.UUIDFromObjectKey(types.NamespacedName{Name: "sample-machine", Namespace: "metal"}): {
					Status: &lcmimachine.MachineStatus{
						LastScanTime:      expectedLastScanTime.UnixNano(),
						LastScanResult:    lcmicommon.ScanResult_SCAN_RESULT_SUCCESS,
						InstalledPackages: expectedInstalledPackages,
					},
					State: lcmicommon.ScanState_SCAN_STATE_FINISHED,
				},
			}),
			request:     ctrl.Request{NamespacedName: types.NamespacedName{Name: "sample-machine", Namespace: "metal"}},
			expectError: false,
		},
	}

	for n, tc := range tests {
		name, testCase := n, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			clientOpts := []clientOption{withRuntimeObject(testCase.target)}
			r := newMachineReconciler(t, schemeOpts, clientOpts)
			r.BrokerClient = testCase.brokerClient
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
				actualMachine := &v1alpha1.Machine{}
				err = r.Get(context.Background(), client.ObjectKeyFromObject(testCase.target), actualMachine)
				assert.NoError(t, err)
				assert.True(t, actualMachine.Status.LastScanResult.IsSuccess())
				assert.Equal(t, expectedLastScanTime.Unix(), actualMachine.Status.LastScanTime.Time.Unix())
				assert.Equal(t, len(expectedInstalledPackages), len(actualMachine.Status.InstalledPackages))
			}
		})
	}
}
