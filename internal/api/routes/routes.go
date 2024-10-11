package routes

import (
	"net/http"

	"github.com/jimmyvallejo/gleamspeak-api/internal/api/middleware"
	"github.com/jimmyvallejo/gleamspeak-api/internal/api/v1/handlers"
)

type Router struct {
	mux        *http.ServeMux
	handlers   *handlers.Handlers
	middleware *middleware.Middleware
}

func NewRouter(h *handlers.Handlers, m *middleware.Middleware) *Router {
	return &Router{
		mux:        http.NewServeMux(),
		handlers:   h,
		middleware: m,
	}
}

func (r *Router) SetupV1Routes() {

	// Test readiness
	r.mux.HandleFunc("GET /v1/healthz", handlers.HandlerReadiness)
	r.mux.HandleFunc("GET /v1/err", handlers.HandlerError)

	// AWS Routes
	r.mux.HandleFunc("POST /v1/s3/url", r.middleware.IsAuthenticated(r.handlers.GetSignedURL))

	// Auth Routes
	r.mux.HandleFunc("POST /v1/login", r.handlers.LoginUserStandard)
	r.mux.HandleFunc("POST /v1/logout", r.handlers.LogoutUserStandard)
	r.mux.HandleFunc("GET /v1/auth", r.middleware.IsAuthenticated(r.handlers.CheckAuthStatus))

	// User Routes
	r.mux.HandleFunc("POST /v1/users", r.handlers.CreateUserStandard)
	r.mux.HandleFunc("PUT /v1/users", r.middleware.IsAuthenticated(r.handlers.UpdateUser))
	r.mux.HandleFunc("PUT /v1/users/avatar", r.middleware.IsAuthenticated(r.handlers.UpdateAvatar))
	r.mux.HandleFunc("GET /v1/users/auth", r.middleware.IsAuthenticated(r.handlers.FetchAuthUser))
	r.mux.HandleFunc("DELETE /v1/users/{userID}", r.middleware.IsAuthenticated(r.handlers.DeleteUser))

	// Server Routes
	r.mux.HandleFunc("POST /v1/servers", r.middleware.IsAuthenticated(r.handlers.CreateServer))
	r.mux.HandleFunc("POST /v1/servers/join", r.middleware.IsAuthenticated(r.handlers.JoinServerByID))
	r.mux.HandleFunc("POST /v1/servers/code", r.middleware.IsAuthenticated(r.handlers.JoinServerByCode))
	r.mux.HandleFunc("PUT /v1/servers", r.middleware.IsAuthenticated(r.handlers.UpdateServer))
	r.mux.HandleFunc("PUT /v1/servers/images", r.middleware.IsAuthenticated(r.handlers.UpdateServerImages))
	r.mux.HandleFunc("DELETE /v1/servers/user", r.middleware.IsAuthenticated(r.handlers.LeaveServer))
	r.mux.HandleFunc("GET /v1/servers/recent", r.handlers.GetRecentServers)
	r.mux.HandleFunc("GET /v1/servers/user/many", r.middleware.IsAuthenticated(r.handlers.GetUserServers))
	r.mux.HandleFunc("GET /v1/servers/{serverID}", r.handlers.GetServerByID)
	r.mux.HandleFunc("DELETE /v1/servers/{serverID}", r.middleware.IsAuthenticated(r.handlers.DeleteServer))

	// Text Channel Routes
	r.mux.HandleFunc("POST /v1/channels/text", r.middleware.IsAuthenticated(r.handlers.CreateTextChannel))
	r.mux.HandleFunc("GET /v1/channels/{serverID}", r.middleware.IsAuthenticated(r.handlers.GetServerTextChannels))

	// Voice Channel Routes
	r.mux.HandleFunc("POST /v1/channels/voice", r.middleware.IsAuthenticated(r.handlers.CreateVoiceChannel))
	r.mux.HandleFunc("GET /v1/channels/voice/{serverID}", r.middleware.IsAuthenticated(r.handlers.GetServerVoiceChannels))
	r.mux.HandleFunc("DELETE /v1/channels/voice/{userID}", r.middleware.IsAuthenticated(r.handlers.LeaveVoiceChannelByUserID))

	// Message Routes
	r.mux.HandleFunc("GET /v1/messages/{channelID}", r.middleware.IsAuthenticated(r.handlers.GetChannelTextMessages))

	// Token Routes
	r.mux.HandleFunc("POST /v1/refresh", r.handlers.RefreshToken)

	// Upgrade to WebSocket
	r.mux.HandleFunc("GET /ws", r.handlers.HandleWebSocketUpgrade)
}

func (r *Router) GetHandler() http.Handler {
	return r.mux
}
