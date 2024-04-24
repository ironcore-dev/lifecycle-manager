// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package job

import (
	"context"
	"log/slog"

	"github.com/ironcore-dev/lifecycle-manager/clientgo/connectrpc/machinetype/v1alpha1/machinetypev1alpha1connect"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type MachineTypeLifecycleWorker struct {
	machinetypev1alpha1connect.MachineTypeServiceClient
	client.Client
	log   *slog.Logger
	jobId string
}

func NewMachineTypeLifecycleWorker(opts Options) *MachineTypeLifecycleWorker {
	return &MachineTypeLifecycleWorker{
		log:    opts.Log,
		jobId:  opts.JobId,
		Client: opts.KubeClient,
	}
}

func (w *MachineTypeLifecycleWorker) WithClient(
	c machinetypev1alpha1connect.MachineTypeServiceClient,
) *MachineTypeLifecycleWorker {
	w.MachineTypeServiceClient = c
	return w
}

func (w *MachineTypeLifecycleWorker) Start(ctx context.Context) error {
	return nil
}
