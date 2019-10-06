package utils

import "gonum.org/v1/gonum/mat"

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
