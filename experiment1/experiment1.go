package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/jamOne-/kiwi-zero/evaluator"
	"github.com/jamOne-/kiwi-zero/monteCarloTreeSearchPlayer"
	"github.com/jamOne-/kiwi-zero/sgd"

	"gonum.org/v1/gonum/mat"

	"github.com/jamOne-/kiwi-zero/game"
	"github.com/jamOne-/kiwi-zero/minMaxPlayer"
	"github.com/jamOne-/kiwi-zero/reversi"
	"github.com/jamOne-/kiwi-zero/runner"
)

var EVALUATOR_GAMES = 10
var GAMES_PER_ITERATION = 20
var ITERATIONS = 100
var MAX_HISTORY_LENGTH = 128
var MINMAX_DEPTH = 4
var TRAINING_SIZE = 256
var COMPARE_WITH_OLD_MINMAX = true

func main() {
	rand.Seed(time.Now().UnixNano())

	// newWeights := mat.NewVecDense(len(AFTER_1000_ITERS_WEIGHTS), AFTER_1000_ITERS_WEIGHTS)
	// previousWeights := mat.NewVecDense(len(PREVIOUS_WEIGHTS), PREVIOUS_WEIGHTS)

	// newVFn := createWeightedReversiFn(newWeights)
	// newMMPlayer := minMaxPlayer.NewMinMaxPlayer(MINMAX_DEPTH, newVFn)

	// previousVFn := createWeightedReversiFn(previousWeights)
	// previousMMPlayer := minMaxPlayer.NewMinMaxPlayer(MINMAX_DEPTH, previousVFn)

	// // randomPlayer := randomPlayer.NewRandomPlayer()

	// newWins := evaluator.ComparePlayers(reversiGameFactory, newMMPlayer, previousMMPlayer, 50)
	// fmt.Println(newWins)

	// bestWeights := getInitialWeights()
	bestWeights := mat.NewVecDense(66, PREVIOUS_WEIGHTS)
	valueFn := createWeightedReversiFn(bestWeights)
	bestPlayer := minMaxPlayer.NewMinMaxPlayer(MINMAX_DEPTH, valueFn)

	historyPositions := make([]*mat.VecDense, 0)
	historyYs := make([]float64, 0)

	mctsPlayer := monteCarloTreeSearchPlayer.NewMonteCarloTreeSearchPlayer(1000)
	// minMaxWins := evaluator.ComparePlayers(reversiGameFactory, bestPlayer, mctsPlayer, 2*EVALUATOR_GAMES)
	// fmt.Printf("MinMax won %d / %d games versus MCTS\n", minMaxWins, 2*EVALUATOR_GAMES)
	minMaxWins := 0

	oldWeights := mat.NewVecDense(len(PREVIOUS_WEIGHTS), PREVIOUS_WEIGHTS)
	oldValueFn := createWeightedReversiFn(oldWeights)
	oldMinMaxPlayer := minMaxPlayer.NewMinMaxPlayer(MINMAX_DEPTH, oldValueFn)

	for iteration := 0; iteration < ITERATIONS; iteration++ {
		fmt.Printf("Iteration %d / %d...\n", iteration+1, ITERATIONS)

		if COMPARE_WITH_OLD_MINMAX && (iteration+1)%10 == 0 {
			numberOfGames := 20
			bestWins := evaluator.ComparePlayers(reversiGameFactory, bestPlayer, oldMinMaxPlayer, numberOfGames)
			fmt.Printf("Current best player won %d / %d games versus OLD MinMax\n", bestWins, numberOfGames)
		}

		results, _ := runner.PlayNGames(reversiGameFactory, bestPlayer, bestPlayer, GAMES_PER_ITERATION)
		// utils.SaveGameResultsToFile(results, strings.Replace(time.Now().String()[:19], ":", "", -1)+".txt")

		Xs, ys := gameResultsToXsAndys2(results, 2)
		historyPositions = append(historyPositions, Xs...)
		historyYs = append(historyYs, ys...)

		if len(historyPositions) > MAX_HISTORY_LENGTH {
			startIndex := len(historyPositions) - MAX_HISTORY_LENGTH
			historyPositions = historyPositions[startIndex:]
			historyYs = historyYs[startIndex:]
		}

		Xs, ys = chooseXsAndys(historyPositions, historyYs, TRAINING_SIZE)

		sgdResult := sgd.SGD(sgd.MeanSquaredErrorWithGradient, bestWeights, Xs, ys, map[string]float64{
			"alpha0":     5e-5,
			"alphaConst": 0,
			"momentum":   0,
			"batch_size": 16,
			"max_epochs": 10000,
			"debug":      1})

		newValueFn := createWeightedReversiFn(sgdResult.BestWeights)
		candidate := minMaxPlayer.NewMinMaxPlayer(MINMAX_DEPTH, newValueFn)
		candidateWins := evaluator.ComparePlayers(reversiGameFactory, candidate, bestPlayer, EVALUATOR_GAMES)

		fmt.Printf("New candidate won %d / %d games\n", candidateWins, EVALUATOR_GAMES)

		if float64(candidateWins)/float64(EVALUATOR_GAMES) >= 0.5 {
			bestPlayer = candidate
			bestWeights = sgdResult.BestWeights
		}

		fmt.Println("")
	}

	fmt.Println(bestWeights.RawVector().Data)

	minMaxWins = evaluator.ComparePlayers(reversiGameFactory, bestPlayer, mctsPlayer, 5*EVALUATOR_GAMES)
	fmt.Printf("MinMax won %d / %d games versus MCTS\n", minMaxWins, 5*EVALUATOR_GAMES)

	SaveWeightsToFile(bestWeights, "weights_"+strings.Replace(time.Now().String()[:19], ":", "", -1)+".txt")
}

// just a wrap
func reversiGameFactory() game.Game {
	return reversi.NewReversiGame()
}

func gameResultsToXsAndys(results []*runner.GameResult, totalPositions int) ([]*mat.VecDense, []float64) {
	Xs := make([]*mat.VecDense, totalPositions)
	ys := make([]float64, totalPositions)

	positionIndex := 0
	for _, gameResult := range results {
		for _, position := range gameResult.History {
			reversiGame := position.(*reversi.ReversiGame)
			Xs[positionIndex] = ReversiToFeatures(reversiGame)
			ys[positionIndex] = float64(gameResult.Winner) * 100.0
			positionIndex += 1
		}
	}

	return Xs, ys
}

func gameResultsToXsAndys2(results []*runner.GameResult, positionsPerGame int) ([]*mat.VecDense, []float64) {
	totalPositions := len(results) * positionsPerGame
	Xs := make([]*mat.VecDense, totalPositions)
	ys := make([]float64, totalPositions)

	for resultIndex, gameResult := range results {
		for i := 0; i < positionsPerGame; i++ {
			positionIndex := resultIndex*positionsPerGame + i
			position := gameResult.History[rand.Intn(len(gameResult.History))]
			reversiGame := position.(*reversi.ReversiGame)
			Xs[positionIndex] = ReversiToFeatures(reversiGame)
			ys[positionIndex] = float64(gameResult.Winner) * 100.0
		}
	}

	return Xs, ys
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
