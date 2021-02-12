package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jamOne-/kiwi-zero/predictor"
	"github.com/jamOne-/kiwi-zero/reversiValueFns"

	"github.com/spf13/viper"

	tfpredictor "github.com/jamOne-/kiwi-zero/TFPredictor"
	"github.com/jamOne-/kiwi-zero/game"
	"github.com/jamOne-/kiwi-zero/reversi"
	"github.com/jamOne-/kiwi-zero/runner"
	"github.com/jamOne-/kiwi-zero/utils"
)

type GamePosition struct {
	gameId   string
	position game.Game
	policy   []float32
	winner   float64
}

type trainingParams struct {
	Xs       []game.Features
	ys       []float64
	policies [][]float32
}

func Optimizer(
	gameResultsChan chan *runner.GameResultsBatch,
	predictorsChan chan predictor.Predictor,
	gameToFeaturesFn game.GameToFeaturesFn,
	resultsDirPath string,
) {
	MAX_HISTORY_LENGTH := viper.GetInt("OPTIMIZER_MAX_HISTORY_LENGTH")
	MAX_POSITIONS_FROM_BATCH := viper.GetInt("OPTIMIZER_MAX_POSITIONS_FROM_BATCH")
	FITS_PER_ITERATION := viper.GetInt("OPTIMIZER_FITS_PER_ITERATION")
	TRAINING_SIZE := viper.GetInt("OPTIMIZER_TRAINING_SIZE")
	TRAINING_SET_SAME_GAMES_ALLOWED := viper.GetBool("OPTIMIZER_TRAINING_SET_SAME_GAMES_ALLOWED")
	FLIP_POSITIONS_PROB := viper.GetFloat64("OPTIMIZER_FLIP_POSITIONS_PROB")
	// TRAINING_TRANSFORM_POSITIONS := viper.GetBool("TRAINING_TRANSFORM_POSITIONS")

	// Create models directory for current run
	modelsDirPath := filepath.Join(resultsDirPath, "models")
	os.Mkdir(modelsDirPath, os.ModePerm)

	gamePositions := make([]*GamePosition, 0)

	pythonOptimizerCmd := exec.Command(
		"python3", "../python/optimizer/optimizer.py",
		"--models_directory", modelsDirPath,
		"--input_shape", reversiValueFns.FEATURES_FN_TO_SHAPE_DICT[viper.GetString("GAME_TO_FEATURES_FN")],
		"--learning_rate", fmt.Sprintf("%f", viper.GetFloat64("OPTIMIZER_LEARNING_RATE")),
		"--momentum", fmt.Sprintf("%f", viper.GetFloat64("OPTIMIZER_MOMENTUM")),
		"--epochs", strconv.Itoa(viper.GetInt("OPTIMIZER_MAX_EPOCHS")),
		"--batch_size", strconv.Itoa(viper.GetInt("OPTIMIZER_BATCH_SIZE")),
		"--regularizer_const", fmt.Sprintf("%f", viper.GetFloat64("OPTIMIZER_REGULARIZER_CONST")),
		"--optimize_policy", strconv.Itoa(utils.BoolToInt(viper.GetBool("OPTIMIZER_OPTIMIZE_POLICY"))),
		"--fully_connected", strconv.Itoa(utils.BoolToInt(viper.GetBool("OPTIMIZER_FULLY_CONNECTED"))),
		"--fc_layers_count", strconv.Itoa(viper.GetInt("OPTIMIZER_FC_LAYERS_COUNT")),
		"--fc_layer_units", strconv.Itoa(viper.GetInt("OPTIMIZER_FC_LAYER_UNITS")),
		"--fc_dropout", fmt.Sprintf("%f", viper.GetFloat64("OPTIMIZER_FC_DROPOUT")),
		"--conv_filters", viper.GetString("OPTIMIZER_CONV_FILTERS"),
	)

	fmt.Println("Args: " + strings.Join(pythonOptimizerCmd.Args, " "))

	optimizerIn, _ := pythonOptimizerCmd.StdinPipe()
	optimizerOut, _ := pythonOptimizerCmd.StdoutPipe()

	trainingChan := make(chan *trainingParams)
	newModelPathChan := make(chan string)
	go trainer(trainingChan, optimizerIn, optimizerOut, newModelPathChan)
	go predictorsCreator(newModelPathChan, predictorsChan, gameToFeaturesFn)

	pythonOptimizerCmd.Start()

	optimizerIteration := 1

	for {
		select {
		case batch := <-gameResultsChan:
			results, totalPositions := batch.Results, batch.TotalPositions
			positions := splitResults(results, totalPositions, optimizerIteration)

			// if TRAINING_TRANSFORM_POSITIONS {
			// 	transformPositions(positions)
			// }

			if FLIP_POSITIONS_PROB > 0 {
				flipPositionsColors(FLIP_POSITIONS_PROB, positions)
			}

			if MAX_POSITIONS_FROM_BATCH != -1 {
				positions = choosePositions(positions, true, MAX_POSITIONS_FROM_BATCH*len(results))
			}

			gamePositions = append(gamePositions, positions...)

			if len(gamePositions) > MAX_HISTORY_LENGTH {
				startIndex := len(gamePositions) - MAX_HISTORY_LENGTH
				gamePositions = gamePositions[startIndex:]
			}

			positions = choosePositions(gamePositions, TRAINING_SET_SAME_GAMES_ALLOWED, FITS_PER_ITERATION*TRAINING_SIZE)
			params := positionsToTrainingParams(gameToFeaturesFn, positions)
			trainingChan <- params

			optimizerIteration += 1

			// select {
			// case trainingChan <- params:
			// 	// try to send

			// default:
			// 	// else skip
			// }
		}
	}
}

func choosePositions(positions []*GamePosition, sameGameAllowed bool, N int) []*GamePosition {
	choice := make([]*GamePosition, N)
	gameChosen := make(map[string]bool)

	for i := 0; i < N; {
		index := rand.Intn(len(positions))
		position := positions[index]

		if _, exists := gameChosen[position.gameId]; !exists || sameGameAllowed {
			gameChosen[position.gameId] = true
			choice[i] = position
			i += 1
		}
	}

	return choice
}

func positionsToTrainingParams(gameToFeaturesFn game.GameToFeaturesFn, positions []*GamePosition) *trainingParams {
	Xs := make([]game.Features, len(positions))
	ys := make([]float64, len(positions))
	policies := make([][]float32, len(positions))

	for i, position := range positions {
		Xs[i] = gameToFeaturesFn(position.position)
		ys[i] = (position.winner + 1) / 2
		policies[i] = position.policy
	}

	return &trainingParams{Xs, ys, policies}
}

func splitResults(results []*runner.GameResult, totalPositions int, iteration int) []*GamePosition {
	gamePositions := make([]*GamePosition, totalPositions)

	index := 0
	for gameIndex, result := range results {
		winner := float64(result.Winner)

		for _, tuple := range result.History {
			game := tuple.Game
			policy := tuple.Policy
			id := fmt.Sprintf("%d-%d", iteration, gameIndex)

			gamePositions[index] = &GamePosition{id, game, policy, winner}
			index += 1
		}
	}

	return gamePositions
}

func randomPositionTransformation(position game.Game) {
	TRANSFORMATIONS := 4 + 4
	reversiPosition := position.(*reversi.ReversiGame) // todo: game.GetBoard()
	transformation := rand.Intn(TRANSFORMATIONS)

	switch transformation {
	case 0:
		fallthrough
	case 1:
		fallthrough
	case 2:
		fallthrough
	case 3:
		utils.RotateSquareVector(reversiPosition.Board, transformation)

	case 4:
		utils.PerformSymmetryVector1(reversiPosition.Board)
	case 5:
		utils.PerformSymmetryVector2(reversiPosition.Board)
	case 6:
		utils.PerformSymmetryVector3(reversiPosition.Board)
	case 7:
		utils.PerformSymmetryVector4(reversiPosition.Board)
	}
}

func transformPositions(positions []game.Game) {
	for _, position := range positions {
		randomPositionTransformation(position)
	}
}

func flipPositionsColors(flipProb float64, positions []*GamePosition) {
	for _, position := range positions {
		if rand.Float64() < flipProb {
			flipGameColors(position.position)
			position.winner *= -1
		}
	}
}

func flipGameColors(g game.Game) {
	reversiGame := g.(*reversi.ReversiGame) // todo: game.GetBoard(), SetCurrentPlayer
	reversiGame.Turn *= -1

	for i, color := range reversiGame.Board {
		reversiGame.Board[i] = color * -1
	}
}

func trainer(
	paramsChan chan *trainingParams,
	optimizerIn io.WriteCloser,
	optimizerOut io.ReadCloser,
	newModelPathChan chan string,
) {
	optimizerOutReader := bufio.NewReader(optimizerOut)
	training_i := 1

	for params := range paramsChan {
		Xs, ys, policies := params.Xs, params.ys, params.policies
		Xs_shape := [4]int{len(Xs), len(Xs[0]), len(Xs[0][0]), len(Xs[0][0][0])}

		optimizerIn.Write([]byte(fmt.Sprintf("%v %v %v %v\n", Xs_shape[0], Xs_shape[1], Xs_shape[2], Xs_shape[3])))

		for _, XXX := range Xs {
			for _, XX := range XXX {
				for _, X := range XX {
					line := utils.FloatsToString(X) + "\n"
					optimizerIn.Write([]byte(line))
				}
			}
		}

		for _, y := range ys {
			line := fmt.Sprintf("%v\n", y)
			optimizerIn.Write([]byte(line))
		}

		for _, policy := range policies {
			line := utils.FloatsToString(policy) + "\n"
			optimizerIn.Write([]byte(line))
		}

		// fmt.Printf("Optimizer (%d): training...\n", training_i)

		trainingSummary := ""
		SUMMARY_LINES := 5
		for i := 0; i < SUMMARY_LINES; i++ {
			summaryLine, _ := optimizerOutReader.ReadString('\n')
			trainingSummary += summaryLine
		}

		newModelPath, _ := optimizerOutReader.ReadString('\n')
		newModelPath = strings.TrimSpace(newModelPath)

		// newWeightsString, _ := optimizerOutReader.ReadString('\n')
		// newWeights := parseWeights(newWeightsString)

		fmt.Printf("Optimizer (%d): %s", training_i, trainingSummary)
		newModelPathChan <- newModelPath
		training_i += 1

		// select {
		// case newModelPathChan <- newModelPath:
		// 	// try to send

		// default:
		// 	// else skip
		// }

		// sgdResult := sgd.SGD(sgd.MeanSquaredErrorWithGradient, bestWeights, Xs, ys, SGD_CONFIG)
		// bestWeights = sgdResult.BestWeights

		// fmt.Printf("Optimizer (%d): totalEpochs=%d, bestEpoch=%d, testSetError=%.2f, validationSetError=%.2f, trainingError=%.2f\n", training_i, sgdResult.TotalEpochs, sgdResult.BestWeightsEpoch, sgdResult.TestSetErrorRate, sgdResult.BestValidErrorRate, sgdResult.TrainErrorsHistory[len(sgdResult.TrainErrorsHistory)-1])
	}
}

func predictorsCreator(
	modelPathChan chan string,
	predictorsChan chan predictor.Predictor,
	gameToFeatures game.GameToFeaturesFn,
) {
	for path := range modelPathChan {
		if path == "" {
			continue
		}

		// fmt.Printf("Optimizer: saved model path=%s", path)

		predictor := tfpredictor.NewTFPredictor(path)
		// valueFn := reversiValueFns.CreateMinMaxValueFn(gameToFeatures, predictor)

		predictorsChan <- predictor
		// valueFns <- valueFn
		// select {
		// case valueFns <- valueFn:
		// 	// try to send

		// default:
		// 	// else skip
		// }
	}
}
