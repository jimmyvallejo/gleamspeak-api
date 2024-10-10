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
	EventSendMessage        = "send_message"
	EventNewMessage         = "new_message"
	EventChangeRoom         = "change_room"
	EventChangeVoiceRoom    = "change_voice_room"
	EventChangeServer       = "change_server"
	EventAddVoiceMember     = "add_voice_member"
	EventAddedVoiceMember   = "added_voice_member"
	EventRemoveVoiceMember  = "remove_voice_member"
	EventRemovedVoiceMember = "removed_voice_member"
)

type SendMessageEvent struct {
	Message string `json:"message"`
	From    string `json:"from"`
	Handle  string `json:"handle"`
	Channel string `json:"channel"`
	Image   string `json:"image"`
	Avatar  string `json:"avatar"`
}

type VoiceMemberEvent struct {
	User    string `json:"user_id"`
	Channel string `json:"channel_id"`
	Server  string `json:"server_id"`
	Handle  string `json:"handle"`
}

type ReturnEvent interface {
	GetType() string
}

type ReturnEventMessage struct {
	Type    string                 `json:"type"`
	Payload handlers.SimpleMessage `json:"payload"`
}

func (r ReturnEventMessage) GetType() string {
	return r.Type
}

type ReturnEventVoiceMember struct {
	Type    string                         `json:"type"`
	Payload handlers.ChannelMemberExpanded `json:"payload"`
}

func (r ReturnEventVoiceMember) GetType() string {
	return r.Type
}

type changeRoomEvent struct {
	ID string `json:"id"`
}
