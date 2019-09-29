package humanPlayer

import (
	"fmt"

	"github.com/jamOne-/kiwi-zero/game"
)

type HumanPlayer struct{}

func NewHumanPlayer() *HumanPlayer {
	return &HumanPlayer{}
}

func (player *HumanPlayer) SelectMove(g game.Game) game.Move {
	var move game.Move

	fmt.Println(g.GetPossibleMoves())
	fmt.Scan(&move)

	return move
}
