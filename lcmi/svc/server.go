// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package svc

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"connectrpc.com/connect"
	"connectrpc.com/grpchealth"
	"connectrpc.com/grpcreflect"
	"connectrpc.com/validate"
	lifecyclev1alpha1 "github.com/ironcore-dev/lifecycle-manager/api/lifecycle/v1alpha1"
	"github.com/ironcore-dev/lifecycle-manager/lcmi/api/machine/v1alpha1/machinev1alpha1connect"
	"github.com/ironcore-dev/lifecycle-manager/lcmi/api/machinetype/v1alpha1/machinetypev1alpha1connect"
	"github.com/ironcore-dev/lifecycle-manager/lcmi/svc/interceptor"
	machinesvcv1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/svc/machine/v1alpha1"
	machinetypesvcv1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/svc/machinetype/v1alpha1"
	"github.com/ironcore-dev/lifecycle-manager/lcmi/svc/scheduler"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"k8s.io/client-go/rest"
)

type GrpcServer struct {
	log                *slog.Logger
	host               string
	port               int
	machineService     *machinesvcv1alpha1.MachineService
	machineTypeService *machinetypesvcv1alpha1.MachineTypeService
}

type Options struct {
	Cfg  *rest.Config
	Log  *slog.Logger
	Host string
	Port int

	Namespace     string
	Workers       uint64
	Horizon       time.Duration
	QueueCapacity uint64
}

func NewGrpcServer(opts Options) *GrpcServer {
	srv := &GrpcServer{
		log:  opts.Log,
		host: opts.Host,
		port: opts.Port,
	}
	srv.machineService = setupMachineService(opts)
	srv.machineTypeService = setupMachineTypeService(opts)
	return srv
}

func (s *GrpcServer) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	reflector := grpcreflect.NewStaticReflector(Names...)
	checker := grpchealth.NewStaticChecker(Names...)

	validator, err := validate.NewInterceptor()
	if err != nil {
		s.log.Error("failed to create validator", "error", err.Error())
		return err
	}
	logger := interceptor.NewLoggerInterceptor(s.log)

	// enable services
	mux.Handle(machinev1alpha1connect.NewMachineServiceHandler(s.machineService,
		connect.WithInterceptors(logger, validator)))
	mux.Handle(machinetypev1alpha1connect.NewMachineTypeServiceHandler(s.machineTypeService,
		connect.WithInterceptors(logger, validator)))

	// enable health checks
	mux.Handle(grpchealth.NewHandler(checker))

	// enable reflection for gRPC server
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	srv := &http2.Server{}

	go func() {
		defer func() {
			s.log.Debug("stopping server", "kind", "lifecycle-service")
			s.log.Info("server stopped")
			os.Exit(0)
		}()
		<-ctx.Done()
	}()

	go s.machineService.StartScheduler(ctx)
	go s.machineTypeService.StartScheduler(ctx)

	s.log.Info("start serving", "addr", fmt.Sprintf("%s:%d", s.host, s.port))
	return http.ListenAndServe(fmt.Sprintf("%s:%d", s.host, s.port), h2c.NewHandler(mux, srv))
}

func setupMachineService(opts Options) *machinesvcv1alpha1.MachineService {
	machineScheduler := scheduler.NewScheduler[*lifecyclev1alpha1.Machine](
		opts.Log.With("scheduler", "Machine"), opts.Cfg, opts.Namespace,
		scheduler.WithWorkerCount[*lifecyclev1alpha1.Machine](opts.Workers),
		scheduler.WithActiveJobCache[*lifecyclev1alpha1.Machine](opts.Workers, opts.Horizon),
		scheduler.WithQueueCapacity[*lifecyclev1alpha1.Machine](opts.QueueCapacity))
	machineService := machinesvcv1alpha1.NewService(opts.Cfg,
		machinesvcv1alpha1.WithNamespace(opts.Namespace),
		machinesvcv1alpha1.WithHorizon(opts.Horizon),
		machinesvcv1alpha1.WithScheduler(machineScheduler))
	return machineService
}

func setupMachineTypeService(opts Options) *machinetypesvcv1alpha1.MachineTypeService {
	machinetypeScheduler := scheduler.NewScheduler[*lifecyclev1alpha1.MachineType](
		opts.Log.With("scheduler", "MachineType"), opts.Cfg, opts.Namespace,
		scheduler.WithWorkerCount[*lifecyclev1alpha1.MachineType](opts.Workers),
		scheduler.WithActiveJobCache[*lifecyclev1alpha1.MachineType](opts.Workers, opts.Horizon),
		scheduler.WithQueueCapacity[*lifecyclev1alpha1.MachineType](opts.QueueCapacity))
	machinetypeService := machinetypesvcv1alpha1.NewService(opts.Cfg,
		machinetypesvcv1alpha1.WithNamespace(opts.Namespace),
		machinetypesvcv1alpha1.WithHorizon(opts.Horizon),
		machinetypesvcv1alpha1.WithScheduler(machinetypeScheduler))
	return machinetypeService
}
