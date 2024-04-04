// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	"context"
	"slices"
	"time"

	"connectrpc.com/connect"
	"github.com/go-logr/logr"
	lifecyclev1alpha1 "github.com/ironcore-dev/lifecycle-manager/api/lifecycle/v1alpha1"
	commonv1alpha1 "github.com/ironcore-dev/lifecycle-manager/api/proto/common/v1alpha1"
	machinetypev1alpha1 "github.com/ironcore-dev/lifecycle-manager/api/proto/machinetype/v1alpha1"
	"github.com/ironcore-dev/lifecycle-manager/clientgo/applyconfiguration/lifecycle/v1alpha1"
	"github.com/ironcore-dev/lifecycle-manager/clientgo/connectrpc/machinetype/v1alpha1/machinetypev1alpha1connect"
	"github.com/ironcore-dev/lifecycle-manager/clientgo/lifecycle"
	"github.com/ironcore-dev/lifecycle-manager/internal/service/scheduler"
	"github.com/ironcore-dev/lifecycle-manager/internal/util/apiutil"
	"github.com/ironcore-dev/lifecycle-manager/internal/util/uuidutil"
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
	scheduler *scheduler.Scheduler[*lifecyclev1alpha1.MachineType]
	horizon   time.Duration
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

func WithScheduler(scheduler *scheduler.Scheduler[*lifecyclev1alpha1.MachineType]) Option {
	return func(svc *MachineTypeService) {
		svc.scheduler = scheduler
	}
}

func (s *MachineTypeService) StartScheduler(ctx context.Context) {
	s.scheduler.Start(ctx)
}

func (s *MachineTypeService) ListMachineTypes(
	ctx context.Context,
	c *connect.Request[machinetypev1alpha1.ListMachineTypesRequest],
) (*connect.Response[machinetypev1alpha1.ListMachineTypesResponse], error) {
	log := logr.FromContextAsSlogLogger(ctx)
	log.Info("request", "request_body", c.Any())
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
	ctx context.Context,
	c *connect.Request[machinetypev1alpha1.ScanRequest],
) (*connect.Response[machinetypev1alpha1.ScanResponse], error) {
	log := logr.FromContextAsSlogLogger(ctx)
	log.Info("request", "request_body", c.Any())
	req := c.Msg
	namespace := req.GetNamespace()
	if namespace == "" {
		namespace = s.namespace
	}
	resp := &machinetypev1alpha1.ScanResponse{
		Result: commonv1alpha1.RequestResult_REQUEST_RESULT_UNSPECIFIED,
	}

	machineType, err := s.c.LifecycleV1alpha1().MachineTypes(namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		errCode := connect.CodeInternal
		if apierrors.IsNotFound(err) {
			errCode = connect.CodeNotFound
		}
		return nil, connect.NewError(errCode, err)
	}
	key := uuidutil.UUIDFromObjectKey(types.NamespacedName{Name: req.Name, Namespace: namespace})
	resp.Result = s.scheduler.Schedule(
		scheduler.NewTask[*lifecyclev1alpha1.MachineType](key, scheduler.ScanJob, machineType))
	return connect.NewResponse(resp), nil
}

func (s *MachineTypeService) UpdateMachineTypeStatus(
	ctx context.Context,
	c *connect.Request[machinetypev1alpha1.UpdateMachineTypeStatusRequest],
) (*connect.Response[machinetypev1alpha1.UpdateMachineTypeStatusResponse], error) {
	log := logr.FromContextAsSlogLogger(ctx)
	log.Info("request", "request_body", c.Any())
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
	key := uuidutil.UUIDFromObjectKey(types.NamespacedName{Name: req.Name, Namespace: namespace})
	s.scheduler.ForgetFinishedJob(key)
	return connect.NewResponse(&machinetypev1alpha1.UpdateMachineTypeStatusResponse{
		Result: commonv1alpha1.RequestResult_REQUEST_RESULT_SUCCESS,
	}), nil
}

func (s *MachineTypeService) AddMachineGroup(
	ctx context.Context,
	c *connect.Request[machinetypev1alpha1.AddMachineGroupRequest],
) (*connect.Response[machinetypev1alpha1.AddMachineGroupResponse], error) {
	log := logr.FromContextAsSlogLogger(ctx)
	log.Info("request", "request_body", c.Any())
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
	machineGroups := apiutil.MachineGroupsToGrpcAPI(machinetype.Spec.MachineGroups)
	if machineGroupIndex(req.MachineGroup.Name, machineGroups) > -1 {
		return connect.NewResponse(&machinetypev1alpha1.AddMachineGroupResponse{
			Reason: AddMachineGroupFailureReason,
			Result: commonv1alpha1.RequestResult_REQUEST_RESULT_FAILURE,
		}), nil
	}
	machineGroups = slices.Grow(machineGroups, 1)
	machineGroups = append(machineGroups, req.MachineGroup)
	machinetypeApply := v1alpha1.MachineType(machinetype.Name, machinetype.Namespace).WithSpec(
		v1alpha1.MachineTypeSpec().WithMachineGroups(apiutil.MachineGroupsToApplyConfiguration(machineGroups)...))
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
	log := logr.FromContextAsSlogLogger(ctx)
	log.Info("request", "request_body", c.Any())
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
	machineGroups := apiutil.MachineGroupsToGrpcAPI(machinetype.Spec.MachineGroups)
	var idx int
	if idx = machineGroupIndex(req.GroupName, machineGroups); idx == -1 {
		return connect.NewResponse(&machinetypev1alpha1.RemoveMachineGroupResponse{
			Result: commonv1alpha1.RequestResult_REQUEST_RESULT_SUCCESS,
		}), nil
	}
	machineGroups = removeMachineGroup(machineGroups, idx)
	machinetypeApply := v1alpha1.MachineType(machinetype.Name, machinetype.Namespace).WithSpec(
		v1alpha1.MachineTypeSpec().WithMachineGroups(apiutil.MachineGroupsToApplyConfiguration(machineGroups)...))
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
