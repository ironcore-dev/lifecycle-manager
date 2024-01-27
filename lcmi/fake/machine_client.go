// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package fake

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"k8s.io/apimachinery/pkg/types"

	commonv1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/common/v1alpha1"
	machinev1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/machine/v1alpha1"
	"github.com/ironcore-dev/lifecycle-manager/util/uuidutil"
)

type MachineClient struct {
	cache map[string]*machinev1alpha1.MachineStatus
}

func NewMachineClient(cache map[string]*machinev1alpha1.MachineStatus) *MachineClient {
	return &MachineClient{cache: cache}
}

func (m *MachineClient) WriteCache(id string, item *machinev1alpha1.MachineStatus) {
	m.cache[id] = item
}

func (m *MachineClient) ReadCache(id string) *machinev1alpha1.MachineStatus {
	return m.cache[id]
}

func (m *MachineClient) ListMachines(
	_ context.Context,
	_ *machinev1alpha1.ListMachinesRequest,
	_ ...grpc.CallOption,
) (*machinev1alpha1.ListMachinesResponse, error) {
	return nil, nil
}

func (m *MachineClient) ScanMachine(
	_ context.Context,
	in *machinev1alpha1.ScanMachineRequest,
	_ ...grpc.CallOption,
) (*machinev1alpha1.ScanMachineResponse, error) {
	if in.Name == "failed-scan" {
		return nil, errors.New("fake error")
	}
	key := types.NamespacedName{Name: in.Name, Namespace: in.Namespace}
	uid := uuidutil.UUIDFromObjectKey(key)
	entry, ok := m.cache[uid]
	if ok {
		return &machinev1alpha1.ScanMachineResponse{
			Status: entry,
			Result: commonv1alpha1.RequestResult_REQUEST_RESULT_SUCCESS,
		}, nil
	}
	m.cache[uid] = &machinev1alpha1.MachineStatus{}
	return &machinev1alpha1.ScanMachineResponse{
		Status: nil,
		Result: commonv1alpha1.RequestResult_REQUEST_RESULT_SCHEDULED,
	}, nil
}

func (m *MachineClient) Install(
	_ context.Context,
	in *machinev1alpha1.InstallRequest,
	_ ...grpc.CallOption,
) (*machinev1alpha1.InstallResponse, error) {
	if in.Name == "failed-install" {
		return nil, errors.New("fake error")
	}
	if in.Name == "fake-failure" {
		return &machinev1alpha1.InstallResponse{Result: commonv1alpha1.RequestResult_REQUEST_RESULT_FAILURE}, nil
	}
	key := types.NamespacedName{Name: in.Name, Namespace: in.Namespace}
	uid := uuidutil.UUIDFromObjectKey(key)
	_, ok := m.cache[uid]
	if ok {
		return &machinev1alpha1.InstallResponse{Result: commonv1alpha1.RequestResult_REQUEST_RESULT_SCHEDULED}, nil
	}
	m.cache[uid] = &machinev1alpha1.MachineStatus{}
	return &machinev1alpha1.InstallResponse{Result: commonv1alpha1.RequestResult_REQUEST_RESULT_UNSPECIFIED}, nil
}

func (m *MachineClient) UpdateMachineStatus(
	_ context.Context,
	_ *machinev1alpha1.UpdateMachineStatusRequest,
	_ ...grpc.CallOption,
) (*machinev1alpha1.UpdateMachineStatusResponse, error) {
	return nil, nil
}

func (m *MachineClient) AddPackageVersion(
	_ context.Context,
	_ *machinev1alpha1.AddPackageVersionRequest,
	_ ...grpc.CallOption,
) (*machinev1alpha1.AddPackageVersionResponse, error) {
	return nil, nil
}

func (m *MachineClient) SetPackageVersion(
	_ context.Context,
	_ *machinev1alpha1.SetPackageVersionRequest,
	_ ...grpc.CallOption,
) (*machinev1alpha1.SetPackageVersionResponse, error) {
	return nil, nil
}

func (m *MachineClient) RemovePackageVersion(
	_ context.Context,
	_ *machinev1alpha1.RemovePackageVersionRequest,
	_ ...grpc.CallOption,
) (*machinev1alpha1.RemovePackageVersionResponse, error) {
	return nil, nil
}
