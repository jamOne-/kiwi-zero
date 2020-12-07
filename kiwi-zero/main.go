package main

import (
	"fmt"
	"math"
	"sort"
	"time"

	tfpredictor "github.com/jamOne-/kiwi-zero/TFPredictor"
	"github.com/jamOne-/kiwi-zero/minMaxPlayer"
	"github.com/jamOne-/kiwi-zero/policyPlayer"
	"github.com/jamOne-/kiwi-zero/reversiValueFns"

	"github.com/jamOne-/kiwi-zero/game"
	"github.com/jamOne-/kiwi-zero/reversi"
)

func main() {
	totalScore := 0
	times1, times2 := make([]time.Duration, 0), make([]time.Duration, 0)
	NUMBER_OF_GAMES := 100

	// 1 layer:
	// tfpredictor := tfpredictor.NewTFPredictor("../experiment2/results/2020-11-04-211132/models/150")
	// 7 layers trained with smaller pool and capped positions
	// tfpredictor := tfpredictor.NewTFPredictor("../experiment2/results/2020-12-01-140644/models/123")
	// conv first try
	tfpredictor := tfpredictor.NewTFPredictor("../experiment2/results/2020-11-26-202747/models/350")
	gameToFeaturesFn := reversiValueFns.ConvertReversiFnToGeneralFeatuersFn(reversiValueFns.ReversiToOneHotBoardMoves)
	gameToDistributionFn := policyPlayer.GameToDistributionFnFromTfPredictor(gameToFeaturesFn, tfpredictor)

	valueFn := reversiValueFns.CreateMinMaxValueFn(gameToFeaturesFn, tfpredictor)

	for gameNumber := 0; gameNumber < NUMBER_OF_GAMES; gameNumber += 1 {
		g := reversi.NewReversiGame()
		// player1 := humanPlayer.NewHumanPlayer()
		// player1 := monteCarloTreeSearchPlayer.NewThreadedMonteCarloTreeSearchPlayer(2000, 1)
		player1 := policyPlayer.NewPolicyPlayer(gameToDistributionFn)
		// player2 := minMaxPlayer.NewPredictorMinMaxPlayer(4)
		// player2 := minMaxPlayer.NewMinMaxPlayer(7)
		// player2 := randomPlayer.NewRandomPlayer()

		player2 := minMaxPlayer.NewMinMaxPlayer(3, valueFn)

		// g.DrawBoard()
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
			// g.DrawBoard()
			// value := tfpredictor.Predict(gameToFeaturesFn(g))
			// policy := tfpredictor.PredictPolicy(gameToFeaturesFn(g))
			// fmt.Println(len(policy))
			// fmt.Println("")
		}

		fmt.Printf("Game %d/%d: %d wins\n", gameNumber+1, NUMBER_OF_GAMES, (-winner+1)/2+1)
		totalScore += int(winner)
	}

	player1Wins := (NUMBER_OF_GAMES + totalScore) / 2
	player2Wins := NUMBER_OF_GAMES - player1Wins

	fmt.Printf("Total score: %d - %d (%d)\n", player1Wins, player2Wins, totalScore)
	fmt.Printf("Player 1 was thinking for %s (average: %s, max: %s, median: %s)\n", sumTimes(times1), time.Duration(float64(sumTimes(times1))/float64(len(times1))), maxTime(times1), median(times1))
	fmt.Printf("Player 2 was thinking for %s (average: %s, max: %s, median: %s)\n", sumTimes(times2), time.Duration(float64(sumTimes(times2))/float64(len(times2))), maxTime(times2), median(times2))
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

type Durations []time.Duration

func (a Durations) Len() int           { return len(a) }
func (a Durations) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Durations) Less(i, j int) bool { return a[i] < a[j] }

func median(times Durations) time.Duration {
	sort.Sort(times)
	return times[int(math.Floor(float64(len(times))/2))]
}
