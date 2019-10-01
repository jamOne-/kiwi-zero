package minMaxPlayer

import (
	"math"

	"github.com/jamOne-/kiwi-zero/game"
)

const INFINITY = 99999999.0

type ValueFn func(game.Game) float64

type MinMaxPlayer struct {
	depth   int
	valueFn ValueFn
}

func NewMinMaxPlayer(depth int, valueFn ValueFn) *MinMaxPlayer {
	return &MinMaxPlayer{depth, valueFn}
}

func (player *MinMaxPlayer) SelectMove(game game.Game) game.Move {
	_, move := negaMax(player.valueFn, game, player.depth, -INFINITY, INFINITY)
	return move
}

func negaMax(valueFn ValueFn, g game.Game, depth int, a float64, b float64) (float64, game.Move) {
	if finished, winner := g.IsGameFinished(); finished {
		return INFINITY * float64(winner*g.GetCurrentPlayerColor()), game.Move(-1)
	}

	if depth == 0 {
		return float64(g.GetCurrentPlayerColor()) * valueFn(g), game.Move(-1)
	}

	moves := g.GetPossibleMoves()
	bestValue, bestMove := -INFINITY, game.Move(-1)

	for _, move := range moves {
		gameCopy := g.Copy()
		gameCopy.MakeMove(move)

		value, _ := negaMax(valueFn, gameCopy, depth-1, -b, -a)
		value = -value

		if value > bestValue {
			bestValue = value
			bestMove = move
		}

		a = math.Max(a, value)

		if a >= b {
			break
		}
	}

	return bestValue, bestMove
}
