package handlers

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jimmyvallejo/gleamspeak-api/utils"
)

func (h *Handlers) RefreshToken(w http.ResponseWriter, r *http.Request) {

	const (
		accessTokenExpirySeconds = 900
	)

	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		if err == http.ErrNoCookie {
			respondWithError(w, http.StatusUnauthorized, "Token not found")
		} else {
			respondWithError(w, http.StatusBadRequest, "Error reading cookie")
		}
		return
	}

	validated, err := utils.ValidateToken(cookie.Value, h.JWT)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error validating token")
		return
	}

	claims, ok := validated.Claims.(*jwt.RegisteredClaims)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Error validating token")
		return
	}

	expirationTime, err := claims.GetExpirationTime()
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error validating token")
		return
	}

	if time.Now().After(expirationTime.Time) {
		respondWithError(w, http.StatusUnauthorized, "Token is expired")
		return
	}

	idStr, err := claims.GetSubject()
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error validating token")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error parsing id")
		return
	}

	accessToken, err := utils.CreateToken(id, h.JWT, 900)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error issuing token")
		return
	}

	utils.SetTokenCookie(w, "access_token", accessToken, accessTokenExpirySeconds)

	respondNoBody(w, http.StatusOK)
}
