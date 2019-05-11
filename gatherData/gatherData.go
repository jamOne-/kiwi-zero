package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jamOne-/kiwi-zero/game"
	"github.com/jamOne-/kiwi-zero/monteCarloTreeSearchPlayer"
	"github.com/jamOne-/kiwi-zero/randomPlayer"

	"github.com/jamOne-/kiwi-zero/reversi"
)

var NUMBER_OF_POSITIONS = 10000
var NUMBER_OF_SIMULATIONS = 20000
var AVERAGE_POSITION_NUMBER = 35

func main() {
	rand.Seed(time.Now().UnixNano())
	averageTime := 0

	file, _ := os.Create(strings.ReplaceAll(time.Now().String()[:19], ":", "") + ".txt")

	for i := 0; i < NUMBER_OF_POSITIONS; i++ {
		debugStep := NUMBER_OF_POSITIONS / 1000
		if debugStep == 0 {
			debugStep = 1
		}

		if i%debugStep == 0 {
			fmt.Printf("%v/%v (%.2f%%)\t%v left\n", i+1, NUMBER_OF_POSITIONS, float32(i+1)*100.0/float32(NUMBER_OF_POSITIONS), time.Duration(averageTime*(NUMBER_OF_POSITIONS-i))*time.Nanosecond)
		}

		timeStart := time.Now()
		game := reversi.NewGame()
		player := monteCarloTreeSearchPlayer.NewThreadedMonteCarloTreeSearchPlayer(100, 1)
		finished, _ := game.IsGameFinished()

		for !finished && rand.Float64() > 1.0/float64(AVERAGE_POSITION_NUMBER) {
			move := player.SelectMove(game)
			finished, _ = game.MakeMove(move)
		}

		if finished {
			i--
			continue
		}

		blackWins, draws, whiteWins := 0, 0, 0
		var waitgroup sync.WaitGroup
		resultChannel := make(chan int8)
		go countersUpdater(resultChannel, &waitgroup, &blackWins, &draws, &whiteWins)

		// game.DrawBoard()

		for simulation := 0; simulation < NUMBER_OF_SIMULATIONS; simulation++ {
			gameCopy := game.Copy()
			waitgroup.Add(1)
			go runSimulation(resultChannel, gameCopy)
		}

		waitgroup.Wait()
		close(resultChannel)

		file.WriteString(game.SerializeBoard() + " " + strconv.Itoa(blackWins) + " " + strconv.Itoa(draws) + " " + strconv.Itoa(whiteWins) + "\n")

		timeDuration := time.Since(timeStart)
		averageTime += (int(timeDuration) - averageTime) / (i + 1)
	}

	file.Close()
}

func countersUpdater(resultChannel chan int8, waitgroup *sync.WaitGroup, blackWins *int, draws *int, whiteWins *int) {
	for result := range resultChannel {
		if result == 1 {
			*blackWins = *blackWins + 1
		} else if result == 0 {
			*draws = *draws + 1
		} else {
			*whiteWins = *whiteWins + 1
		}

		waitgroup.Done()
	}
}

func runSimulation(resultChannel chan int8, game game.Game) {
	finished, winner := false, int8(0)
	player := randomPlayer.NewRandomPlayer()

	for !finished {
		move := player.SelectMove(game)
		finished, winner = game.MakeMove(move)
	}

	resultChannel <- winner
}
