package main

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"

	"github.com/jamOne-/kiwi-zero/edaxPlayer"
	"github.com/jamOne-/kiwi-zero/games"
	"github.com/jamOne-/kiwi-zero/policyPlayer"
	"github.com/jamOne-/kiwi-zero/randomPlayer"

	"github.com/jamOne-/kiwi-zero/predictor"
	"github.com/jamOne-/kiwi-zero/randomPredictor"

	tfpredictor "github.com/jamOne-/kiwi-zero/TFPredictor"
	"github.com/jamOne-/kiwi-zero/minMaxPlayer"
	"github.com/jamOne-/kiwi-zero/player"

	"github.com/jamOne-/kiwi-zero/game"
	"github.com/jamOne-/kiwi-zero/monteCarloTreeSearchPlayer"
	"github.com/jamOne-/kiwi-zero/reversiValueFns"
	"github.com/jamOne-/kiwi-zero/runner"
	"github.com/jamOne-/kiwi-zero/utils"

	"github.com/spf13/viper"
)

func getPlayerFactory(
	gameToFeaturesFn game.GameToFeaturesFn,
	selfPlay bool,
	playerType string,
) player.PlayerFactory {
	prefix := "EVALUATOR_"
	if selfPlay {
		prefix = "SELFPLAY_"
	}

	if playerType == "minmax" {
		depth := viper.GetInt(prefix + "MINMAX_DEPTH")
		return getMinMaxFactory(gameToFeaturesFn, depth)
	} else if playerType == "minmax-e" {
		depth := viper.GetInt(prefix + "MINMAX_DEPTH")
		return getEpsilonGreedyMinMaxFactory(gameToFeaturesFn, depth)
	} else if playerType == "minmax-sm" {
		depth := viper.GetInt(prefix + "MINMAX_DEPTH")
		return getSoftmaxMinMaxFactory(gameToFeaturesFn, depth)
	} else if playerType == "mcts-pred" {
		simulations := viper.GetInt(prefix + "MCTS_SIMULATIONS")
		rolloutDepth := viper.GetInt(prefix + "MCTS_ROLLOUT_DEPTH")
		policyRolloutPlayer := viper.GetBool(prefix + "MCTS_POLICY_ROLLOUT_PLAYER")
		optimizePolicy := viper.GetBool("OPTIMIZER_OPTIMIZE_POLICY")
		return getMCTSFactory(gameToFeaturesFn, simulations, rolloutDepth, policyRolloutPlayer, optimizePolicy)
	}

	return nil
}

func main() {
	rand.Seed(time.Now().UnixNano())
	initConfig()

	GAME := viper.GetString("GAME")
	GAME_TO_FEATURES_FN := viper.GetString("GAME_TO_FEATURES_FN")
	MCTS_SIMULATIONS := viper.GetInt("MCTS_SIMULATIONS")
	MINMAX_DEPTH := viper.GetInt("MINMAX_DEPTH")
	OLD_MINMAX_MODEL_PATH := viper.GetString("OLD_MINMAX_MODEL_PATH")
	OLD_MINMAX_MODEL_GAME_TO_FEATURES_FN := viper.GetString("OLD_MINMAX_MODEL_GAME_TO_FEATURES_FN")
	RESULTS_DIR_NAME := viper.GetString("RESULTS_DIR_NAME")

	resultsDirPath := createResultsDir(RESULTS_DIR_NAME)
	configPath := path.Join(resultsDirPath, "config.yaml")
	viper.WriteConfigAs(configPath)

	initialPredictor := randomPredictor.NewRandomPredictor()
	gameFactory := games.GAME_FACTORY_DICT[GAME]
	gameToFeaturesInfo := games.FEATURES_FNS_DICT[GAME][GAME_TO_FEATURES_FN]

	playersToCompareWith := make([]*PlayerToCompare, 0)
	// mctsPlayer := monteCarloTreeSearchPlayer.NewMonteCarloTreeSearchPlayer(MCTS_SIMULATIONS)
	mctsPlayer := monteCarloTreeSearchPlayer.NewGeneralMCTSPlayer(MCTS_SIMULATIONS, 2.0, 99, randomPlayer.NewRandomPlayer(), nil, nil, nil)
	playersToCompareWith = append(playersToCompareWith, &PlayerToCompare{fmt.Sprintf("MCTS (%d sims)", MCTS_SIMULATIONS), mctsPlayer})
	playersToCompareWith = append(playersToCompareWith, &PlayerToCompare{"Random", randomPlayer.NewRandomPlayer()})

	if OLD_MINMAX_MODEL_PATH != "" {
		oldGameToFeaturesFn := games.FEATURES_FNS_DICT[GAME][OLD_MINMAX_MODEL_GAME_TO_FEATURES_FN].Fn
		oldMinMaxPlayer := loadMinMaxPlayer(oldGameToFeaturesFn, OLD_MINMAX_MODEL_PATH, MINMAX_DEPTH)
		playersToCompareWith = append(playersToCompareWith, &PlayerToCompare{"OLD MinMax", oldMinMaxPlayer})
	}

	var selfPlayTeacherFactory runner.NewPlayerFactory = nil
	if viper.GetString("SELFPLAY_TEACHER") == "edax" {
		selfPlayTeacherFactory = getEdaxRunnerFactory(-64, 64, viper.GetInt("SELFPLAY_TEACHER_EDAX_DEPTH"), 100)
	} else if viper.GetString("SELFPLAY_TEACHER") == "mcts" {
		selfPlayTeacherFactory = runner.FactorizePlayer(monteCarloTreeSearchPlayer.NewGeneralMCTSPlayer(viper.GetInt("SELFPLAY_TEACHER_MCTS_SIMULATIONS"), 2.0, 99, randomPlayer.NewRandomPlayer(), nil, nil, nil))
	}

	selfPlayPlayerFactory := getPlayerFactory(gameToFeaturesInfo.Fn, true, viper.GetString("SELFPLAY_PLAYER_TYPE"))
	evaluatorPlayerFactory := getPlayerFactory(gameToFeaturesInfo.Fn, false, viper.GetString("EVALUATOR_PLAYER_TYPE"))

	policyLength := (gameFactory()).GetMaxPossibleMoves()

	gameResultsChan := make(chan *runner.GameResultsBatch)
	bestPredictorsChan := make(chan predictor.Predictor)
	newPredictorsChan := make(chan predictor.Predictor)

	go SelfPlayLoop(bestPredictorsChan, gameResultsChan, gameFactory, initialPredictor, selfPlayPlayerFactory, selfPlayTeacherFactory)
	go Optimizer(gameResultsChan, newPredictorsChan, gameToFeaturesInfo.Fn, gameToFeaturesInfo.Shape, policyLength, resultsDirPath)
	go Evaluator(newPredictorsChan, bestPredictorsChan, gameFactory, initialPredictor, evaluatorPlayerFactory, playersToCompareWith, resultsDirPath)
	// bestValueFnsChan <- initialValueFn
	bestPredictorsChan <- initialPredictor

	// f, err := os.Create("cpuboard1features.prof")
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

func getMinMaxFactory(gameToFeatures game.GameToFeaturesFn, depth int) player.PlayerFactory {
	return func(predictor predictor.Predictor) player.Player {
		valueFn := reversiValueFns.CreateMinMaxValueFn(gameToFeatures, predictor)
		return minMaxPlayer.NewMinMaxPlayer(depth, valueFn)
	}
}

func getSoftmaxMinMaxFactory(gameToFeatures game.GameToFeaturesFn, depth int) player.PlayerFactory {
	return func(predictor predictor.Predictor) player.Player {
		valueFn := reversiValueFns.CreateMinMaxValueFn(gameToFeatures, predictor)
		return minMaxPlayer.NewSoftMaxMinMaxPlayer(depth, valueFn)
	}
}

func getEpsilonGreedyMinMaxFactory(gameToFeatures game.GameToFeaturesFn, depth int) player.PlayerFactory {
	return func(predictor predictor.Predictor) player.Player {
		valueFn := reversiValueFns.CreateMinMaxValueFn(gameToFeatures, predictor)
		return minMaxPlayer.NewEpsilonGreedyMinMaxPlayer(depth, viper.GetFloat64("EPSILON"), valueFn)
	}
}

func getMCTSFactory(gameToFeatures game.GameToFeaturesFn, maxSimulations int, rolloutDepth int, policyRolloutPlayer bool, optimizePolicy bool) player.PlayerFactory {
	return func(pred predictor.Predictor) player.Player {
		valueFn := reversiValueFns.CreateMinMaxValueFn(gameToFeatures, pred)

		var valueAndPolicyFn monteCarloTreeSearchPlayer.ValueAndPolicyFn = nil
		if optimizePolicy {
			valueAndPolicyFn = predictor.CreateValueAndPolicyFn(gameToFeatures, pred)
		}

		var rolloutPlayer player.Player = randomPlayer.NewRandomPlayer()
		if policyRolloutPlayer {
			distributionFn := policyPlayer.GameToDistributionFnFromPredictor(gameToFeatures, pred)
			rolloutPlayer = policyPlayer.NewPolicyPlayer(distributionFn)
		}

		return monteCarloTreeSearchPlayer.NewGeneralMCTSPlayer(
			maxSimulations,
			2.0,
			rolloutDepth,
			rolloutPlayer,
			valueFn,
			valueAndPolicyFn,
			nil,
		)
	}
}

func getEdaxRunnerFactory(alpha float64, beta float64, depth int, probcut int) runner.NewPlayerFactory {
	return func() player.Player {
		return edaxPlayer.NewEdaxPlayer(alpha, beta, depth, probcut)
	}
}
