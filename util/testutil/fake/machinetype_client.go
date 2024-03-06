// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package fake

import (
	"context"

	"connectrpc.com/connect"
	commonv1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/common/v1alpha1"
	machinetypeapiv1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/machinetype/v1alpha1"
	"github.com/ironcore-dev/lifecycle-manager/util/uuidutil"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/types"
)

type MachineTypeClient struct {
	cache map[string]*machinetypeapiv1alpha1.MachineTypeStatus
}

func NewMachineTypeClient(cache map[string]*machinetypeapiv1alpha1.MachineTypeStatus) *MachineTypeClient {
	return &MachineTypeClient{cache: cache}
}

func (c *MachineTypeClient) WriteCache(id string, item *machinetypeapiv1alpha1.MachineTypeStatus) {
	c.cache[id] = item
}

func (c *MachineTypeClient) ReadCache(id string) *machinetypeapiv1alpha1.MachineTypeStatus {
	return c.cache[id]
}

func (c *MachineTypeClient) ListMachineTypes(
	_ context.Context,
	_ *connect.Request[machinetypeapiv1alpha1.ListMachineTypesRequest],
) (*connect.Response[machinetypeapiv1alpha1.ListMachineTypesResponse], error) {
	return nil, nil
}

func (c *MachineTypeClient) Scan(
	ctx context.Context,
	req *connect.Request[machinetypeapiv1alpha1.ScanRequest],
) (*connect.Response[machinetypeapiv1alpha1.ScanResponse], error) {
	in := req.Msg
	if in.Name == "failed-scan" {
		return nil, connect.NewError(connect.CodeInternal, errors.New("fake error"))
	}
	key := types.NamespacedName{Name: in.Name, Namespace: in.Namespace}
	uid := uuidutil.UUIDFromObjectKey(key)
	_, ok := c.cache[uid]
	if ok {
		return connect.NewResponse(&machinetypeapiv1alpha1.ScanResponse{
			Result: commonv1alpha1.RequestResult_REQUEST_RESULT_SUCCESS,
		}), nil
	}
	c.cache[uid] = &machinetypeapiv1alpha1.MachineTypeStatus{}
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
