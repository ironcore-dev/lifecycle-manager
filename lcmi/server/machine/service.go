package machine

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sigs.k8s.io/controller-runtime/pkg/client"

	machinev1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/machine/v1alpha1"
)

type GrpcService struct {
	client.Client
	cache map[string]string
}

func NewGRPCService() *GrpcService {
	return &GrpcService{}
}

func (s *GrpcService) ListMachines(
	ctx context.Context,
	req *machinev1alpha1.ListMachinesRequest,
) (*machinev1alpha1.ListMachinesResponse, error) {
	// TODO implement me
	err := status.Error(codes.Unimplemented, "ListMachines() is not implemented yet")
	return nil, err
}

func (s *GrpcService) Scan(
	ctx context.Context,
	req *machinev1alpha1.ScanRequest,
) (*machinev1alpha1.ScanResponse, error) {
	// TODO implement me
	err := status.Error(codes.Unimplemented, "Scan() is not implemented yet")
	return nil, err
}

func (s *GrpcService) Install(
	ctx context.Context,
	req *machinev1alpha1.InstallRequest,
) (*machinev1alpha1.InstallResponse, error) {
	// TODO implement me
	err := status.Error(codes.Unimplemented, "Install() is not implemented yet")
	return nil, err
}
