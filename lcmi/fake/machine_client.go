package fake

import (
	"context"

	"google.golang.org/grpc"
	"k8s.io/apimachinery/pkg/types"

	"github.com/ironcore-dev/lifecycle-manager/internal/util/uuidutil"
	"github.com/ironcore-dev/lifecycle-manager/lcmi/api/machine/v1alpha1"
)

type MachineClient struct {
	cache map[string]*v1alpha1.ScanResponse
}

func NewMachineClient(cache map[string]*v1alpha1.ScanResponse) *MachineClient {
	return &MachineClient{cache: cache}
}

func (m *MachineClient) WriteCache(id string, item *v1alpha1.ScanResponse) {
	m.cache[id] = item
}

func (m *MachineClient) ReadCache(id string) *v1alpha1.ScanResponse {
	return m.cache[id]
}

func (m *MachineClient) ListMachines(
	_ context.Context,
	_ *v1alpha1.ListMachinesRequest,
	_ ...grpc.CallOption,
) (*v1alpha1.ListMachinesResponse, error) {
	return nil, nil
}

func (m *MachineClient) Scan(
	_ context.Context,
	in *v1alpha1.ScanRequest,
	_ ...grpc.CallOption,
) (*v1alpha1.ScanResponse, error) {
	key := types.NamespacedName{Name: in.Name, Namespace: in.Namespace}
	uid := uuidutil.UUIDFromObjectKey(key)
	return m.cache[uid], nil
}

func (m *MachineClient) Install(
	ctx context.Context,
	in *v1alpha1.InstallRequest,
	opts ...grpc.CallOption,
) (*v1alpha1.InstallResponse, error) {
	return nil, nil
}
