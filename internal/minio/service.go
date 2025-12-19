package minio

import (
	"bytes"
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)

func Upsert(filename string, content []byte) {
	ctx := context.Background()

	minioClient, err := minio.New(Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(AccessKeyID, SecretAccessKey, Token),
		Secure: UseSSL,
	})
	if err != nil {
		log.Fatal(err)
	}

	bucketInit(minioClient, ctx)

	reader := bytes.NewReader(content)
	_, err = minioClient.PutObject(ctx, BucketName, filename, reader, int64(len(content)), minio.PutObjectOptions{})
	if err != nil {
		log.Fatal(err)
	}
}

func bucketInit(minioClient *minio.Client, ctx context.Context) {
	exists, err := minioClient.BucketExists(ctx, BucketName)
	if err != nil {
		log.Fatal(err)
	}
	if !exists {
		err := minioClient.MakeBucket(ctx, BucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatal(err)
		}
	}
}
