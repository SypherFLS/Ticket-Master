package middlewares

import (
	"context"
	"log"
	"net/http"
	"runtime/debug"
	"time"
	"tmaster/internal/api/helpers"
)

type Middleware func(http.Handler) http.Handler

func Chain(h http.Handler, m ...Middleware) http.Handler {
	for i := len(m) - 1; i >= 0; i-- {
		h = m[i](h)
	}

	return h
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
