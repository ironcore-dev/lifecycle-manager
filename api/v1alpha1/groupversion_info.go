// Copyright 2023 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	SchemeBuilder.Register(&PackageUpdate{}, &PackageUpdateList{})
	runtimeutil.Must(SchemeBuilder.AddToScheme(scheme))
}
