package ports

import (
	"context"
	"io"

	"github.com/abdukhashimov/student_aggregator/internal/core/domain"
	"github.com/minio/minio-go/v7"
)

type StorageService interface {
	SetClient(cl *minio.Client)
	PutFile(ctx context.Context, options domain.PutFileOptions) (stug string, err error)
	GetFile(ctx context.Context, slug string) (io.Reader, int64, error)
}
