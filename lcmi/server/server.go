// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/go-logr/logr"
	"google.golang.org/grpc"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/jellydator/ttlcache/v3"

	machinev1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/machine/v1alpha1"
	machinetypev1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/machinetype/v1alpha1"
	storagev1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/storage/v1alpha1"
	"github.com/ironcore-dev/lifecycle-manager/lcmi/server/machine"
	"github.com/ironcore-dev/lifecycle-manager/lcmi/server/machinetype"
	"github.com/ironcore-dev/lifecycle-manager/lcmi/server/storage"
)

const cacheCapacity = 1024

type LifecycleGRPCServer struct {
	log                    logr.Logger
	machineGrpcService     *machine.GrpcService
	machinetypeGrpcService *machinetype.GrpcService
	storageGrpcService     *storage.GrpcService
	machineCache           *ttlcache.Cache[string, *machinev1alpha1.MachineStatus]
	machinetypeCache       *ttlcache.Cache[string, *machinetypev1alpha1.MachineTypeStatus]
	port                   int
}

type Options struct {
	Cfg        *rest.Config
	Log        logr.Logger
	Port       int
	Namespace  string
	ScanPeriod time.Duration
	Horizon    time.Duration
}

func NewLifecycleGRPCServer(opts Options) *LifecycleGRPCServer {
	srv := &LifecycleGRPCServer{
		log:  opts.Log,
		port: opts.Port,
	}
	machineCache := ttlcache.New[string, *machinev1alpha1.MachineStatus](
		ttlcache.WithTTL[string, *machinev1alpha1.MachineStatus](opts.ScanPeriod),
		ttlcache.WithDisableTouchOnHit[string, *machinev1alpha1.MachineStatus](),
		ttlcache.WithCapacity[string, *machinev1alpha1.MachineStatus](cacheCapacity))
	machinetypeCache := ttlcache.New[string, *machinetypev1alpha1.MachineTypeStatus](
		ttlcache.WithTTL[string, *machinetypev1alpha1.MachineTypeStatus](opts.ScanPeriod),
		ttlcache.WithDisableTouchOnHit[string, *machinetypev1alpha1.MachineTypeStatus](),
		ttlcache.WithCapacity[string, *machinetypev1alpha1.MachineTypeStatus](cacheCapacity))
	machineGrpcService := machine.NewGrpcService(opts.Cfg,
		machine.WithNamespace(opts.Namespace),
		machine.WithHorizon(opts.Horizon),
		machine.WithScanPeriod(opts.ScanPeriod),
		machine.WithCache(machineCache))
	machinetypeGrpcService := machinetype.NewGrpcService()
	storageGrpcService := storage.NewGrpcService()
	srv.machineCache = machineCache
	srv.machinetypeCache = machinetypeCache
	srv.machineGrpcService = machineGrpcService
	srv.machinetypeGrpcService = machinetypeGrpcService
	srv.storageGrpcService = storageGrpcService
	return srv
}

func (s *LifecycleGRPCServer) Start(ctx context.Context) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		s.log.Error(err, "failed to bind port for listener", "port", s.port)
		return err
	}

	// todo: run scheduler for scan/install jobs

	srv := grpc.NewServer(grpc.UnaryInterceptor(s.addLogger))
	machinev1alpha1.RegisterMachineServiceServer(srv, s.machineGrpcService)
	machinetypev1alpha1.RegisterMachineTypeServiceServer(srv, s.machinetypeGrpcService)
	storagev1alpha1.RegisterFirmwareStorageServiceServer(srv, s.storageGrpcService)

	go func() {
		defer func() {
			s.machineCache.Stop()
			s.machinetypeCache.Stop()
			s.log.V(0).Info("stopping server", "kind", "lifecycle-grpc-server")
			srv.GracefulStop()
			s.log.V(0).Info("stopped")
		}()
		<-ctx.Done()
	}()

	go s.machineCache.Start()
	go s.machinetypeCache.Start()

	s.log.V(0).Info("starting server", "kind", "lifecycle-grpc-server", "addr", listener.Addr().String())
	if err = srv.Serve(listener); err != nil {
		s.log.Error(err, "failed to serve")
	}
	return nil
}

func (s *LifecycleGRPCServer) addLogger(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	log := s.log.WithName(info.FullMethod)
	ctx = ctrl.LoggerInto(ctx, log)
	log.V(0).Info("request")
	resp, err := handler(ctx, req)
	if err != nil {
		log.Error(err, "failed to handle request")
	}
	return resp, err
}
