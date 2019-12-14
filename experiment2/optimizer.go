package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os/exec"
	"strconv"
	"strings"

	"github.com/spf13/viper"

	"github.com/jamOne-/kiwi-zero/game"
	"github.com/jamOne-/kiwi-zero/reversi"
	"github.com/jamOne-/kiwi-zero/reversiValueFns"
	"github.com/jamOne-/kiwi-zero/runner"
	"github.com/jamOne-/kiwi-zero/utils"

	"gonum.org/v1/gonum/mat"
)

type trainingParams struct {
	Xs []*mat.VecDense
	ys []float64
}

func Optimizer(
	gameResultsChan chan *runner.GameResultsBatch,
	newWeightsChan chan *mat.VecDense,
	initialWeights *mat.VecDense,
	reversiToFeaturesFn reversiValueFns.ReversiToFeaturesFn) {

	MAX_HISTORY_LENGTH := viper.GetInt("MAX_HISTORY_LENGTH")
	TRAINING_SIZE := viper.GetInt("TRAINING_SIZE")
	TRAINING_FLIP_POSITIONS_PROB := viper.GetFloat64("TRAINING_FLIP_POSITIONS_PROB")
	TRAINING_TRANSFORM_POSITIONS := viper.GetBool("TRAINING_TRANSFORM_POSITIONS")
	SGD_CONFIG := viper.Get("SGD_CONFIG").(map[string]float64)

	// bestWeights := initialWeights
	gamePositions := make([]*mat.VecDense, 0)
	gameWinners := make([]float64, 0)

	pythonOptimizerCmd := exec.Command(
		"python", "../python/optimizer/optimizer.py",
		"--weights", strconv.Itoa(initialWeights.Len()),
		"--learning_rate", fmt.Sprintf("%f", SGD_CONFIG["alpha0"]),
		"--epochs", strconv.Itoa(int(SGD_CONFIG["max_epochs"])),
		"--batch_size", strconv.Itoa(int(SGD_CONFIG["batch_size"])))

	fmt.Println("Args: " + strings.Join(pythonOptimizerCmd.Args, " "))

	optimizerIn, _ := pythonOptimizerCmd.StdinPipe()
	optimizerOut, _ := pythonOptimizerCmd.StdoutPipe()

	trainingChan := make(chan *trainingParams)
	go trainer(trainingChan, optimizerIn, optimizerOut, newWeightsChan)
	// go optimizerReader(optimizerOut, newWeightsChan)
	// go echoReader(optimizerOut)
	pythonOptimizerCmd.Start()

	// training_i := 1

	for {
		select {
		case batch := <-gameResultsChan:
			results, totalPositions := batch.Results, batch.TotalPositions
			positions, winners := splitResults(results, totalPositions)

			if TRAINING_TRANSFORM_POSITIONS {
				transformPositions(positions)
			}

			if TRAINING_FLIP_POSITIONS_PROB > 0 {
				flipPositionsColors(TRAINING_FLIP_POSITIONS_PROB, positions, winners)
			}

			// transform -1, 0, 1 winners to black winning probability
			transformWinnersToProbabilities(winners)

			features := createFeaturesSlice(reversiToFeaturesFn, positions)

			gamePositions = append(gamePositions, features...)
			gameWinners = append(gameWinners, winners...)

			if len(gamePositions) > MAX_HISTORY_LENGTH {
				startIndex := len(gamePositions) - MAX_HISTORY_LENGTH
				gamePositions = gamePositions[startIndex:]
				gameWinners = gameWinners[startIndex:]
			}

			Xs, ys := chooseXsAndys(gamePositions, gameWinners, TRAINING_SIZE)
			params := &trainingParams{Xs, ys}

			select {
			case trainingChan <- params:
				// try to send

			default:
				// else skip
			}
		}
	}
}

func transformWinnersToProbabilities(winners []float64) {
	for i, winner := range winners {
		blackWinProb := (winner + 1) / 2
		winners[i] = blackWinProb
	}
}

func createFeaturesSlice(reversiToFeaturesFn reversiValueFns.ReversiToFeaturesFn, positions []game.Game) []*mat.VecDense {
	featuresSlice := make([]*mat.VecDense, len(positions))

	for i, position := range positions {
		reversiGame := position.(*reversi.ReversiGame)
		featuresSlice[i] = reversiToFeaturesFn(reversiGame)
	}

	return featuresSlice
}

func chooseXsAndys(XsSource []*mat.VecDense, ysSource []float64, N int) ([]*mat.VecDense, []float64) {
	Xs := make([]*mat.VecDense, N)
	ys := make([]float64, N)

	for i := 0; i < N; i++ {
		index := rand.Intn(len(XsSource))
		Xs[i] = XsSource[index]
		ys[i] = ysSource[index]
	}

	return Xs, ys
}

func splitResults(results []*runner.GameResult, totalPositions int) ([]game.Game, []float64) {
	positions := make([]game.Game, totalPositions)
	winners := make([]float64, totalPositions)

	index := 0
	for _, result := range results {
		winner := float64(result.Winner)

		for _, position := range result.History {
			positions[index] = position
			winners[index] = winner
			index += 1
		}
	}

	return positions, winners
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

func flipPositionsColors(flipProb float64, positions []game.Game, winners []float64) {
	for i, position := range positions {
		if rand.Float64() < flipProb {
			flipGameColors(position)
			winners[i] = 1.0 - winners[i]
			// winners[i] *= -1
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

func parseWeights(weightsString string) *mat.VecDense {
	splitted := strings.Split(weightsString, " ")
	weights := mat.NewVecDense(len(splitted), nil)

	for i, s := range splitted {
		weight, _ := strconv.ParseFloat(s, 64)
		weights.SetVec(i, weight)
	}

	return weights
}

func trainer(paramsChan chan *trainingParams, optimizerIn io.WriteCloser, optimizerOut io.ReadCloser, newWeightsChan chan *mat.VecDense) {
	optimizerOutReader := bufio.NewReader(optimizerOut)
	training_i := 1

	for params := range paramsChan {
		Xs, ys := params.Xs, params.ys

		optimizerIn.Write([]byte(fmt.Sprintf("%v\n", len(ys))))

		for _, X := range Xs {
			line := utils.FloatsToString(X.RawVector().Data) + "\n"
			optimizerIn.Write([]byte(line))
		}

		for _, y := range ys {
			line := fmt.Sprintf("%v\n", y)
			optimizerIn.Write([]byte(line))
		}

		// fmt.Printf("Optimizer (%d): training...\n", training_i)

		trainingSummary, _ := optimizerOutReader.ReadString('\n')
		newWeightsString, _ := optimizerOutReader.ReadString('\n')
		newWeights := parseWeights(newWeightsString)
		fmt.Printf("Optimizer (%d): %s", training_i, trainingSummary)
		training_i += 1

		select {
		case newWeightsChan <- newWeights:
			// try to send

		default:
			// else skip
		}

		// sgdResult := sgd.SGD(sgd.MeanSquaredErrorWithGradient, bestWeights, Xs, ys, SGD_CONFIG)
		// bestWeights = sgdResult.BestWeights

		// fmt.Printf("Optimizer (%d): totalEpochs=%d, bestEpoch=%d, testSetError=%.2f, validationSetError=%.2f, trainingError=%.2f\n", training_i, sgdResult.TotalEpochs, sgdResult.BestWeightsEpoch, sgdResult.TestSetErrorRate, sgdResult.BestValidErrorRate, sgdResult.TrainErrorsHistory[len(sgdResult.TrainErrorsHistory)-1])
	}
}
