package monteCarloTreeSearchPlayer

import (
	"math/rand"

	"github.com/jamOne-/kiwi-zero/game"
	"github.com/jamOne-/kiwi-zero/player"
)

type GeneralMCTSPlayer struct {
	maxSimulations int
	rolloutDepth   int
	rolloutPlayer  player.Player
	valueFn        game.ValueFn
}

func NewGeneralMCTSPlayer(
	maxSimulations int,
	rolloutDepth int,
	rolloutPlayer player.Player,
	valueFn game.ValueFn,
) *GeneralMCTSPlayer {
	return &GeneralMCTSPlayer{maxSimulations, rolloutDepth, rolloutPlayer, valueFn}
}

func (player *GeneralMCTSPlayer) SelectMove(game game.Game) game.Move {
	tree := newNode(game, nil)

	for simulation := 0; simulation < player.maxSimulations; simulation += 1 {
		// TODO: know number of moves made and run some undos
		gameCopy := game.Copy()

		selectedNode := selectNode(gameCopy, tree)
		createdNode := selectedNode.expand(gameCopy)

		result := rollout(player.rolloutDepth, player.rolloutPlayer, player.valueFn, gameCopy)
		updateNsAndVs(result, createdNode)
	}

	bestVisitCount, bestNodes := 0, make([]int, 0)
	for i, child := range tree.Nodes {
		if child.N > bestVisitCount {
			bestVisitCount = child.N
			bestNodes = append(make([]int, 0), i)
		} else if child.N == bestVisitCount {
			bestNodes = append(bestNodes, i)
		}
	}

	selectedIndex := bestNodes[rand.Intn(len(bestNodes))]
	return tree.moves[selectedIndex]
}

func rollout(maxDepth int, player player.Player, valueFn game.ValueFn, game game.Game) int8 {
	finished, winner := false, int8(0)

	for depth := 0; depth < maxDepth && !finished; depth += 1 {
		move := player.SelectMove(game)
		finished, winner = game.MakeMove(move)
	}

	if finished {
		return winner
	}

	probability := valueFn(game)

	if probability > 0.5 {
		return 1
	} else {
		return -1
	}
}
