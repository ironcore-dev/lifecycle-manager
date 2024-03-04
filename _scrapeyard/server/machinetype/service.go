// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package machinetype

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	"github.com/ironcore-dev/lifecycle-manager/clientgo/applyconfiguration/lifecycle/v1alpha1"
	"github.com/ironcore-dev/lifecycle-manager/clientgo/lifecycle"
	commonv1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/common/v1alpha1"
	machinetypev1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/machinetype/v1alpha1"
	"github.com/ironcore-dev/lifecycle-manager/util/apiutil"
	"github.com/ironcore-dev/lifecycle-manager/util/uuidutil"
	"github.com/jellydator/ttlcache/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
)

type Scheduler interface {
	ScanScheduled(key types.NamespacedName) bool
	ScheduleScan(key types.NamespacedName, next time.Duration)
	DropJobFromCache(key types.NamespacedName)
}

type GrpcService struct {
	machinetypev1alpha1.UnimplementedMachineTypeServiceServer
	cl         *lifecycle.Clientset
	scheduler  Scheduler
	cache      *ttlcache.Cache[string, *machinetypev1alpha1.MachineTypeStatus]
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

func WithCache(cache *ttlcache.Cache[string, *machinetypev1alpha1.MachineTypeStatus]) ServiceOption {
	return func(svc *GrpcService) {
		svc.cache = cache
	}
}

func (s *GrpcService) ListMachineTypes(
	ctx context.Context,
	req *machinetypev1alpha1.ListMachineTypesRequest,
) (*machinetypev1alpha1.ListMachineTypesResponse, error) {
	log := logr.FromContextOrDiscard(ctx)
	namespace := req.GetNamespace()
	if namespace == "" {
		namespace = s.namespace
	}

	opts := metav1.ListOptions{}
	reqSelector := req.GetLabelSelector()
	if reqSelector != nil {
		opts.LabelSelector = labels.Set(reqSelector.MatchLabels).String()
	}
	machinetypes, err := s.cl.LifecycleV1alpha1().MachineTypes(namespace).List(ctx, opts)
	if err != nil {
		log.Error(err, "failed to list machine types")
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}
	resp := &machinetypev1alpha1.ListMachineTypesResponse{MachineTypes: make([]*machinetypev1alpha1.MachineType,
		len(machinetypes.Items))}
	for i, item := range machinetypes.Items {
		m := item.DeepCopy()
		resp.MachineTypes[i] = apiutil.MachineTypeToGrpcAPI(m)
	}
	return nil, err
}

func (s *GrpcService) Scan(
	ctx context.Context,
	req *machinetypev1alpha1.ScanRequest,
) (*machinetypev1alpha1.ScanResponse, error) {
	resp := &machinetypev1alpha1.ScanResponse{
		Result: commonv1alpha1.RequestResult_REQUEST_RESULT_UNSPECIFIED,
	}
	namespace := req.Namespace
	if namespace == "" {
		namespace = s.namespace
	}
	log := logr.FromContextOrDiscard(ctx).WithValues("name", req.Name, "namespace", namespace)

	key := types.NamespacedName{Name: req.Name, Namespace: namespace}
	cacheItem := s.cache.Get(uuidutil.UUIDFromObjectKey(key))

	switch {
	case cacheItem == nil:
		fallthrough
	case time.Since(cacheItem.Value().LastScanTime.Time) > s.horizon:
		log.V(1).Info("scan scheduled")
		if !s.scheduler.ScanScheduled(key) {
			s.scheduler.ScheduleScan(key, s.scanPeriod)
		}
		resp.Result = commonv1alpha1.RequestResult_REQUEST_RESULT_SCHEDULED
	default:
		resp.Result = commonv1alpha1.RequestResult_REQUEST_RESULT_SUCCESS
	}

	return resp, nil
}

func (s *GrpcService) UpdateMachineTypeStatus(
	ctx context.Context,
	req *machinetypev1alpha1.UpdateMachineTypeStatusRequest,
) (*machinetypev1alpha1.UpdateMachineTypeStatusResponse, error) {
	namespace := req.Namespace
	if namespace == "" {
		namespace = s.namespace
	}
	log := logr.FromContextOrDiscard(ctx).WithValues("name", req.Name, "namespace", namespace)
	machinetype, err := s.cl.LifecycleV1alpha1().MachineTypes(namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		log.Error(err, "failed to get machine type")
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}
	machinetypeApply := v1alpha1.MachineType(machinetype.Name, machinetype.Namespace).
		WithStatus(apiutil.MachineTypeStatusToApplyConfiguration(req.Status))
	if _, err = s.cl.LifecycleV1alpha1().MachineTypes(namespace).ApplyStatus(ctx, machinetypeApply, metav1.ApplyOptions{
		FieldManager: "lifecycle.ironcore.dev/lifecycle-manager",
		Force:        true,
	}); err != nil {
		log.Error(err, "failed to update machine type status")
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}
	key := types.NamespacedName{Name: req.Name, Namespace: namespace}
	s.cache.Set(uuidutil.UUIDFromObjectKey(key), req.Status, ttlcache.DefaultTTL)
	s.scheduler.DropJobFromCache(key)
	return &machinetypev1alpha1.UpdateMachineTypeStatusResponse{
		Result: commonv1alpha1.RequestResult_REQUEST_RESULT_SUCCESS,
	}, nil
}

func (s *GrpcService) AddMachineGroup(
	ctx context.Context,
	req *machinetypev1alpha1.AddMachineGroupRequest,
) (*machinetypev1alpha1.AddMachineGroupResponse, error) {
	// TODO implement me
	err := status.Error(codes.Unimplemented, "AddMachineGroup() is not implemented yet")
	return nil, err
}

func (s *GrpcService) RemoveMachineGroup(
	ctx context.Context,
	req *machinetypev1alpha1.RemoveMachineGroupRequest,
) (*machinetypev1alpha1.RemoveMachineGroupResponse, error) {
	// TODO implement me
	err := status.Error(codes.Unimplemented, "RemoveMachineGroup() is not implemented yet")
	return nil, err
}
