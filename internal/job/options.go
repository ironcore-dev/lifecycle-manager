// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package job

import (
	"log/slog"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Options struct {
	KubeClient client.Client
	Log        *slog.Logger
	JobId      string
}
