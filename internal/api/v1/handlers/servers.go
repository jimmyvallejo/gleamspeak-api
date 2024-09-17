package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
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
		InviteCode: uuid.New().String()[:12],
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

func (h *Handlers) GetServerByID(w http.ResponseWriter, r *http.Request) {
	serverID := strings.TrimPrefix(r.URL.Path, "/v1/servers/")

	serverUUID, err := uuid.Parse(serverID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to parse uuid, possible params error")
		return
	}

	server, err := h.DB.GetOneServerByID(r.Context(), serverUUID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Failed to find server")
	}

	response := SimpleServer{
		ServerID:        server.ID,
		OwnerID:         server.OwnerID,
		ServerName:      server.ServerName,
		Description:     server.Description.String,
		IconURL:         server.IconUrl.String,
		BannerURL:       server.BannerUrl.String,
		IsPublic:        server.IsPublic.Bool,
		InviteCode:      server.InviteCode,
		MemberCount:     server.MemberCount.Int32,
		ServerLevel:     server.ServerLevel.Int32,
		MaxMembers:      server.MaxMembers.Int32,
		ServerCreatedAt: server.CreatedAt,
		ServerUpdatedAt: server.UpdatedAt,
	}

	respondWithJSON(w, http.StatusOK, response)

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
			OwnerID:         server.OwnerID,
			ServerName:      server.ServerName,
			Description:     server.Description.String,
			IconURL:         server.IconUrl.String,
			BannerURL:       server.BannerUrl.String,
			IsPublic:        server.IsPublic.Bool,
			InviteCode:      server.InviteCode,
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

type JoinServerByCodeRequest struct {
	InviteCode string `json:"invite_code"`
}

func (h *Handlers) JoinServerByCode(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(common.UserContextKey).(database.User)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var request JoinServerByCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	foundServer, err := h.DB.GetOneServerByCode(r.Context(), request.InviteCode)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Server not found")
		return
	}

	userServerParams := database.CreateUserServerParams{
		UserID:   user.ID,
		ServerID: foundServer.ID,
		Role:     serverUser,
	}

	var userServer database.UserServer
	if foundServer.IsPublic.Bool {
		userServer, err = h.DB.CreateUserServer(r.Context(), userServerParams)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to join server")
			return
		}
	} else {
		respondWithError(w, http.StatusForbidden, "Server is not public")
		return
	}

	newCount := sql.NullInt32{
		Int32: foundServer.MemberCount.Int32 + 1,
		Valid: true,
	}

	updateMemberCountParams := database.UpdateServerMemberCountParams{
		ID:          foundServer.ID,
		MemberCount: newCount,
	}

	_, err = h.DB.UpdateServerMemberCount(r.Context(), updateMemberCountParams)
	if err != nil {
		log.Printf("Failed to update member count: %v", err)
	}

	respondWithJSON(w, http.StatusCreated, userServer)
}

type JoinServerByIDRequest struct {
	ServerID uuid.UUID `json:"server_id"`
}

func (h *Handlers) JoinServerByID(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(common.UserContextKey).(database.User)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var request JoinServerByIDRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	foundServer, err := h.DB.GetOneServerByID(r.Context(), request.ServerID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Server not found")
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

	newCount := sql.NullInt32{
		Int32: foundServer.MemberCount.Int32 + 1,
		Valid: true,
	}

	updateMemberCountParams := database.UpdateServerMemberCountParams{
		ID:          request.ServerID,
		MemberCount: newCount,
	}

	_, err = h.DB.UpdateServerMemberCount(r.Context(), updateMemberCountParams)
	if err != nil {
		log.Printf("Failed to update member count: %v", err)
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

func (h *Handlers) GetRecentServers(w http.ResponseWriter, r *http.Request) {

	servers, err := h.DB.GetRecentServers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to fetch recent servers: %v", err))
	}

	SimpleRecentServers := make([]SimpleRecentServer, len(servers))

	for i, server := range servers {
		SimpleRecentServers[i] = SimpleRecentServer{
			ServerID:        server.ID,
			ServerName:      server.ServerName,
			Description:     server.Description.String,
			IconURL:         server.IconUrl.String,
			BannerURL:       server.BannerUrl.String,
			MemberCount:     server.MemberCount.Int32,
			ServerCreatedAt: server.CreatedAt,
			ServerUpdatedAt: server.UpdatedAt,
			OwnerHandle:     server.Handle,
			OwnerAvatar:     server.AvatarUrl.String,
		}
	}

	respondWithJSON(w, http.StatusOK, SimpleRecentServers)
}
