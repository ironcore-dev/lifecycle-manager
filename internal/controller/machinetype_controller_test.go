// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/ironcore-dev/lifecycle-manager/util/apiutil"
	"github.com/ironcore-dev/lifecycle-manager/util/testutil"
	"github.com/ironcore-dev/lifecycle-manager/util/uuidutil"

	lifecyclev1alpha1 "github.com/ironcore-dev/lifecycle-manager/api/lifecycle/v1alpha1"
	machinetypev1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/machinetype/v1alpha1"
	"github.com/ironcore-dev/lifecycle-manager/lcmi/fake"
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
			s := testutil.SetupScheme(testutil.WithGroupVersion(lifecyclev1alpha1.AddToScheme))
			c := testutil.SetupClient(s, testutil.WithRuntimeObject(testutil.NewMachineTypeObject("sample", "default",
				testutil.MachineTypeWithDeletionTimestamp(),
				testutil.MachineTypeWithFinalizer(),
			)))
			machinetypeRec := NewMachineTypeReconciler(c, s)
			req := ctrl.Request{NamespacedName: types.NamespacedName{Namespace: "default", Name: "sample"}}
			res, err := machinetypeRec.Reconcile(context.Background(), req)
			Expect(err).NotTo(HaveOccurred())
			Expect(res).To(Equal(ctrl.Result{}))
		})
	})

	Context("When new scan request submitted", func() {
		It("Should update MachineType object's status with corresponding message", func() {
			machinetypeKey := types.NamespacedName{Namespace: "default", Name: "sample"}
			s := testutil.SetupScheme(testutil.WithGroupVersion(lifecyclev1alpha1.AddToScheme))
			c := testutil.SetupClient(s, testutil.WithRuntimeObject(testutil.NewMachineTypeObject("sample", "default")))
			machinetypeRec := NewMachineTypeReconciler(c, s)
			brokerClient := fake.NewMachineTypeClient(map[string]*machinetypev1alpha1.MachineTypeStatus{})
			machinetypeRec.Broker = brokerClient
			req := ctrl.Request{NamespacedName: machinetypeKey}
			res, err := machinetypeRec.Reconcile(context.Background(), req)
			Expect(err).NotTo(HaveOccurred())
			Expect(res).To(Equal(ctrl.Result{}))
			broker, _ := machinetypeRec.Broker.(*fake.MachineTypeClient)
			entry := broker.ReadCache(uuidutil.UUIDFromObjectKey(machinetypeKey))
			Expect(entry).NotTo(BeNil())
			reconciledMachineType := &lifecyclev1alpha1.MachineType{}
			err = machinetypeRec.Get(context.Background(), machinetypeKey, reconciledMachineType)
			Expect(err).NotTo(HaveOccurred())
			Expect(reconciledMachineType.Status.Message).To(Equal(StatusMessageScanRequestSubmitted))
		})
	})

	Context("When failed to send scan request", func() {
		It("Should interrupt reconciliation and return empty result with error", func() {
			machinetypeKey := types.NamespacedName{Namespace: "default", Name: "failed-scan"}
			s := testutil.SetupScheme(testutil.WithGroupVersion(lifecyclev1alpha1.AddToScheme))
			c := testutil.SetupClient(s, testutil.WithRuntimeObject(testutil.NewMachineTypeObject("failed-scan", "default")))
			machinetypeRec := NewMachineTypeReconciler(c, s)
			brokerClient := fake.NewMachineTypeClient(map[string]*machinetypev1alpha1.MachineTypeStatus{})
			machinetypeRec.Broker = brokerClient
			req := ctrl.Request{NamespacedName: machinetypeKey}
			res, err := machinetypeRec.Reconcile(context.Background(), req)
			Expect(err).To(HaveOccurred())
			Expect(res).To(Equal(ctrl.Result{}))
		})
	})

	Context("When scan response received", func() {
		It("Should update MachineType object's status with scan timestamp and result", func() {
			now := time.Unix(time.Now().Unix(), 0)
			machinetypeKey := types.NamespacedName{Namespace: "default", Name: "sample"}
			s := testutil.SetupScheme(testutil.WithGroupVersion(lifecyclev1alpha1.AddToScheme))
			c := testutil.SetupClient(s, testutil.WithRuntimeObject(testutil.NewMachineTypeObject("sample", "default")))
			machinetypeRec := NewMachineTypeReconciler(c, s)
			brokerClient := fake.NewMachineTypeClient(map[string]*machinetypev1alpha1.MachineTypeStatus{
				uuidutil.UUIDFromObjectKey(machinetypeKey): {
					LastScanTime:   &metav1.Time{Time: now},
					LastScanResult: 1,
					AvailablePackages: []*machinetypev1alpha1.AvailablePackageVersions{
						{Name: "raid", Versions: []string{"2.9.0", "3.0.0", "3.2.1"}},
					},
				},
			})
			machinetypeRec.Broker = brokerClient
			req := ctrl.Request{NamespacedName: machinetypeKey}
			res, err := machinetypeRec.Reconcile(context.Background(), req)
			Expect(err).NotTo(HaveOccurred())
			Expect(res).To(Equal(ctrl.Result{}))
			broker, _ := machinetypeRec.Broker.(*fake.MachineTypeClient)
			entry := broker.ReadCache(uuidutil.UUIDFromObjectKey(machinetypeKey))
			Expect(entry).NotTo(BeNil())
			reconciledMachineType := &lifecyclev1alpha1.MachineType{}
			err = machinetypeRec.Get(context.Background(), machinetypeKey, reconciledMachineType)
			Expect(err).NotTo(HaveOccurred())
			Expect(reconciledMachineType.Status).To(Equal(apiutil.MachineTypeStatusFrom(entry)))
		})
	})
})
