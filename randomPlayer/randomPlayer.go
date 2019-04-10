package randomPlayer

import (
	"math/rand"
	"time"

	"github.com/jamOne-/kiwi-zero/game"
)

type RandomPlayer struct{}

func NewRandomPlayer() *RandomPlayer {
	rand.Seed(time.Now().UnixNano())
	return &RandomPlayer{}
}

func (player *RandomPlayer) SelectMove(g game.Game) game.Move {
	possibleMoves := g.GetPossibleMoves()
	return possibleMoves[rand.Intn(len(possibleMoves))]
}
