// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"crypto/tls"
	"flag"
	"net"
	"net/http"
	"os"
	"time"

	"connectrpc.com/connect"
	"github.com/ironcore-dev/lifecycle-manager/lcmi/api/machine/v1alpha1/machinev1alpha1connect"
	"github.com/ironcore-dev/lifecycle-manager/lcmi/api/machinetype/v1alpha1/machinetypev1alpha1connect"
	oobv1alpha1 "github.com/ironcore-dev/oob/api/v1alpha1"
	"golang.org/x/net/http2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/healthz"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	lifecyclev1alpha1 "github.com/ironcore-dev/lifecycle-manager/api/lifecycle/v1alpha1"
	"github.com/ironcore-dev/lifecycle-manager/internal/controllers"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")

	lcmiEndpoint = "http://lifecycle-service-svc:8080"
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(lifecyclev1alpha1.AddToScheme(scheme))
	utilruntime.Must(oobv1alpha1.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	var lcmiServiceAddr string
	var horizon time.Duration
	flag.DurationVar(&horizon, "scan-horizon", time.Minute*10,
		"The period within which scan results considered to be valid.")
	flag.StringVar(&lcmiServiceAddr, "lcmi-address", lcmiEndpoint,
		"The address lifecycle-service running on.")
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080",
		"The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081",
		"The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		Metrics:                metricsserver.Options{BindAddress: metricsAddr},
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "86db2628.ironcore.dev",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = setupControllers(mgr, lcmiServiceAddr, horizon); err != nil {
		os.Exit(1)
	}
	if err = setupHandlers(mgr); err != nil {
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err = mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

func setupControllers(mgr ctrl.Manager, endpoint string, horizon time.Duration) error {
	httpClient := setupHTTPClient()
	if err := (&controllers.MachineReconciler{
		Client:               mgr.GetClient(),
		Scheme:               mgr.GetScheme(),
		Log:                  mgr.GetLogger().WithName("lifecycle-machine-controller"),
		MachineServiceClient: setupMachineClient(endpoint, httpClient),
		Horizon:              horizon,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Machine")
		return err
	}
	if err := (&controllers.MachineTypeReconciler{
		Client:                   mgr.GetClient(),
		Scheme:                   mgr.GetScheme(),
		Log:                      mgr.GetLogger().WithName("lifecycle-machinetype-controller"),
		MachineTypeServiceClient: setupMachineTypeClient(endpoint, httpClient),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "MachineType")
		return err
	}
	if err := (&controllers.OnboardingReconciler{
		Client:        mgr.GetClient(),
		Log:           mgr.GetLogger().WithName("lifecycle-onboarding-controller"),
		Scheme:        mgr.GetScheme(),
		RequeuePeriod: time.Minute,
		ScanPeriod:    v1.Duration{Duration: time.Hour * 24},
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Onboarding")
		return err
	}
	// +kubebuilder:scaffold:builder
	return nil
}

func setupHandlers(mgr ctrl.Manager) error {
	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		return err
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		return err
	}
	return nil
}

func setupHTTPClient() *http.Client {
	return &http.Client{
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLSContext: func(_ context.Context, network, addr string, _ *tls.Config) (net.Conn, error) {
				return net.Dial(network, addr)
			},
		},
	}
}

func setupMachineClient(endpoint string, cl *http.Client) machinev1alpha1connect.MachineServiceClient {
	return machinev1alpha1connect.NewMachineServiceClient(cl, endpoint, connect.WithGRPC())
}

func setupMachineTypeClient(endpoint string, cl *http.Client) machinetypev1alpha1connect.MachineTypeServiceClient {
	return machinetypev1alpha1connect.NewMachineTypeServiceClient(cl, endpoint, connect.WithGRPC())
}
