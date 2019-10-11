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

type ReversiToFeaturesFn func(reversiGame *reversi.ReversiGame) *mat.VecDense

var NUMBER_OF_FEATURES = 8*8 + 2 // fields + count + mobility

func createWeightedReversiFn(reversiToFeaturesFn ReversiToFeaturesFn, weights *mat.VecDense) minMaxPlayer.ValueFn {
	return func(g game.Game) float64 {
		reversiGame := g.(*reversi.ReversiGame) // nieładnie, ale brak generyków to jest jakiś dramat

		features := reversiToFeaturesFn(reversiGame)
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
	features.SetVec(NUMBER_OF_FEATURES-2, countFeature)

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

	features.SetVec(NUMBER_OF_FEATURES-1, mobilityFeature)

	return features
}

func ReversiToFeaturesTriangle(reversiGame *reversi.ReversiGame) *mat.VecDense {
	triangleNumberOfFeatures := 4 + 3 + 2 + 1 + 2
	features := mat.NewVecDense(triangleNumberOfFeatures, nil)
	board := reversiGame.Board

	features.SetVec(0, float64(board[0]+board[7]+board[56]+board[63]))
	features.SetVec(1, float64(board[1]+board[6]+board[8]+board[15]+board[48]+board[57]+board[55]+board[62]))
	features.SetVec(2, float64(board[2]+board[5]+board[16]+board[23]+board[40]+board[47]+board[58]+board[61]))
	features.SetVec(3, float64(board[3]+board[4]+board[24]+board[32]+board[31]+board[39]+board[59]+board[60]))
	features.SetVec(4, float64(board[9]+board[14]+board[49]+board[54]))
	features.SetVec(5, float64(board[10]+board[13]+board[17]+board[22]+board[41]+board[46]+board[50]+board[53]))
	features.SetVec(6, float64(board[11]+board[12]+board[25]+board[33]+board[30]+board[38]+board[51]+board[52]))
	features.SetVec(7, float64(board[18]+board[21]+board[45]+board[42]))
	features.SetVec(8, float64(board[19]+board[20]+board[26]+board[29]+board[34]+board[37]+board[43]+board[44]))
	features.SetVec(9, float64(board[27]+board[28]+board[35]+board[36]))

	countFeature := mat.Sum(features)
	features.SetVec(triangleNumberOfFeatures-2, countFeature)

	currentPlayer := reversiGame.GetCurrentPlayerColor()
	reversiGame.Turn = game.BLACK
	blackMobility := len(reversiGame.GetPossibleMoves())
	reversiGame.Turn = game.WHITE
	whiteMobility := len(reversiGame.GetPossibleMoves())
	reversiGame.Turn = currentPlayer

	mobilityFeature := float64(blackMobility - whiteMobility)
	features.SetVec(triangleNumberOfFeatures-1, mobilityFeature)

	return features
}

func getTriangleInitialWeights() *mat.VecDense {
	triangleNumberOfFeatures := 4 + 3 + 2 + 1 + 2
	weights := make([]float64, triangleNumberOfFeatures)

	for i := 0; i < triangleNumberOfFeatures; i++ {
		weights[i] = 1
	}

	return mat.NewVecDense(triangleNumberOfFeatures, weights)
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

var OLD_MINMAX_WEIGHTS = []float64{20, -3, 11, 8, 8, 11, -3, 20,
	-3, -7, -4, 1, 1, -4, -7, -3,
	11, -4, 2, 2, 2, 2, -4, 11,
	8, 1, 2, -3, -3, 2, 1, 8,
	8, 1, 2, -3, -3, 2, 1, 8,
	11, -4, 2, 2, 2, 2, -4, 11,
	-3, -7, -4, 1, 1, -4, -7, -3,
	20, -3, 11, 8, 8, 11, -3, 20, 0, 0}
