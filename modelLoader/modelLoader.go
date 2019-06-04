package main

import (
	"fmt"

	"github.com/jamOne-/kiwi-zero/game"
	"github.com/jamOne-/kiwi-zero/reversi"

	tg "github.com/galeone/tfgo"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

type Predictor struct {
	model *tg.Model
}

var MODEL_PATH = "../python/saved_models/1559589871"

func main() {
	predictor := NewPredictor(MODEL_PATH)
	game := reversi.NewGame()

	prediction := predictor.PredictBlackWinProb(game)
	fmt.Println(prediction)
}

func NewPredictor(path string) *Predictor {
	return &Predictor{tg.LoadModel(path, []string{"eval"}, nil)}
}

func (predictor *Predictor) PredictBlackWinProb(game game.Game) float32 {
	model := predictor.model
	board := game.OneHotBoard()
	inputTensor, err := tf.NewTensor([1][][][]float32{board})

	if err != nil {
		fmt.Println(err.Error())
	}

	results := model.Exec([]tf.Output{
		model.Op("dense_1/Relu", 0),
	}, map[tf.Output]*tf.Tensor{
		model.Op("conv2d_input", 0): inputTensor,
	})

	prediction := results[0].Value().([][]float32)[0][0]
	return prediction
}
