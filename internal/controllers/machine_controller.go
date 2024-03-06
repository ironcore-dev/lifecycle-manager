// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"
	"reflect"
	"slices"

	"connectrpc.com/connect"
	"github.com/go-logr/logr"
	"github.com/ironcore-dev/lifecycle-manager/lcmi/api/machine/v1alpha1/machinev1alpha1connect"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

// MachineReconciler reconciles a Machine object.
type MachineReconciler struct {
	client.Client
	machinev1alpha1connect.MachineServiceClient

	Namespace string

	Log    logr.Logger
	Scheme *runtime.Scheme
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

func (r *MachineReconciler) reconcileRequired(
	ctx context.Context,
	obj *lifecyclev1alpha1.Machine,
) (ctrl.Result, error) {
	log := logr.FromContextOrDiscard(ctx)
	if obj.GetDeletionTimestamp().IsZero() {
		return r.reconcileScan(ctx, obj)
	}
	log.V(2).Info("object is being deleted")
	return ctrl.Result{}, nil
}

func (r *MachineReconciler) reconcileScan(ctx context.Context, obj *lifecyclev1alpha1.Machine) (ctrl.Result, error) {
	log := logr.FromContextOrDiscard(ctx)
	resp, err := r.ScanMachine(ctx, connect.NewRequest(&machinev1alpha1.ScanMachineRequest{
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
	return r.reconcileInstall(ctx, obj)
}

func (r *MachineReconciler) reconcileInstall(ctx context.Context, obj *lifecyclev1alpha1.Machine) (ctrl.Result, error) {
	log := logr.FromContextOrDiscard(ctx)
	packagesToInstall, err := r.packagesToInstall(ctx, obj)
	if err != nil {
		return ctrl.Result{}, err
	}
	if len(packagesToInstall) == 0 {
		log.V(1).Info("install packages versions match desired state")
		return ctrl.Result{}, nil
	}
	resp, err := r.Install(ctx, connect.NewRequest(&machinev1alpha1.InstallRequest{
		Name:      obj.Name,
		Namespace: obj.Namespace,
	}))
	installResponse := resp.Msg
	if err != nil {
		log.Error(err, "failed to send install request")
		return ctrl.Result{}, err
	}
	if LCIMRequestResultToString[installResponse.Result].IsScheduled() {
		obj.Status.Message = StatusMessageInstallationScheduled
		return ctrl.Result{}, nil
	}
	if LCIMRequestResultToString[installResponse.Result].IsFailure() {
		log.V(1).Info(StatusMessageInstallRequestFailed)
		obj.Status.Message = StatusMessageInstallRequestFailed
	}
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
		Named("machine").
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
	if oldMachineType.Namespace != r.Namespace || newMachineType.Namespace != r.Namespace {
		return
	}
	if reflect.DeepEqual(oldMachineType.Spec.MachineGroups, newMachineType.Spec.MachineGroups) {
		return
	}
	for _, group := range newMachineType.Spec.MachineGroups {
		ls := group.MachineSelector
		selector, err := metav1.LabelSelectorAsSelector(&ls)
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

func (r *MachineReconciler) packagesToInstall(
	ctx context.Context,
	obj *lifecyclev1alpha1.Machine,
) ([]lifecyclev1alpha1.PackageVersion, error) {
	log := logr.FromContextOrDiscard(ctx)
	machineType := &lifecyclev1alpha1.MachineType{}
	key := types.NamespacedName{Name: obj.Spec.MachineTypeRef.Name, Namespace: obj.Namespace}
	if err := r.Get(ctx, key, machineType); err != nil {
		log.Error(err, "failed to get referenced machine type object")
		return nil, err
	}

	defaultPackages := make([]lifecyclev1alpha1.PackageVersion, 0)
	for _, entry := range machineType.Spec.MachineGroups {
		selector := entry.MachineSelector.DeepCopy()
		groupLabels, err := metav1.LabelSelectorAsMap(selector)
		if err != nil {
			log.Error(err, "failed to convert label selector to map")
			return nil, err
		}
		if machineLabelsCompliant(obj.Labels, groupLabels) {
			defaultPackages = append(defaultPackages, entry.Packages...)
			break
		}
	}

	resultPackages := filterPackageVersion(obj.Spec.Packages, obj.Status.InstalledPackages, defaultPackages)
	return resultPackages, nil
}

func machineLabelsCompliant(machineLabels, selectorLabels map[string]string) bool {
	for k, sv := range selectorLabels {
		mv, ok := machineLabels[k]
		if !ok {
			return false
		}
		if mv != sv {
			return false
		}
	}
	return true
}

func filterPackageVersion(
	desiredPackages, installedPackages, defaultPackages []lifecyclev1alpha1.PackageVersion,
) []lifecyclev1alpha1.PackageVersion {
	tempPackages := make([]lifecyclev1alpha1.PackageVersion, 0)
	resultPackages := make([]lifecyclev1alpha1.PackageVersion, 0)
	tempPackages = append(tempPackages, desiredPackages...)

	for _, pv := range defaultPackages {
		idx := slices.IndexFunc(desiredPackages, func(packageVersion lifecyclev1alpha1.PackageVersion) bool {
			return pv.Name == packageVersion.Name
		})
		if idx < 0 {
			tempPackages = append(tempPackages, pv)
			continue
		}
	}

	for _, pv := range tempPackages {
		idx := slices.IndexFunc(installedPackages, func(packageVersion lifecyclev1alpha1.PackageVersion) bool {
			return pv.Name == packageVersion.Name
		})
		if idx < 0 {
			resultPackages = append(resultPackages, pv)
			continue
		}
		if installedPackages[idx].Version != pv.Version {
			resultPackages = append(resultPackages, pv)
		}
	}
	return resultPackages
}
