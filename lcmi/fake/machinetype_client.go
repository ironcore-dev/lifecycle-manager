// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package fake

import (
	"context"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	lcmimachinetype "github.com/ironcore-dev/lifecycle-manager/lcmi/api/machinetype/v1alpha1"
)

type MachineTypeClient struct {
	sync.Mutex

	scans map[string]*lcmimachinetype.ScanResponse
}

func NewFakeMachineTypeClient() *MachineTypeClient {
	return &MachineTypeClient{scans: make(map[string]*lcmimachinetype.ScanResponse)}
}

func NewFakeMachineTypeClientWithScans(scans map[string]*lcmimachinetype.ScanResponse) *MachineTypeClient {
	return &MachineTypeClient{scans: scans}
}

func (f *MachineTypeClient) ListMachineTypes(
	_ context.Context,
	_ *lcmimachinetype.ListMachineTypesRequest,
	_ ...grpc.CallOption,
) (*lcmimachinetype.ListMachineTypesResponse, error) {
	return nil, nil
}

func (f *MachineTypeClient) Scan(
	_ context.Context,
	in *lcmimachinetype.ScanRequest,
	_ ...grpc.CallOption,
) (*lcmimachinetype.ScanResponse, error) {
	f.Lock()
	defer f.Unlock()

	stored, ok := f.scans[in.Id]
	if !ok {
		return nil, status.Error(codes.NotFound, "no scans were done for current object yet")
	}

	return stored, nil
}

func (f *MachineTypeClient) Status(
	_ context.Context,
	in *lcmimachinetype.StatusRequest,
	_ ...grpc.CallOption,
) (*lcmimachinetype.StatusResponse, error) {
	f.Lock()
	defer f.Unlock()

	stored, ok := f.scans[in.Id]
	if !ok {
		return nil, status.Error(codes.NotFound, "scan result not found")
	}

	return &lcmimachinetype.StatusResponse{
		Status: stored.Status,
	}, nil
}
