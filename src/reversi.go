package main

import "fmt"

type Player = int8
type Field = int8
type Move = Field
type Game struct {
	turn  Player
	board []Field
}

const BOARD_SIZE = 8
const TOTAL_SIZE = BOARD_SIZE * BOARD_SIZE
const EMPTY = Field(0)
const WHITE = Field(-1)
const BLACK = Field(1)

func GetYX(field Field) (int8, int8) {
	return field / BOARD_SIZE, field % BOARD_SIZE
}

func YXToField(y int8, x int8) Field {
	return y*BOARD_SIZE + x
}

func NewGame() *Game {
	turn := BLACK
	board := make([]Field, TOTAL_SIZE)
	board[YXToField(3, 3)], board[YXToField(4, 4)] = WHITE, WHITE
	board[YXToField(3, 4)], board[YXToField(4, 3)] = BLACK, BLACK

	return &Game{turn, board}
}

func (game *Game) MakeMove(move Move) (bool, Player) {
	currentPlayer := game.turn
	game.turn *= -1

	if move == -1 {
		return false, EMPTY
	}

	game.board[move] = currentPlayer

	for _, field := range GetKilledPawns(game.board, move, currentPlayer) {
		game.board[field] = currentPlayer
	}

	return game.IsGameFinished()
}

func (game *Game) GetPossibleMoves() []Move {
	result := make([]Field, 0, 8)
	result = append(result, -1)

	for field := int8(0); field < TOTAL_SIZE; field++ {
		if len(GetKilledPawns(game.board, field, game.turn)) > 0 {
			result = append(result, field)
		}
	}

	return result
}

func (game *Game) IsGameFinished() (bool, Player) {
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

func GetKilledPawns(board []Field, start Field, player Player) []Field {
	opponent := player * -1
	result := make([]Field, 0, 4)
	startY, startX := GetYX(start)
	deltas := []int8{-1, 0, 1}

	for _, dy := range deltas {
		for _, dx := range deltas {
			if dx == 0 && dy == 0 {
				continue
			}

			candidates := make([]Field, 0, 4)

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
