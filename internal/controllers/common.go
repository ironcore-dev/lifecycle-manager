// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"time"

	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/ironcore-dev/lifecycle-manager/api/lifecycle/v1alpha1"
	commonv1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/common/v1alpha1"
)

type RequestResult string

const requeuePeriod = time.Second * 5

const (
	InstallScheduled = "Scheduled"
	InstallFailed    = "Failed"

	RequestResultScheduled RequestResult = "Scheduled"
	RequestResultSuccess   RequestResult = "Success"
	RequestResultFailure   RequestResult = "Failure"
)

const (
	StatusMessageScanRequestSubmitted  = "scan request submitted"
	StatusMessageScanRequestFailed     = "scan request failed"
	StatusMessageInstallRequestFailed  = "install request failed"
	StatusMessageInstallationScheduled = "packages installation scheduled"
)

func (r RequestResult) IsScheduled() bool {
	return r == RequestResultScheduled
}

func (r RequestResult) IsSuccess() bool {
	return r == RequestResultSuccess
}

func (r RequestResult) IsFailure() bool {
	return r == RequestResultFailure
}

var patchOpts = &client.SubResourcePatchOptions{
	PatchOptions: client.PatchOptions{
		Force:        ptr.To(true),
		FieldManager: "lifecycle-manager",
	},
}

var LCIMScanResultToString = map[commonv1alpha1.ScanResult]v1alpha1.ScanResult{
	commonv1alpha1.ScanResult_SCAN_RESULT_UNSPECIFIED: "",
	commonv1alpha1.ScanResult_SCAN_RESULT_SUCCESS:     v1alpha1.ScanSuccess,
	commonv1alpha1.ScanResult_SCAN_RESULT_FAILURE:     v1alpha1.ScanFailure,
}

var LCIMRequestResultToString = map[commonv1alpha1.RequestResult]RequestResult{
	commonv1alpha1.RequestResult_REQUEST_RESULT_UNSPECIFIED: "",
	commonv1alpha1.RequestResult_REQUEST_RESULT_SCHEDULED:   RequestResultScheduled,
	commonv1alpha1.RequestResult_REQUEST_RESULT_SUCCESS:     RequestResultSuccess,
	commonv1alpha1.RequestResult_REQUEST_RESULT_FAILURE:     RequestResultFailure,
}
