package router

import (
	"log/slog"
	"tmaster/internal/api/handlers"
	"tmaster/internal/api/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"tmaster/docs"
	"tmaster/internal/auth"
	"tmaster/internal/config"
)

func NewRouter(h *handlers.Handler, jwtManager *auth.JWTManager, cfg config.Config, logger *slog.Logger) *gin.Engine {
	router := gin.New()

	docs.SwaggerInfo.Title = "Ticket Master API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Host = "localhost:8080"

	router.Use(
		middlewares.TraceMiddleware(),
		middlewares.LoggingMiddleware(logger),
		middlewares.RecoverMiddleware(logger),
		middlewares.TimeoutMiddleware(cfg.Server.Timeout),
	)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", h.RegisterHandler)
		authGroup.POST("/login", h.LoginHandler)
	}

	apiGroup := router.Group("/api")
	apiGroup.Use(
		middlewares.AuthMiddleware(jwtManager),
	)
	{
		apiGroup.POST("/create", h.CreateTicketHandler)
		apiGroup.GET("/tickets", h.GetOwnTicketsHandler)
		apiGroup.POST("/claim", h.ClaimNextTicketHandler)
		apiGroup.PATCH("/close", h.CloseTicketHandler)
		apiGroup.GET("/tickets_queue", h.GetNewTicketsHandler)
	}

	return router
}
