package main

import (
	"gonum.org/v1/gonum/mat"

	"github.com/jamOne-/kiwi-zero/game"
	"github.com/jamOne-/kiwi-zero/minMaxPlayer"
	"github.com/jamOne-/kiwi-zero/reversi"
)

var NUMBER_OF_FEATURES = 8*8 + 2 // fields + count + mobility

func createWeightedReversiFn(weights *mat.VecDense) minMaxPlayer.ValueFn {
	return func(g game.Game) float64 {
		reversiGame := g.(*reversi.ReversiGame) // nieładnie, ale brak generyków to jest jakiś dramat

		features := ReversiToFeatures(reversiGame)
		totalScore := mat.Dot(features, weights)

		return totalScore
	}
}

func getInitialWeights() *mat.VecDense {
	weights := make([]float64, NUMBER_OF_FEATURES)

	for i := 0; i < NUMBER_OF_FEATURES; i++ {
		weights[i] = 1
	}

	return mat.NewVecDense(NUMBER_OF_FEATURES, weights)
}

func ReversiToFeatures(reversiGame *reversi.ReversiGame) *mat.VecDense {
	features := mat.NewVecDense(NUMBER_OF_FEATURES, nil)

	for i, field := range reversiGame.Board {
		features.SetVec(i, float64(field))
	}

	countDifference := mat.Sum(features)
	features.SetVec(64, countDifference)

	currentPlayer := reversiGame.GetCurrentPlayerColor()
	reversiGame.Turn = game.BLACK
	blackMobility := len(reversiGame.GetPossibleMoves())
	reversiGame.Turn = game.WHITE
	whiteMobility := len(reversiGame.GetPossibleMoves())
	reversiGame.Turn = currentPlayer
	features.SetVec(65, float64(blackMobility-whiteMobility))

	return features
}
