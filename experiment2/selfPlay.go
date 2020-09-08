package main

import (
	"fmt"

	"github.com/jamOne-/kiwi-zero/game"
	"github.com/jamOne-/kiwi-zero/minMaxPlayer"
	"github.com/jamOne-/kiwi-zero/player"
	"github.com/jamOne-/kiwi-zero/reversiValueFns"
	"github.com/jamOne-/kiwi-zero/runner"
	"github.com/spf13/viper"
	"gonum.org/v1/gonum/mat"
)

func SelfPlayLoop(
	bestWeights chan *mat.VecDense,
	gameResults chan *runner.GameResultsBatch,
	initialWeights *mat.VecDense,
	gameFactory runner.NewGameFactory,
	gameToFeaturesFn game.GameToFeaturesFn,
) {
	EPSILON := viper.GetFloat64("EPSILON")
	GAMES_PER_ITERATION := viper.GetInt("GAMES_PER_ITERATION")
	MINMAX_DEPTH := viper.GetInt("MINMAX_DEPTH")
	SELFPLAY_GAMES_AT_ONCE := viper.GetInt("SELFPLAY_GAMES_AT_ONCE")

	if SELFPLAY_GAMES_AT_ONCE == 0 {
		SELFPLAY_GAMES_AT_ONCE = GAMES_PER_ITERATION
	}

	selfPlay_i := 1
	selfPlayPlayer := createSelfPlayPlayer(gameToFeaturesFn, initialWeights, MINMAX_DEPTH, EPSILON)

	for {
		select {
		case newWeights := <-bestWeights:
			selfPlayPlayer = createSelfPlayPlayer(gameToFeaturesFn, newWeights, MINMAX_DEPTH, EPSILON)

		default:
			results, totalPositions := runner.PlayNGamesAsync(
				gameFactory,
				gameToFeaturesFn,
				selfPlayPlayer,
				selfPlayPlayer,
				GAMES_PER_ITERATION,
				SELFPLAY_GAMES_AT_ONCE,
			)

			fmt.Printf("Selfplay (%d): finished %d games", selfPlay_i, GAMES_PER_ITERATION)

			selfPlay_i += 1
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

func createSelfPlayPlayer(
	gameToFeaturesFn game.GameToFeaturesFn,
	weights *mat.VecDense,
	depth int,
	epsilon float64,
) player.Player {
	valueFn := reversiValueFns.CreateWeightedReversiFn(gameToFeaturesFn, weights)
	selfPlayPlayer := minMaxPlayer.NewEpsilonGreedyMinMaxPlayer(depth, epsilon, valueFn)

	return selfPlayPlayer
}
