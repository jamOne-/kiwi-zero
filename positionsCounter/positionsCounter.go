package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/jamOne-/kiwi-zero/games"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	RANDOM_MOVES := 2
	RUNS := 20000000
	gameFactory := games.GetRandomStartGomokuGameFactory(RANDOM_MOVES)

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
