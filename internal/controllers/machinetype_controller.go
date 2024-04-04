// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"
	"time"

	"connectrpc.com/connect"
	"github.com/go-logr/logr"
	"github.com/ironcore-dev/lifecycle-manager/clientgo/connectrpc/machinetype/v1alpha1/machinetypev1alpha1connect"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	lifecyclev1alpha1 "github.com/ironcore-dev/lifecycle-manager/api/lifecycle/v1alpha1"
	machinetypev1alpha1 "github.com/ironcore-dev/lifecycle-manager/api/proto/machinetype/v1alpha1"
)

// MachineTypeReconciler reconciles a MachineType object.
type MachineTypeReconciler struct {
	client.Client
	machinetypev1alpha1connect.MachineTypeServiceClient

	Horizon time.Duration

	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=lifecycle.ironcore.dev,resources=machinetypes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=lifecycle.ironcore.dev,resources=machinetypes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=lifecycle.ironcore.dev,resources=machinetypes/finalizers,verbs=update
// +kubebuilder:rbac:groups=ironcore.dev,resources=oobs,verbs=get;list;watch
// +kubebuilder:rbac:groups=ironcore.dev,resources=oobs/status,verbs=get;list;watch

func (r *MachineTypeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (reconcile.Result, error) {
	log := logr.FromContextOrDiscard(ctx)
	log.V(1).Info("reconciliation started")

	obj := &lifecyclev1alpha1.MachineType{}
	if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	result, err := r.reconcileRequired(ctx, obj)
	if err != nil {
		log.V(1).Info("reconciliation interrupted by an error")
		return result, err
	}
	if err = r.Status().Patch(ctx, obj, client.Merge); err != nil {
		if apierrors.IsConflict(err) {
			return reconcile.Result{RequeueAfter: requeuePeriod}, nil
		}
		log.Error(err, "failed to update object status")
		return reconcile.Result{}, err
	}
	log.V(1).Info("reconciliation finished")
	return result, nil
}

func (r *MachineTypeReconciler) reconcileRequired(
	ctx context.Context,
	obj *lifecyclev1alpha1.MachineType,
) (reconcile.Result, error) {
	log := logr.FromContextOrDiscard(ctx)
	if obj.GetDeletionTimestamp().IsZero() {
		return r.reconcile(ctx, obj)
	}
	log.V(1).Info("object is being deleted")
	return reconcile.Result{}, nil
}

func (r *MachineTypeReconciler) reconcile(
	ctx context.Context,
	obj *lifecyclev1alpha1.MachineType,
) (reconcile.Result, error) {
	if obj.Status.LastScanTime.IsZero() {
		return r.scan(ctx, obj)
	}
	if time.Since(obj.Status.LastScanTime.Time) > r.Horizon {
		return r.scan(ctx, obj)
	}
	obj.Status.Message = ""
	return reconcile.Result{}, nil
}

func (r *MachineTypeReconciler) scan(
	ctx context.Context,
	obj *lifecyclev1alpha1.MachineType,
) (reconcile.Result, error) {
	log := logr.FromContextOrDiscard(ctx)
	resp, err := r.Scan(ctx, connect.NewRequest(&machinetypev1alpha1.ScanRequest{
		Name:      obj.Name,
		Namespace: obj.Namespace,
	}))
	if err != nil {
		log.Error(err, "failed to send scan request")
		return reconcile.Result{}, err
	}
	scanResponse := resp.Msg
	result := reconcile.Result{}
	switch {
	case LCIMRequestResultToString[scanResponse.Result].IsFailure():
		obj.Status.Message = StatusMessageScanRequestFailed
		result.RequeueAfter = requeuePeriod
	case LCIMRequestResultToString[scanResponse.Result].IsScheduled():
		obj.Status.Message = StatusMessageScanRequestProcessing
	case LCIMRequestResultToString[scanResponse.Result].IsSuccess():
		obj.Status.Message = StatusMessageScanRequestSuccessful
	}
	return result, nil
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
