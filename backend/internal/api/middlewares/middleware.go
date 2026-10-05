package middlewares

import (
	"context"
	"log/slog"
	"net/http"

	// "runtime/debug"
	"strings"
	"time"
	"tmaster/internal/api/utils/helpers"
	"tmaster/internal/api/utils/selfwriter"
	"tmaster/internal/auth"
	"tmaster/internal/constants"
	"tmaster/internal/constants/params"
	"uuid"
)

type Middleware func(http.Handler) http.Handler

func Chain(h http.Handler, m ...Middleware) http.Handler {
	for i := len(m) - 1; i >= 0; i-- {
		h = m[i](h)
	}

	return h
}

func CommonChain(h http.Handler, timeout int, logger *slog.Logger) http.Handler {
	return Chain(
		h,
		TraceMiddleware,
		LoggingMiddleware(logger),
		RecoverMiddleware(logger),
		TimeoutMiddleware(timeout),
	)
}

func TraceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()
		ctx := context.WithValue(
			r.Context(),
			constants.RequestIDKey,
			requestID,
		)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func TimeoutMiddleware(timeout int) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), time.Duration(timeout)*time.Second)
			defer cancel()

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func LoggingMiddleware(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			sw := &selfwriter.SelfWriter{
				ResponseWriter: w,
				Code:           200,
			}
			requestID := params.GetRequestID(r.Context())
			requestLogger := logger.With(
				"request_id", requestID,
			)

			ctx := context.WithValue(r.Context(), constants.RequestLogger, requestLogger)
			next.ServeHTTP(sw, r.WithContext(ctx))

			requestLogger.Info(
				"request completed",
				"method", r.Method,
				"path", r.URL.Path,
				"status", sw.Code,
				"duration", time.Since(start),
			)
		})
	}
}

func RecoverMiddleware(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			defer func() {
				if err := recover(); err != nil {
					requestID := params.GetRequestID(r.Context())
					logger.Error(
						"panic recovered",
						"request_id", requestID,
						"error", err,
					)
					// debug.PrintStack()

					helpers.WriteError(w, 500, "panic")
					return
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

const bearer = "Bearer "

func AuthMiddleware(jwtManager *auth.JWTManager) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				helpers.WriteError(w, http.StatusUnauthorized, "invalid authorization header")
				return
			}
			if !strings.HasPrefix(authHeader, bearer) {
				helpers.WriteError(w, http.StatusUnauthorized, "invalid authorization header")
				return
			}

			token := strings.TrimPrefix(authHeader, bearer)

			userID, userRole, err := jwtManager.Validate(token)
			if err != nil {
				helpers.WriteError(w, http.StatusUnauthorized, "invalid authorization header")
				return
			}
			ctx := context.WithValue(r.Context(), constants.UserIDKey, userID)
			ctx2 := context.WithValue(ctx, constants.UserRoleKey, userRole)
			next.ServeHTTP(w, r.WithContext(ctx2))
		})
	}
}
