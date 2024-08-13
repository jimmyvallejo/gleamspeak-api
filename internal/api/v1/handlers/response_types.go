package handlers

import (
	"time"

	"github.com/google/uuid"
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

type SimpleServer struct {
	ServerID        uuid.UUID    `json:"server_id"`
	ServerName      string    `json:"server_name"`
	Description     string    `json:"description"` 
	IconURL         string    `json:"icon_url"`    
	BannerURL       string    `json:"banner_url"`  
	IsPublic        bool      `json:"is_public"`
	MemberCount     int32     `json:"member_count"`
	ServerLevel     int32     `json:"server_level"`
	MaxMembers      int32     `json:"max_members"`
	ServerCreatedAt time.Time `json:"server_created_at"`
	ServerUpdatedAt time.Time `json:"server_updated_at"`
}

type SimpleDisplayServerResponse struct {
	UserID  uuid.UUID      `json:"user_id"`
	Servers []SimpleServer `json:"servers"`
}
