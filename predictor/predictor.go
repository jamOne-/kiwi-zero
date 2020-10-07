package predictor

type Features = [][][]float32

type Predictor interface {
	Predict(features Features) float32
}
