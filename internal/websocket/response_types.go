package websocket

import (
	"time"

	"github.com/google/uuid"
)

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

type ChannelMemberExpanded struct {
	ChannelMember ChannelMember `json:"member"`
	Channel       uuid.UUID     `json:"channel_id"`
}

type ChannelMember struct {
	UserID uuid.UUID `json:"user_id"`
	Handle string    `json:"handle"`
}
