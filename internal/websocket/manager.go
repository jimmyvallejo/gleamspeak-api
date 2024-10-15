package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jimmyvallejo/gleamspeak-api/internal/database"
	"github.com/jimmyvallejo/gleamspeak-api/internal/redis"
)

var webSocketUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		log.Printf("Incoming WebSocket connection attempt from origin: %s", origin)
		allowedOrigins := []string{"http://localhost:5173", "http://localhost:3000"}
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				return true
			}
		}

		log.Printf("Rejected WebSocket connection from origin: %s", origin)
		return false
	},
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

type Manager struct {
	clients ClientList
	DB      *database.Queries
	RDB     *redis.RedisClient
	sync.RWMutex

	handlers map[string]EventHandler
}

func NewManager(db *database.Queries, rdb *redis.RedisClient) *Manager {
	m := &Manager{
		clients:  make(ClientList),
		DB:       db,
		RDB:      rdb,
		handlers: make(map[string]EventHandler),
	}
	m.setupEventHandlers()
	return m
}

func (m *Manager) setupEventHandlers() {
	m.handlers[EventSendMessage] = SendMessage
	m.handlers[EventChangeRoom] = ChatRoomHandler
	m.handlers[EventChangeVoiceRoom] = VoiceRoomHandler
	m.handlers[EventChangeServer] = ServerChangeHandler
	m.handlers[EventAddVoiceMember] = AddVoiceMember
	m.handlers[EventRemoveVoiceMember] = RemoveVoiceMember
}

func (m *Manager) routeEvent(event Event, c *Client) error {
	if handler, ok := m.handlers[event.Type]; ok {
		if err := handler(event, c); err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("no such event type")
	}
}

func ChatRoomHandler(event Event, c *Client) error {
	var changeRoomEvent changeRoomEvent

	if err := json.Unmarshal(event.Payload, &changeRoomEvent); err != nil {
		return fmt.Errorf("bad payoad in req: %v", err)
	}
	c.chatroom = changeRoomEvent.ID
	return nil
}

func VoiceRoomHandler(event Event, c *Client) error {
	var changeRoomEvent changeRoomEvent

	if err := json.Unmarshal(event.Payload, &changeRoomEvent); err != nil {
		return fmt.Errorf("bad payoad in req: %v", err)
	}
	c.voiceroom = changeRoomEvent.ID
	return nil
}

func ServerChangeHandler(event Event, c *Client) error {
	payload := VoiceMemberEvent{
		User:    c.userID,
		Channel: c.voiceroom,
		Server:  c.server,
		Handle:  c.handle,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling payload: %v", err)
	}

	removeEvent := Event{
		Type:    EventRemoveVoiceMember,
		Payload: json.RawMessage(jsonPayload),
	}

	err = RemoveVoiceMember(removeEvent, c)
	if err != nil {
		log.Printf("Error removing voice member: %v", err)
	}

	var changeRoomEvent changeRoomEvent

	if err := json.Unmarshal(event.Payload, &changeRoomEvent); err != nil {
		return fmt.Errorf("bad payoad in req: %v", err)
	}
	c.server = changeRoomEvent.ID
	log.Printf("Changed Server to %v", changeRoomEvent.ID)
	return nil
}

func SendMessage(event Event, c *Client) error {
	var chatEvent SendMessageEvent
	if err := json.Unmarshal(event.Payload, &chatEvent); err != nil {
		return fmt.Errorf("bad payload in request: %v", err)
	}

	channelID, err := uuid.Parse(chatEvent.Channel)
	if err != nil {
		return fmt.Errorf("invalid UUID format for chatroom: %v", err)
	}

	ownerID, err := uuid.Parse(chatEvent.From)
	if err != nil {
		return fmt.Errorf("invalid UUID format for chatroom: %v", err)
	}

	var createParams = database.CreateTextMessageParams{
		ID:        uuid.New(),
		OwnerID:   ownerID,
		ChannelID: channelID,
		Message:   chatEvent.Message,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	createdMessage, err := c.manager.DB.CreateTextMessage(context.Background(), createParams)
	if err != nil {
		return fmt.Errorf("failed to add message to database: %v", err)
	}

	var response = SimpleMessage{
		ID:          createdMessage.ID,
		ChannelID:   createdMessage.ChannelID,
		OwnerID:     createdMessage.OwnerID,
		OwnerHandle: chatEvent.Handle,
		OwnerImage:  chatEvent.Avatar,
		Message:     createdMessage.Message,
		Image:       createdMessage.Image.String,
		CreatedAt:   createdMessage.CreatedAt,
		UpdatedAt:   createdMessage.UpdatedAt,
	}

	payload, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("error marshaling json for response: %v", err)
	}

	outgoingEvent := Event{
		Payload: payload,
		Type:    EventNewMessage,
	}

	for client := range c.manager.clients {
		if client.chatroom == c.chatroom {
			client.egress <- outgoingEvent
		}
	}
	return nil
}

func AddVoiceMember(event Event, c *Client) error {

	err := RemoveVoiceMember(event, c)
	if err != nil {
		return fmt.Errorf("error removing user prior to add: %v", err)
	}

	var memberEvent VoiceMemberEvent
	if err := json.Unmarshal(event.Payload, &memberEvent); err != nil {
		return fmt.Errorf("bad payload in request: %v", err)
	}

	userUUID, err := uuid.Parse(memberEvent.User)
	if err != nil {
		return fmt.Errorf("invalid UUID format for user: %v", err)
	}

	channelUUID, err := uuid.Parse(memberEvent.Channel)
	if err != nil {
		return fmt.Errorf("invalid UUID format for channel: %v", err)
	}

	serverUUID, err := uuid.Parse(memberEvent.Server)
	if err != nil {
		return fmt.Errorf("invalid UUID format for server: %v", err)
	}

	var createParams = database.JoinVoiceChannelParams{
		UserID:    userUUID,
		ChannelID: channelUUID,
		ServerID:  serverUUID,
	}

	_, err = c.manager.DB.JoinVoiceChannel(context.Background(), createParams)
	if err != nil {
		return fmt.Errorf("failed to add voice room member to database: %v", err)
	}

	m := ChannelMember{
		UserID: userUUID,
		Handle: memberEvent.Handle,
	}

	response := ChannelMemberExpanded{
		ChannelMember: m,
		Channel:       channelUUID,
	}

	payload, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("error marshaling json for response: %v", err)
	}

	outgoingAdd := Event{
		Payload: payload,
		Type:    EventAddedVoiceMember,
	}

	for client := range c.manager.clients {
		if client.server == memberEvent.Server {
			client.egress <- outgoingAdd
		}
	}

	return nil
}

func RemoveVoiceMember(event Event, c *Client) error {
	var memberEvent VoiceMemberEvent
	if err := json.Unmarshal(event.Payload, &memberEvent); err != nil {
		return fmt.Errorf("bad payload in request: %v", err)
	}
	userUUID, err := uuid.Parse(memberEvent.User)
	if err != nil {
		return fmt.Errorf("invalid UUID format for user: %v", err)
	}

	channelUUID, err := uuid.Parse(memberEvent.Channel)
	if err != nil {
		return fmt.Errorf("invalid UUID format for channel: %v", err)
	}

	m := ChannelMember{
		UserID: userUUID,
		Handle: memberEvent.Handle,
	}

	response := ChannelMemberExpanded{
		ChannelMember: m,
		Channel:       channelUUID,
	}

	payload, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("error marshaling json for response: %v", err)
	}

	outgoingRemove := Event{
		Payload: payload,
		Type:    EventRemovedVoiceMember,
	}

	for client := range c.manager.clients {
		if client.server == memberEvent.Server {
			client.egress <- outgoingRemove
		}
	}
	return nil
}

func (m *Manager) ServeWs(w http.ResponseWriter, r *http.Request, userId, handle string) {
	log.Println("New WebSocket connection attempt")

	conn, err := webSocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection to WebSocket: %v", err)
		return
	}

	log.Println("WebSocket connection established successfully")

	client := NewClient(conn, m, userId, handle)
	m.addClient(client)

	go client.readMessages()
	go client.writeMessages()
}

func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	m.clients[client] = true
}

func (m *Manager) removeClient(client *Client) {
	payload := VoiceMemberEvent{
		User:    client.userID,
		Channel: client.voiceroom,
		Server:  client.server,
		Handle:  client.handle,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling payload: %v", err)
	}

	event := Event{
		Type:    EventRemoveVoiceMember,
		Payload: json.RawMessage(jsonPayload),
	}

	err = RemoveVoiceMember(event, client)
	if err != nil {
		log.Printf("Error removing voice member: %v", err)
	}

	m.Lock()
	defer m.Unlock()

	if _, ok := m.clients[client]; ok {
		if client.egress != nil {
			close(client.egress)
		}

		if err := client.connection.Close(); err != nil {
			log.Printf("Error closing client connection: %v", err)
		}

		delete(m.clients, client)

		log.Printf("Client %s removed from manager", client.userID)
	} else {
		log.Printf("Client %s not found in manager", client.userID)
	}
}

// func checkOrigin(r *http.Request) bool {
// 	origin := r.Header.Get("Origin")

// 	switch origin {
// 	case "http://localhost:5173":
// 		return true
// 	default:
// 		return false
// 	}
// }
