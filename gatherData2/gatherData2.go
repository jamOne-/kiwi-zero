package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	mcts "github.com/jamOne-/kiwi-zero/monteCarloTreeSearchPlayer"

	"github.com/jamOne-/kiwi-zero/reversi"
)

type config struct {
	games           int
	mctsSimulations int
	mctsThreads     int
}

func get_config() *config {
	var gamesFlag = flag.Int("games", 10000, "number of games to be simulated")
	var mctsSimulationsFlag = flag.Int("mctsSimulations", 1000, "MCTS Simulations")
	var mctsThreadsFlag = flag.Int("mctsThreads", 4, "MCTS Threads")
	flag.Parse()

	return &config{*gamesFlag, *mctsSimulationsFlag, *mctsThreadsFlag}
}

func main() {
	config := get_config()
	NUMBER_OF_GAMES := config.games
	MCTS_SIMULATIONS := config.mctsSimulations
	MCTS_THREADS := config.mctsThreads

	file, _ := os.Create(strings.Replace(time.Now().String()[:19], ":", "", -1) + ".txt")
	defer file.Close()

	rand.Seed(time.Now().UnixNano())
	averageTime := 0

	for i := 0; i < NUMBER_OF_GAMES; i++ {
		debugStep := NUMBER_OF_GAMES / 1000
		if debugStep == 0 {
			debugStep = 1
		}

		if i%debugStep == 0 {
			fmt.Printf("%v/%v (%.2f%%)\t%v left\n", i+1, NUMBER_OF_GAMES, float32(i+1)*100.0/float32(NUMBER_OF_GAMES), time.Duration(averageTime*(NUMBER_OF_GAMES-i))*time.Nanosecond)
		}

		timeStart := time.Now()
		game := reversi.NewGame()
		player := mcts.NewThreadedMonteCarloTreeSearchPlayer(MCTS_SIMULATIONS, MCTS_THREADS)

		boards := make([]string, 100)
		probs := make([][]float64, 100)
		turns := make([]int8, 100)
		values := make([]float64, 100)

		finished, result := game.IsGameFinished()
		for !finished {
			moveRoot := player.SelectMoveWithRoot(game)
			move := moveRoot.Move
			root := moveRoot.Root

			moveProbs := getMoveProbs(game, root)
			board := game.SerializeBoard(game.Turn == reversi.WHITE)
			turn := len(game.History)

			if turn >= 100 {
				fmt.Println(turn)
			}

			boards[turn] = board
			probs[turn] = moveProbs
			turns[turn] = game.Turn
			values[turn] = float64(root.V) / float64(root.N)

			finished, result = game.MakeMove(move)
		}

		index := rand.Intn(len(game.History))

		fmt.Fprintf(file, "%s ", boards[index])

		for _, prob := range probs[index] {
			fmt.Fprintf(file, "%.5f ", prob)
		}

		fmt.Fprintf(file, "%.5f ", values[index])
		fmt.Fprintf(file, "%d ", turns[index])
		fmt.Fprintf(file, "%d\n", result)

		timeDuration := time.Since(timeStart)
		averageTime += (int(timeDuration) - averageTime) / (i + 1)
	}
}

func getMoveProbs(game *reversi.Game, root *mcts.Node) []float64 {
	probs := make([]float64, 65, 65)
	moves := game.GetPossibleMoves()

	for i := 0; i < len(moves); i++ {
		move := moves[i]
		node := root.Nodes[i]

		if move == -1 {
			move = 64
		}

		probs[move] = float64(node.N) / float64(root.N)
	}

	return probs
}
