package websocket

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jimmyvallejo/gleamspeak-api/internal/api/v1/handlers"
)

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type EventHandler func(event Event, c *Client) error

const (
	EventSendMessage = "send_message"
	EventNewMessage  = "new_message"
	EventChangeRoom  = "change_room"
)

type SendMessageEvent struct {
	Message string `json:"message"`
	From    string `json:"from"`
	Handle  string `json:"handle"`
	Channel string `json:"channel"`
	Image   string `json:"image"`
}

type SendMessageResponse struct {
	OwnerID   uuid.UUID `json:"owner_id"`
	Handle    string    `json:"handle"`
	ChannelID uuid.UUID `json:"channel_id"`
	Message   string    `json:"message"`
	Image     string    `json:"image"`
}

type ReturnEvent struct {
	Type    string                 `json:"type"`
	Payload handlers.SimpleMessage `json:"payload"`
}

type NewMessageEvent struct {
	SendMessageEvent
	Sent time.Time `json:"sent"`
}

type changeRoomEvent struct {
	ID string `json:"id"`
}
