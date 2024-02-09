// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"

	"github.com/ironcore-dev/lifecycle-manager/cmd/lifecycle-service/app"
	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	controllerruntime "sigs.k8s.io/controller-runtime"
	// +kubebuilder:scaffold:imports
)

func main() {
	ctx := controllerruntime.SetupSignalHandler()

	if err := app.Command().ExecuteContext(ctx); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
