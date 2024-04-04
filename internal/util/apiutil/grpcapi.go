// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package apiutil

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	lifecyclev1alpha1 "github.com/ironcore-dev/lifecycle-manager/api/lifecycle/v1alpha1"
	commonv1alpha1 "github.com/ironcore-dev/lifecycle-manager/api/proto/common/v1alpha1"
	machinev1alpha1 "github.com/ironcore-dev/lifecycle-manager/api/proto/machine/v1alpha1"
	machinetypev1alpha1 "github.com/ironcore-dev/lifecycle-manager/api/proto/machinetype/v1alpha1"
	lifecycleapplyv1alpha1 "github.com/ironcore-dev/lifecycle-manager/clientgo/applyconfiguration/lifecycle/v1alpha1"
	v1 "github.com/ironcore-dev/lifecycle-manager/clientgo/applyconfiguration/meta/v1"
)

var LCIMScanResultToString = map[commonv1alpha1.ScanResult]lifecyclev1alpha1.ScanResult{
	commonv1alpha1.ScanResult_SCAN_RESULT_UNSPECIFIED: lifecyclev1alpha1.Unspecified,
	commonv1alpha1.ScanResult_SCAN_RESULT_SUCCESS:     lifecyclev1alpha1.ScanSuccess,
	commonv1alpha1.ScanResult_SCAN_RESULT_FAILURE:     lifecyclev1alpha1.ScanFailure,
}

func MachineToKubeAPI(src *machinev1alpha1.Machine) *lifecyclev1alpha1.Machine {
	return nil
}

func MachineStatusToKubeAPI(src *machinev1alpha1.MachineStatus) lifecyclev1alpha1.MachineStatus {
	if src == nil {
		return lifecyclev1alpha1.MachineStatus{}
	}
	status := lifecyclev1alpha1.MachineStatus{
		LastScanTime:      metav1.Time{Time: time.Unix(src.LastScanTime.Seconds, int64(src.LastScanTime.Nanos))},
		LastScanResult:    lifecyclev1alpha1.ScanResult(src.LastScanResult),
		InstalledPackages: PackageVersionsToKubeAPI(src.InstalledPackages),
		Message:           src.Message,
	}
	return status
}

func MachineStatusToApplyConfiguration(
	src *machinev1alpha1.MachineStatus,
) *lifecycleapplyv1alpha1.MachineStatusApplyConfiguration {
	apply := lifecycleapplyv1alpha1.MachineStatus().
		WithMessage(src.Message).
		WithConditions(ConditionsToApplyConfiguration(src.Conditions)...).
		WithInstalledPackages(PackageVersionsToApplyConfiguration(src.InstalledPackages)...).
		WithLastScanResult(LCIMScanResultToString[src.LastScanResult])
	if src.LastScanTime != nil {
		apply = apply.WithLastScanTime(metav1.Time{
			Time: time.Unix(src.LastScanTime.Seconds, int64(src.LastScanTime.Nanos))})
	}
	return apply
}

func ConditionsToApplyConfiguration(
	src []*commonv1alpha1.Condition,
) []*v1.ConditionApplyConfiguration {
	result := make([]*v1.ConditionApplyConfiguration, len(src))
	for i, item := range src {
		result[i] = v1.Condition().
			WithMessage(item.Message).
			WithReason(item.Reason).
			WithType(item.Type).
			WithObservedGeneration(item.ObservedGeneration).
			WithLastTransitionTime(metav1.Time{
				Time: time.Unix(item.LastTransitionTime.Seconds, int64(item.LastTransitionTime.Nanos))}).
			WithStatus(metav1.ConditionStatus(item.Status))
	}
	return result
}

func PackageVersionsToApplyConfiguration(
	src []*commonv1alpha1.PackageVersion,
) []*lifecycleapplyv1alpha1.PackageVersionApplyConfiguration {
	result := make([]*lifecycleapplyv1alpha1.PackageVersionApplyConfiguration, len(src))
	for i, item := range src {
		result[i] = lifecycleapplyv1alpha1.PackageVersion().
			WithName(item.Name).
			WithVersion(item.Version)
	}
	return result
}

func PackageVersionsToKubeAPI(src []*commonv1alpha1.PackageVersion) []lifecyclev1alpha1.PackageVersion {
	if src == nil {
		return []lifecyclev1alpha1.PackageVersion{}
	}
	result := make([]lifecyclev1alpha1.PackageVersion, len(src))
	for i, item := range src {
		el := lifecyclev1alpha1.PackageVersion{
			Name:    item.Name,
			Version: item.Version,
		}
		result[i] = el
	}
	return result
}

func MachineTypeToKubeAPI(src *machinetypev1alpha1.MachineType) *lifecyclev1alpha1.MachineType {
	return nil
}

func MachineTypeStatusToKubeAPI(src *machinetypev1alpha1.MachineTypeStatus) lifecyclev1alpha1.MachineTypeStatus {
	if src == nil {
		return lifecyclev1alpha1.MachineTypeStatus{}
	}
	status := lifecyclev1alpha1.MachineTypeStatus{
		LastScanTime:      metav1.Time{Time: time.Unix(src.LastScanTime.Seconds, int64(src.LastScanTime.Nanos))},
		LastScanResult:    lifecyclev1alpha1.ScanResult(src.LastScanResult),
		AvailablePackages: AvailablePackageVersionsToKubeAPI(src.AvailablePackages),
		Message:           src.Message,
	}
	return status
}

func MachineTypeStatusToApplyConfiguration(
	src *machinetypev1alpha1.MachineTypeStatus,
) *lifecycleapplyv1alpha1.MachineTypeStatusApplyConfiguration {
	apply := lifecycleapplyv1alpha1.MachineTypeStatus().
		WithMessage(src.Message).
		WithAvailablePackages(AvailablePackagesToApplyConfiguration(src.AvailablePackages)...).
		WithLastScanResult(lifecyclev1alpha1.ScanResult(src.LastScanResult))
	if src.LastScanTime != nil {
		apply = apply.WithLastScanTime(metav1.Time{
			Time: time.Unix(src.LastScanTime.Seconds, int64(src.LastScanTime.Nanos))})
	}
	return apply
}

func AvailablePackagesToApplyConfiguration(
	src []*machinetypev1alpha1.AvailablePackageVersions,
) []*lifecycleapplyv1alpha1.AvailablePackageVersionsApplyConfiguration {
	result := make([]*lifecycleapplyv1alpha1.AvailablePackageVersionsApplyConfiguration, len(src))
	for i, item := range src {
		result[i] = lifecycleapplyv1alpha1.AvailablePackageVersions().
			WithName(item.Name).
			WithVersions(item.Versions...)
	}
	return result
}

func AvailablePackageVersionsToKubeAPI(
	src []*machinetypev1alpha1.AvailablePackageVersions,
) []lifecyclev1alpha1.AvailablePackageVersions {
	if src == nil {
		return []lifecyclev1alpha1.AvailablePackageVersions{}
	}
	result := make([]lifecyclev1alpha1.AvailablePackageVersions, len(src))
	for i, item := range src {
		el := lifecyclev1alpha1.AvailablePackageVersions{
			Name:     item.Name,
			Versions: item.Versions,
		}
		result[i] = el
	}
	return result
}

func MachineGroupsToApplyConfiguration(
	src []*machinetypev1alpha1.MachineGroup,
) []*lifecycleapplyv1alpha1.MachineGroupApplyConfiguration {
	result := make([]*lifecycleapplyv1alpha1.MachineGroupApplyConfiguration, len(src))
	for i, item := range src {
		result[i] = lifecycleapplyv1alpha1.MachineGroup().
			WithName(item.Name).
			WithPackages(PackageVersionsToApplyConfiguration(item.Packages)...).
			WithMachineSelector(LabelSelectorToApplyConfiguration(item.MachineSelector))
	}
	return result
}

func LabelSelectorToApplyConfiguration(
	src *metav1.LabelSelector,
) *v1.LabelSelectorApplyConfiguration {
	return v1.LabelSelector().
		WithMatchLabels(src.MatchLabels).
		WithMatchExpressions(MatchExpressionToApplyConfiguration(src.MatchExpressions)...)
}

func MatchExpressionToApplyConfiguration(
	src []metav1.LabelSelectorRequirement,
) []*v1.LabelSelectorRequirementApplyConfiguration {
	result := make([]*v1.LabelSelectorRequirementApplyConfiguration, len(src))
	for i, item := range src {
		result[i] = v1.LabelSelectorRequirement().
			WithKey(item.Key).
			WithOperator(item.Operator).
			WithValues(item.Values...)
	}
	return result
}
