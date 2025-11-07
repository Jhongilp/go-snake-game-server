package main

import "encoding/json"

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type EventHandler func(event Event, c *Client) error

const (
	EventSendMessage = "send_message"
	EventPlayerId    = "player_id"
)

type SendMessageEvent struct {
	Message  string `json:"message"`
	From     string `json:"from"`
	PlayerId string `json:"playerId"`
}

type PlayerIdEvent struct {
	Id string `json:"id"`
}
