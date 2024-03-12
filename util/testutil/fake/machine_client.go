// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package fake

import (
	"context"

	"connectrpc.com/connect"
	commonv1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/common/v1alpha1"
	machineapiv1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/machine/v1alpha1"
	"github.com/pkg/errors"
)

type MachineClient struct {
}

func NewMachineClient() *MachineClient {
	return &MachineClient{}
}

func (c *MachineClient) ScanMachine(
	_ context.Context,
	req *connect.Request[machineapiv1alpha1.ScanMachineRequest],
) (*connect.Response[machineapiv1alpha1.ScanMachineResponse], error) {
	in := req.Msg
	switch {
	case in.Name == "sample-scan-submitted":
		return connect.NewResponse(&machineapiv1alpha1.ScanMachineResponse{
			Result: commonv1alpha1.RequestResult_REQUEST_RESULT_SUCCESS,
		}), nil
	case in.Name == "failed-scan":
		return nil, errors.New("fake error")
	}
	return connect.NewResponse(&machineapiv1alpha1.ScanMachineResponse{
		Result: commonv1alpha1.RequestResult_REQUEST_RESULT_SCHEDULED,
	}), nil
}

func (c *MachineClient) Install(
	_ context.Context,
	req *connect.Request[machineapiv1alpha1.InstallRequest],
) (*connect.Response[machineapiv1alpha1.InstallResponse], error) {
	in := req.Msg
	switch {
	case in.Name == "sample-install-submitted":
		return connect.NewResponse(&machineapiv1alpha1.InstallResponse{
			Result: commonv1alpha1.RequestResult_REQUEST_RESULT_SCHEDULED,
		}), nil
	case in.Name == "failed-install":
		return nil, connect.NewError(connect.CodeInternal, errors.New("fake error"))
	case in.Name == "fake-failure":
		return connect.NewResponse(&machineapiv1alpha1.InstallResponse{
			Result: commonv1alpha1.RequestResult_REQUEST_RESULT_FAILURE,
		}), nil
	}
	return connect.NewResponse(&machineapiv1alpha1.InstallResponse{
		Result: commonv1alpha1.RequestResult_REQUEST_RESULT_UNSPECIFIED,
	}), nil
}

func (c *MachineClient) UpdateMachineStatus(
	_ context.Context,
	_ *connect.Request[machineapiv1alpha1.UpdateMachineStatusRequest],
) (*connect.Response[machineapiv1alpha1.UpdateMachineStatusResponse], error) {
	return nil, nil
}

func (c *MachineClient) ListMachines(
	_ context.Context,
	_ *connect.Request[machineapiv1alpha1.ListMachinesRequest],
) (*connect.Response[machineapiv1alpha1.ListMachinesResponse], error) {
	return nil, nil
}

func (c *MachineClient) AddPackageVersion(
	_ context.Context,
	_ *connect.Request[machineapiv1alpha1.AddPackageVersionRequest],
) (*connect.Response[machineapiv1alpha1.AddPackageVersionResponse], error) {
	return nil, nil
}

func (c *MachineClient) SetPackageVersion(
	_ context.Context,
	_ *connect.Request[machineapiv1alpha1.SetPackageVersionRequest],
) (*connect.Response[machineapiv1alpha1.SetPackageVersionResponse], error) {
	return nil, nil
}

func (c *MachineClient) RemovePackageVersion(
	_ context.Context,
	_ *connect.Request[machineapiv1alpha1.RemovePackageVersionRequest],
) (*connect.Response[machineapiv1alpha1.RemovePackageVersionResponse], error) {
	return nil, nil
}
