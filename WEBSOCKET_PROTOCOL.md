# WebSocket Protocol for Multiplayer Snake Tron

This document describes the expected WebSocket message format between the client and the Go server.

## Connection
- **URL**: `ws://localhost:8080/ws`
- **Format**: All messages are JSON strings

## Message Types

### 1. Server → Client: Player ID Assignment
When a client connects, the server should send their player ID.

```json
{
  "type": "player_id",
  "payload": {
    "id": "unique-player-id-string"
  }
}
```

### 2. Client → Server: Player Input
When the player changes direction, the client sends:

```json
{
  "type": "player_input",
  "payload": {
    "direction": "up" | "down" | "left" | "right"
  }
}
```

### 3. Server → Client: Full Game State
Complete game state update (typically sent when game starts or resets):

```json
{
  "type": "game_state",
  "payload": {
    "players": [
      {
        "id": "player-id-string",
        "x": 100,
        "y": 150,
        "color": "#00FFFF",
        "trail": [
          { "x": 95, "y": 150 },
          { "x": 100, "y": 150 }
        ],
        "alive": true,
        "score": 0,
        "name": "Player 1",
        "length": 3,
        "maxLength": 3
      }
    ],
    "powerUps": [
      {
        "x": 200,
        "y": 300,
        "type": "length",
        "active": true
      }
    ],
    "gameRunning": true
  }
}
```

### 4. Server → Client: Game Update
Incremental game state updates (sent each game tick):

```json
{
  "type": "game_update",
  "payload": {
    "players": [
      {
        "id": "player-id-string",
        "x": 105,
        "y": 150,
        "color": "#00FFFF",
        "trail": [
          { "x": 100, "y": 150 },
          { "x": 105, "y": 150 }
        ],
        "alive": true,
        "score": 0,
        "name": "Player 1",
        "length": 3,
        "maxLength": 3
      }
    ],
    "powerUps": [
      {
        "x": 200,
        "y": 300,
        "type": "length",
        "active": false
      }
    ],
    "gameRunning": true
  }
}
```

## Game State Properties

### Player Object
- `id` (string): Unique identifier for the player
- `x` (number): Current X position (grid-based, multiples of 5)
- `y` (number): Current Y position (grid-based, multiples of 5)
- `color` (string): Hex color code for the player (e.g., "#00FFFF", "#FF00FF")
- `trail` (array): Array of {x, y} points representing the player's trail
- `alive` (boolean): Whether the player is still alive
- `score` (number): Player's current score
- `name` (string): Player's display name
- `length` (number): Current trail length
- `maxLength` (number): Maximum allowed trail length

### PowerUp Object
- `x` (number): X position of the power-up
- `y` (number): Y position of the power-up
- `type` (string): Type of power-up (currently only "length")
- `active` (boolean): Whether the power-up is still available to collect

### Game State
- `players` (array): List of all players in the game
- `powerUps` (array, optional): List of power-ups on the grid
- `gameRunning` (boolean): Whether the game is currently active

## Expected Server Behavior

1. **On Connection**: Send `player_id` message to the new client
2. **On Player Input**: Update that player's direction in the game state
3. **Game Loop**: Every ~125ms (8 ticks per second):
   - Update all player positions based on their directions
   - Check for collisions (walls, trails)
   - Check for power-up collection
   - Send `game_update` to all connected clients
4. **Game Over**: Set `gameRunning: false` and send final state
5. **Restart**: After game over, reset and send new `game_state`

## Client Behavior

- Connects to WebSocket on load
- Sends direction changes on keyboard input (WASD or Arrow keys)
- Renders game state received from server
- Highlights current player with yellow glow
- Shows disconnection message if WebSocket closes
