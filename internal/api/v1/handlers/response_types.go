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
	Avatar string    `json:"avatar"`
}

type FullUserResponse struct {
	ID         uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	Handle     string    `json:"handle"`
	IsActive   bool      `json:"is_active"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Bio        string    `json:"bio"`
	AvatarURL  string    `json:"avatar_url"`
	IsVerified bool      `json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CreateServerResponse struct {
	ID         uuid.UUID `json:"id"`
	OwnerID    uuid.UUID `json:"owner_id"`
	ServerName string    `json:"server_name"`
}

type SimpleServer struct {
	ServerID        uuid.UUID `json:"server_id"`
	OwnerID         uuid.UUID `json:"owner_id"`
	ServerName      string    `json:"server_name"`
	Description     string    `json:"description"`
	IconURL         string    `json:"icon_url"`
	BannerURL       string    `json:"banner_url"`
	IsPublic        bool      `json:"is_public"`
	InviteCode      string    `json:"invite_code"`
	MemberCount     int32     `json:"member_count"`
	ServerLevel     int32     `json:"server_level"`
	MaxMembers      int32     `json:"max_members"`
	ServerCreatedAt time.Time `json:"server_created_at"`
	ServerUpdatedAt time.Time `json:"server_updated_at"`
}

type SimpleRecentServer struct {
	ServerID        uuid.UUID `json:"server_id"`
	ServerName      string    `json:"server_name"`
	Description     string    `json:"description"`
	IconURL         string    `json:"icon_url"`
	BannerURL       string    `json:"banner_url"`
	MemberCount     int32     `json:"member_count"`
	ServerCreatedAt time.Time `json:"server_created_at"`
	ServerUpdatedAt time.Time `json:"server_updated_at"`
	OwnerHandle     string    `json:"owner_handle"`
	OwnerAvatar     string    `json:"owner_avatar"`
}

type SimpleDisplayServerResponse struct {
	UserID  uuid.UUID      `json:"user_id"`
	Servers []SimpleServer `json:"servers"`
}

type CreateTextChannelResponse struct {
	ID          uuid.UUID `json:"channel_id"`
	OwnerID     uuid.UUID `json:"owner_id"`
	ServerID    uuid.UUID `json:"server_id"`
	ChannelName string    `json:"channel_name"`
}

type SimpleChannel struct {
	ChannelID        uuid.UUID `json:"channel_id"`
	OwnerID          uuid.UUID `json:"owner_id"`
	ServerID         uuid.UUID `json:"server_id"`
	LanguageID       uuid.UUID `json:"language_id"`
	ChannelName      string    `json:"channel_name"`
	LastActive       time.Time `json:"last_active,omitempty"`
	IsLocked         bool      `json:"is_locked"`
	ChannelCreatedAt time.Time `json:"channel_created_at"`
	ChannelUpdatedAt time.Time `json:"channel_updated_at"`
}

type ChannelMember struct {
	UserID uuid.UUID `json:"user_id"`
	Handle string    `json:"handle"`
}

type ChannelMemberExpanded struct {
	ChannelMember ChannelMember `json:"member"`
	Channel       uuid.UUID     `json:"channel_id"`
}

type SimpleChannelWithMembers struct {
	SimpleChannel
	Members []ChannelMember `json:"members"`
}

type GetServerTextChannelResponse struct {
	ServerID uuid.UUID       `json:"server_id"`
	Channels []SimpleChannel `json:"channels"`
}

type GetServerVoiceChannelResponse struct {
	ServerID uuid.UUID                  `json:"server_id"`
	Channels []SimpleChannelWithMembers `json:"channels"`
}

type SimpleMessage struct {
	ID          uuid.UUID `json:"id"`
	OwnerID     uuid.UUID `json:"owner_id"`
	OwnerHandle string    `json:"handle"`
	OwnerImage  string    `json:"owner_image"`
	ChannelID   uuid.UUID `json:"channel_id"`
	Message     string    `json:"message"`
	Image       string    `json:"image"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SignedURLResponse struct {
	URL       string `json:"url"`
	PublicURL string `json:"public_url"`
}
