package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PackageVersion defines the concrete package version item.
type PackageVersion struct {
	// Name defines the name of the firmware package.
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// Version defines the version of the firmware package.
	// +kubebuilder:validation:Required
	Version string `json:"version"`
}

// InstalledPackage defines the concrete installed package version
// with installation timestamp.
type InstalledPackage struct {
	PackageVersion `json:",inline"`

	// InstallTime reflects timestamp when the package was installed.
	// +kubebuilder:validation:Required
	InstallTime metav1.Time `json:"installTime"`
}

type ScanResult string

const (
	ScanFailure ScanResult = "Failure"
	ScanSuccess ScanResult = "Success"
)

func (in ScanResult) IsSuccess() bool {
	return in == ScanSuccess
}

func (in ScanResult) IsFailure() bool {
	return in == ScanFailure
}
