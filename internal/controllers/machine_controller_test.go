// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"

	"github.com/ironcore-dev/lifecycle-manager/util/testutil/fake"
	"github.com/ironcore-dev/lifecycle-manager/util/testutil/mock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"

	lifecyclev1alpha1 "github.com/ironcore-dev/lifecycle-manager/api/lifecycle/v1alpha1"
	"github.com/ironcore-dev/lifecycle-manager/util/testutil"
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
				now := metav1.Now()
				machine := mock.NewUnstructuredBuilder().
					WithName("sample").
					WithNamespace("default").
					WithDeletionTimestamp(&now).
					WithFinalizers([]string{"test-suite-finalizer"}).
					MachineFromUnstructured().Complete()
				Expect(machine).NotTo(BeNil())
				s := testutil.SetupScheme(testutil.WithGroupVersion(lifecyclev1alpha1.AddToScheme))
				c := testutil.SetupClient(s, testutil.WithRuntimeObject(machine))
				machineRec := NewMachineReconciler(c, s)
				req := ctrl.Request{NamespacedName: types.NamespacedName{Namespace: "default", Name: "sample"}}
				res, err := machineRec.Reconcile(context.Background(), req)
				Expect(err).NotTo(HaveOccurred())
				Expect(res).To(Equal(ctrl.Result{}))
			})
		})

		Context("When referred MachineType object is not found", func() {
			It("Should interrupt reconciliation and return an error", func() {
				machine := mock.NewUnstructuredBuilder().
					WithName("sample").
					WithNamespace("default").
					MachineFromUnstructured().WithMachineTypeRef("sample").
					WithDesiredPackages(lifecyclev1alpha1.PackageVersion{Name: "bios", Version: "0.0.1"}).
					WithLastScanTime(metav1.Now()).
					Complete()
				Expect(machine).NotTo(BeNil())
				machineKey := types.NamespacedName{Namespace: "default", Name: "sample"}
				s := testutil.SetupScheme(testutil.WithGroupVersion(lifecyclev1alpha1.AddToScheme))
				c := testutil.SetupClient(s, testutil.WithRuntimeObject(machine))
				machineRec := NewMachineReconciler(c, s)
				brokerClient := fake.NewMachineClient()
				machineRec.MachineServiceClient = brokerClient
				req := ctrl.Request{NamespacedName: machineKey}
				_, err := machineRec.Reconcile(context.Background(), req)
				Expect(err).To(HaveOccurred())
				Expect(apierrors.IsNotFound(err)).To(BeTrue())
			})
		})

		Context("When new scan request submitted", func() {
			It("Should update Machine object's status with corresponding message", func() {
				machineKey := types.NamespacedName{Namespace: "default", Name: "sample-scan-submitted"}
				machine := mock.NewUnstructuredBuilder().
					WithName("sample-scan-submitted").
					WithNamespace("default").
					WithLabels(map[string]string{"env": "test"}).
					MachineFromUnstructured().WithMachineTypeRef("sample").Complete()
				Expect(machine).NotTo(BeNil())
				machineType := mock.NewUnstructuredBuilder().
					WithName("sample").
					WithNamespace("default").
					MachineTypeFromUnstructured().
					Complete()
				Expect(machineType).NotTo(BeNil())
				s := testutil.SetupScheme(testutil.WithGroupVersion(lifecyclev1alpha1.AddToScheme))
				c := testutil.SetupClient(s,
					testutil.WithRuntimeObject(machine),
					testutil.WithRuntimeObject(machineType))
				machineRec := NewMachineReconciler(c, s)
				brokerClient := fake.NewMachineClient()
				machineRec.MachineServiceClient = brokerClient
				req := ctrl.Request{NamespacedName: machineKey}
				res, err := machineRec.Reconcile(context.Background(), req)
				Expect(err).NotTo(HaveOccurred())
				Expect(res).To(Equal(ctrl.Result{}))
				reconciledMachine := &lifecyclev1alpha1.Machine{}
				err = machineRec.Get(context.Background(), machineKey, reconciledMachine)
				Expect(err).NotTo(HaveOccurred())
				Expect(reconciledMachine.Status.Message).To(Equal(StatusMessageScanRequestSuccessful))
			})
		})

		Context("When failed to send scan request", func() {
			It("Should interrupt reconciliation and return empty result with error", func() {
				machineKey := types.NamespacedName{Namespace: "default", Name: "failed-scan"}
				machine := mock.NewUnstructuredBuilder().
					WithName("failed-scan").
					WithNamespace("default").
					MachineFromUnstructured().Complete()
				Expect(machine).NotTo(BeNil())
				s := testutil.SetupScheme(testutil.WithGroupVersion(lifecyclev1alpha1.AddToScheme))
				c := testutil.SetupClient(s, testutil.WithRuntimeObject(machine))
				machineRec := NewMachineReconciler(c, s)
				brokerClient := fake.NewMachineClient()
				machineRec.MachineServiceClient = brokerClient
				req := ctrl.Request{NamespacedName: machineKey}
				res, err := machineRec.Reconcile(context.Background(), req)
				Expect(err).To(HaveOccurred())
				Expect(res).To(Equal(ctrl.Result{}))
			})
		})

		Context("When packages installation scheduled", func() {
			It("Should update Machine object's status with corresponding message", func() {
				desiredPackages := []lifecyclev1alpha1.PackageVersion{{Name: "bios", Version: "1.0.0"}}
				machine := mock.NewUnstructuredBuilder().
					WithName("sample-install-submitted").
					WithNamespace("default").
					WithLabels(map[string]string{"env": "test"}).
					MachineFromUnstructured().WithMachineTypeRef("sample").
					WithLastScanTime(metav1.Now()).
					Complete()
				Expect(machine).NotTo(BeNil())
				machineType := mock.NewUnstructuredBuilder().
					WithName("sample").
					WithNamespace("default").
					MachineTypeFromUnstructured().
					WithMachineGroups([]lifecyclev1alpha1.MachineGroup{{
						MachineSelector: metav1.LabelSelector{MatchLabels: map[string]string{"env": "test"}},
						Packages:        desiredPackages}}).
					Complete()
				Expect(machineType).NotTo(BeNil())
				machineKey := types.NamespacedName{Namespace: "default", Name: "sample-install-submitted"}
				s := testutil.SetupScheme(testutil.WithGroupVersion(lifecyclev1alpha1.AddToScheme))
				c := testutil.SetupClient(s,
					testutil.WithRuntimeObject(machine),
					testutil.WithRuntimeObject(machineType))
				machineRec := NewMachineReconciler(c, s)
				brokerClient := fake.NewMachineClient()
				machineRec.MachineServiceClient = brokerClient
				req := ctrl.Request{NamespacedName: machineKey}
				res, err := machineRec.Reconcile(context.Background(), req)
				Expect(err).NotTo(HaveOccurred())
				Expect(res).To(Equal(ctrl.Result{}))

				reconciledMachine := &lifecyclev1alpha1.Machine{}
				err = machineRec.Get(context.Background(), machineKey, reconciledMachine)
				Expect(err).NotTo(HaveOccurred())
				Expect(reconciledMachine.Status.Message).To(Equal(StatusMessageInstallRequestProcessing))
			})
		})

		Context("When failed to send install request", func() {
			It("Should interrupt reconciliation and return empty result with error", func() {
				desiredPackages := []lifecyclev1alpha1.PackageVersion{{Name: "bios", Version: "1.0.0"}}
				machine := mock.NewUnstructuredBuilder().
					WithName("failed-install").
					WithNamespace("default").
					WithLabels(map[string]string{"env": "test"}).
					MachineFromUnstructured().WithMachineTypeRef("sample").
					WithLastScanTime(metav1.Now()).
					Complete()
				Expect(machine).NotTo(BeNil())
				machineType := mock.NewUnstructuredBuilder().
					WithName("sample").WithNamespace("default").MachineTypeFromUnstructured().
					WithMachineGroups([]lifecyclev1alpha1.MachineGroup{{
						MachineSelector: metav1.LabelSelector{MatchLabels: map[string]string{"env": "test"}},
						Packages:        desiredPackages}}).
					Complete()
				Expect(machineType).NotTo(BeNil())
				machineKey := types.NamespacedName{Namespace: "default", Name: "failed-install"}
				s := testutil.SetupScheme(testutil.WithGroupVersion(lifecyclev1alpha1.AddToScheme))
				c := testutil.SetupClient(s,
					testutil.WithRuntimeObject(machine),
					testutil.WithRuntimeObject(machineType))
				machineRec := NewMachineReconciler(c, s)
				brokerClient := fake.NewMachineClient()
				machineRec.MachineServiceClient = brokerClient
				req := ctrl.Request{NamespacedName: machineKey}
				res, err := machineRec.Reconcile(context.Background(), req)
				Expect(err).To(HaveOccurred())
				Expect(res).To(Equal(ctrl.Result{}))
			})
		})
	})
})
