package services

import (
	"context"
	"io"

	"github.com/abdukhashimov/student_aggregator/internal/config"
	"github.com/abdukhashimov/student_aggregator/internal/core/domain"
	"github.com/abdukhashimov/student_aggregator/internal/core/ports"
	"github.com/minio/minio-go/v7"
)

type StorageService struct {
	cfg    *config.Config
	client *minio.Client
}

func NewStorageService(cfg *config.Config) ports.StorageService {
	return &StorageService{
		cfg: cfg,
	}
}

func (s *StorageService) SetClient(cl *minio.Client) {
	s.client = cl
}

func (s *StorageService) PutFile(ctx context.Context, options domain.PutFileOptions) (string, error) {

	info, err := s.client.PutObject(
		ctx, s.cfg.Storage.BucketName,
		options.ObjectName,
		options.Body,
		options.Size,
		minio.PutObjectOptions{
			ContentType: options.ContentType,
		},
	)
	if err != nil {
		return "", err
	}

	return info.Key, nil
}

func (s *StorageService) GetFile(ctx context.Context, objectName string) (content io.Reader, size int64, err error) {

	content, err = s.client.GetObject(ctx, s.cfg.Storage.BucketName, objectName, minio.GetObjectOptions{})

	if err != nil {
		return
	}

	objectInfo, err := content.(*minio.Object).Stat()

	if err != nil {
		return
	}

	return content, objectInfo.Size, err
}
