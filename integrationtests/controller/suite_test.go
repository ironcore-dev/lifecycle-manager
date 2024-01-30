// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	oobv1alpha1 "github.com/onmetal/oob-operator/api/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	lifecyclev1alpha1 "github.com/ironcore-dev/lifecycle-manager/api/lifecycle/v1alpha1"
	"github.com/ironcore-dev/lifecycle-manager/internal/controllers"
	"github.com/ironcore-dev/lifecycle-manager/util/testutil"
	// +kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var (
	lifecycleCRDPath string
	oobCRDPath       string
	// todo: scheme *runtime.Scheme
	cfg        *rest.Config
	k8sClient  client.Client
	k8sManager manager.Manager
	testEnv    *envtest.Environment
	ctx        context.Context
	cancel     context.CancelFunc
	err        error
)

var scanPeriod = metav1.Duration{Duration: time.Second}

const (
	timeout       = time.Second * 3
	interval      = time.Millisecond * 50
	requeuePeriod = time.Second
	namespace     = "default"
)

func TestControllers(t *testing.T) {
	t.Parallel()
	RegisterFailHandler(Fail)
	RunSpecs(t, "Controller Suite")
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	By("bootstrapping test environment")
	lifecycleCRDPath = filepath.Join("..", "..", "config", "crd", "bases")
	oobCRDPath, err = testutil.GetCrdPath(oobv1alpha1.OOB{})
	Expect(err).NotTo(HaveOccurred())
	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{lifecycleCRDPath, oobCRDPath},
		ErrorIfCRDPathMissing: true,
	}
	ctx, cancel = context.WithCancel(context.TODO())

	// cfg is defined in this file globally.
	cfg, err = testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	scheme := runtime.NewScheme()
	Expect(lifecyclev1alpha1.AddToScheme(scheme)).To(Succeed())
	Expect(oobv1alpha1.AddToScheme(scheme)).To(Succeed())

	// +kubebuilder:scaffold:scheme

	k8sManager, err = ctrl.NewManager(cfg, ctrl.Options{
		Scheme: scheme,
	})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sManager).NotTo(BeNil())

	k8sClient = k8sManager.GetClient()
	Expect(k8sClient).NotTo(BeNil())

	Expect((&controllers.OnboardingReconciler{
		Client:        k8sClient,
		Scheme:        scheme,
		RequeuePeriod: requeuePeriod,
		ScanPeriod:    scanPeriod,
	}).SetupWithManager(k8sManager)).To(Succeed())
	Expect((&controllers.MachineTypeReconciler{
		Client: k8sClient,
		Scheme: scheme,
		Broker: nil, // todo: setup broker client
	}).SetupWithManager(k8sManager)).To(Succeed())
	Expect((&controllers.MachineReconciler{
		Client:               k8sClient,
		Scheme:               scheme,
		MachineServiceClient: nil, // todo: setup broker client
		Namespace:            namespace,
	}).SetupWithManager(k8sManager)).To(Succeed())

	go func() {
		defer GinkgoRecover()
		Expect(k8sManager.Start(ctx)).To(Succeed())
	}()
})

var _ = AfterSuite(func() {
	cancel()
	By("tearing down the test environment")
	Expect(testEnv.Stop()).To(Succeed())
})
