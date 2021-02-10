package monteCarloTreeSearchPlayer

import (
	"math"
	"math/rand"

	"github.com/jamOne-/kiwi-zero/game"
)

type LazyExpandMonteCarloTreeSearchPlayer struct {
	maxSimulations int
}

type LazyNode struct {
	parent        *LazyNode
	N             int
	V             int
	currentPlayer int8
	leaf          bool
	moves         []game.Move
	Nodes         []*LazyNode
}

func NewLazyExpandMonteCarloTreeSearchPlayer(maxSimulations int) *LazyExpandMonteCarloTreeSearchPlayer {
	return &LazyExpandMonteCarloTreeSearchPlayer{maxSimulations}
}

func (player *LazyExpandMonteCarloTreeSearchPlayer) SelectMove(game game.Game) game.Move {
	tree := newLazyNode(game, nil)

	for simulation := 0; simulation < player.maxSimulations; simulation += 1 {
		gameCopy := game.Copy()

		selectedNode, _ := selectLazyNode(gameCopy, tree, 0)
		createdNode := selectedNode.expand(gameCopy)

		result := randomSampleFromState(gameCopy)
		updateLazyNsAndVs(result, createdNode)
	}

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
	return tree.moves[selectedIndex]
}

func newLazyNode(game game.Game, parent *LazyNode) *LazyNode {
	moves := game.GetPossibleMoves()
	nodes := make([]*LazyNode, len(moves))

	return &LazyNode{parent, 0, 0, game.GetCurrentPlayerColor(), true, moves, nodes}
}

func selectLazyNode(game game.Game, node *LazyNode, stepsAccumulator int) (*LazyNode, int) {
	if node.leaf {
		return node, stepsAccumulator
	}

	bestScore, bestNodes := -999999.0, make([]int, 0)
	C := 2.0

	for i, child := range node.Nodes {
		v := float64(child.V)
		n := math.Max(float64(child.N), 1)
		N := float64(node.N)

		score := -v/n + C*math.Sqrt(math.Log(N)/n)

		if score > bestScore {
			bestScore = score
			bestNodes = append(make([]int, 0), i)
		} else if score == bestScore {
			bestNodes = append(bestNodes, i)
		}
	}

	selectedIndex := bestNodes[rand.Intn(len(bestNodes))]
	game.MakeMove(node.moves[selectedIndex])

	return selectLazyNode(game, node.Nodes[selectedIndex], stepsAccumulator+1)
}

func updateLazyNsAndVs(result int8, node *LazyNode) {
	for node != nil {
		node.N += 1
		node.V += int(node.currentPlayer * result)
		node = node.parent
	}
}

func (node *LazyNode) expand(game game.Game) *LazyNode {
	for i, move := range node.moves {
		game.MakeMove(move)
		node.Nodes[i] = newLazyNode(game, node)
		game.UndoLastMove()
	}

	node.leaf = false
	selectedIndex := rand.Intn(len(node.moves))
	game.MakeMove(node.moves[selectedIndex])

	return node.Nodes[selectedIndex]
}
