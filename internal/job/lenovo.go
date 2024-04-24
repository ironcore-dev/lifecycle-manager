// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package job

import (
	"context"

	commonv1alpha1 "github.com/ironcore-dev/lifecycle-manager/api/proto/common/v1alpha1"
	machinev1alpha1 "github.com/ironcore-dev/lifecycle-manager/api/proto/machine/v1alpha1"
	"github.com/ironcore-dev/lifecycle-manager/internal/util/convertutil"
	oobv1alpha1 "github.com/ironcore-dev/oob/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

const onecli = "/lenovo/onecli"

func (w *MachineLifecycleWorker) LenovoScanMachine(ctx context.Context, machine *machinev1alpha1.Machine) error {
	oob := &oobv1alpha1.OOB{}
	key := types.NamespacedName{
		Namespace: machine.ObjectMeta.Namespace,
		Name:      machine.Spec.OobMachineRef.Name,
	}
	if err := w.Get(ctx, key, oob); err != nil {
		return err
	}

	// todo:
	//  0. remove dumb status update below
	//  1. implement execution of onecli app with args to collect info about installed software (output format XML)
	//  2. parse XML and update machine's status
	machine.Status.InstalledPackages = []*commonv1alpha1.PackageVersion{
		{
			Name:    "bios",
			Version: "2.8.1",
		},
		{
			Name:    "bmc",
			Version: "1.1.0",
		},
		{
			Name:    "raid",
			Version: "7.0.0+rc1",
		},
	}
	machine.Status.LastScanTime = convertutil.TimeToTimestampPtr(metav1.Now())
	machine.Status.LastScanResult = commonv1alpha1.ScanResult_SCAN_RESULT_SUCCESS
	return nil
}
