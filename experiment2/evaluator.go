package main

import (
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/spf13/viper"

	"github.com/jamOne-/kiwi-zero/player"
	"github.com/jamOne-/kiwi-zero/predictor"
	"github.com/jamOne-/kiwi-zero/runner"
)

type PlayerToCompare struct {
	name   string
	player player.Player
}

func Evaluator(
	newPredictors chan predictor.Predictor,
	bestPredictors chan predictor.Predictor,
	gameFactory runner.NewGameFactory,
	initialPredictor predictor.Predictor,
	evaluatorPlayerFactory player.PlayerFactory,
	playersToCompareWith []*PlayerToCompare,
	resultsDirPath string,
) {
	CHECKPOINT_EVERY := viper.GetInt("CHECKPOINT_EVERY")
	EVALUATOR_GAMES := viper.GetInt("EVALUATOR_GAMES")
	EVALUATOR_GAMES_AT_ONCE := viper.GetInt("EVALUATOR_GAMES_AT_ONCE")
	MAX_BEST_PLAYERS_POOL_LENGTH := viper.GetInt("MAX_BEST_PLAYERS_POOL_LENGTH")
	// MINMAX_DEPTH := viper.GetInt("MINMAX_DEPTH")

	// bestPlayer := minMaxPlayer.NewMinMaxPlayer(MINMAX_DEPTH, initialValueFn)
	bestPlayer := evaluatorPlayerFactory(initialPredictor)
	bestPlayersPool := []player.Player{bestPlayer}
	bestPredictor := initialPredictor

	evaluator_i := 1

	for newPredictor := range newPredictors {
		// if newModelPath == initialWeights {
		// 	continue
		// }

		// newPlayer := minMaxPlayer.NewMinMaxPlayer(MINMAX_DEPTH, newPredictor)
		newPlayer := evaluatorPlayerFactory(newPredictor)
		// newPlayerWins := runner.ComparePlayersAsync(gameFactory, newPlayer, bestPlayer, EVALUATOR_GAMES, EVALUATOR_GAMES_AT_ONCE)
		newPlayerWins, gamesPlayed := runner.ComparePlayerWithOthersAsync(gameFactory, newPlayer, bestPlayersPool, EVALUATOR_GAMES, EVALUATOR_GAMES_AT_ONCE)

		fmt.Printf("Evaluator (%d): New candidate won %d/%d games\n", evaluator_i, newPlayerWins, gamesPlayed)

		if float64(newPlayerWins)/float64(gamesPlayed) >= 0.6 {
			fmt.Printf("Evaluator (%d): 🎉 New candidate is the new best player 🎉\n", evaluator_i)

			bestPlayer = newPlayer
			bestPlayersPool = append(bestPlayersPool, bestPlayer)
			bestPredictor = newPredictor

			if len(bestPlayersPool) > MAX_BEST_PLAYERS_POOL_LENGTH {
				bestPlayersPool = bestPlayersPool[len(bestPlayersPool)-MAX_BEST_PLAYERS_POOL_LENGTH:]
			}

			bestPredictors <- newPredictor
		} else {
			bestPredictors <- nil
		}

		if CHECKPOINT_EVERY > 0 && evaluator_i%CHECKPOINT_EVERY == 0 {
			go evaluatorCheckpoint(
				evaluator_i,
				resultsDirPath,
				playersToCompareWith,
				gameFactory,
				bestPlayer,
				bestPredictor,
			)
		}

		evaluator_i += 1
	}
}

func comparePlayersAndSaveResults(
	filePath string,
	gameFactory runner.NewGameFactory,
	player1 player.Player,
	player1Name string,
	player2 player.Player,
	player2Name string,
	numberOfGames int,
	maxGamesAtOnce int,
) {
	resultsFile, _ := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer resultsFile.Close()

	player1Wins, gamesPlayed := runner.ComparePlayerWithOthersAsync(gameFactory, player1, []player.Player{player2}, numberOfGames, maxGamesAtOnce)
	resultsInfo := fmt.Sprintf("%s won %d/%d games versus %s\n", player1Name, player1Wins, gamesPlayed, player2Name)
	fmt.Print(resultsInfo)
	fmt.Fprint(resultsFile, resultsInfo)
}

func evaluatorCheckpoint(
	evaluator_i int,
	resultsDirPath string,
	playersToCompareWith []*PlayerToCompare,
	gameFactory runner.NewGameFactory,
	bestPlayer player.Player,
	bestPredictor predictor.Predictor,
) {
	COMPARE_AT_CHECKPOINTS := viper.GetBool("COMPARE_AT_CHECKPOINTS")
	COMPARE_AT_CHECKPOINTS_GAMES := viper.GetInt("COMPARE_AT_CHECKPOINTS_GAMES")
	EVALUATOR_GAMES_AT_ONCE := viper.GetInt("EVALUATOR_GAMES_AT_ONCE")

	evaluator_iString := strconv.Itoa(evaluator_i)
	// checkpointWeightsPath := path.Join(resultsDirPath, evaluator_iString+"_weights.txt")

	// reversiValueFns.SaveWeightsToFile(newWeights, checkpointWeightsPath)

	if COMPARE_AT_CHECKPOINTS {
		resultsPath := path.Join(resultsDirPath, evaluator_iString+"_results.txt")
		resultsFile, _ := os.OpenFile(resultsPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		fmt.Fprintf(resultsFile, "Predictor id (version=%d): %s\n", evaluator_i, bestPredictor.GetId())
		resultsFile.Close()

		for _, playerToCompareWith := range playersToCompareWith {
			comparePlayersAndSaveResults(
				resultsPath,
				gameFactory,
				bestPlayer,
				fmt.Sprintf("Reinforced (version=%d)", evaluator_i),
				playerToCompareWith.player,
				playerToCompareWith.name,
				COMPARE_AT_CHECKPOINTS_GAMES,
				EVALUATOR_GAMES_AT_ONCE,
			)
		}
	}
}
