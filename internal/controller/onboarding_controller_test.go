// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"context"

	oobv1alpha1 "github.com/onmetal/oob-operator/api/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	lifecyclev1alpha1 "github.com/ironcore-dev/lifecycle-manager/api/lifecycle/v1alpha1"
	"github.com/ironcore-dev/lifecycle-manager/util/testutil"
)

var _ = Describe("Onboarding controller", func() {
	Context("When OOB object is not found", func() {
		It("Should interrupt reconciliation and return empty result with no error", func() {
			s := testutil.SetupScheme(testutil.WithGroupVersion(oobv1alpha1.AddToScheme))
			c := testutil.SetupClient(s)
			onboardingRec := NewOnboardingReconciler(c, s)
			req := ctrl.Request{NamespacedName: types.NamespacedName{Namespace: "default", Name: "sample"}}
			res, err := onboardingRec.Reconcile(context.Background(), req)
			Expect(err).NotTo(HaveOccurred())
			Expect(res).To(Equal(ctrl.Result{}))
		})
	})

	Context("When OOB object's status is empty", func() {
		It("Should requeue OOB object's reconciliation", func() {
			s := testutil.SetupScheme(testutil.WithGroupVersion(oobv1alpha1.AddToScheme))
			c := testutil.SetupClient(s, testutil.WithRuntimeObject(testutil.NewOOBObject("sample", "default")))
			onboardingRec := NewOnboardingReconciler(c, s)
			req := ctrl.Request{NamespacedName: types.NamespacedName{Namespace: "default", Name: "sample"}}
			res, err := onboardingRec.Reconcile(context.Background(), req)
			Expect(err).NotTo(HaveOccurred())
			Expect(res).To(Equal(ctrl.Result{RequeueAfter: onboardingRec.RequeuePeriod}))
		})
	})

	Context("When OOB object is being deleted", func() {
		It("Should interrupt reconciliation and return empty result with no error", func() {
			s := testutil.SetupScheme(testutil.WithGroupVersion(oobv1alpha1.AddToScheme))
			c := testutil.SetupClient(s, testutil.WithRuntimeObject(testutil.NewOOBObject("sample", "default",
				testutil.OOBWithDeletionTimestamp(),
				testutil.OOBWithFinalizer(),
			)))
			onboardingRec := NewOnboardingReconciler(c, s)
			req := ctrl.Request{NamespacedName: types.NamespacedName{Namespace: "default", Name: "sample"}}
			res, err := onboardingRec.Reconcile(context.Background(), req)
			Expect(err).NotTo(HaveOccurred())
			Expect(res).To(Equal(ctrl.Result{}))
		})
	})

	Context("When OOB object is being processed normally", func() {
		It("Should create MachineType object if absent", func() {
			expectedMachineTypeKey := types.NamespacedName{Namespace: "default", Name: "Sample-0000"}
			s := testutil.SetupScheme(
				testutil.WithGroupVersion(oobv1alpha1.AddToScheme),
				testutil.WithGroupVersion(lifecyclev1alpha1.AddToScheme),
			)
			c := testutil.SetupClient(s, testutil.WithRuntimeObject(
				testutil.NewOOBObject("sample", "default", testutil.OOBWithStatus())))
			onboardingRec := NewOnboardingReconciler(c, s)
			onboardedMachineType := &lifecyclev1alpha1.MachineType{}
			err := onboardingRec.Get(context.Background(), expectedMachineTypeKey, onboardedMachineType)
			Expect(err).To(HaveOccurred())
			Expect(apierrors.IsNotFound(err)).To(BeTrue())
			req := ctrl.Request{NamespacedName: types.NamespacedName{Namespace: "default", Name: "sample"}}
			_, err = onboardingRec.Reconcile(context.Background(), req)
			Expect(err).NotTo(HaveOccurred())
			err = onboardingRec.Get(context.Background(), expectedMachineTypeKey, onboardedMachineType)
			Expect(err).NotTo(HaveOccurred())
			Expect(client.ObjectKeyFromObject(onboardedMachineType)).To(Equal(expectedMachineTypeKey))
		})
	})

	Context("When OOB object is being processed normally", func() {
		It("Should create Machine object if absent", func() {
			expectedMachineKey := types.NamespacedName{Namespace: "default", Name: "sample"}
			s := testutil.SetupScheme(
				testutil.WithGroupVersion(oobv1alpha1.AddToScheme),
				testutil.WithGroupVersion(lifecyclev1alpha1.AddToScheme),
			)
			c := testutil.SetupClient(s, testutil.WithRuntimeObject(
				testutil.NewOOBObject("sample", "default", testutil.OOBWithStatus())))
			onboardingRec := NewOnboardingReconciler(c, s)
			onboardedMachine := &lifecyclev1alpha1.Machine{}
			err := onboardingRec.Get(context.Background(), expectedMachineKey, onboardedMachine)
			Expect(err).To(HaveOccurred())
			Expect(apierrors.IsNotFound(err)).To(BeTrue())
			req := ctrl.Request{NamespacedName: types.NamespacedName{Namespace: "default", Name: "sample"}}
			_, err = onboardingRec.Reconcile(context.Background(), req)
			Expect(err).NotTo(HaveOccurred())
			err = onboardingRec.Get(context.Background(), expectedMachineKey, onboardedMachine)
			Expect(err).NotTo(HaveOccurred())
			Expect(client.ObjectKeyFromObject(onboardedMachine)).To(Equal(expectedMachineKey))
		})
	})
})
