// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"cmp"
	"context"
	"fmt"
	"slices"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/tools/reference"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/ironcore-dev/lifecycle-manager/api/v1alpha1"
)

const (
	DanglingReference   = "DanglingReference"
	MachineTypeMismatch = "MachineTypeMismatch"
)

var (
	danglingReferenceMessage   = "object not found by reference; kind: %s; reference: %s"
	machineTypeMismatchMessage = "referenced object has different MachineType; kind: %s; reference: %s"

	patchOpts = &client.SubResourcePatchOptions{
		PatchOptions: client.PatchOptions{
			Force:        pointer.Bool(true),
			FieldManager: "lifecycle-manager",
		},
	}
)

// UpdateTaskReconciler reconciles UpdateTask objects.
type UpdateTaskReconciler struct {
	client.Client

	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=metal.ironcore.dev,resources=updatetasks,verbs=watch;get;list;patch;update
// +kubebuilder:rbac:groups=metal.ironcore.dev,resources=updatetasks/status,verbs=get;patch;update
// +kubebuilder:rbac:groups=metal.ironcore.dev,resources=machineupdatejobs,verbs=watch;get;list;create
// +kubebuilder:rbac:groups=metal.ironcore.dev,resources=machineupdatejobs/status,verbs=get
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

func (r *UpdateTaskReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var (
		result ctrl.Result
		ref    *corev1.ObjectReference
		err    error
	)

	obj := &v1alpha1.UpdateTask{}
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

func (r *UpdateTaskReconciler) reconcileRequired(ctx context.Context, obj *v1alpha1.UpdateTask) (ctrl.Result, error) {
	log := logr.FromContextOrDiscard(ctx)
	if !obj.GetDeletionTimestamp().IsZero() {
		log.V(2).Info("object is being deleted")
		return ctrl.Result{}, nil
	}
	if obj.Status.JobsTotal == 0 {
		return r.reconcileNewTask(ctx, obj)
	}
	return r.reconcileExistingTask(ctx, obj)
}

func (r *UpdateTaskReconciler) reconcileNewTask(ctx context.Context, obj *v1alpha1.UpdateTask) (ctrl.Result, error) {
	log := logr.FromContextOrDiscard(ctx)
	var validatedMachinesRefs uint32 = 0
	var validatedPackagesRefs uint32 = 0
	for _, mRef := range obj.Spec.Machines {
		key := types.NamespacedName{Name: mRef.Name, Namespace: obj.Namespace}
		ok, err := r.validateMachineLifecycleRef(ctx, key, obj.Spec.MachineTypeRef.Name)
		if err != nil {
			r.Recorder.Event(obj, corev1.EventTypeWarning, DanglingReference,
				fmt.Sprintf(danglingReferenceMessage, "MachineLifecycle", mRef.Name))
			log.V(2).Info(fmt.Sprintf(danglingReferenceMessage, "MachineLifecycle", mRef.Name))
			return ctrl.Result{}, err
		}
		if !ok {
			r.Recorder.Event(obj, corev1.EventTypeWarning, MachineTypeMismatch,
				fmt.Sprintf(machineTypeMismatchMessage, "MachineLifecycle", mRef.Name))
			log.V(2).Info(fmt.Sprintf(machineTypeMismatchMessage, "MachineLifecycle", mRef.Name))
			return ctrl.Result{}, nil
		}
		validatedMachinesRefs += 1
		processedPackages, err := r.processPackages(ctx, obj, mRef.Name)
		if err != nil {
			return ctrl.Result{}, err
		}
		validatedPackagesRefs += processedPackages
	}
	totalJobs := validatedMachinesRefs * validatedPackagesRefs
	obj.Status.JobsTotal = totalJobs
	log.V(1).WithValues("jobs", totalJobs).Info("object processing completed")
	return ctrl.Result{}, nil
}

func (r *UpdateTaskReconciler) validateMachineLifecycleRef(
	ctx context.Context,
	key types.NamespacedName,
	mType string,
) (bool, error) {
	log := logr.FromContextOrDiscard(ctx)
	machineLifecycle := &v1alpha1.MachineLifecycle{}
	if err := r.Get(ctx, key, machineLifecycle); err != nil {
		log.Error(err, "failed to get referenced machineLifecycle object")
		return false, err
	}
	if machineLifecycle.Spec.MachineTypeRef.Name != mType {
		return false, nil
	}
	return true, nil
}

func (r *UpdateTaskReconciler) processPackages(
	ctx context.Context,
	obj *v1alpha1.UpdateTask,
	mRef string,
) (uint32, error) {
	log := logr.FromContextOrDiscard(ctx)
	var validPackageRefs uint32 = 0
	for _, pRef := range obj.Spec.Packages {
		key := types.NamespacedName{Name: pRef.Name, Namespace: obj.Namespace}
		ok, err := r.validateFirmwarePackageRef(ctx, key, obj.Spec.MachineTypeRef.Name)
		if err != nil {
			r.Recorder.Event(obj, corev1.EventTypeWarning, DanglingReference,
				fmt.Sprintf(danglingReferenceMessage, "FirmwarePackage", pRef.Name))
			log.V(2).Info(fmt.Sprintf(danglingReferenceMessage, "FirmwarePackage", pRef.Name))
			return 0, err
		}
		if !ok {
			r.Recorder.Event(obj, corev1.EventTypeWarning, MachineTypeMismatch,
				fmt.Sprintf(machineTypeMismatchMessage, "FirmwarePackage", pRef.Name))
			log.V(2).Info(fmt.Sprintf(machineTypeMismatchMessage, "FirmwarePackage", pRef.Name))
			continue
		}
		validPackageRefs += 1
		updateJob := &v1alpha1.MachineUpdateJob{
			ObjectMeta: metav1.ObjectMeta{
				Name:      string(uuid.NewUUID()),
				Namespace: obj.Namespace,
				OwnerReferences: []metav1.OwnerReference{
					ownerRefFromTask(obj),
				},
			},
			Spec: v1alpha1.MachineUpdateJobSpec{
				MachineLifecycleRef: corev1.LocalObjectReference{Name: mRef},
				FirmwarePackageRef:  corev1.LocalObjectReference{Name: pRef.Name},
			},
		}
		if err := r.Create(ctx, updateJob); err != nil {
			log.Error(err, "failed to create MachineUpdateJob object")
			return 0, err
		}
	}
	return validPackageRefs, nil
}

func (r *UpdateTaskReconciler) validateFirmwarePackageRef(
	ctx context.Context,
	key types.NamespacedName,
	mType string,
) (bool, error) {
	log := logr.FromContextOrDiscard(ctx)
	firmwarePackage := &v1alpha1.FirmwarePackage{}
	if err := r.Get(ctx, key, firmwarePackage); err != nil {
		log.Error(err, "failed to get referenced firmwarePackage object")
		return false, err
	}
	if firmwarePackage.Spec.MachineTypeRef.Name != mType {
		return false, nil
	}
	return true, nil
}

func (r *UpdateTaskReconciler) reconcileExistingTask(ctx context.Context, obj *v1alpha1.UpdateTask) (ctrl.Result, error) {
	log := logr.FromContextOrDiscard(ctx)
	jobs := &v1alpha1.MachineUpdateJobList{}
	if err := r.List(ctx, jobs); err != nil {
		log.Error(err, "failed to list machineUpdateJob objects")
		return ctrl.Result{}, err
	}
	for _, job := range jobs.Items {
		slices.SortFunc(job.GetOwnerReferences(), func(a, b metav1.OwnerReference) int {
			return cmp.Compare(a.Name, b.Name)
		})
		_, found := slices.BinarySearchFunc(
			job.GetOwnerReferences(), ownerRefFromTask(obj), func(a, b metav1.OwnerReference) int {
				return cmp.Compare(a.Name, b.Name)
			})
		if !found {
			continue
		}
		if job.Status.State == v1alpha1.UpdateJobStateFailure {
			obj.Status.JobsFailed += 1
		}
		if job.Status.State == v1alpha1.UpdateJobStateSuccess {
			obj.Status.JobsSuccessful += 1
		}
	}
	log.V(2).
		WithValues("jobs_failed", obj.Status.JobsFailed).
		WithValues("jobs_successful", obj.Status.JobsSuccessful).
		Info("object processing completed")
	return ctrl.Result{}, nil
}

func (r *UpdateTaskReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.UpdateTask{}).
		Owns(&v1alpha1.MachineUpdateJob{}, builder.WithPredicates(r.machineUpdateJobPredicates())).
		WithOptions(controller.Options{
			MaxConcurrentReconciles: 10,
		}).
		Complete(r)
}

func (r *UpdateTaskReconciler) machineUpdateJobPredicates() predicate.Predicate {
	return predicate.Funcs{
		UpdateFunc: r.enqueueByOwner,
	}
}

func (r *UpdateTaskReconciler) enqueueByOwner(e event.UpdateEvent) bool {
	objOld, okOld := e.ObjectOld.(*v1alpha1.MachineUpdateJob)
	objNew, okNew := e.ObjectNew.(*v1alpha1.MachineUpdateJob)
	if !okOld || !okNew {
		return false
	}
	return objOld.Status.State != objNew.Status.State
}

func ownerRefFromTask(obj *v1alpha1.UpdateTask) metav1.OwnerReference {
	return metav1.OwnerReference{
		APIVersion: obj.APIVersion,
		Kind:       obj.Kind,
		Name:       obj.Name,
		UID:        obj.UID,
		Controller: pointer.Bool(true),
	}
}
