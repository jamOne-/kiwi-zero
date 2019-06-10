package minMaxPlayer

import (
	"github.com/jamOne-/kiwi-zero/game"
	"github.com/jamOne-/kiwi-zero/predictor"
	"github.com/jamOne-/kiwi-zero/reversi"
)

type PredictorMinMaxPlayer struct {
	depth     int
	predictor *predictor.Predictor
}

func NewPredictorMinMaxPlayer(depth int) *PredictorMinMaxPlayer {
	predictor := predictor.NewPredictor(predictor.MODEL_PATH)

	return &PredictorMinMaxPlayer{depth, predictor}
}

func (player *PredictorMinMaxPlayer) SelectMove(game *reversi.Game) game.Move {
	_, move := negaMax(player.predictorHeuristicFn, game, player.depth, -INFINITY, INFINITY)
	return move
}

func (player *PredictorMinMaxPlayer) predictorHeuristicFn(game *reversi.Game) int {
	value := player.predictor.PredictBlackWinProb(game) - 0.5
	value *= 10000 * float32(game.Turn)
	return int(value)
}
