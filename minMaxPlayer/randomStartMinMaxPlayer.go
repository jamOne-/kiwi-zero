package minMaxPlayer

import (
	"math/rand"

	"github.com/jamOne-/kiwi-zero/game"
)

type RandomMinMaxPlayer struct {
	depth            int
	randomFirstMoves int
	valueFn          game.ValueFn
}

func NewRandomMinMaxPlayer(depth int, randomFirstMoves int, valueFn game.ValueFn) *RandomMinMaxPlayer {
	return &RandomMinMaxPlayer{depth, randomFirstMoves, valueFn}
}

func (player *RandomMinMaxPlayer) SelectMove(g game.Game) game.Move {
	var move game.Move = -1

	// Compare number of moves that this player has already done
	if g.GetTurnNumber()/2 < player.randomFirstMoves {
		possibleMoves := g.GetPossibleMoves()

		// Random move different than pass
		// TODO: assuming -1 is pass move
		for move == -1 {
			move = possibleMoves[rand.Intn(len(possibleMoves))]
		}

	} else {
		_, move = negaMax(player.valueFn, g, player.depth, -INFINITY, INFINITY)
	}

	return move
}
