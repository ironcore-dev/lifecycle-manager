package controller

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	lifecyclev1alpha1 "github.com/ironcore-dev/lifecycle-manager/api/lifecycle/v1alpha1"
	"github.com/ironcore-dev/lifecycle-manager/internal/util/uuidutil"
	"github.com/ironcore-dev/lifecycle-manager/lcmi/fake"
)

func TestMachineTypeReconciler_Reconcile(t *testing.T) {
	t.Parallel()

	schemeOpts := []schemeOption{withGroupVersion(lifecyclev1alpha1.AddToScheme)}
	now := time.Unix(time.Now().Unix(), 0)

	tests := map[string]struct {
		machineType  *lifecyclev1alpha1.MachineType
		request      ctrl.Request
		brokerClient *fake.MachineTypeClient
	}{
		"absent-machine-type": {
			machineType: &lifecyclev1alpha1.MachineType{},
			request: ctrl.Request{NamespacedName: types.NamespacedName{
				Name:      "absent-machinetype",
				Namespace: "ironcore",
			}},
		},
		"machinetype-being-deleted": {
			machineType: &lifecyclev1alpha1.MachineType{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "sample-machinetype",
					Namespace:         "ironcore",
					DeletionTimestamp: &metav1.Time{Time: time.Now()},
					Finalizers:        []string{"test-finalizer"},
				},
				Spec: lifecyclev1alpha1.MachineTypeSpec{
					Manufacturer: "Sample",
					Type:         "abc",
					ScanPeriod:   metav1.Duration{Duration: time.Hour * 24},
					MachineGroups: []lifecyclev1alpha1.MachineGroup{
						{
							MachineSelector: metav1.LabelSelector{
								MatchLabels: map[string]string{"env": "test"},
							},
							Packages: []lifecyclev1alpha1.PackageVersion{
								{Name: "raid", Version: "3.2.1"},
							},
						},
					},
				},
			},
			request: ctrl.Request{NamespacedName: types.NamespacedName{
				Name:      "sample-machinetype",
				Namespace: "ironcore",
			}},
		},
		"scan-request-submitted": {
			machineType: &lifecyclev1alpha1.MachineType{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sample-machinetype",
					Namespace: "ironcore",
				},
				Spec: lifecyclev1alpha1.MachineTypeSpec{
					Manufacturer: "Sample",
					Type:         "abc",
					ScanPeriod:   metav1.Duration{Duration: time.Hour * 24},
					MachineGroups: []lifecyclev1alpha1.MachineGroup{
						{
							MachineSelector: metav1.LabelSelector{
								MatchLabels: map[string]string{"env": "test"},
							},
							Packages: []lifecyclev1alpha1.PackageVersion{
								{Name: "raid", Version: "3.2.1"},
							},
						},
					},
				},
			},
			request: ctrl.Request{NamespacedName: types.NamespacedName{
				Name:      "sample-machinetype",
				Namespace: "ironcore",
			}},
			brokerClient: fake.NewMachineTypeClient(map[string]*lifecyclev1alpha1.MachineTypeStatus{}),
		},
		"scan-response-received": {
			machineType: &lifecyclev1alpha1.MachineType{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sample-machinetype",
					Namespace: "ironcore",
				},
				Spec: lifecyclev1alpha1.MachineTypeSpec{
					Manufacturer: "Sample",
					Type:         "abc",
					ScanPeriod:   metav1.Duration{Duration: time.Hour * 24},
					MachineGroups: []lifecyclev1alpha1.MachineGroup{
						{
							MachineSelector: metav1.LabelSelector{
								MatchLabels: map[string]string{"env": "test"},
							},
							Packages: []lifecyclev1alpha1.PackageVersion{
								{Name: "raid", Version: "3.2.1"},
							},
						},
					},
				},
			},
			request: ctrl.Request{NamespacedName: types.NamespacedName{
				Name:      "sample-machinetype",
				Namespace: "ironcore",
			}},
			brokerClient: fake.NewMachineTypeClient(map[string]*lifecyclev1alpha1.MachineTypeStatus{
				uuidutil.UUIDFromObjectKey(types.NamespacedName{
					Name:      "sample-machinetype",
					Namespace: "ironcore",
				}): {
					LastScanTime:   metav1.Time{Time: now},
					LastScanResult: lifecyclev1alpha1.ScanSuccess,
					AvailablePackages: []lifecyclev1alpha1.AvailablePackageVersions{
						{Name: "raid", Versions: []string{"2.9.0", "3.0.0", "3.2.1"}},
					},
				},
			}),
		},
	}

	for n, tc := range tests {
		name, tt := n, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			clientOpts := []clientOption{withRuntimeObject(tt.machineType)}
			r := newMachineTypeReconciler(t, schemeOpts, clientOpts)
			r.Broker = tt.brokerClient
			resp, err := r.Reconcile(context.Background(), tt.request)
			assert.NoError(t, err)
			assert.Empty(t, resp)
			switch name {
			case "scan-request-submitted":
				broker, _ := r.Broker.(*fake.MachineTypeClient)
				entry := broker.ReadCache(uuidutil.UUIDFromObjectKey(client.ObjectKeyFromObject(tt.machineType)))
				reconciledMachineType := &lifecyclev1alpha1.MachineType{}
				_ = r.Get(context.Background(), client.ObjectKeyFromObject(tt.machineType), reconciledMachineType)
				assert.NotNil(t, entry)
				assert.Equal(t, StatusMessageScanRequestSubmitted, reconciledMachineType.Status.Message)
			case "scan-response-received":
				broker, _ := r.Broker.(*fake.MachineTypeClient)
				entry := broker.ReadCache(uuidutil.UUIDFromObjectKey(client.ObjectKeyFromObject(tt.machineType)))
				reconciledMachineType := &lifecyclev1alpha1.MachineType{}
				_ = r.Get(context.Background(), client.ObjectKeyFromObject(tt.machineType), reconciledMachineType)
				assert.Equal(t, *entry, reconciledMachineType.Status)
			}
		})
	}
}
