// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package apiutil

import (
	"slices"

	lifecyclev1alpha1 "github.com/ironcore-dev/lifecycle-manager/api/lifecycle/v1alpha1"
	commonv1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/common/v1alpha1"
	machinev1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/machine/v1alpha1"
	machinetypev1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/machinetype/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var ScanResultToInt = map[lifecyclev1alpha1.ScanResult]commonv1alpha1.ScanResult{
	lifecyclev1alpha1.Unspecified: commonv1alpha1.ScanResult_SCAN_RESULT_UNSPECIFIED,
	lifecyclev1alpha1.ScanSuccess: commonv1alpha1.ScanResult_SCAN_RESULT_SUCCESS,
	lifecyclev1alpha1.ScanFailure: commonv1alpha1.ScanResult_SCAN_RESULT_FAILURE,
}

func MachineToGrpcAPI(src *lifecyclev1alpha1.Machine) *machinev1alpha1.Machine {
	m := &machinev1alpha1.Machine{
		TypeMeta:   &src.TypeMeta,
		ObjectMeta: &src.ObjectMeta,
		Spec:       MachineSpecToGrpcAPI(src.Spec),
		Status:     MachineStatusToGrpcAPI(src.Status),
	}
	return m
}

func MachineSpecToGrpcAPI(src lifecyclev1alpha1.MachineSpec) *machinev1alpha1.MachineSpec {
	s := &machinev1alpha1.MachineSpec{
		MachineTypeRef: &corev1.LocalObjectReference{Name: src.MachineTypeRef.Name},
		OobMachineRef:  &corev1.LocalObjectReference{Name: src.OOBMachineRef.Name},
		ScanPeriod:     &src.ScanPeriod,
		Packages:       PackageVersionsToGrpcAPI(src.Packages),
	}
	return s
}

func MachineStatusToGrpcAPI(src lifecyclev1alpha1.MachineStatus) *machinev1alpha1.MachineStatus {
	s := &machinev1alpha1.MachineStatus{
		LastScanTime:      &metav1.Timestamp{Seconds: src.LastScanTime.Unix()},
		LastScanResult:    ScanResultToInt[src.LastScanResult],
		InstalledPackages: PackageVersionsToGrpcAPI(src.InstalledPackages),
		Message:           src.Message,
		Conditions:        ConditionsToGrpcAPI(src.Conditions),
	}
	return s
}

func PackageVersionsToGrpcAPI(src []lifecyclev1alpha1.PackageVersion) []*commonv1alpha1.PackageVersion {
	result := make([]*commonv1alpha1.PackageVersion, len(src))
	for i, item := range src {
		el := &commonv1alpha1.PackageVersion{
			Name:    item.Name,
			Version: item.Version,
		}
		result[i] = el
	}
	return result
}

func ConditionsToGrpcAPI(src []metav1.Condition) []*commonv1alpha1.Condition {
	result := make([]*commonv1alpha1.Condition, len(src))
	for i, item := range src {
		el := &commonv1alpha1.Condition{
			Type:               item.Type,
			Status:             string(item.Status),
			ObservedGeneration: item.ObservedGeneration,
			LastTransitionTime: &metav1.Timestamp{Seconds: item.LastTransitionTime.Unix()},
			Reason:             item.Reason,
			Message:            item.Message,
		}
		result[i] = el
	}
	return result
}

func MachineTypeToGrpcAPI(src *lifecyclev1alpha1.MachineType) *machinetypev1alpha1.MachineType {
	m := &machinetypev1alpha1.MachineType{
		TypeMeta:   &src.TypeMeta,
		ObjectMeta: &src.ObjectMeta,
		Spec:       MachineTypeSpecToGrpcAPI(src.Spec),
		Status:     MachineTypeStatusToGrpcAPI(src.Status),
	}
	return m
}

func MachineTypeSpecToGrpcAPI(src lifecyclev1alpha1.MachineTypeSpec) *machinetypev1alpha1.MachineTypeSpec {
	s := &machinetypev1alpha1.MachineTypeSpec{
		Manufacturer:  src.Manufacturer,
		Type:          src.Type,
		ScanPeriod:    &src.ScanPeriod,
		MachineGroups: MachineGroupsToGrpcAPI(src.MachineGroups),
	}
	return s
}

func MachineGroupsToGrpcAPI(src []lifecyclev1alpha1.MachineGroup) []*machinetypev1alpha1.MachineGroup {
	result := make([]*machinetypev1alpha1.MachineGroup, len(src))
	for i, item := range src {
		el := &machinetypev1alpha1.MachineGroup{
			Name:            item.Name,
			MachineSelector: item.MachineSelector.DeepCopy(),
			Packages:        PackageVersionsToGrpcAPI(item.Packages),
		}
		result[i] = el
	}
	return result
}

func MachineTypeStatusToGrpcAPI(src lifecyclev1alpha1.MachineTypeStatus) *machinetypev1alpha1.MachineTypeStatus {
	s := &machinetypev1alpha1.MachineTypeStatus{
		LastScanTime:      &metav1.Timestamp{Seconds: src.LastScanTime.Unix()},
		LastScanResult:    ScanResultToInt[src.LastScanResult],
		AvailablePackages: AvailablePackagesToGrpcAPI(src.AvailablePackages),
		Message:           src.Message,
	}
	return s
}

func AvailablePackagesToGrpcAPI(
	src []lifecyclev1alpha1.AvailablePackageVersions,
) []*machinetypev1alpha1.AvailablePackageVersions {
	result := make([]*machinetypev1alpha1.AvailablePackageVersions, len(src))
	for i, item := range src {
		el := &machinetypev1alpha1.AvailablePackageVersions{
			Name:     item.Name,
			Versions: slices.Clone(item.Versions),
		}
		result[i] = el
	}
	return result
}
