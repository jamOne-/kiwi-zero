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

func main() {
	rand.Seed(time.Now().UnixNano())
	var simulationsFlag = flag.Int("simulations", 10000, "number of simulations")
	flag.Parse()

	var NUMBER_OF_SIMULATIONS = *simulationsFlag

	averageTime := 0

	file, _ := os.Create(strings.Replace(time.Now().String()[:19], ":", "", -1) + ".txt")

	for i := 0; i < NUMBER_OF_SIMULATIONS; i++ {
		debugStep := NUMBER_OF_SIMULATIONS / 10000
		if debugStep == 0 || true {
			debugStep = 1
		}

		if i%debugStep == 0 {
			fmt.Printf("%v/%v (%.2f%%)\t%v left\n", i+1, NUMBER_OF_SIMULATIONS, float32(i+1)*100.0/float32(NUMBER_OF_SIMULATIONS), time.Duration(averageTime*(NUMBER_OF_SIMULATIONS-i))*time.Nanosecond)
		}

		timeStart := time.Now()
		game := reversi.NewGame()
		player := mcts.NewThreadedMonteCarloTreeSearchPlayer(1000, 4)

		boards := make([]string, 70)
		probs := make([][]float64, 70)
		turns := make([]int8, 70)

		finished, result := game.IsGameFinished()
		for !finished {
			moveRoot := player.SelectMoveWithRoot(game)
			move := moveRoot.Move
			root := moveRoot.Root

			moveProbs := getMoveProbs(game, root)
			board := game.SerializeBoard(game.Turn == reversi.WHITE)
			turn := len(game.History)

			boards[turn] = board
			probs[turn] = moveProbs
			turns[turn] = game.Turn

			finished, result = game.MakeMove(move)
		}

		index := rand.Intn(len(game.History))

		fmt.Fprintf(file, "%s ", boards[index])

		for _, prob := range probs[index] {
			fmt.Fprintf(file, "%.5f ", prob)
		}

		fmt.Fprintf(file, "%d\n", result*turns[index])

		timeDuration := time.Since(timeStart)
		averageTime += (int(timeDuration) - averageTime) / (i + 1)
	}

	file.Close()
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
