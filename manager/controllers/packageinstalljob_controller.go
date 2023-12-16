// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"
	"slices"
	"time"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/reference"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	"github.com/ironcore-dev/lifecycle-manager/api/v1alpha1"
	lcmimeta "github.com/ironcore-dev/lifecycle-manager/lcmi/api/meta/v1alpha1"
	lcmipackageinstalljob "github.com/ironcore-dev/lifecycle-manager/lcmi/api/packageinstalljob/v1alpha1"
	"github.com/ironcore-dev/lifecycle-manager/util/uuidutil"
)

// PackageInstallJobReconciler reconciles PackageInstallJob objects.
type PackageInstallJobReconciler struct {
	client.Client

	BrokerClient lcmipackageinstalljob.PackageInstallJobBrokerServiceClient

	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=lifecycle.ironcore.dev,resources=packageinstalljobs,verbs=get;list;watch;update;patch
// +kubebuilder:rbac:groups=lifecycle.ironcore.dev,resources=packageinstalljobs/status,verbs=get;update;patch

func (r *PackageInstallJobReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var (
		result ctrl.Result
		ref    *corev1.ObjectReference
		err    error
	)

	obj := &v1alpha1.PackageInstallJob{}
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

func (r *PackageInstallJobReconciler) reconcileRequired(
	ctx context.Context,
	obj *v1alpha1.PackageInstallJob,
) (ctrl.Result, error) {
	log := logr.FromContextOrDiscard(ctx)
	if !obj.GetDeletionTimestamp().IsZero() {
		log.V(2).Info("object is being deleted")
		return ctrl.Result{}, nil
	}
	return r.reconcile(ctx, obj)
}

func (r *PackageInstallJobReconciler) reconcile(
	ctx context.Context,
	obj *v1alpha1.PackageInstallJob,
) (ctrl.Result, error) {
	log := logr.FromContextOrDiscard(ctx)
	resp, err := r.BrokerClient.ListPackageInstallJobs(ctx, &lcmipackageinstalljob.ListPackageInstallJobsRequest{
		Filter: &lcmipackageinstalljob.PackageInstallJobFilter{
			Id: uuidutil.UUIDFromObjectKey(client.ObjectKeyFromObject(obj)),
		},
	})
	if err != nil {
		log.Error(err, "failed to list package install job runners")
		return ctrl.Result{}, err
	}
	if len(resp.PackageInstallJob) == 0 {
		return r.createNewJob(ctx, obj)
	}
	return r.reconcileExistingJob(ctx, obj, resp.PackageInstallJob[0])
}

func (r *PackageInstallJobReconciler) createNewJob(
	ctx context.Context,
	obj *v1alpha1.PackageInstallJob,
) (ctrl.Result, error) {
	log := logr.FromContextOrDiscard(ctx)
	resp, err := r.BrokerClient.Execute(ctx, &lcmipackageinstalljob.ExecuteRequest{
		PackageInstallJob: &lcmipackageinstalljob.PackageInstallJob{
			Metadata: &lcmimeta.ObjectMetadata{
				Id:          uuidutil.UUIDFromObjectKey(client.ObjectKeyFromObject(obj)),
				Annotations: obj.GetAnnotations(),
				Labels:      obj.GetLabels(),
			},
			Spec: &lcmipackageinstalljob.PackageInstallJobSpec{
				MachineRef: &lcmimeta.LocalObjectReference{
					Name: obj.Spec.MachineRef.Name,
				},
				FirmwarePackageRef: &lcmimeta.LocalObjectReference{
					Name: obj.Spec.FirmwarePackageRef.Name,
				},
			},
		},
	})
	if err != nil {
		log.Error(err, "failed to create package install job runner")
		return ctrl.Result{}, err
	}
	idx := slices.IndexFunc(resp.Status.Conditions, func(condition *lcmimeta.Condition) bool {
		return condition.Type == 1
	})
	respCondition := resp.Status.Conditions[idx]
	condition := v1alpha1.PackageInstallJobCondition{Condition: metav1.Condition{
		Type:               LCMIConditionTypeToString[respCondition.Type],
		Status:             metav1.ConditionStatus(respCondition.Status),
		ObservedGeneration: obj.Generation,
		LastTransitionTime: metav1.Time{Time: time.Unix(0, respCondition.LastTransitionTime)},
		Reason:             respCondition.Reason,
		Message:            respCondition.Message,
	}}
	obj.Status.Conditions = append(obj.Status.Conditions, condition)
	obj.Status.Message = respCondition.Message
	return ctrl.Result{}, nil
}

func (r *PackageInstallJobReconciler) reconcileExistingJob(
	ctx context.Context,
	obj *v1alpha1.PackageInstallJob,
	lcmiJob *lcmipackageinstalljob.PackageInstallJob,
) (ctrl.Result, error) {
	log := logr.FromContextOrDiscard(ctx)
	conditions := make([]v1alpha1.PackageInstallJobCondition, len(lcmiJob.Status.Conditions))
	for _, jobCondition := range lcmiJob.Status.Conditions {
		condition := v1alpha1.PackageInstallJobCondition{Condition: metav1.Condition{
			Type:               LCMIConditionTypeToString[jobCondition.Type],
			Status:             metav1.ConditionStatus(jobCondition.Status),
			ObservedGeneration: obj.Generation,
			LastTransitionTime: metav1.Time{Time: time.Unix(0, jobCondition.LastTransitionTime)},
			Reason:             jobCondition.Reason,
			Message:            jobCondition.Message,
		}}
		conditions = append(conditions, condition)
	}
	obj.Status.Conditions = conditions
	obj.Status.State = LCMIJobStateToString[lcmiJob.Status.State]
	obj.Status.Message = lcmiJob.Status.Message
	if obj.Status.State.IsFailure() || obj.Status.State.IsSuccess() {
		log.V(1).WithValues("job_state", obj.Status.State).Info("package install job finished")
	}
	return ctrl.Result{}, nil
}

func (r *PackageInstallJobReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.PackageInstallJob{}).
		WithOptions(controller.Options{
			MaxConcurrentReconciles: 10,
		}).
		Complete(r)
}
