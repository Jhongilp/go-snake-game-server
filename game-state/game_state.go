package gamestate

type GameState struct {
	Players []Player `json:"players"`
}

type Player struct {
	Id        string    `json:"id"`
	X         int       `json:"x"`
	Y         int       `json:"y"`
	Direction Direction `json:"-"` // Don't send to client
	Color     string    `json:"color"`
	Trail     []Point   `json:"trail"`
	Alive     bool      `json:"alive"`
	Score     int       `json:"score"`
	Name      string    `json:"name"`
	Length    int       `json:"length"`
	MaxLength int       `json:"max_length"`
}

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Direction string

const (
	DirectionUp    Direction = "up"
	DirectionDown  Direction = "down"
	DirectionLeft  Direction = "left"
	DirectionRight Direction = "right"
)

func NewGameState() *GameState {
	return &GameState{
		Players: make([]Player, 0),
	}
}

func NewPlayer(id string, x, y int, color, name string) *Player {
	trail := make([]Point, 0)
	trail = append(trail, Point{X: x + 5, Y: y}, Point{X: x + 10, Y: y})
	return &Player{
		Id:        id,
		X:         x,
		Y:         y,
		Direction: DirectionRight, // Default direction
		Color:     color,
		Trail:     trail,
		Alive:     true,
		Score:     0,
		Name:      name,
		Length:    1,
		MaxLength: 10,
	}
}

func (gs *GameState) AddPlayer(player Player) {
	gs.Players = append(gs.Players, player)
}

func (gs *GameState) RemovePlayer(playerID string) {
	for i, p := range gs.Players {
		if p.Id == playerID {
			gs.Players = append(gs.Players[:i], gs.Players[i+1:]...)
			return
		}
	}
}

func (gs *GameState) UpdatePlayerPosition(playerID string, direction Direction) {
	for i, p := range gs.Players {
		if p.Id == playerID {
			switch direction {
			case DirectionUp:
				p.Y--
			case DirectionDown:
				p.Y++
			case DirectionLeft:
				p.X--
			case DirectionRight:
				p.X++
			}
			gs.Players[i] = p
			return
		}
	}
}

func (gs *GameState) SetPlayerDirection(playerID string, direction Direction) {
	for i := range gs.Players {
		if gs.Players[i].Id == playerID {
			gs.Players[i].Direction = direction
			return
		}
	}
}

// UpdateAllPlayers moves all players based on their current direction
func (gs *GameState) UpdateAllPlayers() {
	for i := range gs.Players {
		if !gs.Players[i].Alive {
			continue
		}

		// Move player based on direction (grid-based, multiples of 5)
		switch gs.Players[i].Direction {
		case DirectionUp:
			gs.Players[i].Y -= 5
		case DirectionDown:
			gs.Players[i].Y += 5
		case DirectionLeft:
			gs.Players[i].X -= 5
		case DirectionRight:
			gs.Players[i].X += 5
		}

		// Add current position to trail
		newPoint := Point{X: gs.Players[i].X, Y: gs.Players[i].Y}
		gs.Players[i].Trail = append(gs.Players[i].Trail, newPoint)

		// Keep trail at max length
		if len(gs.Players[i].Trail) > gs.Players[i].MaxLength {
			gs.Players[i].Trail = gs.Players[i].Trail[1:]
		}
	}
}
