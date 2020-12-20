package policyPlayer

import (
	tfpredictor "github.com/jamOne-/kiwi-zero/TFPredictor"
	"github.com/jamOne-/kiwi-zero/game"
	"github.com/jamOne-/kiwi-zero/utils"
)

type Distribution []float32
type GameToDistributionFn func(game game.Game) Distribution

type PolicyPlayer struct {
	distributionFn GameToDistributionFn
}

func NewPolicyPlayer(distributionFn GameToDistributionFn) *PolicyPlayer {
	return &PolicyPlayer{distributionFn}
}

func (player *PolicyPlayer) SelectMove(game game.Game) game.Move {
	distribution := player.distributionFn(game)
	moves := game.GetPossibleMoves()

	movesValues := make([]float32, len(moves))
	for i, move := range moves {
		movesValues[i] = distribution[move+1] // TODO. +1 because of pass move
	}

	valuesSum := utils.SumFloats32(movesValues)
	scaleFactor := float32(1.0) / valuesSum
	for i, _ := range movesValues {
		movesValues[i] *= scaleFactor
	}

	moveIndex := utils.RandomFromDistribution(movesValues)
	return moves[moveIndex]
}

func GameToDistributionFnFromTfPredictor(
	gameToFeaturesFn game.GameToFeaturesFn,
	tfpredictor *tfpredictor.TFPredictor,
) GameToDistributionFn {
	return func(game game.Game) Distribution {
		features := gameToFeaturesFn(game)
		return tfpredictor.PredictPolicy(features)
	}
}
