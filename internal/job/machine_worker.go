// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package job

import (
	"context"
	"fmt"
	"log/slog"

	"connectrpc.com/connect"
	lifecyclev1alpha1 "github.com/ironcore-dev/lifecycle-manager/api/lifecycle/v1alpha1"
	commonv1alpha1 "github.com/ironcore-dev/lifecycle-manager/api/proto/common/v1alpha1"
	machinev1alpha1 "github.com/ironcore-dev/lifecycle-manager/api/proto/machine/v1alpha1"
	"github.com/ironcore-dev/lifecycle-manager/clientgo/connectrpc/machine/v1alpha1/machinev1alpha1connect"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type MachineLifecycleWorker struct {
	machinev1alpha1connect.MachineServiceClient
	client.Client
	log   *slog.Logger
	jobId string
}

func NewMachineLifecycleWorker(opts Options) *MachineLifecycleWorker {
	return &MachineLifecycleWorker{
		log:    opts.Log,
		jobId:  opts.JobId,
		Client: opts.KubeClient,
	}
}

func (w *MachineLifecycleWorker) WithClient(c machinev1alpha1connect.MachineServiceClient) *MachineLifecycleWorker {
	w.MachineServiceClient = c
	return w
}

func (w *MachineLifecycleWorker) Start(ctx context.Context) error {
	var (
		getJobResponse              *connect.Response[machinev1alpha1.GetJobResponse]
		updateMachineStatusResponse *connect.Response[machinev1alpha1.UpdateMachineStatusResponse]
		err                         error
	)
	getJobResponse, err = w.GetJob(ctx, connect.NewRequest(&machinev1alpha1.GetJobRequest{Id: w.jobId}))
	if err != nil {
		w.log.Error("error getting job", "error", err)
		return err
	}
	task := getJobResponse.Msg
	target := task.Target
	switch task.JobType {
	case "scan":
		err = w.scan(ctx, target)
	case "install":
		err = w.install(ctx, target)
	}
	if err != nil {
		return err
	}
	updateMachineStatusResponse, err = w.UpdateMachineStatus(ctx, connect.NewRequest(
		&machinev1alpha1.UpdateMachineStatusRequest{
			Name:      target.ObjectMeta.Name,
			Namespace: target.ObjectMeta.Namespace,
			Status:    target.Status,
		}))
	if err != nil {
		return err
	}
	if updateMachineStatusResponse.Msg.Result != commonv1alpha1.RequestResult_REQUEST_RESULT_SUCCESS {
		return fmt.Errorf("update machine status failed: %s", updateMachineStatusResponse.Msg.Result)
	}
	return nil
}

func (w *MachineLifecycleWorker) scan(ctx context.Context, target *machinev1alpha1.Machine) error {
	machineType := &lifecyclev1alpha1.MachineType{}
	key := types.NamespacedName{
		Namespace: target.ObjectMeta.Namespace,
		Name:      target.Spec.MachineTypeRef.Name,
	}
	if err := w.Get(ctx, key, machineType); err != nil {
		return err
	}

	switch machineType.Spec.Manufacturer {
	case "Lenovo":
		return w.LenovoScanMachine(ctx, target)
	}

	return nil
}

func (w *MachineLifecycleWorker) install(ctx context.Context, target *machinev1alpha1.Machine) error {
	return nil
}
