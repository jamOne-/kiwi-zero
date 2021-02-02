package predictor

type Features = [][][]float32
type Distribution = []float32

type Predictor interface {
	GetId() string
	PredictValue(features Features) float32
	PredictPolicy(features Features) Distribution
}
