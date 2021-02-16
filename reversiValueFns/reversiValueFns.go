package reversiValueFns

import (
	"github.com/jamOne-/kiwi-zero/game"
	"github.com/jamOne-/kiwi-zero/predictor"
	"github.com/jamOne-/kiwi-zero/reversi"
	"github.com/jamOne-/kiwi-zero/utils"
)

// var NUMBER_OF_FEATURES = 8*8 + 2 // fields + count + mobility

// func CreateWeightedReversiFn(gameToFeaturesFn game.GameToFeaturesFn, weights *mat.VecDense) minMaxPlayer.ValueFn {
// 	return func(g game.Game) float64 {
// 		features := gameToFeaturesFn(g)
// 		totalScore := mat.Dot(features, weights)
// 		afterSigmoid := 1.0 / (1 + math.Exp(-totalScore))

// 		finalValue := afterSigmoid*2.0 - 1.0

// 		return finalValue
// 	}
// }

// func GetInitialWeights() *mat.VecDense {
// 	return utils.CreateFilledVector(NUMBER_OF_FEATURES, 1)
// }

// func GetTriangleInitialWeights() *mat.VecDense {
// 	triangleNumberOfFeatures := 4 + 3 + 2 + 1 + 2
// 	return utils.CreateFilledVector(triangleNumberOfFeatures, 1)
// }

// func GetExtendedInitialWeights() *mat.VecDense {
// 	extendedNumberOfFeatures := 8*8 + 4 + 4
// 	return utils.CreateFilledVector(extendedNumberOfFeatures, 1)
// }

func ConvertReversiFnToGeneralFeatuersFn(reversiFn func(reversiGame *reversi.ReversiGame) game.Features) game.GameToFeaturesFn {
	return func(g game.Game) game.Features {
		reversiGame := g.(*reversi.ReversiGame)
		return reversiFn(reversiGame)
	}
}

func ReversiToOneHotBoard3(reversiGame *reversi.ReversiGame) game.Features {
	return reversiGame.OneHotBoard()
}

func ReversiToOneHotBoardMoves(game *reversi.ReversiGame) game.Features {
	emptyOneHot := []float32{1, 0, 0, 0}
	whiteOneHot := []float32{0, 1, 0, 0}
	blackOneHot := []float32{0, 0, 1, 0}
	features := make([][][]float32, reversi.BOARD_SIZE)

	for row := int8(0); row < reversi.BOARD_SIZE; row++ {
		features[row] = make([][]float32, reversi.BOARD_SIZE)

		for col := int8(0); col < reversi.BOARD_SIZE; col++ {
			field := game.Board[reversi.YXToField(row, col)]

			oneHotField := emptyOneHot
			if field == reversi.WHITE {
				oneHotField = whiteOneHot
			} else if field == reversi.BLACK {
				oneHotField = blackOneHot
			}

			features[row][col] = oneHotField
		}
	}

	possibleMoves := game.GetPossibleMoves()
	for _, move := range possibleMoves {
		if move == -1 {
			continue
		}

		y, x := reversi.GetYX(move)
		features[y][x][3] = 1
	}

	return features
}

func ReversiToOneHotBoardPaddedMoves(game *reversi.ReversiGame) game.Features {
	emptyOneHot := []float32{1, 0, 0, 0, 0}
	borderOneHot := []float32{0, 1, 0, 0, 0}
	whiteOneHot := []float32{0, 0, 1, 0, 0}
	blackOneHot := []float32{0, 0, 0, 1, 0}

	boardSize := int8(reversi.BOARD_SIZE + 2)
	features := make([][][]float32, boardSize)
	features[0] = make([][]float32, boardSize)
	features[boardSize-1] = make([][]float32, boardSize)

	for col := int8(0); col < boardSize; col++ {
		features[0][col] = borderOneHot
		features[boardSize-1][col] = borderOneHot
	}

	for row := int8(0); row < reversi.BOARD_SIZE; row++ {
		features[row+1] = make([][]float32, boardSize)
		features[row+1][0] = borderOneHot
		features[row+1][boardSize-1] = borderOneHot

		for col := int8(0); col < reversi.BOARD_SIZE; col++ {
			field := game.Board[reversi.YXToField(row, col)]

			oneHotField := emptyOneHot
			if field == reversi.WHITE {
				oneHotField = whiteOneHot
			} else if field == reversi.BLACK {
				oneHotField = blackOneHot
			}

			features[row+1][col+1] = oneHotField
		}
	}

	possibleMoves := game.GetPossibleMoves()
	for _, move := range possibleMoves {
		if move == -1 {
			continue
		}

		y, x := reversi.GetYX(move)
		features[y+1][x+1][4] = 1
	}

	return features
}

func CreateMinMaxValueFn(gameToFeaturesFn game.GameToFeaturesFn, predictor predictor.Predictor) game.ValueFn {
	return func(g game.Game) float64 {
		features := gameToFeaturesFn(g)
		prediction := predictor.PredictValue(features)

		return float64(prediction)*2.0 - 1.0
	}
}

func ReversiToFeatures(reversiGame *reversi.ReversiGame) game.Features {
	numberOfFeatures := 8 * 8
	features := make([][][]float32, numberOfFeatures)
	for row := 0; row < numberOfFeatures; row++ {
		features[row] = make([][]float32, 1)
		features[row][0] = make([]float32, 1)
	}

	for i, field := range reversiGame.Board {
		features[i][0][0] = float32(field)
	}

	return features
}

// func ReversiToFeaturesTriangle(reversiGame *reversi.ReversiGame) *mat.VecDense {
// 	triangleNumberOfFeatures := 4 + 3 + 2 + 1 + 2
// 	features := mat.NewVecDense(triangleNumberOfFeatures, nil)
// 	board := reversiGame.Board

// 	features.SetVec(0, float64(board[0]+board[7]+board[56]+board[63]))
// 	features.SetVec(1, float64(board[1]+board[6]+board[8]+board[15]+board[48]+board[57]+board[55]+board[62]))
// 	features.SetVec(2, float64(board[2]+board[5]+board[16]+board[23]+board[40]+board[47]+board[58]+board[61]))
// 	features.SetVec(3, float64(board[3]+board[4]+board[24]+board[32]+board[31]+board[39]+board[59]+board[60]))
// 	features.SetVec(4, float64(board[9]+board[14]+board[49]+board[54]))
// 	features.SetVec(5, float64(board[10]+board[13]+board[17]+board[22]+board[41]+board[46]+board[50]+board[53]))
// 	features.SetVec(6, float64(board[11]+board[12]+board[25]+board[33]+board[30]+board[38]+board[51]+board[52]))
// 	features.SetVec(7, float64(board[18]+board[21]+board[45]+board[42]))
// 	features.SetVec(8, float64(board[19]+board[20]+board[26]+board[29]+board[34]+board[37]+board[43]+board[44]))
// 	features.SetVec(9, float64(board[27]+board[28]+board[35]+board[36]))

// 	countFeature := mat.Sum(features)
// 	features.SetVec(triangleNumberOfFeatures-2, countFeature)

// 	currentPlayer := reversiGame.GetCurrentPlayerColor()
// 	reversiGame.Turn = game.BLACK
// 	blackMobility := len(reversiGame.GetPossibleMoves())
// 	reversiGame.Turn = game.WHITE
// 	whiteMobility := len(reversiGame.GetPossibleMoves())
// 	reversiGame.Turn = currentPlayer

// 	mobilityFeature := float64(blackMobility - whiteMobility)
// 	features.SetVec(triangleNumberOfFeatures-1, mobilityFeature)

// 	return features
// }

func ReversiToFeaturesExtended(reversiGame *reversi.ReversiGame) game.Features {
	numberOfFeaturesExtended := 8*8 + 4 + 4
	features := make([][][]float32, numberOfFeaturesExtended)
	for row := 0; row < numberOfFeaturesExtended; row++ {
		features[row] = make([][]float32, 1)
		features[row][0] = make([]float32, 1)
	}

	blackCount, whiteCount := 0.0, 0.0

	for i, field := range reversiGame.Board {
		features[i][0][0] = float32(field)

		if field == game.BLACK {
			blackCount += 1.0
		} else if field == game.WHITE {
			whiteCount += 1.0
		}
	}

	countDifference := blackCount - whiteCount
	countQuotient := calculateQuotient(blackCount, whiteCount)
	features[64][0][0] = float32(blackCount)
	features[65][0][0] = float32(whiteCount)
	features[66][0][0] = float32(countDifference)
	features[67][0][0] = float32(countQuotient)

	currentPlayer := reversiGame.GetCurrentPlayerColor()
	reversiGame.Turn = game.BLACK
	blackMobility := float64(len(reversiGame.GetPossibleMoves()))
	reversiGame.Turn = game.WHITE
	whiteMobility := float64(len(reversiGame.GetPossibleMoves()))
	reversiGame.Turn = currentPlayer

	mobilityDifference := blackMobility - whiteMobility
	mobilityQuotient := calculateQuotient(blackMobility, whiteMobility)

	features[68][0][0] = float32(blackMobility)
	features[69][0][0] = float32(whiteMobility)
	features[70][0][0] = float32(mobilityDifference)
	features[71][0][0] = float32(mobilityQuotient)

	return features
}

func calculateQuotient(a float64, b float64) float64 {
	if a > b {
		return a / (a + b)
	} else {
		return -b / (a + b)
	}
}

// func SaveWeightsToFile(weights *mat.VecDense, fileName string) {
// 	file, _ := os.Create(fileName)
// 	defer file.Close()

// 	for _, weight := range weights.RawVector().Data {
// 		fmt.Fprintf(file, "%f ", weight)
// 	}

// 	fmt.Fprintf(file, "\n")
// }

// func LoadWeightsFromFile(fileName string) *mat.VecDense {
// 	weights := []float64{}
// 	file, _ := os.Open(fileName)
// 	defer file.Close()

// 	aux := 0.0
// 	for {
// 		_, err := fmt.Fscan(file, &aux)

// 		if err == io.EOF {
// 			break
// 		}

// 		weights = append(weights, aux)
// 	}

// 	return mat.NewVecDense(len(weights), weights)
// }

// var OLD_MINMAX_WEIGHTS = []float64{20, -3, 11, 8, 8, 11, -3, 20,
// 	-3, -7, -4, 1, 1, -4, -7, -3,
// 	11, -4, 2, 2, 2, 2, -4, 11,
// 	8, 1, 2, -3, -3, 2, 1, 8,
// 	8, 1, 2, -3, -3, 2, 1, 8,
// 	11, -4, 2, 2, 2, 2, -4, 11,
// 	-3, -7, -4, 1, 1, -4, -7, -3,
// 	20, -3, 11, 8, 8, 11, -3, 20, 0, 0}

func ReversiToOneHotBoardMovesTurn(game *reversi.ReversiGame) game.Features {
	emptyOneHot := []float32{1, 0, 0, 0, 0, 0}
	whiteOneHot := []float32{0, 1, 0, 0, 0, 0}
	blackOneHot := []float32{0, 0, 1, 0, 0, 0}
	features := make([][][]float32, reversi.BOARD_SIZE)

	turnDim := 4 + utils.BoolToInt(game.Turn == reversi.WHITE)
	emptyOneHot[turnDim] = 1
	whiteOneHot[turnDim] = 1
	blackOneHot[turnDim] = 1

	for row := int8(0); row < reversi.BOARD_SIZE; row++ {
		features[row] = make([][]float32, reversi.BOARD_SIZE)

		for col := int8(0); col < reversi.BOARD_SIZE; col++ {
			field := game.Board[reversi.YXToField(row, col)]

			oneHotField := emptyOneHot
			if field == reversi.WHITE {
				oneHotField = whiteOneHot
			} else if field == reversi.BLACK {
				oneHotField = blackOneHot
			}

			features[row][col] = oneHotField
		}
	}

	possibleMoves := game.GetPossibleMoves()
	for _, move := range possibleMoves {
		if move == -1 {
			continue
		}

		y, x := reversi.GetYX(move)
		features[y][x][3] = 1
	}

	return features
}

func ReversiToBoardTurn(game *reversi.ReversiGame) game.Features {
	emptyField := []float32{0, 1, 0, 0, 0}
	whiteField := []float32{1, 0, 0, 0, 0}
	blackField := []float32{0, 0, 1, 0, 0}
	oneHotBoard := make([][][]float32, reversi.BOARD_SIZE)

	turnDim := 3 + utils.BoolToInt(game.Turn == reversi.WHITE)
	emptyField[turnDim] = 1
	whiteField[turnDim] = 1
	blackField[turnDim] = 1

	for row := int8(0); row < reversi.BOARD_SIZE; row++ {
		oneHotBoard[row] = make([][]float32, reversi.BOARD_SIZE)

		for col := int8(0); col < reversi.BOARD_SIZE; col++ {
			field := game.Board[row*reversi.BOARD_SIZE+col]

			oneHotField := emptyField
			if field == reversi.WHITE {
				oneHotField = whiteField
			} else if field == reversi.BLACK {
				oneHotField = blackField
			}

			oneHotBoard[row][col] = oneHotField
		}
	}

	return oneHotBoard
}
