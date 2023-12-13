package controllers

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

	"github.com/ironcore-dev/lifecycle-manager/api/v1alpha1"
)

// MachineTypeReconciler reconciles MachineType objects.
type MachineTypeReconciler struct {
	client.Client

	Log     logr.Logger
	Scheme  *runtime.Scheme
	Horizon time.Duration
}

// +kubebuilder:rbac:groups=metal.ironcore.dev,resources=machinetypes,verbs=get;list;watch;update;patch
// +kubebuilder:rbac:groups=metal.ironcore.dev,resources=machinetypes/status,verbs=get;update;patch

func (r *MachineTypeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var (
		result ctrl.Result
		ref    *corev1.ObjectReference
		err    error
	)

	obj := &v1alpha1.MachineType{}
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

func (r *MachineTypeReconciler) reconcileRequired(ctx context.Context, obj *v1alpha1.MachineType) (ctrl.Result, error) {
	log := logr.FromContextOrDiscard(ctx)
	if !obj.GetDeletionTimestamp().IsZero() {
		log.V(2).Info("object is being deleted")
		return ctrl.Result{}, nil
	}
	return r.reconcile(ctx, obj)
}

func (r *MachineTypeReconciler) reconcile(_ context.Context, obj *v1alpha1.MachineType) (ctrl.Result, error) {
	// 1. discover downloader service
	// 2. check object status
	// 2. -> 3. if empty, then it is a brand-new machine type, so force scan required
	// 2. -> 4. otherwise check if the last scan timestamp is in the horizon relative to current time & whether the last scan was successful
	// 4. -> 5. if not, then request force scan

	if r.scanRequired(obj) {
		// request scan -> get response immediately if scan was performed within horizon -> update status -> return
		// request scan -> get response that scan is scheduled -> return (wait for event spawned by downloader service)
		return ctrl.Result{}, nil
	}
	return ctrl.Result{}, nil
}

func (r *MachineTypeReconciler) lastScanNotInHorizon(obj *v1alpha1.MachineType) bool {
	return obj.Status.LastScanTime.IsZero() || time.Since(obj.Status.LastScanTime.Time) > r.Horizon
}

func (r *MachineTypeReconciler) scanRequired(obj *v1alpha1.MachineType) bool {
	if obj.Status.LastScanTime.IsZero() {
		return true
	}
	return r.lastScanNotInHorizon(obj) || obj.Status.LastScanResult == v1alpha1.ScanFailure
}

func (r *MachineTypeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.MachineType{}).
		WithOptions(controller.Options{
			MaxConcurrentReconciles: 10,
		}).
		Complete(r)
}
