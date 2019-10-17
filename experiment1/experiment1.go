package main

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/jamOne-/kiwi-zero/game"
	"github.com/jamOne-/kiwi-zero/minMaxPlayer"
	"github.com/jamOne-/kiwi-zero/monteCarloTreeSearchPlayer"
	"github.com/jamOne-/kiwi-zero/player"
	"github.com/jamOne-/kiwi-zero/reversi"
	"github.com/jamOne-/kiwi-zero/reversiValueFns"
	"github.com/jamOne-/kiwi-zero/runner"
	"github.com/jamOne-/kiwi-zero/sgd"
	"github.com/jamOne-/kiwi-zero/utils"
	"github.com/spf13/viper"

	"gonum.org/v1/gonum/mat"
)

var INITIAL_WEIGHTS_BY_MODE = map[string](func() *mat.VecDense){
	"normal":   reversiValueFns.GetInitialWeights,
	"triangle": reversiValueFns.GetTriangleInitialWeights}

var REVERSI_TO_FEATURES_BY_MODE = map[string]reversiValueFns.ReversiToFeaturesFn{
	"normal":   reversiValueFns.ReversiToFeatures,
	"triangle": reversiValueFns.ReversiToFeaturesTriangle}

func main() {
	rand.Seed(time.Now().UnixNano())
	initConfig()

	BREAK_AFTER_NO_CHANGES := viper.GetInt("BREAK_AFTER_NO_CHANGES")
	CHECKPOINT_EVERY := viper.GetInt("CHECKPOINT_EVERY")
	COMPARE_AT_CHECKPOINTS := viper.GetBool("COMPARE_AT_CHECKPOINTS")
	COMPARE_AT_CHECKPOINTS_GAMES := viper.GetInt("COMPARE_AT_CHECKPOINTS_GAMES")
	EPSILON := viper.GetFloat64("EPSILON")
	EVALUATOR_GAMES := viper.GetInt("EVALUATOR_GAMES")
	FINISH_COMPARISON_GAMES := viper.GetInt("FINISH_COMPARISON_GAMES")
	GAMES_PER_ITERATION := viper.GetInt("GAMES_PER_ITERATION")
	INITIAL_WEIGHTS_PATH := viper.GetString("INITIAL_WEIGHTS_PATH")
	ITERATIONS := viper.GetInt("ITERATIONS")
	MAX_HISTORY_LENGTH := viper.GetInt("MAX_HISTORY_LENGTH")
	MCTS_SIMULATIONS := viper.GetInt("MCTS_SIMULATIONS")
	MINMAX_DEPTH := viper.GetInt("MINMAX_DEPTH")
	TRAINING_SIZE := viper.GetInt("TRAINING_SIZE")
	TRAINING_MODE := viper.GetString("TRAINING_MODE")
	RESULTS_DIR_NAME := viper.GetString("RESULTS_DIR_NAME")
	OLD_MINMAX_WEIGHTS_PATH := viper.GetString("OLD_MINMAX_WEIGHTS_PATH")
	OLD_MINMAX_MODE := viper.GetString("OLD_MINMAX_MODE")
	SGD_CONFIG := viper.Get("SGD_CONFIG").(map[string]float64)

	resultsDirPath := createResultsDir(RESULTS_DIR_NAME)
	configPath := path.Join(resultsDirPath, "config.yaml")
	viper.WriteConfigAs(configPath)

	initialWeights := INITIAL_WEIGHTS_BY_MODE[TRAINING_MODE]()
	if INITIAL_WEIGHTS_PATH != "" {
		initialWeights = reversiValueFns.LoadWeightsFromFile(INITIAL_WEIGHTS_PATH)
	}

	reversiToFeaturesFn := REVERSI_TO_FEATURES_BY_MODE[TRAINING_MODE]
	valueFn := reversiValueFns.CreateWeightedReversiFn(reversiToFeaturesFn, initialWeights)
	bestWeights := initialWeights
	bestPlayer := minMaxPlayer.NewMinMaxPlayer(MINMAX_DEPTH, valueFn)
	selfPlayPlayer := minMaxPlayer.NewEpsilonGreedyMinMaxPlayer(MINMAX_DEPTH, EPSILON, valueFn)

	oldWeights := reversiValueFns.LoadWeightsFromFile(OLD_MINMAX_WEIGHTS_PATH)
	oldValueFn := reversiValueFns.CreateWeightedReversiFn(REVERSI_TO_FEATURES_BY_MODE[OLD_MINMAX_MODE], oldWeights)
	oldMinMaxPlayer := minMaxPlayer.NewMinMaxPlayer(MINMAX_DEPTH, oldValueFn)

	mctsPlayer := monteCarloTreeSearchPlayer.NewMonteCarloTreeSearchPlayer(MCTS_SIMULATIONS)

	avgIterationTime := 0
	lastIterationChange := -1
	historyPositions := make([]*mat.VecDense, 0)
	historyYs := make([]float64, 0)

	for iteration := 0; iteration < ITERATIONS && iteration-lastIterationChange < BREAK_AFTER_NO_CHANGES; iteration++ {
		fmt.Printf("%d/%d (lastIterationChange=%d => %d/%d)\t%v left\n", iteration+1, ITERATIONS, lastIterationChange+1, iteration-lastIterationChange, BREAK_AFTER_NO_CHANGES, time.Duration(avgIterationTime*(ITERATIONS-iteration))*time.Nanosecond)

		iterationTimeStart := time.Now()
		results, totalPositions := runner.PlayNGames(reversiGameFactory, selfPlayPlayer, selfPlayPlayer, GAMES_PER_ITERATION)
		// utils.SaveGameResultsToFile(results, path.Join(resultsDirPath, iteration+"_results.txt")

		Xs, ys := gameResultsToXsAndys(reversiToFeaturesFn, results, totalPositions)
		historyPositions = append(historyPositions, Xs...)
		historyYs = append(historyYs, ys...)

		if len(historyPositions) > MAX_HISTORY_LENGTH {
			startIndex := len(historyPositions) - MAX_HISTORY_LENGTH
			historyPositions = historyPositions[startIndex:]
			historyYs = historyYs[startIndex:]
		}

		Xs, ys = chooseXsAndys(historyPositions, historyYs, TRAINING_SIZE)
		sgdResult := sgd.SGD(sgd.MeanSquaredErrorWithGradient, bestWeights, Xs, ys, SGD_CONFIG)

		newValueFn := reversiValueFns.CreateWeightedReversiFn(reversiToFeaturesFn, sgdResult.BestWeights)
		candidate := minMaxPlayer.NewMinMaxPlayer(MINMAX_DEPTH, newValueFn)
		candidateWins := runner.ComparePlayers(reversiGameFactory, candidate, bestPlayer, EVALUATOR_GAMES)

		fmt.Printf("New candidate won %d/%d games\n", candidateWins, EVALUATOR_GAMES)

		if float64(candidateWins)/float64(EVALUATOR_GAMES) > 0.5 {
			bestPlayer = candidate
			bestWeights = sgdResult.BestWeights
			selfPlayPlayer = minMaxPlayer.NewEpsilonGreedyMinMaxPlayer(MINMAX_DEPTH, EPSILON, newValueFn)
			lastIterationChange = iteration
		}

		if CHECKPOINT_EVERY > 0 && iteration > 0 && iteration%CHECKPOINT_EVERY == 0 {
			iterationString := strconv.Itoa(iteration)
			checkpointWeightsPath := path.Join(resultsDirPath, iterationString+"_weights.txt")

			reversiValueFns.SaveWeightsToFile(bestWeights, checkpointWeightsPath)

			if COMPARE_AT_CHECKPOINTS {
				resultsPath := path.Join(resultsDirPath, iterationString+"_results.txt")
				comparePlayersAndSaveResults(resultsPath, bestPlayer, "MinMax", mctsPlayer, "MCTS", COMPARE_AT_CHECKPOINTS_GAMES)
				comparePlayersAndSaveResults(resultsPath, bestPlayer, "MinMax", oldMinMaxPlayer, "OLD MinMax", COMPARE_AT_CHECKPOINTS_GAMES)
			}
		}

		avgIterationTime += (int(time.Since(iterationTimeStart)) - avgIterationTime) / (iteration + 1)
		fmt.Print("\n")
	}

	fmt.Println(bestWeights.RawVector().Data)
	bestWeightsPath := path.Join(resultsDirPath, "best_weights.txt")
	reversiValueFns.SaveWeightsToFile(bestWeights, bestWeightsPath)

	bestResultsPath := path.Join(resultsDirPath, "best_results.txt")
	comparePlayersAndSaveResults(bestResultsPath, bestPlayer, "MinMax", mctsPlayer, "MCTS", FINISH_COMPARISON_GAMES)
	comparePlayersAndSaveResults(bestResultsPath, bestPlayer, "MinMax", oldMinMaxPlayer, "OLD MinMax", FINISH_COMPARISON_GAMES)
}

// just a wrap
func reversiGameFactory() game.Game {
	return reversi.NewReversiGame()
}

func gameResultsToXsAndys(reversiToFeaturesFn reversiValueFns.ReversiToFeaturesFn, results []*runner.GameResult, totalPositions int) ([]*mat.VecDense, []float64) {
	Xs := make([]*mat.VecDense, totalPositions)
	ys := make([]float64, totalPositions)

	positionIndex := 0
	for _, gameResult := range results {
		for _, position := range gameResult.History {
			reversiGame := position.(*reversi.ReversiGame)
			Xs[positionIndex] = reversiToFeaturesFn(reversiGame)
			ys[positionIndex] = float64(gameResult.Winner) * 100.0
			positionIndex += 1
		}
	}

	return Xs, ys
}

func gameResultsToXsAndys2(reversiToFeaturesFn reversiValueFns.ReversiToFeaturesFn, results []*runner.GameResult, positionsPerGame int) ([]*mat.VecDense, []float64) {
	totalPositions := len(results) * positionsPerGame
	Xs := make([]*mat.VecDense, totalPositions)
	ys := make([]float64, totalPositions)

	for resultIndex, gameResult := range results {
		for i := 0; i < positionsPerGame; i++ {
			positionIndex := resultIndex*positionsPerGame + i
			position := gameResult.History[rand.Intn(len(gameResult.History))]
			reversiGame := position.(*reversi.ReversiGame)
			Xs[positionIndex] = reversiToFeaturesFn(reversiGame)
			ys[positionIndex] = float64(gameResult.Winner) * 100.0
		}
	}

	return Xs, ys
}

func chooseXsAndys(XsSource []*mat.VecDense, ysSource []float64, N int) ([]*mat.VecDense, []float64) {
	Xs := make([]*mat.VecDense, N)
	ys := make([]float64, N)

	for i := 0; i < N; i++ {
		index := rand.Intn(len(XsSource))
		Xs[i] = XsSource[index]
		ys[i] = ysSource[index]
	}

	return Xs, ys
}

func createResultsDir(resultsDirName string) string {
	if resultsDirName == "" {
		resultsDirName = utils.TimeNowString()
	}

	path := filepath.Join("./results/", resultsDirName)
	os.Mkdir(path, os.ModePerm)

	return path
}

func comparePlayersAndSaveResults(filePath string, player1 player.Player, player1Name string, player2 player.Player, player2Name string, numberOfGames int) {
	resultsFile, _ := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer resultsFile.Close()

	player1Wins := runner.ComparePlayers(reversiGameFactory, player1, player2, numberOfGames)
	resultsInfo := fmt.Sprintf("%s won %d/%d games versus %s\n", player1Name, player1Wins, numberOfGames, player2Name)
	fmt.Print(resultsInfo)
	fmt.Fprint(resultsFile, resultsInfo)
}
