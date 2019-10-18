package minMaxPlayer

import (
	"math"
	"math/rand"

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
		return valueFn(g), game.Move(-1)
	}

	moves := g.GetPossibleMoves()
	bestValue, bestMoves := -INFINITY, []game.Move{game.Move(-1)}

	for _, move := range moves {
		gameCopy := g.Copy()
		gameCopy.MakeMove(move)

		value, _ := negaMax(valueFn, gameCopy, depth-1, -b, -a)
		value = -value

		if value > bestValue {
			bestValue = value
			bestMoves = []game.Move{move}
		} else if value == bestValue {
			bestMoves = append(bestMoves, move)
		}

		a = math.Max(a, value)

		if a >= b {
			break
		}
	}

	bestMoveIndex := rand.Intn(len(bestMoves))
	bestMove := bestMoves[bestMoveIndex]
	return bestValue, bestMove
}
