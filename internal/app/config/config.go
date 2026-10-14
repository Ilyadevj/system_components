package config

import "os"

type Config struct {
	DatabaseDSN     string
	MinioEndpoint   string
	MinioAccessKey  string
	MinioSecretKey  string
	MinioBucketName string
}

func Load() Config {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=system_components port=5432 sslmode=disable"
	}

	minioEndpoint := os.Getenv("MINIO_ENDPOINT")
	if minioEndpoint == "" {
		minioEndpoint = "localhost:9000"
	}

	minioAccessKey := os.Getenv("MINIO_ACCESS_KEY")
	if minioAccessKey == "" {
		minioAccessKey = "root"
	}

	minioSecretKey := os.Getenv("MINIO_SECRET_KEY")
	if minioSecretKey == "" {
		minioSecretKey = "rootpassword"
	}

	minioBucket := os.Getenv("MINIO_BUCKET_NAME")
	if minioBucket == "" {
		minioBucket = "system-components"
	}

	return Config{
		DatabaseDSN:     dsn,
		MinioEndpoint:   minioEndpoint,
		MinioAccessKey:  minioAccessKey,
		MinioSecretKey:  minioSecretKey,
		MinioBucketName: minioBucket,
	}
}