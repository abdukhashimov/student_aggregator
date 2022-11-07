package services

import (
	"github.com/abdukhashimov/student_aggregator/internal/config"
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

func (s *StorageService) PushFile() {

}

func (s *StorageService) PutFile() {
	panic("implement me")
	// s.client.FPutObject()
}
