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
	machineapiv1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/machine/v1alpha1"
	"github.com/ironcore-dev/lifecycle-manager/lcmi/api/machine/v1alpha1/machinev1alpha1connect"
	machinetypeapiv1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/machinetype/v1alpha1"
	"github.com/ironcore-dev/lifecycle-manager/lcmi/api/machinetype/v1alpha1/machinetypev1alpha1connect"
	"github.com/ironcore-dev/lifecycle-manager/lcmi/svc/interceptor"
	machinesvcv1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/svc/machine/v1alpha1"
	machinetypesvcv1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/svc/machinetype/v1alpha1"
	"github.com/jellydator/ttlcache/v3"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"k8s.io/client-go/rest"
)

const cacheCapacity = 1024

type GrpcServer struct {
	log                *slog.Logger
	host               string
	port               int
	machineService     *machinesvcv1alpha1.MachineService
	machineCache       *ttlcache.Cache[string, *machineapiv1alpha1.MachineStatus]
	machinetypeCache   *ttlcache.Cache[string, *machinetypeapiv1alpha1.MachineTypeStatus]
	machineTypeService *machinetypesvcv1alpha1.MachineTypeService
}

type Options struct {
	Cfg        *rest.Config
	Log        *slog.Logger
	Host       string
	Port       int
	Namespace  string
	ScanPeriod time.Duration
	Horizon    time.Duration
}

func NewGrpcServer(opts Options) *GrpcServer {
	srv := &GrpcServer{
		log:  opts.Log,
		host: opts.Host,
		port: opts.Port,
	}
	machineCache := ttlcache.New[string, *machineapiv1alpha1.MachineStatus](
		ttlcache.WithTTL[string, *machineapiv1alpha1.MachineStatus](opts.ScanPeriod),
		ttlcache.WithDisableTouchOnHit[string, *machineapiv1alpha1.MachineStatus](),
		ttlcache.WithCapacity[string, *machineapiv1alpha1.MachineStatus](cacheCapacity))
	machinetypeCache := ttlcache.New[string, *machinetypeapiv1alpha1.MachineTypeStatus](
		ttlcache.WithTTL[string, *machinetypeapiv1alpha1.MachineTypeStatus](opts.ScanPeriod),
		ttlcache.WithDisableTouchOnHit[string, *machinetypeapiv1alpha1.MachineTypeStatus](),
		ttlcache.WithCapacity[string, *machinetypeapiv1alpha1.MachineTypeStatus](cacheCapacity))
	machineService := machinesvcv1alpha1.NewService(opts.Cfg,
		machinesvcv1alpha1.WithNamespace(opts.Namespace),
		machinesvcv1alpha1.WithHorizon(opts.Horizon),
		machinesvcv1alpha1.WithScanPeriod(opts.ScanPeriod),
		machinesvcv1alpha1.WithCache(machineCache))
	srv.machineCache = machineCache
	srv.machinetypeCache = machinetypeCache
	srv.machineService = machineService
	return srv
}

func (s *GrpcServer) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	reflector := grpcreflect.NewStaticReflector(machinesvcv1alpha1.Names...)
	checker := grpchealth.NewStaticChecker(machinesvcv1alpha1.Names...)

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
			s.machineCache.Stop()
			s.machinetypeCache.Stop()
			s.log.Info("server stopped")
			os.Exit(0)
		}()
		<-ctx.Done()
	}()

	go s.machineCache.Start()
	go s.machinetypeCache.Start()

	s.log.Info("start serving", "addr", fmt.Sprintf("%s:%d", s.host, s.port))
	return http.ListenAndServe(fmt.Sprintf("%s:%d", s.host, s.port), h2c.NewHandler(mux, srv))
}
