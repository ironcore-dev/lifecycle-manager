// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	lifecyclev1alpha1 "github.com/ironcore-dev/lifecycle-manager/api/lifecycle/v1alpha1"
	"github.com/ironcore-dev/lifecycle-manager/lcmi/fake"
	"github.com/ironcore-dev/lifecycle-manager/util/uuidutil"
)

func TestMachineReconciler_Reconcile(t *testing.T) {
	t.Parallel()

	schemeOpts := []schemeOption{withGroupVersion(lifecyclev1alpha1.AddToScheme)}
	now := time.Unix(time.Now().Unix(), 0)

	tests := map[string]struct {
		machine      *lifecyclev1alpha1.Machine
		machineType  *lifecyclev1alpha1.MachineType
		request      ctrl.Request
		brokerClient *fake.MachineClient
		expectError  bool
	}{
		"absent-machine": {
			machine:     &lifecyclev1alpha1.Machine{},
			machineType: &lifecyclev1alpha1.MachineType{},
			request: ctrl.Request{NamespacedName: types.NamespacedName{
				Name:      "absent-machine",
				Namespace: "ironcore",
			}},
			expectError: false,
		},
		"machine-being-deleted": {
			machine: &lifecyclev1alpha1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "sample-machine",
					Namespace:         "ironcore",
					DeletionTimestamp: &metav1.Time{Time: time.Now()},
					Finalizers:        []string{"test-finalizer"},
				},
				Spec: lifecyclev1alpha1.MachineSpec{
					MachineTypeRef: corev1.LocalObjectReference{Name: "sample-machine-type"},
					OOBMachineRef:  corev1.LocalObjectReference{Name: "sample-machine"},
					ScanPeriod:     metav1.Duration{Duration: time.Hour * 24},
				},
			},
			machineType: &lifecyclev1alpha1.MachineType{},
			request: ctrl.Request{NamespacedName: types.NamespacedName{
				Name:      "sample-machine",
				Namespace: "ironcore",
			}},
			expectError: false,
		},
		"absent-machine-type": {
			machine: &lifecyclev1alpha1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sample-machine",
					Namespace: "ironcore",
				},
				Spec: lifecyclev1alpha1.MachineSpec{
					MachineTypeRef: corev1.LocalObjectReference{Name: "absent-machine-type"},
					OOBMachineRef:  corev1.LocalObjectReference{Name: "sample-machine"},
					ScanPeriod:     metav1.Duration{Duration: time.Hour * 24},
				},
			},
			machineType: &lifecyclev1alpha1.MachineType{},
			request: ctrl.Request{NamespacedName: types.NamespacedName{
				Name:      "sample-machine",
				Namespace: "ironcore",
			}},
			brokerClient: fake.NewMachineClient(map[string]*lifecyclev1alpha1.MachineStatus{
				uuidutil.UUIDFromObjectKey(types.NamespacedName{
					Name:      "sample-machine",
					Namespace: "ironcore",
				}): {},
			}),
			expectError: true,
		},
		"scan-request-submitted": {
			machine: &lifecyclev1alpha1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sample-machine",
					Namespace: "ironcore",
				},
				Spec: lifecyclev1alpha1.MachineSpec{
					MachineTypeRef: corev1.LocalObjectReference{Name: "sample-machine-type"},
					OOBMachineRef:  corev1.LocalObjectReference{Name: "sample-machine"},
					ScanPeriod:     metav1.Duration{Duration: time.Hour * 24},
				},
			},
			machineType: &lifecyclev1alpha1.MachineType{},
			request: ctrl.Request{NamespacedName: types.NamespacedName{
				Name:      "sample-machine",
				Namespace: "ironcore",
			}},
			brokerClient: fake.NewMachineClient(map[string]*lifecyclev1alpha1.MachineStatus{}),
			expectError:  false,
		},
		"packages-installed": {
			machine: &lifecyclev1alpha1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sample-machine",
					Namespace: "ironcore",
					Labels:    map[string]string{"env": "test"},
				},
				Spec: lifecyclev1alpha1.MachineSpec{
					MachineTypeRef: corev1.LocalObjectReference{Name: "sample-machine-type"},
					OOBMachineRef:  corev1.LocalObjectReference{Name: "sample-machine"},
					ScanPeriod:     metav1.Duration{Duration: time.Hour * 24},
					Packages: []lifecyclev1alpha1.PackageVersion{
						{Name: "bios", Version: "1.0.0"},
						{Name: "bmc", Version: "1.0.0"},
					},
				},
				Status: lifecyclev1alpha1.MachineStatus{
					InstalledPackages: []lifecyclev1alpha1.PackageVersion{
						{Name: "bios", Version: "1.0.0"},
						{Name: "bmc", Version: "1.0.0"},
						{Name: "raid", Version: "3.2.1"},
						{Name: "net", Version: "0.1.0"},
					},
				},
			},
			machineType: &lifecyclev1alpha1.MachineType{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sample-machine-type",
					Namespace: "ironcore",
				},
				Spec: lifecyclev1alpha1.MachineTypeSpec{
					Manufacturer: "Sample",
					Type:         "abc",
					ScanPeriod:   metav1.Duration{Duration: time.Hour * 24},
					MachineGroups: []lifecyclev1alpha1.MachineGroup{
						{
							MachineSelector: metav1.LabelSelector{
								MatchLabels: map[string]string{"env": "test"},
							},
							Packages: []lifecyclev1alpha1.PackageVersion{
								{Name: "raid", Version: "3.2.1"},
							},
						},
					},
				},
			},
			request: ctrl.Request{NamespacedName: types.NamespacedName{
				Name:      "sample-machine",
				Namespace: "ironcore",
			}},
			brokerClient: fake.NewMachineClient(map[string]*lifecyclev1alpha1.MachineStatus{
				uuidutil.UUIDFromObjectKey(types.NamespacedName{
					Name:      "sample-machine",
					Namespace: "ironcore",
				}): {
					LastScanTime:   metav1.Time{Time: now},
					LastScanResult: lifecyclev1alpha1.ScanSuccess,
					InstalledPackages: []lifecyclev1alpha1.PackageVersion{
						{Name: "bios", Version: "1.0.0"},
						{Name: "bmc", Version: "1.0.0"},
						{Name: "raid", Version: "3.2.1"},
						{Name: "net", Version: "0.1.0"},
					},
				},
			}),
			expectError: false,
		},
		"installation-scheduled": {
			machine: &lifecyclev1alpha1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sample-machine",
					Namespace: "ironcore",
					Labels:    map[string]string{"env": "test"},
				},
				Spec: lifecyclev1alpha1.MachineSpec{
					MachineTypeRef: corev1.LocalObjectReference{Name: "sample-machine-type"},
					OOBMachineRef:  corev1.LocalObjectReference{Name: "sample-machine"},
					ScanPeriod:     metav1.Duration{Duration: time.Hour * 24},
					Packages: []lifecyclev1alpha1.PackageVersion{
						{Name: "bios", Version: "1.0.0"},
						{Name: "bmc", Version: "1.0.0"},
					},
				},
				Status: lifecyclev1alpha1.MachineStatus{
					InstalledPackages: []lifecyclev1alpha1.PackageVersion{
						{Name: "bios", Version: "1.0.0"},
						{Name: "bmc", Version: "1.0.0"},
						{Name: "raid", Version: "3.2.1"},
						{Name: "net", Version: "0.1.0"},
					},
				},
			},
			machineType: &lifecyclev1alpha1.MachineType{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sample-machine-type",
					Namespace: "ironcore",
				},
				Spec: lifecyclev1alpha1.MachineTypeSpec{
					Manufacturer: "Sample",
					Type:         "abc",
					ScanPeriod:   metav1.Duration{Duration: time.Hour * 24},
					MachineGroups: []lifecyclev1alpha1.MachineGroup{
						{
							MachineSelector: metav1.LabelSelector{
								MatchLabels: map[string]string{"env": "test"},
							},
							Packages: []lifecyclev1alpha1.PackageVersion{
								{Name: "raid", Version: "3.2.1"},
							},
						},
					},
				},
			},
			request: ctrl.Request{NamespacedName: types.NamespacedName{
				Name:      "sample-machine",
				Namespace: "ironcore",
			}},
			brokerClient: fake.NewMachineClient(map[string]*lifecyclev1alpha1.MachineStatus{
				uuidutil.UUIDFromObjectKey(types.NamespacedName{
					Name:      "sample-machine",
					Namespace: "ironcore",
				}): {
					LastScanTime:   metav1.Time{Time: now},
					LastScanResult: lifecyclev1alpha1.ScanSuccess,
					InstalledPackages: []lifecyclev1alpha1.PackageVersion{
						{Name: "bios", Version: "2.0.0"},
						{Name: "bmc", Version: "1.0.0"},
						{Name: "raid", Version: "3.2.1"},
						{Name: "net", Version: "0.1.0"},
					},
				},
			}),
			expectError: false,
		},
	}

	for n, tc := range tests {
		name, tt := n, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			clientOpts := []clientOption{
				withRuntimeObject(tt.machine),
				withRuntimeObject(tt.machineType),
			}
			r := newMachineReconciler(t, schemeOpts, clientOpts)
			r.Broker = tt.brokerClient
			resp, err := r.Reconcile(context.Background(), tt.request)
			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, resp)
				return
			}
			assert.NoError(t, err)
			switch name {
			case "scan-request-submitted":
				broker, _ := r.Broker.(*fake.MachineClient)
				entry := broker.ReadCache(uuidutil.UUIDFromObjectKey(client.ObjectKeyFromObject(tt.machine)))
				reconciledMachine := &lifecyclev1alpha1.Machine{}
				_ = r.Get(context.Background(), client.ObjectKeyFromObject(tt.machine), reconciledMachine)
				assert.NotNil(t, entry)
				assert.Equal(t, StatusMessageScanRequestSubmitted, reconciledMachine.Status.Message)
			case "packages-installed":
				broker, _ := r.Broker.(*fake.MachineClient)
				entry := broker.ReadCache(uuidutil.UUIDFromObjectKey(client.ObjectKeyFromObject(tt.machine)))
				reconciledMachine := &lifecyclev1alpha1.Machine{}
				_ = r.Get(context.Background(), client.ObjectKeyFromObject(tt.machine), reconciledMachine)
				assert.Equal(t, *entry, reconciledMachine.Status)
			case "installation-scheduled":
				reconciledMachine := &lifecyclev1alpha1.Machine{}
				_ = r.Get(context.Background(), client.ObjectKeyFromObject(tt.machine), reconciledMachine)
				assert.Equal(t, StatusMessageInstallationScheduled, reconciledMachine.Status.Message)
			}
		})
	}
}

func TestMachineReconciler_PackagesToInstall(t *testing.T) {
	t.Parallel()

	schemeOpts := []schemeOption{withGroupVersion(lifecyclev1alpha1.AddToScheme)}

	machineType := &lifecyclev1alpha1.MachineType{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sample-machine-type",
			Namespace: "ironcore",
		},
		Spec: lifecyclev1alpha1.MachineTypeSpec{
			Manufacturer: "Sample",
			Type:         "abc",
			ScanPeriod:   metav1.Duration{Duration: time.Hour * 24},
			MachineGroups: []lifecyclev1alpha1.MachineGroup{
				{
					MachineSelector: metav1.LabelSelector{
						MatchLabels: map[string]string{"env": "test"},
					},
				},
			},
		},
	}

	machine := &lifecyclev1alpha1.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sample-machine",
			Namespace: "ironcore",
			Labels:    map[string]string{"env": "test"},
		},
		Spec: lifecyclev1alpha1.MachineSpec{
			MachineTypeRef: corev1.LocalObjectReference{Name: "sample-machine-type"},
			OOBMachineRef:  corev1.LocalObjectReference{Name: "sample-machine"},
			ScanPeriod:     metav1.Duration{Duration: time.Hour * 24},
		},
	}

	tests := map[string]struct {
		desiredPackages   []lifecyclev1alpha1.PackageVersion
		installedPackages []lifecyclev1alpha1.PackageVersion
		defaultPackages   []lifecyclev1alpha1.PackageVersion
		expectedPackages  []lifecyclev1alpha1.PackageVersion
	}{
		"no-packages-to-install": {
			desiredPackages: []lifecyclev1alpha1.PackageVersion{
				{Name: "bios", Version: "1.0.0"},
				{Name: "bmc", Version: "1.0.0"},
			},
			installedPackages: []lifecyclev1alpha1.PackageVersion{
				{Name: "bios", Version: "1.0.0"},
				{Name: "bmc", Version: "1.0.0"},
				{Name: "raid", Version: "3.2.1"},
				{Name: "net", Version: "0.1.0"},
			},
			defaultPackages: []lifecyclev1alpha1.PackageVersion{
				{Name: "raid", Version: "3.2.1"},
			},
			expectedPackages: []lifecyclev1alpha1.PackageVersion{},
		},
		"desired-to-update": {
			desiredPackages: []lifecyclev1alpha1.PackageVersion{
				{Name: "bios", Version: "2.0.0"},
				{Name: "bmc", Version: "1.0.0"},
			},
			installedPackages: []lifecyclev1alpha1.PackageVersion{
				{Name: "bios", Version: "1.0.0"},
				{Name: "bmc", Version: "1.0.0"},
				{Name: "raid", Version: "3.2.1"},
				{Name: "net", Version: "0.1.0"},
			},
			defaultPackages: []lifecyclev1alpha1.PackageVersion{
				{Name: "raid", Version: "3.2.1"},
			},
			expectedPackages: []lifecyclev1alpha1.PackageVersion{
				{Name: "bios", Version: "2.0.0"},
			},
		},
		"default-to-update": {
			desiredPackages: []lifecyclev1alpha1.PackageVersion{
				{Name: "bios", Version: "1.0.0"},
				{Name: "bmc", Version: "1.0.0"},
			},
			installedPackages: []lifecyclev1alpha1.PackageVersion{
				{Name: "bios", Version: "1.0.0"},
				{Name: "bmc", Version: "1.0.0"},
				{Name: "raid", Version: "3.2.1"},
				{Name: "net", Version: "0.1.0"},
			},
			defaultPackages: []lifecyclev1alpha1.PackageVersion{
				{Name: "raid", Version: "4.0.0"},
			},
			expectedPackages: []lifecyclev1alpha1.PackageVersion{
				{Name: "raid", Version: "4.0.0"},
			},
		},
		"desired-and-default-to-update": {
			desiredPackages: []lifecyclev1alpha1.PackageVersion{
				{Name: "bios", Version: "1.0.3"},
				{Name: "bmc", Version: "1.7.0"},
			},
			installedPackages: []lifecyclev1alpha1.PackageVersion{
				{Name: "bios", Version: "1.0.0"},
				{Name: "bmc", Version: "1.0.0"},
				{Name: "raid", Version: "3.2.1"},
				{Name: "net", Version: "0.1.0"},
			},
			defaultPackages: []lifecyclev1alpha1.PackageVersion{
				{Name: "raid", Version: "4.0.0"},
				{Name: "net", Version: "1.5.0"},
			},
			expectedPackages: []lifecyclev1alpha1.PackageVersion{
				{Name: "raid", Version: "4.0.0"},
				{Name: "bios", Version: "1.0.3"},
				{Name: "bmc", Version: "1.7.0"},
				{Name: "net", Version: "1.5.0"},
			},
		},
		"desired-overrides-default": {
			desiredPackages: []lifecyclev1alpha1.PackageVersion{
				{Name: "bios", Version: "2.0.0"},
				{Name: "bmc", Version: "1.0.0"},
				{Name: "raid", Version: "5.5.0"},
			},
			installedPackages: []lifecyclev1alpha1.PackageVersion{
				{Name: "bios", Version: "1.0.0"},
				{Name: "bmc", Version: "1.0.0"},
				{Name: "raid", Version: "3.2.1"},
				{Name: "net", Version: "0.1.0"},
			},
			defaultPackages: []lifecyclev1alpha1.PackageVersion{
				{Name: "raid", Version: "3.2.1"},
			},
			expectedPackages: []lifecyclev1alpha1.PackageVersion{
				{Name: "bios", Version: "2.0.0"},
				{Name: "raid", Version: "5.5.0"},
			},
		},
	}

	for n, tc := range tests {
		name, tt := n, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			testMachineType := machineType.DeepCopy()
			testMachineType.Spec.MachineGroups[0].Packages = tt.defaultPackages
			testMachine := machine.DeepCopy()
			testMachine.Spec.Packages = tt.desiredPackages
			testMachine.Status.InstalledPackages = tt.installedPackages
			clientOpts := []clientOption{withRuntimeObject(testMachineType)}
			r := newMachineReconciler(t, schemeOpts, clientOpts)
			resultPackages, err := r.packagesToInstall(context.Background(), testMachine)
			assert.NoError(t, err)
			slices.SortFunc(tt.expectedPackages, func(a, b lifecyclev1alpha1.PackageVersion) bool {
				return a.Name < b.Name
			})
			slices.SortFunc(resultPackages, func(a, b lifecyclev1alpha1.PackageVersion) bool {
				return a.Name < b.Name
			})
			assert.Equal(t, tt.expectedPackages, resultPackages)
		})
	}
}

func TestMachineReconciler_MachineLabelsCompliant(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		selector       *metav1.LabelSelector
		labels         map[string]string
		expectedResult bool
	}{
		"exact-match": {
			selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"group": "green",
					"env":   "prod",
				},
			},
			labels: map[string]string{
				"env":   "prod",
				"group": "green",
			},
			expectedResult: true,
		},
		"match": {
			selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"env":   "prod",
					"group": "green",
				},
			},
			labels: map[string]string{
				"region": "eu-west1",
				"env":    "prod",
				"az":     "A",
				"group":  "green",
			},
			expectedResult: true,
		},
		"partial-match": {
			selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"group": "green",
					"env":   "prod",
				},
			},
			labels: map[string]string{
				"group": "green",
			},
			expectedResult: false,
		},
		"no-match": {
			selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"env":   "prod",
					"group": "green",
				},
			},
			labels: map[string]string{
				"region": "eu-west1",
				"az":     "A",
			},
			expectedResult: false,
		},
	}

	for n, tc := range tests {
		name, tt := n, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			labelsFromSelector, err := metav1.LabelSelectorAsMap(tt.selector)
			assert.NoError(t, err)
			result := machineLabelsCompliant(tt.labels, labelsFromSelector)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestMachineReconciler_FilterPackageVersion(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		desiredPackages   []lifecyclev1alpha1.PackageVersion
		installedPackages []lifecyclev1alpha1.PackageVersion
		defaultPackages   []lifecyclev1alpha1.PackageVersion
		expectedPackages  []lifecyclev1alpha1.PackageVersion
	}{
		"no-packages-to-install": {
			desiredPackages: []lifecyclev1alpha1.PackageVersion{
				{Name: "bios", Version: "1.0.0"},
				{Name: "bmc", Version: "1.0.0"},
			},
			installedPackages: []lifecyclev1alpha1.PackageVersion{
				{Name: "bios", Version: "1.0.0"},
				{Name: "bmc", Version: "1.0.0"},
				{Name: "raid", Version: "3.2.1"},
				{Name: "net", Version: "0.1.0"},
			},
			defaultPackages: []lifecyclev1alpha1.PackageVersion{
				{Name: "raid", Version: "3.2.1"},
			},
			expectedPackages: []lifecyclev1alpha1.PackageVersion{},
		},
		"desired-to-update": {
			desiredPackages: []lifecyclev1alpha1.PackageVersion{
				{Name: "bios", Version: "2.0.0"},
				{Name: "bmc", Version: "1.0.0"},
			},
			installedPackages: []lifecyclev1alpha1.PackageVersion{
				{Name: "bios", Version: "1.0.0"},
				{Name: "bmc", Version: "1.0.0"},
				{Name: "raid", Version: "3.2.1"},
				{Name: "net", Version: "0.1.0"},
			},
			defaultPackages: []lifecyclev1alpha1.PackageVersion{
				{Name: "raid", Version: "3.2.1"},
			},
			expectedPackages: []lifecyclev1alpha1.PackageVersion{
				{Name: "bios", Version: "2.0.0"},
			},
		},
		"default-to-update": {
			desiredPackages: []lifecyclev1alpha1.PackageVersion{
				{Name: "bios", Version: "1.0.0"},
				{Name: "bmc", Version: "1.0.0"},
			},
			installedPackages: []lifecyclev1alpha1.PackageVersion{
				{Name: "bios", Version: "1.0.0"},
				{Name: "bmc", Version: "1.0.0"},
				{Name: "raid", Version: "3.2.1"},
				{Name: "net", Version: "0.1.0"},
			},
			defaultPackages: []lifecyclev1alpha1.PackageVersion{
				{Name: "raid", Version: "4.0.0"},
			},
			expectedPackages: []lifecyclev1alpha1.PackageVersion{
				{Name: "raid", Version: "4.0.0"},
			},
		},
		"desired-and-default-to-update": {
			desiredPackages: []lifecyclev1alpha1.PackageVersion{
				{Name: "bios", Version: "1.0.3"},
				{Name: "bmc", Version: "1.7.0"},
			},
			installedPackages: []lifecyclev1alpha1.PackageVersion{
				{Name: "bios", Version: "1.0.0"},
				{Name: "bmc", Version: "1.0.0"},
				{Name: "raid", Version: "3.2.1"},
				{Name: "net", Version: "0.1.0"},
			},
			defaultPackages: []lifecyclev1alpha1.PackageVersion{
				{Name: "raid", Version: "4.0.0"},
				{Name: "net", Version: "1.5.0"},
			},
			expectedPackages: []lifecyclev1alpha1.PackageVersion{
				{Name: "raid", Version: "4.0.0"},
				{Name: "bios", Version: "1.0.3"},
				{Name: "bmc", Version: "1.7.0"},
				{Name: "net", Version: "1.5.0"},
			},
		},
		"desired-overrides-default": {
			desiredPackages: []lifecyclev1alpha1.PackageVersion{
				{Name: "bios", Version: "2.0.0"},
				{Name: "bmc", Version: "1.0.0"},
				{Name: "raid", Version: "5.5.0"},
			},
			installedPackages: []lifecyclev1alpha1.PackageVersion{
				{Name: "bios", Version: "1.0.0"},
				{Name: "bmc", Version: "1.0.0"},
				{Name: "raid", Version: "3.2.1"},
				{Name: "net", Version: "0.1.0"},
			},
			defaultPackages: []lifecyclev1alpha1.PackageVersion{
				{Name: "raid", Version: "3.2.1"},
			},
			expectedPackages: []lifecyclev1alpha1.PackageVersion{
				{Name: "bios", Version: "2.0.0"},
				{Name: "raid", Version: "5.5.0"},
			},
		},
	}

	for n, tc := range tests {
		name, tt := n, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			resultPackages := filterPackageVersion(
				tt.desiredPackages,
				tt.installedPackages,
				tt.defaultPackages)
			slices.SortFunc(tt.expectedPackages, func(a, b lifecyclev1alpha1.PackageVersion) bool {
				return a.Name < b.Name
			})
			slices.SortFunc(resultPackages, func(a, b lifecyclev1alpha1.PackageVersion) bool {
				return a.Name < b.Name
			})
			assert.Equal(t, tt.expectedPackages, resultPackages)
		})
	}
}
