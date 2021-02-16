package games

import (
	"github.com/jamOne-/kiwi-zero/connectFour"
	"github.com/jamOne-/kiwi-zero/game"
	"github.com/jamOne-/kiwi-zero/gomoku"
	"github.com/jamOne-/kiwi-zero/randomPlayer"
	"github.com/jamOne-/kiwi-zero/reversi"
	"github.com/jamOne-/kiwi-zero/runner"
)

var GAME_FACTORY_DICT = map[string]runner.NewGameFactory{
	"reversi":        ReversiGameFactory,
	"reversirandom":  GetRandomStartReversiGameFactory(4),
	"connect4":       ConnectFourGameFactory,
	"connect4random": GetRandomStartConnectFourGameFactory(3),
	"gomoku":         GomokuGameFactory,
	"gomokurandom":   GetRandomStartGomokuGameFactory(3),
}

func ReversiGameFactory() game.Game {
	return reversi.NewReversiGame()
}

func GetRandomStartReversiGameFactory(randomMoves int) runner.NewGameFactory {
	return func() game.Game {
		g := reversi.NewReversiGame()

		for i := 0; i < randomMoves; i += 1 {
			g.MakeMove(randomPlayer.SelectRandomMoveDifferentThan(g, reversi.PASS_MOVE))
			g.MakeMove(randomPlayer.SelectRandomMoveDifferentThan(g, reversi.PASS_MOVE))
		}

		return g
	}
}

func ConnectFourGameFactory() game.Game {
	return connectFour.NewConnectFourGame()
}

func GetRandomStartConnectFourGameFactory(randomMoves int) runner.NewGameFactory {
	return func() game.Game {
		var g game.Game
		finished := true

		for finished {
			g = connectFour.NewConnectFourGame()

			for i := 0; i < randomMoves; i += 1 {
				g.MakeMove(randomPlayer.SelectRandomMove(g))
				g.MakeMove(randomPlayer.SelectRandomMove(g))
			}

			finished, _ = g.IsGameFinished()
		}

		return g
	}
}

func GomokuGameFactory() game.Game {
	return gomoku.NewGomokuGame()
}

func GetRandomStartGomokuGameFactory(randomMoves int) runner.NewGameFactory {
	return func() game.Game {
		var g game.Game
		finished := true

		for finished {
			g = gomoku.NewGomokuGame()

			for i := 0; i < randomMoves; i += 1 {
				g.MakeMove(randomPlayer.SelectRandomMove(g))
				g.MakeMove(randomPlayer.SelectRandomMove(g))
			}

			finished, _ = g.IsGameFinished()
		}

		return g
	}
}
