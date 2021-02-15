package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"sort"
	"time"

	tfpredictor "github.com/jamOne-/kiwi-zero/TFPredictor"

	"github.com/jamOne-/kiwi-zero/policyPlayer"
	"github.com/jamOne-/kiwi-zero/predictor"
	"github.com/jamOne-/kiwi-zero/reversiValueFns"

	"github.com/jamOne-/kiwi-zero/humanPlayer"
	"github.com/jamOne-/kiwi-zero/minMaxPlayer"

	"github.com/jamOne-/kiwi-zero/edaxPlayer"
	"github.com/jamOne-/kiwi-zero/games"
	"github.com/jamOne-/kiwi-zero/monteCarloTreeSearchPlayer"
	"github.com/jamOne-/kiwi-zero/player"
	"github.com/jamOne-/kiwi-zero/randomPlayer"

	"github.com/jamOne-/kiwi-zero/game"
)

func NeutralValueFunction(game game.Game) float64 {
	return 0
}

func ReversiTempFn(turn int) float64 {
	return -0.05 + 1.3*math.Exp(-0.11*float64(turn))
}

func instantiatePlayer(
	gameName string,
	playerType string,
	mmDepth int,
	mmEpsilon float64,
	pred predictor.Predictor,
	modelFn string,
	mctsSims int,
	mctsDepth int,
) player.Player {
	switch playerType {
	case "random":
		return randomPlayer.NewRandomPlayer()
	case "human":
		return humanPlayer.NewHumanPlayer()
	case "edax":
		return edaxPlayer.NewEdaxPlayer(-64, 64, mmDepth, 100)
	case "minmax":
		featuresFn := games.FEATURES_FNS_DICT[gameName][modelFn].Fn
		valueFn := reversiValueFns.CreateMinMaxValueFn(featuresFn, pred)
		return minMaxPlayer.NewMinMaxPlayer(mmDepth, valueFn)
	case "minmax-e":
		featuresFn := games.FEATURES_FNS_DICT[gameName][modelFn].Fn
		valueFn := reversiValueFns.CreateMinMaxValueFn(featuresFn, pred)
		return minMaxPlayer.NewEpsilonGreedyMinMaxPlayer(mmDepth, mmEpsilon, valueFn)
	case "mcts":
		var valueFn game.ValueFn = nil
		var valueAndPolicyFn monteCarloTreeSearchPlayer.ValueAndPolicyFn = nil

		if pred != nil {
			featuresFn := games.FEATURES_FNS_DICT[gameName][modelFn].Fn
			valueFn = reversiValueFns.CreateMinMaxValueFn(featuresFn, pred)
			valueAndPolicyFn = predictor.CreateValueAndPolicyFn(featuresFn, pred)
		}

		return monteCarloTreeSearchPlayer.NewGeneralMCTSPlayer(
			mctsSims,
			2.0,
			mctsDepth,
			randomPlayer.NewRandomPlayer(),
			valueFn,
			valueAndPolicyFn,
			nil,
		)
	case "policy":
		featuresFn := games.FEATURES_FNS_DICT[gameName][modelFn].Fn
		gameToDistributionFn := policyPlayer.GameToDistributionFnFromPredictor(featuresFn, pred)
		return policyPlayer.NewPolicyPlayer(gameToDistributionFn)
	}

	log.Fatal("Wrong player type!")
	return nil
}

func main() {
	rand.Seed(time.Now().UnixNano())

	NUMBER_OF_GAMES := flag.Int("games", 10, "Number of games to play")
	GAME_NAME := flag.String("game", "reversirandom", "reversi|reversirandom|connect4|connect4random|gomoku|gomokurandom")
	DRAW_BOARD := flag.Bool("draw-board", false, "If set, draws board after every move")
	P1_TYPE := flag.String("p1", "", "minmax|minmax-e|mcts|policy|random|human|edax")
	P2_TYPE := flag.String("p2", "", "minmax|minmax-e|mcts|policy|random|human|edax")
	P1_MINMAX_DEPTH := flag.Int("p1-minmax-depth", 3, "Depth of Minmax search for Player 1")
	P2_MINMAX_DEPTH := flag.Int("p2-minmax-depth", 3, "Depth of Minmax search for Player 2")
	P1_MINMAX_EPSILON := flag.Float64("p1-minmax-e", 3, "Epsilon value of Minmax-e for Player 1")
	P2_MINMAX_EPSILON := flag.Float64("p2-minmax-e", 3, "Epsilon value of Minmax-e for Player 2")
	P1_MODEL_PATH := flag.String("p1-model-path", "", "Path of serialized model for Player 1")
	P2_MODEL_PATH := flag.String("p2-model-path", "", "Path of serialized model for Player 2")
	P1_MODEL_FN := flag.String("p1-model-fn", "board3", "board1|board1features|b1turn|board3|b3turn|bmt|boardmoves|paddedmoves")
	P2_MODEL_FN := flag.String("p2-model-fn", "board3", "board1|board1features|b1turn|board3|b3turn|bmt|boardmoves|paddedmoves")
	P1_MCTS_SIMS := flag.Int("p1-mcts-sims", 1000, "Number of simulations for MCTS for Player 1")
	P2_MCTS_SIMS := flag.Int("p2-mcts-sims", 1000, "Number of simulations for MCTS for Player 2")
	P1_MCTS_DEPTH := flag.Int("p1-mcts-depth", 99, "Rollout depth of MCTS simulations for Player 1")
	P2_MCTS_DEPTH := flag.Int("p2-mcts-depth", 99, "Rollout depth of MCTS simulations for Player 2")
	flag.Parse()

	var pred1 predictor.Predictor = nil
	if *P1_MODEL_PATH != "" {
		pred1 = tfpredictor.NewTFPredictor(*P1_MODEL_PATH)
	}
	var pred2 predictor.Predictor = nil
	if *P2_MODEL_PATH != "" {
		pred2 = tfpredictor.NewTFPredictor(*P2_MODEL_PATH)
	}

	gameFactory := games.GAME_FACTORY_DICT[*GAME_NAME]
	p1Score, p2Score := 0, 0
	times1, times2 := &[]time.Duration{}, &[]time.Duration{}

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

	gamesToPlay := make([]game.Game, *NUMBER_OF_GAMES)
	for i := 0; i < *NUMBER_OF_GAMES/2; i++ {
		g := gameFactory()
		gamesToPlay[2*i] = g
		gamesToPlay[2*i+1] = g.Copy()
	}
	if *NUMBER_OF_GAMES%2 == 1 {
		gamesToPlay[*NUMBER_OF_GAMES-1] = gameFactory()
	}

	for gameNumber := 0; gameNumber < *NUMBER_OF_GAMES; gameNumber += 1 {
		g := gamesToPlay[gameNumber]
		finished, winner := false, int8(0)
		blackPlayer := instantiatePlayer(
			*GAME_NAME,
			*P1_TYPE,
			*P1_MINMAX_DEPTH,
			*P1_MINMAX_EPSILON,
			pred1,
			*P1_MODEL_FN,
			*P1_MCTS_SIMS,
			*P1_MCTS_DEPTH,
		)
		whitePlayer := instantiatePlayer(
			*GAME_NAME,
			*P2_TYPE,
			*P2_MINMAX_DEPTH,
			*P2_MINMAX_EPSILON,
			pred2,
			*P2_MODEL_FN,
			*P2_MCTS_SIMS,
			*P2_MCTS_DEPTH,
		)
		blackTimes, whiteTimes := times1, times2
		blackScore, whiteScore := &p1Score, &p2Score

		if gameNumber%2 == 1 {
			blackPlayer, whitePlayer = whitePlayer, blackPlayer
			blackTimes, whiteTimes = whiteTimes, blackTimes
			blackScore, whiteScore = whiteScore, blackScore
		}

		for !finished {
			var move game.Move
			start := time.Now()

			if g.GetCurrentPlayerColor() == game.BLACK {
				move = blackPlayer.SelectMove(g)
				*blackTimes = append(*blackTimes, time.Since(start))
			} else {
				move = whitePlayer.SelectMove(g)
				*whiteTimes = append(*whiteTimes, time.Since(start))
			}
			finished, winner = g.MakeMove(move)

			if *DRAW_BOARD {
				g.DrawBoard()
				fmt.Printf("\n")
			}
		}

		if winner == 0 {
			fmt.Printf("Game %d/%d: ended with a draw\n", gameNumber+1, *NUMBER_OF_GAMES)
		} else {
			if winner == game.BLACK {
				*blackScore += 1
			} else if winner == game.WHITE {
				*whiteScore += 1
			}

			if gameNumber%2 == 1 {
				winner *= -1
			}

			fmt.Printf("Game %d/%d: %d wins\n", gameNumber+1, *NUMBER_OF_GAMES, (-winner+1)/2+1)
		}
	}

	fmt.Printf("Total score: %d - %d\n", p1Score, p2Score)
	fmt.Printf("Player 1 was thinking for %s (average: %s, max: %s, median: %s)\n", sumTimes(*times1), time.Duration(float64(sumTimes(*times1))/float64(len(*times1))), maxTime(*times1), median(*times1))
	fmt.Printf("Player 2 was thinking for %s (average: %s, max: %s, median: %s)\n", sumTimes(*times2), time.Duration(float64(sumTimes(*times2))/float64(len(*times2))), maxTime(*times2), median(*times2))
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
