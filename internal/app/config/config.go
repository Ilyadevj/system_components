package config

import (
	"os"
	"time"
)

type Config struct {
	DatabaseDSN     string
	MinioEndpoint   string
	MinioAccessKey  string
	MinioSecretKey  string
	MinioBucketName string
	RedisAddr       string
	RedisPassword   string
	JWTSecret       string
	JWTExpiresIn    time.Duration
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

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "super_secret_jwt_key_system_components"
	}

	return Config{
		DatabaseDSN:     dsn,
		MinioEndpoint:   minioEndpoint,
		MinioAccessKey:  minioAccessKey,
		MinioSecretKey:  minioSecretKey,
		MinioBucketName: minioBucket,
		RedisAddr:       redisAddr,
		RedisPassword:   os.Getenv("REDIS_PASSWORD"),
		JWTSecret:       jwtSecret,
		JWTExpiresIn:    24 * time.Hour,
	}
}