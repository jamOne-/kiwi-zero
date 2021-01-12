package monteCarloTreeSearchPlayer

import (
	"math/rand"
	"sync"
	"time"

	"github.com/jamOne-/kiwi-zero/game"
)

type ThreadedMonteCarloTreeSearchPlayer struct {
	maxSimulations   int
	maxParallelGames int
}

func NewThreadedMonteCarloTreeSearchPlayer(maxSimulations int, maxParallelGames int) *ThreadedMonteCarloTreeSearchPlayer {
	rand.Seed(time.Now().UnixNano())
	return &ThreadedMonteCarloTreeSearchPlayer{maxSimulations, maxParallelGames}
}

type MoveAndRoot struct {
	Move game.Move
	Root *Node
}

func (player *ThreadedMonteCarloTreeSearchPlayer) SelectMove(game game.Game) game.Move {
	return player.SelectMoveWithRoot(game).Move
}

func (player *ThreadedMonteCarloTreeSearchPlayer) SelectMoveWithRoot(game game.Game) *MoveAndRoot {
	tree := newNode(game, nil)

	var waitgroup sync.WaitGroup
	resultChannel := make(chan *resultTuple)
	gamesChannel := make(chan *gameRequestTuple, player.maxParallelGames)
	go vsUpdater(resultChannel, &waitgroup)
	go gamesRequester(gamesChannel, resultChannel)

	for simulation := 0; simulation < player.maxSimulations; simulation += 1 {
		gameCopy := game.Copy()

		selectedNode, _ := selectNode(gameCopy, tree, 0)
		createdNode := selectedNode.expand(gameCopy)
		updateNs(createdNode)

		waitgroup.Add(1)
		gamesChannel <- &gameRequestTuple{createdNode, gameCopy}
	}

	waitgroup.Wait()
	close(resultChannel)
	close(gamesChannel)

	bestVisitCount, bestNodes := 0, make([]int, 0, 10)
	for i, child := range tree.Nodes {
		if child.N > bestVisitCount {
			bestVisitCount = child.N
			bestNodes = append(make([]int, 0, 10), i)
		} else if child.N == bestVisitCount {
			bestNodes = append(bestNodes, i)
		}
	}

	selectedIndex := bestNodes[rand.Intn(len(bestNodes))]
	return &MoveAndRoot{tree.moves[selectedIndex], tree}
}

type resultTuple struct {
	node   *Node
	result int8
}

func vsUpdater(ch chan *resultTuple, waitgroup *sync.WaitGroup) {
	for tuple := range ch {
		updateVs(tuple.result, tuple.node)
		waitgroup.Done()
	}
}

type gameRequestTuple struct {
	node *Node
	game game.Game
}

func gamesRequester(gameRequestsChannel chan *gameRequestTuple, resultChannel chan *resultTuple) {
	for tuple := range gameRequestsChannel {
		go gamesPlayer(resultChannel, tuple.node, tuple.game)
	}
}

func gamesPlayer(resultChannel chan *resultTuple, node *Node, game game.Game) {
	result := randomSampleFromState(game)
	resultChannel <- &resultTuple{node, result}
}
