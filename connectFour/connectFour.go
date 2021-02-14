package connectFour

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/jamOne-/kiwi-zero/game"
	"github.com/jamOne-/kiwi-zero/utils"
)

const WIDTH = 7
const HEIGHT = 6
const TOTAL_SIZE = WIDTH * HEIGHT
const EMPTY = game.Field(0)
const WHITE = game.Field(-1)
const BLACK = game.Field(1)

type ConnectFourGame struct {
	Board   []game.Field
	History []int8
	Turn    game.PlayerColor
}

func NewConnectFourGame() *ConnectFourGame {
	turn := game.BLACK
	board := make([]game.Field, TOTAL_SIZE)
	history := make([]int8, 0, TOTAL_SIZE)

	return &ConnectFourGame{board, history, turn}
}

func (connectFour *ConnectFourGame) Copy() game.Game {
	turn := connectFour.Turn
	board := make([]game.Field, TOTAL_SIZE)
	history := make([]int8, len(connectFour.History), cap(connectFour.History))
	copy(board, connectFour.Board)
	copy(history, connectFour.History)

	return &ConnectFourGame{board, history, turn}
}

func (connectFour *ConnectFourGame) MakeMove(move game.Move) (bool, game.PlayerColor) {
	y := HEIGHT - 1

	for connectFour.Board[y*WIDTH+int(move)] != EMPTY {
		y--
	}
	connectFour.Board[y*WIDTH+int(move)] = connectFour.Turn
	connectFour.History = append(connectFour.History, move)
	connectFour.Turn *= -1

	return connectFour.IsGameFinished()
}

func (connectFour *ConnectFourGame) IsGameFinished() (bool, game.PlayerColor) {
	for y := HEIGHT - 1; y >= 0; y-- {
		for x := 0; x < WIDTH; x++ {
			field := connectFour.Board[y*WIDTH+x]

			if field == EMPTY {
				continue
			}

			if connectFour.countInDirection(x, y, -1, 0) == 4 || // W
				connectFour.countInDirection(x, y, 0, -1) == 4 || // N
				connectFour.countInDirection(x, y, -1, -1) == 4 || // NW
				connectFour.countInDirection(x, y, 1, -1) == 4 { // NE
				return true, field
			}
		}
	}

	for x := 0; x < WIDTH; x++ {
		if connectFour.Board[x] == EMPTY {
			return false, EMPTY
		}
	}

	return true, EMPTY
}

func (connectFour *ConnectFourGame) countInDirection(x int, y int, dx int, dy int) int {
	field := connectFour.Board[y*WIDTH+x]
	count := 0

	for x >= 0 && x < WIDTH && y >= 0 && y < HEIGHT && connectFour.Board[y*WIDTH+x] == field {
		count += 1
		x += dx
		y += dy
	}

	return count
}

func (connectFour *ConnectFourGame) GetPossibleMoves() []game.Move {
	result := make([]game.Field, 0, WIDTH)

	for x := 0; x < WIDTH; x++ {
		if connectFour.Board[x] == EMPTY {
			result = append(result, game.Move(x))
		}
	}

	return result
}

func (connectFour *ConnectFourGame) UndoLastMove() {
	historyLen := len(connectFour.History)
	move := connectFour.History[historyLen-1]

	x := int(move)
	y := 0
	for connectFour.Board[y*WIDTH+x] == EMPTY {
		y += 1
	}
	connectFour.Board[y*WIDTH+x] = EMPTY
	connectFour.History = connectFour.History[:historyLen-1]
	connectFour.Turn *= -1
}

func (connectFour *ConnectFourGame) GetCurrentPlayerColor() game.PlayerColor {
	return connectFour.Turn
}

func (game *ConnectFourGame) SerializeBoard(flipColors bool) string {
	stringsBoard := make([]string, len(game.Board))
	factor := 1
	if flipColors {
		factor = -1
	}

	for i, field := range game.Board {
		stringsBoard[i] = strconv.Itoa(int(field) * factor)
	}

	board := strings.Join(stringsBoard, " ")
	return fmt.Sprintf("%s %d", board, game.Turn)
}

func (game *ConnectFourGame) GetTurnNumber() int {
	return len(game.History)
}

func (game *ConnectFourGame) GetMaxPossibleMoves() int {
	return WIDTH + 1
}

func (game *ConnectFourGame) EncodeMoveToPolicy(move game.Move) []float32 {
	policy := make([]float32, game.GetMaxPossibleMoves())
	policy[move+1] = 1 // move + 1, because pass is -1
	return policy
}

func (game *ConnectFourGame) FlipColors() {
	game.Turn *= -1

	for i, color := range game.Board {
		game.Board[i] = color * -1
	}
}

func (game *ConnectFourGame) RandomPositionTransformation(policy []float32) {
	TRANSFORMATIONS := 2
	transformation := rand.Intn(TRANSFORMATIONS)

	switch transformation {
	case 0:
		return
	case 1:
		utils.PerformSymmetryVector2Float32(WIDTH, HEIGHT, policy)
		if policy != nil {
			utils.PerformSymmetryVector2Float32(WIDTH, HEIGHT, policy)
		}
	}
}

func (game *ConnectFourGame) OneHotBoard() [][][]float32 {
	emptyField := []float32{0, 1, 0}
	whiteField := []float32{1, 0, 0}
	blackField := []float32{0, 0, 1}
	oneHotBoard := make([][][]float32, HEIGHT)

	for row := int8(0); row < HEIGHT; row++ {
		oneHotBoard[row] = make([][]float32, WIDTH)

		for col := int8(0); col < WIDTH; col++ {
			field := game.Board[row*WIDTH+col]

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

func (game *ConnectFourGame) DrawBoard() {
	board := ""

	for field := 0; field < TOTAL_SIZE; field++ {
		switch game.Board[field] {
		case EMPTY:
			board += "."
		case BLACK:
			board += "x"
		case WHITE:
			board += "o"
		}
	}

	fmt.Println("0 1 2 3 4 5 6")

	for row := 0; row < HEIGHT; row++ {
		start := row * WIDTH

		line := board[start : start+WIDTH]

		for _, field := range line {
			fmt.Printf("%c ", field)
		}

		fmt.Printf("\n")
	}

	fmt.Println("0 1 2 3 4 5 6")
}
