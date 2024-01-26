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
	machinev1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/machine/v1alpha1"
	"github.com/ironcore-dev/lifecycle-manager/util/apiutil"
)

type GrpcService struct {
	machinev1alpha1.UnimplementedMachineServiceServer
	cl        *lifecycle.Clientset
	cache     map[string]lifecyclev1alpha1.Machine
	horizon   time.Duration
	namespace string
}

func NewGRPCService(cfg *rest.Config) *GrpcService {
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
	// TODO implement me
	// lookup cache for stored entry
	// if cache miss - get machine and invoke scan providing necessary info
	// if cache hit && timestamp out of horizon - invoke scan
	// otherwise send entry in response
	err := status.Error(codes.Unimplemented, "Scan() is not implemented yet")
	return nil, err
}

func (s *GrpcService) Install(
	ctx context.Context,
	req *machinev1alpha1.InstallRequest,
) (*machinev1alpha1.InstallResponse, error) {
	// TODO implement me
	// invoke install with provided params
	err := status.Error(codes.Unimplemented, "Install() is not implemented yet")
	return nil, err
}

func (s *GrpcService) UpdateMachine(
	ctx context.Context,
	req *machinev1alpha1.UpdateMachineRequest,
) (*machinev1alpha1.UpdateMachineResponse, error) {
	// TODO implement me
	// invoke install with provided params
	err := status.Error(codes.Unimplemented, "UpdateMachine() is not implemented yet")
	return nil, err
}
