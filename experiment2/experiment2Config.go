package main

import "github.com/spf13/viper"

func initConfig() {
	viper.SetDefault("CHECKPOINT_EVERY", 20)
	viper.SetDefault("COMPARE_AT_CHECKPOINTS", true)
	viper.SetDefault("COMPARE_AT_CHECKPOINTS_GAMES", 20)
	viper.SetDefault("EPSILON", 0.1)
	viper.SetDefault("EVALUATOR_GAMES", 16)
	viper.SetDefault("GAMES_PER_ITERATION", 50)
	viper.SetDefault("INITIAL_WEIGHTS_PATH", "")
	viper.SetDefault("MAX_HISTORY_LENGTH", 10000)
	viper.SetDefault("MCTS_SIMULATIONS", 2000)
	viper.SetDefault("MINMAX_DEPTH", 4)
	viper.SetDefault("TRAINING_SIZE", 512)
	viper.SetDefault("TRAINING_MODE", "normal") // "normal" | "triangle"
	viper.SetDefault("RESULTS_DIR_NAME", "")
	viper.SetDefault("OLD_MINMAX_WEIGHTS_PATH", "../experiment1/weights_2019-10-10 231145.txt")
	viper.SetDefault("OLD_MINMAX_WEIGHTS_MODE", "triangle")

	viper.SetDefault("SGD_CONFIG", map[string]float64{
		"alpha0":     1e-4,
		"alphaConst": 1e-5,
		"momentum":   0.1,
		"batch_size": 16,
		"max_epochs": 10000,
		"debug":      0})
}
