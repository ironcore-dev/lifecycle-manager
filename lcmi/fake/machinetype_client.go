// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package fake

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"k8s.io/apimachinery/pkg/types"

	commonv1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/common/v1alpha1"
	machinetypev1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/machinetype/v1alpha1"
	"github.com/ironcore-dev/lifecycle-manager/util/uuidutil"
)

type MachineTypeClient struct {
	cache map[string]*machinetypev1alpha1.MachineTypeStatus
}

func NewMachineTypeClient(cache map[string]*machinetypev1alpha1.MachineTypeStatus) *MachineTypeClient {
	return &MachineTypeClient{cache: cache}
}

func (m *MachineTypeClient) WriteCache(id string, item *machinetypev1alpha1.MachineTypeStatus) {
	m.cache[id] = item
}

func (m *MachineTypeClient) ReadCache(id string) *machinetypev1alpha1.MachineTypeStatus {
	return m.cache[id]
}

func (m *MachineTypeClient) ListMachineTypes(
	_ context.Context,
	_ *machinetypev1alpha1.ListMachineTypesRequest,
	_ ...grpc.CallOption,
) (*machinetypev1alpha1.ListMachineTypesResponse, error) {
	return nil, nil
}

func (m *MachineTypeClient) Scan(
	_ context.Context,
	in *machinetypev1alpha1.ScanRequest,
	_ ...grpc.CallOption,
) (*machinetypev1alpha1.ScanResponse, error) {
	if in.Name == "failed-scan" {
		return nil, errors.New("fake error")
	}
	key := types.NamespacedName{Name: in.Name, Namespace: in.Namespace}
	uid := uuidutil.UUIDFromObjectKey(key)
	_, ok := m.cache[uid]
	if ok {
		return &machinetypev1alpha1.ScanResponse{
			Result: commonv1alpha1.RequestResult_REQUEST_RESULT_SUCCESS,
		}, nil
	}
	m.cache[uid] = &machinetypev1alpha1.MachineTypeStatus{}
	return &machinetypev1alpha1.ScanResponse{
		Result: commonv1alpha1.RequestResult_REQUEST_RESULT_SCHEDULED,
	}, nil
}

func (m *MachineTypeClient) UpdateMachineTypeStatus(
	_ context.Context,
	_ *machinetypev1alpha1.UpdateMachineTypeStatusRequest,
	_ ...grpc.CallOption,
) (*machinetypev1alpha1.UpdateMachineTypeStatusResponse, error) {
	return nil, nil
}

func (m *MachineTypeClient) AddMachineGroup(
	_ context.Context,
	_ *machinetypev1alpha1.AddMachineGroupRequest,
	_ ...grpc.CallOption,
) (*machinetypev1alpha1.AddMachineGroupResponse, error) {
	return nil, nil
}

func (m *MachineTypeClient) RemoveMachineGroup(
	_ context.Context,
	_ *machinetypev1alpha1.RemoveMachineGroupRequest,
	_ ...grpc.CallOption,
) (*machinetypev1alpha1.RemoveMachineGroupResponse, error) {
	return nil, nil
}
