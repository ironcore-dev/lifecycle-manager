// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"

	"github.com/ironcore-dev/lifecycle-manager/util/convertutil"
	"github.com/ironcore-dev/lifecycle-manager/util/testutil/fake"
	"github.com/ironcore-dev/lifecycle-manager/util/testutil/mock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"

	lifecyclev1alpha1 "github.com/ironcore-dev/lifecycle-manager/api/lifecycle/v1alpha1"
	commonv1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/common/v1alpha1"
	machinev1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/machine/v1alpha1"
	"github.com/ironcore-dev/lifecycle-manager/util/apiutil"
	"github.com/ironcore-dev/lifecycle-manager/util/testutil"
	"github.com/ironcore-dev/lifecycle-manager/util/uuidutil"
)

var _ = Describe("Machine controller", func() {
	Context("Reconciliation workflow", func() {
		Context("When Machine object is not found", func() {
			It("Should interrupt reconciliation and return empty result with no error", func() {
				s := testutil.SetupScheme(testutil.WithGroupVersion(lifecyclev1alpha1.AddToScheme))
				c := testutil.SetupClient(s)
				machineRec := NewMachineReconciler(c, s)
				req := ctrl.Request{NamespacedName: types.NamespacedName{Namespace: "default", Name: "sample"}}
				res, err := machineRec.Reconcile(context.Background(), req)
				Expect(err).NotTo(HaveOccurred())
				Expect(res).To(Equal(ctrl.Result{}))
			})
		})

		Context("When Machine object is being deleted", func() {
			It("Should interrupt reconciliation and return empty result with no error", func() {
				s := testutil.SetupScheme(testutil.WithGroupVersion(lifecyclev1alpha1.AddToScheme))
				c := testutil.SetupClient(s, testutil.WithRuntimeObject(mock.NewMachineObject("sample", "default",
					mock.MachineWithDeletionTimestamp(),
					mock.MachineWithFinalizer(),
				)))
				machineRec := NewMachineReconciler(c, s)
				req := ctrl.Request{NamespacedName: types.NamespacedName{Namespace: "default", Name: "sample"}}
				res, err := machineRec.Reconcile(context.Background(), req)
				Expect(err).NotTo(HaveOccurred())
				Expect(res).To(Equal(ctrl.Result{}))
			})
		})

		Context("When referred MachineType object is not found", func() {
			It("Should interrupt reconciliation and return empty result with error", func() {
				now := metav1.Now()
				machineKey := types.NamespacedName{Namespace: "default", Name: "sample"}
				s := testutil.SetupScheme(testutil.WithGroupVersion(lifecyclev1alpha1.AddToScheme))
				c := testutil.SetupClient(s, testutil.WithRuntimeObject(mock.NewMachineObject("sample", "default",
					mock.MachineWithMachineTypeRef("sample"),
				)))
				machineRec := NewMachineReconciler(c, s)
				brokerClient := fake.NewMachineClient(map[string]*machinev1alpha1.MachineStatus{
					uuidutil.UUIDFromObjectKey(machineKey): {
						LastScanTime:      convertutil.TimeToTimestampPtr(now),
						LastScanResult:    commonv1alpha1.ScanResult_SCAN_RESULT_SUCCESS,
						InstalledPackages: nil,
						Message:           "",
					},
				})
				machineRec.MachineServiceClient = brokerClient
				req := ctrl.Request{NamespacedName: machineKey}
				res, err := machineRec.Reconcile(context.Background(), req)
				Expect(err).To(HaveOccurred())
				Expect(res).To(Equal(ctrl.Result{}))
			})
		})

		Context("When new scan request submitted", func() {
			It("Should update Machine object's status with corresponding message", func() {
				machineKey := types.NamespacedName{Namespace: "default", Name: "sample"}
				s := testutil.SetupScheme(testutil.WithGroupVersion(lifecyclev1alpha1.AddToScheme))
				c := testutil.SetupClient(s, testutil.WithRuntimeObject(mock.NewMachineObject("sample", "default")))
				machineRec := NewMachineReconciler(c, s)
				brokerClient := fake.NewMachineClient(map[string]*machinev1alpha1.MachineStatus{})
				machineRec.MachineServiceClient = brokerClient
				req := ctrl.Request{NamespacedName: machineKey}
				res, err := machineRec.Reconcile(context.Background(), req)
				Expect(err).NotTo(HaveOccurred())
				Expect(res).To(Equal(ctrl.Result{}))

				broker, _ := machineRec.MachineServiceClient.(*fake.MachineClient)
				entry := broker.ReadCache(uuidutil.UUIDFromObjectKey(machineKey))
				Expect(entry).NotTo(BeNil())
				reconciledMachine := &lifecyclev1alpha1.Machine{}
				err = machineRec.Get(context.Background(), machineKey, reconciledMachine)
				Expect(err).NotTo(HaveOccurred())
				Expect(reconciledMachine.Status.Message).To(Equal(StatusMessageScanRequestSubmitted))
			})
		})

		Context("When failed to send scan request", func() {
			It("Should interrupt reconciliation and return empty result with error", func() {
				machineKey := types.NamespacedName{Namespace: "default", Name: "failed-scan"}
				s := testutil.SetupScheme(testutil.WithGroupVersion(lifecyclev1alpha1.AddToScheme))
				c := testutil.SetupClient(s, testutil.WithRuntimeObject(mock.NewMachineObject("failed-scan", "default")))
				machineRec := NewMachineReconciler(c, s)
				brokerClient := fake.NewMachineClient(map[string]*machinev1alpha1.MachineStatus{})
				machineRec.MachineServiceClient = brokerClient
				req := ctrl.Request{NamespacedName: machineKey}
				res, err := machineRec.Reconcile(context.Background(), req)
				Expect(err).To(HaveOccurred())
				Expect(res).To(Equal(ctrl.Result{}))
			})
		})

		Context("When installed packages match desired state", func() {
			It("Should update Machine object's status with scan timestamp and result", func() {
				now := metav1.Now()
				expectedPackages := []*commonv1alpha1.PackageVersion{{Name: "bios", Version: "1.0.0"}}
				machine := mock.NewMachineObject("sample", "default",
					mock.MachineWithMachineTypeRef("sample"),
					mock.MachineWithLabels(map[string]string{"env": "test"}),
					mock.MachineStatusWithInstalledPackages(apiutil.PackageVersionsToKubeAPI(expectedPackages)))
				machineType := mock.NewMachineTypeObject("sample", "default",
					mock.MachineTypeWithGroup(lifecyclev1alpha1.MachineGroup{
						MachineSelector: metav1.LabelSelector{MatchLabels: map[string]string{"env": "test"}},
						Packages:        apiutil.PackageVersionsToKubeAPI(expectedPackages)}))

				machineKey := types.NamespacedName{Namespace: "default", Name: "sample"}
				s := testutil.SetupScheme(testutil.WithGroupVersion(lifecyclev1alpha1.AddToScheme))
				c := testutil.SetupClient(s,
					testutil.WithRuntimeObject(machine),
					testutil.WithRuntimeObject(machineType))
				machineRec := NewMachineReconciler(c, s)
				brokerClient := fake.NewMachineClient(map[string]*machinev1alpha1.MachineStatus{
					uuidutil.UUIDFromObjectKey(machineKey): {
						LastScanTime:      convertutil.TimeToTimestampPtr(now),
						LastScanResult:    1,
						InstalledPackages: expectedPackages,
					}})
				machineRec.MachineServiceClient = brokerClient
				req := ctrl.Request{NamespacedName: machineKey}
				res, err := machineRec.Reconcile(context.Background(), req)
				Expect(err).NotTo(HaveOccurred())
				Expect(res).To(Equal(ctrl.Result{}))

				reconciledMachine := &lifecyclev1alpha1.Machine{}
				err = machineRec.Get(context.Background(), machineKey, reconciledMachine)
				Expect(err).NotTo(HaveOccurred())
				Expect(reconciledMachine.Status.InstalledPackages).To(Equal(apiutil.PackageVersionsToKubeAPI(expectedPackages)))
			})
		})

		Context("When packages installation scheduled", func() {
			It("Should update Machine object's status with corresponding message", func() {
				now := metav1.Now()
				desiredPackages := []lifecyclev1alpha1.PackageVersion{{Name: "bios", Version: "1.0.0"}}
				machine := mock.NewMachineObject("sample", "default",
					mock.MachineWithMachineTypeRef("sample"),
					mock.MachineWithLabels(map[string]string{"env": "test"}))
				machineType := mock.NewMachineTypeObject("sample", "default",
					mock.MachineTypeWithGroup(lifecyclev1alpha1.MachineGroup{
						MachineSelector: metav1.LabelSelector{MatchLabels: map[string]string{"env": "test"}},
						Packages:        desiredPackages}))

				machineKey := types.NamespacedName{Namespace: "default", Name: "sample"}
				s := testutil.SetupScheme(testutil.WithGroupVersion(lifecyclev1alpha1.AddToScheme))
				c := testutil.SetupClient(s,
					testutil.WithRuntimeObject(machine),
					testutil.WithRuntimeObject(machineType))
				machineRec := NewMachineReconciler(c, s)
				brokerClient := fake.NewMachineClient(map[string]*machinev1alpha1.MachineStatus{
					uuidutil.UUIDFromObjectKey(machineKey): {
						LastScanTime:   convertutil.TimeToTimestampPtr(now),
						LastScanResult: 1,
					}})
				machineRec.MachineServiceClient = brokerClient
				req := ctrl.Request{NamespacedName: machineKey}
				res, err := machineRec.Reconcile(context.Background(), req)
				Expect(err).NotTo(HaveOccurred())
				Expect(res).To(Equal(ctrl.Result{}))

				reconciledMachine := &lifecyclev1alpha1.Machine{}
				err = machineRec.Get(context.Background(), machineKey, reconciledMachine)
				Expect(err).NotTo(HaveOccurred())
				Expect(reconciledMachine.Status.Message).To(Equal(StatusMessageInstallationScheduled))
			})
		})

		Context("When failed to send install request", func() {
			It("Should interrupt reconciliation and return empty result with error", func() {
				now := metav1.Now()
				desiredPackages := []lifecyclev1alpha1.PackageVersion{{Name: "bios", Version: "1.0.0"}}
				machine := mock.NewMachineObject("failed-install", "default",
					mock.MachineWithMachineTypeRef("sample"),
					mock.MachineWithLabels(map[string]string{"env": "test"}))
				machineType := mock.NewMachineTypeObject("sample", "default",
					mock.MachineTypeWithGroup(lifecyclev1alpha1.MachineGroup{
						MachineSelector: metav1.LabelSelector{MatchLabels: map[string]string{"env": "test"}},
						Packages:        desiredPackages}))

				machineKey := types.NamespacedName{Namespace: "default", Name: "failed-install"}
				s := testutil.SetupScheme(testutil.WithGroupVersion(lifecyclev1alpha1.AddToScheme))
				c := testutil.SetupClient(s,
					testutil.WithRuntimeObject(machine),
					testutil.WithRuntimeObject(machineType))
				machineRec := NewMachineReconciler(c, s)
				brokerClient := fake.NewMachineClient(map[string]*machinev1alpha1.MachineStatus{
					uuidutil.UUIDFromObjectKey(machineKey): {
						LastScanTime:   convertutil.TimeToTimestampPtr(now),
						LastScanResult: 1,
					}})
				machineRec.MachineServiceClient = brokerClient
				req := ctrl.Request{NamespacedName: machineKey}
				res, err := machineRec.Reconcile(context.Background(), req)
				Expect(err).To(HaveOccurred())
				Expect(res).To(Equal(ctrl.Result{}))
			})
		})
	})
})
