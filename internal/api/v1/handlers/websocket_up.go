package handlers

import (
	"net/http"
)

func (h *Handlers) HandleWebSocketUpgrade(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("userId")
	handle := r.URL.Query().Get("handle")
	if userId == "" || handle == "" {
		http.Error(w, "Unauthorized: User ID is required", http.StatusUnauthorized)
		return
	}

	h.Ws.ServeWs(w, r, userId, handle)
}
