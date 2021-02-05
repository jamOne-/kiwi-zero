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

type HistoryTuple struct {
	Game game.Game
	Move game.Move
}

type GameResult struct {
	History []*HistoryTuple
	Winner  game.PlayerColor
}

type NewGameFactory func() game.Game
type NewPlayerFactory func() player.Player

func FactorizePlayer(p player.Player) NewPlayerFactory {
	return func() player.Player {
		return p
	}
}

func PlayGame(
	g game.Game,
	blackPlayer player.Player,
	whitePlayer player.Player,
	saveHistory bool,
) *GameResult {
	finished, winner := false, game.PlayerColor(0)
	history := make([]*HistoryTuple, 0)

	for !finished {
		currentPlayer := g.GetCurrentPlayerColor()
		var move game.Move

		if currentPlayer == game.BLACK {
			move = blackPlayer.SelectMove(g)
		} else {
			move = whitePlayer.SelectMove(g)
		}

		if saveHistory {
			tuple := &HistoryTuple{g.Copy(), move}
			history = append(history, tuple)
		}

		finished, winner = g.MakeMove(move)
	}

	return &GameResult{history, winner}
}

func PlayGameAsyncWrapper(
	g game.Game,
	blackPlayer player.Player,
	whitePlayer player.Player,
	saveHistory bool,
	resultChan chan *GameResult,
	waitGroup *sync.WaitGroup,
) {
	resultChan <- PlayGame(g, blackPlayer, whitePlayer, saveHistory)

	if waitGroup != nil {
		waitGroup.Done()
	}
}

func PlayNGames(
	newGameFactory NewGameFactory,
	saveHistory bool,
	player1 player.Player,
	player2 player.Player,
	n int,
) ([]*GameResult, int) {
	results := make([]*GameResult, n)
	totalPositions := 0

	for i := 0; i < n; i++ {
		newGame := newGameFactory()
		result := PlayGame(newGame, player1, player2, saveHistory)
		results[i] = result
		totalPositions += len(result.History)
	}

	return results, totalPositions
}

func PlayTwoGamesAsync(
	newGameFactory NewGameFactory,
	saveHistory bool,
	player1Factory NewPlayerFactory,
	player2Factory NewPlayerFactory,
	blackResultsChan chan *GameResult,
	whiteResultsChan chan *GameResult,
) {
	newGame1 := newGameFactory()
	newGame2 := newGame1.Copy()

	go PlayGameAsyncWrapper(newGame1, player1Factory(), player2Factory(), saveHistory, blackResultsChan, nil)
	go PlayGameAsyncWrapper(newGame2, player2Factory(), player1Factory(), saveHistory, whiteResultsChan, nil)
}

func PlayNGamesAsync(
	newGameFactory NewGameFactory,
	saveHistory bool,
	player1Factory NewPlayerFactory,
	player2Factory NewPlayerFactory,
	n int,
	maxGamesAtOnce int,
) ([]*GameResult, int) {
	resultsChan := make(chan *GameResult, n)
	maxGamesAtOnce = int(math.Min(float64(n), float64(maxGamesAtOnce)))

	for i := 0; i < maxGamesAtOnce/2; i++ {
		PlayTwoGamesAsync(newGameFactory, saveHistory, player1Factory, player2Factory, resultsChan, resultsChan)
	}

	gamesRan := maxGamesAtOnce
	results := make([]*GameResult, n)
	totalPositions := 0

	for i := 0; i < n; i++ {
		result := <-resultsChan
		results[i] = result
		totalPositions += len(result.History)

		if gamesRan < n-1 {
			PlayTwoGamesAsync(newGameFactory, saveHistory, player1Factory, player2Factory, resultsChan, resultsChan)
			gamesRan += 2
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
		result := PlayGame(newGame, player1, player2, false)

		if result.Winner == game.BLACK {
			player1Wins += 1
		}
	}

	restOfGames := numberOfGames - halfOfGames
	for i := 0; i < restOfGames; i++ {
		newGame := gameFactory()
		result := PlayGame(newGame, player2, player1, false)

		if result.Winner == game.WHITE {
			player1Wins += 1
		}
	}

	return player1Wins
}

func ComparePlayersAsync(
	gameFactory NewGameFactory,
	player1Factory NewPlayerFactory,
	player2Factory NewPlayerFactory,
	numberOfGames int,
	maxGamesAtOnce int,
) int {
	halfOfGames := numberOfGames / 2
	player1Wins := 0

	results, _ := PlayNGamesAsync(gameFactory, false, player1Factory, player2Factory, halfOfGames, maxGamesAtOnce)

	for _, result := range results {
		if result.Winner == game.BLACK {
			player1Wins += 1
		}
	}

	results, _ = PlayNGamesAsync(gameFactory, false, player2Factory, player1Factory, numberOfGames-halfOfGames, maxGamesAtOnce)

	for _, result := range results {
		if result.Winner == game.WHITE {
			player1Wins += 1
		}
	}

	return player1Wins
}

func ComparePlayersAsyncWrapper(
	gameFactory NewGameFactory,
	player1Factory NewPlayerFactory,
	player2Factory NewPlayerFactory,
	numberOfGames int,
	maxGamesAtOnce int,
	resultChan chan int,
	waitGroup *sync.WaitGroup,
) {
	resultChan <- ComparePlayersAsync(gameFactory, player1Factory, player2Factory, numberOfGames, maxGamesAtOnce)

	if waitGroup != nil {
		waitGroup.Done()
	}
}

func ComparePlayerWithOthersAsync(
	gameFactory NewGameFactory,
	player player.Player,
	players []player.Player,
	numberOfGames int,
	maxGamesAtOnce int,
) (int, int) {
	playerFactory := FactorizePlayer(player)
	asBlackResultsChan := make(chan *GameResult, 2)
	asWhiteResultsChan := make(chan *GameResult, 2)

	gamesPlayed := 0
	playerWins := 0
	opponent_i := 0

	for gamesPlayed < numberOfGames {
		opponentFactory := FactorizePlayer(players[opponent_i])

		PlayTwoGamesAsync(gameFactory, false, playerFactory, opponentFactory, asBlackResultsChan, asWhiteResultsChan)

		asBlackResult := <-asBlackResultsChan
		asWhiteResult := <-asWhiteResultsChan

		if asBlackResult.Winner == game.BLACK {
			playerWins += 1
		}

		if asWhiteResult.Winner == game.WHITE {
			playerWins += 1
		}

		opponent_i = (opponent_i + 1) % len(players)
		gamesPlayed += 2
	}

	return playerWins, gamesPlayed
}
