package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"

	"github.com/jamOne-/kiwi-zero/monteCarloTreeSearchPlayer"

	"github.com/jamOne-/kiwi-zero/randomPlayer"

	"github.com/jamOne-/kiwi-zero/edaxPlayer"

	"github.com/jamOne-/kiwi-zero/game"
	"github.com/jamOne-/kiwi-zero/reversi"
)

func NeutralValueFunction(game game.Game) float64 {
	return 0
}

func ReversiTempFn(turn int) float64 {
	return -0.05 + 1.3*math.Exp(-0.11*float64(turn))
}

func main() {
	rand.Seed(time.Now().UnixNano())

	totalScore := 0
	times1, times2 := make([]time.Duration, 0), make([]time.Duration, 0)
	NUMBER_OF_GAMES := 10

	// 1 layer:
	// tfpredictor := tfpredictor.NewTFPredictor("../experiment2/results/2020-11-04-211132/models/150")
	// 7 layers trained with smaller pool and capped positions
	// tfpredictor := tfpredictor.NewTFPredictor("../experiment2/results/2020-12-01-140644/models/100")
	// conv first try
	// tfpredictor := tfpredictor.NewTFPredictor("../experiment2/results/2020-11-26-202747/models/350")
	// conv softmax minmax 3
	// tfpredictor := tfpredictor.NewTFPredictor("../experiment2/results/2020-12-07-222844/models/250")
	// gameToFeaturesFn := reversiValueFns.ConvertReversiFnToGeneralFeatuersFn(reversiValueFns.ReversiToOneHotBoardMoves)
	//
	// tfpredictor := tfpredictor.NewTFPredictor("../experiment2/results/2021-01-23-205906/models/2525")
	// gameToFeaturesFn := reversiValueFns.ConvertReversiFnToGeneralFeatuersFn(reversiValueFns.ReversiToFeaturesExtended)

	// tfpredictor := tfpredictor.NewTFPredictor("../experiment2/results/2021-02-02-204309/models/150")
	// gameToFeaturesFn := reversiValueFns.ConvertReversiFnToGeneralFeatuersFn(reversiValueFns.ReversiToOneHotBoardMoves)

	// gameToDistributionFn := policyPlayer.GameToDistributionFnFromPredictor(gameToFeaturesFn, tfpredictor)

	// valueFn := reversiValueFns.CreateMinMaxValueFn(gameToFeaturesFn, tfpredictor)

	for gameNumber := 0; gameNumber < NUMBER_OF_GAMES; gameNumber += 1 {
		g := reversi.NewReversiGame()
		// player1 := humanPlayer.NewHumanPlayer()
		// player2 := monteCarloTreeSearchPlayer.NewThreadedMonteCarloTreeSearchPlayer(500, 1)
		// player2 := randomPlayer.NewRandomPlayer()
		// player1 := policyPlayer.NewPolicyPlayer(gameToDistributionFn)
		// player2 := minMaxPlayer.NewPredictorMinMaxPlayer(4)
		// player2 := minMaxPlayer.NewMinMaxPlayer(7)
		// player2 := randomPlayer.NewRandomPlayer()
		// player2 := monteCarloTreeSearchPlayer.NewGeneralMCTSPlayer(2000, 99, randomPlayer.NewRandomPlayer(), valueFn)

		// player1 := minMaxPlayer.NewMinMaxPlayer(4, valueFn)
		player1 := monteCarloTreeSearchPlayer.NewGeneralMCTSPlayer(
			100,
			2.0,
			99,
			randomPlayer.NewRandomPlayer(),
			NeutralValueFunction,
			nil,
			ReversiTempFn,
		)
		player2 := monteCarloTreeSearchPlayer.NewMonteCarloTreeSearchPlayer(100)
		// player1 := monteCarloTreeSearchPlayer.NewLazyExpandMonteCarloTreeSearchPlayer(1000)
		edax := edaxPlayer.NewEdaxPlayer(-64, 64, 1, 100)
		random := randomPlayer.NewRandomPlayer()
		// player2 := edax

		// g.DrawBoard()
		// fmt.Println("")

		finished, winner := false, int8(0)

		START_RANDOM_MOVES := 0
		for i := 0; i < START_RANDOM_MOVES; i++ {
			g.MakeMove(random.SelectMoveDifferentThan(g, reversi.PASS_MOVE))
			g.MakeMove(random.SelectMoveDifferentThan(g, reversi.PASS_MOVE))
		}

		for !finished {
			var move game.Move
			start := time.Now()

			if g.Turn > 0 {
				// var policy []float32
				// move, policy = player1.SelectMoveWithPolicy(g)
				move = player1.SelectMove(g)
				times1 = append(times1, time.Since(start))
				// fmt.Println(policy)
			} else {
				move = player2.SelectMove(g)
				times2 = append(times2, time.Since(start))
			}

			// g.DrawBoard()

			// fmt.Printf("player %d was thinking for %s\n", (-g.Turn+1)/2+1, time.Since(start))
			// fmt.Println(move)
			finished, winner = g.MakeMove(move)
			// g.DrawBoard()
			// value := tfpredictor.Predict(gameToFeaturesFn(g))
			// policy := tfpredictor.PredictPolicy(gameToFeaturesFn(g))
			// fmt.Println(len(policy))
			// fmt.Println("")
		}
		g.DrawBoard()

		edax.ClosePlayer()

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
