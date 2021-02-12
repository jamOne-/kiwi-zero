package runner

import (
	"math"

	"github.com/jamOne-/kiwi-zero/game"
	"github.com/jamOne-/kiwi-zero/player"
	"github.com/jamOne-/kiwi-zero/utils"
)

type GameResultsBatch struct {
	Results        []*GameResult
	TotalPositions int
}

type HistoryTuple struct {
	Game   game.Game
	Policy []float32
}

type GameResult struct {
	History []*HistoryTuple
	Winner  game.PlayerColor
	Label   string
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
	label string,
) *GameResult {
	finished, winner := false, game.PlayerColor(0)
	history := make([]*HistoryTuple, 0)

	for !finished {
		currentColor := g.GetCurrentPlayerColor()
		var currentPlayer player.Player

		if currentColor == game.BLACK {
			currentPlayer = blackPlayer
		} else {
			currentPlayer = whitePlayer
		}

		var move game.Move

		if saveHistory {
			var policy []float32

			if p, ok := currentPlayer.(player.PlayerWithPolicy); ok {
				move, policy = p.SelectMoveWithPolicy(g)
			} else {
				move = currentPlayer.SelectMove(g)
				policy = g.EncodeMoveToPolicy(move)
			}

			tuple := &HistoryTuple{g.Copy(), policy}
			history = append(history, tuple)
		} else {
			move = currentPlayer.SelectMove(g)
		}

		finished, winner = g.MakeMove(move)
	}

	return &GameResult{history, winner, label}
}

func PlayGameAsyncWrapper(
	g game.Game,
	blackPlayer player.Player,
	whitePlayer player.Player,
	saveHistory bool,
	resultChan chan *GameResult,
	label string,
) {
	resultChan <- PlayGame(g, blackPlayer, whitePlayer, saveHistory, label)
}

func PlayNGamesAsync(
	newGameFactory NewGameFactory,
	saveHistory bool,
	player1Factory NewPlayerFactory,
	player2Factory NewPlayerFactory,
	n int,
	maxGamesAtOnce int,
) (int, []*GameResult, int) {
	resultsChan := make(chan *GameResult, n)

	gamesToPlay := make([]game.Game, n/2)
	for i := 0; i < n/2; i++ {
		gamesToPlay[i] = newGameFactory()
	}

	maxGamesAtOnce = int(math.Min(float64(n), float64(maxGamesAtOnce)))

	for i := 0; i < maxGamesAtOnce; i++ {
		game := gamesToPlay[i/2].Copy()

		if i%2 == 0 {
			go PlayGameAsyncWrapper(game, player1Factory(), player2Factory(), saveHistory, resultsChan, "black")
		} else {
			go PlayGameAsyncWrapper(game, player2Factory(), player1Factory(), saveHistory, resultsChan, "white")
		}
	}

	gamesRan := maxGamesAtOnce
	results := make([]*GameResult, n)
	totalPositions := 0
	player1Wins := 0

	for i := 0; i < n; i++ {
		result := <-resultsChan
		results[i] = result
		totalPositions += len(result.History)
		player1Wins += utils.BoolToInt(result.Label == "black" && result.Winner == game.BLACK || result.Label == "white" && result.Winner == game.WHITE)

		if gamesRan < n {
			game := gamesToPlay[gamesRan/2].Copy()

			if gamesRan%2 == 0 {
				go PlayGameAsyncWrapper(game, player1Factory(), player2Factory(), saveHistory, resultsChan, "black")
			} else {
				go PlayGameAsyncWrapper(game, player2Factory(), player1Factory(), saveHistory, resultsChan, "white")
			}

			gamesRan += 1
		}
	}

	return player1Wins, results, totalPositions
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
		result := PlayGame(newGame, player1, player2, false, "")

		if result.Winner == game.BLACK {
			player1Wins += 1
		}
	}

	restOfGames := numberOfGames - halfOfGames
	for i := 0; i < restOfGames; i++ {
		newGame := gameFactory()
		result := PlayGame(newGame, player2, player1, false, "")

		if result.Winner == game.WHITE {
			player1Wins += 1
		}
	}

	return player1Wins
}

func ComparePlayerWithOthersAsync(
	gameFactory NewGameFactory,
	player player.Player,
	players []player.Player,
	numberOfGames int,
	maxGamesAtOnce int,
) (int, int) {
	playerFactory := FactorizePlayer(player)

	gamesPlayed := 0
	playerWins := 0
	opponent_i := 0

	for gamesPlayed < numberOfGames {
		opponentFactory := FactorizePlayer(players[opponent_i])
		wins, _, _ := PlayNGamesAsync(gameFactory, false, playerFactory, opponentFactory, 2, maxGamesAtOnce)
		playerWins += wins

		opponent_i = (opponent_i + 1) % len(players)
		gamesPlayed += 2
	}

	return playerWins, gamesPlayed
}
