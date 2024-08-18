package websocket

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jimmyvallejo/gleamspeak-api/internal/api/v1/handlers"
	"github.com/jimmyvallejo/gleamspeak-api/internal/database"
)

var (
	pongWait = 10 * time.Second

	pingInterval = (pongWait * 9 / 10)
)

type ClientList map[*Client]bool

type Client struct {
	connection *websocket.Conn
	manager    *Manager
	egress     chan Event
	chatroom   string
}

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		egress:     make(chan Event),
	}
}

func (c *Client) readMessages() {
	defer func() {
		c.manager.removeClient(c)
	}()

	if err := c.connection.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Println(err)
		return
	}

	c.connection.SetReadLimit(512)

	c.connection.SetPongHandler(c.pongHandler)

	for {
		_, payload, err := c.connection.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}

			break
		}

		var request Event
		if err := json.Unmarshal(payload, &request); err != nil {
			log.Printf("error marshalling")
			break
		}

		log.Printf("payload: %v", string(request.Type))
		if err := c.manager.routeEvent(request, c); err != nil {
			log.Println("error handling message")
		}
	}
}

func (c *Client) writeMessages() {
	defer func() {
		c.manager.removeClient(c)
	}()

	ticker := time.NewTicker(pingInterval)

	for {
		select {
		case message, ok := <-c.egress:
			if !ok {
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("connection closes:", err)
				}
				return
			}

			var wsParams = SendMessageResponse{}

			if err := json.Unmarshal(message.Payload, &wsParams); err != nil {
				log.Println("error unmarshaling message")
				continue
			}

			owner, err := uuid.Parse(wsParams.OwnerID.String())
			if err != nil {
				log.Println("failed to parse owner id")
				continue
			}
			channel, err := uuid.Parse(wsParams.ChannelID.String())
			if err != nil {
				log.Println("failed to parse channel id")
				continue
			}

			var createParams = database.CreateTextMessageParams{
				ID:        uuid.New(),
				OwnerID:   owner,
				ChannelID: channel,
				Message:   wsParams.Message,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			createdMessage, err := c.manager.DB.CreateTextMessage(context.Background(), createParams)
			if err != nil {
				log.Printf("failed to save message to database")
				continue
			}

			var response = handlers.SimpleMessage{
				ID:          createdMessage.ID,
				ChannelID:   createdMessage.ChannelID,
				OwnerID:     createdMessage.OwnerID,
				OwnerHandle: wsParams.Handle, 
				Message:     createdMessage.Message,
				Image:       createdMessage.Image.String, 
				CreatedAt:   createdMessage.CreatedAt,
				UpdatedAt:   createdMessage.UpdatedAt,
			}

			var sentEvent = ReturnEvent{
				Type:    message.Type,
				Payload: response,
			}

			data, err := json.Marshal(sentEvent)
			if err != nil {
				log.Println("error marshaling message")
				continue

			}

			log.Printf("Data before WriteMessage: %s", string(data))
			if err := c.connection.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Printf("failed to send message: %v", err)

			}
			log.Println("message sent")
		case <-ticker.C:
			log.Println("ping")

			if err := c.connection.WriteMessage(websocket.PingMessage, []byte(``)); err != nil {
				log.Println("writemsg err: ", err)
				return
			}
		}
	}
}

func (c *Client) pongHandler(pongMsg string) error {
	log.Println("pong")
	return c.connection.SetReadDeadline(time.Now().Add(pongWait))
}
