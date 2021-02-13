package connectFour

import (
	"testing"

	"github.com/jamOne-/kiwi-zero/utils"

	"github.com/stretchr/testify/assert"
)

const e, b, w = EMPTY, BLACK, WHITE

func SetFullBoard(game *ConnectFourGame) {
	for x := 0; x < WIDTH; x++ {
		for y := 0; y < HEIGHT; y++ {
			if (x+y/3)%2 == 0 {
				game.Board[y*WIDTH+x] = WHITE
			} else {
				game.Board[y*WIDTH+x] = BLACK
			}
		}
	}
}

func TestIsGameFinishedFullBoard(t *testing.T) {
	game := NewConnectFourGame()
	SetFullBoard(game)
	finished, winner := game.IsGameFinished()

	assert.Equal(t, true, finished)
	assert.Equal(t, EMPTY, winner)
}

func TestIsGameFinishedAlmostFullBoard(t *testing.T) {
	game := NewConnectFourGame()
	SetFullBoard(game)
	game.Board[5] = EMPTY

	finished, _ := game.IsGameFinished()

	assert.Equal(t, false, finished)
}

func TestIsGameFinishedRowWin(t *testing.T) {
	game := NewConnectFourGame()

	y := HEIGHT - 2
	row := []int8{WHITE, WHITE, BLACK, BLACK, BLACK, BLACK, WHITE}

	for x := 0; x < WIDTH; x++ {
		game.Board[y*WIDTH+x] = row[x]
	}

	finished, winner := game.IsGameFinished()

	assert.Equal(t, true, finished)
	assert.Equal(t, BLACK, winner)
}

func TestIsGameFinishedColumnWin(t *testing.T) {
	game := NewConnectFourGame()
	game.Board = []int8{
		e, e, e, e, e, e, e,
		e, e, e, e, b, e, e,
		e, e, e, e, b, e, e,
		e, e, e, b, b, e, e,
		e, w, e, w, b, e, e,
		e, b, e, w, w, e, e,
	}

	finished, winner := game.IsGameFinished()

	assert.Equal(t, true, finished)
	assert.Equal(t, BLACK, winner)
}

func TestIsGameFinishedNWWin(t *testing.T) {
	game := NewConnectFourGame()
	game.Board = []int8{
		e, e, e, e, e, e, e,
		e, w, e, e, e, e, e,
		e, b, w, e, e, e, e,
		e, b, w, w, e, e, e,
		e, w, w, b, w, e, e,
		e, b, b, w, w, b, e,
	}

	finished, winner := game.IsGameFinished()

	assert.Equal(t, true, finished)
	assert.Equal(t, WHITE, winner)
}

func TestIsGameFinishedNEWin(t *testing.T) {
	game := NewConnectFourGame()
	game.Board = []int8{
		e, e, e, e, e, e, e,
		e, w, e, e, e, e, e,
		e, b, w, e, e, e, e,
		e, b, w, w, e, e, e,
		e, w, w, b, w, e, e,
		e, b, b, w, w, b, e,
	}
	utils.PerformSymmetryVector2(WIDTH, HEIGHT, game.Board)

	finished, winner := game.IsGameFinished()

	assert.Equal(t, true, finished)
	assert.Equal(t, WHITE, winner)
}
