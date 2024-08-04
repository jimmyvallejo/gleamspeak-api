package handlers

import "github.com/google/uuid"

type StatusResponse struct {
	Status string `json:"status"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type CreateUserResponse struct {
	ID     uuid.UUID `json:"id"`
	Email  string    `json:"email"`
	Handle string    `json:"handle"`
}
