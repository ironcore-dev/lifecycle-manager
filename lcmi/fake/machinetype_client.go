package fake

import (
	"context"

	"google.golang.org/grpc"
	"k8s.io/apimachinery/pkg/types"

	lifecyclev1alpha1 "github.com/ironcore-dev/lifecycle-manager/api/lifecycle/v1alpha1"
	"github.com/ironcore-dev/lifecycle-manager/internal/util/uuidutil"
	machinetypev1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/machinetype/v1alpha1"
)

type MachineTypeClient struct {
	cache map[string]*lifecyclev1alpha1.MachineTypeStatus
}

func NewMachineTypeClient(cache map[string]*lifecyclev1alpha1.MachineTypeStatus) *MachineTypeClient {
	return &MachineTypeClient{cache: cache}
}

func (m *MachineTypeClient) WriteCache(id string, item *lifecyclev1alpha1.MachineTypeStatus) {
	m.cache[id] = item
}

func (m *MachineTypeClient) ReadCache(id string) *lifecyclev1alpha1.MachineTypeStatus {
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
	key := types.NamespacedName{Name: in.Name, Namespace: in.Namespace}
	uid := uuidutil.UUIDFromObjectKey(key)
	entry, ok := m.cache[uid]
	if ok {
		return &machinetypev1alpha1.ScanResponse{Response: entry}, nil
	}
	m.cache[uid] = &lifecyclev1alpha1.MachineTypeStatus{}
	return &machinetypev1alpha1.ScanResponse{Response: nil}, nil
}
