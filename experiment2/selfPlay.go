package main

import (
	"fmt"
	"math"

	"github.com/jamOne-/kiwi-zero/player"
	"github.com/jamOne-/kiwi-zero/predictor"

	"github.com/jamOne-/kiwi-zero/runner"
	"github.com/spf13/viper"
)

func SelfPlayLoop(
	bestPredictors chan predictor.Predictor,
	gameResults chan *runner.GameResultsBatch,
	gameFactory runner.NewGameFactory,
	initialPredictor predictor.Predictor,
	selfPlayFactory player.PlayerFactory,
) {
	GAMES_PER_ITERATION := viper.GetInt("GAMES_PER_ITERATION")
	SELFPLAY_GAMES_AT_ONCE := viper.GetInt("SELFPLAY_GAMES_AT_ONCE")

	if SELFPLAY_GAMES_AT_ONCE == 0 {
		SELFPLAY_GAMES_AT_ONCE = GAMES_PER_ITERATION
	}

	selfPlay_i := 1
	selfPlayPlayer := selfPlayFactory(initialPredictor)

	for {
		select {
		case predictor := <-bestPredictors:
			if predictor != nil {
				selfPlayPlayer = selfPlayFactory(predictor)
			}

			// default:
			gamesCount := GAMES_PER_ITERATION
			if selfPlay_i == 1 {
				gamesCount = GAMES_PER_ITERATION * int(2+math.Ceil(viper.GetFloat64("OPTIMIZER_TRAINING_SIZE")/float64(GAMES_PER_ITERATION)))
			}

			results, totalPositions := runner.PlayNGamesAsync(
				gameFactory,
				/* saveHistory */ true,
				selfPlayPlayer,
				selfPlayPlayer,
				gamesCount,
				SELFPLAY_GAMES_AT_ONCE,
			)

			fmt.Printf("Selfplay (%d): finished %d games\n", selfPlay_i, gamesCount)

			selfPlay_i += 1
			resultsBatch := &runner.GameResultsBatch{Results: results, TotalPositions: totalPositions}

			gameResults <- resultsBatch
			// select {
			// case gameResults <- resultsBatch:
			// 	// try to send
			// default:
			// 	// else skip
			// }
		}
	}
}
