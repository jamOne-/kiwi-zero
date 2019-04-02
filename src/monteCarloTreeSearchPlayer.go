package main

import (
	"math"
	"math/rand"
	"time"
)

type MonteCarloTreeSearchPlayer struct {
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

func NewMonteCarloTreeSearchPlayer() *MonteCarloTreeSearchPlayer {
	rand.Seed(time.Now().UnixNano())
	return &MonteCarloTreeSearchPlayer{}
}

func (player *MonteCarloTreeSearchPlayer) SelectMove(game *Game) Move {
	MAX_SIMULATIONS := 300
	tree := newNode(game, nil)

	for simulation := 0; simulation < MAX_SIMULATIONS; simulation += 1 {
		gameCopy := game.Copy()

		selectedNode := selectNode(gameCopy, tree)
		createdNode := selectedNode.expand(gameCopy)
		updateNs(createdNode)

		result := randomSampleFromState(gameCopy)
		updateVs(result, createdNode)
	}

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
