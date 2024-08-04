package handlers

import (
	"net/http"
)

func HandlerReadiness(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, StatusResponse{Status: http.StatusText(http.StatusOK)})
}

func HandlerError(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}
