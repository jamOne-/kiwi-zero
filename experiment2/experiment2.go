package main

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"

	"github.com/jamOne-/kiwi-zero/game"
	"github.com/jamOne-/kiwi-zero/monteCarloTreeSearchPlayer"
	"github.com/jamOne-/kiwi-zero/reversi"
	"github.com/jamOne-/kiwi-zero/reversiValueFns"
	"github.com/jamOne-/kiwi-zero/runner"
	"github.com/jamOne-/kiwi-zero/utils"

	"github.com/spf13/viper"
)

// var INITIAL_WEIGHTS_BY_MODE = map[string](func() *mat.VecDense){
// 	"normal":   reversiValueFns.GetInitialWeights,
// 	"triangle": reversiValueFns.GetTriangleInitialWeights,
// 	"extended": reversiValueFns.GetExtendedInitialWeights}

// var REVERSI_TO_FEATURES_BY_MODE = map[string]game.GameToFeaturesFn{
// 	"normal":   reversiValueFns.ConvertToReversiFn(reversiValueFns.ReversiToFeatures),
// 	"triangle": reversiValueFns.ConvertToReversiFn(reversiValueFns.ReversiToFeaturesTriangle),
// 	"extended": reversiValueFns.ConvertToReversiFn(reversiValueFns.ReversiToFeaturesExtended)}

func main() {
	rand.Seed(time.Now().UnixNano())
	initConfig()

	// INITIAL_WEIGHTS_PATH := viper.GetString("INITIAL_WEIGHTS_PATH")
	MCTS_SIMULATIONS := viper.GetInt("MCTS_SIMULATIONS")
	// MINMAX_DEPTH := viper.GetInt("MINMAX_DEPTH")
	// OLD_MINMAX_WEIGHTS_PATH := viper.GetString("OLD_MINMAX_WEIGHTS_PATH")
	// OLD_MINMAX_WEIGHTS_MODE := viper.GetString("OLD_MINMAX_WEIGHTS_MODE")
	RESULTS_DIR_NAME := viper.GetString("RESULTS_DIR_NAME")
	// TRAINING_MODE := viper.GetString("TRAINING_MODE")

	resultsDirPath := createResultsDir(RESULTS_DIR_NAME)
	configPath := path.Join(resultsDirPath, "config.yaml")
	viper.WriteConfigAs(configPath)

	// initialWeights := INITIAL_WEIGHTS_BY_MODE[TRAINING_MODE]()
	// if INITIAL_WEIGHTS_PATH != "" {
	// 	initialWeights = reversiValueFns.LoadWeightsFromFile(INITIAL_WEIGHTS_PATH)
	// }

	initialValueFn := getInitialValueFn()

	playersToCompareWith := make([]*PlayerToCompare, 0)

	mctsPlayer := monteCarloTreeSearchPlayer.NewMonteCarloTreeSearchPlayer(MCTS_SIMULATIONS)
	playersToCompareWith = append(playersToCompareWith, &PlayerToCompare{fmt.Sprintf("MCTS (%d sims)", MCTS_SIMULATIONS), mctsPlayer})

	// TODO: Add support for loading saved versions of players to compare with
	// if OLD_MINMAX_WEIGHTS_PATH != "" {
	// 	oldMinMaxWeights := reversiValueFns.LoadWeightsFromFile(OLD_MINMAX_WEIGHTS_PATH)
	// 	oldMinMaxValueFn := reversiValueFns.CreateWeightedReversiFn(REVERSI_TO_FEATURES_BY_MODE[OLD_MINMAX_WEIGHTS_MODE], oldMinMaxWeights)
	// 	oldMinMaxPlayer := minMaxPlayer.NewMinMaxPlayer(MINMAX_DEPTH, oldMinMaxValueFn)

	// 	playersToCompareWith = append(playersToCompareWith, &PlayerToCompare{"OLD MinMax", oldMinMaxPlayer})
	// }

	bestValueFnsChan := make(chan game.ValueFn)
	gameResultsChan := make(chan *runner.GameResultsBatch)
	newValueFnsChan := make(chan game.ValueFn)
	// reversiToFeaturesFn := REVERSI_TO_FEATURES_BY_MODE[TRAINING_MODE]
	gameToFeaturesFn := reversiValueFns.ConvertReversiFnToGeneralFeatuersFn(reversiValueFns.ReversiToOneHotBoard)

	go SelfPlayLoop(bestValueFnsChan, gameResultsChan, reversiGameFactory, initialValueFn)
	go Optimizer(gameResultsChan, newValueFnsChan, gameToFeaturesFn, resultsDirPath)
	go Evaluator(newValueFnsChan, bestValueFnsChan, reversiGameFactory, initialValueFn, playersToCompareWith, resultsDirPath)

	var waitGroup sync.WaitGroup
	waitGroup.Add(1)
	waitGroup.Wait()
}

func reversiGameFactory() game.Game {
	return reversi.NewReversiGame()
}

func createResultsDir(resultsDirName string) string {
	if resultsDirName == "" {
		resultsDirName = utils.TimeNowString()
	}

	path := filepath.Join("./results/", resultsDirName)
	os.Mkdir(path, os.ModePerm)

	return path
}

func getInitialValueFn() game.ValueFn {
	// TODO
	return func(game game.Game) float64 {
		return 0.5
	}
}
