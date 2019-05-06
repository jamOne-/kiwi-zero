package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/jamOne-/kiwi-zero/game"
	"github.com/jamOne-/kiwi-zero/monteCarloTreeSearchPlayer"
	"github.com/jamOne-/kiwi-zero/randomPlayer"

	"github.com/jamOne-/kiwi-zero/reversi"
)

var NUMBER_OF_POSITIONS = 1000
var NUMBER_OF_SIMULATIONS = 100

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
		player := monteCarloTreeSearchPlayer.NewThreadedMonteCarloTreeSearchPlayer(1000, 4)
		finished, _ := game.IsGameFinished()

		for !finished && rand.Float64() > 1.0/35.0 {
			move := player.SelectMove(game)
			finished, _ = game.MakeMove(move)
		}

		if finished {
			i--
			continue
		}

		counter := 0
		var waitgroup sync.WaitGroup
		resultChannel := make(chan int8)
		go counterUpdater(resultChannel, &waitgroup, &counter)

		game.DrawBoard()

		for simulation := 0; simulation < NUMBER_OF_SIMULATIONS; simulation++ {
			gameCopy := game.Copy()
			waitgroup.Add(1)
			go runSimulation(resultChannel, gameCopy)
		}

		waitgroup.Wait()
		close(resultChannel)

		likelyWinner := 0
		if counter > 0 {
			likelyWinner = 1
		} else {
			likelyWinner = -1
		}

		fmt.Println(likelyWinner)
		file.WriteString(game.SerializeBoard() + " " + string(likelyWinner))

		timeDuration := time.Since(timeStart)
		averageTime += (int(timeDuration) - averageTime) / (i + 1)
	}

	file.Close()
}

func counterUpdater(resultChannel chan int8, waitgroup *sync.WaitGroup, counter *int) {
	for result := range resultChannel {
		*counter = *counter + int(result)
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
