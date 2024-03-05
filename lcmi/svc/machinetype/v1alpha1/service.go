// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	"context"
	"slices"
	"time"

	"connectrpc.com/connect"
	"github.com/ironcore-dev/lifecycle-manager/clientgo/applyconfiguration/lifecycle/v1alpha1"
	"github.com/ironcore-dev/lifecycle-manager/clientgo/lifecycle"
	commonv1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/common/v1alpha1"
	machinetypev1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/machinetype/v1alpha1"
	"github.com/ironcore-dev/lifecycle-manager/lcmi/api/machinetype/v1alpha1/machinetypev1alpha1connect"
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
	AddMachineGroupFailureReason = "machine group already in the list"
)

type MachineTypeService struct {
	machinetypev1alpha1connect.UnimplementedMachineTypeServiceHandler
	c         *lifecycle.Clientset
	cache     *ttlcache.Cache[string, *machinetypev1alpha1.MachineTypeStatus]
	horizon   time.Duration
	period    time.Duration
	namespace string
}

type Option func(service *MachineTypeService)

func NewService(cfg *rest.Config, opts ...Option) *MachineTypeService {
	c := lifecycle.NewForConfigOrDie(cfg)
	svc := &MachineTypeService{c: c}
	for _, opt := range opts {
		opt(svc)
	}
	return svc
}

func WithNamespace(namespace string) Option {
	return func(svc *MachineTypeService) {
		svc.namespace = namespace
	}
}

func WithHorizon(horizon time.Duration) Option {
	return func(svc *MachineTypeService) {
		svc.horizon = horizon
	}
}

func WithScanPeriod(period time.Duration) Option {
	return func(svc *MachineTypeService) {
		svc.period = period
	}
}

func WithCache(cache *ttlcache.Cache[string, *machinetypev1alpha1.MachineTypeStatus]) Option {
	return func(svc *MachineTypeService) {
		svc.cache = cache
	}
}

func (s *MachineTypeService) ListMachineTypes(
	ctx context.Context,
	c *connect.Request[machinetypev1alpha1.ListMachineTypesRequest],
) (*connect.Response[machinetypev1alpha1.ListMachineTypesResponse], error) {
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
	machinetypes, err := s.c.LifecycleV1alpha1().MachineTypes(namespace).List(ctx, opts)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	resp := &machinetypev1alpha1.ListMachineTypesResponse{
		MachineTypes: make([]*machinetypev1alpha1.MachineType, len(machinetypes.Items))}
	for i, item := range machinetypes.Items {
		m := item.DeepCopy()
		resp.MachineTypes[i] = apiutil.MachineTypeToGrpcAPI(m)
	}
	return connect.NewResponse(resp), nil
}

func (s *MachineTypeService) Scan(
	_ context.Context,
	c *connect.Request[machinetypev1alpha1.ScanRequest],
) (*connect.Response[machinetypev1alpha1.ScanResponse], error) {
	req := c.Msg
	namespace := req.GetNamespace()
	if namespace == "" {
		namespace = s.namespace
	}
	resp := &machinetypev1alpha1.ScanResponse{
		Result: commonv1alpha1.RequestResult_REQUEST_RESULT_UNSPECIFIED,
	}

	key := types.NamespacedName{Name: req.Name, Namespace: namespace}
	cacheItem := s.cache.Get(uuidutil.UUIDFromObjectKey(key))
	switch {
	case cacheItem == nil:
		fallthrough
	case time.Since(cacheItem.Value().LastScanTime.Time) > s.horizon:
		// todo: result should be returned by the call to scheduler, like
		//  resp.Result = s.scheduleJob()
		resp.Result = commonv1alpha1.RequestResult_REQUEST_RESULT_SCHEDULED
	default:
		resp.Result = commonv1alpha1.RequestResult_REQUEST_RESULT_SUCCESS
	}
	return connect.NewResponse(resp), nil
}

func (s *MachineTypeService) UpdateMachineTypeStatus(
	ctx context.Context,
	c *connect.Request[machinetypev1alpha1.UpdateMachineTypeStatusRequest],
) (*connect.Response[machinetypev1alpha1.UpdateMachineTypeStatusResponse], error) {
	req := c.Msg
	namespace := req.GetNamespace()
	if namespace == "" {
		namespace = s.namespace
	}

	machinetype, err := s.c.LifecycleV1alpha1().MachineTypes(namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		errCode := connect.CodeInternal
		if apierrors.IsNotFound(err) {
			errCode = connect.CodeNotFound
		}
		return nil, connect.NewError(errCode, err)
	}
	machinetypeApply := v1alpha1.MachineType(machinetype.Name, machinetype.Namespace).
		WithStatus(apiutil.MachineTypeStatusToApplyConfiguration(req.Status))
	if _, err = s.c.LifecycleV1alpha1().MachineTypes(namespace).ApplyStatus(ctx, machinetypeApply, metav1.ApplyOptions{
		FieldManager: "lifecycle.ironcore.dev/lifecycle-manager",
		Force:        true,
	}); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	key := types.NamespacedName{Name: req.Name, Namespace: namespace}
	s.cache.Set(uuidutil.UUIDFromObjectKey(key), req.Status, ttlcache.DefaultTTL)
	return connect.NewResponse(&machinetypev1alpha1.UpdateMachineTypeStatusResponse{
		Result: commonv1alpha1.RequestResult_REQUEST_RESULT_SUCCESS,
	}), nil
}

func (s *MachineTypeService) AddMachineGroup(
	ctx context.Context,
	c *connect.Request[machinetypev1alpha1.AddMachineGroupRequest],
) (*connect.Response[machinetypev1alpha1.AddMachineGroupResponse], error) {
	req := c.Msg
	namespace := req.GetNamespace()
	if namespace == "" {
		namespace = s.namespace
	}

	machinetype, err := s.c.LifecycleV1alpha1().MachineTypes(namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		errCode := connect.CodeInternal
		if apierrors.IsNotFound(err) {
			errCode = connect.CodeNotFound
		}
		return nil, connect.NewError(errCode, err)
	}
	mgroups := apiutil.MachineGroupsToGrpcAPI(machinetype.Spec.MachineGroups)
	if machineGroupIndex(req.MachineGroup.Name, mgroups) > -1 {
		return connect.NewResponse(&machinetypev1alpha1.AddMachineGroupResponse{
			Reason: AddMachineGroupFailureReason,
			Result: commonv1alpha1.RequestResult_REQUEST_RESULT_FAILURE,
		}), nil
	}
	mgroups = slices.Grow(mgroups, 1)
	mgroups = append(mgroups, req.MachineGroup)
	machinetypeApply := v1alpha1.MachineType(machinetype.Name, machinetype.Namespace).WithSpec(
		v1alpha1.MachineTypeSpec().WithMachineGroups(apiutil.MachineGroupsToApplyConfiguration(mgroups)...))
	if _, err = s.c.LifecycleV1alpha1().MachineTypes(namespace).Apply(ctx, machinetypeApply, metav1.ApplyOptions{
		FieldManager: "lifecycle.ironcore.dev/lifecycle-manager",
		Force:        true,
	}); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&machinetypev1alpha1.AddMachineGroupResponse{
		Result: commonv1alpha1.RequestResult_REQUEST_RESULT_SUCCESS,
	}), nil
}

func (s *MachineTypeService) RemoveMachineGroup(
	ctx context.Context,
	c *connect.Request[machinetypev1alpha1.RemoveMachineGroupRequest],
) (*connect.Response[machinetypev1alpha1.RemoveMachineGroupResponse], error) {
	req := c.Msg
	namespace := req.GetNamespace()
	if namespace == "" {
		namespace = s.namespace
	}

	machinetype, err := s.c.LifecycleV1alpha1().MachineTypes(namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		errCode := connect.CodeInternal
		if apierrors.IsNotFound(err) {
			errCode = connect.CodeNotFound
		}
		return nil, connect.NewError(errCode, err)
	}
	mgroups := apiutil.MachineGroupsToGrpcAPI(machinetype.Spec.MachineGroups)
	var idx int
	if idx = machineGroupIndex(req.GroupName, mgroups); idx == -1 {
		return connect.NewResponse(&machinetypev1alpha1.RemoveMachineGroupResponse{
			Result: commonv1alpha1.RequestResult_REQUEST_RESULT_SUCCESS,
		}), nil
	}
	mgroups = removeMachineGroup(mgroups, idx)
	machinetypeApply := v1alpha1.MachineType(machinetype.Name, machinetype.Namespace).WithSpec(
		v1alpha1.MachineTypeSpec().WithMachineGroups(apiutil.MachineGroupsToApplyConfiguration(mgroups)...))
	if _, err = s.c.LifecycleV1alpha1().MachineTypes(namespace).Apply(ctx, machinetypeApply, metav1.ApplyOptions{
		FieldManager: "lifecycle.ironcore.dev/lifecycle-manager",
		Force:        true,
	}); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&machinetypev1alpha1.RemoveMachineGroupResponse{
		Result: commonv1alpha1.RequestResult_REQUEST_RESULT_SUCCESS,
	}), nil
}

func machineGroupIndex(name string, dst []*machinetypev1alpha1.MachineGroup) int {
	return slices.IndexFunc(dst, func(g *machinetypev1alpha1.MachineGroup) bool {
		return name == g.Name
	})
}

func removeMachineGroup(src []*machinetypev1alpha1.MachineGroup, index int) []*machinetypev1alpha1.MachineGroup {
	result := make([]*machinetypev1alpha1.MachineGroup, 0, len(src)-1)
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
