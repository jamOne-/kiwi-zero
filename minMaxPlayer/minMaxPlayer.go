package minMaxPlayer

import (
	"math"
	"math/rand"

	"github.com/jamOne-/kiwi-zero/game"
)

const INFINITY = 99999999.0

type MinMaxPlayer struct {
	depth   int
	valueFn game.ValueFn
}

func NewMinMaxPlayer(
	depth int,
	valueFn game.ValueFn,
) *MinMaxPlayer {
	return &MinMaxPlayer{depth, valueFn}
}

func (player *MinMaxPlayer) SelectMove(game game.Game) game.Move {
	_, move := negaMax(player.valueFn, game, player.depth, -INFINITY, INFINITY)
	return move
}

func negaMax(
	valueFn game.ValueFn,
	g game.Game,
	depth int,
	a float64,
	b float64,
) (float64, game.Move) {
	if finished, winner := g.IsGameFinished(); finished {
		return INFINITY * float64(winner*g.GetCurrentPlayerColor()), game.Move(-1)
	}

	if depth == 0 {
		return float64(g.GetCurrentPlayerColor()) * valueFn(g), game.Move(-1)
	}

	// movesScores := getMovesAndScores(valueFn, g)
	// sort.Sort(ByScore(movesScores))

	bestValue, bestMoves := -INFINITY, []game.Move{game.Move(-1)}
	// for _, moveScore := range movesScores {
	// 	move := moveScore.move
	moves := g.GetPossibleMoves()
	for _, move := range moves {
		g.MakeMove(move)

		value, _ := negaMax(valueFn, g, depth-1, -b, -a)
		value = -value

		g.UndoLastMove()

		if value > bestValue {
			bestValue = value
			bestMoves = []game.Move{move}
		} else if value == bestValue {
			bestMoves = append(bestMoves, move)
		}

		a = math.Max(a, value)

		if a >= b {
			break
		}
	}

	bestMoveIndex := rand.Intn(len(bestMoves))
	bestMove := bestMoves[bestMoveIndex]
	return bestValue, bestMove
}

type MoveAndScore struct {
	move  game.Move
	score float64
}

type ByScore []*MoveAndScore

func (a ByScore) Len() int           { return len(a) }
func (a ByScore) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByScore) Less(i, j int) bool { return a[j].score < a[i].score }

func getMovesAndScores(valueFn game.ValueFn, g game.Game) []*MoveAndScore {
	scores := make([]*MoveAndScore, 0)
	moves := g.GetPossibleMoves()
	// color := g.GetCurrentPlayerColor()

	for _, move := range moves {
		g.MakeMove(move)
		// score := valueFn(g) * float64(color)
		score := float64(len(g.GetPossibleMoves()))
		g.UndoLastMove()

		scores = append(scores, &MoveAndScore{move, score})
	}

	return scores
}
