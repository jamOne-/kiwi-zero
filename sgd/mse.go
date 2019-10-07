package sgd

import (
	"gonum.org/v1/gonum/mat"
)

func MeanSquaredErrorWithGradient(Xs []*mat.VecDense, ys []float64, weights *mat.VecDense) (float64, *mat.VecDense) {
	rows, cols := len(Xs), Xs[0].Len()
	X := mat.NewDense(rows, cols, nil)

	for row := 0; row < rows; row++ {
		X.SetRow(row, Xs[row].RawVector().Data)
	}

	Y := mat.NewVecDense(rows, ys)

	aux := mat.NewVecDense(rows, nil)
	aux.MulVec(X, weights)
	aux.SubVec(Y, aux) // aux contains prediction error now

	gradient := mat.NewVecDense(cols, nil)
	gradient.MulVec(X.T(), aux)
	gradient.ScaleVec(-1.0/float64(rows), gradient)

	aux.MulElemVec(aux, aux)

	mse := mat.Sum(aux) / 2.0 / float64(rows)
	return mse, gradient
}
