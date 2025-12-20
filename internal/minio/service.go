package minio

import (
	"bytes"
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"os"
)

func Upsert(filename string, content []byte) {
	ctx := context.Background()
	minioClient := getMinIOClient()

	bucketInit(minioClient, ctx)

	reader := bytes.NewReader(content)
	_, err := minioClient.PutObject(ctx, BucketName, filename, reader, int64(len(content)), minio.PutObjectOptions{})
	if err != nil {
		log.Fatal(err)
	}
}

func getMinIOClient() *minio.Client {
	minioClient, err := minio.New(os.Getenv("minio_address"), &minio.Options{
		Creds: credentials.NewStaticV4(
			os.Getenv("minio_access_key_id"),
			os.Getenv("minio_secret_access_key"),
			os.Getenv("minio_token"),
		),
		Secure: UseSSL,
	})
	if err != nil {
		log.Fatal(err)
	}
	return minioClient
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
