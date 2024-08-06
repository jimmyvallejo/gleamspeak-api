package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jimmyvallejo/gleamspeak-api/utils"
	"golang.org/x/crypto/bcrypt"
)

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handlers) LoginUserStandard(w http.ResponseWriter, r *http.Request) {
	const (
		accessTokenExpirySeconds  = 900
		refreshTokenExpirySeconds = 604800
	)

	request := LoginUserRequest{}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	user, err := h.DB.GetUserByEmail(r.Context(), request.Email)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Email or password incorrect")
		return
	}

	if !user.Password.Valid {
		respondWithError(w, http.StatusInternalServerError, "Invalid user data")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password.String), []byte(request.Password))
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Email or password incorrect")
		return
	}

	accessToken, err := utils.CreateAccessToken(user.ID, h.JWT, 900)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error issuing token")
		return
	}
	refreshToken, err := utils.CreateRefreshToken(h.JWT, 604800)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error issuing token")
		return
	}

	utils.SetTokenCookie(w, "access_token", accessToken, accessTokenExpirySeconds)
	utils.SetTokenCookie(w, "refresh_token", refreshToken, refreshTokenExpirySeconds)

	respondNoBody(w, http.StatusOK)
}
