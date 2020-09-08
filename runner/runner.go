package runner

import (
	"math"
	"sync"

	"github.com/jamOne-/kiwi-zero/game"
	"github.com/jamOne-/kiwi-zero/player"
)

type GameResultsBatch struct {
	Results        []*GameResult
	TotalPositions int
}

type GameResult struct {
	FeaturesList []game.Features
	Winner       game.PlayerColor
}

type NewGameFactory func() game.Game

func EmptyFeaturesFn(g game.Game) game.Features {
	return nil
}

func PlayGame(
	g game.Game,
	blackPlayer player.Player,
	whitePlayer player.Player,
	featuresFn game.GameToFeaturesFn,
) *GameResult {
	if featuresFn == nil {
		featuresFn = EmptyFeaturesFn
	}

	finished, winner := false, game.PlayerColor(0)
	featuresList := make([]game.Features, 0)
	featuresList = append(featuresList, featuresFn(g))

	for !finished {
		currentPlayer := g.GetCurrentPlayerColor()
		var move game.Move

		if currentPlayer == game.BLACK {
			move = blackPlayer.SelectMove(g)
		} else {
			move = whitePlayer.SelectMove(g)
		}

		finished, winner = g.MakeMove(move)
		featuresList = append(featuresList, featuresFn(g))
	}

	return &GameResult{featuresList, winner}
}

func PlayGameAsyncWrapper(
	g game.Game,
	blackPlayer player.Player,
	whitePlayer player.Player,
	featuresFn game.GameToFeaturesFn,
	resultChan chan *GameResult,
	waitGroup *sync.WaitGroup,
) {
	resultChan <- PlayGame(g, blackPlayer, whitePlayer, featuresFn)

	if waitGroup != nil {
		waitGroup.Done()
	}
}

func PlayNGames(
	newGameFactory NewGameFactory,
	featuresFn game.GameToFeaturesFn,
	player1 player.Player,
	player2 player.Player,
	n int,
) ([]*GameResult, int) {
	results := make([]*GameResult, n)
	totalPositions := 0

	for i := 0; i < n; i++ {
		newGame := newGameFactory()
		result := PlayGame(newGame, player1, player2, featuresFn)
		results[i] = result
		totalPositions += len(result.FeaturesList)
	}

	return results, totalPositions
}

func PlayNGamesAsync(
	newGameFactory NewGameFactory,
	featuresFn game.GameToFeaturesFn,
	player1 player.Player,
	player2 player.Player,
	n int,
	maxGamesAtOnce int,
) ([]*GameResult, int) {
	resultsChan := make(chan *GameResult, n)
	maxGamesAtOnce = int(math.Min(float64(n), float64(maxGamesAtOnce)))

	for i := 0; i < maxGamesAtOnce; i++ {
		newGame := newGameFactory()
		go PlayGameAsyncWrapper(newGame, player1, player2, featuresFn, resultsChan, nil)
	}

	gamesRan := maxGamesAtOnce
	results := make([]*GameResult, n)
	totalPositions := 0

	for i := 0; i < n; i++ {
		result := <-resultsChan
		results[i] = result
		totalPositions += len(result.FeaturesList)

		if gamesRan < n {
			newGame := newGameFactory()
			go PlayGameAsyncWrapper(newGame, player1, player2, featuresFn, resultsChan, nil)
			gamesRan += 1
		}
	}

	return results, totalPositions
}

func ComparePlayers(
	gameFactory NewGameFactory,
	player1 player.Player,
	player2 player.Player,
	numberOfGames int,
) int {
	player1Wins := 0
	halfOfGames := numberOfGames / 2

	for i := 0; i < halfOfGames; i++ {
		newGame := gameFactory()
		result := PlayGame(newGame, player1, player2, nil)

		if result.Winner == game.BLACK {
			player1Wins += 1
		}
	}

	restOfGames := numberOfGames - halfOfGames
	for i := 0; i < restOfGames; i++ {
		newGame := gameFactory()
		result := PlayGame(newGame, player2, player1, nil)

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
			go PlayGameAsyncWrapper(newGame, player1, player2, nil, blackResults, &waitGroup)
		} else {
			go PlayGameAsyncWrapper(newGame, player2, player1, nil, whiteResults, &waitGroup)
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
