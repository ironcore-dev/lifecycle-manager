// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package testutil

import (
	"go/build"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"golang.org/x/mod/modfile"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

type ClientOption func(*fake.ClientBuilder)

func WithRuntimeObject(object client.Object) ClientOption {
	return func(b *fake.ClientBuilder) {
		b.WithRuntimeObjects(object)
		b.WithStatusSubresource(object)
	}
}

type SchemeOption func(*runtime.Scheme)

func WithGroupVersion(addToScheme func(*runtime.Scheme) error) SchemeOption {
	return func(s *runtime.Scheme) {
		if err := addToScheme(s); err != nil {
			panic("Fatal")
		}
	}
}

func SetupScheme(opts ...SchemeOption) *runtime.Scheme {
	s := runtime.NewScheme()
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func SetupClient(scheme *runtime.Scheme, opts ...ClientOption) client.Client {
	builder := fake.NewClientBuilder()
	builder.WithScheme(scheme)
	for _, opt := range opts {
		opt(builder)
	}
	cl := builder.Build()
	return cl
}

func GetCrdPath(crdPackageScheme interface{}) (string, error) {
	globalPackagePath := reflect.TypeOf(crdPackageScheme).PkgPath()
	goModData, err := os.ReadFile(filepath.Join("..", "..", "go.mod"))
	if err != nil {
		return "", err
	}
	goModFile, err := modfile.ParseLax("", goModData, nil)
	if err != nil {
		return "", err
	}
	globalModulePath := ""
	for _, req := range goModFile.Require {
		if strings.HasPrefix(globalPackagePath, req.Mod.Path) {
			globalModulePath = req.Mod.String()
			break
		}
	}
	return filepath.Join(build.Default.GOPATH, "pkg", "mod", globalModulePath, "config", "crd", "bases"), nil
}
