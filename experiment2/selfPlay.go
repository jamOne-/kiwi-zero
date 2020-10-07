package main

import (
	"fmt"

	"github.com/jamOne-/kiwi-zero/game"
	"github.com/jamOne-/kiwi-zero/minMaxPlayer"
	"github.com/jamOne-/kiwi-zero/runner"
	"github.com/spf13/viper"
)

func SelfPlayLoop(
	bestValueFns chan game.ValueFn,
	gameResults chan *runner.GameResultsBatch,
	gameFactory runner.NewGameFactory,
	initialValueFn game.ValueFn,
) {
	EPSILON := viper.GetFloat64("EPSILON")
	GAMES_PER_ITERATION := viper.GetInt("GAMES_PER_ITERATION")
	MINMAX_DEPTH := viper.GetInt("MINMAX_DEPTH")
	SELFPLAY_GAMES_AT_ONCE := viper.GetInt("SELFPLAY_GAMES_AT_ONCE")

	if SELFPLAY_GAMES_AT_ONCE == 0 {
		SELFPLAY_GAMES_AT_ONCE = GAMES_PER_ITERATION
	}

	selfPlay_i := 1
	selfPlayPlayer := minMaxPlayer.NewEpsilonGreedyMinMaxPlayer(MINMAX_DEPTH, EPSILON, initialValueFn)

	for {
		select {
		case valueFn := <-bestValueFns:
			selfPlayPlayer = minMaxPlayer.NewEpsilonGreedyMinMaxPlayer(MINMAX_DEPTH, EPSILON, valueFn)

		default:
			results, totalPositions := runner.PlayNGamesAsync(
				gameFactory,
				/* saveHistory */ true,
				selfPlayPlayer,
				selfPlayPlayer,
				GAMES_PER_ITERATION,
				SELFPLAY_GAMES_AT_ONCE,
			)

			fmt.Printf("Selfplay (%d): finished %d games\n", selfPlay_i, GAMES_PER_ITERATION)

			selfPlay_i += 1
			resultsBatch := &runner.GameResultsBatch{Results: results, TotalPositions: totalPositions}

			// gameResults <- resultsBatch
			select {
			case gameResults <- resultsBatch:
				// try to send
			default:
				// else skip
			}
		}
	}
}
