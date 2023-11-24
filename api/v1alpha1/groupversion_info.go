package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	runtimeutil "k8s.io/apimachinery/pkg/util/runtime"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

var (
	// GroupVersion is a group & version definition for provided API types
	GroupVersion = schema.GroupVersion{Group: "metal.ironcore.dev", Version: "v1alpha1"}

	// SchemeBuilder builds a new scheme to map provided API types to kubernetes GroupVersionKind
	SchemeBuilder = &scheme.Builder{GroupVersion: GroupVersion}
)

// AddToScheme adds provided API types to scheme.Builder and adds them to runtime.Scheme
func AddToScheme(scheme *runtime.Scheme) {
	SchemeBuilder.Register(&MachineLifecycle{}, &MachineLifecycleList{})
	SchemeBuilder.Register(&MachineType{}, &MachineTypeList{})
	SchemeBuilder.Register(&PackageUpdate{}, &PackageUpdateList{})
	runtimeutil.Must(SchemeBuilder.AddToScheme(scheme))
}
