package main

import "github.com/spf13/viper"

func initConfig() {
	viper.SetDefault("CHECKPOINT_EVERY", 50)
	viper.SetDefault("COMPARE_AT_CHECKPOINTS", true)
	viper.SetDefault("COMPARE_AT_CHECKPOINTS_GAMES", 50)
	viper.SetDefault("EPSILON", 0.1)
	viper.SetDefault("EVALUATOR_GAMES", 40)
	viper.SetDefault("EVALUATOR_GAMES_AT_ONCE", 4)
	viper.SetDefault("GAMES_PER_ITERATION", 50)
	viper.SetDefault("INITIAL_WEIGHTS_PATH", "")
	viper.SetDefault("MAX_BEST_PLAYERS_POOL_LENGTH", 20)
	viper.SetDefault("MAX_HISTORY_LENGTH", 25000)
	viper.SetDefault("MCTS_SIMULATIONS", 2000)
	viper.SetDefault("MINMAX_DEPTH", 4)
	viper.SetDefault("RESULTS_DIR_NAME", "")
	viper.SetDefault("SELFPLAY_GAMES_AT_ONCE", 4)
	viper.SetDefault("TRAINING_MODE", "extended") // "normal" | "triangle" | "extended"
	// viper.SetDefault("TRAINING_TRANSFORM_POSITIONS", true)
	// viper.SetDefault("OLD_MINMAX_MODEL_PATH", "./results/2020-11-06-174428-4-3-128-d/models/50")

	viper.SetDefault("OPTIMIZER_BATCH_SIZE", 16)
	viper.SetDefault("OPTIMIZER_FLIP_POSITIONS_PROB", 0.5)
	viper.SetDefault("OPTIMIZER_INPUT_SHAPE", "(8,8,3)")
	viper.SetDefault("OPTIMIZER_LEARNING_RATE", 1e-4)
	viper.SetDefault("OPTIMIZER_MAX_EPOCHS", 1000)
	viper.SetDefault("OPTIMIZER_TRAINING_SIZE", 256)

	viper.SetDefault("OPTIMIZER_FULLY_CONNECTED", true)
	viper.SetDefault("OPTIMIZER_FC_LAYERS_COUNT", 7)
	viper.SetDefault("OPTIMIZER_FC_LAYER_UNITS", 128)
	viper.SetDefault("OPTIMIZER_FC_DROPOUT", 0.25)
}
