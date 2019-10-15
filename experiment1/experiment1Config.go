package main

import (
	"github.com/spf13/viper"
)

func initConfig() {
	viper.SetDefault("BREAK_AFTER_NO_CHANGES", 50)
	viper.SetDefault("CHECKPOINT_EVERY", 100)
	viper.SetDefault("COMPARE_AT_CHECKPOINTS", true)
	viper.SetDefault("COMPARE_AT_CHECKPOINTS_GAMES", 20)
	viper.SetDefault("EPSILON", 0.1)
	viper.SetDefault("EVALUATOR_GAMES", 16)
	viper.SetDefault("FINISH_COMPARISON_GAMES", 100)
	viper.SetDefault("GAMES_PER_ITERATION", 20)
	viper.SetDefault("INITIAL_WEIGHTS_PATH", "./results/2019-10-14 200020/best_weights.txt")
	viper.SetDefault("ITERATIONS", 5000)
	viper.SetDefault("MAX_HISTORY_LENGTH", 30000)
	viper.SetDefault("MCTS_SIMULATIONS", 1000)
	viper.SetDefault("MINMAX_DEPTH", 4)
	viper.SetDefault("TRAINING_SIZE", 256)
	viper.SetDefault("TRAINING_MODE", "normal") // "normal" | "triangle"
	viper.SetDefault("RESULTS_DIR_NAME", "")
	viper.SetDefault("OLD_MINMAX_WEIGHTS_PATH", "./weights_2019-10-10 231145.txt")
	viper.SetDefault("OLD_MINMAX_MODE", "triangle")

	viper.SetDefault("SGD_CONFIG", map[string]float64{
		"alpha0":     5e-5,
		"alphaConst": 0,
		"momentum":   0.9,
		"batch_size": 16,
		"max_epochs": 10000,
		"debug":      1})
}
