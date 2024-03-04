// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	"context"

	"connectrpc.com/connect"
	machinetypev1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/machinetype/v1alpha1"
)

type MachineTypeService struct{}

func (m MachineTypeService) ListMachineTypes(ctx context.Context, c *connect.Request[machinetypev1alpha1.ListMachineTypesRequest]) (*connect.Response[machinetypev1alpha1.ListMachineTypesResponse], error) {
	// TODO implement me
	panic("implement me")
}

func (m MachineTypeService) Scan(ctx context.Context, c *connect.Request[machinetypev1alpha1.ScanRequest]) (*connect.Response[machinetypev1alpha1.ScanResponse], error) {
	// TODO implement me
	panic("implement me")
}

func (m MachineTypeService) UpdateMachineTypeStatus(ctx context.Context, c *connect.Request[machinetypev1alpha1.UpdateMachineTypeStatusRequest]) (*connect.Response[machinetypev1alpha1.UpdateMachineTypeStatusResponse], error) {
	// TODO implement me
	panic("implement me")
}

func (m MachineTypeService) AddMachineGroup(ctx context.Context, c *connect.Request[machinetypev1alpha1.AddMachineGroupRequest]) (*connect.Response[machinetypev1alpha1.AddMachineGroupResponse], error) {
	// TODO implement me
	panic("implement me")
}

func (m MachineTypeService) RemoveMachineGroup(ctx context.Context, c *connect.Request[machinetypev1alpha1.RemoveMachineGroupRequest]) (*connect.Response[machinetypev1alpha1.RemoveMachineGroupResponse], error) {
	// TODO implement me
	panic("implement me")
}
