package reversi

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jamOne-/kiwi-zero/game"
)

type Game struct {
	Turn    game.PlayerColor
	Board   []game.Field
	History []game.Move
}

const BOARD_SIZE = 8
const TOTAL_SIZE = BOARD_SIZE * BOARD_SIZE
const EMPTY = game.Field(0)
const WHITE = game.Field(-1)
const BLACK = game.Field(1)

func getYX(field game.Field) (int8, int8) {
	return field / BOARD_SIZE, field % BOARD_SIZE
}

func yXToField(y int8, x int8) game.Field {
	return y*BOARD_SIZE + x
}

func NewGame() *Game {
	turn := BLACK
	board := make([]game.Field, TOTAL_SIZE)
	history := make([]game.Move, 0, TOTAL_SIZE)
	board[yXToField(3, 3)], board[yXToField(4, 4)] = WHITE, WHITE
	board[yXToField(3, 4)], board[yXToField(4, 3)] = BLACK, BLACK

	return &Game{turn, board, history}
}

func (g *Game) Copy() game.Game {
	turn := g.Turn
	board := make([]game.Field, TOTAL_SIZE)
	history := make([]game.Move, len(g.History), cap(g.History))
	copy(board, g.Board)
	copy(history, g.History)

	return &Game{turn, board, history}
}

func (g *Game) MakeMove(move game.Move) (bool, game.PlayerColor) {
	currentPlayer := g.Turn
	g.Turn *= -1
	g.History = append(g.History, move)

	if move != -1 {
		g.Board[move] = currentPlayer

		for _, field := range getKilledPawns(g.Board, move, currentPlayer) {
			g.Board[field] = currentPlayer
		}
	}

	return g.IsGameFinished()
}

func (g *Game) GetPossibleMoves() []game.Move {
	result := make([]game.Field, 0, 8)
	result = append(result, -1)

	for field := int8(0); field < TOTAL_SIZE; field++ {
		if g.Board[field] == EMPTY && len(getKilledPawns(g.Board, field, g.Turn)) > 0 {
			result = append(result, field)
		}
	}

	return result
}

func (g *Game) GetCurrentPlayerColor() game.PlayerColor {
	return g.Turn
}

func (g *Game) IsGameFinished() (bool, game.PlayerColor) {
	turns := len(g.History)
	if turns < 2 || g.History[turns-2] != -1 || g.History[turns-1] != -1 {
		currentPlayerMoves := g.GetPossibleMoves()

		if len(currentPlayerMoves) > 1 {
			return false, EMPTY
		}

		g.Turn *= -1
		nextPlayerMoves := g.GetPossibleMoves()
		g.Turn *= -1

		if len(nextPlayerMoves) > 1 {
			return false, EMPTY
		}
	}

	// evaluating winner

	blacks, whites := g.CountPawns()
	winner := EMPTY

	if blacks > whites {
		winner = BLACK
	} else if whites > blacks {
		winner = WHITE
	}

	return true, winner
}

func (game *Game) CountPawns() (int8, int8) {
	black, white := int8(0), int8(0)

	for field := 0; field < TOTAL_SIZE; field++ {
		pawn := game.Board[field]

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
	result := make([]game.Field, 0, 4)
	startY, startX := getYX(start)
	deltas := []int8{-1, 0, 1}

	for _, dy := range deltas {
		for _, dx := range deltas {
			if dx == 0 && dy == 0 {
				continue
			}

			candidates := make([]game.Field, 0, 4)

			for y, x := startY+dy, startX+dx; x >= 0 && x < BOARD_SIZE && y >= 0 && y < BOARD_SIZE; y, x = y+dy, x+dx {
				field := yXToField(y, x)
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

func (game *Game) DrawBoard() {
	output := ""

	for field := 0; field < TOTAL_SIZE; field++ {
		switch game.Board[field] {
		case EMPTY:
			output += "."
		case BLACK:
			output += "x"
		case WHITE:
			output += "o"
		}
	}

	for _, move := range game.GetPossibleMoves() {
		if move != -1 {
			output = output[:move] + "+" + output[move+1:]
		}
	}

	for line := 0; line < BOARD_SIZE; line++ {
		start := line * BOARD_SIZE

		fmt.Println(output[start : start+BOARD_SIZE])
	}
}

func (game *Game) SerializeBoard() string {
	stringsBoard := make([]string, len(game.Board))

	for i, field := range game.Board {
		stringsBoard[i] = strconv.Itoa(int(field))
	}

	return strings.Join(stringsBoard, " ")
}

func (game *Game) OneHotBoard() [][][]int8 {
	oneHotBoard := make([][][]int8, BOARD_SIZE)

	for row := int8(0); row < BOARD_SIZE; row++ {
		oneHotBoard[row] = make([][]int8, BOARD_SIZE)

		for col := int8(0); col < BOARD_SIZE; col++ {
			field := game.Board[yXToField(row, col)]

			oneHotField := []int8{1, 0, 0}
			if field == -1 {
				oneHotField = []int8{0, 1, 0}
			} else if field == 1 {
				oneHotField = []int8{0, 0, 1}
			}

			oneHotBoard[row][col] = oneHotField
		}
	}

	return oneHotBoard
}
