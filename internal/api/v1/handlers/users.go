package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/jimmyvallejo/gleamspeak-api/internal/database"
)

type createUserRequest struct {
	Email  string `json:"email"`
	Handle string `json:"handle"`
}

func (h *Handlers) CreateUser(w http.ResponseWriter, r *http.Request) {

	request := createUserRequest{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	params := database.CreateUserParams{
		ID:        uuid.New(),
		Email:     request.Email,
		Handle:    request.Handle,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	user, err := h.DB.CreateUser(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	response := CreateUserResponse{
		ID:     user.ID,
		Email:  user.Email,
		Handle: user.Handle,
	}

	respondWithJSON(w, http.StatusCreated, response)
}
