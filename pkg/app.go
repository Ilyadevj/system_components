package pkg

import (
	"os"
	"system_components/internal/app/config"
	"system_components/internal/app/handler"

	_ "system_components/docs"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	// CORS middleware для поддержки локальной сети, мобильных клиентов и PWA
	a.Router.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Swagger UI
	a.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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

	// ==================== 10 МЕТОДОВ REST API (ЛР3–ЛР4) ====================
	apiGroup := a.Router.Group("/api")
	{
		// Публичные методы чтения и входа (работают независимо от токена/куки)
		apiGroup.GET("/system-components", a.Handler.OptionalAuthMiddleware(), a.Handler.GetComponentsGridAPI)
		apiGroup.GET("/system-components/feed", a.Handler.OptionalAuthMiddleware(), a.Handler.GetComponentFeedAPI)
		apiGroup.POST("/users/register", a.Handler.RegisterUserAPI)
		apiGroup.POST("/users/login", a.Handler.LoginUserAPI)

		// Защищенные методы создателя (для гостя возвращают 401/403)
		authGroup := apiGroup.Group("")
		authGroup.Use(a.Handler.RequireAuthMiddleware())
		{
			authGroup.GET("/system-components/draft", a.Handler.GetDraftAPI)
			authGroup.POST("/system-components", a.Handler.CreateDraftAPI)
			authGroup.PUT("/system-components/:id", a.Handler.PublishComponentAPI)
			authGroup.DELETE("/system-components/:id", a.Handler.DeleteComponentAPI)
			authGroup.POST("/system-components/:id/like", a.Handler.LikeComponentAPI)
			authGroup.POST("/users/logout", a.Handler.LogoutUserAPI)
		}
	}

	serverAddr := os.Getenv("BACKEND_LISTEN_ADDR")
	if serverAddr == "" {
		serverAddr = "127.0.0.1:8080"
	}
	logrus.Infof("Server listening on %s", serverAddr)
	if err := a.Router.Run(serverAddr); err != nil {
		logrus.Fatalf("Failed to run server: %v", err)
	}
}
