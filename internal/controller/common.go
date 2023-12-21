package controller

import (
	"github.com/ironcore-dev/lifecycle-manager/api/lifecycle/v1alpha1"
	commonv1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/common/v1alpha1"
)

var LCIMScanResultToString = map[commonv1alpha1.ScanResult]v1alpha1.ScanResult{
	commonv1alpha1.ScanResult_SCAN_RESULT_UNSPECIFIED: "",
	commonv1alpha1.ScanResult_SCAN_RESULT_SUCCESS:     v1alpha1.ScanSuccess,
	commonv1alpha1.ScanResult_SCAN_RESULT_FAILURE:     v1alpha1.ScanFailure,
}
