// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package machine

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/rest"

	lifecyclev1alpha1 "github.com/ironcore-dev/lifecycle-manager/api/lifecycle/v1alpha1"
	"github.com/ironcore-dev/lifecycle-manager/clientgo/lifecycle"
	commonv1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/common/v1alpha1"
	machinev1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/machine/v1alpha1"
	"github.com/ironcore-dev/lifecycle-manager/util/apiutil"
)

type GrpcService struct {
	machinev1alpha1.UnimplementedMachineServiceServer
	cl         *lifecycle.Clientset
	cache      map[string]lifecyclev1alpha1.Machine
	horizon    time.Duration
	scanPeriod time.Duration
	namespace  string
}

func NewGrpcService(cfg *rest.Config) *GrpcService {
	cl := lifecycle.NewForConfigOrDie(cfg)
	return &GrpcService{cl: cl}
}

func (s *GrpcService) ListMachines(
	ctx context.Context,
	req *machinev1alpha1.ListMachinesRequest,
) (*machinev1alpha1.ListMachinesResponse, error) {
	log := logr.FromContextOrDiscard(ctx)
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
		return nil, err
	}
	resp := &machinev1alpha1.ListMachinesResponse{Machines: make([]*machinev1alpha1.Machine, len(machines.Items))}
	for i, item := range machines.Items {
		m := item.DeepCopy()
		resp.Machines[i] = apiutil.MachineTo(m)
	}
	return resp, nil
}

func (s *GrpcService) ScanMachine(
	ctx context.Context,
	req *machinev1alpha1.ScanMachineRequest,
) (*machinev1alpha1.ScanMachineResponse, error) {
	log := logr.FromContextOrDiscard(ctx)
	resp := &machinev1alpha1.ScanMachineResponse{
		Result: commonv1alpha1.RequestResult_REQUEST_RESULT_UNSPECIFIED,
		Status: nil,
	}
	// todo:
	//  - check whether machine is already cached
	//  - if machine is in cache, check whether last scan timestamp in horizon
	//  - if yes then send Success response with cached data
	//  - otherwise (machine is not in cache, last scan is outdated), send machine to scheduler

	_, err := s.cl.LifecycleV1alpha1().Machines(req.Namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		log.Error(err, "failed to get machine", "name", req.Name, "namespace", req.Namespace)
		resp.Result = commonv1alpha1.RequestResult_REQUEST_RESULT_FAILURE
		return resp, err
	}

	resp.Result = commonv1alpha1.RequestResult_REQUEST_RESULT_SCHEDULED
	return resp, nil
}

func (s *GrpcService) Install(
	ctx context.Context,
	req *machinev1alpha1.InstallRequest,
) (*machinev1alpha1.InstallResponse, error) {
	log := logr.FromContextOrDiscard(ctx)
	resp := &machinev1alpha1.InstallResponse{Result: commonv1alpha1.RequestResult_REQUEST_RESULT_UNSPECIFIED}

	// todo: send data to scheduler
	log.V(1).Info("packages installation scheduled for machine", "name", req.Name, "namespace", req.Namespace)
	resp.Result = commonv1alpha1.RequestResult_REQUEST_RESULT_SCHEDULED
	return resp, nil
}

func (s *GrpcService) UpdateMachineStatus(
	ctx context.Context,
	req *machinev1alpha1.UpdateMachineStatusRequest,
) (*machinev1alpha1.UpdateMachineStatusResponse, error) {
	// TODO implement me
	// invoke install with provided params
	err := status.Error(codes.Unimplemented, "UpdateMachine() is not implemented yet")
	return nil, err
}

func (s *GrpcService) AddPackageVersion(
	ctx context.Context,
	req *machinev1alpha1.AddPackageVersionRequest,
) (*machinev1alpha1.AddPackageVersionResponse, error) {
	// TODO implement me
	// invoke install with provided params
	err := status.Error(codes.Unimplemented, "AddPackageVersion() is not implemented yet")
	return nil, err
}

func (s *GrpcService) SetPackageVersion(
	ctx context.Context,
	req *machinev1alpha1.SetPackageVersionRequest,
) (*machinev1alpha1.SetPackageVersionResponse, error) {
	// TODO implement me
	// invoke install with provided params
	err := status.Error(codes.Unimplemented, "SetPackageVersion() is not implemented yet")
	return nil, err
}

func (s *GrpcService) RemovePackageVersion(
	ctx context.Context,
	req *machinev1alpha1.RemovePackageVersionRequest,
) (*machinev1alpha1.RemovePackageVersionResponse, error) {
	// TODO implement me
	// invoke install with provided params
	err := status.Error(codes.Unimplemented, "RemovePackageVersion() is not implemented yet")
	return nil, err
}
