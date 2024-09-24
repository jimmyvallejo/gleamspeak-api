package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jimmyvallejo/gleamspeak-api/internal/api/common"
	"github.com/jimmyvallejo/gleamspeak-api/internal/database"
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

	accessToken, err := utils.CreateToken(user.ID, h.JWT, 900)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error issuing token")
		return
	}
	refreshToken, err := utils.CreateToken(user.ID, h.JWT, 604800)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error issuing token")
		return
	}

	utils.SetTokenCookie(w, "access_token", accessToken, accessTokenExpirySeconds)
	utils.SetTokenCookie(w, "refresh_token", refreshToken, refreshTokenExpirySeconds)

	respondNoBody(w, http.StatusOK)
}

func (h *Handlers) CheckAuthStatus(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(common.UserContextKey).(database.User)

	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unathorized")
		return
	}

	response := UserResponse{
		ID:     user.ID,
		Email:  user.Email,
		Handle: user.Handle,
		Avatar: user.AvatarUrl.String,
	}

	respondWithJSON(w, http.StatusOK, response)

}

func (h *Handlers) LogoutUserStandard(w http.ResponseWriter, r *http.Request) {

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	respondNoBody(w, http.StatusOK)
}
