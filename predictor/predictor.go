package predictor

import "github.com/jamOne-/kiwi-zero/game"

type Features = [][][]float32
type Distribution = []float32

type Predictor interface {
	GetId() string
	PredictValue(features Features) float32
	PredictPolicy(features Features) Distribution
	PredictValueAndPolicy(features Features) (float32, Distribution)
}

func CreateValueAndPolicyFn(gameToFeaturesFn game.GameToFeaturesFn, predictor Predictor) func(game game.Game) (float64, []float32) {
	return func(game game.Game) (float64, []float32) {
		features := gameToFeaturesFn(game)
		value, policy := predictor.PredictValueAndPolicy(features)

		return float64(value)*2.0 - 1.0, policy
	}
}
