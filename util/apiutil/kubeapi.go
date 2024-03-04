// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package apiutil

import (
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
		LastScanTime:      &src.LastScanTime,
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

func ConditionsToGrpcAPI(src []metav1.Condition) []*metav1.Condition {
	result := make([]*metav1.Condition, len(src))
	for i, item := range src {
		el := &metav1.Condition{
			Type:               item.Type,
			Status:             item.Status,
			ObservedGeneration: item.ObservedGeneration,
			LastTransitionTime: item.LastTransitionTime,
			Reason:             item.Reason,
			Message:            item.Message,
		}
		result[i] = el
	}
	return result
}

func MachineTypeToGrpcAPI(src *lifecyclev1alpha1.MachineType) *machinetypev1alpha1.MachineType {
	return nil
}

func MachineTypeStatusToGrpcAPI(src lifecyclev1alpha1.MachineTypeStatus) *machinetypev1alpha1.MachineTypeStatus {
	return nil
}
