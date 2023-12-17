// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/ironcore-dev/lifecycle-manager/api/v1alpha1"
	lcmicommon "github.com/ironcore-dev/lifecycle-manager/lcmi/api/common/v1alpha1"
	lcmimeta "github.com/ironcore-dev/lifecycle-manager/lcmi/api/meta/v1alpha1"
	lcmipackageinstalljob "github.com/ironcore-dev/lifecycle-manager/lcmi/api/packageinstalljob/v1alpha1"
	"github.com/ironcore-dev/lifecycle-manager/lcmi/fake"
	"github.com/ironcore-dev/lifecycle-manager/util/uuidutil"
)

func TestPackageInstallJobReconciler_Reconcile(t *testing.T) {
	t.Parallel()

	tn := time.Now()
	jobPendingTimestamp := tn.Add(-10 * time.Minute)
	jobScheduledTimestamp := tn.Add(-5 * time.Minute)
	jobInProgressTimestamp := tn.Add(-1 * time.Minute)

	schemeOpts := []schemeOption{withGroupVersion(v1alpha1.SchemeBuilder)}
	tests := map[string]struct {
		target       *v1alpha1.PackageInstallJob
		brokerClient *fake.PackageInstallJobClient
		request      ctrl.Request
		expectError  bool
	}{
		"absent-job": {
			target:      &v1alpha1.PackageInstallJob{},
			request:     ctrl.Request{NamespacedName: types.NamespacedName{Name: "absent-job"}},
			expectError: false,
		},
		"deleting-job": {
			target: &v1alpha1.PackageInstallJob{
				TypeMeta: metav1.TypeMeta{
					Kind:       "PackageInstallJob",
					APIVersion: "v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:              "sample-job",
					Namespace:         "metal",
					DeletionTimestamp: &metav1.Time{Time: time.Now()},
					Finalizers:        []string{"fake-finalizer"},
				},
				Spec: v1alpha1.PackageInstallJobSpec{
					MachineRef:         corev1.LocalObjectReference{Name: "sample-machine"},
					FirmwarePackageRef: corev1.LocalObjectReference{Name: "sample-package"},
				},
			},
			request:     ctrl.Request{NamespacedName: types.NamespacedName{Name: "sample-job", Namespace: "metal"}},
			expectError: false,
		},
		"new-job-runner": {
			target: &v1alpha1.PackageInstallJob{
				TypeMeta: metav1.TypeMeta{
					Kind:       "PackageInstallJob",
					APIVersion: "v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sample-job",
					Namespace: "metal",
				},
				Spec: v1alpha1.PackageInstallJobSpec{
					MachineRef:         corev1.LocalObjectReference{Name: "sample-machine"},
					FirmwarePackageRef: corev1.LocalObjectReference{Name: "sample-package"},
				},
			},
			brokerClient: fake.NewPackageInstallJobClient(),
			request:      ctrl.Request{NamespacedName: types.NamespacedName{Name: "sample-job", Namespace: "metal"}},
			expectError:  false,
		},
		"update-job": {
			target: &v1alpha1.PackageInstallJob{
				TypeMeta: metav1.TypeMeta{
					Kind:       "PackageInstallJob",
					APIVersion: "v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sample-job",
					Namespace: "metal",
				},
				Spec: v1alpha1.PackageInstallJobSpec{
					MachineRef:         corev1.LocalObjectReference{Name: "sample-machine"},
					FirmwarePackageRef: corev1.LocalObjectReference{Name: "sample-package"},
				},
				Status: v1alpha1.PackageInstallJobStatus{
					Conditions: []v1alpha1.PackageInstallJobCondition{
						{
							Condition: metav1.Condition{
								Type:               v1alpha1.ConditionTypePending,
								Status:             "True",
								ObservedGeneration: 1,
								LastTransitionTime: metav1.Time{Time: time.Now()},
								Reason:             "",
								Message:            "",
							},
						},
					},
					State:   "",
					Message: "",
				},
			},
			brokerClient: fake.NewPackageInstallJobClientWithJobs(map[string]*lcmipackageinstalljob.PackageInstallJob{
				uuidutil.UUIDFromObjectKey(types.NamespacedName{Name: "sample-job", Namespace: "metal"}): {
					Metadata: &lcmimeta.ObjectMetadata{
						Id: uuidutil.UUIDFromObjectKey(types.NamespacedName{Name: "sample-job", Namespace: "metal"}),
					},
					Spec: &lcmipackageinstalljob.PackageInstallJobSpec{
						MachineRef:         &lcmimeta.LocalObjectReference{Name: "sample-machine"},
						FirmwarePackageRef: &lcmimeta.LocalObjectReference{Name: "sample-package"},
					},
					Status: &lcmipackageinstalljob.PackageInstallJobStatus{
						Conditions: []*lcmimeta.Condition{
							{
								Type:               lcmicommon.ConditionType_CONDITION_TYPE_PENDING,
								Status:             "True",
								LastTransitionTime: jobPendingTimestamp.UnixNano(),
							},
							{
								Type:               lcmicommon.ConditionType_CONDITION_TYPE_SCHEDULED,
								Status:             "True",
								LastTransitionTime: jobScheduledTimestamp.UnixNano(),
							},
							{
								Type:               lcmicommon.ConditionType_CONDITION_TYPE_IN_PROGRESS,
								Status:             "True",
								LastTransitionTime: jobInProgressTimestamp.UnixNano(),
							},
						},
						State:   0,
						Message: "",
					},
				},
			}),
			request:     ctrl.Request{NamespacedName: types.NamespacedName{Name: "sample-job", Namespace: "metal"}},
			expectError: false,
		},
	}

	for n, tc := range tests {
		name, testCase := n, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			clientOpts := []clientOption{withRuntimeObject(testCase.target)}
			r := newPackageInstallJobReconciler(t, schemeOpts, clientOpts)
			r.BrokerClient = testCase.brokerClient
			resp, err := r.Reconcile(context.Background(), testCase.request)
			if testCase.expectError {
				assert.Error(t, err)
				assert.Empty(t, resp)
				return
			}
			assert.NoError(t, err)
			assert.Empty(t, resp)
			switch name {
			case "new-job-runner":
				actualJob := &v1alpha1.PackageInstallJob{}
				err = r.Get(context.Background(), client.ObjectKeyFromObject(testCase.target), actualJob)
				assert.NoError(t, err)
				assert.Equal(t, 1, len(actualJob.Status.Conditions))
				lcmiJobs, err := r.BrokerClient.ListPackageInstallJobs(
					context.Background(),
					&lcmipackageinstalljob.ListPackageInstallJobsRequest{
						Filter: &lcmipackageinstalljob.PackageInstallJobFilter{
							Id: uuidutil.UUIDFromObjectKey(client.ObjectKeyFromObject(actualJob))},
					},
				)
				assert.NoError(t, err)
				assert.Equal(t, 1, len(lcmiJobs.PackageInstallJob))
			case "update-job":
				actualJob := &v1alpha1.PackageInstallJob{}
				err = r.Get(context.Background(), client.ObjectKeyFromObject(testCase.target), actualJob)
				assert.NoError(t, err)
				assert.Equal(t, 3, len(actualJob.Status.Conditions))
				lcmiJobs, err := r.BrokerClient.ListPackageInstallJobs(
					context.Background(),
					&lcmipackageinstalljob.ListPackageInstallJobsRequest{
						Filter: &lcmipackageinstalljob.PackageInstallJobFilter{
							Id: uuidutil.UUIDFromObjectKey(client.ObjectKeyFromObject(actualJob))},
					},
				)
				assert.NoError(t, err)
				assert.Equal(t, 1, len(lcmiJobs.PackageInstallJob))
				for _, condition := range actualJob.Status.Conditions {
					assert.True(t, condition.IsTrue())
				}
			}
		})
	}
}
