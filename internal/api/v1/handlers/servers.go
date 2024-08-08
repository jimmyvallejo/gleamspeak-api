package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jimmyvallejo/gleamspeak-api/internal/api/common"
	"github.com/jimmyvallejo/gleamspeak-api/internal/database"
)

type CreateServerRequest struct {
	ServerName string `json:"server_name"`
}

func (h *Handlers) CreateServer(w http.ResponseWriter, r *http.Request) {

	user, ok := r.Context().Value(common.UserContextKey).(database.User)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unathorized")
		return
	}

	request := CreateServerRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	serverParams := database.CreateServerParams{
		ID:         uuid.New(),
		OwnerID:    user.ID,
		ServerName: request.ServerName,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	server, err := h.DB.CreateServer(r.Context(), serverParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create server")
		return
	}

	userServerParams := database.CreateUserServerParams{
		UserID:   user.ID,
		ServerID: server.ID,
	}

	_, err = h.DB.CreateUserServer(r.Context(), userServerParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create server")
		return
	}

	response := ServerResponse{
		ID: server.ID,
		OwnerID: user.ID,
		ServerName: server.ServerName,
	}

	respondWithJSON(w, http.StatusCreated, response)

}
