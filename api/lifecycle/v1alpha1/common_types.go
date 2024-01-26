// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	_ "github.com/gogo/protobuf/gogoproto"
	_ "golang.org/x/crypto/cryptobyte"
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
