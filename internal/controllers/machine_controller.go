// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"
	"reflect"
	"slices"
	"time"

	"connectrpc.com/connect"
	"github.com/go-logr/logr"
	"github.com/ironcore-dev/lifecycle-manager/clientgo/connectrpc/machine/v1alpha1/machinev1alpha1connect"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	lifecyclev1alpha1 "github.com/ironcore-dev/lifecycle-manager/api/lifecycle/v1alpha1"
	machinev1alpha1 "github.com/ironcore-dev/lifecycle-manager/api/proto/machine/v1alpha1"
)

// MachineReconciler reconciles a Machine object.
type MachineReconciler struct {
	client.Client
	machinev1alpha1connect.MachineServiceClient

	Namespace string
	Horizon   time.Duration

	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=lifecycle.ironcore.dev,resources=machines,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=lifecycle.ironcore.dev,resources=machines/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=lifecycle.ironcore.dev,resources=machines/finalizers,verbs=update

func (r *MachineReconciler) Reconcile(ctx context.Context, req ctrl.Request) (reconcile.Result, error) {
	log := logr.FromContextOrDiscard(ctx)
	log.V(1).Info("reconciliation started")

	obj := &lifecyclev1alpha1.Machine{}
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
	return result, err
}

func (r *MachineReconciler) reconcileRequired(
	ctx context.Context,
	obj *lifecyclev1alpha1.Machine,
) (reconcile.Result, error) {
	log := logr.FromContextOrDiscard(ctx)
	if obj.GetDeletionTimestamp().IsZero() {
		return r.reconcile(ctx, obj)
	}
	log.V(1).Info("object is being deleted")
	return reconcile.Result{}, nil
}

func (r *MachineReconciler) reconcile(ctx context.Context, obj *lifecyclev1alpha1.Machine) (reconcile.Result, error) {
	if obj.Status.LastScanTime.IsZero() {
		return r.scan(ctx, obj)
	}
	if time.Since(obj.Status.LastScanTime.Time) > r.Horizon {
		return r.scan(ctx, obj)
	}
	diff, err := r.packagesToInstall(ctx, obj)
	if err != nil {
		return reconcile.Result{}, err
	}
	if len(diff) == 0 {
		obj.Status.Message = ""
		return reconcile.Result{}, nil
	}
	obj.Spec.Packages = mergePackageVersion(obj.Spec.Packages, diff)
	if err = r.Patch(ctx, obj, client.Merge); err != nil {
		return reconcile.Result{}, err
	}
	return r.install(ctx, obj)
}

func (r *MachineReconciler) scan(ctx context.Context, obj *lifecyclev1alpha1.Machine) (reconcile.Result, error) {
	log := logr.FromContextOrDiscard(ctx)
	resp, err := r.ScanMachine(ctx, connect.NewRequest(&machinev1alpha1.ScanMachineRequest{
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

func (r *MachineReconciler) install(ctx context.Context, obj *lifecyclev1alpha1.Machine) (reconcile.Result, error) {
	log := logr.FromContextOrDiscard(ctx)
	resp, err := r.Install(ctx, connect.NewRequest(&machinev1alpha1.InstallRequest{
		Name:      obj.Name,
		Namespace: obj.Namespace,
	}))
	if err != nil {
		log.Error(err, "failed to send install request")
		return reconcile.Result{}, err
	}
	installResponse := resp.Msg
	result := reconcile.Result{}
	switch {
	case LCIMRequestResultToString[installResponse.Result].IsFailure():
		obj.Status.Message = StatusMessageInstallRequestFailed
		result.RequeueAfter = requeuePeriod
	case LCIMRequestResultToString[installResponse.Result].IsScheduled():
		obj.Status.Message = StatusMessageInstallRequestProcessing
	case LCIMRequestResultToString[installResponse.Result].IsSuccess():
		obj.Status.Message = StatusMessageInstallRequestSuccessful
	}
	return result, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MachineReconciler) SetupWithManager(mgr ctrl.Manager) error {
	labelPredicate, err := predicate.LabelSelectorPredicate(metav1.LabelSelector{
		MatchExpressions: []metav1.LabelSelectorRequirement{
			{
				Key:      "lifecycle.ironcore.dev/exclude",
				Operator: "DoesNotExist",
				Values:   nil,
			},
		},
	})
	if err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&lifecyclev1alpha1.Machine{}, builder.WithPredicates(labelPredicate)).
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
	tempPackages := make([]lifecyclev1alpha1.PackageVersion,
		len(desiredPackages), len(desiredPackages)+len(defaultPackages))
	copy(tempPackages, desiredPackages)

	for _, pv := range defaultPackages {
		idx := slices.IndexFunc(desiredPackages, func(packageVersion lifecyclev1alpha1.PackageVersion) bool {
			return pv.Name == packageVersion.Name
		})
		if idx >= 0 {
			continue
		}
		tempPackages = append(tempPackages, pv)
	}

	resultPackages := make([]lifecyclev1alpha1.PackageVersion, 0, len(tempPackages)+len(installedPackages))
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
	return slices.Clip(resultPackages)
}

func mergePackageVersion(base, add []lifecyclev1alpha1.PackageVersion) []lifecyclev1alpha1.PackageVersion {
	result := make([]lifecyclev1alpha1.PackageVersion, len(base), len(base)+len(add))
	copy(result, base)
	for _, pv := range add {
		idx := slices.IndexFunc(base, func(entry lifecyclev1alpha1.PackageVersion) bool {
			return entry.Name == pv.Name
		})
		if idx >= 0 {
			result[idx] = pv
			continue
		}
		result = append(result, pv)
	}
	return slices.Clip(result)
}
