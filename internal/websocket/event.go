package websocket

import (
	"encoding/json"

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
	Avatar  string `json:"avatar"`
}

type ReturnEvent struct {
	Type    string                 `json:"type"`
	Payload handlers.SimpleMessage `json:"payload"`
}

type changeRoomEvent struct {
	ID string `json:"id"`
}
