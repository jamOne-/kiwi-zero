package utils

import (
	"math"
	"strconv"
	"strings"
	"time"

	"gonum.org/v1/gonum/mat"
)

func SumFloats32(xs []float32) float32 {
	var sum float32 = 0.0

	for _, x := range xs {
		sum += x
	}

	return sum
}

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

// func SaveGameResultsToFile(gameResults []*runner.GameResult, fileName string) {
// 	file, _ := os.Create(fileName)
// 	defer file.Close()

// 	for _, result := range gameResults {
// 		historyLength := len(result.History)
// 		winner := result.Winner

// 		fmt.Fprintf(file, "%d %d\n", winner, historyLength)

// 		for _, state := range result.History {
// 			currentPlayer := state.GetCurrentPlayerColor()
// 			board := state.SerializeBoard(false)

// 			fmt.Fprintf(file, "%d %s\n", currentPlayer, board)
// 		}
// 	}
// }

func TimeNowString() string {
	ret := time.Now().String()[:19]
	ret = strings.ReplaceAll(ret, ":", "")
	ret = strings.ReplaceAll(ret, " ", "-")

	return ret
}

func CreateFilledVector(length int, value float64) *mat.VecDense {
	vec := mat.NewVecDense(length, nil)

	for i := 0; i < length; i++ {
		vec.SetVec(i, value)
	}

	return vec
}

func RotateSquareVector(vec []int8, rotates int) {
	N := int(math.Sqrt(float64(len(vec))))

	for y := 0; y < N/2; y++ {
		for x := y; x < N-y-1; x++ {
			indices := []int{
				y*N + x,
				(N-1-x)*N + y,
				(N-1-y)*N + (N - 1 - x),
				x*N + (N - 1 - y)}

			for rotation := 0; rotation < rotates; rotation++ {
				aux := vec[indices[0]]
				vec[indices[0]] = vec[indices[1]]
				vec[indices[1]] = vec[indices[2]]
				vec[indices[2]] = vec[indices[3]]
				vec[indices[3]] = aux
			}
		}
	}
}

func PerformSymmetryVector1(vec []int8) {
	N := int(math.Sqrt(float64(len(vec))))

	for y := 0; y < N/2; y++ {
		for x := 0; x < N; x++ {
			i1, i2 := y*N+x, (N-1-y)*N+x
			vec[i1], vec[i2] = vec[i2], vec[i1]
		}
	}
}

func PerformSymmetryVector2(vec []int8) {
	N := int(math.Sqrt(float64(len(vec))))

	for x := 0; x < N/2; x++ {
		for y := 0; y < N; y++ {
			i1, i2 := y*N+x, y*N+(N-1-x)
			vec[i1], vec[i2] = vec[i2], vec[i1]
		}
	}
}

func PerformSymmetryVector3(vec []int8) {
	N := int(math.Sqrt(float64(len(vec))))

	for y := 0; y < N; y++ {
		for x := y + 1; x < N; x++ {
			i1, i2 := y*N+x, x*N+y
			vec[i1], vec[i2] = vec[i2], vec[i1]
		}
	}
}

func PerformSymmetryVector4(vec []int8) {
	N := int(math.Sqrt(float64(len(vec))))

	for y := 0; y < N; y++ {
		for x := 0; x < N-y-1; x++ {
			i1, i2 := y*N+x, (N-1-x)*N+(N-1-y)
			vec[i1], vec[i2] = vec[i2], vec[i1]
		}
	}
}

func FloatsToString(xs []float32) string {
	xs_string := make([]string, len(xs))

	for i, x := range xs {
		xs_string[i] = strconv.FormatFloat(float64(x), 'f', -1, 64)
	}

	return strings.Join(xs_string, " ")
}
