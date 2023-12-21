// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/reference"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	lifecyclev1alpha1 "github.com/ironcore-dev/lifecycle-manager/api/lifecycle/v1alpha1"
)

// MachineTypeReconciler reconciles a MachineType object.
type MachineTypeReconciler struct {
	client.Client
	Log     logr.Logger
	Scheme  *runtime.Scheme
	Horizon time.Duration
}

// +kubebuilder:rbac:groups=lifecycle.ironcore.dev,resources=machinetypes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=lifecycle.ironcore.dev,resources=machinetypes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=lifecycle.ironcore.dev,resources=machinetypes/finalizers,verbs=update

func (r *MachineTypeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var (
		result ctrl.Result
		ref    *corev1.ObjectReference
		err    error
	)

	obj := &lifecyclev1alpha1.MachineType{}
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

	base := obj.DeepCopy()
	recCtx := logr.NewContext(ctx, log)
	result, err = r.reconcileRequired(recCtx, obj)
	if err != nil {
		log.V(1).Info("reconciliation interrupted by an error")
		return result, err
	}
	if err = r.Status().Patch(ctx, obj, client.MergeFrom(base)); err != nil {
		log.Error(err, "failed to update object status")
		return ctrl.Result{}, err
	}
	log.V(1).Info("reconciliation finished")
	return result, err
}

func (r *MachineTypeReconciler) reconcileRequired(
	ctx context.Context,
	obj *lifecyclev1alpha1.MachineType,
) (ctrl.Result, error) {
	log := logr.FromContextOrDiscard(ctx)
	if obj.GetDeletionTimestamp().IsZero() {
		return r.reconcile(ctx, obj)
	}
	log.V(2).Info("object is being deleted")
	return ctrl.Result{}, nil
}

func (r *MachineTypeReconciler) reconcile(
	ctx context.Context,
	obj *lifecyclev1alpha1.MachineType,
) (ctrl.Result, error) {
	// todo: implement me
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MachineTypeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&lifecyclev1alpha1.MachineType{}).
		WithOptions(controller.Options{
			MaxConcurrentReconciles: 10,
		}).
		Complete(r)
}
