// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"

	"connectrpc.com/connect"
	lifecyclev1alpha1 "github.com/ironcore-dev/lifecycle-manager/api/lifecycle/v1alpha1"
	"github.com/ironcore-dev/lifecycle-manager/clientgo/connectrpc/machine/v1alpha1/machinev1alpha1connect"
	"github.com/ironcore-dev/lifecycle-manager/clientgo/connectrpc/machinetype/v1alpha1/machinetypev1alpha1connect"
	"github.com/ironcore-dev/lifecycle-manager/internal/job"
	oobv1alpha1 "github.com/ironcore-dev/oob/api/v1alpha1"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/net/http2"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

type LogFormat string

const (
	JSON LogFormat = "json"
	Text LogFormat = "text"
)

var logLevelMapping = map[string]slog.Leveler{
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
}

var scheme = runtime.NewScheme()

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(lifecyclev1alpha1.AddToScheme(scheme))
	utilruntime.Must(oobv1alpha1.AddToScheme(scheme))
}

type Worker interface {
	Start(ctx context.Context) error
}

var lcmEndpoint = "http://lifecycle-service-svc:8080"

type Options struct {
	kubeconfig  string
	logLevel    string
	logFormat   string
	lcmEndpoint string
	targetType  string
	jobId       string
	dev         bool
}

func (o *Options) addFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	fs.StringVar(&o.logLevel, "log-level", "info", "logging level")
	fs.StringVar(&o.logFormat, "log-format", "json", "logging format")
	fs.StringVar(&o.lcmEndpoint, "lcm-endpoint", lcmEndpoint, "lcm endpoint")
	fs.StringVar(&o.jobId, "job-id", "", "job id")
	fs.StringVar(&o.targetType, "target-type", "", "target type")
	fs.BoolVar(&o.dev, "dev", false, "development mode")
}

func Command() *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use: "lcmjob",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			return Run(ctx, opts)
		},
	}

	fs := pflag.NewFlagSet("", 0)
	cmd.PersistentFlags().AddFlagSet(fs)
	opts.addFlags(cmd.Flags())

	return cmd
}

func Run(ctx context.Context, opts Options) error {
	var w Worker
	cfg := config.GetConfigOrDie()
	cl, err := client.New(cfg, client.Options{Scheme: scheme})
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	workerOpts := job.Options{
		KubeClient: cl,
		Log:        setupLogger(LogFormat(opts.logFormat), logLevelMapping[opts.logLevel], opts.dev),
		JobId:      opts.jobId,
	}
	switch opts.targetType {
	case "machine":
		w = job.NewMachineLifecycleWorker(workerOpts).
			WithClient(setupMachineClient(opts.lcmEndpoint, setupHTTPClient()))
	case "machinetype":
		w = job.NewMachineTypeLifecycleWorker(workerOpts).
			WithClient(setupMachineTypeClient(opts.lcmEndpoint, setupHTTPClient()))
	}
	if w == nil {
		return fmt.Errorf("no worker implementation")
	}
	return w.Start(ctx)
}

func setupLogger(format LogFormat, level slog.Leveler, dev bool) *slog.Logger {
	switch format {
	case JSON:
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: dev,
			Level:     level,
		}))
	case Text:
		return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: dev,
			Level:     level,
		}))
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
