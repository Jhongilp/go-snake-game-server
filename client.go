package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

var (
	pongWait     = 10 * time.Second
	pingInterval = (pongWait * 9) / 10
)

type Client struct {
	connection *websocket.Conn
	manager    *Manager
	playerId   string

	// egress is used to avoid concurrent writes on the web socket connection
	egress chan Event
}

type ClientList map[*Client]bool

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		egress:     make(chan Event),
	}
}

func (c *Client) readMessages() {
	defer func() {
		// cleanup connection
		log.Println("removing connection")
		c.manager.RemoveClient(c)
	}()

	if err := c.connection.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Println(err)
		return
	}

	c.connection.SetReadLimit(512)
	c.connection.SetPongHandler(c.PongHandler)

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
			log.Printf("error marshalling event: %v", err)
			break
		}

		log.Printf("Event received: %+v", request)

		if err := c.manager.routeEvent(request, c); err != nil {
			log.Println("Error handling message: ", err)
		}

	}

}

func (c *Client) writeMessages() {
	defer func() {
		log.Println("[write message] removing client")
		c.manager.RemoveClient(c)
	}()

	ticket := time.NewTicker(pingInterval)

	for {
		select {
		case message, ok := <-c.egress:
			if !ok {
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("connection closed: ", err)
				}
				return
			}

			data, err := json.Marshal(message)
			if err != nil {
				log.Println("[write message] error marshalling message: ", err)
				return
			}

			if err := c.connection.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Printf("failed to send message %v", err)
			}
			log.Println("message sent")

		case <-ticket.C:
			log.Println("ping")
			if err := c.connection.WriteMessage(websocket.PingMessage, []byte(``)); err != nil {
				log.Println("closing, ping failed: ", err)
				return
			}
		}
	}
}

func (c *Client) PongHandler(pongMsg string) error {
	log.Println("pong")
	return c.connection.SetReadDeadline(time.Now().Add(pongWait))
}
