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
	"reversirandom":  RandomStartReversiGameFactory,
	"connect4":       ConnectFourGameFactory,
	"connect4random": RandomStartConnectFourGameFactory,
	"gomoku":         GomokuGameFactory,
	"gomokurandom":   RandomStartGomokuGameFactory,
}

func ReversiGameFactory() game.Game {
	return reversi.NewReversiGame()
}

func RandomStartReversiGameFactory() game.Game {
	g := reversi.NewReversiGame()
	NUMBER_OF_RANDOM_MOVES := 4

	for i := 0; i < NUMBER_OF_RANDOM_MOVES; i += 1 {
		g.MakeMove(randomPlayer.SelectRandomMoveDifferentThan(g, reversi.PASS_MOVE))
		g.MakeMove(randomPlayer.SelectRandomMoveDifferentThan(g, reversi.PASS_MOVE))
	}

	return g
}

func ConnectFourGameFactory() game.Game {
	return connectFour.NewConnectFourGame()
}

func RandomStartConnectFourGameFactory() game.Game {
	NUMBER_OF_RANDOM_MOVES := 3

	var g game.Game
	finished := true

	for finished {
		g = connectFour.NewConnectFourGame()

		for i := 0; i < NUMBER_OF_RANDOM_MOVES; i += 1 {
			g.MakeMove(randomPlayer.SelectRandomMove(g))
			g.MakeMove(randomPlayer.SelectRandomMove(g))
		}

		finished, _ = g.IsGameFinished()
	}

	return g
}

func GomokuGameFactory() game.Game {
	return gomoku.NewGomokuGame()
}

func RandomStartGomokuGameFactory() game.Game {
	NUMBER_OF_RANDOM_MOVES := 3

	var g game.Game
	finished := true

	for finished {
		g = gomoku.NewGomokuGame()

		for i := 0; i < NUMBER_OF_RANDOM_MOVES; i += 1 {
			g.MakeMove(randomPlayer.SelectRandomMove(g))
			g.MakeMove(randomPlayer.SelectRandomMove(g))
		}

		finished, _ = g.IsGameFinished()
	}

	return g
}
