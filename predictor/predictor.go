package predictor

type Features = [][][]float32
type Distribution = []float32

type Predictor interface {
	PredictValue(features Features) float32
	PredictPolicy(features Features) Distribution
}
