package tfpredictor

import (
	"fmt"

	tg "github.com/galeone/tfgo"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

type TFPredictor struct {
	model *tg.Model
}

func NewTFPredictor(path string) *TFPredictor {
	model := tg.LoadModel(path, []string{"eval"}, nil)

	return &TFPredictor{model}
}

func (predictor *TFPredictor) Predict(features preidctor.Features) float32 {
	model := predictor.model
	inputTensor, err := tf.NewTensor([1][][][]float32{features})

	if err != nil {
		fmt.Println(err.Error())
	}

	results := model.Exec([]tf.Output{
		model.Op("value_out", 0),
	}, map[tf.Output]*tf.Tensor{
		model.Op("conv2d_input", 0): inputTensor,
	})

	prediction := results[0].Value().([][]float32)[0][0]
	return prediction
}
