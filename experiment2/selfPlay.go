package main

import (
	"github.com/jamOne-/kiwi-zero/minMaxPlayer"
	"github.com/jamOne-/kiwi-zero/player"
	"github.com/jamOne-/kiwi-zero/reversiValueFns"
	"github.com/jamOne-/kiwi-zero/runner"
	"github.com/spf13/viper"
	"gonum.org/v1/gonum/mat"
)

func SelfPlayLoop(bestWeights chan *mat.VecDense, gameResults chan *runner.GameResultsBatch, initialWeights *mat.VecDense, gameFactory runner.NewGameFactory) {
	EPSILON := viper.GetFloat64("EPSILON")
	GAMES_PER_ITERATION := viper.GetInt("GAMES_PER_ITERATION")
	MINMAX_DEPTH := viper.GetInt("MINMAX_DEPTH")
	SELFPLAY_GAMES_AT_ONCE := viper.GetInt("SELFPLAY_GAMES_AT_ONCE")

	if SELFPLAY_GAMES_AT_ONCE == 0 {
		SELFPLAY_GAMES_AT_ONCE = GAMES_PER_ITERATION
	}

	selfPlayPlayer := createSelfPlayPlayer(initialWeights, MINMAX_DEPTH, EPSILON)

	for {
		select {
		case newWeights := <-bestWeights:
			selfPlayPlayer = createSelfPlayPlayer(newWeights, MINMAX_DEPTH, EPSILON)

		default:
			results, totalPositions := runner.PlayNGamesAsync(gameFactory, selfPlayPlayer, selfPlayPlayer, GAMES_PER_ITERATION, SELFPLAY_GAMES_AT_ONCE)
			resultsBatch := &runner.GameResultsBatch{Results: results, TotalPositions: totalPositions}

			select {
			case gameResults <- resultsBatch:
				// try to send
			default:
				// else skip
			}
		}
	}
}

func createSelfPlayPlayer(weights *mat.VecDense, depth int, epsilon float64) player.Player {
	reversiToFeatures := reversiValueFns.ReversiToFeatures
	valueFn := reversiValueFns.CreateWeightedReversiFn(reversiToFeatures, weights)
	selfPlayPlayer := minMaxPlayer.NewEpsilonGreedyMinMaxPlayer(depth, epsilon, valueFn)

	return selfPlayPlayer
}
