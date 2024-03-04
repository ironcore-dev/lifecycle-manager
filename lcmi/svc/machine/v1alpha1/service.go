// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	"context"
	"reflect"
	"slices"
	"time"

	"connectrpc.com/connect"
	lifecyclev1alpha1 "github.com/ironcore-dev/lifecycle-manager/api/lifecycle/v1alpha1"
	"github.com/ironcore-dev/lifecycle-manager/clientgo/applyconfiguration/lifecycle/v1alpha1"
	"github.com/ironcore-dev/lifecycle-manager/clientgo/lifecycle"
	commonv1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/common/v1alpha1"
	machinev1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/machine/v1alpha1"
	"github.com/ironcore-dev/lifecycle-manager/util/apiutil"
	"github.com/ironcore-dev/lifecycle-manager/util/uuidutil"
	"github.com/jellydator/ttlcache/v3"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
)

const (
	AddPackageFailureReason = "package already in the list"
	SetPackageFailureReason = "package is not in the list"
)

type MachineService struct {
	machinev1alpha1.UnimplementedMachineServiceServer
	c         *lifecycle.Clientset
	cache     *ttlcache.Cache[string, *machinev1alpha1.MachineStatus]
	horizon   time.Duration
	period    time.Duration
	namespace string
}

type Option func(service *MachineService)

func NewService(cfg *rest.Config, opts ...Option) *MachineService {
	c := lifecycle.NewForConfigOrDie(cfg)
	svc := &MachineService{c: c}
	for _, opt := range opts {
		opt(svc)
	}
	return svc
}

func WithNamespace(namespace string) Option {
	return func(svc *MachineService) {
		svc.namespace = namespace
	}
}

func WithHorizon(horizon time.Duration) Option {
	return func(svc *MachineService) {
		svc.horizon = horizon
	}
}

func WithScanPeriod(period time.Duration) Option {
	return func(svc *MachineService) {
		svc.period = period
	}
}

func WithCache(cache *ttlcache.Cache[string, *machinev1alpha1.MachineStatus]) Option {
	return func(svc *MachineService) {
		svc.cache = cache
	}
}

// ScanMachine runs scan job for target Machine. First it checks whether Machine
// state is in service's cache. If cache entry is found, it checks whether last
// scan timestamp is within defined horizon. If not, scan job will be spawned.
// Otherwise, cached state will be returned in response.
func (s MachineService) ScanMachine(
	ctx context.Context,
	c *connect.Request[machinev1alpha1.ScanMachineRequest],
) (*connect.Response[machinev1alpha1.ScanMachineResponse], error) {
	req := c.Msg
	namespace := req.GetNamespace()
	if namespace == "" {
		namespace = s.namespace
	}
	resp := &machinev1alpha1.ScanMachineResponse{
		Result: commonv1alpha1.RequestResult_REQUEST_RESULT_UNSPECIFIED,
	}

	machine, err := s.c.LifecycleV1alpha1().Machines(namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		errCode := connect.CodeInternal
		if apierrors.IsNotFound(err) {
			errCode = connect.CodeNotFound
		}
		return nil, connect.NewError(errCode, err)
	}

	key := types.NamespacedName{Name: req.Name, Namespace: namespace}
	cacheItem := s.cache.Get(uuidutil.UUIDFromObjectKey(key))

	switch {
	case cacheItem == nil:
		fallthrough
	case time.Since(cacheItem.Value().LastScanTime.Time) > s.horizon:
		fallthrough
	case !installedPackagesEqual(machine.Status, cacheItem.Value()):
		// todo: result should be returned by the call to scheduler, like
		//  resp.Result = s.scheduleJob()
		resp.Result = commonv1alpha1.RequestResult_REQUEST_RESULT_SCHEDULED
	default:
		resp.Result = commonv1alpha1.RequestResult_REQUEST_RESULT_SUCCESS
	}
	return connect.NewResponse(resp), nil
}

// Install schedules package installation for target Machine.
func (s MachineService) Install(
	ctx context.Context,
	c *connect.Request[machinev1alpha1.InstallRequest],
) (*connect.Response[machinev1alpha1.InstallResponse], error) {
	req := c.Msg
	namespace := req.Namespace
	if namespace == "" {
		namespace = s.namespace
	}
	resp := &machinev1alpha1.InstallResponse{
		Result: commonv1alpha1.RequestResult_REQUEST_RESULT_UNSPECIFIED,
	}

	_, err := s.c.LifecycleV1alpha1().Machines(namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		errCode := connect.CodeInternal
		if apierrors.IsNotFound(err) {
			errCode = connect.CodeNotFound
		}
		return nil, connect.NewError(errCode, err)
	}
	// todo: result should be returned by the call to scheduler, like
	//  resp.Result = s.scheduleJob()
	resp.Result = commonv1alpha1.RequestResult_REQUEST_RESULT_SCHEDULED
	return connect.NewResponse(resp), nil
}

// UpdateMachineStatus request initialized by the spawned Job and should update
// the status of processed Machine. If request succeed, Job exits with exit code 0,
// otherwise, Job will stop with non-zero exit code.
func (s MachineService) UpdateMachineStatus(
	ctx context.Context,
	c *connect.Request[machinev1alpha1.UpdateMachineStatusRequest],
) (*connect.Response[machinev1alpha1.UpdateMachineStatusResponse], error) {
	req := c.Msg
	namespace := req.Namespace
	if namespace == "" {
		namespace = s.namespace
	}

	machine, err := s.c.LifecycleV1alpha1().Machines(namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		errCode := connect.CodeInternal
		if apierrors.IsNotFound(err) {
			errCode = connect.CodeNotFound
		}
		return nil, connect.NewError(errCode, err)
	}
	machineApply := v1alpha1.Machine(machine.Name, machine.Namespace).
		WithStatus(apiutil.MachineStatusToApplyConfiguration(req.Status))
	if _, err = s.c.LifecycleV1alpha1().Machines(namespace).ApplyStatus(ctx, machineApply, metav1.ApplyOptions{
		FieldManager: "lifecycle.ironcore.dev/lifecycle-manager",
		Force:        true,
	}); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	key := types.NamespacedName{Name: req.Name, Namespace: namespace}
	s.cache.Set(uuidutil.UUIDFromObjectKey(key), req.Status, ttlcache.DefaultTTL)
	return connect.NewResponse(&machinev1alpha1.UpdateMachineStatusResponse{
		Result: commonv1alpha1.RequestResult_REQUEST_RESULT_SUCCESS,
	}), nil
}

// ListMachines returns the list of Machine objects.
func (s MachineService) ListMachines(
	ctx context.Context,
	c *connect.Request[machinev1alpha1.ListMachinesRequest],
) (*connect.Response[machinev1alpha1.ListMachinesResponse], error) {
	req := c.Msg
	namespace := req.GetNamespace()
	if namespace == "" {
		namespace = s.namespace
	}

	opts := metav1.ListOptions{}
	reqSelector := req.GetLabelSelector()
	if reqSelector != nil {
		opts.LabelSelector = labels.Set(reqSelector.MatchLabels).String()
	}
	machines, err := s.c.LifecycleV1alpha1().Machines(namespace).List(ctx, opts)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	resp := &machinev1alpha1.ListMachinesResponse{Machines: make([]*machinev1alpha1.Machine, len(machines.Items))}
	for i, item := range machines.Items {
		m := item.DeepCopy()
		resp.Machines[i] = apiutil.MachineToGrpcAPI(m)
	}
	return connect.NewResponse(resp), nil
}

func (s MachineService) AddPackageVersion(
	ctx context.Context,
	c *connect.Request[machinev1alpha1.AddPackageVersionRequest],
) (*connect.Response[machinev1alpha1.AddPackageVersionResponse], error) {
	req := c.Msg
	namespace := req.GetNamespace()
	if namespace == "" {
		namespace = s.namespace
	}

	machine, err := s.c.LifecycleV1alpha1().Machines(namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		errCode := connect.CodeInternal
		if apierrors.IsNotFound(err) {
			errCode = connect.CodeNotFound
		}
		return nil, connect.NewError(errCode, err)
	}
	pkg := apiutil.PackageVersionsToGrpcAPI(machine.Spec.Packages)
	if packageIndex(req.Package.Name, pkg) > -1 {
		return connect.NewResponse(&machinev1alpha1.AddPackageVersionResponse{
			Reason: AddPackageFailureReason,
			Result: commonv1alpha1.RequestResult_REQUEST_RESULT_FAILURE,
		}), nil
	}
	pkg = slices.Grow(pkg, 1)
	pkg = append(pkg, req.Package)
	machineApply := v1alpha1.Machine(machine.Name, machine.Namespace).WithSpec(
		v1alpha1.MachineSpec().WithPackages(apiutil.PackageVersionsToApplyConfiguration(pkg)...))
	if _, err = s.c.LifecycleV1alpha1().Machines(namespace).Apply(ctx, machineApply, metav1.ApplyOptions{
		FieldManager: "lifecycle.ironcore.dev/lifecycle-manager",
		Force:        true,
	}); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&machinev1alpha1.AddPackageVersionResponse{
		Result: commonv1alpha1.RequestResult_REQUEST_RESULT_SUCCESS,
	}), nil
}

func (s MachineService) SetPackageVersion(
	ctx context.Context,
	c *connect.Request[machinev1alpha1.SetPackageVersionRequest],
) (*connect.Response[machinev1alpha1.SetPackageVersionResponse], error) {
	req := c.Msg
	namespace := req.GetNamespace()
	if namespace == "" {
		namespace = s.namespace
	}

	machine, err := s.c.LifecycleV1alpha1().Machines(namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		errCode := connect.CodeInternal
		if apierrors.IsNotFound(err) {
			errCode = connect.CodeNotFound
		}
		return nil, connect.NewError(errCode, err)
	}
	pkg := apiutil.PackageVersionsToGrpcAPI(machine.Spec.Packages)
	var idx int
	if idx = packageIndex(req.Package.Name, pkg); idx == -1 {
		return connect.NewResponse(&machinev1alpha1.SetPackageVersionResponse{
			Reason: SetPackageFailureReason,
			Result: commonv1alpha1.RequestResult_REQUEST_RESULT_FAILURE,
		}), nil
	}
	pkg[idx] = req.Package
	machineApply := v1alpha1.Machine(machine.Name, machine.Namespace).WithSpec(
		v1alpha1.MachineSpec().WithPackages(apiutil.PackageVersionsToApplyConfiguration(pkg)...))
	if _, err = s.c.LifecycleV1alpha1().Machines(namespace).Apply(ctx, machineApply, metav1.ApplyOptions{
		FieldManager: "lifecycle.ironcore.dev/lifecycle-manager",
		Force:        true,
	}); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&machinev1alpha1.SetPackageVersionResponse{
		Result: commonv1alpha1.RequestResult_REQUEST_RESULT_SUCCESS,
	}), nil
}

func (s MachineService) RemovePackageVersion(
	ctx context.Context,
	c *connect.Request[machinev1alpha1.RemovePackageVersionRequest],
) (*connect.Response[machinev1alpha1.RemovePackageVersionResponse], error) {
	req := c.Msg
	namespace := req.GetNamespace()
	if namespace == "" {
		namespace = s.namespace
	}

	machine, err := s.c.LifecycleV1alpha1().Machines(namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		errCode := connect.CodeInternal
		if apierrors.IsNotFound(err) {
			errCode = connect.CodeNotFound
		}
		return nil, connect.NewError(errCode, err)
	}
	pkg := apiutil.PackageVersionsToGrpcAPI(machine.Spec.Packages)
	var idx int
	if idx = packageIndex(req.PackageName, pkg); idx == -1 {
		return connect.NewResponse(&machinev1alpha1.RemovePackageVersionResponse{
			Result: commonv1alpha1.RequestResult_REQUEST_RESULT_SUCCESS,
		}), nil
	}
	pkg = removePackage(pkg, idx)
	machineApply := v1alpha1.Machine(machine.Name, machine.Namespace).WithSpec(
		v1alpha1.MachineSpec().WithPackages(apiutil.PackageVersionsToApplyConfiguration(pkg)...))
	if _, err = s.c.LifecycleV1alpha1().Machines(namespace).Apply(ctx, machineApply, metav1.ApplyOptions{
		FieldManager: "lifecycle.ironcore.dev/lifecycle-manager",
		Force:        true,
	}); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&machinev1alpha1.RemovePackageVersionResponse{
		Result: commonv1alpha1.RequestResult_REQUEST_RESULT_SUCCESS,
	}), nil
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
