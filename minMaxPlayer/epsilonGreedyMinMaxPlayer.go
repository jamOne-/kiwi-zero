package minMaxPlayer

import (
	"math/rand"

	"github.com/jamOne-/kiwi-zero/game"
)

type EpsilonGreedyMinMaxPlayer struct {
	depth   int
	epsilon float64
	valueFn ValueFn
}

func NewEpsilonGreedyMinMaxPlayer(depth int, epsilon float64, valueFn ValueFn) *EpsilonGreedyMinMaxPlayer {
	return &EpsilonGreedyMinMaxPlayer{depth, epsilon, valueFn}
}

func (player *EpsilonGreedyMinMaxPlayer) SelectMove(g game.Game) game.Move {
	var move game.Move

	if rand.Float64() < player.epsilon {
		possibleMoves := g.GetPossibleMoves()
		move = possibleMoves[rand.Intn(len(possibleMoves))]
	} else {
		_, move = negaMax(player.valueFn, g, player.depth, -INFINITY, INFINITY)
		return move
	}

	return move
}
