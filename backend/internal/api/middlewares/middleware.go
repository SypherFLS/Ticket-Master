package middlewares

import (
	"context"
	"log/slog"
	"net/http"

	// "runtime/debug"
	"strings"
	"time"
	"tmaster/internal/auth"
	"tmaster/internal/constants"
	"tmaster/internal/constants/params"
	"uuid"

	"github.com/gin-gonic/gin"
)

func TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.New().String()

		ctx := context.WithValue(
			c.Request.Context(),
			constants.RequestIDKey,
			requestID,
		)

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func TimeoutMiddleware(timeout int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(timeout)*time.Second)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func LoggingMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		requestID := params.GetRequestID(c.Request.Context())
		requestLogger := logger.With(
			"request_id", requestID,
		)

		ctx := context.WithValue(c.Request.Context(), constants.RequestLogger, requestLogger)
		c.Request = c.Request.WithContext(ctx)
		c.Next()

		requestLogger.Info(
			"request completed",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"duration", time.Since(start),
		)
	}
}

func RecoverMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				requestID := params.GetRequestID(c.Request.Context())
				logger.Error(
					"panic recovered",
					"request_id", requestID,
					"error", err,
				)
				// debug.PrintStack()

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "internal server error",
				})
				return
			}
		}()

		c.Next()
	}
}

const bearer = "Bearer "

func AuthMiddleware(jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, bearer) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error" : "invalid authorization header",
			})
			return
		}

		token := strings.TrimPrefix(authHeader, bearer)

		userID, userRole, err := jwtManager.Validate(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error" : "invalid authorization header",
			})
			return
		}
		ctx := context.WithValue(c.Request.Context(), constants.UserIDKey, userID)
		ctx2 := context.WithValue(ctx, constants.UserRoleKey, userRole)
		
		c.Request = c.Request.WithContext(ctx2)
		c.Next()
	}
}
