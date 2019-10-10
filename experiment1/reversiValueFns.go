package main

import (
	"fmt"
	"io"
	"os"

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
	blackCount, whiteCount := 0, 0

	for i, field := range reversiGame.Board {
		features.SetVec(i, float64(field))

		if field == game.BLACK {
			blackCount += 1
		} else if field == game.WHITE {
			whiteCount += 1
		}
	}

	// countFeature := float64(blackCount) / float64(blackCount+whiteCount)
	// if whiteCount > blackCount {
	// 	countFeature = -float64(whiteCount) / float64(blackCount+whiteCount)
	// }

	countFeature := mat.Sum(features)
	features.SetVec(64, countFeature)

	currentPlayer := reversiGame.GetCurrentPlayerColor()
	reversiGame.Turn = game.BLACK
	blackMobility := len(reversiGame.GetPossibleMoves())
	reversiGame.Turn = game.WHITE
	whiteMobility := len(reversiGame.GetPossibleMoves())
	reversiGame.Turn = currentPlayer

	mobilityFeature := float64(blackMobility - whiteMobility)
	// mobilityFeature := float64(blackMobility) / float64(blackMobility+whiteMobility)
	// if whiteMobility > blackMobility {
	// 	mobilityFeature = -float64(whiteMobility) / float64(blackMobility+whiteMobility)
	// }

	features.SetVec(65, mobilityFeature)

	return features
}

func SaveWeightsToFile(weights *mat.VecDense, fileName string) {
	file, _ := os.Create(fileName)
	defer file.Close()

	for _, weight := range weights.RawVector().Data {
		fmt.Fprintf(file, "%f ", weight)
	}

	fmt.Fprintf(file, "\n")
}

func LoadWeightsFromFile(fileName string) *mat.VecDense {
	weights := []float64{}
	file, _ := os.Open(fileName)
	defer file.Close()

	aux := 0.0
	for {
		_, err := fmt.Fscan(file, &aux)

		if err == io.EOF {
			break
		}

		weights = append(weights, aux)
	}

	return mat.NewVecDense(len(weights), weights)
}

var PREVIOUS_WEIGHTS = []float64{20, -3, 11, 8, 8, 11, -3, 20,
	-3, -7, -4, 1, 1, -4, -7, -3,
	11, -4, 2, 2, 2, 2, -4, 11,
	8, 1, 2, -3, -3, 2, 1, 8,
	8, 1, 2, -3, -3, 2, 1, 8,
	11, -4, 2, 2, 2, 2, -4, 11,
	-3, -7, -4, 1, 1, -4, -7, -3,
	20, -3, 11, 8, 8, 11, -3, 20, 0, 0}
