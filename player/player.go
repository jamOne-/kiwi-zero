package player

import (
	"github.com/jamOne-/kiwi-zero/game"
	"github.com/jamOne-/kiwi-zero/predictor"
)

type Player interface {
	SelectMove(game game.Game) game.Move
}

type PlayerFactory func(predictor predictor.Predictor) Player
