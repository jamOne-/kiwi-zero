package reversi

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jamOne-/kiwi-zero/game"
)

type ReversiGame struct {
	Turn    game.PlayerColor
	Board   []game.Field
	History []*ReversiGameHistoryItem
}

type ReversiGameHistoryItem struct {
	Move   game.Move
	Killed []game.Field
}

const BOARD_SIZE = 8
const TOTAL_SIZE = BOARD_SIZE * BOARD_SIZE
const EMPTY = game.Field(0)
const WHITE = game.Field(-1)
const BLACK = game.Field(1)
const PASS_MOVE = -1

func GetYX(field game.Field) (int8, int8) {
	return field / BOARD_SIZE, field % BOARD_SIZE
}

func YXToField(y int8, x int8) game.Field {
	return y*BOARD_SIZE + x
}

func NewReversiGame() *ReversiGame {
	turn := BLACK
	board := make([]game.Field, TOTAL_SIZE)
	history := make([]*ReversiGameHistoryItem, 0, TOTAL_SIZE)
	board[YXToField(3, 3)], board[YXToField(4, 4)] = WHITE, WHITE
	board[YXToField(3, 4)], board[YXToField(4, 3)] = BLACK, BLACK

	return &ReversiGame{turn, board, history}
}

func (reversiGame *ReversiGame) Copy() game.Game {
	turn := reversiGame.Turn
	board := make([]game.Field, TOTAL_SIZE)
	history := make([]*ReversiGameHistoryItem, len(reversiGame.History), cap(reversiGame.History))
	copy(board, reversiGame.Board)
	copy(history, reversiGame.History)

	return &ReversiGame{turn, board, history}
}

func (reversiGame *ReversiGame) MakeMove(move game.Move) (bool, game.PlayerColor) {
	currentPlayer := reversiGame.Turn
	killed := make([]game.Field, 0)

	if move != PASS_MOVE {
		reversiGame.Board[move] = currentPlayer
		killed = getKilledPawns(reversiGame.Board, move, currentPlayer)

		for _, field := range killed {
			reversiGame.Board[field] = currentPlayer
		}
	}

	historyItem := &ReversiGameHistoryItem{move, killed}
	reversiGame.History = append(reversiGame.History, historyItem)
	reversiGame.Turn *= -1

	return reversiGame.IsGameFinished()
}

func (reversiGame *ReversiGame) GetPossibleMoves() []game.Move {
	result := make([]game.Field, 0)

	for field := int8(0); field < TOTAL_SIZE; field++ {
		if reversiGame.Board[field] == EMPTY && len(getKilledPawns(reversiGame.Board, field, reversiGame.Turn)) > 0 {
			result = append(result, field)
		}
	}

	if len(result) == 0 {
		result = append(result, PASS_MOVE)
	}

	return result
}

func (reversiGame *ReversiGame) UndoLastMove() {
	historyLen := len(reversiGame.History)
	historyItem := reversiGame.History[historyLen-1]
	move, killed := historyItem.Move, historyItem.Killed

	if move != PASS_MOVE {
		reversiGame.Board[move] = EMPTY

		for _, field := range killed {
			reversiGame.Board[field] *= -1
		}
	}

	reversiGame.History = reversiGame.History[:historyLen-1]
	reversiGame.Turn *= -1
}

func (reversiGame *ReversiGame) GetCurrentPlayerColor() game.PlayerColor {
	return reversiGame.Turn
}

func (reversiGame *ReversiGame) IsGameFinished() (bool, game.PlayerColor) {
	turns := len(reversiGame.History)
	if turns < 2 || reversiGame.History[turns-2].Move != PASS_MOVE || reversiGame.History[turns-1].Move != PASS_MOVE {
		currentPlayerMoves := reversiGame.GetPossibleMoves()

		if currentPlayerMoves[0] != PASS_MOVE {
			return false, EMPTY
		}

		reversiGame.Turn *= -1
		nextPlayerMoves := reversiGame.GetPossibleMoves()
		reversiGame.Turn *= -1

		if nextPlayerMoves[0] != PASS_MOVE {
			return false, EMPTY
		}
	}

	// evaluating winner

	blacks, whites := reversiGame.CountPawns()
	winner := EMPTY

	if blacks > whites {
		winner = BLACK
	} else if whites > blacks {
		winner = WHITE
	}

	return true, winner
}

func (reversiGame *ReversiGame) CountPawns() (int8, int8) {
	black, white := int8(0), int8(0)

	for field := 0; field < TOTAL_SIZE; field++ {
		pawn := reversiGame.Board[field]

		if pawn == BLACK {
			black++
		} else if pawn == WHITE {
			white++
		}
	}

	return black, white
}

func getKilledPawns(board []game.Field, start game.Field, player game.PlayerColor) []game.Field {
	opponent := player * -1
	result := make([]game.Field, 0)
	startY, startX := GetYX(start)
	deltas := []int8{-1, 0, 1}

	for _, dy := range deltas {
		for _, dx := range deltas {
			if dx == 0 && dy == 0 {
				continue
			}

			candidates := make([]game.Field, 0)

			for y, x := startY+dy, startX+dx; x >= 0 && x < BOARD_SIZE && y >= 0 && y < BOARD_SIZE; y, x = y+dy, x+dx {
				field := YXToField(y, x)
				pawn := board[field]

				if pawn == opponent {
					candidates = append(candidates, field)
				} else if pawn == player {
					result = append(result, candidates...)
					break
				} else {
					break
				}
			}
		}
	}

	return result
}

func (game *ReversiGame) DrawBoard() {
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

	for _, move := range game.GetPossibleMoves() {
		if move != -1 {
			board = board[:move] + "+" + board[move+1:]
		}
	}

	fmt.Println("   0 1 2 3 4 5 6 7")

	for row := 0; row < BOARD_SIZE; row++ {
		start := row * BOARD_SIZE

		if row < 2 {
			fmt.Printf(" %d ", row*8)
		} else {
			fmt.Printf("%d ", row*8)
		}

		line := board[start : start+BOARD_SIZE]

		for _, field := range line {
			fmt.Printf("%c ", field)
		}

		fmt.Printf("\n")
	}
}

func (game *ReversiGame) SerializeBoard(flipColors bool) string {
	stringsBoard := make([]string, len(game.Board))
	factor := 1
	if flipColors {
		factor = -1
	}

	for i, field := range game.Board {
		stringsBoard[i] = strconv.Itoa(int(field) * factor)
	}

	board := strings.Join(stringsBoard, " ")
	return fmt.Sprintf("%s %d %d", board, game.GetTurnNumber(), game.Turn)
}

func (game *ReversiGame) OneHotBoard() [][][]float32 {
	oneHotBoard := make([][][]float32, BOARD_SIZE)

	for row := int8(0); row < BOARD_SIZE; row++ {
		oneHotBoard[row] = make([][]float32, BOARD_SIZE)

		for col := int8(0); col < BOARD_SIZE; col++ {
			field := game.Board[YXToField(row, col)]

			oneHotField := []float32{0, 1, 0}
			if field == -1 {
				oneHotField = []float32{1, 0, 0}
			} else if field == 1 {
				oneHotField = []float32{0, 0, 1}
			}

			oneHotBoard[row][col] = oneHotField
		}
	}

	return oneHotBoard
}

func (game *ReversiGame) GetTurnNumber() int {
	return len(game.History)
}

func (game *ReversiGame) GetMaxPossibleMoves() int {
	return TOTAL_SIZE + 1
}

func (game *ReversiGame) EncodeMoveToPolicy(move game.Move) []float32 {
	policy := make([]float32, game.GetMaxPossibleMoves())
	policy[move+1] = 1 // move + 1, because pass is -1
	return policy
}
