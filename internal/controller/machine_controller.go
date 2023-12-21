// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"context"
	"reflect"
	"time"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/reference"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	lifecyclev1alpha1 "github.com/ironcore-dev/lifecycle-manager/api/lifecycle/v1alpha1"
	machinev1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/machine/v1alpha1"
)

const (
	StatusMessageScanPending = "scan job pending"
)

// MachineReconciler reconciles a Machine object.
type MachineReconciler struct {
	client.Client
	machinev1alpha1.MachineServiceClient

	Namespace string

	Log     logr.Logger
	Scheme  *runtime.Scheme
	Horizon time.Duration
}

// +kubebuilder:rbac:groups=lifecycle.ironcore.dev,resources=machines,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=lifecycle.ironcore.dev,resources=machines/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=lifecycle.ironcore.dev,resources=machines/finalizers,verbs=update

func (r *MachineReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var (
		result ctrl.Result
		ref    *corev1.ObjectReference
		err    error
	)

	obj := &lifecyclev1alpha1.Machine{}
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

func (r *MachineReconciler) reconcileRequired(
	ctx context.Context,
	obj *lifecyclev1alpha1.Machine,
) (ctrl.Result, error) {
	log := logr.FromContextOrDiscard(ctx)
	if obj.GetDeletionTimestamp().IsZero() {
		return r.reconcile(ctx, obj)
	}
	log.V(2).Info("object is being deleted")
	return ctrl.Result{}, nil
}

func (r *MachineReconciler) reconcile(ctx context.Context, obj *lifecyclev1alpha1.Machine) (ctrl.Result, error) {
	log := logr.FromContextOrDiscard(ctx)
	resp, err := r.Scan(ctx, &machinev1alpha1.ScanRequest{Name: obj.Name, Namespace: obj.Namespace})
	if err != nil {
		log.Error(err, "failed to get scan results")
		return ctrl.Result{}, err
	}
	if resp.Response == nil {
		obj.Status.Message = StatusMessageScanPending
		return ctrl.Result{}, nil
	}
	obj.Status = *resp.Response
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MachineReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&lifecyclev1alpha1.Machine{}).
		WithOptions(controller.Options{
			MaxConcurrentReconciles: 10,
		}).
		Watches(&lifecyclev1alpha1.MachineType{}, handler.Funcs{
			UpdateFunc: r.enqueueOnMachineTypeUpdate,
		}).
		Complete(r)
}

func (r *MachineReconciler) enqueueOnMachineTypeUpdate(
	ctx context.Context,
	e event.UpdateEvent,
	q workqueue.RateLimitingInterface,
) {
	oldMachineType, castOldOk := e.ObjectOld.(*lifecyclev1alpha1.MachineType)
	newMachineType, castNewOk := e.ObjectNew.(*lifecyclev1alpha1.MachineType)
	if !castOldOk || !castNewOk {
		return
	}
	if reflect.DeepEqual(oldMachineType.Spec.MachineGroups, newMachineType.Spec.MachineGroups) {
		return
	}
	for _, group := range newMachineType.Spec.MachineGroups {
		ls := group.MachineSelector
		selector, err := v1.LabelSelectorAsSelector(&ls)
		if err != nil {
			r.Log.Error(err, "failed to process machinetype update")
			return
		}
		opts := &client.ListOptions{
			LabelSelector: selector,
			Namespace:     r.Namespace,
			Limit:         1000,
		}
		machines := &lifecyclev1alpha1.MachineList{}
		if err = r.List(ctx, machines, opts); err != nil {
			r.Log.Error(err, "failed to list machines")
			return
		}
		for _, item := range machines.Items {
			q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
				Namespace: item.Namespace,
				Name:      item.Name,
			}})
		}
	}
}
