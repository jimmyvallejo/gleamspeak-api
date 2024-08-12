package handlers

import (
	"github.com/google/uuid"
	"github.com/jimmyvallejo/gleamspeak-api/internal/database"
)

type StatusResponse struct {
	Status string `json:"status"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type UserResponse struct {
	ID     uuid.UUID `json:"id"`
	Email  string    `json:"email"`
	Handle string    `json:"handle"`
}

type CreateServerResponse struct {
	ID         uuid.UUID `json:"id"`
	OwnerID    uuid.UUID `json:"owner_id"`
	ServerName string    `json:"server_name"`
}

type DisplayServerResponse struct {
	UserID      uuid.UUID         `json:"user_id"`
	Servers []database.GetUserServersRow `json:"servers"`
}
