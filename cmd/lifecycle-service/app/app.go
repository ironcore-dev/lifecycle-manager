package app

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/ironcore-dev/lifecycle-manager/lcmi/server"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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

type Options struct {
	kubeconfig string
	logLevel   string
	logFormat  string
	port       int
	namespace  string
	scanPeriod time.Duration
	horizon    time.Duration
}

func (o *Options) addFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.kubeconfig, "kubeconfig", "", "path to kubeconfig file")
	fs.StringVar(&o.logLevel, "log-level", "info", "logging level")
	fs.StringVar(&o.logFormat, "log-format", "json", "logging format")
	fs.IntVar(&o.port, "port", 26500, "bind port")
	fs.StringVar(&o.namespace, "namespace", "default", "default namespace name")
	fs.DurationVar(&o.scanPeriod, "scan-period", time.Hour*24, "scan period")
	fs.DurationVar(&o.horizon, "horizon", time.Minute*30, "allowed lag for scan period check")
}

func Command() *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use: "lifecycle-service",
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(context.TODO(), opts)
		},
	}

	fs := pflag.NewFlagSet("", 0)
	cmd.PersistentFlags().AddFlagSet(fs)
	opts.addFlags(cmd.Flags())

	return cmd
}

func Run(ctx context.Context, opts Options) error {
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	srvOpts := server.Options{
		Cfg:        cfg,
		Log:        setupLogger(LogFormat(opts.logFormat), logLevelMapping[opts.logLevel]),
		Port:       opts.port,
		Namespace:  opts.namespace,
		ScanPeriod: opts.scanPeriod,
		Horizon:    opts.horizon,
	}
	srv := server.NewLifecycleGRPCServer(srvOpts)
	return srv.Start(ctx)
}

func setupLogger(format LogFormat, level slog.Leveler) *slog.Logger {
	switch format {
	case JSON:
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     level,
		}))
	case Text:
		return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     level,
		}))
	}
	return nil
}
