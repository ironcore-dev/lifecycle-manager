// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"github.com/ironcore-dev/lifecycle-manager/api/v1alpha1"
	lcmicommon "github.com/ironcore-dev/lifecycle-manager/lcmi/api/common/v1alpha1"
)

var LCMIScanResultToString = map[lcmicommon.ScanResult]v1alpha1.ScanResult{
	lcmicommon.ScanResult_SCAN_RESULT_UNSPECIFIED: "",
	lcmicommon.ScanResult_SCAN_RESULT_SUCCESS:     v1alpha1.ScanSuccess,
	lcmicommon.ScanResult_SCAN_RESULT_FAILURE:     v1alpha1.ScanFailure,
}

var LCMIScanStateToString = map[lcmicommon.ScanState]v1alpha1.ScanState{
	lcmicommon.ScanState_SCAN_STATE_UNSPECIFIED: "",
	lcmicommon.ScanState_SCAN_STATE_SCHEDULED:   v1alpha1.ScanScheduled,
	lcmicommon.ScanState_SCAN_STATE_FINISHED:    v1alpha1.ScanFinished,
}
