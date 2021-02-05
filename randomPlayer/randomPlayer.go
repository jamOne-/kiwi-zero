package randomPlayer

import (
	"math/rand"

	"github.com/jamOne-/kiwi-zero/game"
)

type RandomPlayer struct{}

func NewRandomPlayer() *RandomPlayer {
	return &RandomPlayer{}
}

func (player *RandomPlayer) SelectMove(g game.Game) game.Move {
	return SelectRandomMove(g)
}

func (player *RandomPlayer) SelectMoveDifferentThan(g game.Game, move game.Move) game.Move {
	return SelectRandomMoveDifferentThan(g, move)
}

func SelectRandomMove(g game.Game) game.Move {
	possibleMoves := g.GetPossibleMoves()
	return possibleMoves[rand.Intn(len(possibleMoves))]
}

func SelectRandomMoveDifferentThan(g game.Game, move game.Move) game.Move {
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
