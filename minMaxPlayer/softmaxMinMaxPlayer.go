package minMaxPlayer

import (
	"math"

	"github.com/jamOne-/kiwi-zero/game"
	"github.com/jamOne-/kiwi-zero/utils"
)

type SoftMaxMinMaxPlayer struct {
	depth   int
	valueFn game.ValueFn
}

func NewSoftMaxMinMaxPlayer(depth int, valueFn game.ValueFn) *SoftMaxMinMaxPlayer {
	return &SoftMaxMinMaxPlayer{depth, valueFn}
}

func (player *SoftMaxMinMaxPlayer) SelectMove(g game.Game) game.Move {
	moves := g.GetPossibleMoves()
	values := make([]float32, len(moves))

	for i, move := range moves {
		g.MakeMove(move)

		value, _ := negaMax(player.valueFn, g, player.depth-1, -INFINITY, INFINITY)
		values[i] = float32(-value)

		g.UndoLastMove()
	}

	index := utils.RandomFromDistribution(SoftMax32(values))
	return moves[index]
}

func SoftMax32(xs []float32) []float32 {
	result := make([]float32, len(xs))
	sum := float32(0)

	for i, x := range xs {
		result[i] = float32(math.Exp(float64(x * 10)))
		sum += result[i]
	}

	for i, _ := range xs {
		result[i] /= sum
	}

	return result
}
