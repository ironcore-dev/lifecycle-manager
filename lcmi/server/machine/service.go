// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package machine

import (
	"context"
	"reflect"
	"slices"
	"time"

	"github.com/go-logr/logr"
	"github.com/jellydator/ttlcache/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"

	lifecyclev1alpha1 "github.com/ironcore-dev/lifecycle-manager/api/lifecycle/v1alpha1"
	"github.com/ironcore-dev/lifecycle-manager/clientgo/applyconfiguration/lifecycle/v1alpha1"
	"github.com/ironcore-dev/lifecycle-manager/clientgo/lifecycle"
	commonv1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/common/v1alpha1"
	machinev1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/machine/v1alpha1"
	"github.com/ironcore-dev/lifecycle-manager/util/apiutil"
	"github.com/ironcore-dev/lifecycle-manager/util/uuidutil"
)

const (
	AddPackageFailureReason = "package already in the list"
	SetPackageFailureReason = "package is not in the list"
)

type GrpcService struct {
	machinev1alpha1.UnimplementedMachineServiceServer
	cl         *lifecycle.Clientset
	cache      *ttlcache.Cache[string, *machinev1alpha1.MachineStatus]
	horizon    time.Duration
	scanPeriod time.Duration
	namespace  string
}

type ServiceOption func(service *GrpcService)

func NewGrpcService(cfg *rest.Config, opts ...ServiceOption) *GrpcService {
	cl := lifecycle.NewForConfigOrDie(cfg)
	svc := &GrpcService{
		cl: cl,
	}
	for _, opt := range opts {
		opt(svc)
	}
	return svc
}

func WithNamespace(namespace string) ServiceOption {
	return func(svc *GrpcService) {
		svc.namespace = namespace
	}
}

func WithHorizon(horizon time.Duration) ServiceOption {
	return func(svc *GrpcService) {
		svc.horizon = horizon
	}
}

func WithScanPeriod(period time.Duration) ServiceOption {
	return func(svc *GrpcService) {
		svc.scanPeriod = period
	}
}

func WithCache(cache *ttlcache.Cache[string, *machinev1alpha1.MachineStatus]) ServiceOption {
	return func(svc *GrpcService) {
		svc.cache = cache
	}
}

// ListMachines returns the list of Machine objects.
func (s *GrpcService) ListMachines(
	ctx context.Context,
	req *machinev1alpha1.ListMachinesRequest,
) (*machinev1alpha1.ListMachinesResponse, error) {
	log := logr.FromContextOrDiscard(ctx)
	if req.Filter == nil {
		return nil, status.Error(codes.InvalidArgument, "filter is mandatory field")
	}
	namespace := req.Filter.Namespace
	if namespace == "" {
		namespace = s.namespace
	}
	opts := metav1.ListOptions{}
	if req.Filter.LabelSelector != nil {
		opts.LabelSelector = labels.Set(req.Filter.LabelSelector.MatchLabels).String()
	}
	machines, err := s.cl.LifecycleV1alpha1().Machines(namespace).List(ctx, opts)
	if err != nil {
		log.Error(err, "failed to list machines")
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}
	resp := &machinev1alpha1.ListMachinesResponse{Machines: make([]*machinev1alpha1.Machine, len(machines.Items))}
	for i, item := range machines.Items {
		m := item.DeepCopy()
		resp.Machines[i] = apiutil.MachineToGrpcAPI(m)
	}
	return resp, nil
}

// ScanMachine runs scan job for target Machine. First it checks whether Machine
// state is in service's cache. If cache entry is found, it checks whether last
// scan timestamp is within defined horizon. If not, scan job will be spawned.
// Otherwise, cached state will be returned in response.
func (s *GrpcService) ScanMachine(
	ctx context.Context,
	req *machinev1alpha1.ScanMachineRequest,
) (*machinev1alpha1.ScanMachineResponse, error) {
	resp := &machinev1alpha1.ScanMachineResponse{
		Result: commonv1alpha1.RequestResult_REQUEST_RESULT_UNSPECIFIED,
	}
	namespace := req.Namespace
	if namespace == "" {
		namespace = s.namespace
	}
	log := logr.FromContextOrDiscard(ctx).WithValues("name", req.Name, "namespace", namespace)

	machine, err := s.cl.LifecycleV1alpha1().Machines(namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		log.Error(err, "failed to get machine")
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}

	key := types.NamespacedName{Name: req.Name, Namespace: namespace}
	cacheItem := s.cache.Get(uuidutil.UUIDFromObjectKey(key))
	if cacheItem == nil {
		// todo: if scan job already scheduled, return, otherwise send to scheduler
		log.V(1).Info("scan scheduled")
		resp.Result = commonv1alpha1.RequestResult_REQUEST_RESULT_SCHEDULED
	}
	entry := cacheItem.Value()
	switch {
	case !lastScanTimeInHorizon(s.horizon, entry.LastScanTime), !installedPackagesEqual(machine.Status, entry):
		// todo: if scan job already scheduled, return, otherwise send to scheduler
		log.V(1).Info("scan scheduled")
		resp.Result = commonv1alpha1.RequestResult_REQUEST_RESULT_SCHEDULED
	default:
		resp.Result = commonv1alpha1.RequestResult_REQUEST_RESULT_SUCCESS
	}

	return resp, nil
}

// Install schedules package installation for target Machine.
func (s *GrpcService) Install(
	ctx context.Context,
	req *machinev1alpha1.InstallRequest,
) (*machinev1alpha1.InstallResponse, error) {
	resp := &machinev1alpha1.InstallResponse{
		Result: commonv1alpha1.RequestResult_REQUEST_RESULT_UNSPECIFIED,
	}
	namespace := req.Namespace
	if namespace == "" {
		namespace = s.namespace
	}
	log := logr.FromContextOrDiscard(ctx).WithValues("name", req.Name, "namespace", namespace)

	// todo: if install task already scheduled, return, otherwise send to scheduler
	log.V(1).Info("packages installation scheduled")
	resp.Result = commonv1alpha1.RequestResult_REQUEST_RESULT_SCHEDULED
	return resp, nil
}

// UpdateMachineStatus request initialized by the spawned Job and should update
// the status of processed Machine. If request succeed, Job exits with exit code 0,
// otherwise, Job will stop with non-zero exit code.
func (s *GrpcService) UpdateMachineStatus(
	ctx context.Context,
	req *machinev1alpha1.UpdateMachineStatusRequest,
) (*machinev1alpha1.UpdateMachineStatusResponse, error) {
	namespace := req.Namespace
	if namespace == "" {
		namespace = s.namespace
	}
	log := logr.FromContextOrDiscard(ctx).WithValues("name", req.Name, "namespace", namespace)
	machine, err := s.cl.LifecycleV1alpha1().Machines(namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		log.Error(err, "failed to get machine")
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}
	machineApply := v1alpha1.Machine(machine.Name, machine.Namespace).
		WithStatus(apiutil.MachineStatusToApplyConfiguration(req.Status))
	if _, err = s.cl.LifecycleV1alpha1().Machines(namespace).ApplyStatus(ctx, machineApply, metav1.ApplyOptions{
		FieldManager: "lifecycle.ironcore.dev/lifecycle-manager",
		Force:        true,
	}); err != nil {
		log.Error(err, "failed to update machine status")
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}
	key := types.NamespacedName{Name: req.Name, Namespace: namespace}
	s.cache.Set(uuidutil.UUIDFromObjectKey(key), req.Status, ttlcache.DefaultTTL)
	return &machinev1alpha1.UpdateMachineStatusResponse{
		Result: commonv1alpha1.RequestResult_REQUEST_RESULT_SUCCESS,
	}, nil
}

func (s *GrpcService) AddPackageVersion(
	ctx context.Context,
	req *machinev1alpha1.AddPackageVersionRequest,
) (*machinev1alpha1.AddPackageVersionResponse, error) {
	namespace := req.Namespace
	if namespace == "" {
		namespace = s.namespace
	}
	log := logr.FromContextOrDiscard(ctx).WithValues("name", req.Name, "namespace", namespace)
	machine, err := s.cl.LifecycleV1alpha1().Machines(namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		log.Error(err, "failed to get machine")
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}
	pkg := apiutil.PackageVersionsToGrpcAPI(machine.Spec.Packages)
	if packageIndex(req.Package.Name, pkg) == -1 {
		return &machinev1alpha1.AddPackageVersionResponse{
			Reason: AddPackageFailureReason,
			Result: commonv1alpha1.RequestResult_REQUEST_RESULT_FAILURE,
		}, nil
	}
	pkg = slices.Grow(pkg, 1)
	pkg = append(pkg, req.Package)
	machineApply := v1alpha1.Machine(machine.Name, machine.Namespace).WithSpec(
		v1alpha1.MachineSpec().WithPackages(apiutil.PackageVersionsToApplyConfiguration(pkg)...))
	if _, err = s.cl.LifecycleV1alpha1().Machines(namespace).Apply(ctx, machineApply, metav1.ApplyOptions{
		FieldManager: "lifecycle.ironcore.dev/lifecycle-manager",
		Force:        true,
	}); err != nil {
		log.Error(err, "failed to update machine")
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}
	return &machinev1alpha1.AddPackageVersionResponse{
		Reason: "",
		Result: commonv1alpha1.RequestResult_REQUEST_RESULT_SUCCESS,
	}, nil
}

func (s *GrpcService) SetPackageVersion(
	ctx context.Context,
	req *machinev1alpha1.SetPackageVersionRequest,
) (*machinev1alpha1.SetPackageVersionResponse, error) {
	namespace := req.Namespace
	if namespace == "" {
		namespace = s.namespace
	}
	log := logr.FromContextOrDiscard(ctx).WithValues("name", req.Name, "namespace", namespace)
	machine, err := s.cl.LifecycleV1alpha1().Machines(namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		log.Error(err, "failed to get machine")
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}
	pkg := apiutil.PackageVersionsToGrpcAPI(machine.Spec.Packages)
	idx := packageIndex(req.Package.Name, pkg)
	if idx < 0 {
		return &machinev1alpha1.SetPackageVersionResponse{
			Reason: SetPackageFailureReason,
			Result: commonv1alpha1.RequestResult_REQUEST_RESULT_FAILURE,
		}, nil
	}
	pkg[idx] = req.Package
	machineApply := v1alpha1.Machine(machine.Name, machine.Namespace).WithSpec(
		v1alpha1.MachineSpec().WithPackages(apiutil.PackageVersionsToApplyConfiguration(pkg)...))
	if _, err = s.cl.LifecycleV1alpha1().Machines(namespace).Apply(ctx, machineApply, metav1.ApplyOptions{
		FieldManager: "lifecycle.ironcore.dev/lifecycle-manager",
		Force:        true,
	}); err != nil {
		log.Error(err, "failed to update machine")
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}
	return &machinev1alpha1.SetPackageVersionResponse{
		Reason: "",
		Result: commonv1alpha1.RequestResult_REQUEST_RESULT_SUCCESS,
	}, nil
}

func (s *GrpcService) RemovePackageVersion(
	ctx context.Context,
	req *machinev1alpha1.RemovePackageVersionRequest,
) (*machinev1alpha1.RemovePackageVersionResponse, error) {
	namespace := req.Namespace
	if namespace == "" {
		namespace = s.namespace
	}
	log := logr.FromContextOrDiscard(ctx).WithValues("name", req.Name, "namespace", namespace)
	machine, err := s.cl.LifecycleV1alpha1().Machines(namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		log.Error(err, "failed to get machine")
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}
	pkg := apiutil.PackageVersionsToGrpcAPI(machine.Spec.Packages)
	idx := packageIndex(req.PackageName, pkg)
	if idx < 0 {
		return &machinev1alpha1.RemovePackageVersionResponse{
			Reason: "",
			Result: commonv1alpha1.RequestResult_REQUEST_RESULT_SUCCESS,
		}, nil
	}
	pkg = removePackage(pkg, idx)
	machineApply := v1alpha1.Machine(machine.Name, machine.Namespace).WithSpec(
		v1alpha1.MachineSpec().WithPackages(apiutil.PackageVersionsToApplyConfiguration(pkg)...))
	if _, err = s.cl.LifecycleV1alpha1().Machines(namespace).Apply(ctx, machineApply, metav1.ApplyOptions{
		FieldManager: "lifecycle.ironcore.dev/lifecycle-manager",
		Force:        true,
	}); err != nil {
		log.Error(err, "failed to update machine")
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}
	return &machinev1alpha1.RemovePackageVersionResponse{
		Reason: "",
		Result: commonv1alpha1.RequestResult_REQUEST_RESULT_SUCCESS,
	}, nil
}

func lastScanTimeInHorizon(horizon time.Duration, timestamp *metav1.Time) bool {
	return time.Since(timestamp.Time) < horizon
}

func installedPackagesEqual(
	src lifecyclev1alpha1.MachineStatus,
	tgt *machinev1alpha1.MachineStatus,
) bool {
	conv := apiutil.MachineStatusToKubeAPI(tgt)
	return reflect.DeepEqual(src, conv)
}

func packageIndex(pkg string, dst []*commonv1alpha1.PackageVersion) int {
	return slices.IndexFunc(dst, func(pv *commonv1alpha1.PackageVersion) bool {
		return pkg == pv.Name
	})
}

func removePackage(src []*commonv1alpha1.PackageVersion, index int) []*commonv1alpha1.PackageVersion {
	result := make([]*commonv1alpha1.PackageVersion, 0, len(src)-1)
	if len(src) == 1 {
		return result
	}
	result = append(result, src[:index]...)
	if len(src) == 2 && index == 1 {
		return result
	}
	result = append(result, src[index+1:]...)
	return result
}
