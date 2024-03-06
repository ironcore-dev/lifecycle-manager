// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"

	"connectrpc.com/connect"
	"github.com/go-logr/logr"
	"github.com/ironcore-dev/lifecycle-manager/lcmi/api/machinetype/v1alpha1/machinetypev1alpha1connect"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/reference"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	lifecyclev1alpha1 "github.com/ironcore-dev/lifecycle-manager/api/lifecycle/v1alpha1"
	machinetypev1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/machinetype/v1alpha1"
)

// MachineTypeReconciler reconciles a MachineType object.
type MachineTypeReconciler struct {
	client.Client
	machinetypev1alpha1connect.MachineTypeServiceClient

	Log    logr.Logger
	Scheme *runtime.Scheme
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

	recCtx := logr.NewContext(ctx, log)
	result, err = r.reconcileRequired(recCtx, obj)
	if err != nil {
		log.V(1).Info("reconciliation interrupted by an error")
		return result, err
	}
	if err = r.Status().Update(ctx, obj); err != nil {
		if apierrors.IsConflict(err) {
			return ctrl.Result{RequeueAfter: requeuePeriod}, nil
		}
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
	log := logr.FromContextOrDiscard(ctx)
	resp, err := r.Scan(ctx, connect.NewRequest(&machinetypev1alpha1.ScanRequest{
		Name:      obj.Name,
		Namespace: obj.Namespace,
	}))
	scanResponse := resp.Msg
	if err != nil {
		log.Error(err, "failed to send scan request")
		return ctrl.Result{}, err
	}
	if LCIMRequestResultToString[scanResponse.Result].IsScheduled() {
		obj.Status.Message = StatusMessageScanRequestSubmitted
		return ctrl.Result{}, nil
	}
	if LCIMRequestResultToString[scanResponse.Result].IsFailure() {
		obj.Status.Message = StatusMessageScanRequestFailed
		return ctrl.Result{}, nil
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MachineTypeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&lifecyclev1alpha1.MachineType{}).
		WithOptions(controller.Options{
			MaxConcurrentReconciles: 10,
		}).
		Named("machinetype").
		Complete(r)
}
