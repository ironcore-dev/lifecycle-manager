// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

// Package v1alpha1 contains API Schema definitions for the lifecycle v1alpha1 API group
// +groupName=lifecycle.ironcore.dev

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

var (
	// GroupVersion is a group & version definition for provided API types.
	GroupVersion = schema.GroupVersion{Group: "lifecycle.ironcore.dev", Version: "v1alpha1"}

	// SchemeBuilder builds a new scheme to map provided API types to kubernetes GroupVersionKind.
	SchemeBuilder = &scheme.Builder{GroupVersion: GroupVersion}

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToScheme = SchemeBuilder.AddToScheme
)

func init() {
	SchemeBuilder.Register(&Machine{}, &MachineList{})
	SchemeBuilder.Register(&MachineType{}, &MachineTypeList{})
	SchemeBuilder.Register(&FirmwarePackage{}, &FirmwarePackageList{})
	SchemeBuilder.Register(&PackageInstallTask{}, &PackageInstallTaskList{})
	SchemeBuilder.Register(&PackageInstallJob{}, &PackageInstallJobList{})
}
