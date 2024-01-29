package storage

import (
	"context"

	storagev1alpha1 "github.com/ironcore-dev/lifecycle-manager/lcmi/api/storage/v1alpha1"
)

type GrpcService struct {
	storagev1alpha1.UnimplementedFirmwareStorageServiceServer
}

func NewGrpcService() *GrpcService {
	return &GrpcService{}
}

func (g *GrpcService) InitUpload(ctx context.Context, req *storagev1alpha1.InitUploadRequest) (*storagev1alpha1.InitUploadResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (g *GrpcService) Upload(srv storagev1alpha1.FirmwareStorageService_UploadServer) error {
	// TODO implement me
	panic("implement me")
}

func (g *GrpcService) InitDownload(ctx context.Context, req *storagev1alpha1.InitDownloadRequest) (*storagev1alpha1.InitDownloadResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (g *GrpcService) Download(req *storagev1alpha1.DownloadRequest, srv storagev1alpha1.FirmwareStorageService_DownloadServer) error {
	// TODO implement me
	panic("implement me")
}
