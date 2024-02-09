// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package apiutil

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	lifecyclev1alpha1 "github.com/ironcore-dev/lifecycle-manager/api/lifecycle/v1alpha1"
	lifecycleapplyv1alpha1 "github.com/ironcore-dev/lifecycle-manager/clientgo/applyconfiguration/lifecycle/v1alpha1"
	v1 "github.com/ironcore-dev/lifecycle-manager/clientgo/applyconfiguration/meta/v1"
	commonv1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/common/v1alpha1"
	machinev1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/machine/v1alpha1"
	machinetypev1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/machinetype/v1alpha1"
)

func MachineToKubeAPI(src *machinev1alpha1.Machine) *lifecyclev1alpha1.Machine {
	return nil
}

func MachineStatusToKubeAPI(src *machinev1alpha1.MachineStatus) lifecyclev1alpha1.MachineStatus {
	if src == nil {
		return lifecyclev1alpha1.MachineStatus{}
	}
	status := lifecyclev1alpha1.MachineStatus{
		LastScanTime:      *src.LastScanTime,
		LastScanResult:    lifecyclev1alpha1.ScanResult(src.LastScanResult),
		InstalledPackages: PackageVersionsToKubeAPI(src.InstalledPackages),
		Message:           src.Message,
	}
	return status
}

func MachineStatusToApplyConfiguration(
	src *machinev1alpha1.MachineStatus,
) *lifecycleapplyv1alpha1.MachineStatusApplyConfiguration {
	return lifecycleapplyv1alpha1.MachineStatus().
		WithMessage(src.Message).
		WithConditions(ConditionsToApplyConfiguration(src.Conditions)...).
		WithInstalledPackages(PackageVersionsToApplyConfiguration(src.InstalledPackages)...).
		WithLastScanResult(lifecyclev1alpha1.ScanResult(src.LastScanResult)).
		WithLastScanTime(*src.LastScanTime)
}

func ConditionsToApplyConfiguration(
	src []*metav1.Condition,
) []*v1.ConditionApplyConfiguration {
	result := make([]*v1.ConditionApplyConfiguration, len(src))
	for i, item := range src {
		result[i] = v1.Condition().
			WithMessage(item.Message).
			WithReason(item.Reason).
			WithType(item.Type).
			WithObservedGeneration(item.ObservedGeneration).
			WithLastTransitionTime(item.LastTransitionTime).
			WithStatus(item.Status)
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
		LastScanTime:      *src.LastScanTime,
		LastScanResult:    lifecyclev1alpha1.ScanResult(src.LastScanResult),
		AvailablePackages: AvailablePackageVersionsToKubeAPI(src.AvailablePackages),
		Message:           src.Message,
	}
	return status
}

func MachineTypeStatusToApplyConfiguration(
	src *machinetypev1alpha1.MachineTypeStatus,
) *lifecycleapplyv1alpha1.MachineTypeStatusApplyConfiguration {
	return lifecycleapplyv1alpha1.MachineTypeStatus().
		WithMessage(src.Message).
		WithAvailablePackages(AvailablePackagesToApplyConfiguration(src.AvailablePackages)...).
		WithLastScanResult(lifecyclev1alpha1.ScanResult(src.LastScanResult)).
		WithLastScanTime(*src.LastScanTime)
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
