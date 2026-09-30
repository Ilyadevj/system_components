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

	// 1. Статические файлы и шаблоны
	a.Router.Static("/static", "./resources")
	a.Router.LoadHTMLGlob("templates/*")

	// 2. Редирект с корня в каталог
	a.Router.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/system-components")
	})

	// 3. Регистрация обязательных 3 GET-маршрутов
	a.Router.GET("/system-components", a.Handler.GetComponentsGrid)
	a.Router.GET("/system-components/feed", a.Handler.GetComponentFeed)
	a.Router.GET("/system-components/editor", a.Handler.GetComponentEditor)

	// 4. Регистрация обязательных 3 POST-маршрутов
	a.Router.POST("/system-components/draft", a.Handler.CreateDraft)
	a.Router.POST("/system-components/publish", a.Handler.PublishComponent)
	a.Router.POST("/system-components/delete", a.Handler.DeleteComponent)

	// 5. Запуск HTTP-сервера
	serverAddr := ":8080"
	logrus.Infof("Server listening on http://localhost%s", serverAddr)
	if err := a.Router.Run(serverAddr); err != nil {
		logrus.Fatalf("Failed to run server: %v", err)
	}
}
