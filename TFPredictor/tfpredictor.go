package tfpredictor

import (
	"fmt"

	tg "github.com/galeone/tfgo"
	"github.com/jamOne-/kiwi-zero/predictor"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

type TFPredictor struct {
	model *tg.Model
	path  string
}

func NewTFPredictor(path string) *TFPredictor {
	options := &tf.SessionOptions{Target: "", Config: []byte("2\x02 \x01")} // gpu_options.allow_growth = True
	model := tg.LoadModel(path, []string{"serve"}, options)

	return &TFPredictor{model, path}
}

func (predictor *TFPredictor) GetId() string {
	return predictor.path
}

func (predictor *TFPredictor) PredictValue(features predictor.Features) float32 {
	model := predictor.model
	inputTensor, err := tf.NewTensor([1][][][]float32{features})

	if err != nil {
		fmt.Println(err.Error())
	}

	results := model.Exec([]tf.Output{
		model.Op("StatefulPartitionedCall", 1), // TODO: 0 -- policy, 1 -- value
	}, map[tf.Output]*tf.Tensor{
		model.Op("serving_default_input_1", 0): inputTensor, // TODO
	})

	prediction := results[0].Value().([][]float32)[0][0]
	return prediction
}

func (predictor *TFPredictor) PredictPolicy(features predictor.Features) []float32 {
	model := predictor.model
	inputTensor, err := tf.NewTensor([1][][][]float32{features})

	if err != nil {
		fmt.Println(err.Error())
	}

	policy := model.Exec([]tf.Output{
		model.Op("StatefulPartitionedCall", 0), // TODO: 0 -- policy, 1 -- value
	}, map[tf.Output]*tf.Tensor{
		model.Op("serving_default_input_1", 0): inputTensor, // TODO
	})

	prediction := policy[0].Value().([][]float32)[0]
	return prediction
}
