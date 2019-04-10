package main

import (
	"fmt"
	"time"

	"github.com/jamOne-/kiwi-zero/game"
	"github.com/jamOne-/kiwi-zero/minMaxPlayer"
	"github.com/jamOne-/kiwi-zero/monteCarloTreeSearchPlayer"
	"github.com/jamOne-/kiwi-zero/reversi"
)

func main() {
	totalScore := 0
	NUMBER_OF_GAMES := 10

	for gameNumber := 0; gameNumber < NUMBER_OF_GAMES; gameNumber += 1 {
		g := reversi.NewGame()
		player1 := monteCarloTreeSearchPlayer.NewThreadedMonteCarloTreeSearchPlayer(2000, 4)
		player2 := minMaxPlayer.NewMinMaxPlayer(7)
		// player2 := NewRandomPlayer()

		// game.DrawBoard()
		// fmt.Println("")

		finished, winner := false, int8(0)

		for !finished {
			var move game.Move
			start := time.Now()

			if g.Turn > 0 {
				move = player1.SelectMove(g)
			} else {
				move = player2.SelectMove(g)
			}

			fmt.Printf("player %d was thinking for %s\n", g.Turn, time.Since(start))
			// fmt.Println(move)
			finished, winner = g.MakeMove(move)
			// game.DrawBoard()
			// fmt.Println("")
		}

		fmt.Printf("Game %d/%d: %d wins\n", gameNumber+1, NUMBER_OF_GAMES, winner)
		totalScore += int(winner)
	}

	player1Wins := (NUMBER_OF_GAMES + totalScore) / 2
	player2Wins := NUMBER_OF_GAMES - player1Wins

	fmt.Printf("Total score: %d - %d (%d)\n", player1Wins, player2Wins, totalScore)
}
