package connectFour

import (
	"github.com/jamOne-/kiwi-zero/game"
	"github.com/jamOne-/kiwi-zero/utils"
)

func ConvertConnect4FnToGeneralFeatuersFn(connect4Fn func(c4game *ConnectFourGame) game.Features) game.GameToFeaturesFn {
	return func(g game.Game) game.Features {
		connect4game := g.(*ConnectFourGame)
		return connect4Fn(connect4game)
	}
}

func Connect4ToBoard1(c4game *ConnectFourGame) game.Features {
	features := make([][][]float32, TOTAL_SIZE)
	for row := 0; row < TOTAL_SIZE; row++ {
		features[row] = make([][]float32, 1)
		features[row][0] = make([]float32, 1)
	}

	for i, field := range c4game.Board {
		features[i][0][0] = float32(field)
	}

	return features
}

func Connect4ToBoardTurn(c4game *ConnectFourGame) game.Features {
	emptyField := []float32{0, 1, 0, 0, 0}
	whiteField := []float32{1, 0, 0, 0, 0}
	blackField := []float32{0, 0, 1, 0, 0}
	oneHotBoard := make([][][]float32, HEIGHT)

	turnDim := 3 + utils.BoolToInt(c4game.Turn == WHITE)
	emptyField[turnDim] = 1
	whiteField[turnDim] = 1
	blackField[turnDim] = 1

	for row := int8(0); row < HEIGHT; row++ {
		oneHotBoard[row] = make([][]float32, WIDTH)

		for col := int8(0); col < WIDTH; col++ {
			field := c4game.Board[row*WIDTH+col]

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

func Connect4ToBoardMovesTurn(c4game *ConnectFourGame) game.Features {
	emptyField := []float32{0, 1, 0, 0, 0, 0}
	whiteField := []float32{1, 0, 0, 0, 0, 0}
	blackField := []float32{0, 0, 1, 0, 0, 0}
	oneHotBoard := make([][][]float32, HEIGHT)

	turnDim := 4 + utils.BoolToInt(c4game.Turn == WHITE)
	emptyField[turnDim] = 1
	whiteField[turnDim] = 1
	blackField[turnDim] = 1

	for row := int8(0); row < HEIGHT; row++ {
		oneHotBoard[row] = make([][]float32, WIDTH)

		for col := int8(0); col < WIDTH; col++ {
			field := c4game.Board[row*WIDTH+col]

			oneHotField := emptyField
			if field == WHITE {
				oneHotField = whiteField
			} else if field == BLACK {
				oneHotField = blackField
			}

			oneHotBoard[row][col] = oneHotField
		}
	}

	possiblex := c4game.GetPossibleMoves()
	for _, x := range possiblex {
		y := HEIGHT - 1
		for c4game.Board[y*WIDTH+int(x)] != EMPTY {
			y--
		}

		oneHotBoard[y][x][3] = 1
	}

	return oneHotBoard
}

func Connect4ToBoard1Turn(c4game *ConnectFourGame) game.Features {
	size := TOTAL_SIZE + 1
	features := make([][][]float32, size)
	for row := 0; row < size; row++ {
		features[row] = make([][]float32, 1)
		features[row][0] = make([]float32, 1)
	}

	for i, field := range c4game.Board {
		features[i][0][0] = float32(field)
	}

	features[size-1][0][0] = float32(c4game.Turn)

	return features
}
