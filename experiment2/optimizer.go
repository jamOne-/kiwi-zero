package main

import (
	"fmt"
	"math/rand"

	"github.com/spf13/viper"

	"github.com/jamOne-/kiwi-zero/game"
	"github.com/jamOne-/kiwi-zero/reversi"
	"github.com/jamOne-/kiwi-zero/reversiValueFns"
	"github.com/jamOne-/kiwi-zero/runner"
	"github.com/jamOne-/kiwi-zero/sgd"

	"gonum.org/v1/gonum/mat"
)

func Optimizer(gameResultsChan chan *runner.GameResultsBatch, newWeightsChan chan *mat.VecDense, initialWeights *mat.VecDense) {
	MAX_HISTORY_LENGTH := viper.GetInt("MAX_HISTORY_LENGTH")
	TRAINING_SIZE := viper.GetInt("TRAINING_SIZE")
	SGD_CONFIG := viper.Get("SGD_CONFIG").(map[string]float64)

	bestWeights := initialWeights
	gamePositions := make([]*mat.VecDense, 0)
	gameWinners := make([]float64, 0)

	training_i := 1

	for {
		select {
		case batch := <-gameResultsChan:
			results, totalPositions := batch.Results, batch.TotalPositions
			positions, winners := splitResults(results, totalPositions)
			features := createFeaturesSlice(positions)

			gamePositions = append(gamePositions, features...)
			gameWinners = append(gameWinners, winners...)

			if len(gamePositions) > MAX_HISTORY_LENGTH {
				startIndex := len(gamePositions) - MAX_HISTORY_LENGTH
				gamePositions = gamePositions[startIndex:]
				gameWinners = gameWinners[startIndex:]
			}

			Xs, ys := chooseXsAndys(gamePositions, gameWinners, TRAINING_SIZE)
			sgdResult := sgd.SGD(sgd.MeanSquaredErrorWithGradient, bestWeights, Xs, ys, SGD_CONFIG)
			bestWeights = sgdResult.BestWeights

			fmt.Printf("Optimizer (%d): totalEpochs=%d, bestEpoch=%d, testSetError=%.2f, validationSetError=%.2f, trainingError=%.2f\n", training_i, sgdResult.TotalEpochs, sgdResult.BestWeightsEpoch, sgdResult.TestSetErrorRate, sgdResult.BestValidErrorRate, sgdResult.TrainErrorsHistory[len(sgdResult.TrainErrorsHistory)-1])
			training_i += 1

			select {
			case newWeightsChan <- bestWeights:
				// try to send

			default:
				// else skip
			}
		}
	}
}

func createFeaturesSlice(positions []game.Game) []*mat.VecDense {
	featuresSlice := make([]*mat.VecDense, len(positions))

	for i, position := range positions {
		reversiGame := position.(*reversi.ReversiGame)
		featuresSlice[i] = reversiValueFns.ReversiToFeatures(reversiGame)
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
		winner := float64(result.Winner) * 100.0

		for _, position := range result.History {
			positions[index] = position
			winners[index] = winner
			index += 1
		}
	}

	return positions, winners
}
