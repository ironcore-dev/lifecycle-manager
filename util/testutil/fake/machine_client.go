// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package fake

import (
	"context"

	"connectrpc.com/connect"
	commonv1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/common/v1alpha1"
	machineapiv1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/machine/v1alpha1"
	"github.com/ironcore-dev/lifecycle-manager/util/uuidutil"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/types"
)

type MachineClient struct {
	cache map[string]*machineapiv1alpha1.MachineStatus
}

func NewMachineClient(cache map[string]*machineapiv1alpha1.MachineStatus) *MachineClient {
	return &MachineClient{cache: cache}
}

func (c *MachineClient) WriteCache(id string, item *machineapiv1alpha1.MachineStatus) {
	c.cache[id] = item
}

func (c *MachineClient) ReadCache(id string) *machineapiv1alpha1.MachineStatus {
	return c.cache[id]
}

func (c *MachineClient) ScanMachine(
	_ context.Context,
	req *connect.Request[machineapiv1alpha1.ScanMachineRequest],
) (*connect.Response[machineapiv1alpha1.ScanMachineResponse], error) {
	in := req.Msg
	if in.Name == "failed-scan" {
		return nil, errors.New("fake error")
	}
	key := types.NamespacedName{Name: in.Name, Namespace: in.Namespace}
	uid := uuidutil.UUIDFromObjectKey(key)
	_, ok := c.cache[uid]
	if ok {
		return connect.NewResponse(&machineapiv1alpha1.ScanMachineResponse{
			Result: commonv1alpha1.RequestResult_REQUEST_RESULT_SUCCESS,
		}), nil
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
	if in.Name == "failed-install" {
		return nil, connect.NewError(connect.CodeInternal, errors.New("fake error"))
	}
	if in.Name == "fake-failure" {
		return connect.NewResponse(&machineapiv1alpha1.InstallResponse{
			Result: commonv1alpha1.RequestResult_REQUEST_RESULT_FAILURE,
		}), nil
	}
	key := types.NamespacedName{Name: in.Name, Namespace: in.Namespace}
	uid := uuidutil.UUIDFromObjectKey(key)
	_, ok := c.cache[uid]
	if ok {
		return connect.NewResponse(&machineapiv1alpha1.InstallResponse{
			Result: commonv1alpha1.RequestResult_REQUEST_RESULT_SCHEDULED,
		}), nil
	}
	c.cache[uid] = &machineapiv1alpha1.MachineStatus{}
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
