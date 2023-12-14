package fake

import (
	"context"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	lcmi "github.com/ironcore-dev/lifecycle-manager/lcmi/api/machine_type/v1alpha1"
)

type MachineTypeClient struct {
	sync.Mutex

	scans map[string]*lcmi.ScanResponse
}

func NewFakeMachineTypeClient() *MachineTypeClient {
	return &MachineTypeClient{scans: make(map[string]*lcmi.ScanResponse)}
}

func NewFakeMachineTypeClientWithScans(scans map[string]*lcmi.ScanResponse) *MachineTypeClient {
	return &MachineTypeClient{scans: scans}
}

func (f *MachineTypeClient) ListMachineTypes(_ context.Context, _ *lcmi.ListMachineTypesRequest, _ ...grpc.CallOption) (*lcmi.ListMachineTypesResponse, error) {
	return nil, nil
}

func (f *MachineTypeClient) Scan(_ context.Context, in *lcmi.ScanRequest, _ ...grpc.CallOption) (*lcmi.ScanResponse, error) {
	f.Lock()
	defer f.Unlock()

	stored, ok := f.scans[in.Id]
	if !ok {
		return nil, status.Error(codes.NotFound, "no scans were done for current object yet")
	}

	return stored, nil
}

func (f *MachineTypeClient) Status(_ context.Context, in *lcmi.StatusRequest, _ ...grpc.CallOption) (*lcmi.StatusResponse, error) {
	f.Lock()
	defer f.Unlock()

	stored, ok := f.scans[in.Id]
	if !ok {
		return nil, status.Error(codes.NotFound, "scan result not found")
	}

	return &lcmi.StatusResponse{
		Status: stored.Status,
	}, nil
}
