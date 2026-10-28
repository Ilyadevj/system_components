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

// @title System Components API
// @version 1.0
// @description Веб-сервис расчета времени отклика веб-приложения (system_components)
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Введите токен в формате: Bearer <ваш_jwt_токен>
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