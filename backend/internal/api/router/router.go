package router

import (
	"net/http"
	"tmaster/internal/api/handlers"
	"tmaster/internal/auth"
)

func NewRouter(h *handlers.Handler, jwtm *auth.JWTManager) http.Handler{
	mux := http.NewServeMux()

	return mux
}