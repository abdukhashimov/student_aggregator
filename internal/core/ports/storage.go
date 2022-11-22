package ports

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
)

type StorageService interface {
	SetClient(cl *minio.Client)
	PutFile(ctx context.Context, objectName string, body io.Reader, size int64) (stug string, err error)
	GetFile(ctx context.Context, slug string) (io.Reader, int64, error)
}
