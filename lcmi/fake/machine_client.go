// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package fake

import (
	"context"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	lcmimachine "github.com/ironcore-dev/lifecycle-manager/lcmi/api/machine/v1alpha1"
)

type MachineClient struct {
	sync.Mutex

	scans map[string]*lcmimachine.ScanResponse
}

func NewMachineClient() *MachineClient {
	return &MachineClient{scans: make(map[string]*lcmimachine.ScanResponse)}
}

func NewMachineClientWithScans(scans map[string]*lcmimachine.ScanResponse) *MachineClient {
	return &MachineClient{scans: scans}
}

func (f *MachineClient) ListMachines(
	_ context.Context,
	_ *lcmimachine.ListMachinesRequest,
	_ ...grpc.CallOption,
) (*lcmimachine.ListMachinesResponse, error) {
	return nil, nil
}

func (f *MachineClient) Scan(
	_ context.Context,
	in *lcmimachine.ScanRequest,
	_ ...grpc.CallOption,
) (*lcmimachine.ScanResponse, error) {
	f.Lock()
	defer f.Unlock()

	stored, ok := f.scans[in.Id]
	if !ok {
		return nil, status.Error(codes.NotFound, "no scans were done for current object yet")
	}

	return stored, nil
}

func (f *MachineClient) Status(
	_ context.Context,
	in *lcmimachine.StatusRequest,
	_ ...grpc.CallOption,
) (*lcmimachine.StatusResponse, error) {
	f.Lock()
	defer f.Unlock()

	stored, ok := f.scans[in.Id]
	if !ok {
		return nil, status.Error(codes.NotFound, "scan result not found")
	}

	return &lcmimachine.StatusResponse{
		Status: stored.Status,
	}, nil
}
