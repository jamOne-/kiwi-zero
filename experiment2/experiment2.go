package main

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"

	tfpredictor "github.com/jamOne-/kiwi-zero/TFPredictor"
	"github.com/jamOne-/kiwi-zero/minMaxPlayer"
	"github.com/jamOne-/kiwi-zero/player"

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

var REVERSI_TO_FEATURES_FN_DICT = map[string]game.GameToFeaturesFn{
	"board3":     reversiValueFns.ConvertReversiFnToGeneralFeatuersFn(reversiValueFns.ReversiToOneHotBoard3),
	"boardmoves": reversiValueFns.ConvertReversiFnToGeneralFeatuersFn(reversiValueFns.ReversiToOneHotBoardMoves),
}

func main() {
	rand.Seed(time.Now().UnixNano())
	initConfig()

	// INITIAL_WEIGHTS_PATH := viper.GetString("INITIAL_WEIGHTS_PATH")
	GAME_TO_FEATURES_FN := viper.GetString("GAME_TO_FEATURES_FN")
	MCTS_SIMULATIONS := viper.GetInt("MCTS_SIMULATIONS")
	MINMAX_DEPTH := viper.GetInt("MINMAX_DEPTH")
	OLD_MINMAX_MODEL_PATH := viper.GetString("OLD_MINMAX_MODEL_PATH")
	OLD_MINMAX_MODEL_GAME_TO_FEATURES_FN := viper.GetString("OLD_MINMAX_MODEL_GAME_TO_FEATURES_FN")
	RESULTS_DIR_NAME := viper.GetString("RESULTS_DIR_NAME")

	resultsDirPath := createResultsDir(RESULTS_DIR_NAME)
	configPath := path.Join(resultsDirPath, "config.yaml")
	viper.WriteConfigAs(configPath)

	// initialWeights := INITIAL_WEIGHTS_BY_MODE[TRAINING_MODE]()
	// if INITIAL_WEIGHTS_PATH != "" {
	// 	initialWeights = reversiValueFns.LoadWeightsFromFile(INITIAL_WEIGHTS_PATH)
	// }

	initialValueFn := getInitialValueFn()
	gameToFeaturesFn := REVERSI_TO_FEATURES_FN_DICT[GAME_TO_FEATURES_FN]

	playersToCompareWith := make([]*PlayerToCompare, 0)
	mctsPlayer := monteCarloTreeSearchPlayer.NewMonteCarloTreeSearchPlayer(MCTS_SIMULATIONS)
	playersToCompareWith = append(playersToCompareWith, &PlayerToCompare{fmt.Sprintf("MCTS (%d sims)", MCTS_SIMULATIONS), mctsPlayer})

	if OLD_MINMAX_MODEL_PATH != "" {
		oldGameToFeaturesFn := REVERSI_TO_FEATURES_FN_DICT[OLD_MINMAX_MODEL_GAME_TO_FEATURES_FN]
		oldMinMaxPlayer := loadMinMaxPlayer(oldGameToFeaturesFn, OLD_MINMAX_MODEL_PATH, MINMAX_DEPTH)
		playersToCompareWith = append(playersToCompareWith, &PlayerToCompare{"OLD MinMax", oldMinMaxPlayer})
	}

	bestValueFnsChan := make(chan game.ValueFn)
	gameResultsChan := make(chan *runner.GameResultsBatch)
	newValueFnsChan := make(chan game.ValueFn)

	go SelfPlayLoop(bestValueFnsChan, gameResultsChan, reversiGameFactory, initialValueFn)
	go Optimizer(gameResultsChan, newValueFnsChan, gameToFeaturesFn, resultsDirPath)
	go Evaluator(newValueFnsChan, bestValueFnsChan, reversiGameFactory, initialValueFn, playersToCompareWith, resultsDirPath)
	bestValueFnsChan <- initialValueFn

	// f, err := os.Create("cpu.prof")
	// if err != nil {
	// 	log.Fatal("could not create CPU profile: ", err)
	// }
	// if err := pprof.StartCPUProfile(f); err != nil {
	// 	log.Fatal("could not start CPU profile: ", err)
	// }

	// time.Sleep(20 * time.Minute)
	// pprof.StopCPUProfile()
	// f.Close()

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

func loadMinMaxPlayer(gameToFeatures game.GameToFeaturesFn, path string, depth int) player.Player {
	predictor := tfpredictor.NewTFPredictor(path)
	valueFn := reversiValueFns.CreateMinMaxValueFn(gameToFeatures, predictor)

	return minMaxPlayer.NewMinMaxPlayer(depth, valueFn)
}
