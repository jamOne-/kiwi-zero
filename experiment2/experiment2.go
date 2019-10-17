package main

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"

	"github.com/jamOne-/kiwi-zero/utils"

	"github.com/jamOne-/kiwi-zero/game"
	"github.com/jamOne-/kiwi-zero/minMaxPlayer"
	"github.com/jamOne-/kiwi-zero/reversi"

	"github.com/jamOne-/kiwi-zero/monteCarloTreeSearchPlayer"

	"github.com/jamOne-/kiwi-zero/reversiValueFns"
	"github.com/jamOne-/kiwi-zero/runner"
	"github.com/spf13/viper"
	"gonum.org/v1/gonum/mat"
)

var REVERSI_TO_FEATURES_BY_MODE = map[string]reversiValueFns.ReversiToFeaturesFn{
	"normal":   reversiValueFns.ReversiToFeatures,
	"triangle": reversiValueFns.ReversiToFeaturesTriangle}

func main() {
	rand.Seed(time.Now().UnixNano())
	initConfig()

	INITIAL_WEIGHTS_PATH := viper.GetString("INITIAL_WEIGHTS_PATH")
	MCTS_SIMULATIONS := viper.GetInt("MCTS_SIMULATIONS")
	MINMAX_DEPTH := viper.GetInt("MINMAX_DEPTH")
	OLD_MINMAX_WEIGHTS_PATH := viper.GetString("OLD_MINMAX_WEIGHTS_PATH")
	OLD_MINMAX_WEIGHTS_MODE := viper.GetString("OLD_MINMAX_WEIGHTS_MODE")
	RESULTS_DIR_NAME := viper.GetString("RESULTS_DIR_NAME")

	resultsDirPath := createResultsDir(RESULTS_DIR_NAME)
	configPath := path.Join(resultsDirPath, "config.yaml")
	viper.WriteConfigAs(configPath)

	initialWeights := reversiValueFns.GetInitialWeights()
	if INITIAL_WEIGHTS_PATH != "" {
		initialWeights = reversiValueFns.LoadWeightsFromFile(INITIAL_WEIGHTS_PATH)
	}

	playersToCompareWith := make([]*PlayerToCompare, 0)

	mctsPlayer := monteCarloTreeSearchPlayer.NewMonteCarloTreeSearchPlayer(MCTS_SIMULATIONS)
	playersToCompareWith = append(playersToCompareWith, &PlayerToCompare{fmt.Sprintf("MCTS (%d sims)", MCTS_SIMULATIONS), mctsPlayer})

	if OLD_MINMAX_WEIGHTS_PATH != "" {
		oldMinMaxWeights := reversiValueFns.LoadWeightsFromFile(OLD_MINMAX_WEIGHTS_PATH)
		oldMinMaxValueFn := reversiValueFns.CreateWeightedReversiFn(REVERSI_TO_FEATURES_BY_MODE[OLD_MINMAX_WEIGHTS_MODE], oldMinMaxWeights)
		oldMinMaxPlayer := minMaxPlayer.NewMinMaxPlayer(MINMAX_DEPTH, oldMinMaxValueFn)

		playersToCompareWith = append(playersToCompareWith, &PlayerToCompare{"OLD MinMax", oldMinMaxPlayer})
	}

	bestWeightsChan := make(chan *mat.VecDense)
	gameResultsChan := make(chan *runner.GameResultsBatch)
	newWeightsChan := make(chan *mat.VecDense)

	go SelfPlayLoop(bestWeightsChan, gameResultsChan, initialWeights, reversiGameFactory)
	go Optimizer(gameResultsChan, newWeightsChan, initialWeights)
	go Evaluator(newWeightsChan, bestWeightsChan, reversiGameFactory, initialWeights, playersToCompareWith, resultsDirPath)

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
