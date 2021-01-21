package main

import "github.com/spf13/viper"

func initConfig() {
	viper.SetDefault("CHECKPOINT_EVERY", 25)
	viper.SetDefault("COMPARE_AT_CHECKPOINTS", true)
	viper.SetDefault("COMPARE_AT_CHECKPOINTS_GAMES", 50)
	viper.SetDefault("EPSILON", 0.1)
	viper.SetDefault("EVALUATOR_GAMES", 20)
	viper.SetDefault("EVALUATOR_GAMES_AT_ONCE", 4)
	viper.SetDefault("GAMES_PER_ITERATION", 50)
	viper.SetDefault("GAME_TO_FEATURES_FN", "paddedmoves")
	viper.SetDefault("INITIAL_WEIGHTS_PATH", "")
	viper.SetDefault("MAX_BEST_PLAYERS_POOL_LENGTH", 10)
	viper.SetDefault("MCTS_SIMULATIONS", 2000)
	viper.SetDefault("MINMAX_DEPTH", 3)
	viper.SetDefault("RESULTS_DIR_NAME", "")
	viper.SetDefault("SELFPLAY_GAMES_AT_ONCE", 4)
	// viper.SetDefault("TRAINING_TRANSFORM_POSITIONS", true)
	// viper.SetDefault("OLD_MINMAX_MODEL_PATH", "./results/2020-11-06-174428-4-3-128-d/models/50")
	viper.SetDefault("OLD_MINMAX_MODEL_PATH", "./results/2020-11-19-225216-4-7-128-d25/models/300")
	viper.SetDefault("OLD_MINMAX_MODEL_GAME_TO_FEATURES_FN", "board3")

	viper.SetDefault("OPTIMIZER_BATCH_SIZE", 16)
	viper.SetDefault("OPTIMIZER_FLIP_POSITIONS_PROB", 0.5)
	viper.SetDefault("OPTIMIZER_LEARNING_RATE", 1e-4)
	viper.SetDefault("OPTIMIZER_MAX_EPOCHS", 200)
	viper.SetDefault("OPTIMIZER_TRAINING_SIZE", 256)
	viper.SetDefault("OPTIMIZER_MAX_POSITIONS_FROM_BATCH", 10)
	viper.SetDefault("OPTIMIZER_MAX_HISTORY_LENGTH", 25000)

	viper.SetDefault("OPTIMIZER_FULLY_CONNECTED", false)
	viper.SetDefault("OPTIMIZER_FC_LAYERS_COUNT", 7)
	viper.SetDefault("OPTIMIZER_FC_LAYER_UNITS", 128)
	viper.SetDefault("OPTIMIZER_FC_DROPOUT", 0.25)

	viper.SetDefault("OPTIMIZER_CONV_FILTERS", "[32,32,64,64]")
}
