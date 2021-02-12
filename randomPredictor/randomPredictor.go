package randomPredictor

import (
	"github.com/jamOne-/kiwi-zero/predictor"
)

type RandomPredictor struct{}

func NewRandomPredictor() *RandomPredictor {
	return &RandomPredictor{}
}

func (predictor *RandomPredictor) GetId() string {
	return "Random predictor =)"
}

func (predictor *RandomPredictor) PredictValue(features predictor.Features) float32 {
	return 0.5
}

func (predictor *RandomPredictor) PredictPolicy(features predictor.Features) predictor.Distribution {
	distribution := make([]float32, 65) // TODO
	// distribution[rand.Intn(65)] = 1

	for i := 0; i < 65; i++ {
		distribution[i] = 1
	}

	return distribution
}

func (predictor *RandomPredictor) PredictValueAndPolicy(features predictor.Features) (float32, []float32) {
	return predictor.PredictValue(features), predictor.PredictPolicy(features)
}
