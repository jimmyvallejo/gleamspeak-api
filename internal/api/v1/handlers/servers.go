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
		Role:     serverAdmin,
	}

	_, err = h.DB.CreateUserServer(r.Context(), userServerParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to join server")
		return
	}

	response := CreateServerResponse{
		ID:         server.ID,
		OwnerID:    user.ID,
		ServerName: server.ServerName,
	}

	respondWithJSON(w, http.StatusCreated, response)

}

func (h *Handlers) GetUserServers(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(common.UserContextKey).(database.User)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	servers, err := h.DB.GetUserServers(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch servers")
		return
	}

	simpleServers := make([]SimpleServer, len(servers))
	for i, server := range servers {
		simpleServers[i] = SimpleServer{
			ServerID:        server.ServerID,
			ServerName:      server.ServerName,
			Description:     server.Description.String,
			IconURL:         server.IconUrl.String,
			BannerURL:       server.BannerUrl.String,
			IsPublic:        server.IsPublic.Bool,
			MemberCount:     server.MemberCount.Int32,
			ServerLevel:     server.ServerLevel.Int32,
			MaxMembers:      server.MaxMembers.Int32,
			ServerCreatedAt: server.ServerCreatedAt,
			ServerUpdatedAt: server.ServerUpdatedAt,
		}
	}

	response := SimpleDisplayServerResponse{
		UserID:  user.ID,
		Servers: simpleServers,
	}

	respondWithJSON(w, http.StatusOK, response)
}

type JoinServerRequest struct {
	ServerID uuid.UUID `json:"server_id"`
}

func (h *Handlers) JoinServer(w http.ResponseWriter, r *http.Request) {

	user, ok := r.Context().Value(common.UserContextKey).(database.User)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unathorized")
		return
	}

	request := JoinServerRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	userServerParams := database.CreateUserServerParams{
		UserID:   user.ID,
		ServerID: request.ServerID,
		Role:     serverUser,
	}

	userServer, err := h.DB.CreateUserServer(r.Context(), userServerParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to join server")
		return
	}

	respondWithJSON(w, http.StatusCreated, userServer)
}

type leaveServerRequest struct {
	ServerID uuid.UUID `json:"server_id"`
}

func (h *Handlers) LeaveServer(w http.ResponseWriter, r *http.Request) {

	user, ok := r.Context().Value(common.UserContextKey).(database.User)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unathorized")
		return
	}

	request := leaveServerRequest{}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	userServerParams := database.DeleteUserServerParams{
		UserID:   user.ID,
		ServerID: request.ServerID,
	}

	err = h.DB.DeleteUserServer(r.Context(), userServerParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to leave server")
		return
	}

	respondNoBody(w, http.StatusOK)
}
