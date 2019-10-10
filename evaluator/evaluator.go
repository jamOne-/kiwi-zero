package evaluator

import (
	"github.com/jamOne-/kiwi-zero/game"
	"github.com/jamOne-/kiwi-zero/player"
	"github.com/jamOne-/kiwi-zero/runner"
)

func ComparePlayers(gameFactory runner.NewGameFactory, player1 player.Player, player2 player.Player, numberOfGames int) int {
	player1Wins := 0
	halfOfGames := numberOfGames / 2

	results, _ := runner.PlayNGames(gameFactory, player1, player2, halfOfGames)

	for _, result := range results {
		if result.Winner == game.BLACK {
			player1Wins += 1
		}
	}

	restOfGames := numberOfGames - halfOfGames
	results, _ = runner.PlayNGames(gameFactory, player2, player1, restOfGames)

	for _, result := range results {
		if result.Winner == game.WHITE {
			player1Wins += 1
		}
	}

	return player1Wins
}
