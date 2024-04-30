// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"

	"github.com/ironcore-dev/lifecycle-manager/cmd/lifecycle-job/app"
)

func main() {
	ctx := app.SetupSignalHandler()

	if err := app.Command().ExecuteContext(ctx); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
