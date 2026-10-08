package router

import (
	"log/slog"
	"tmaster/internal/api/handlers"
	"tmaster/internal/api/middlewares"

	"github.com/gin-gonic/gin"

	"tmaster/internal/auth"
	"tmaster/internal/config"
)

func NewRouter(h *handlers.Handler, jwtManager *auth.JWTManager, cfg config.Config, logger *slog.Logger) *gin.Engine {
	router := gin.New()

	router.Use(
		middlewares.TraceMiddleware(),
		middlewares.LoggingMiddleware(logger),
		middlewares.RecoverMiddleware(logger),
		middlewares.TimeoutMiddleware(cfg.Server.Timeout),
	)

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
