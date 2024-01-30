package apiutil

import (
	lifecyclev1alpha1 "github.com/ironcore-dev/lifecycle-manager/api/lifecycle/v1alpha1"
	commonv1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/common/v1alpha1"
	machinev1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/machine/v1alpha1"
	machinetypev1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/machinetype/v1alpha1"
)

func MachineToGrpcAPI(src *lifecyclev1alpha1.Machine) *machinev1alpha1.Machine {
	return nil
}

func MachineStatusToGrpcAPI(src lifecyclev1alpha1.MachineStatus) *machinev1alpha1.MachineStatus {
	return nil
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

func MachineTypeToGrpcAPI(src *lifecyclev1alpha1.MachineType) *machinetypev1alpha1.MachineType {
	return nil
}

func MachineTypeStatusToGrpcAPI(src lifecyclev1alpha1.MachineTypeStatus) *machinetypev1alpha1.MachineTypeStatus {
	return nil
}
