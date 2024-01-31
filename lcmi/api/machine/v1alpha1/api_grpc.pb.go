// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: machine/v1alpha1/api.proto

package machinev1alpha1

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	MachineService_ScanMachine_FullMethodName          = "/machine.v1alpha1.MachineService/ScanMachine"
	MachineService_Install_FullMethodName              = "/machine.v1alpha1.MachineService/Install"
	MachineService_UpdateMachineStatus_FullMethodName  = "/machine.v1alpha1.MachineService/UpdateMachineStatus"
	MachineService_ListMachines_FullMethodName         = "/machine.v1alpha1.MachineService/ListMachines"
	MachineService_AddPackageVersion_FullMethodName    = "/machine.v1alpha1.MachineService/AddPackageVersion"
	MachineService_SetPackageVersion_FullMethodName    = "/machine.v1alpha1.MachineService/SetPackageVersion"
	MachineService_RemovePackageVersion_FullMethodName = "/machine.v1alpha1.MachineService/RemovePackageVersion"
)

// MachineServiceClient is the client API for MachineService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MachineServiceClient interface {
	ScanMachine(ctx context.Context, in *ScanMachineRequest, opts ...grpc.CallOption) (*ScanMachineResponse, error)
	Install(ctx context.Context, in *InstallRequest, opts ...grpc.CallOption) (*InstallResponse, error)
	UpdateMachineStatus(ctx context.Context, in *UpdateMachineStatusRequest, opts ...grpc.CallOption) (*UpdateMachineStatusResponse, error)
	ListMachines(ctx context.Context, in *ListMachinesRequest, opts ...grpc.CallOption) (*ListMachinesResponse, error)
	AddPackageVersion(ctx context.Context, in *AddPackageVersionRequest, opts ...grpc.CallOption) (*AddPackageVersionResponse, error)
	SetPackageVersion(ctx context.Context, in *SetPackageVersionRequest, opts ...grpc.CallOption) (*SetPackageVersionResponse, error)
	RemovePackageVersion(ctx context.Context, in *RemovePackageVersionRequest, opts ...grpc.CallOption) (*RemovePackageVersionResponse, error)
}

type machineServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMachineServiceClient(cc grpc.ClientConnInterface) MachineServiceClient {
	return &machineServiceClient{cc}
}

func (c *machineServiceClient) ScanMachine(ctx context.Context, in *ScanMachineRequest, opts ...grpc.CallOption) (*ScanMachineResponse, error) {
	out := new(ScanMachineResponse)
	err := c.cc.Invoke(ctx, MachineService_ScanMachine_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *machineServiceClient) Install(ctx context.Context, in *InstallRequest, opts ...grpc.CallOption) (*InstallResponse, error) {
	out := new(InstallResponse)
	err := c.cc.Invoke(ctx, MachineService_Install_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *machineServiceClient) UpdateMachineStatus(ctx context.Context, in *UpdateMachineStatusRequest, opts ...grpc.CallOption) (*UpdateMachineStatusResponse, error) {
	out := new(UpdateMachineStatusResponse)
	err := c.cc.Invoke(ctx, MachineService_UpdateMachineStatus_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *machineServiceClient) ListMachines(ctx context.Context, in *ListMachinesRequest, opts ...grpc.CallOption) (*ListMachinesResponse, error) {
	out := new(ListMachinesResponse)
	err := c.cc.Invoke(ctx, MachineService_ListMachines_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *machineServiceClient) AddPackageVersion(ctx context.Context, in *AddPackageVersionRequest, opts ...grpc.CallOption) (*AddPackageVersionResponse, error) {
	out := new(AddPackageVersionResponse)
	err := c.cc.Invoke(ctx, MachineService_AddPackageVersion_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *machineServiceClient) SetPackageVersion(ctx context.Context, in *SetPackageVersionRequest, opts ...grpc.CallOption) (*SetPackageVersionResponse, error) {
	out := new(SetPackageVersionResponse)
	err := c.cc.Invoke(ctx, MachineService_SetPackageVersion_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *machineServiceClient) RemovePackageVersion(ctx context.Context, in *RemovePackageVersionRequest, opts ...grpc.CallOption) (*RemovePackageVersionResponse, error) {
	out := new(RemovePackageVersionResponse)
	err := c.cc.Invoke(ctx, MachineService_RemovePackageVersion_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MachineServiceServer is the server API for MachineService service.
// All implementations must embed UnimplementedMachineServiceServer
// for forward compatibility
type MachineServiceServer interface {
	ScanMachine(context.Context, *ScanMachineRequest) (*ScanMachineResponse, error)
	Install(context.Context, *InstallRequest) (*InstallResponse, error)
	UpdateMachineStatus(context.Context, *UpdateMachineStatusRequest) (*UpdateMachineStatusResponse, error)
	ListMachines(context.Context, *ListMachinesRequest) (*ListMachinesResponse, error)
	AddPackageVersion(context.Context, *AddPackageVersionRequest) (*AddPackageVersionResponse, error)
	SetPackageVersion(context.Context, *SetPackageVersionRequest) (*SetPackageVersionResponse, error)
	RemovePackageVersion(context.Context, *RemovePackageVersionRequest) (*RemovePackageVersionResponse, error)
	mustEmbedUnimplementedMachineServiceServer()
}

// UnimplementedMachineServiceServer must be embedded to have forward compatible implementations.
type UnimplementedMachineServiceServer struct {
}

func (UnimplementedMachineServiceServer) ScanMachine(context.Context, *ScanMachineRequest) (*ScanMachineResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ScanMachine not implemented")
}
func (UnimplementedMachineServiceServer) Install(context.Context, *InstallRequest) (*InstallResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Install not implemented")
}
func (UnimplementedMachineServiceServer) UpdateMachineStatus(context.Context, *UpdateMachineStatusRequest) (*UpdateMachineStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateMachineStatus not implemented")
}
func (UnimplementedMachineServiceServer) ListMachines(context.Context, *ListMachinesRequest) (*ListMachinesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListMachines not implemented")
}
func (UnimplementedMachineServiceServer) AddPackageVersion(context.Context, *AddPackageVersionRequest) (*AddPackageVersionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddPackageVersion not implemented")
}
func (UnimplementedMachineServiceServer) SetPackageVersion(context.Context, *SetPackageVersionRequest) (*SetPackageVersionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetPackageVersion not implemented")
}
func (UnimplementedMachineServiceServer) RemovePackageVersion(context.Context, *RemovePackageVersionRequest) (*RemovePackageVersionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemovePackageVersion not implemented")
}
func (UnimplementedMachineServiceServer) mustEmbedUnimplementedMachineServiceServer() {}

// UnsafeMachineServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MachineServiceServer will
// result in compilation errors.
type UnsafeMachineServiceServer interface {
	mustEmbedUnimplementedMachineServiceServer()
}

func RegisterMachineServiceServer(s grpc.ServiceRegistrar, srv MachineServiceServer) {
	s.RegisterService(&MachineService_ServiceDesc, srv)
}

func _MachineService_ScanMachine_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ScanMachineRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MachineServiceServer).ScanMachine(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MachineService_ScanMachine_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MachineServiceServer).ScanMachine(ctx, req.(*ScanMachineRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MachineService_Install_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InstallRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MachineServiceServer).Install(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MachineService_Install_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MachineServiceServer).Install(ctx, req.(*InstallRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MachineService_UpdateMachineStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateMachineStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MachineServiceServer).UpdateMachineStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MachineService_UpdateMachineStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MachineServiceServer).UpdateMachineStatus(ctx, req.(*UpdateMachineStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MachineService_ListMachines_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListMachinesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MachineServiceServer).ListMachines(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MachineService_ListMachines_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MachineServiceServer).ListMachines(ctx, req.(*ListMachinesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MachineService_AddPackageVersion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddPackageVersionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MachineServiceServer).AddPackageVersion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MachineService_AddPackageVersion_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MachineServiceServer).AddPackageVersion(ctx, req.(*AddPackageVersionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MachineService_SetPackageVersion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetPackageVersionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MachineServiceServer).SetPackageVersion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MachineService_SetPackageVersion_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MachineServiceServer).SetPackageVersion(ctx, req.(*SetPackageVersionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MachineService_RemovePackageVersion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemovePackageVersionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MachineServiceServer).RemovePackageVersion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MachineService_RemovePackageVersion_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MachineServiceServer).RemovePackageVersion(ctx, req.(*RemovePackageVersionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MachineService_ServiceDesc is the grpc.ServiceDesc for MachineService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MachineService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "machine.v1alpha1.MachineService",
	HandlerType: (*MachineServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ScanMachine",
			Handler:    _MachineService_ScanMachine_Handler,
		},
		{
			MethodName: "Install",
			Handler:    _MachineService_Install_Handler,
		},
		{
			MethodName: "UpdateMachineStatus",
			Handler:    _MachineService_UpdateMachineStatus_Handler,
		},
		{
			MethodName: "ListMachines",
			Handler:    _MachineService_ListMachines_Handler,
		},
		{
			MethodName: "AddPackageVersion",
			Handler:    _MachineService_AddPackageVersion_Handler,
		},
		{
			MethodName: "SetPackageVersion",
			Handler:    _MachineService_SetPackageVersion_Handler,
		},
		{
			MethodName: "RemovePackageVersion",
			Handler:    _MachineService_RemovePackageVersion_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "machine/v1alpha1/api.proto",
}
