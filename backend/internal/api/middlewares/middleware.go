package middlewares

import (
	"context"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
	"time"
	"tmaster/internal/api/utils/helpers"
	"tmaster/internal/api/utils/params"
	"tmaster/internal/api/utils/selfwriter"
	"tmaster/internal/auth"
)

type Middleware func(http.Handler) http.Handler

func Chain(h http.Handler, m ...Middleware) http.Handler {
	for i := len(m) - 1; i >= 0; i-- {
		h = m[i](h)
	}

	return h
}

func CommonChain(h http.Handler, timeout int) http.Handler {
	return Chain(
		h,
		LoggingMiddleware,
		RecoverMiddleware,
		TimeoutMiddleware(timeout),
	)
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

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		current_time := time.Now()
		sw := &selfwriter.SelfWriter{
			ResponseWriter: w,
			Code: 200,
		}
		log.Printf("handler %v started \n", r.URL.String())
		next.ServeHTTP(sw, r)
		log.Printf("handler %v finished with code %v and time %v \n", r.URL.Path, sw.Code, time.Since(current_time))
	})
}

func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := recover(); err != nil {
				log.Printf("PANIC: %v\n", err)
				debug.PrintStack()

				helpers.WriteError(w, 500, "panic")
				return
			}
		}()

		next.ServeHTTP(w, r)
	})
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
			ctx := context.WithValue(r.Context(), params.UserIDKey, userID)
			ctx2 := context.WithValue(ctx, params.UserRoleKey, userRole)
			next.ServeHTTP(w, r.WithContext(ctx2))
		})
	}
}