package tfpredictor

import (
	"fmt"

	tg "github.com/galeone/tfgo"
	"github.com/jamOne-/kiwi-zero/predictor"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

type TFPredictor struct {
	model *tg.Model
}

func NewTFPredictor(path string) *TFPredictor {
	model := tg.LoadModel(path, []string{"serve"}, nil)

	return &TFPredictor{model}
}

func (predictor *TFPredictor) Predict(features predictor.Features) float32 {
	model := predictor.model
	inputTensor, err := tf.NewTensor([1][][][]float32{features})

	if err != nil {
		fmt.Println(err.Error())
	}

	results := model.Exec([]tf.Output{
		model.Op("StatefulPartitionedCall", 0), // TODO
	}, map[tf.Output]*tf.Tensor{
		model.Op("serving_default_input_1", 0): inputTensor, // TODO
	})

	prediction := results[0].Value().([][]float32)[0][0]
	return prediction
}
