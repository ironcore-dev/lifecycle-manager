// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	"context"
	"reflect"
	"time"

	"connectrpc.com/connect"
	"github.com/go-logr/logr"
	lifecyclev1alpha1 "github.com/ironcore-dev/lifecycle-manager/api/lifecycle/v1alpha1"
	"github.com/ironcore-dev/lifecycle-manager/clientgo/lifecycle"
	commonv1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/common/v1alpha1"
	machinev1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/machine/v1alpha1"
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

func (s MachineService) ScanMachine(ctx context.Context, c *connect.Request[machinev1alpha1.ScanMachineRequest]) (*connect.Response[machinev1alpha1.ScanMachineResponse], error) {
	req := c.Msg
	log := logr.FromContextAsSlogLogger(ctx)
	namespace := req.GetNamespace()
	if namespace == "" {
		namespace = s.namespace
	}
	resp := &machinev1alpha1.ScanMachineResponse{
		Result: commonv1alpha1.RequestResult_REQUEST_RESULT_UNSPECIFIED,
	}

	machine, err := s.c.LifecycleV1alpha1().Machines(namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		log.Error(err.Error())
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}

	key := types.NamespacedName{Name: req.Name, Namespace: namespace}
	cacheItem := s.cache.Get(uuidutil.UUIDFromObjectKey(key))

	switch {
	case cacheItem == nil:
		fallthrough
	case time.Since(cacheItem.Value().LastScanTime.Time) > s.horizon:
		fallthrough
	case !installedPackagesEqual(machine.Status, cacheItem.Value()):
		log.Info("scan scheduled")
		resp.Result = commonv1alpha1.RequestResult_REQUEST_RESULT_SCHEDULED
	default:
		resp.Result = commonv1alpha1.RequestResult_REQUEST_RESULT_SUCCESS
	}
	return connect.NewResponse(resp), nil
}

func (s MachineService) Install(ctx context.Context, c *connect.Request[machinev1alpha1.InstallRequest]) (*connect.Response[machinev1alpha1.InstallResponse], error) {
	// TODO implement me
	panic("implement me")
}

func (s MachineService) UpdateMachineStatus(ctx context.Context, c *connect.Request[machinev1alpha1.UpdateMachineStatusRequest]) (*connect.Response[machinev1alpha1.UpdateMachineStatusResponse], error) {
	// TODO implement me
	panic("implement me")
}

func (s MachineService) ListMachines(ctx context.Context, c *connect.Request[machinev1alpha1.ListMachinesRequest]) (*connect.Response[machinev1alpha1.ListMachinesResponse], error) {
	req := c.Msg
	log := logr.FromContextAsSlogLogger(ctx)
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
		log.Error(err.Error())
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}
	resp := &machinev1alpha1.ListMachinesResponse{Machines: make([]*machinev1alpha1.Machine, len(machines.Items))}
	for i, item := range machines.Items {
		m := item.DeepCopy()
		resp.Machines[i] = apiutil.MachineToGrpcAPI(m)
	}
	return connect.NewResponse(resp), nil
}

func (s MachineService) AddPackageVersion(ctx context.Context, c *connect.Request[machinev1alpha1.AddPackageVersionRequest]) (*connect.Response[machinev1alpha1.AddPackageVersionResponse], error) {
	// TODO implement me
	panic("implement me")
}

func (s MachineService) SetPackageVersion(ctx context.Context, c *connect.Request[machinev1alpha1.SetPackageVersionRequest]) (*connect.Response[machinev1alpha1.SetPackageVersionResponse], error) {
	// TODO implement me
	panic("implement me")
}

func (s MachineService) RemovePackageVersion(ctx context.Context, c *connect.Request[machinev1alpha1.RemovePackageVersionRequest]) (*connect.Response[machinev1alpha1.RemovePackageVersionResponse], error) {
	// TODO implement me
	panic("implement me")
}

func installedPackagesEqual(
	src lifecyclev1alpha1.MachineStatus,
	tgt *machinev1alpha1.MachineStatus,
) bool {
	conv := apiutil.MachineStatusToKubeAPI(tgt)
	return reflect.DeepEqual(src, conv)
}
