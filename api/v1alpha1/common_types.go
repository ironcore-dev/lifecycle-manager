// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

type ScanResult = string

const (
	ScanFailure ScanResult = "Failure"
	ScanSuccess ScanResult = "Success"
)

type UpdateJobState = string

const (
	UpdateJobStateFailure UpdateJobState = "Failure"
	UpdateJobStateSuccess UpdateJobState = "Success"
)
