// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package convertutil

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TimeToTimestamp(in metav1.Time) metav1.Timestamp {
	return metav1.Timestamp{Seconds: in.Unix()}
}

func TimeToTimestampPtr(in metav1.Time) *metav1.Timestamp {
	return &metav1.Timestamp{Seconds: in.Unix()}
}

func TimestampToTime(in metav1.Timestamp) metav1.Time {
	return metav1.Time{Time: time.Unix(in.Seconds, int64(in.Nanos))}
}

func TimestampToTimePtr(in metav1.Timestamp) *metav1.Time {
	return &metav1.Time{Time: time.Unix(in.Seconds, int64(in.Nanos))}
}
