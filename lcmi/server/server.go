package server

import (
	"context"
	"fmt"
	"net"

	"github.com/go-logr/logr"
	"google.golang.org/grpc"
	ctrl "sigs.k8s.io/controller-runtime"

	machinev1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/machine/v1alpha1"
	machinetypev1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/machinetype/v1alpha1"
	"github.com/ironcore-dev/lifecycle-manager/lcmi/server/machine"
	"github.com/ironcore-dev/lifecycle-manager/lcmi/server/machinetype"
)

type LifecycleGRPCServer struct {
	log                    logr.Logger
	machineGrpcService     *machine.GrpcService
	machinetypeGrpcService *machinetype.GrpcService
	port                   int
}

type Options struct {
	Port int
}

func NewLifecycleGRPCServer(log logr.Logger, opts Options) *LifecycleGRPCServer {
	return &LifecycleGRPCServer{
		log:                    log,
		machineGrpcService:     machine.NewGRPCService(),
		machinetypeGrpcService: machinetype.NewGRPCService(),
		port:                   opts.Port,
	}
}

func (s *LifecycleGRPCServer) Start(ctx context.Context) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		s.log.Error(err, "failed to bind port for listener", "port", s.port)
		return err
	}

	srv := grpc.NewServer(grpc.UnaryInterceptor(s.addLogger))
	machinev1alpha1.RegisterMachineServiceServer(srv, s.machineGrpcService)
	machinetypev1alpha1.RegisterMachineTypeServiceServer(srv, s.machinetypeGrpcService)

	go func() {
		defer func() {
			s.log.V(0).Info("stopping server", "kind", "lifecycle-grpc-server")
			srv.GracefulStop()
			s.log.V(0).Info("stopped")
		}()
		<-ctx.Done()
	}()
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