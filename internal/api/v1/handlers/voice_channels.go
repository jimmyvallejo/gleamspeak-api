package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jimmyvallejo/gleamspeak-api/internal/api/common"
	"github.com/jimmyvallejo/gleamspeak-api/internal/database"
)

type CreateVoiceChannelRequest struct {
	ServerID    string `json:"server_id"`
	Language    string `json:"language"`
	ChannelName string `json:"channel_name"`
}

func (h *Handlers) CreateVoiceChannel(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(common.UserContextKey).(database.User)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	request := CreateVoiceChannelRequest{}
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

	channelParams := database.CreateVoiceChannelParams{
		ID:          uuid.New(),
		OwnerID:     user.ID,
		ServerID:    serverUUID,
		LanguageID:  languageID,
		ChannelName: request.ChannelName,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	channel, err := h.DB.CreateVoiceChannel(r.Context(), channelParams)
	if err != nil {
		log.Printf("Error creating voice channel: %v", err)
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

func (h *Handlers) LeaveVoiceChannelByUserID(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(common.UserContextKey).(database.User)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unathorized")
		return
	}

	userID := strings.TrimPrefix(r.URL.Path, "/v1/channels/voice/")

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to parse uuid, possible params error")
		return
	}

	err = h.DB.LeaveVoiceChannelByUser(r.Context(), userUUID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Failed to find channel")
		return
	}

	respondNoBody(w, http.StatusOK)
}

func (h *Handlers) GetServerVoiceChannels(w http.ResponseWriter, r *http.Request) {
	serverID := strings.TrimPrefix(r.URL.Path, "/v1/channels/voice/")
	serverUUID, err := uuid.Parse(serverID)
	if err != nil {
		log.Printf("Invalid server ID: %v", err)
		respondWithError(w, http.StatusBadRequest, "Invalid server ID")
		return
	}

	channels, err := h.DB.GetServerVoiceChannels(r.Context(), serverUUID)
	if err != nil {
		log.Printf("Failed to get voice channels for server %s: %v", serverUUID, err)
		respondWithError(w, http.StatusInternalServerError, "Failed to get voice channels")
		return
	}

	log.Printf("Retrieved %d voice channels for server %s", len(channels), serverUUID)

	simpleChannels := make([]SimpleChannelWithMembers, 0, len(channels))

	for _, channel := range channels {
		log.Printf("Processing channel %s", channel.ChannelID)

		var members []ChannelMember

		switch v := channel.Members.(type) {
		case []ChannelMember:
			members = v
		case []uint8:
			if err := json.Unmarshal(v, &members); err != nil {
				log.Printf("Failed to unmarshal members for channel %s: %v", channel.ChannelID, err)
				continue
			}
		case []interface{}:
			for _, m := range v {
				if member, ok := m.(ChannelMember); ok {
					members = append(members, member)
				} else {
					log.Printf("Unexpected member type for channel %s: %T", channel.ChannelID, m)
				}
			}
		case nil:
			log.Printf("Channel %s has nil Members, treating as empty", channel.ChannelID)
			members = []ChannelMember{}
		default:
			log.Printf("Unexpected type for channel.Members: %T", channel.Members)
			continue
		}

		if len(members) == 0 {
			log.Printf("Channel %s has no members after processing", channel.ChannelID)
		}

		simpleChannels = append(simpleChannels, SimpleChannelWithMembers{
			SimpleChannel: SimpleChannel{
				ChannelID:        channel.ChannelID,
				OwnerID:          channel.OwnerID,
				ServerID:         channel.ServerID,
				LanguageID:       channel.LanguageID,
				ChannelName:      channel.ChannelName,
				LastActive:       channel.LastActive.Time,
				IsLocked:         channel.IsLocked.Bool,
				ChannelCreatedAt: channel.ChannelCreatedAt,
				ChannelUpdatedAt: channel.ChannelUpdatedAt,
			},
			Members: members,
		})
	}

	log.Printf("Processed %d channels successfully", len(simpleChannels))

	response := GetServerVoiceChannelResponse{
		ServerID: serverUUID,
		Channels: simpleChannels,
	}

	respondWithJSON(w, http.StatusOK, response)
}
