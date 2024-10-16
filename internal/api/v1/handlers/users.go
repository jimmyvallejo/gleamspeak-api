package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jimmyvallejo/gleamspeak-api/internal/api/common"
	"github.com/jimmyvallejo/gleamspeak-api/internal/database"
	"golang.org/x/crypto/bcrypt"
)

type createUserRequest struct {
	Email    string `json:"email"`
	Handle   string `json:"handle"`
	Password string `json:"password"`
}

func (h *Handlers) CreateUserStandard(w http.ResponseWriter, r *http.Request) {
	request := createUserRequest{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	passwordBytes := []byte(request.Password)
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword(passwordBytes, 12)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error hashing password")
		return
	}

	params := database.CreateUserStandardParams{
		ID:     uuid.New(),
		Email:  request.Email,
		Handle: request.Handle,
		Password: sql.NullString{
			String: string(hashedPasswordBytes),
			Valid:  true,
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	user, err := h.DB.CreateUserStandard(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	roleID, err := h.DB.GetRoleIDByName(r.Context(), "member")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting role ID")
		return
	}

	userRolesParams := database.CreateUserRolesParams{
		UserID: user.ID,
		RoleID: roleID,
	}

	_, err = h.DB.CreateUserRoles(r.Context(), userRolesParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error assigning role")
		return
	}

	response := UserResponse{
		ID:     user.ID,
		Email:  user.Email,
		Handle: user.Handle,
	}

	respondWithJSON(w, http.StatusCreated, response)
}

func (h *Handlers) FetchAuthUser(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(common.UserContextKey).(database.User)

	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unathorized")
		return
	}

	response := FullUserResponse{
		ID:         user.ID,
		Email:      user.Email,
		Handle:     user.Handle,
		IsActive:   user.IsActive.Bool,
		FirstName:  user.FirstName.String,
		LastName:   user.LastName.String,
		Bio:        user.Bio.String,
		AvatarURL:  user.AvatarUrl.String,
		IsVerified: user.IsVerified.Bool,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}

	respondWithJSON(w, http.StatusOK, response)
}

type UpdateUserRequest struct {
	Email     string `json:"email"`
	Handle    string `json:"handle"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Bio       string `json:"bio"`
}

func (h *Handlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(common.UserContextKey).(database.User)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	request := UpdateUserRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	params := database.UpdateUserByIDParams{
		Email:     user.Email,
		Handle:    user.Handle,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Bio:       user.Bio,
		UpdatedAt: time.Now().UTC(),
		ID:        user.ID,
	}

	if request.Email != "" {
		params.Email = request.Email
	}
	if request.Handle != "" {
		params.Handle = request.Handle
	}
	if request.FirstName != "" {
		params.FirstName = sql.NullString{
			String: request.FirstName,
			Valid:  true,
		}
	} else {
		params.FirstName = sql.NullString{Valid: false}
	}
	if request.LastName != "" {
		params.LastName = sql.NullString{
			String: request.LastName,
			Valid:  true,
		}
	} else {
		params.LastName = sql.NullString{Valid: false}
	}
	if request.Bio != "" {
		params.Bio = sql.NullString{
			String: request.Bio,
			Valid:  true,
		}
	} else {
		params.Bio = sql.NullString{Valid: false}
	}

	updatedUser, err := h.DB.UpdateUserByID(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	userIDStr := user.ID.String()

	err = h.RDB.SetJson("user"+userIDStr, updatedUser, time.Hour)
	if err != nil {
		log.Printf("Failed to save user to cache: %v", err)
	} else {
		log.Printf("User saved to cache: ID=%s", user.ID)
	}

	response := UserResponse{
		ID:     updatedUser.ID,
		Email:  updatedUser.Email,
		Handle: updatedUser.Handle,
	}

	respondWithJSON(w, http.StatusOK, response)
}

type UpdateUserAvatarRequest struct {
	URL *string `json:"url"`
}

func (h *Handlers) UpdateAvatar(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(common.UserContextKey).(database.User)

	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unathorized")
		return
	}

	request := UpdateUserAvatarRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	var avatarURL sql.NullString
	if request.URL != nil {
		avatarURL = sql.NullString{
			String: *request.URL,
			Valid:  true,
		}
	}

	params := database.UpdateUserAvatarByIDParams{
		ID:        user.ID,
		AvatarUrl: avatarURL,
	}

	user, err = h.DB.UpdateUserAvatarByID(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	userIDStr := user.ID.String()

	err = h.RDB.SetJson("user"+userIDStr, user, time.Hour)
	if err != nil {
		log.Printf("Failed to save user to cache: %v", err)
	} else {
		log.Printf("User saved to cache: ID=%s", user.ID)
	}

	respondNoBody(w, http.StatusOK)
}

func (h *Handlers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(common.UserContextKey).(database.User)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unathorized")
		return
	}

	userID := strings.TrimPrefix(r.URL.Path, "/v1/users/")

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to parse uuid, possible params error")
		return
	}

	if user.ID != userUUID {
		respondWithError(w, http.StatusForbidden, "Forbidden")
		return
	}

	err = h.DB.DeleteUser(r.Context(), userUUID)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "Failed to delete server")
		return
	}

	respondNoBody(w, http.StatusOK)
}
