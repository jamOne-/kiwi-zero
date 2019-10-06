package main

import (
	"gonum.org/v1/gonum/mat"

	"github.com/jamOne-/kiwi-zero/game"
	"github.com/jamOne-/kiwi-zero/minMaxPlayer"
	"github.com/jamOne-/kiwi-zero/reversi"
	"github.com/jamOne-/kiwi-zero/utils"
)

func createWeightedReversiFn(weights *mat.VecDense) minMaxPlayer.ValueFn {
	return func(g game.Game) float64 {
		reversiGame := g.(*reversi.ReversiGame) // nieładnie, ale brak generyków to jest jakiś dramat
		board := utils.Int8SliceToVecDense(reversiGame.Board)
		totalScore := 0.0

		countWeight := weights.AtVec(64)
		countDifference := mat.Sum(board)
		totalScore += countWeight * countDifference

		weightedPositions := mat.Dot(board, weights.SliceVec(0, 64))
		totalScore += weightedPositions

		mobilityWeight := weights.AtVec(64)
		currentPlayer := reversiGame.GetCurrentPlayerColor()
		reversiGame.Turn = game.BLACK
		blackMobility := len(reversiGame.GetPossibleMoves())
		reversiGame.Turn = game.WHITE
		whiteMobility := len(reversiGame.GetPossibleMoves())
		reversiGame.Turn = currentPlayer
		totalScore += mobilityWeight * float64(blackMobility-whiteMobility)

		return totalScore
	}
}

func getInitialWeights() *mat.VecDense {
	weightsLength := 8*8 + 2 // fields + countWeight + mobilityWeight
	weights := make([]float64, weightsLength)

	for i := 0; i < weightsLength; i++ {
		weights[i] = 1
	}

	return mat.NewVecDense(weightsLength, weights)
}
