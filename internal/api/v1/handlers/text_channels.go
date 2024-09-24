package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jimmyvallejo/gleamspeak-api/internal/api/common"
	"github.com/jimmyvallejo/gleamspeak-api/internal/database"
)

type CreateTextChannelRequest struct {
	ServerID    string `json:"server_id"`
	Language    string `json:"language"`
	ChannelName string `json:"channel_name"`
}

func (h *Handlers) CreateTextChannel(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(common.UserContextKey).(database.User)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unathorized")
		return
	}

	request := CreateTextChannelRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	serverStr := request.ServerID

	serverUUID, err := uuid.Parse(serverStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	languageID, err := h.DB.GetLanguageIDByName(r.Context(), request.Language)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	channelParams := database.CreateTextChannelParams{
		ID:          uuid.New(),
		OwnerID:     user.ID,
		ServerID:    serverUUID,
		LanguageID:  languageID,
		ChannelName: request.ChannelName,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	channel, err := h.DB.CreateTextChannel(r.Context(), channelParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create channel")
		return
	}

	response := CreateTextChannelResponse{
		ID:          channel.ID,
		OwnerID:     channel.OwnerID,
		ServerID:    channel.ServerID,
		ChannelName: channel.ChannelName,
	}

	respondWithJSON(w, http.StatusCreated, response)

}

func (h *Handlers) GetServerTextChannels(w http.ResponseWriter, r *http.Request) {
	serverID := strings.TrimPrefix(r.URL.Path, "/v1/channels/")

	serverUUID, err := uuid.Parse(serverID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}
	channels, err := h.DB.GetServerTextChannels(r.Context(), serverUUID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get servers belonging to channel")
	}

	simpleChannels := make([]SimpleChannel, len(channels))

	for i, channel := range channels {
		simpleChannels[i] = SimpleChannel{
			ChannelID:        channel.ID,
			OwnerID:          channel.OwnerID,
			ServerID:         channel.ServerID,
			LanguageID:       channel.LanguageID,
			ChannelName:      channel.ChannelName,
			LastActive:       channel.LastActive.Time,
			IsLocked:         channel.IsLocked.Bool,
			ChannelCreatedAt: channel.CreatedAt,
			ChannelUpdatedAt: channel.UpdatedAt,
		}
	}

	response := GetServerTextChannelResponse{
		ServerID: serverUUID,
		Channels: simpleChannels,
	}
	respondWithJSON(w, http.StatusOK, response)
}

func (h *Handlers) GetChannelTextMessages(w http.ResponseWriter, r *http.Request) {
	channnelID := strings.TrimPrefix(r.URL.Path, "/v1/messages/")

	channelUUID, err := uuid.Parse(channnelID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to parse uuid, possible params")
		return
	}

	messages, err := h.DB.GetChannelTextMessages(r.Context(), channelUUID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("%v:", err))
	}

	normalizedMessages := make([]SimpleMessage, len(messages))

	for i, message := range messages {
		normalizedMessages[i] = SimpleMessage{
			ID:          message.ID,
			ChannelID:   message.ChannelID,
			OwnerID:     message.OwnerID,
			OwnerHandle: message.Handle,
			OwnerImage:  message.AvatarUrl.String,
			Message:     message.Message,
			Image:       message.Image.String,
			CreatedAt:   message.CreatedAt,
			UpdatedAt:   message.UpdatedAt,
		}

	}

	respondWithJSON(w, http.StatusOK, normalizedMessages)
}
