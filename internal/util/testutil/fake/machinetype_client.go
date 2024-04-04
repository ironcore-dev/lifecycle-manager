// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package fake

import (
	"context"

	"connectrpc.com/connect"
	commonv1alpha1 "github.com/ironcore-dev/lifecycle-manager/api/proto/common/v1alpha1"
	machinetypeapiv1alpha1 "github.com/ironcore-dev/lifecycle-manager/api/proto/machinetype/v1alpha1"
	"github.com/pkg/errors"
)

type MachineTypeClient struct {
}

func NewMachineTypeClient() *MachineTypeClient {
	return &MachineTypeClient{}
}

func (c *MachineTypeClient) ListMachineTypes(
	_ context.Context,
	_ *connect.Request[machinetypeapiv1alpha1.ListMachineTypesRequest],
) (*connect.Response[machinetypeapiv1alpha1.ListMachineTypesResponse], error) {
	return nil, nil
}

func (c *MachineTypeClient) Scan(
	_ context.Context,
	req *connect.Request[machinetypeapiv1alpha1.ScanRequest],
) (*connect.Response[machinetypeapiv1alpha1.ScanResponse], error) {
	in := req.Msg
	switch {
	case in.Name == "failed-scan":
		return nil, connect.NewError(connect.CodeInternal, errors.New("fake error"))
	case in.Name == "sample-scan-submitted":
		return connect.NewResponse(&machinetypeapiv1alpha1.ScanResponse{
			Result: commonv1alpha1.RequestResult_REQUEST_RESULT_SUCCESS,
		}), nil
	}
	return connect.NewResponse(&machinetypeapiv1alpha1.ScanResponse{
		Result: commonv1alpha1.RequestResult_REQUEST_RESULT_SCHEDULED,
	}), nil
}

func (c *MachineTypeClient) UpdateMachineTypeStatus(
	_ context.Context,
	_ *connect.Request[machinetypeapiv1alpha1.UpdateMachineTypeStatusRequest],
) (*connect.Response[machinetypeapiv1alpha1.UpdateMachineTypeStatusResponse], error) {
	return nil, nil
}

func (c *MachineTypeClient) AddMachineGroup(
	_ context.Context,
	_ *connect.Request[machinetypeapiv1alpha1.AddMachineGroupRequest],
) (*connect.Response[machinetypeapiv1alpha1.AddMachineGroupResponse], error) {
	return nil, nil
}

func (c *MachineTypeClient) RemoveMachineGroup(
	_ context.Context,
	_ *connect.Request[machinetypeapiv1alpha1.RemoveMachineGroupRequest],
) (*connect.Response[machinetypeapiv1alpha1.RemoveMachineGroupResponse], error) {
	return nil, nil
}
