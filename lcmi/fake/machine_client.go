package fake

import (
	"context"

	"google.golang.org/grpc"
	"k8s.io/apimachinery/pkg/types"

	lifecyclev1alpha1 "github.com/ironcore-dev/lifecycle-manager/api/lifecycle/v1alpha1"
	"github.com/ironcore-dev/lifecycle-manager/internal/util/uuidutil"
	commonv1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/common/v1alpha1"
	machinev1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/machine/v1alpha1"
)

type MachineClient struct {
	cache map[string]*lifecyclev1alpha1.MachineStatus
}

func NewMachineClient(cache map[string]*lifecyclev1alpha1.MachineStatus) *MachineClient {
	return &MachineClient{cache: cache}
}

func (m *MachineClient) WriteCache(id string, item *lifecyclev1alpha1.MachineStatus) {
	m.cache[id] = item
}

func (m *MachineClient) ReadCache(id string) *lifecyclev1alpha1.MachineStatus {
	return m.cache[id]
}

func (m *MachineClient) ListMachines(
	_ context.Context,
	_ *machinev1alpha1.ListMachinesRequest,
	_ ...grpc.CallOption,
) (*machinev1alpha1.ListMachinesResponse, error) {
	return nil, nil
}

func (m *MachineClient) Scan(
	_ context.Context,
	in *machinev1alpha1.ScanRequest,
	_ ...grpc.CallOption,
) (*machinev1alpha1.ScanResponse, error) {
	key := types.NamespacedName{Name: in.Name, Namespace: in.Namespace}
	uid := uuidutil.UUIDFromObjectKey(key)
	entry, ok := m.cache[uid]
	if ok {
		return &machinev1alpha1.ScanResponse{Response: entry}, nil
	}
	m.cache[uid] = &lifecyclev1alpha1.MachineStatus{}
	return &machinev1alpha1.ScanResponse{Response: nil}, nil
}

func (m *MachineClient) Install(
	_ context.Context,
	in *machinev1alpha1.InstallRequest,
	_ ...grpc.CallOption,
) (*machinev1alpha1.InstallResponse, error) {
	key := types.NamespacedName{Name: in.Name, Namespace: in.Namespace}
	uid := uuidutil.UUIDFromObjectKey(key)
	_, ok := m.cache[uid]
	if ok {
		return &machinev1alpha1.InstallResponse{Result: commonv1alpha1.InstallResult_INSTALL_RESULT_SCHEDULED}, nil
	}
	m.cache[uid] = &lifecyclev1alpha1.MachineStatus{}
	return &machinev1alpha1.InstallResponse{Result: commonv1alpha1.InstallResult_INSTALL_RESULT_FAILURE}, nil
}
