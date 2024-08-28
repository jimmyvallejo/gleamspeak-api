package routes

import (
	"net/http"

	"github.com/jimmyvallejo/gleamspeak-api/internal/api/middleware"
	"github.com/jimmyvallejo/gleamspeak-api/internal/api/v1/handlers"
	"github.com/jimmyvallejo/gleamspeak-api/internal/websocket"
)

type Router struct {
	mux        *http.ServeMux
	handlers   *handlers.Handlers
	middleware *middleware.Middleware
	websocket  *websocket.Manager
}

func NewRouter(h *handlers.Handlers, m *middleware.Middleware, w *websocket.Manager) *Router {
	return &Router{
		mux:        http.NewServeMux(),
		handlers:   h,
		middleware: m,
		websocket:  w,
	}
}

func (r *Router) SetupRoutes() {
	// Test readiness
	r.mux.HandleFunc("GET /v1/healthz", handlers.HandlerReadiness)
	r.mux.HandleFunc("GET /v1/err", handlers.HandlerError)

	// Auth Routes
	r.mux.HandleFunc("POST /v1/login", r.handlers.LoginUserStandard)
	r.mux.HandleFunc("POST /v1/logout", r.handlers.LogoutUserStandard)
	r.mux.HandleFunc("GET /v1/auth", r.middleware.IsAuthenticated(r.handlers.CheckAuthStatus))

	// User Routes
	r.mux.HandleFunc("POST /v1/users", r.handlers.CreateUserStandard)
	r.mux.HandleFunc("PUT /v1/users", r.middleware.IsAuthenticated(r.handlers.UpdateUser))

	// Server Routes
	r.mux.HandleFunc("POST /v1/servers", r.middleware.IsAuthenticated(r.handlers.CreateServer))
	r.mux.HandleFunc("POST /v1/servers/join", r.middleware.IsAuthenticated(r.handlers.JoinServerByID))
	r.mux.HandleFunc("POST /v1/servers/code", r.middleware.IsAuthenticated(r.handlers.JoinServerByCode))
	r.mux.HandleFunc("DELETE /v1/servers/user", r.middleware.IsAuthenticated(r.handlers.LeaveServer))
	r.mux.HandleFunc("GET /v1/servers/recent", r.handlers.GetRecentServers)
	r.mux.HandleFunc("GET /v1/servers/user/many", r.middleware.IsAuthenticated(r.handlers.GetUserServers))

	// Text Channel Routes
	r.mux.HandleFunc("POST /v1/channels/text", r.middleware.IsAuthenticated(r.handlers.CreateTextChannel))
	r.mux.HandleFunc("GET /v1/channels/{serverID}", r.middleware.IsAuthenticated(r.handlers.GetServerTextChannels))

	// Message Routes
	r.mux.HandleFunc("GET /v1/messages/{channelID}", r.middleware.IsAuthenticated(r.handlers.GetChannelTextMessages))

	// Token Routes
	r.mux.HandleFunc("POST /v1/refresh", r.handlers.RefreshToken)

	// Upgrade to WebSocket
	r.mux.HandleFunc("/ws", r.websocket.ServeWs)
}

func (r *Router) GetHandler() http.Handler {
	return r.mux
}
