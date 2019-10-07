package sgd

import (
	"fmt"
	"math"
	"math/rand"

	"gonum.org/v1/gonum/mat"

	"github.com/jamOne-/kiwi-zero/utils"
)

type OptimizeFn func(Xs []*mat.VecDense, ys []float64, weights *mat.VecDense) (float64, *mat.VecDense)
type SGDReturn struct {
	BestWeights             *mat.VecDense
	TestSetErrorRate        float64
	BestValidErrorRate      float64
	TotalEpochs             int
	BestWeightsEpoch        int
	TrainErrorsHistory      []float64
	ValidationErrorsHistory []float64
}

var DEFAULT_PARAMETERS = map[string]float64{
	"alpha0":               1e-4,
	"alpha_const":          1e-5,
	"batch_size":           32,
	"momentum":             0.9,
	"epochs":               50.0,
	"max_epochs":           20000.0,
	"patience_expansion":   1.5,
	"validation_set_ratio": 0.2,
	"test_set_ratio":       0.2,
	"weights_decay":        0.0,
	"debug":                1}

func SGD(f OptimizeFn, weights *mat.VecDense, Xs []*mat.VecDense, ys []float64, parameters map[string]float64) *SGDReturn {
	parameters = utils.MergeMaps(DEFAULT_PARAMETERS, parameters)
	alpha0, alphaConst := parameters["alpha0"], parameters["alpha_const"]
	batchSize := int(parameters["batch_size"])
	momentum := parameters["momentum"]
	numberOfEpochs := int(parameters["epochs"])
	maxEpochs := int(parameters["max_epochs"])
	patienceExpansion := parameters["patience_expansion"]
	// weightsDecay := parameters["weightsDecay"]

	i := 0
	epoch := 0
	velocities := mat.NewVecDense(weights.Len(), nil)

	bestValidErrorRate := math.MaxFloat64
	bestWeights := mat.NewVecDense(weights.Len(), nil)
	bestWeights.CloneVec(weights)
	bestWeightsEpoch := 0

	trainErrors := make([]float64, 0)
	// trainLoss := make([]float64, 0)
	validationErrors := make([]float64, 0)

	rand.Shuffle(len(Xs), func(i int, j int) {
		Xs[i], Xs[j] = Xs[j], Xs[i]
		ys[i], ys[j] = ys[j], ys[i]
	})

	testSetSize := int(math.Floor(parameters["test_set_ratio"] * float64(len(Xs))))
	testX, testy := Xs[:testSetSize], ys[:testSetSize]
	restX, resty := Xs[testSetSize:], ys[testSetSize:]
	validationSetSize := int(math.Floor(parameters["validation_set_ratio"] * float64(len(restX))))
	validationX, validationy := restX[:validationSetSize], resty[:validationSetSize]
	trainX, trainy := restX[validationSetSize:], resty[validationSetSize:]
	debugMode := int(parameters["debug"]) == 1

	numberOfBatches := int(math.Ceil(float64(len(trainX)) / float64(batchSize)))

	for epoch < numberOfEpochs {
		epoch += 1

		for batchIndex := 0; batchIndex < numberOfBatches; batchIndex++ {
			i += 1

			batchStart := batchIndex * batchSize
			batchEnd := int(math.Min(float64(batchStart+batchSize), float64(len(trainX))))
			batchX, batchy := trainX[batchStart:batchEnd], trainy[batchStart:batchEnd]

			errorRate, gradient := f(batchX, batchy, weights)
			trainErrors = append(trainErrors, errorRate)

			alpha := alpha0 / (1.0 + alphaConst*float64(i))
			velocities.ScaleVec(momentum, velocities)
			velocities.AddScaledVec(velocities, alpha, gradient)

			weights.SubVec(weights, velocities)
		}

		validationError, _ := f(validationX, validationy, weights)
		validationErrors = append(validationErrors, validationError)

		if validationError < bestValidErrorRate {
			numberOfEpochs = int(math.Max(float64(numberOfEpochs), float64(epoch)*patienceExpansion+1.0))
			numberOfEpochs = int(math.Min(float64(maxEpochs), float64(numberOfEpochs)))
			bestValidErrorRate = validationError
			bestWeights.CloneVec(weights)
			bestWeightsEpoch = epoch
		}

		if debugMode {
			fmt.Printf("After epoch %d: validationError: %f currently going to do %d epochs\n", epoch, validationError, numberOfEpochs)
		}
	}

	testErrorRate, _ := f(testX, testy, bestWeights)

	if debugMode {
		fmt.Printf("SGD ended after %d epochs having %f error on test set\n", numberOfEpochs, testErrorRate)
	}

	return &SGDReturn{
		bestWeights,
		testErrorRate,
		bestValidErrorRate,
		numberOfEpochs,
		bestWeightsEpoch,
		trainErrors,
		validationErrors}
}
