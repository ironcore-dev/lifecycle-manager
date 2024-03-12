// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"

	"github.com/ironcore-dev/lifecycle-manager/util/testutil/fake"
	"github.com/ironcore-dev/lifecycle-manager/util/testutil/mock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"

	lifecyclev1alpha1 "github.com/ironcore-dev/lifecycle-manager/api/lifecycle/v1alpha1"
	"github.com/ironcore-dev/lifecycle-manager/util/testutil"
)

var _ = Describe("MachineType controller", func() {
	Context("When MachineType object is not found", func() {
		It("Should interrupt reconciliation and return empty result with no error", func() {
			s := testutil.SetupScheme(testutil.WithGroupVersion(lifecyclev1alpha1.AddToScheme))
			c := testutil.SetupClient(s)
			machinetypeRec := NewMachineTypeReconciler(c, s)
			req := ctrl.Request{NamespacedName: types.NamespacedName{Namespace: "default", Name: "sample"}}
			res, err := machinetypeRec.Reconcile(context.Background(), req)
			Expect(err).NotTo(HaveOccurred())
			Expect(res).To(Equal(ctrl.Result{}))
		})
	})

	Context("When MachineType object is being deleted", func() {
		It("Should interrupt reconciliation and return empty result with no error", func() {
			now := metav1.Now()
			machineType := mock.NewUnstructuredBuilder().
				WithName("sample").
				WithNamespace("default").
				WithDeletionTimestamp(&now).
				WithFinalizers([]string{"test-suite-finalizer"}).
				MachineTypeFromUnstructured().Complete()
			Expect(machineType).NotTo(BeNil())
			s := testutil.SetupScheme(testutil.WithGroupVersion(lifecyclev1alpha1.AddToScheme))
			c := testutil.SetupClient(s, testutil.WithRuntimeObject(machineType))
			machinetypeRec := NewMachineTypeReconciler(c, s)
			req := ctrl.Request{NamespacedName: types.NamespacedName{Namespace: "default", Name: "sample"}}
			res, err := machinetypeRec.Reconcile(context.Background(), req)
			Expect(err).NotTo(HaveOccurred())
			Expect(res).To(Equal(ctrl.Result{}))
		})
	})

	Context("When new scan request submitted", func() {
		It("Should update MachineType object's status with corresponding message", func() {
			machinetypeKey := types.NamespacedName{Namespace: "default", Name: "sample-scan-submitted"}
			machineType := mock.NewUnstructuredBuilder().
				WithName("sample-scan-submitted").
				WithNamespace("default").
				MachineTypeFromUnstructured().Complete()
			Expect(machineType).NotTo(BeNil())
			s := testutil.SetupScheme(testutil.WithGroupVersion(lifecyclev1alpha1.AddToScheme))
			c := testutil.SetupClient(s, testutil.WithRuntimeObject(machineType))
			machinetypeRec := NewMachineTypeReconciler(c, s)
			brokerClient := fake.NewMachineTypeClient()
			machinetypeRec.MachineTypeServiceClient = brokerClient
			req := ctrl.Request{NamespacedName: machinetypeKey}
			res, err := machinetypeRec.Reconcile(context.Background(), req)
			Expect(err).NotTo(HaveOccurred())
			Expect(res).To(Equal(ctrl.Result{}))
			reconciledMachineType := &lifecyclev1alpha1.MachineType{}
			err = machinetypeRec.Get(context.Background(), machinetypeKey, reconciledMachineType)
			Expect(err).NotTo(HaveOccurred())
			Expect(reconciledMachineType.Status.Message).To(Equal(StatusMessageScanRequestSuccessful))
		})
	})

	Context("When failed to send scan request", func() {
		It("Should interrupt reconciliation and return empty result with error", func() {
			machinetypeKey := types.NamespacedName{Namespace: "default", Name: "failed-scan"}
			machineType := mock.NewUnstructuredBuilder().
				WithName("failed-scan").
				WithNamespace("default").
				MachineTypeFromUnstructured().Complete()
			Expect(machineType).NotTo(BeNil())
			s := testutil.SetupScheme(testutil.WithGroupVersion(lifecyclev1alpha1.AddToScheme))
			c := testutil.SetupClient(s, testutil.WithRuntimeObject(machineType))
			machinetypeRec := NewMachineTypeReconciler(c, s)
			brokerClient := fake.NewMachineTypeClient()
			machinetypeRec.MachineTypeServiceClient = brokerClient
			req := ctrl.Request{NamespacedName: machinetypeKey}
			res, err := machinetypeRec.Reconcile(context.Background(), req)
			Expect(err).To(HaveOccurred())
			Expect(res).To(Equal(ctrl.Result{}))
		})
	})
})
