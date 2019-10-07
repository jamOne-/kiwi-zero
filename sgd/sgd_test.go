package sgd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/mat"
)

func TestMeanSquaredError(t *testing.T) {
	Xs := []*mat.VecDense{
		mat.NewVecDense(3, []float64{1, 2, 3}),
		mat.NewVecDense(3, []float64{4, 5, 6})}
	weights := mat.NewVecDense(3, []float64{3, 2, 1})
	ys := []float64{5, 10}

	firstResult := float64(1*3 + 2*2 + 3*1)
	secondResult := float64(4*3 + 5*2 + 6*1)
	firstError := ys[0] - firstResult
	secondError := ys[1] - secondResult
	mse := (firstError*firstError + secondError*secondError) / 2.0 / 2.0

	g0 := -0.5 * (firstError*1 + secondError*4)
	g1 := -0.5 * (firstError*2 + secondError*5)
	g2 := -0.5 * (firstError*3 + secondError*6)
	gradient := []float64{g0, g1, g2}

	mseResult, mseGradient := MeanSquaredErrorWithGradient(Xs, ys, weights)
	assert.Equal(t, mse, mseResult)
	assert.Equal(t, gradient, mseGradient.RawVector().Data)
}

func TestSGDSimple(t *testing.T) {
	Xs := []*mat.VecDense{
		mat.NewVecDense(2, []float64{1, 10}),
		mat.NewVecDense(2, []float64{1, 0}),
		mat.NewVecDense(2, []float64{1, 100}),
		mat.NewVecDense(2, []float64{1, 2}),
		mat.NewVecDense(2, []float64{1, 4}),
		mat.NewVecDense(2, []float64{1, -8}),
		mat.NewVecDense(2, []float64{1, 0}),
		mat.NewVecDense(2, []float64{1, 0})}

	ys := []float64{57, 7, 507, 17, 27, -47, 7, 7}
	initialWeights := []float64{0, 0}

	sgdReturn := SGD(MeanSquaredErrorWithGradient, mat.NewVecDense(2, initialWeights), Xs, ys, map[string]float64{
		"batch_size": 2,
		"momentum":   0.1,
		"alpha0":     1e-4,
		"epochs":     20,
		"max_epochs": 1000,
		"debug":      0})

	weights := sgdReturn.BestWeights.RawVector().Data
	// assert.Equal(t, []float64{7, 5}, sgdReturn.bestWeights.RawVector().Data)
	assert.Greater(t, weights[1], 4.8)
	assert.Less(t, weights[1], 5.2)
}
