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

var LCMIConditionMapping = map[string]int32{
	v1alpha1.ConditionTypePending:    1,
	v1alpha1.ConditionTypeScheduled:  2,
	v1alpha1.ConditionTypeInProgress: 3,
	v1alpha1.ConditionTypeFinished:   4,
}

var LCMIConditionTypeToString = map[lcmicommon.ConditionType]string{
	lcmicommon.ConditionType_CONDITION_TYPE_UNSPECIFIED: "",
	lcmicommon.ConditionType_CONDITION_TYPE_PENDING:     v1alpha1.ConditionTypePending,
	lcmicommon.ConditionType_CONDITION_TYPE_SCHEDULED:   v1alpha1.ConditionTypeScheduled,
	lcmicommon.ConditionType_CONDITION_TYPE_IN_PROGRESS: v1alpha1.ConditionTypeInProgress,
	lcmicommon.ConditionType_CONDITION_TYPE_FINISHED:    v1alpha1.ConditionTypeFinished,
}

var LCMIJobStateToString = map[lcmicommon.PackageInstallJobState]v1alpha1.PackageInstallJobState{
	lcmicommon.PackageInstallJobState_PACKAGE_INSTALL_JOB_STATE_UNSPECIFIED: "",
	lcmicommon.PackageInstallJobState_PACKAGE_INSTALL_JOB_STATE_SUCCESS:     v1alpha1.PackageInstallJobStateSuccess,
	lcmicommon.PackageInstallJobState_PACKAGE_INSTALL_JOB_STATE_FAILURE:     v1alpha1.PackageInstallJobStateFailure,
}
