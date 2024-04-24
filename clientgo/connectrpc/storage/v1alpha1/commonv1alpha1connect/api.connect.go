// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: storage/v1alpha1/api.proto

package commonv1alpha1connect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	v1alpha1 "github.com/ironcore-dev/lifecycle-manager/api/proto/storage/v1alpha1"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion1_13_0

const (
	// FirmwareStorageServiceName is the fully-qualified name of the FirmwareStorageService service.
	FirmwareStorageServiceName = "common.v1alpha1.FirmwareStorageService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// FirmwareStorageServiceInitUploadProcedure is the fully-qualified name of the
	// FirmwareStorageService's InitUpload RPC.
	FirmwareStorageServiceInitUploadProcedure = "/common.v1alpha1.FirmwareStorageService/InitUpload"
	// FirmwareStorageServiceUploadProcedure is the fully-qualified name of the FirmwareStorageService's
	// Upload RPC.
	FirmwareStorageServiceUploadProcedure = "/common.v1alpha1.FirmwareStorageService/Upload"
	// FirmwareStorageServiceInitDownloadProcedure is the fully-qualified name of the
	// FirmwareStorageService's InitDownload RPC.
	FirmwareStorageServiceInitDownloadProcedure = "/common.v1alpha1.FirmwareStorageService/InitDownload"
	// FirmwareStorageServiceDownloadProcedure is the fully-qualified name of the
	// FirmwareStorageService's Download RPC.
	FirmwareStorageServiceDownloadProcedure = "/common.v1alpha1.FirmwareStorageService/Download"
)

// These variables are the protoreflect.Descriptor objects for the RPCs defined in this package.
var (
	firmwareStorageServiceServiceDescriptor            = v1alpha1.File_storage_v1alpha1_api_proto.Services().ByName("FirmwareStorageService")
	firmwareStorageServiceInitUploadMethodDescriptor   = firmwareStorageServiceServiceDescriptor.Methods().ByName("InitUpload")
	firmwareStorageServiceUploadMethodDescriptor       = firmwareStorageServiceServiceDescriptor.Methods().ByName("Upload")
	firmwareStorageServiceInitDownloadMethodDescriptor = firmwareStorageServiceServiceDescriptor.Methods().ByName("InitDownload")
	firmwareStorageServiceDownloadMethodDescriptor     = firmwareStorageServiceServiceDescriptor.Methods().ByName("Download")
)

// FirmwareStorageServiceClient is a client for the common.v1alpha1.FirmwareStorageService service.
type FirmwareStorageServiceClient interface {
	InitUpload(context.Context, *connect.Request[v1alpha1.InitUploadRequest]) (*connect.Response[v1alpha1.InitUploadResponse], error)
	Upload(context.Context) *connect.ClientStreamForClient[v1alpha1.UploadRequest, v1alpha1.UploadResponse]
	InitDownload(context.Context, *connect.Request[v1alpha1.InitDownloadRequest]) (*connect.Response[v1alpha1.InitDownloadResponse], error)
	Download(context.Context, *connect.Request[v1alpha1.DownloadRequest]) (*connect.ServerStreamForClient[v1alpha1.DownloadResponse], error)
}

// NewFirmwareStorageServiceClient constructs a client for the
// common.v1alpha1.FirmwareStorageService service. By default, it uses the Connect protocol with the
// binary Protobuf Codec, asks for gzipped responses, and sends uncompressed requests. To use the
// gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewFirmwareStorageServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) FirmwareStorageServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &firmwareStorageServiceClient{
		initUpload: connect.NewClient[v1alpha1.InitUploadRequest, v1alpha1.InitUploadResponse](
			httpClient,
			baseURL+FirmwareStorageServiceInitUploadProcedure,
			connect.WithSchema(firmwareStorageServiceInitUploadMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		upload: connect.NewClient[v1alpha1.UploadRequest, v1alpha1.UploadResponse](
			httpClient,
			baseURL+FirmwareStorageServiceUploadProcedure,
			connect.WithSchema(firmwareStorageServiceUploadMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		initDownload: connect.NewClient[v1alpha1.InitDownloadRequest, v1alpha1.InitDownloadResponse](
			httpClient,
			baseURL+FirmwareStorageServiceInitDownloadProcedure,
			connect.WithSchema(firmwareStorageServiceInitDownloadMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		download: connect.NewClient[v1alpha1.DownloadRequest, v1alpha1.DownloadResponse](
			httpClient,
			baseURL+FirmwareStorageServiceDownloadProcedure,
			connect.WithSchema(firmwareStorageServiceDownloadMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
	}
}

// firmwareStorageServiceClient implements FirmwareStorageServiceClient.
type firmwareStorageServiceClient struct {
	initUpload   *connect.Client[v1alpha1.InitUploadRequest, v1alpha1.InitUploadResponse]
	upload       *connect.Client[v1alpha1.UploadRequest, v1alpha1.UploadResponse]
	initDownload *connect.Client[v1alpha1.InitDownloadRequest, v1alpha1.InitDownloadResponse]
	download     *connect.Client[v1alpha1.DownloadRequest, v1alpha1.DownloadResponse]
}

// InitUpload calls common.v1alpha1.FirmwareStorageService.InitUpload.
func (c *firmwareStorageServiceClient) InitUpload(ctx context.Context, req *connect.Request[v1alpha1.InitUploadRequest]) (*connect.Response[v1alpha1.InitUploadResponse], error) {
	return c.initUpload.CallUnary(ctx, req)
}

// Upload calls common.v1alpha1.FirmwareStorageService.Upload.
func (c *firmwareStorageServiceClient) Upload(ctx context.Context) *connect.ClientStreamForClient[v1alpha1.UploadRequest, v1alpha1.UploadResponse] {
	return c.upload.CallClientStream(ctx)
}

// InitDownload calls common.v1alpha1.FirmwareStorageService.InitDownload.
func (c *firmwareStorageServiceClient) InitDownload(ctx context.Context, req *connect.Request[v1alpha1.InitDownloadRequest]) (*connect.Response[v1alpha1.InitDownloadResponse], error) {
	return c.initDownload.CallUnary(ctx, req)
}

// Download calls common.v1alpha1.FirmwareStorageService.Download.
func (c *firmwareStorageServiceClient) Download(ctx context.Context, req *connect.Request[v1alpha1.DownloadRequest]) (*connect.ServerStreamForClient[v1alpha1.DownloadResponse], error) {
	return c.download.CallServerStream(ctx, req)
}

// FirmwareStorageServiceHandler is an implementation of the common.v1alpha1.FirmwareStorageService
// service.
type FirmwareStorageServiceHandler interface {
	InitUpload(context.Context, *connect.Request[v1alpha1.InitUploadRequest]) (*connect.Response[v1alpha1.InitUploadResponse], error)
	Upload(context.Context, *connect.ClientStream[v1alpha1.UploadRequest]) (*connect.Response[v1alpha1.UploadResponse], error)
	InitDownload(context.Context, *connect.Request[v1alpha1.InitDownloadRequest]) (*connect.Response[v1alpha1.InitDownloadResponse], error)
	Download(context.Context, *connect.Request[v1alpha1.DownloadRequest], *connect.ServerStream[v1alpha1.DownloadResponse]) error
}

// NewFirmwareStorageServiceHandler builds an HTTP handler from the service implementation. It
// returns the path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewFirmwareStorageServiceHandler(svc FirmwareStorageServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	firmwareStorageServiceInitUploadHandler := connect.NewUnaryHandler(
		FirmwareStorageServiceInitUploadProcedure,
		svc.InitUpload,
		connect.WithSchema(firmwareStorageServiceInitUploadMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	firmwareStorageServiceUploadHandler := connect.NewClientStreamHandler(
		FirmwareStorageServiceUploadProcedure,
		svc.Upload,
		connect.WithSchema(firmwareStorageServiceUploadMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	firmwareStorageServiceInitDownloadHandler := connect.NewUnaryHandler(
		FirmwareStorageServiceInitDownloadProcedure,
		svc.InitDownload,
		connect.WithSchema(firmwareStorageServiceInitDownloadMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	firmwareStorageServiceDownloadHandler := connect.NewServerStreamHandler(
		FirmwareStorageServiceDownloadProcedure,
		svc.Download,
		connect.WithSchema(firmwareStorageServiceDownloadMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	return "/common.v1alpha1.FirmwareStorageService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case FirmwareStorageServiceInitUploadProcedure:
			firmwareStorageServiceInitUploadHandler.ServeHTTP(w, r)
		case FirmwareStorageServiceUploadProcedure:
			firmwareStorageServiceUploadHandler.ServeHTTP(w, r)
		case FirmwareStorageServiceInitDownloadProcedure:
			firmwareStorageServiceInitDownloadHandler.ServeHTTP(w, r)
		case FirmwareStorageServiceDownloadProcedure:
			firmwareStorageServiceDownloadHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedFirmwareStorageServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedFirmwareStorageServiceHandler struct{}

func (UnimplementedFirmwareStorageServiceHandler) InitUpload(context.Context, *connect.Request[v1alpha1.InitUploadRequest]) (*connect.Response[v1alpha1.InitUploadResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("common.v1alpha1.FirmwareStorageService.InitUpload is not implemented"))
}

func (UnimplementedFirmwareStorageServiceHandler) Upload(context.Context, *connect.ClientStream[v1alpha1.UploadRequest]) (*connect.Response[v1alpha1.UploadResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("common.v1alpha1.FirmwareStorageService.Upload is not implemented"))
}

func (UnimplementedFirmwareStorageServiceHandler) InitDownload(context.Context, *connect.Request[v1alpha1.InitDownloadRequest]) (*connect.Response[v1alpha1.InitDownloadResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("common.v1alpha1.FirmwareStorageService.InitDownload is not implemented"))
}

func (UnimplementedFirmwareStorageServiceHandler) Download(context.Context, *connect.Request[v1alpha1.DownloadRequest], *connect.ServerStream[v1alpha1.DownloadResponse]) error {
	return connect.NewError(connect.CodeUnimplemented, errors.New("common.v1alpha1.FirmwareStorageService.Download is not implemented"))
}
