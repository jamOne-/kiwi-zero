package player

import "github.com/jamOne-/kiwi-zero/game"

type Player interface {
	SelectMove(game *game.Game) game.Move
}
