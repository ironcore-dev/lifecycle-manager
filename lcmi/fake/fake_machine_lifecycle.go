// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package fake

import (
	"context"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	lcmimachinelifecycle "github.com/ironcore-dev/lifecycle-manager/lcmi/api/machine_lifecycle/v1alpha1"
)

type MachineLifecycleClient struct {
	sync.Mutex

	scans map[string]*lcmimachinelifecycle.ScanResponse
}

func NewMachineLifecycleClient() *MachineLifecycleClient {
	return &MachineLifecycleClient{scans: make(map[string]*lcmimachinelifecycle.ScanResponse)}
}

func NewMachineLifecycleClientWithScans(scans map[string]*lcmimachinelifecycle.ScanResponse) *MachineLifecycleClient {
	return &MachineLifecycleClient{scans: scans}
}

func (f *MachineLifecycleClient) ListMachineLifecycles(_ context.Context, _ *lcmimachinelifecycle.ListMachineLifecyclesRequest, _ ...grpc.CallOption) (*lcmimachinelifecycle.ListMachineLifecyclesResponse, error) {
	return nil, nil
}

func (f *MachineLifecycleClient) Scan(_ context.Context, in *lcmimachinelifecycle.ScanRequest, _ ...grpc.CallOption) (*lcmimachinelifecycle.ScanResponse, error) {
	f.Lock()
	defer f.Unlock()

	stored, ok := f.scans[in.Id]
	if !ok {
		return nil, status.Error(codes.NotFound, "no scans were done for current object yet")
	}

	return stored, nil
}

func (f *MachineLifecycleClient) Status(_ context.Context, in *lcmimachinelifecycle.StatusRequest, _ ...grpc.CallOption) (*lcmimachinelifecycle.StatusResponse, error) {
	f.Lock()
	defer f.Unlock()

	stored, ok := f.scans[in.Id]
	if !ok {
		return nil, status.Error(codes.NotFound, "scan result not found")
	}

	return &lcmimachinelifecycle.StatusResponse{
		Status: stored.Status,
	}, nil
}
