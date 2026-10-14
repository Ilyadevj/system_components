package repository

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gorm.io/gorm"
)

type RepositorySettings struct {
	MinioEndpoint   string
	MinioAccessKey  string
	MinioSecretKey  string
	MinioBucketName string
}

type Repository struct {
	db              *gorm.DB
	minio           *minio.Client
	minioBucketName string
	minioEndpoint   string
}

func NewRepository(db *gorm.DB, settings *RepositorySettings) (*Repository, error) {
	minioClient, err := minio.New(settings.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(settings.MinioAccessKey, settings.MinioSecretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	exists, err := minioClient.BucketExists(ctx, settings.MinioBucketName)
	if err == nil && !exists {
		_ = minioClient.MakeBucket(ctx, settings.MinioBucketName, minio.MakeBucketOptions{})
		policy := fmt.Sprintf(`{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::%s/*"]}]}`, settings.MinioBucketName)
		_ = minioClient.SetBucketPolicy(ctx, settings.MinioBucketName, policy)
	}

	return &Repository{
		db:              db,
		minio:           minioClient,
		minioBucketName: settings.MinioBucketName,
		minioEndpoint:   settings.MinioEndpoint,
	}, nil
}