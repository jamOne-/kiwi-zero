package runner

import (
	"sync"

	"github.com/jamOne-/kiwi-zero/game"
	"github.com/jamOne-/kiwi-zero/player"
)

type GameResultsBatch struct {
	Results        []*GameResult
	TotalPositions int
}

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

func PlayGameAsyncWrapper(g game.Game, blackPlayer player.Player, whitePlayer player.Player, resultChan chan *GameResult, waitGroup *sync.WaitGroup) {
	resultChan <- PlayGame(g, blackPlayer, whitePlayer)
	waitGroup.Done()
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

func PlayNGamesAsync(newGameFactory NewGameFactory, player1 player.Player, player2 player.Player, n int) ([]*GameResult, int) {
	var waitGroup sync.WaitGroup
	waitGroup.Add(n)

	resultsChan := make(chan *GameResult, n)

	for i := 0; i < n; i++ {
		newGame := newGameFactory()
		go PlayGameAsyncWrapper(newGame, player1, player2, resultsChan, &waitGroup)
	}

	waitGroup.Wait()

	results := make([]*GameResult, n)
	totalPositions := 0

	for i := 0; i < n; i++ {
		result := <-resultsChan
		results[i] = result
		totalPositions += len(result.History)
	}

	return results, totalPositions
}

func ComparePlayers(gameFactory NewGameFactory, player1 player.Player, player2 player.Player, numberOfGames int) int {
	player1Wins := 0
	halfOfGames := numberOfGames / 2

	for i := 0; i < halfOfGames; i++ {
		newGame := gameFactory()
		result := PlayGame(newGame, player1, player2)

		if result.Winner == game.BLACK {
			player1Wins += 1
		}
	}

	restOfGames := numberOfGames - halfOfGames
	for i := 0; i < restOfGames; i++ {
		newGame := gameFactory()
		result := PlayGame(newGame, player2, player1)

		if result.Winner == game.WHITE {
			player1Wins += 1
		}
	}

	return player1Wins
}

func ComparePlayersAsync(gameFactory NewGameFactory, player1 player.Player, player2 player.Player, numberOfGames int) int {
	var waitGroup sync.WaitGroup
	waitGroup.Add(numberOfGames)

	halfOfGames := numberOfGames / 2
	restOfGames := numberOfGames - halfOfGames

	blackResults := make(chan *GameResult, halfOfGames)
	whiteResults := make(chan *GameResult, restOfGames)

	for i := 0; i < numberOfGames; i++ {
		newGame := gameFactory()

		if i < halfOfGames {
			go PlayGameAsyncWrapper(newGame, player1, player2, blackResults, &waitGroup)
		} else {
			go PlayGameAsyncWrapper(newGame, player2, player1, whiteResults, &waitGroup)
		}
	}

	waitGroup.Wait()

	player1Wins := 0

	for i := 0; i < numberOfGames; i++ {
		if i < halfOfGames {
			result := <-blackResults

			if result.Winner == game.BLACK {
				player1Wins += 1
			}
		} else {
			result := <-whiteResults

			if result.Winner == game.WHITE {
				player1Wins += 1
			}
		}
	}

	return player1Wins
}
