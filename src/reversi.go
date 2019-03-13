package main

import "fmt"

type PlayerColor = int8
type Field = int8
type Move = Field
type Game struct {
	turn    PlayerColor
	board   []Field
	history []Move
}

const BOARD_SIZE = 8
const TOTAL_SIZE = BOARD_SIZE * BOARD_SIZE
const EMPTY = Field(0)
const WHITE = Field(-1)
const BLACK = Field(1)

func getYX(field Field) (int8, int8) {
	return field / BOARD_SIZE, field % BOARD_SIZE
}

func yXToField(y int8, x int8) Field {
	return y*BOARD_SIZE + x
}

func NewGame() *Game {
	turn := BLACK
	board := make([]Field, TOTAL_SIZE)
	history := make([]Move, 0, TOTAL_SIZE)
	board[yXToField(3, 3)], board[yXToField(4, 4)] = WHITE, WHITE
	board[yXToField(3, 4)], board[yXToField(4, 3)] = BLACK, BLACK

	return &Game{turn, board, history}
}

func (game *Game) Copy() *Game {
	turn := game.turn
	board := make([]Field, TOTAL_SIZE)
	history := make([]Move, len(game.history), TOTAL_SIZE)
	copy(board, game.board)
	copy(history, game.history)

	return &Game{turn, board, history}
}

func (game *Game) MakeMove(move Move) (bool, PlayerColor) {
	currentPlayer := game.turn
	game.turn *= -1
	game.history = append(game.history, move)

	if move != -1 {
		game.board[move] = currentPlayer

		for _, field := range getKilledPawns(game.board, move, currentPlayer) {
			game.board[field] = currentPlayer
		}
	}

	return game.IsGameFinished()
}

func (game *Game) GetPossibleMoves() []Move {
	result := make([]Field, 0, 8)
	result = append(result, -1)

	for field := int8(0); field < TOTAL_SIZE; field++ {
		if game.board[field] == EMPTY && len(getKilledPawns(game.board, field, game.turn)) > 0 {
			result = append(result, field)
		}
	}

	return result
}

func (game *Game) IsGameFinished() (bool, PlayerColor) {
	turns := len(game.history)
	if turns < 2 || game.history[turns-2] != -1 || game.history[turns-1] != -1 {
		currentPlayerMoves := game.GetPossibleMoves()

		if len(currentPlayerMoves) > 1 {
			return false, EMPTY
		}

		game.turn *= -1
		nextPlayerMoves := game.GetPossibleMoves()
		game.turn *= -1

		if len(nextPlayerMoves) > 1 {
			return false, EMPTY
		}
	}

	// evaluating winner

	blacks, whites := game.CountPawns()
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
		pawn := game.board[field]

		if pawn == BLACK {
			black++
		} else if pawn == WHITE {
			white++
		}
	}

	return black, white
}

func getKilledPawns(board []Field, start Field, player PlayerColor) []Field {
	opponent := player * -1
	result := make([]Field, 0, 4)
	startY, startX := getYX(start)
	deltas := []int8{-1, 0, 1}

	for _, dy := range deltas {
		for _, dx := range deltas {
			if dx == 0 && dy == 0 {
				continue
			}

			candidates := make([]Field, 0, 4)

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
		switch game.board[field] {
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
