package game

type PlayerColor = int8
type Field = int8
type Move = Field
type Features = [][][]float32
type GameToFeaturesFn func(game Game) Features
type ValueFn func(game Game) float64

type Game interface {
	Copy() Game
	MakeMove(move Move) (bool, PlayerColor)
	UndoLastMove()
	GetPossibleMoves() []Move
	GetCurrentPlayerColor() PlayerColor
	GetTurnNumber() int
	GetMaxPossibleMoves() int
	IsGameFinished() (bool, PlayerColor)
	DrawBoard()
	SerializeBoard(flipColors bool) string
	OneHotBoard() [][][]float32
	EncodeMoveToPolicy(move Move) []float32
}

const WHITE = PlayerColor(-1)
const BLACK = PlayerColor(1)
