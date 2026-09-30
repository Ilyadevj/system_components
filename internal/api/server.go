package api

import (
	"system_components/internal/app/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(h *handler.Handler) *gin.Engine {
	r := gin.Default()

	// Статические файлы стилей и изображений
	r.Static("/static", "./resources")

	// HTML-шаблоны
	r.LoadHTMLGlob("templates/*")

	// 3 обязательных GET-маршрута
	r.GET("/system-components", h.GetComponentsGrid)
	r.GET("/system-components/feed", h.GetComponentFeed)
	r.GET("/system-components/editor", h.GetComponentEditor)

	return r
}