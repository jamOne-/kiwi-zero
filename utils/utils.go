package utils

import (
	"fmt"
	"os"

	"github.com/jamOne-/kiwi-zero/runner"
	"gonum.org/v1/gonum/mat"
)

func Int8SliceToVecDense(xs []int8) *mat.VecDense {
	n := len(xs)
	vec := mat.NewVecDense(n, nil)

	for i := 0; i < n; i++ {
		vec.SetVec(i, float64(xs[i]))
	}

	return vec
}

func MergeMaps(m1 map[string]float64, m2 map[string]float64) map[string]float64 {
	resultMap := make(map[string]float64)

	for k, v := range m1 {
		resultMap[k] = v
	}

	for k, v := range m2 {
		resultMap[k] = v
	}

	return resultMap
}

func SaveGameResultsToFile(gameResults []*runner.GameResult, fileName string) {
	file, _ := os.Create(fileName)
	defer file.Close()

	for _, result := range gameResults {
		historyLength := len(result.History)
		winner := result.Winner

		fmt.Fprintf(file, "%d %d\n", winner, historyLength)

		for _, state := range result.History {
			currentPlayer := state.GetCurrentPlayerColor()
			board := state.SerializeBoard(false)

			fmt.Fprintf(file, "%d %s\n", currentPlayer, board)
		}
	}
}
