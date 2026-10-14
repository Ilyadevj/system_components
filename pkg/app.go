package pkg

import (
	"system_components/internal/app/config"
	"system_components/internal/app/handler"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Application struct {
	Config  *config.Config
	Router  *gin.Engine
	Handler *handler.Handler
}

func NewApp(c *config.Config, r *gin.Engine, h *handler.Handler) *Application {
	return &Application{
		Config:  c,
		Router:  r,
		Handler: h,
	}
}

func (a *Application) RunApp() {
	logrus.Info("Server is starting...")

	// Статические файлы и шаблоны (ЛР1–ЛР2)
	a.Router.Static("/static", "./resources")
	a.Router.LoadHTMLGlob("templates/*")

	a.Router.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/system-components")
	})

	// HTML-маршруты (ЛР2)
	a.Router.GET("/system-components", a.Handler.GetComponentsGrid)
	a.Router.GET("/system-components/feed", a.Handler.GetComponentFeed)
	a.Router.GET("/system-components/editor", a.Handler.GetComponentEditor)
	a.Router.POST("/system-components/draft", a.Handler.CreateDraft)
	a.Router.POST("/system-components/publish", a.Handler.PublishComponent)
	a.Router.POST("/system-components/delete", a.Handler.DeleteComponent)

	// ==================== 10 МЕТОДОВ REST API (ЛР3) ====================
	apiGroup := a.Router.Group("/api")
	{
		// Домен услуги (system-components) — 7 методов:
		apiGroup.GET("/system-components", a.Handler.GetComponentsGridAPI)          // 1. Список с фильтрацией
		apiGroup.GET("/system-components/feed", a.Handler.GetComponentFeedAPI)      // 2. Лента (1 карточка)
		apiGroup.GET("/system-components/draft", a.Handler.GetDraftAPI)             // 3. Черновик пользователя
		apiGroup.POST("/system-components", a.Handler.CreateDraftAPI)               // 4. Добавление черновика + файлы в MinIO
		apiGroup.PUT("/system-components/:id", a.Handler.PublishComponentAPI)       // 5. Публикация черновика
		apiGroup.DELETE("/system-components/:id", a.Handler.DeleteComponentAPI)     // 6. Мягкое удаление (soft delete)
		apiGroup.POST("/system-components/:id/like", a.Handler.LikeComponentAPI)    // 7. Постановка/отмена лайка (0/1)

		// Домен пользователя (users) — 3 метода:
		apiGroup.POST("/users/register", a.Handler.RegisterUserAPI)                 // 8. Регистрация
		apiGroup.POST("/users/login", a.Handler.LoginUserAPI)                       // 9. Аутентификация (заглушка)
		apiGroup.POST("/users/logout", a.Handler.LogoutUserAPI)                     // 10. Деавторизация (заглушка)
	}

	serverAddr := ":8080"
	logrus.Infof("Server listening on http://localhost%s", serverAddr)
	if err := a.Router.Run(serverAddr); err != nil {
		logrus.Fatalf("Failed to run server: %v", err)
	}
}