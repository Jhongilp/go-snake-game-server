package main

import (
	"encoding/json"
	gamestate "github/jhongilp/snake-game/game-state"
)

type Event struct {
	Type     string          `json:"type"`
	Payload  json.RawMessage `json:"payload"`
	PlayerId string          `json:"player_id"`
}

type EventHandler func(event Event, c *Client) error

const (
	EventSendMessage = "send_message"
	EventPlayerId    = "player_id"
	EventGameState   = "game_state"
	EventGameUpdate  = "game_update"
	EventPlayerInput = "player_input"
)

type SendMessageEvent struct {
	Message string `json:"message"`
	From    string `json:"from"`
}

type PlayerIdEvent struct {
	Id string `json:"id"`
}

type GameStateEvent struct {
	players []gamestate.Player
}

type PlayerInputEvent struct {
	Direction string `json:"direction"`
}
