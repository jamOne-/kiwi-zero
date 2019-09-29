package main

import (
	"fmt"

	"github.com/jamOne-/kiwi-zero/game"
	"github.com/jamOne-/kiwi-zero/humanPlayer"
	"github.com/jamOne-/kiwi-zero/minMaxPlayer"
	"github.com/jamOne-/kiwi-zero/reversi"
)

func main() {
	reversiGame := reversi.NewReversiGame()
	mmPlayer := minMaxPlayer.NewMinMaxPlayer(3, simpleReversiValueFn)
	humanPlayer := humanPlayer.NewHumanPlayer()

	finished, winner := false, game.PlayerColor(0)

	for !finished {
		currentPlayer := reversiGame.GetCurrentPlayerColor()
		var move game.Move

		if currentPlayer == game.BLACK {
			move = mmPlayer.SelectMove(reversiGame)
		} else {
			move = humanPlayer.SelectMove(reversiGame)
		}

		finished, winner = reversiGame.MakeMove(move)

		reversiGame.DrawBoard()
		fmt.Println("")
	}

	fmt.Println(winner)
}

func simpleReversiValueFn(game game.Game) float64 {
	reversiGame := game.(*reversi.ReversiGame) // nieładnie, ale brak generyków to jest jakiś dramat

	blacks, whites := 0.0, 0.0
	blackScore, whiteScore := 0.0, 0.0

	for i, pawn := range reversiGame.Board {
		if pawn == reversi.BLACK {
			blacks += 1
			blackScore += SCORING[i]
		} else if pawn == reversi.WHITE {
			whites += 1
			whiteScore += SCORING[i]
		}
	}

	p := 0.0
	if blacks > whites {
		p = 100.0 * blacks / (blacks + whites)
	} else if blacks < whites {
		p = -100.0 * whites / (blacks + whites)
	}

	return float64(game.GetCurrentPlayerColor()) * (p + blacks - whites)
}

var SCORING = []float64{
	20, -3, 11, 8, 8, 11, -3, 20,
	-3, -7, -4, 1, 1, -4, -7, -3,
	11, -4, 2, 2, 2, 2, -4, 11,
	8, 1, 2, -3, -3, 2, 1, 8,
	8, 1, 2, -3, -3, 2, 1, 8,
	11, -4, 2, 2, 2, 2, -4, 11,
	-3, -7, -4, 1, 1, -4, -7, -3,
	20, -3, 11, 8, 8, 11, -3, 20}
