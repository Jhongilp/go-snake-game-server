package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			return origin == "http://localhost:5173"
		},
	}
)

type Manager struct {
	clients ClientList
	sync.RWMutex
	handlers map[string]EventHandler
}

func NewManager() *Manager {
	m := &Manager{
		clients:  make(ClientList),
		handlers: make(map[string]EventHandler),
	}

	m.SetupEventHandlers()

	return m
}

func (m *Manager) SetupEventHandlers() {
	m.handlers[EventSendMessage] = SendMessage
}

func (m *Manager) AddClient(c *Client) {
	m.Lock()
	defer m.Unlock()

	m.clients[c] = true
}

func (m *Manager) RemoveClient(c *Client) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.clients[c]; ok {
		c.connection.Close()
		delete(m.clients, c)
	}
}

func SendMessage(event Event, c *Client) error {
	fmt.Printf("Sending message: %+v\n", event.Payload)
	return nil
}

func (m *Manager) routeEvent(event Event, c *Client) error {
	log.Printf("[route] event: %+v\n", event)
	if handler, ok := m.handlers[event.Type]; ok {
		if err := handler(event, c); err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("there is no such event type")
	}
}

func (m *Manager) ServeWS(w http.ResponseWriter, r *http.Request) {
	log.Println("New connection")

	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := NewClient(conn, m)
	m.AddClient(client)

	// start client processes
	go client.readMessages()
	go client.writeMessages()
}
