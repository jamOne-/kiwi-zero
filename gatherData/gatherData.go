package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/jamOne-/kiwi-zero/monteCarloTreeSearchPlayer"

	"github.com/jamOne-/kiwi-zero/reversi"
)

var NUMBER_OF_GAMES = 1000

func main() {
	rand.Seed(time.Now().UnixNano())
	games := make([]string, NUMBER_OF_GAMES)
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
		player := monteCarloTreeSearchPlayer.NewThreadedMonteCarloTreeSearchPlayer(1000, 4)
		finished, winner := game.IsGameFinished()

		for !finished && rand.Float64() > 1.0/50.0 {
			move := player.SelectMove(game)
			finished, winner = game.MakeMove(move)
		}

		if finished {
			i--
			continue
		}

		gamePosition := game.SerializeBoard()

		for !finished {
			move := player.SelectMove(game)
			finished, winner = game.MakeMove(move)
		}

		games[i] = gamePosition + " " + strconv.Itoa(int(winner))

		timeDuration := time.Since(timeStart)
		averageTime += (int(timeDuration) - averageTime) / (i + 1)
	}
}
