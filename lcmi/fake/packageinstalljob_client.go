// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package fake

import (
	"context"
	"sync"
	"time"

	"golang.org/x/exp/maps"
	"google.golang.org/grpc"

	lcmicommon "github.com/ironcore-dev/lifecycle-manager/lcmi/api/common/v1alpha1"
	lcmimeta "github.com/ironcore-dev/lifecycle-manager/lcmi/api/meta/v1alpha1"
	lcmipackageinstalljob "github.com/ironcore-dev/lifecycle-manager/lcmi/api/packageinstalljob/v1alpha1"
)

type PackageInstallJobClient struct {
	sync.Mutex

	jobs map[string]*lcmipackageinstalljob.PackageInstallJob
}

func NewPackageInstallJobClient() *PackageInstallJobClient {
	return &PackageInstallJobClient{jobs: make(map[string]*lcmipackageinstalljob.PackageInstallJob)}
}

func NewPackageInstallJobClientWithJobs(
	jobs map[string]*lcmipackageinstalljob.PackageInstallJob) *PackageInstallJobClient {
	return &PackageInstallJobClient{jobs: jobs}
}

func (f *PackageInstallJobClient) ListPackageInstallJobs(
	_ context.Context,
	_ *lcmipackageinstalljob.ListPackageInstallJobsRequest,
	_ ...grpc.CallOption,
) (*lcmipackageinstalljob.ListPackageInstallJobsResponse, error) {
	return &lcmipackageinstalljob.ListPackageInstallJobsResponse{
		PackageInstallJob: maps.Values(f.jobs),
	}, nil
}

func (f *PackageInstallJobClient) Status(
	_ context.Context,
	_ *lcmipackageinstalljob.StatusRequest,
	_ ...grpc.CallOption,
) (*lcmipackageinstalljob.StatusResponse, error) {
	return nil, nil
}

func (f *PackageInstallJobClient) Execute(
	_ context.Context,
	in *lcmipackageinstalljob.ExecuteRequest,
	_ ...grpc.CallOption,
) (*lcmipackageinstalljob.ExecuteResponse, error) {
	f.jobs[in.PackageInstallJob.Metadata.Id] = in.PackageInstallJob
	return &lcmipackageinstalljob.ExecuteResponse{
		Status: &lcmipackageinstalljob.PackageInstallJobStatus{
			Conditions: []*lcmimeta.Condition{
				{
					Type:               lcmicommon.ConditionType_CONDITION_TYPE_PENDING,
					Status:             "True",
					LastTransitionTime: time.Now().UnixNano(),
					Message:            "",
					Reason:             "",
				},
			},
			State:   lcmicommon.PackageInstallJobState_PACKAGE_INSTALL_JOB_STATE_UNSPECIFIED,
			Message: "",
		},
	}, nil
}
