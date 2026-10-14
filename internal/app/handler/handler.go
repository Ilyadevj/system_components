package handler

import (
	"fmt"
	"mime/multipart"
	"net/http"
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

const CurrentUserID uint = 1

func GetCurrentUserID() uint {
	return CurrentUserID
}

func (h *Handler) ErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	c.HTML(statusCode, "system_components_grid.html", gin.H{
		"Error": message,
	})
}

func (h *Handler) errorHandlerAPI(c *gin.Context, statusCode int, err error) {
	logrus.Error(err.Error())
	c.JSON(statusCode, gin.H{
		"error": err.Error(),
	})
}

func isAllowedType(contentType string, allowed []string) bool {
	for _, t := range allowed {
		if contentType == t {
			return true
		}
	}
	return false
}

func validateMediaUpload(header *multipart.FileHeader, allowedTypes []string) (int, error) {
	file, err := header.Open()
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("не удалось открыть файл")
	}
	defer file.Close()

	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("не удалось прочитать файл")
	}

	contentType := http.DetectContentType(buffer[:n])
	if !isAllowedType(contentType, allowedTypes) {
		return http.StatusBadRequest, fmt.Errorf("недопустимый формат файла: %s", contentType)
	}

	if _, err := file.Seek(0, 0); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("ошибка чтения потока файла")
	}

	return http.StatusOK, nil
}