// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package machinetype

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sigs.k8s.io/controller-runtime/pkg/client"

	machinetypev1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/machinetype/v1alpha1"
)

type GrpcService struct {
	machinetypev1alpha1.UnimplementedMachineTypeServiceServer
	client.Client
	cache map[string]string
}

func NewGRPCService() *GrpcService {
	return &GrpcService{}
}

func (s *GrpcService) ListMachineTypes(
	ctx context.Context,
	req *machinetypev1alpha1.ListMachineTypesRequest,
) (*machinetypev1alpha1.ListMachineTypesResponse, error) {
	// TODO implement me
	err := status.Error(codes.Unimplemented, "ListMachineTypes() is not implemented yet")
	return nil, err
}

func (s *GrpcService) Scan(
	ctx context.Context,
	req *machinetypev1alpha1.ScanRequest,
) (*machinetypev1alpha1.ScanResponse, error) {
	// TODO implement me
	// lookup cache for stored entry
	// if cache miss - invoke scan
	// if cache hit && timestamp out of horizon - invoke scan
	// otherwise send entry in response
	err := status.Error(codes.Unimplemented, "Scan() is not implemented yet")
	return nil, err
}

func (s *GrpcService) UpdateMachineType(
	ctx context.Context,
	req *machinetypev1alpha1.UpdateMachineTypeRequest,
) (*machinetypev1alpha1.UpdateMachineTypeResponse, error) {
	// TODO implement me
	err := status.Error(codes.Unimplemented, "UpdateMachineType() is not implemented yet")
	return nil, err
}
