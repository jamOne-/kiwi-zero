package minMaxPlayer

import (
	"github.com/jamOne-/kiwi-zero/game"
	"github.com/jamOne-/kiwi-zero/reversi"
)

const INFINITY = 99999999

type MinMaxPlayer struct {
	depth int
}

func max(x int, y int) int {
	if x >= y {
		return x
	} else {
		return y
	}
}

func NewMinMaxPlayer(depth int) *MinMaxPlayer {
	return &MinMaxPlayer{depth}
}

func (player *MinMaxPlayer) SelectMove(game *reversi.Game) game.Move {
	_, move := negaMax(game, player.depth, -INFINITY, INFINITY)
	return move
}

func negaMax(g *reversi.Game, depth int, a int, b int) (int, game.Move) {
	if finished, winner := g.IsGameFinished(); finished {
		return INFINITY * int(winner) * int(g.Turn), game.Move(-1)
	}

	if depth == 0 {
		return heuristicValueFunction(g), game.Move(-1)
	}

	moves := g.GetPossibleMoves()
	bestValue, bestMove := -INFINITY, game.Move(-1)

	for _, move := range moves {
		gameCopy := g.Copy().(*reversi.Game)
		gameCopy.MakeMove(move)

		value, _ := negaMax(gameCopy, depth-1, -b, -a)
		value = -value

		if value > bestValue {
			bestValue = value
			bestMove = move
		}

		a = max(a, value)

		if a >= b {
			break
		}
	}

	return bestValue, bestMove
}

func heuristicValueFunction(game *reversi.Game) int {
	blacks, whites := 0, 0
	blackScore, whiteScore := 0, 0

	for i, pawn := range game.Board {
		if pawn == reversi.BLACK {
			blacks += 1
			blackScore += SCORING[i]
		} else if pawn == reversi.WHITE {
			whites += 1
			whiteScore += SCORING[i]
		}
	}

	p := 0
	if blacks > whites {
		p = 100.0 * blacks / (blacks + whites)
	} else if blacks < whites {
		p = -100.0 * whites / (blacks + whites)
	}

	return int(game.Turn) * (p + blacks - whites)
}

var SCORING = []int{
	20, -3, 11, 8, 8, 11, -3, 20,
	-3, -7, -4, 1, 1, -4, -7, -3,
	11, -4, 2, 2, 2, 2, -4, 11,
	8, 1, 2, -3, -3, 2, 1, 8,
	8, 1, 2, -3, -3, 2, 1, 8,
	11, -4, 2, 2, 2, 2, -4, 11,
	-3, -7, -4, 1, 1, -4, -7, -3,
	20, -3, 11, 8, 8, 11, -3, 20}
