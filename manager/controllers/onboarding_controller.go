// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	oobv1alpha1 "github.com/onmetal/oob-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/tools/reference"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	"github.com/ironcore-dev/lifecycle-manager/api/v1alpha1"
)

// OnboardingReconciler watches for OOB objects and creates
// corresponding MachineType and Machine objects.
type OnboardingReconciler struct {
	client.Client

	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder

	RequeuePeriod time.Duration
	ScanPeriod    metav1.Duration
}

// +kubebuilder:rbac:groups=lifecycle.ironcore.dev,resources=machinetypes,verbs=get;list;create
// +kubebuilder:rbac:groups=lifecycle.ironcore.dev,resources=machines,verbs=get;list;create
// +kubebuilder:rbac:groups=onmetal.de,resources=oobs,verbs=watch;get;list
// +kubebuilder:rbac:groups=onmetal.de,resources=oobs/status,verbs=get

func (r *OnboardingReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var (
		result ctrl.Result
		ref    *corev1.ObjectReference
		err    error
	)

	obj := &oobv1alpha1.OOB{}
	if err = r.Get(ctx, req.NamespacedName, obj); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	ref, err = reference.GetReference(r.Scheme, obj)
	if err != nil {
		r.Log.WithValues("request", req.NamespacedName).Error(err, "failed to construct reference")
		return ctrl.Result{}, err
	}
	log := r.Log.WithValues("object", *ref)
	log.V(1).Info("reconciliation started")

	recCtx := logr.NewContext(ctx, log)
	result, err = r.reconcileRequired(recCtx, obj)
	if err != nil {
		log.V(1).Info("reconciliation interrupted by an error")
	}
	log.V(1).Info("reconciliation finished")
	return result, err
}

func (r *OnboardingReconciler) reconcileRequired(ctx context.Context, obj *oobv1alpha1.OOB) (ctrl.Result, error) {
	log := logr.FromContextOrDiscard(ctx)
	if !obj.GetDeletionTimestamp().IsZero() {
		log.V(2).Info("object is being deleted")
		return ctrl.Result{}, nil
	}
	if obj.Status.Manufacturer == "" || obj.Status.SKU == "" {
		log.V(2).Info("object is not processed yet")
		return ctrl.Result{RequeueAfter: r.RequeuePeriod}, nil
	}
	return r.reconcile(ctx, obj)
}

func (r *OnboardingReconciler) reconcile(ctx context.Context, obj *oobv1alpha1.OOB) (ctrl.Result, error) {
	var err error
	log := logr.FromContextOrDiscard(ctx)
	log.V(2).Info("onboarding machineType object")
	if err = r.onboardMachineType(ctx, obj); err != nil {
		return ctrl.Result{}, err
	}
	log.V(2).Info("onboarding machine object")
	if err = r.onboardMachine(ctx, obj); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (r *OnboardingReconciler) onboardMachineType(ctx context.Context, obj *oobv1alpha1.OOB) error {
	log := logr.FromContextOrDiscard(ctx)
	manufacturer := obj.Status.Manufacturer
	typeName := obj.Status.SKU[:4]
	machineType := &v1alpha1.MachineType{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", manufacturer, typeName),
			Namespace: obj.Namespace,
		},
		Spec: v1alpha1.MachineTypeSpec{
			Manufacturer: manufacturer,
			Type:         typeName,
			ScanPeriod:   r.ScanPeriod,
		},
	}
	if err := r.Create(ctx, machineType); err != nil {
		if apierrors.IsAlreadyExists(err) {
			log.V(2).Info("machineType has been already onboarded")
			return nil
		}
		log.Error(err, "failed to onboard machineType")
		return err
	}
	log.V(1).Info("machineType onboarded successfully")
	return nil
}

func (r *OnboardingReconciler) onboardMachine(ctx context.Context, obj *oobv1alpha1.OOB) error {
	log := logr.FromContextOrDiscard(ctx)
	manufacturer := obj.Status.Manufacturer
	typeName := obj.Status.SKU[:4]
	machine := &v1alpha1.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      obj.Name,
			Namespace: obj.Namespace,
		},
		Spec: v1alpha1.MachineSpec{
			MachineTypeRef: corev1.LocalObjectReference{Name: fmt.Sprintf("%s-%s", manufacturer, typeName)},
			OOBMachineRef:  corev1.LocalObjectReference{Name: obj.Name},
			ScanPeriod:     r.ScanPeriod,
		},
	}
	if err := r.Create(ctx, machine); err != nil {
		if apierrors.IsAlreadyExists(err) {
			log.V(2).Info("machine has been already onboarded")
			return nil
		}
		log.Error(err, "failed to onboard machine")
		return err
	}
	log.V(1).Info("machine onboarded successfully")
	return nil
}

func (r *OnboardingReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&oobv1alpha1.OOB{}).
		WithOptions(controller.Options{
			MaxConcurrentReconciles: 10,
		}).
		Complete(r)
}
