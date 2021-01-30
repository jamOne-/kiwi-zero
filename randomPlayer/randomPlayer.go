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

func (player *RandomPlayer) SelectMoveDifferentThan(g game.Game, move game.Move) game.Move {
	possibleMoves := g.GetPossibleMoves()
	possibilities := len(possibleMoves)

	if possibilities == 1 {
		return possibleMoves[0]
	}

	selectedMove := move
	for selectedMove == move {
		selectedMove = possibleMoves[rand.Intn(possibilities)]
	}

	return selectedMove
}
