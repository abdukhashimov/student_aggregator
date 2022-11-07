package minio

import (
	"context"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const timeout = 10 * time.Second

func NewClient(uri, accessKeyId, secretKey string) *minio.Client {
	storageClient, err := minio.New(uri, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyId, secretKey, ""),
		Secure: false,
	})

	if err != nil {
		panic(err)
	}

	listContext, cancelList := context.WithTimeout(context.Background(), timeout)
	defer cancelList()

	_, err = storageClient.ListBuckets(listContext)
	if err != nil {
		panic(err)
	}

	return storageClient
}
