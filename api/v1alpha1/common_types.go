// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

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

type PackageInstallJobState string

const (
	PackageInstallJobStateFailure PackageInstallJobState = "Failure"
	PackageInstallJobStateSuccess PackageInstallJobState = "Success"
)

func (in PackageInstallJobState) IsSuccess() bool {
	return in == PackageInstallJobStateSuccess
}

func (in PackageInstallJobState) IsFailure() bool {
	return in == PackageInstallJobStateFailure
}

type ScanState string

const (
	ScanScheduled ScanState = "Scheduled"
	ScanFinished  ScanState = "Finished"
)

func (in ScanState) IsScheduled() bool {
	return in == ScanScheduled
}

func (in ScanState) IsFinished() bool {
	return in == ScanFinished
}
