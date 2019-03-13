package main

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

func (player *MinMaxPlayer) SelectMove(game *Game) Move {
	_, move := negaMax(game, player.depth, -INFINITY, INFINITY)
	return move
}

func negaMax(game *Game, depth int, a int, b int) (int, Move) {
	if finished, winner := game.IsGameFinished(); finished {
		return INFINITY * int(winner) * int(game.turn), Move(-1)
	}

	if depth == 0 {
		return heuristicValueFunction(game), Move(-1)
	}

	moves := game.GetPossibleMoves()
	bestValue, bestMove := -INFINITY, Move(-1)

	for _, move := range moves {
		gameCopy := game.Copy()
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

func heuristicValueFunction(game *Game) int {
	blacks, whites := 0, 0
	blackScore, whiteScore := 0, 0

	for i, pawn := range game.board {
		if pawn == BLACK {
			blacks += 1
			blackScore += SCORING[i]
		} else if pawn == WHITE {
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

	return int(game.turn) * (p + blacks - whites)
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
