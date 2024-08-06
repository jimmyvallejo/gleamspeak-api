package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
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
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
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

type UpdateUserRequest struct {
	Email  string `json:"email"`
	Handle string `json:"handle"`
}

func (h *Handlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(common.UserContextKey).(database.User)
	
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unathorized")
		return
	}

	request := UpdateUserRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	params := database.UpdateUserByIDParams{
		Email:     request.Email,
		Handle:    request.Handle,
		UpdatedAt: time.Now(),
		ID:        user.ID,
	}
	updatedUser, err := h.DB.UpdateUserByID(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	response := UserResponse{
		ID:     updatedUser.ID,
		Email:  updatedUser.Email,
		Handle: updatedUser.Handle,
	}

	respondWithJSON(w, http.StatusOK, response)

}
