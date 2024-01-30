// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: storage/v1alpha1/api.proto

package v1alpha1

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
	FirmwareStorageService_InitUpload_FullMethodName   = "/common.v1alpha1.FirmwareStorageService/InitUpload"
	FirmwareStorageService_Upload_FullMethodName       = "/common.v1alpha1.FirmwareStorageService/Upload"
	FirmwareStorageService_InitDownload_FullMethodName = "/common.v1alpha1.FirmwareStorageService/InitDownload"
	FirmwareStorageService_Download_FullMethodName     = "/common.v1alpha1.FirmwareStorageService/Download"
)

// FirmwareStorageServiceClient is the client API for FirmwareStorageService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FirmwareStorageServiceClient interface {
	InitUpload(ctx context.Context, in *InitUploadRequest, opts ...grpc.CallOption) (*InitUploadResponse, error)
	Upload(ctx context.Context, opts ...grpc.CallOption) (FirmwareStorageService_UploadClient, error)
	InitDownload(ctx context.Context, in *InitDownloadRequest, opts ...grpc.CallOption) (*InitDownloadResponse, error)
	Download(ctx context.Context, in *DownloadRequest, opts ...grpc.CallOption) (FirmwareStorageService_DownloadClient, error)
}

type firmwareStorageServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFirmwareStorageServiceClient(cc grpc.ClientConnInterface) FirmwareStorageServiceClient {
	return &firmwareStorageServiceClient{cc}
}

func (c *firmwareStorageServiceClient) InitUpload(ctx context.Context, in *InitUploadRequest, opts ...grpc.CallOption) (*InitUploadResponse, error) {
	out := new(InitUploadResponse)
	err := c.cc.Invoke(ctx, FirmwareStorageService_InitUpload_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *firmwareStorageServiceClient) Upload(ctx context.Context, opts ...grpc.CallOption) (FirmwareStorageService_UploadClient, error) {
	stream, err := c.cc.NewStream(ctx, &FirmwareStorageService_ServiceDesc.Streams[0], FirmwareStorageService_Upload_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &firmwareStorageServiceUploadClient{stream}
	return x, nil
}

type FirmwareStorageService_UploadClient interface {
	Send(*UploadRequest) error
	CloseAndRecv() (*UploadResponse, error)
	grpc.ClientStream
}

type firmwareStorageServiceUploadClient struct {
	grpc.ClientStream
}

func (x *firmwareStorageServiceUploadClient) Send(m *UploadRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *firmwareStorageServiceUploadClient) CloseAndRecv() (*UploadResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(UploadResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *firmwareStorageServiceClient) InitDownload(ctx context.Context, in *InitDownloadRequest, opts ...grpc.CallOption) (*InitDownloadResponse, error) {
	out := new(InitDownloadResponse)
	err := c.cc.Invoke(ctx, FirmwareStorageService_InitDownload_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *firmwareStorageServiceClient) Download(ctx context.Context, in *DownloadRequest, opts ...grpc.CallOption) (FirmwareStorageService_DownloadClient, error) {
	stream, err := c.cc.NewStream(ctx, &FirmwareStorageService_ServiceDesc.Streams[1], FirmwareStorageService_Download_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &firmwareStorageServiceDownloadClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type FirmwareStorageService_DownloadClient interface {
	Recv() (*DownloadResponse, error)
	grpc.ClientStream
}

type firmwareStorageServiceDownloadClient struct {
	grpc.ClientStream
}

func (x *firmwareStorageServiceDownloadClient) Recv() (*DownloadResponse, error) {
	m := new(DownloadResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// FirmwareStorageServiceServer is the server API for FirmwareStorageService service.
// All implementations must embed UnimplementedFirmwareStorageServiceServer
// for forward compatibility
type FirmwareStorageServiceServer interface {
	InitUpload(context.Context, *InitUploadRequest) (*InitUploadResponse, error)
	Upload(FirmwareStorageService_UploadServer) error
	InitDownload(context.Context, *InitDownloadRequest) (*InitDownloadResponse, error)
	Download(*DownloadRequest, FirmwareStorageService_DownloadServer) error
	mustEmbedUnimplementedFirmwareStorageServiceServer()
}

// UnimplementedFirmwareStorageServiceServer must be embedded to have forward compatible implementations.
type UnimplementedFirmwareStorageServiceServer struct {
}

func (UnimplementedFirmwareStorageServiceServer) InitUpload(context.Context, *InitUploadRequest) (*InitUploadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InitUpload not implemented")
}
func (UnimplementedFirmwareStorageServiceServer) Upload(FirmwareStorageService_UploadServer) error {
	return status.Errorf(codes.Unimplemented, "method Upload not implemented")
}
func (UnimplementedFirmwareStorageServiceServer) InitDownload(context.Context, *InitDownloadRequest) (*InitDownloadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InitDownload not implemented")
}
func (UnimplementedFirmwareStorageServiceServer) Download(*DownloadRequest, FirmwareStorageService_DownloadServer) error {
	return status.Errorf(codes.Unimplemented, "method Download not implemented")
}
func (UnimplementedFirmwareStorageServiceServer) mustEmbedUnimplementedFirmwareStorageServiceServer() {
}

// UnsafeFirmwareStorageServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FirmwareStorageServiceServer will
// result in compilation errors.
type UnsafeFirmwareStorageServiceServer interface {
	mustEmbedUnimplementedFirmwareStorageServiceServer()
}

func RegisterFirmwareStorageServiceServer(s grpc.ServiceRegistrar, srv FirmwareStorageServiceServer) {
	s.RegisterService(&FirmwareStorageService_ServiceDesc, srv)
}

func _FirmwareStorageService_InitUpload_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InitUploadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FirmwareStorageServiceServer).InitUpload(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FirmwareStorageService_InitUpload_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FirmwareStorageServiceServer).InitUpload(ctx, req.(*InitUploadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FirmwareStorageService_Upload_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(FirmwareStorageServiceServer).Upload(&firmwareStorageServiceUploadServer{stream})
}

type FirmwareStorageService_UploadServer interface {
	SendAndClose(*UploadResponse) error
	Recv() (*UploadRequest, error)
	grpc.ServerStream
}

type firmwareStorageServiceUploadServer struct {
	grpc.ServerStream
}

func (x *firmwareStorageServiceUploadServer) SendAndClose(m *UploadResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *firmwareStorageServiceUploadServer) Recv() (*UploadRequest, error) {
	m := new(UploadRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _FirmwareStorageService_InitDownload_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InitDownloadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FirmwareStorageServiceServer).InitDownload(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FirmwareStorageService_InitDownload_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FirmwareStorageServiceServer).InitDownload(ctx, req.(*InitDownloadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FirmwareStorageService_Download_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(DownloadRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(FirmwareStorageServiceServer).Download(m, &firmwareStorageServiceDownloadServer{stream})
}

type FirmwareStorageService_DownloadServer interface {
	Send(*DownloadResponse) error
	grpc.ServerStream
}

type firmwareStorageServiceDownloadServer struct {
	grpc.ServerStream
}

func (x *firmwareStorageServiceDownloadServer) Send(m *DownloadResponse) error {
	return x.ServerStream.SendMsg(m)
}

// FirmwareStorageService_ServiceDesc is the grpc.ServiceDesc for FirmwareStorageService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FirmwareStorageService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "common.v1alpha1.FirmwareStorageService",
	HandlerType: (*FirmwareStorageServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "InitUpload",
			Handler:    _FirmwareStorageService_InitUpload_Handler,
		},
		{
			MethodName: "InitDownload",
			Handler:    _FirmwareStorageService_InitDownload_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Upload",
			Handler:       _FirmwareStorageService_Upload_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "Download",
			Handler:       _FirmwareStorageService_Download_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "storage/v1alpha1/api.proto",
}
