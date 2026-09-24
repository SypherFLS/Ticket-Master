package router

import (
	"net/http"
	"tmaster/internal/api/handlers"
	"tmaster/internal/api/middlewares"

	"tmaster/internal/auth"
	"tmaster/internal/config"
)

func NewRouter(h *handlers.Handler, jwtManager *auth.JWTManager, cfg config.Config) http.Handler {
	root := http.NewServeMux()

	public := http.NewServeMux()
	private := http.NewServeMux()

	public.Handle("POST /register", http.HandlerFunc(h.RegisterHandler))
	public.Handle("POST /login", http.HandlerFunc(h.LoginHandler))

	private.Handle("POST /create", http.HandlerFunc(h.CreateTicketHandler))
	private.Handle("GET /tickets", http.HandlerFunc(h.GetOwnTicketsHandler))
	private.Handle("POST /claime", http.HandlerFunc(h.ClaimNextTicketHandler))
	private.Handle("PATCH /close", http.HandlerFunc(h.CloseTicketHandler))
	private.Handle("GET /tickets_queue", http.HandlerFunc(h.GetNewTicketsHandler))

	publicChain := middlewares.CommonChain(
		public,
		cfg.Server.Timeout,
	)
	privateChain := middlewares.CommonChain(
		middlewares.AuthMiddleware(jwtManager)(
			private,
		),
		cfg.Server.Timeout,
	)

	root.Handle("/api/", http.StripPrefix("/api", privateChain))
	root.Handle("/auth/", http.StripPrefix("/auth", publicChain))
	return root
}
