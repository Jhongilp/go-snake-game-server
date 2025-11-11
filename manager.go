package main

import (
	"encoding/json"
	"errors"
	"fmt"
	gamestate "github/jhongilp/snake-game/game-state"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
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
	clients   ClientList
	gameState *gamestate.GameState
	sync.RWMutex
	handlers map[string]EventHandler
}

func NewManager() *Manager {
	m := &Manager{
		clients:   make(ClientList),
		gameState: gamestate.NewGameState(),
		handlers:  make(map[string]EventHandler),
	}

	m.SetupEventHandlers()

	return m
}

func (m *Manager) SetupEventHandlers() {
	m.handlers[EventSendMessage] = SendMessage
	m.handlers[EventPlayerInput] = HandlePlayerInput
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

		// Remove player from game state
		m.gameState.RemovePlayer(c.playerId)
	}
}

// BroadcastGameState sends the game state to all connected clients
func (m *Manager) BroadcastGameState() {
	m.RLock()
	gameStatePayload, err := json.Marshal(m.gameState)
	m.RUnlock()

	if err != nil {
		log.Printf("error marshalling game state: %v", err)
		return
	}

	event := Event{
		Type:    EventGameUpdate,
		Payload: gameStatePayload,
	}

	m.RLock()
	for client := range m.clients {
		select {
		case client.egress <- event:
		default:
			// Client's egress channel is full, skip
		}
	}
	m.RUnlock()
}

// StartGameLoop runs the game loop at 8 ticks per second (~125ms)
func (m *Manager) StartGameLoop() {
	ticker := time.NewTicker(125 * time.Millisecond)
	defer ticker.Stop()

	log.Println("Game loop started")

	for range ticker.C {
		m.Lock()
		m.gameState.UpdateAllPlayers()
		m.Unlock()

		m.BroadcastGameState()
	}
}

func HandlePlayerInput(event Event, c *Client) error {
	var inputEvent PlayerInputEvent
	if err := json.Unmarshal(event.Payload, &inputEvent); err != nil {
		return fmt.Errorf("bad payload in request in player input: %v", err)
	}

	log.Printf("[handle player input] player: %s, direction: %s\n", c.playerId, inputEvent.Direction)

	// Update the player's direction (movement happens in game loop)
	c.manager.Lock()
	c.manager.gameState.SetPlayerDirection(c.playerId, gamestate.Direction(inputEvent.Direction))
	c.manager.Unlock()

	return nil
}

func SendMessage(event Event, c *Client) error {
	var messageEvent SendMessageEvent
	if err := json.Unmarshal(event.Payload, &messageEvent); err != nil {
		return fmt.Errorf("bad payload in request: %v", err)
	}

	log.Printf("Received message from client: %+v\n", messageEvent)

	// Create response to send back to client
	response := SendMessageEvent{
		Message: "Server received: " + messageEvent.Message,
		From:    "server",
	}

	responsePayload, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %v", err)
	}

	responseEvent := Event{
		Type:    EventSendMessage,
		Payload: responsePayload,
	}

	// Send to client through egress channel
	c.egress <- responseEvent

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

	// Generate unique player ID
	playerId := uuid.New().String()

	client := NewClient(conn, m)
	client.playerId = playerId
	m.AddClient(client)

	// Start client processes FIRST
	go client.readMessages()
	go client.writeMessages()

	// Send player ID to client after goroutines are running
	playerIdPayload, err := json.Marshal(PlayerIdEvent{Id: playerId})
	if err != nil {
		log.Printf("error marshalling player ID: %v", err)
		return
	}

	playerIdEvent := Event{
		Type:    EventPlayerId,
		Payload: playerIdPayload,
	}

	log.Printf("Assigned player ID: %s", playerId)
	client.egress <- playerIdEvent

	// create Player
	player := gamestate.NewPlayer(playerId, 100, 100, "red", "Player 1")

	// Add player to the shared game state
	m.Lock()
	m.gameState.AddPlayer(*player)
	m.Unlock()

	gameStatePayload, err := json.Marshal(m.gameState)
	if err != nil {
		log.Printf("error marshalling game state: %v", err)
		return
	}

	client.egress <- Event{
		Type:    EventGameState,
		Payload: gameStatePayload,
	}
}
