// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/ironcore-dev/lifecycle-manager/api/lifecycle/v1alpha1"
	commonv1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/common/v1alpha1"
)

const (
	InstallScheduled = "Scheduled"
	InstallFailed    = "Failed"
)

var patchOpts = &client.SubResourcePatchOptions{
	PatchOptions: client.PatchOptions{
		Force:        pointer.Bool(true),
		FieldManager: "lifecycle-manager",
	},
}

var LCIMScanResultToString = map[commonv1alpha1.ScanResult]v1alpha1.ScanResult{
	commonv1alpha1.ScanResult_SCAN_RESULT_UNSPECIFIED: "",
	commonv1alpha1.ScanResult_SCAN_RESULT_SUCCESS:     v1alpha1.ScanSuccess,
	commonv1alpha1.ScanResult_SCAN_RESULT_FAILURE:     v1alpha1.ScanFailure,
}

var LCIMInstallResultToString = map[commonv1alpha1.InstallResult]string{
	commonv1alpha1.InstallResult_INSTALL_RESULT_UNSPECIFIED: "",
	commonv1alpha1.InstallResult_INSTALL_RESULT_SCHEDULED:   InstallScheduled,
	commonv1alpha1.InstallResult_INSTALL_RESULT_FAILURE:     InstallFailed,
}
