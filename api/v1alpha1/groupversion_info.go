// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	runtimeutil "k8s.io/apimachinery/pkg/util/runtime"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

var (
	// GroupVersion is a group & version definition for provided API types.
	GroupVersion = schema.GroupVersion{Group: "metal.ironcore.dev", Version: "v1alpha1"}

	// SchemeBuilder builds a new scheme to map provided API types to kubernetes GroupVersionKind.
	SchemeBuilder = &scheme.Builder{GroupVersion: GroupVersion}
)

// AddToScheme adds provided API types to scheme.Builder and adds them to runtime.Scheme.
func AddToScheme(scheme *runtime.Scheme) {
	SchemeBuilder.Register(&MachineLifecycle{}, &MachineLifecycleList{})
	SchemeBuilder.Register(&MachineType{}, &MachineTypeList{})
	SchemeBuilder.Register(&FirmwarePackage{}, &FirmwarePackageList{})
	SchemeBuilder.Register(&UpdateTask{}, &UpdateTaskList{})
	SchemeBuilder.Register(&MachineUpdateJob{}, &MachineUpdateJobList{})
	runtimeutil.Must(SchemeBuilder.AddToScheme(scheme))
}
