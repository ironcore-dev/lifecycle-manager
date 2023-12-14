// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/reference"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	"github.com/ironcore-dev/lifecycle-manager/api/v1alpha1"
	lcmimachinelifecycle "github.com/ironcore-dev/lifecycle-manager/lcmi/api/machine_lifecycle/v1alpha1"
	"github.com/ironcore-dev/lifecycle-manager/util/uuidutil"
)

// MachineLifecycleReconciler reconciles MachineLifecycle objects.
type MachineLifecycleReconciler struct {
	client.Client

	MachineLifecycleBroker lcmimachinelifecycle.MachineLifecycleBrokerServiceClient

	Log     logr.Logger
	Scheme  *runtime.Scheme
	Horizon time.Duration
}

func (r *MachineLifecycleReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var (
		result ctrl.Result
		ref    *corev1.ObjectReference
		err    error
	)

	obj := &v1alpha1.MachineLifecycle{}
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
	if err = r.Status().Patch(ctx, obj, client.Apply, patchOpts); err != nil {
		log.Error(err, "failed to update object status")
		return ctrl.Result{}, err
	}
	log.V(1).Info("reconciliation finished")
	return result, err
}

func (r *MachineLifecycleReconciler) reconcileRequired(ctx context.Context, obj *v1alpha1.MachineLifecycle) (ctrl.Result, error) {
	log := logr.FromContextOrDiscard(ctx)
	if !obj.GetDeletionTimestamp().IsZero() {
		log.V(2).Info("object is being deleted")
		return ctrl.Result{}, nil
	}
	return r.reconcile(ctx, obj)
}

func (r *MachineLifecycleReconciler) reconcile(ctx context.Context, obj *v1alpha1.MachineLifecycle) (ctrl.Result, error) {
	log := logr.FromContextOrDiscard(ctx)
	if r.scanRequired(obj) {
		resp, err := r.MachineLifecycleBroker.Scan(ctx, &lcmimachinelifecycle.ScanRequest{
			Id: uuidutil.UUIDFromObjectKey(client.ObjectKeyFromObject(obj)),
		})
		if err != nil {
			if status.Code(err) == codes.NotFound {
				log.V(1).Info(err.Error())
				return ctrl.Result{}, nil
			}
			log.Error(err, "failed to get scan result")
			return ctrl.Result{}, err
		}
		if state, ok := LCMIScanStateToString[resp.State]; ok && state.IsFinished() {
			obj.Status.LastScanResult = LCMIScanResultToString[resp.Status.LastScanResult]
			obj.Status.LastScanTime = metav1.Time{Time: time.Unix(0, resp.Status.LastScanTime)}
			obj.Status.InstalledPackages = make([]corev1.LocalObjectReference, len(resp.Status.InstalledPackages))
			for i, item := range resp.Status.InstalledPackages {
				obj.Status.InstalledPackages[i] = corev1.LocalObjectReference{Name: item.Name}
			}
		}
	}
	return ctrl.Result{}, nil
}

func (r *MachineLifecycleReconciler) lastScanNotInHorizon(obj *v1alpha1.MachineLifecycle) bool {
	return obj.Status.LastScanTime.IsZero() || time.Since(obj.Status.LastScanTime.Time) > r.Horizon
}

func (r *MachineLifecycleReconciler) scanRequired(obj *v1alpha1.MachineLifecycle) bool {
	if obj.Status.LastScanTime.IsZero() {
		return true
	}
	return r.lastScanNotInHorizon(obj) || obj.Status.LastScanResult == v1alpha1.ScanFailure
}

func (r *MachineLifecycleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.MachineLifecycle{}).
		WithOptions(controller.Options{
			MaxConcurrentReconciles: 10,
		}).
		Complete(r)
}
