package gomoku

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/jamOne-/kiwi-zero/game"
	"github.com/jamOne-/kiwi-zero/utils"
)

const WIDTH = 8
const HEIGHT = 8
const TOTAL_SIZE = WIDTH * HEIGHT
const EMPTY = game.Field(0)
const WHITE = game.Field(-1)
const BLACK = game.Field(1)

type GomokuGame struct {
	Board   []game.Field
	History []int8
	Turn    game.PlayerColor
}

func NewGomokuGame() *GomokuGame {
	turn := game.BLACK
	board := make([]game.Field, TOTAL_SIZE)
	history := make([]int8, 0, TOTAL_SIZE)

	return &GomokuGame{board, history, turn}
}

func (gomoku *GomokuGame) Copy() game.Game {
	turn := gomoku.Turn
	board := make([]game.Field, TOTAL_SIZE)
	history := make([]int8, len(gomoku.History), cap(gomoku.History))
	copy(board, gomoku.Board)
	copy(history, gomoku.History)

	return &GomokuGame{board, history, turn}
}

func (gomoku *GomokuGame) MakeMove(move game.Move) (bool, game.PlayerColor) {
	gomoku.Board[move] = gomoku.Turn
	gomoku.History = append(gomoku.History, move)
	gomoku.Turn *= -1

	return gomoku.IsGameFinished()
}

func (gomoku *GomokuGame) IsGameFinished() (bool, game.PlayerColor) {
	for y := HEIGHT - 1; y >= 0; y-- {
		for x := 0; x < WIDTH; x++ {
			field := gomoku.Board[y*WIDTH+x]

			if field == EMPTY {
				continue
			}

			if gomoku.countInDirection(x, y, -1, 0)+gomoku.countInDirection(x, y, 1, 0) == 4 || // -
				gomoku.countInDirection(x, y, 0, -1)+gomoku.countInDirection(x, y, 0, 1) == 4 || // |
				gomoku.countInDirection(x, y, -1, -1)+gomoku.countInDirection(x, y, 1, 1) == 4 || // \
				gomoku.countInDirection(x, y, -1, 1)+gomoku.countInDirection(x, y, 1, -1) == 4 { // /
				return true, field
			}
		}
	}

	return gomoku.GetTurnNumber() == TOTAL_SIZE, EMPTY
}

func (gomoku *GomokuGame) countInDirection(x int, y int, dx int, dy int) int {
	field := gomoku.Board[y*WIDTH+x]
	count := 0

	x += dx
	y += dy

	for x >= 0 && x < WIDTH && y >= 0 && y < HEIGHT && gomoku.Board[y*WIDTH+x] == field {
		count += 1
		x += dx
		y += dy
	}

	return count
}

func (gomoku *GomokuGame) GetPossibleMoves() []game.Move {
	result := make([]game.Field, 0, TOTAL_SIZE)

	for y := 0; y < HEIGHT; y++ {
		for x := 0; x < WIDTH; x++ {
			move := game.Move(y*WIDTH + x)

			if gomoku.Board[move] == EMPTY {
				result = append(result, move)
			}
		}
	}

	return result
}

func (gomoku *GomokuGame) UndoLastMove() {
	historyLen := len(gomoku.History)
	move := gomoku.History[historyLen-1]

	gomoku.Board[move] = EMPTY
	gomoku.History = gomoku.History[:historyLen-1]
	gomoku.Turn *= -1
}

func (gomoku *GomokuGame) GetCurrentPlayerColor() game.PlayerColor {
	return gomoku.Turn
}

func (gomoku *GomokuGame) SerializeBoard(flipColors bool) string {
	stringsBoard := make([]string, len(gomoku.Board))
	factor := 1
	if flipColors {
		factor = -1
	}

	for i, field := range gomoku.Board {
		stringsBoard[i] = strconv.Itoa(int(field) * factor)
	}

	board := strings.Join(stringsBoard, " ")
	return fmt.Sprintf("%s %d", board, gomoku.Turn)
}

func (gomoku *GomokuGame) GetTurnNumber() int {
	return len(gomoku.History)
}

func (gomoku *GomokuGame) GetMaxPossibleMoves() int {
	return TOTAL_SIZE + 1
}

func (gomoku *GomokuGame) EncodeMoveToPolicy(move game.Move) []float32 {
	policy := make([]float32, gomoku.GetMaxPossibleMoves())
	policy[move+1] = 1 // move + 1, because pass is -1
	return policy
}

func (gomoku *GomokuGame) FlipColors() {
	gomoku.Turn *= -1

	for i, color := range gomoku.Board {
		gomoku.Board[i] = color * -1
	}
}

func (gomoku *GomokuGame) RandomPositionTransformation() {
	TRANSFORMATIONS := 4 + 4
	transformation := rand.Intn(TRANSFORMATIONS)

	switch transformation {
	case 0:
		return
	case 1:
		fallthrough
	case 2:
		fallthrough
	case 3:
		utils.RotateSquareVector(gomoku.Board, transformation)

	case 4:
		utils.PerformSymmetryVector1(gomoku.Board)
	case 5:
		utils.PerformSymmetryVector2(WIDTH, HEIGHT, gomoku.Board)
	case 6:
		utils.PerformSymmetryVector3(gomoku.Board)
	case 7:
		utils.PerformSymmetryVector4(gomoku.Board)
	}
}

func (gomoku *GomokuGame) OneHotBoard() [][][]float32 {
	emptyField := []float32{0, 1, 0}
	whiteField := []float32{1, 0, 0}
	blackField := []float32{0, 0, 1}
	oneHotBoard := make([][][]float32, HEIGHT)

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

func (gomoku *GomokuGame) DrawBoard() {
	board := ""

	for field := 0; field < TOTAL_SIZE; field++ {
		switch gomoku.Board[field] {
		case EMPTY:
			board += "."
		case BLACK:
			board += "x"
		case WHITE:
			board += "o"
		}
	}

	fmt.Println("   0 1 2 3 4 5 6 7")

	for row := 0; row < HEIGHT; row++ {
		start := row * WIDTH

		if row < 2 {
			fmt.Printf(" %d ", row*8)
		} else {
			fmt.Printf("%d ", row*8)
		}

		line := board[start : start+WIDTH]

		for _, field := range line {
			fmt.Printf("%c ", field)
		}

		fmt.Printf("\n")
	}
}
