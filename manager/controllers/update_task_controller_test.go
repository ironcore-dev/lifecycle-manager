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
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/ironcore-dev/lifecycle-manager/api/v1alpha1"
)

const (
	validObjectsKey = iota
	newUpdateTask
	existingUpdateTask
	missingLifecycleKey
	missingPackageKey
	mismatchObjectsKey
)

var (
	reconcileRequest = ctrl.Request{NamespacedName: types.NamespacedName{Namespace: "metal", Name: "sample-task"}}
	updateTasks      = map[int]*v1alpha1.UpdateTask{
		newUpdateTask: {
			TypeMeta: metav1.TypeMeta{
				Kind:       "UpdateTask",
				APIVersion: "v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "sample-task",
				Namespace: "metal",
			},
			Spec: v1alpha1.UpdateTaskSpec{
				MachineTypeRef: corev1.LocalObjectReference{Name: "Dell-R440"},
				Packages:       []corev1.LocalObjectReference{{Name: "sample-package"}},
				Machines:       []corev1.LocalObjectReference{{Name: "sample-lifecycle"}},
			},
			Status: v1alpha1.UpdateTaskStatus{},
		},
		existingUpdateTask: {
			TypeMeta: metav1.TypeMeta{
				Kind:       "UpdateTask",
				APIVersion: "v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "sample-task",
				Namespace: "metal",
				UID:       "11112222-3333-4444-5555-666677778888",
			},
			Spec: v1alpha1.UpdateTaskSpec{
				MachineTypeRef: corev1.LocalObjectReference{Name: "Dell-R440"},
				Packages:       []corev1.LocalObjectReference{{Name: "sample-package-BIOS"}, {Name: "sample-package-NIC"}},
				Machines:       []corev1.LocalObjectReference{{Name: "sample-lifecycle"}},
			},
			Status: v1alpha1.UpdateTaskStatus{
				JobsTotal:      2,
				JobsFailed:     0,
				JobsSuccessful: 0,
			},
		},
	}
	machineTypes = map[int]*v1alpha1.MachineType{
		validObjectsKey: {
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
		},
	}
	machineLifecycles = map[int][]*v1alpha1.MachineLifecycle{
		validObjectsKey: {
			{
				TypeMeta: metav1.TypeMeta{
					Kind:       "MachineLifecycle",
					APIVersion: "v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sample-lifecycle",
					Namespace: "metal",
				},
				Spec: v1alpha1.MachineLifecycleSpec{
					MachineTypeRef: corev1.LocalObjectReference{Name: "Dell-R440"},
					OOBMachineRef:  corev1.LocalObjectReference{Name: "sample-lifecycle"},
					ScanPeriod:     metav1.Duration{Duration: time.Hour * 24},
				},
			},
		},
		missingLifecycleKey: {},
		mismatchObjectsKey: {
			{
				TypeMeta: metav1.TypeMeta{
					Kind:       "MachineLifecycle",
					APIVersion: "v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sample-lifecycle",
					Namespace: "metal",
				},
				Spec: v1alpha1.MachineLifecycleSpec{
					MachineTypeRef: corev1.LocalObjectReference{Name: "Lenovo-7X21"},
					OOBMachineRef:  corev1.LocalObjectReference{Name: "sample-lifecycle"},
					ScanPeriod:     metav1.Duration{Duration: time.Hour * 24},
				},
			},
		},
	}
	firmwarePackages = map[int][]*v1alpha1.FirmwarePackage{
		validObjectsKey: {
			{
				TypeMeta: metav1.TypeMeta{
					Kind:       "FirmwarePackage",
					APIVersion: "v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sample-package",
					Namespace: "metal",
				},
				Spec: v1alpha1.FirmwarePackageSpec{
					MachineTypeRef: corev1.LocalObjectReference{Name: "Dell-R440"},
					Name:           "BIOS",
					Version:        "v1.0.0",
					Source:         "https://fake.package.store.com",
					RebootRequired: false,
				},
			},
		},
		missingPackageKey: {},
		mismatchObjectsKey: {
			{
				TypeMeta: metav1.TypeMeta{
					Kind:       "FirmwarePackage",
					APIVersion: "v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sample-package",
					Namespace: "metal",
				},
				Spec: v1alpha1.FirmwarePackageSpec{
					MachineTypeRef: corev1.LocalObjectReference{Name: "Lenovo-7X21"},
					Name:           "BIOS",
					Version:        "v1.0.0",
					Source:         "https://fake.package.store.com",
					RebootRequired: false,
				},
			},
		},
	}
	updateJobs = map[int][]*v1alpha1.MachineUpdateJob{
		validObjectsKey: {
			{
				TypeMeta: metav1.TypeMeta{
					Kind:       "FirmwarePackage",
					APIVersion: "v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sample-job-success",
					Namespace: "metal",
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion: "v1alpha1",
							Kind:       "UpdateTask",
							Name:       "sample-task",
							UID:        "11112222-3333-4444-5555-666677778888",
							Controller: pointer.Bool(true),
						},
					},
				},
				Spec: v1alpha1.MachineUpdateJobSpec{
					MachineLifecycleRef: corev1.LocalObjectReference{Name: "sample-lifecycle"},
					FirmwarePackageRef:  corev1.LocalObjectReference{Name: "sample-package-BIOS"},
				},
				Status: v1alpha1.MachineUpdateJobStatus{
					Conditions: nil,
					State:      "Success",
					Message:    "",
				},
			},
			{
				TypeMeta: metav1.TypeMeta{
					Kind:       "FirmwarePackage",
					APIVersion: "v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sample-job-failure",
					Namespace: "metal",
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion: "v1alpha1",
							Kind:       "UpdateTask",
							Name:       "sample-task",
							UID:        "11112222-3333-4444-5555-666677778888",
							Controller: pointer.Bool(true),
						},
					},
				},
				Spec: v1alpha1.MachineUpdateJobSpec{
					MachineLifecycleRef: corev1.LocalObjectReference{Name: "sample-lifecycle"},
					FirmwarePackageRef:  corev1.LocalObjectReference{Name: "sample-package-NIC"},
				},
				Status: v1alpha1.MachineUpdateJobStatus{
					Conditions: nil,
					State:      "Failure",
					Message:    "Failed to download update package from remote",
				},
			},
		},
	}
)

func TestUpdateTaskReconciler_Reconcile(t *testing.T) {
	t.Parallel()

	schemeOpts := []schemeOption{withGroupVersion(v1alpha1.SchemeBuilder)}

	tests := map[string]struct {
		target      *v1alpha1.UpdateTask
		request     ctrl.Request
		expectError bool
	}{
		"deleting-task": {
			target: &v1alpha1.UpdateTask{
				TypeMeta: metav1.TypeMeta{
					Kind:       "UpdateTask",
					APIVersion: "v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:              "deleting-task",
					Namespace:         "metal",
					DeletionTimestamp: &metav1.Time{Time: time.Now()},
					Finalizers:        []string{"fake-finalizer"},
				},
				Spec: v1alpha1.UpdateTaskSpec{
					MachineTypeRef: corev1.LocalObjectReference{Name: "fake-type"},
					Packages:       []corev1.LocalObjectReference{{Name: "fake-package"}},
					Machines:       []corev1.LocalObjectReference{{Name: "fake-lifecycle"}},
				},
				Status: v1alpha1.UpdateTaskStatus{},
			},
			request:     ctrl.Request{NamespacedName: types.NamespacedName{Name: "deleting-task", Namespace: "metal"}},
			expectError: false,
		},
		"absent-machine": {
			target: &v1alpha1.UpdateTask{
				TypeMeta: metav1.TypeMeta{
					Kind:       "UpdateTask",
					APIVersion: "v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sample-task",
					Namespace: "metal",
				},
				Spec: v1alpha1.UpdateTaskSpec{
					MachineTypeRef: corev1.LocalObjectReference{Name: "fake-type"},
					Packages:       []corev1.LocalObjectReference{{Name: "fake-package"}},
					Machines:       []corev1.LocalObjectReference{{Name: "fake-lifecycle"}},
				},
				Status: v1alpha1.UpdateTaskStatus{},
			},
			request:     ctrl.Request{NamespacedName: types.NamespacedName{Name: "sample-task", Namespace: "metal"}},
			expectError: true,
		},
	}

	for n, tc := range tests {
		name, testCase := n, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			clientOpts := []clientOption{withRuntimeObject(testCase.target)}
			r := newUpdateTaskReconciler(t, schemeOpts, clientOpts)
			resp, err := r.Reconcile(context.Background(), testCase.request)
			if testCase.expectError {
				assert.Error(t, err)
				assert.Empty(t, resp)
				return
			}
			assert.NoError(t, err)
			assert.Empty(t, resp)
		})
	}
}

func TestUpdateTaskReconciler_Reconcile2(t *testing.T) {
	t.Parallel()

	schemeOpts := []schemeOption{withGroupVersion(v1alpha1.SchemeBuilder)}

	tests := map[string]struct {
		updateTaskKey        int
		machineTypeKey       int
		machineLifecyclesKey int
		firmwarePackagesKey  int
		updateJobsKey        int
		request              ctrl.Request
		sampleKey            string
		expectError          bool
	}{
		"happy-path-new-task": {
			updateTaskKey:        newUpdateTask,
			machineTypeKey:       validObjectsKey,
			machineLifecyclesKey: validObjectsKey,
			firmwarePackagesKey:  validObjectsKey,
			request:              reconcileRequest,
			expectError:          false,
		},
		"dangling-lifecycle-ref": {
			updateTaskKey:        newUpdateTask,
			machineTypeKey:       validObjectsKey,
			machineLifecyclesKey: missingLifecycleKey,
			firmwarePackagesKey:  validObjectsKey,
			request:              reconcileRequest,
			expectError:          true,
		},
		"dangling-package-ref": {
			updateTaskKey:        newUpdateTask,
			machineTypeKey:       validObjectsKey,
			machineLifecyclesKey: validObjectsKey,
			firmwarePackagesKey:  missingPackageKey,
			request:              reconcileRequest,
			expectError:          true,
		},
		"lifecycle-type-mismatch": {
			updateTaskKey:        newUpdateTask,
			machineTypeKey:       validObjectsKey,
			machineLifecyclesKey: mismatchObjectsKey,
			firmwarePackagesKey:  validObjectsKey,
			request:              reconcileRequest,
			expectError:          false,
		},
		"package-type-mismatch": {
			updateTaskKey:        newUpdateTask,
			machineTypeKey:       validObjectsKey,
			machineLifecyclesKey: validObjectsKey,
			firmwarePackagesKey:  mismatchObjectsKey,
			request:              reconcileRequest,
			expectError:          false,
		},
		"happy-path-existing-task": {
			updateTaskKey:        existingUpdateTask,
			machineTypeKey:       validObjectsKey,
			machineLifecyclesKey: validObjectsKey,
			firmwarePackagesKey:  validObjectsKey,
			updateJobsKey:        validObjectsKey,
			request:              reconcileRequest,
			expectError:          false,
		},
	}

	for n, tc := range tests {
		name, testCase := n, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			clientOpts := []clientOption{
				withRuntimeObject(updateTasks[testCase.updateTaskKey]),
				withRuntimeObject(machineTypes[testCase.machineTypeKey]),
			}
			for _, o := range machineLifecycles[testCase.machineLifecyclesKey] {
				clientOpts = append(clientOpts, withRuntimeObject(o))
			}
			for _, o := range firmwarePackages[testCase.firmwarePackagesKey] {
				clientOpts = append(clientOpts, withRuntimeObject(o))
			}
			for _, o := range updateJobs[testCase.updateJobsKey] {
				clientOpts = append(clientOpts, withRuntimeObject(o))
			}
			r := newUpdateTaskReconciler(t, schemeOpts, clientOpts)
			_, err := r.Reconcile(context.Background(), testCase.request)
			if testCase.expectError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			actualUpdateTask := &v1alpha1.UpdateTask{}
			err = r.Get(context.Background(), testCase.request.NamespacedName, actualUpdateTask)
			assert.NoError(t, err)
			switch name {
			case "lifecycle-type-mismatch", "package-type-mismatch":
				assert.Equal(t, uint32(0), actualUpdateTask.Status.JobsTotal)
			case "happy-path-existing-task":
				assert.Equal(t, uint32(1), actualUpdateTask.Status.JobsFailed)
				assert.Equal(t, uint32(1), actualUpdateTask.Status.JobsSuccessful)
			default:
				assert.Equal(t, uint32(1), actualUpdateTask.Status.JobsTotal)
			}
		})
	}
}
