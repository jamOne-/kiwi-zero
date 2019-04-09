package main

import (
	"math"
	"math/rand"
	"sync"
	"time"
)

type ThreadedMonteCarloTreeSearchPlayer struct {
	maxSimulations   int
	maxParallelGames int
}

type node struct {
	parent        *node
	N             int
	V             int
	turn          int8
	childrenCount int8
	moves         []Move
	nodes         []*node
}

func NewThreadedMonteCarloTreeSearchPlayer(maxSimulations int, maxParallelGames int) *ThreadedMonteCarloTreeSearchPlayer {
	rand.Seed(time.Now().UnixNano())
	return &ThreadedMonteCarloTreeSearchPlayer{maxSimulations, maxParallelGames}
}

func (player *ThreadedMonteCarloTreeSearchPlayer) SelectMove(game *Game) Move {
	tree := newNode(game, nil)

	var waitgroup sync.WaitGroup
	resultChannel := make(chan *resultTuple)
	gamesChannel := make(chan *gameRequestTuple, player.maxParallelGames)
	go vsUpdater(resultChannel, &waitgroup)
	go gamesRequester(gamesChannel, resultChannel)

	for simulation := 0; simulation < player.maxSimulations; simulation += 1 {
		gameCopy := game.Copy()

		selectedNode := selectNode(gameCopy, tree)
		createdNode := selectedNode.expand(gameCopy)
		updateNs(createdNode)

		waitgroup.Add(1)
		gamesChannel <- &gameRequestTuple{createdNode, gameCopy}
	}

	waitgroup.Wait()
	close(resultChannel)
	close(gamesChannel)

	bestVisitCount, bestNodes := 0, make([]int, 0, 10)
	for i, child := range tree.nodes {
		if child.N > bestVisitCount {
			bestVisitCount = child.N
			bestNodes = append(make([]int, 0, 10), i)
		} else if child.N == bestVisitCount {
			bestNodes = append(bestNodes, i)
		}
	}

	selectedIndex := bestNodes[rand.Intn(len(bestNodes))]
	return tree.moves[selectedIndex]
}

func newNode(game *Game, parent *node) *node {
	moves := game.GetPossibleMoves()
	nodes := make([]*node, len(moves))

	return &node{parent, 0, 0, game.turn, 0, moves, nodes}
}

func selectNode(game *Game, node *node) *node {
	if node.isLeaf() {
		return node
	}

	bestScore, bestNodes := -999999.0, make([]int, 0, 10)
	C := 2.0

	for i, child := range node.nodes {
		v := float64(child.V)
		n := float64(child.N)
		N := float64(node.N)

		score := -v/n + C*math.Sqrt(math.Log(N)/n)

		if score > bestScore {
			bestScore = score
			bestNodes = append(make([]int, 0, 10), i)
		} else if score == bestScore {
			bestNodes = append(bestNodes, i)
		}
	}

	selectedIndex := bestNodes[rand.Intn(len(bestNodes))]
	game.MakeMove(node.moves[selectedIndex])

	return selectNode(game, node.nodes[selectedIndex])
}

func updateNs(node *node) {
	for node != nil {
		node.N += 1
		node = node.parent
	}
}

func updateVs(result int8, node *node) {
	for node != nil {
		node.V += int(node.turn * result)
		node = node.parent
	}
}

func (node *node) isLeaf() bool {
	return int(node.childrenCount) < len(node.moves)
}

func (node *node) expand(game *Game) *node {
	possibilities := len(node.moves) - int(node.childrenCount)
	selected := rand.Intn(possibilities) + 1

	index := 0
	for ; selected > 0; index += 1 {
		if node.nodes[index] == nil {
			selected -= 1
		}
	}
	index -= 1

	game.MakeMove(node.moves[index])
	createdNode := newNode(game, node)
	node.nodes[index] = createdNode
	node.childrenCount += 1

	return createdNode
}

func randomSampleFromState(game *Game) int8 {
	randomPlayer := NewRandomPlayer()
	finished, winner := false, int8(0)

	for !finished {
		move := randomPlayer.SelectMove(game)
		finished, winner = game.MakeMove(move)
	}

	return winner
}

type resultTuple struct {
	node   *node
	result int8
}

func vsUpdater(ch chan *resultTuple, waitgroup *sync.WaitGroup) {
	for tuple := range ch {
		updateVs(tuple.result, tuple.node)
		waitgroup.Done()
	}
}

type gameRequestTuple struct {
	node *node
	game *Game
}

func gamesRequester(gameRequestsChannel chan *gameRequestTuple, resultChannel chan *resultTuple) {
	for tuple := range gameRequestsChannel {
		go gamesPlayer(resultChannel, tuple.node, tuple.game)
	}
}

func gamesPlayer(resultChannel chan *resultTuple, node *node, game *Game) {
	result := randomSampleFromState(game)
	resultChannel <- &resultTuple{node, result}
}
