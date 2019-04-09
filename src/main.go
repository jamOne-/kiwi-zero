package main

import (
	"fmt"
	"time"
)

func main() {
	totalScore := 0
	NUMBER_OF_GAMES := 10

	for gameNumber := 0; gameNumber < NUMBER_OF_GAMES; gameNumber += 1 {
		game := NewGame()
		player1 := NewMonteCarloTreeSearchPlayer()
		player2 := NewMinMaxPlayer(7)
		// player2 := NewRandomPlayer()

		// game.DrawBoard()
		// fmt.Println("")

		finished, winner := false, int8(0)

		for !finished {
			var move Move
			start := time.Now()

			if game.turn > 0 {
				move = player1.SelectMove(game)
			} else {
				move = player2.SelectMove(game)
			}

			fmt.Printf("player %d was thinking for %s\n", game.turn, time.Since(start))
			// fmt.Println(move)
			finished, winner = game.MakeMove(move)
			// game.DrawBoard()
			// fmt.Println("")
		}

		fmt.Printf("Game %d/%d: %d wins\n", gameNumber+1, NUMBER_OF_GAMES, winner)
		totalScore += int(winner)
	}

	fmt.Println("Total score:", totalScore)
}
