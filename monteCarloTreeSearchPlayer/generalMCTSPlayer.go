package monteCarloTreeSearchPlayer

import (
	"fmt"
	"math"

	"github.com/jamOne-/kiwi-zero/game"
	"github.com/jamOne-/kiwi-zero/player"
	"github.com/jamOne-/kiwi-zero/utils"
)

type TempFn func(turnNumber int) float64
type ValueAndPolicyFn func(game game.Game) (float64, []float32)

func DefaultTempFn(turnNumber int) float64 {
	return 0.1
}

func DefaultValueAndPolicyFn(game game.Game) (float64, []float32) {
	policyLength := game.GetMaxPossibleMoves()
	policy := make([]float32, policyLength)
	for i := 0; i < policyLength; i++ {
		policy[i] = 1
	}

	return 0, policy
}

type StateAction struct {
	State  string
	Action game.Move
}

type GeneralMCTSPlayer struct {
	C                float64
	maxSimulations   int
	rolloutDepth     int
	rolloutPlayer    player.Player
	valueFn          game.ValueFn
	valueAndPolicyFn ValueAndPolicyFn
	tempFn           TempFn
}

func NewGeneralMCTSPlayer(
	maxSimulations int,
	C float64,
	rolloutDepth int,
	rolloutPlayer player.Player,
	valueFn game.ValueFn,
	valueAndPolicyFn ValueAndPolicyFn,
	tempFn TempFn,
) *GeneralMCTSPlayer {
	if valueAndPolicyFn == nil {
		valueAndPolicyFn = DefaultValueAndPolicyFn
	}
	if tempFn == nil {
		tempFn = DefaultTempFn
	}

	return &GeneralMCTSPlayer{C, maxSimulations, rolloutDepth, rolloutPlayer, valueFn, valueAndPolicyFn, tempFn}
}

func (player *GeneralMCTSPlayer) SelectMove(game game.Game) game.Move {
	move, _ := player.SelectMoveWithPolicy(game)
	return move
}

func (player *GeneralMCTSPlayer) SelectMoveWithPolicy(g game.Game) (game.Move, []float32) {
	Ns := make(map[string]int)
	Nsa := make(map[StateAction]int)
	Qsa := make(map[StateAction]float64)
	Psa := make(map[StateAction]float32)
	Moves := make(map[string][]game.Move)

	for simulation := 0; simulation < player.maxSimulations; simulation += 1 {
		player.OneSimulation(g, Ns, Nsa, Qsa, Psa, Moves)
	}

	t := player.tempFn(g.GetTurnNumber())
	t = math.Max(t, 0.1)
	T := 1.0 / t

	policy := make([]float32, g.GetMaxPossibleMoves())
	repr := g.SerializeBoard(false)
	moves := Moves[repr]
	for _, move := range moves {
		key := StateAction{repr, move}
		policy[move+1] = float32(math.Pow(float64(Nsa[key])/float64(Ns[repr]), T))
	}

	sum := utils.SumFloats32(policy)
	scaleFactor := float32(1.0) / sum
	for i, _ := range policy {
		policy[i] *= scaleFactor
	}

	move := utils.RandomFromDistribution(policy) - 1

	return int8(move), policy
}

func (player *GeneralMCTSPlayer) OneSimulation(
	g game.Game,
	Ns map[string]int,
	Nsa map[StateAction]int,
	Qsa map[StateAction]float64,
	Psa map[StateAction]float32,
	Moves map[string][]game.Move,
) float64 {
	repr := g.SerializeBoard(false)

	if _, exists := Ns[repr]; !exists {
		if finished, winner := g.IsGameFinished(); finished {
			return float64(g.GetCurrentPlayerColor() * winner)
		}

		return player.NewNode(g, repr, Ns, Psa, Moves)
	}

	if finished, _ := g.IsGameFinished(); finished {
		fmt.Println("FINISHED WTF")
	}

	bestScore, bestMove := -9999999.0, game.Move(-1)
	N := float64(Ns[repr])

	for _, move := range Moves[repr] {
		saKey := StateAction{repr, move}
		q := float64(Qsa[saKey])
		n := float64(Nsa[saKey])
		p := float64(Psa[saKey])

		score := q + player.C*p*math.Sqrt(N)/(1+n)

		if score > bestScore {
			bestScore, bestMove = score, move
		}
	}

	g.MakeMove(bestMove)
	v := -player.OneSimulation(g, Ns, Nsa, Qsa, Psa, Moves)
	g.UndoLastMove()

	saKey := StateAction{repr, bestMove}
	nsa := float64(Nsa[saKey])
	Qsa[saKey] = (Qsa[saKey]*nsa + v) / (nsa + 1)
	Nsa[saKey] = int(nsa + 1)
	Ns[repr] += 1

	return v
}

func (player *GeneralMCTSPlayer) NewNode(
	game game.Game,
	repr string,
	Ns map[string]int,
	Psa map[StateAction]float32,
	Moves map[string][]game.Move,
) float64 {
	v, policy := player.valueAndPolicyFn(game)
	moves := game.GetPossibleMoves()

	Ns[repr] = 1
	Moves[repr] = moves

	validSum := float32(0)
	for _, move := range moves {
		validSum += policy[move+1] // todo: pass
	}
	scaleFactor := float32(1.0) / validSum

	for _, move := range moves {
		key := StateAction{repr, move}
		Psa[key] = policy[move+1] * scaleFactor // todo: pass
	}

	if player.rolloutDepth != 0 {
		v = rollout(player.rolloutDepth, player.rolloutPlayer, player.valueFn, game)
	}

	return v * float64(game.GetCurrentPlayerColor())
}

func rollout(maxDepth int, player player.Player, valueFn game.ValueFn, game game.Game) float64 {
	finished, winner := false, int8(0)
	steps := 0

	for depth := 0; depth < maxDepth && !finished; depth += 1 {
		move := player.SelectMove(game)
		finished, winner = game.MakeMove(move)
		steps += 1
	}

	for i := 0; i < steps; i += 1 {
		game.UndoLastMove()
	}

	if finished {
		return float64(winner)
	}

	value := valueFn(game) // returns value from [-1;1]
	return value
}
