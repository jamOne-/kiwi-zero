package main

import (
	"github.com/jamOne-/kiwi-zero/game"
	"github.com/jamOne-/kiwi-zero/minMaxPlayer"
	"github.com/jamOne-/kiwi-zero/reversi"
)

func createWeightedReversiFn(weights []float64) minMaxPlayer.ValueFn {
	return func(g game.Game) float64 {
		reversiGame := g.(*reversi.ReversiGame) // nieładnie, ale brak generyków to jest jakiś dramat
		blacks, whites := 0, 0
		totalScore := 0.0

		for i, pawn := range reversiGame.Board {
			if pawn == reversi.BLACK {
				blacks += 1
				totalScore += weights[i]
			} else if pawn == reversi.WHITE {
				whites += 1
				totalScore -= weights[i]
			}
		}

		countWeight := weights[64]
		totalScore += countWeight * float64(blacks-whites)

		mobilityWeight := weights[65]
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

func getInitialWeights() []float64 {
	weightsLength := 8*8 + 2 // fields + countWeight + mobilityWeight
	weights := make([]float64, weightsLength)

	for i := 0; i < weightsLength; i++ {
		weights[i] = 1
	}

	return weights
}
