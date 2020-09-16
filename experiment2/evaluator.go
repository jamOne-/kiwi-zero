package main

import (
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/spf13/viper"

	"github.com/jamOne-/kiwi-zero/game"
	"github.com/jamOne-/kiwi-zero/minMaxPlayer"
	"github.com/jamOne-/kiwi-zero/player"
	"github.com/jamOne-/kiwi-zero/reversiValueFns"
	"github.com/jamOne-/kiwi-zero/runner"
	"gonum.org/v1/gonum/mat"
)

type PlayerToCompare struct {
	name   string
	player player.Player
}

func Evaluator(
	newWeightsChan chan *mat.VecDense,
	bestWeightsChan chan *mat.VecDense,
	gameFactory runner.NewGameFactory,
	initialWeights *mat.VecDense,
	gameToFeaturesFn game.GameToFeaturesFn,
	playersToCompareWith []*PlayerToCompare,
	resultsDirPath string,
) {

	CHECKPOINT_EVERY := viper.GetInt("CHECKPOINT_EVERY")
	EVALUATOR_GAMES := viper.GetInt("EVALUATOR_GAMES")
	EVALUATOR_GAMES_AT_ONCE := viper.GetInt("EVALUATOR_GAMES_AT_ONCE")
	MAX_BEST_PLAYERS_POOL_LENGTH := viper.GetInt("MAX_BEST_PLAYERS_POOL_LENGTH")
	MINMAX_DEPTH := viper.GetInt("MINMAX_DEPTH")

	bestPlayer := createPlayer(gameToFeaturesFn, initialWeights, MINMAX_DEPTH)
	bestPlayersPool := []player.Player{bestPlayer}

	evaluator_i := 1

	for newWeights := range newWeightsChan {
		if newWeights == initialWeights {
			continue
		}

		newPlayer := createPlayer(gameToFeaturesFn, newWeights, MINMAX_DEPTH)
		// newPlayerWins := runner.ComparePlayersAsync(gameFactory, newPlayer, bestPlayer, EVALUATOR_GAMES, EVALUATOR_GAMES_AT_ONCE)
		newPlayerWins := runner.ComparePlayerWithOthersAsync(gameFactory, newPlayer, bestPlayersPool, EVALUATOR_GAMES, EVALUATOR_GAMES_AT_ONCE)

		fmt.Printf("Evaluator (%d): New candidate won %d/%d games\n", evaluator_i, newPlayerWins, EVALUATOR_GAMES)

		if float64(newPlayerWins)/float64(EVALUATOR_GAMES) >= 0.55 {
			fmt.Printf("Evaluator (%d): 🎉 New candidate is the new best player 🎉\n", evaluator_i)

			bestPlayer = newPlayer
			bestPlayersPool = append(bestPlayersPool, bestPlayer)

			if len(bestPlayersPool) > MAX_BEST_PLAYERS_POOL_LENGTH {
				bestPlayersPool = bestPlayersPool[len(bestPlayersPool)-MAX_BEST_PLAYERS_POOL_LENGTH:]
			}

			bestWeightsChan <- newWeights
		}

		if CHECKPOINT_EVERY > 0 && evaluator_i%CHECKPOINT_EVERY == 0 {
			go evaluatorCheckpoint(
				evaluator_i,
				resultsDirPath,
				newWeights,
				playersToCompareWith,
				gameFactory,
				bestPlayer,
			)
		}

		evaluator_i += 1
	}
}

func createPlayer(
	gameToFeaturesFn game.GameToFeaturesFn,
	weights *mat.VecDense,
	depth int,
) player.Player {
	valueFn := reversiValueFns.CreateWeightedReversiFn(gameToFeaturesFn, weights)
	player := minMaxPlayer.NewMinMaxPlayer(depth, valueFn)

	return player
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

	player1Wins := runner.ComparePlayersAsync(gameFactory, player1, player2, numberOfGames, maxGamesAtOnce)
	resultsInfo := fmt.Sprintf("%s won %d/%d games versus %s\n", player1Name, player1Wins, numberOfGames, player2Name)
	fmt.Print(resultsInfo)
	fmt.Fprint(resultsFile, resultsInfo)
}

func evaluatorCheckpoint(
	bestPlayer_i int,
	resultsDirPath string,
	newWeights *mat.VecDense,
	playersToCompareWith []*PlayerToCompare,
	gameFactory runner.NewGameFactory,
	bestPlayer player.Player) {

	COMPARE_AT_CHECKPOINTS := viper.GetBool("COMPARE_AT_CHECKPOINTS")
	COMPARE_AT_CHECKPOINTS_GAMES := viper.GetInt("COMPARE_AT_CHECKPOINTS_GAMES")
	EVALUATOR_GAMES_AT_ONCE := viper.GetInt("EVALUATOR_GAMES_AT_ONCE")

	bestPlayer_iString := strconv.Itoa(bestPlayer_i)
	checkpointWeightsPath := path.Join(resultsDirPath, bestPlayer_iString+"_weights.txt")

	reversiValueFns.SaveWeightsToFile(newWeights, checkpointWeightsPath)

	if COMPARE_AT_CHECKPOINTS {
		resultsPath := path.Join(resultsDirPath, bestPlayer_iString+"_results.txt")

		for _, playerToCompareWith := range playersToCompareWith {
			comparePlayersAndSaveResults(
				resultsPath,
				gameFactory,
				bestPlayer,
				fmt.Sprintf("MinMax (version=%d)", bestPlayer_i),
				playerToCompareWith.player,
				playerToCompareWith.name,
				COMPARE_AT_CHECKPOINTS_GAMES,
				EVALUATOR_GAMES_AT_ONCE,
			)
		}
	}
}
