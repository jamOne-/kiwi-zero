package gomoku

import (
	"github.com/jamOne-/kiwi-zero/game"
	"github.com/jamOne-/kiwi-zero/utils"
)

func ConvertGomokuFnToGeneralFeatuersFn(gomokuFn func(gomoku *GomokuGame) game.Features) game.GameToFeaturesFn {
	return func(g game.Game) game.Features {
		gomokuGame := g.(*GomokuGame)
		return gomokuFn(gomokuGame)
	}
}

func GomokuToBoard1(gomoku *GomokuGame) game.Features {
	features := make([][][]float32, TOTAL_SIZE)
	for row := 0; row < TOTAL_SIZE; row++ {
		features[row] = make([][]float32, 1)
		features[row][0] = make([]float32, 1)
	}

	for i, field := range gomoku.Board {
		features[i][0][0] = float32(field)
	}

	return features
}

func GomokuToBoardTurn(gomoku *GomokuGame) game.Features {
	emptyField := []float32{0, 1, 0, 0, 0}
	whiteField := []float32{1, 0, 0, 0, 0}
	blackField := []float32{0, 0, 1, 0, 0}
	oneHotBoard := make([][][]float32, HEIGHT)

	turnDim := 3 + utils.BoolToInt(gomoku.Turn == WHITE)
	emptyField[turnDim] = 1
	whiteField[turnDim] = 1
	blackField[turnDim] = 1

	for row := int8(0); row < HEIGHT; row++ {
		oneHotBoard[row] = make([][]float32, WIDTH)

		for col := int8(0); col < WIDTH; col++ {
			field := gomoku.Board[row*WIDTH+col]

			oneHotField := emptyField
			if field == WHITE {
				oneHotField = whiteField
			} else if field == BLACK {
				oneHotField = blackField
			}

			oneHotBoard[row][col] = oneHotField
		}
	}

	return oneHotBoard
}

func GomokuToBoard1Turn(gomoku *GomokuGame) game.Features {
	numberOfFeatures := WIDTH*HEIGHT + 1
	features := make([][][]float32, numberOfFeatures)
	for row := 0; row < numberOfFeatures; row++ {
		features[row] = make([][]float32, 1)
		features[row][0] = make([]float32, 1)
	}

	for i, field := range gomoku.Board {
		features[i][0][0] = float32(field)
	}

	features[numberOfFeatures-1][0][0] = float32(gomoku.GetCurrentPlayerColor())

	return features
}
