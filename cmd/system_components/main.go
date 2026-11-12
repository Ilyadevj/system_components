package main

import (
	"context"
	"system_components/internal/app/config"
	"system_components/internal/app/dsn"
	"system_components/internal/app/handler"
	"system_components/internal/app/redis"
	"system_components/internal/app/repository"
	"system_components/pkg"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)


func main() {
	ctx := context.Background()
	cfg := config.Load()

	db, err := dsn.Connect(cfg)
	if err != nil {
		logrus.Fatalf("Database connection error: %v", err)
	}

	redisClient, err := redis.NewClient(ctx, cfg.RedisAddr, cfg.RedisPassword)
	if err != nil {
		logrus.Fatalf("Redis connection error: %v", err)
	}

	repo, err := repository.NewRepository(db, &repository.RepositorySettings{
		MinioEndpoint:   cfg.MinioEndpoint,
		MinioAccessKey:  cfg.MinioAccessKey,
		MinioSecretKey:  cfg.MinioSecretKey,
		MinioBucketName: cfg.MinioBucketName,
	})
	if err != nil {
		logrus.Fatalf("Repository initialization error: %v", err)
	}

	h := handler.NewHandler(repo, redisClient, &cfg)
	router := gin.Default()

	app := pkg.NewApp(&cfg, router, h)
	app.RunApp()
}