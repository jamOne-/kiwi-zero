package main

import "github.com/spf13/viper"

func initConfig() {
	viper.SetDefault("CHECKPOINT_EVERY", 100)
	viper.SetDefault("COMPARE_AT_CHECKPOINTS", true)
	viper.SetDefault("COMPARE_AT_CHECKPOINTS_GAMES", 50)
	viper.SetDefault("EPSILON", 0.1)
	viper.SetDefault("EVALUATOR_GAMES", 16)
	viper.SetDefault("GAMES_PER_ITERATION", 50)
	viper.SetDefault("INITIAL_WEIGHTS_PATH", "")
	viper.SetDefault("MAX_HISTORY_LENGTH", 100000)
	viper.SetDefault("MCTS_SIMULATIONS", 1000)
	viper.SetDefault("MINMAX_DEPTH", 4)
	viper.SetDefault("RESULTS_DIR_NAME", "")
	viper.SetDefault("SELFPLAY_GAMES_AT_ONCE", 15)
	viper.SetDefault("TRAINING_FLIP_POSITIONS_PROB", 0.0)
	viper.SetDefault("TRAINING_MODE", "extended") // "normal" | "triangle" | "extended"
	viper.SetDefault("TRAINING_SIZE", 512)
	viper.SetDefault("TRAINING_TRANSFORM_POSITIONS", true)
	// viper.SetDefault("OLD_MINMAX_WEIGHTS_PATH", "../experiment1/weights_2019-10-10 231145.txt")
	viper.SetDefault("OLD_MINMAX_WEIGHTS_PATH", "./results/2019-12-14 151304/100_weights.txt")
	// viper.SetDefault("OLD_MINMAX_WEIGHTS_MODE", "triangle")
	viper.SetDefault("OLD_MINMAX_WEIGHTS_MODE", "extended")

	viper.SetDefault("SGD_CONFIG", map[string]float64{
		"alpha0":        1e-4,
		"alphaConst":    1e-5,
		"momentum":      0.9,
		"batch_size":    16,
		"max_epochs":    1000,
		"weights_decay": 0,
		"debug":         0})
}
