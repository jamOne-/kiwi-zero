package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/jamOne-/kiwi-zero/game"
	"github.com/jamOne-/kiwi-zero/randomPlayer"
	"github.com/jamOne-/kiwi-zero/reversi"
	"github.com/jamOne-/kiwi-zero/runner"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	RANDOM_MOVES := 4
	RUNS := 1000000

	gameFactory := getRandomStartReversiGameFactory(RANDOM_MOVES)
	positions := make(map[string]bool)
	finishedCounter := 0

	for i := 0; i < RUNS; i++ {
		game := gameFactory()

		finished, _ := game.IsGameFinished()
		if finished {
			finishedCounter += 1
		}

		boardString := game.SerializeBoard(false)
		positions[boardString] = true
	}

	fmt.Printf("Found %d unique positions and %d were already finished\n", len(positions), finishedCounter)
}

/*
	0 -> 1
	1 -> 12
	2 -> 236
	3 -> 7092
	4 -> >230k, 0
*/
func getRandomStartReversiGameFactory(moves int) runner.NewGameFactory {
	return func() game.Game {
		g := reversi.NewReversiGame()

		for i := 0; i < moves; i += 1 {
			g.MakeMove(randomPlayer.SelectRandomMoveDifferentThan(g, reversi.PASS_MOVE))
			g.MakeMove(randomPlayer.SelectRandomMoveDifferentThan(g, reversi.PASS_MOVE))
		}

		return g
	}
}
