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
	times1, times2 := make([]time.Duration, 32), make([]time.Duration, 32)
	NUMBER_OF_GAMES := 10

	for gameNumber := 0; gameNumber < NUMBER_OF_GAMES; gameNumber += 1 {
		g := reversi.NewGame()
		player1 := monteCarloTreeSearchPlayer.NewThreadedMonteCarloTreeSearchPlayer(2000, 4)
		player2 := minMaxPlayer.NewPredictorMinMaxPlayer(4)
		// player2 := minMaxPlayer.NewMinMaxPlayer(7)
		// player2 := randomPlayer.NewRandomPlayer()

		// game.DrawBoard()
		// fmt.Println("")

		finished, winner := false, int8(0)

		for !finished {
			var move game.Move
			start := time.Now()

			if g.Turn > 0 {
				move = player1.SelectMove(g)
				times1 = append(times1, time.Since(start))
			} else {
				move = player2.SelectMove(g)
				times2 = append(times2, time.Since(start))
			}

			// fmt.Printf("player %d was thinking for %s\n", (-g.Turn+1)/2+1, time.Since(start))
			// fmt.Println(move)
			finished, winner = g.MakeMove(move)
			// game.DrawBoard()
			// fmt.Println("")
		}

		fmt.Printf("Game %d/%d: %d wins\n", gameNumber+1, NUMBER_OF_GAMES, (-winner+1)/2+1)
		totalScore += int(winner)
	}

	player1Wins := (NUMBER_OF_GAMES + totalScore) / 2
	player2Wins := NUMBER_OF_GAMES - player1Wins

	fmt.Printf("Total score: %d - %d (%d)\n", player1Wins, player2Wins, totalScore)
	fmt.Printf("Player 1 was thinking for %s (average: %s, max: %s)\n", sumTimes(times1), time.Duration(float64(sumTimes(times1))/float64(len(times1))), maxTime(times1))
	fmt.Printf("Player 2 was thinking for %s (average: %s, max: %s)\n", sumTimes(times2), time.Duration(float64(sumTimes(times2))/float64(len(times2))), maxTime(times2))
}

func sumTimes(times []time.Duration) time.Duration {
	sum := time.Duration(0)

	for _, time := range times {
		sum += time
	}

	return sum
}

func maxTime(times []time.Duration) time.Duration {
	max := time.Duration(0)

	for _, time := range times {
		if time > max {
			max = time
		}
	}

	return max
}
