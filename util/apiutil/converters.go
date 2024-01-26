// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package apiutil

import (
	lifecyclev1alpha1 "github.com/ironcore-dev/lifecycle-manager/api/lifecycle/v1alpha1"
	commonv1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/common/v1alpha1"
	machinev1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/machine/v1alpha1"
	machinetypev1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/machinetype/v1alpha1"
)

// todo: conversion k8s api -> grpc api
// todo: conversion grpc api -> k8s api

func MachineTo(src *lifecyclev1alpha1.Machine) *machinev1alpha1.Machine {
	return nil
}

func MachineFrom(src *machinev1alpha1.Machine) *lifecyclev1alpha1.Machine {
	return nil
}

func MachineStatusTo(src lifecyclev1alpha1.MachineStatus) *machinev1alpha1.MachineStatus {
	return nil
}

func MachineStatusFrom(src *machinev1alpha1.MachineStatus) lifecyclev1alpha1.MachineStatus {
	if src == nil {
		return lifecyclev1alpha1.MachineStatus{}
	}
	status := lifecyclev1alpha1.MachineStatus{
		LastScanTime:      *src.LastScanTime,
		LastScanResult:    lifecyclev1alpha1.ScanResult(src.LastScanResult),
		InstalledPackages: PackageVersionsFrom(src.InstalledPackages),
		Message:           src.Message,
	}
	return status
}

func PackageVersionsTo(src []lifecyclev1alpha1.PackageVersion) []*commonv1alpha1.PackageVersion {
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

func PackageVersionsFrom(src []*commonv1alpha1.PackageVersion) []lifecyclev1alpha1.PackageVersion {
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

func MachineTypeTo(src *lifecyclev1alpha1.MachineType) *machinetypev1alpha1.MachineType {
	return nil
}

func MachineTypeFrom(src *machinetypev1alpha1.MachineType) *lifecyclev1alpha1.MachineType {
	return nil
}

func MachineTypeStatusTo(src lifecyclev1alpha1.MachineTypeStatus) *machinetypev1alpha1.MachineTypeStatus {
	return nil
}

func MachineTypeStatusFrom(src *machinetypev1alpha1.MachineTypeStatus) lifecyclev1alpha1.MachineTypeStatus {
	if src == nil {
		return lifecyclev1alpha1.MachineTypeStatus{}
	}
	status := lifecyclev1alpha1.MachineTypeStatus{
		LastScanTime:      *src.LastScanTime,
		LastScanResult:    lifecyclev1alpha1.ScanResult(src.LastScanResult),
		AvailablePackages: AvailablePackageVersionsFrom(src.AvailablePackages),
		Message:           src.Message,
	}
	return status
}

func AvailablePackageVersionsFrom(
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
