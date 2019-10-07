package runner

import (
	"github.com/jamOne-/kiwi-zero/game"
	"github.com/jamOne-/kiwi-zero/player"
)

type GameResult struct {
	History []game.Game
	Winner  game.PlayerColor
}

type NewGameFactory func() game.Game

func PlayGame(g game.Game, blackPlayer player.Player, whitePlayer player.Player) *GameResult {
	finished, winner := false, game.PlayerColor(0)
	history := make([]game.Game, 0)
	history = append(history, g.Copy())

	for !finished {
		currentPlayer := g.GetCurrentPlayerColor()
		var move game.Move

		if currentPlayer == game.BLACK {
			move = blackPlayer.SelectMove(g)
		} else {
			move = whitePlayer.SelectMove(g)
		}

		finished, winner = g.MakeMove(move)
		history = append(history, g.Copy())
	}

	return &GameResult{history, winner}
}

func PlayNGames(newGameFactory NewGameFactory, player1 player.Player, player2 player.Player, n int) ([]*GameResult, int) {
	results := make([]*GameResult, n)
	totalPositions := 0

	for i := 0; i < n; i++ {
		newGame := newGameFactory()
		result := PlayGame(newGame, player1, player2)
		results[i] = result
		totalPositions += len(result.History)
	}

	return results, totalPositions
}
