package handler

import (
	"system_components/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	repo *repository.Repository
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{repo: repo}
}

const CurrentUserID uint = 1 // Захардкоженный ID пользователя до авторизации

func (h *Handler) ErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	c.HTML(statusCode, "system_components_grid.html", gin.H{
		"Error": message,
	})
}