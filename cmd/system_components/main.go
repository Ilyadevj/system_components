package main

import (
	"system_components/internal/app/config"
	"system_components/internal/app/dsn"
	"system_components/internal/app/handler"
	"system_components/internal/app/repository"
	"system_components/pkg"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg := config.Load()

	db, err := dsn.Connect(cfg)
	if err != nil {
		logrus.Fatalf("Database connection error: %v", err)
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

	h := handler.NewHandler(repo)
	router := gin.Default()

	app := pkg.NewApp(&cfg, router, h)
	app.RunApp()
}