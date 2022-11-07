package ports

import "github.com/minio/minio-go/v7"

type StorageService interface {
	SetClient(cl *minio.Client)
	PutFile()
}
