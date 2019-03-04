package main

import "fmt"

func main() {
	game := NewGame()
	player1 := NewHumanPlayer()
	player2 := NewRandomPlayer()

	game.DrawBoard()
	fmt.Println("")

	finished, winner := false, int8(0)

	for !finished {
		var move Move

		if game.turn > 0 {
			move = player1.SelectMove(game)
		} else {
			move = player2.SelectMove(game)
		}

		finished, winner = game.MakeMove(move)
		game.DrawBoard()
		fmt.Println("")
	}

	fmt.Println(winner)
}
