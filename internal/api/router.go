package api

import (
	"system_components/internal/app/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(h *handler.Handler) *gin.Engine {
	r := gin.Default()

	// Статика
	r.Static("/static", "./resources")
	r.LoadHTMLGlob("templates/*")

	// Главная перенаправляет в каталог
	r.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/system-components")
	})

	// 3 обязательных GET-маршрута
	r.GET("/system-components", h.GetComponentsGrid)
	r.GET("/system-components/feed", h.GetComponentFeed)
	r.GET("/system-components/editor", h.GetComponentEditor)

	// 3 обязательных POST-маршрута
	r.POST("/system-components/draft", h.CreateDraft)      // 1. Создание черновика (ORM)
	r.POST("/system-components/publish", h.PublishComponent) // 2. Публикация черновика (ORM)
	r.POST("/system-components/delete", h.DeleteComponent)   // 3. Удаление через SQL UPDATE
	return r
}