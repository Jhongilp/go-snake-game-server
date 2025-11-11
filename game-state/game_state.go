package gamestate

type GameState struct {
	Players []Player `json:"players"`
}

type Player struct {
	Id        string  `json:"id"`
	X         int     `json:"x"`
	Y         int     `json:"y"`
	Color     string  `json:"color"`
	Trail     []Point `json:"trail"`
	Alive     bool    `json:"alive"`
	Score     int     `json:"score"`
	Name      string  `json:"name"`
	Length    int     `json:"length"`
	MaxLength int     `json:"max_length"`
}

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

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

func (gs *GameState) UpdatePlayerPosition(playerID string, newX, newY int) {
	for i, p := range gs.Players {
		if p.Id == playerID {
			gs.Players[i].X = newX
			gs.Players[i].Y = newY
			// Add current position to trail
			gs.Players[i].Trail = append(gs.Players[i].Trail, Point{X: newX, Y: newY})
			// Trim trail if it exceeds MaxLength
			if len(gs.Players[i].Trail) > gs.Players[i].Length {
				gs.Players[i].Trail = gs.Players[i].Trail[1:]
			}
			return
		}
	}
}
