package minio

import (
	"context"
	"time"

	"github.com/abdukhashimov/student_aggregator/internal/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const timeout = 10 * time.Second

func NewClient(conf config.StorageConfig) *minio.Client {
	storageClient, err := minio.New(conf.URI, &minio.Options{
		Creds:  credentials.NewStaticV4(conf.AccessKeyID, conf.SecretAccessKey, ""),
		Secure: false,
	})

	if err != nil {
		panic(err)
	}

	bucketCtx, cancelCtx := context.WithTimeout(context.Background(), timeout)
	defer cancelCtx()

	exists, errBucketExists := storageClient.BucketExists(bucketCtx, conf.BucketName)
	if errBucketExists != nil {
		panic(err)
	}

	if !exists {
		err := storageClient.MakeBucket(bucketCtx, conf.BucketName, minio.MakeBucketOptions{})

		if err != nil {
			panic(err)
		}
	}

	return storageClient
}
