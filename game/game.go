package game

type PlayerColor = int8
type Field = int8
type Move = Field
type Features = [][][]float32
type GameToFeaturesFn func(game Game) Features
type ValueFn func(game Game) float32

type Game interface {
	Copy() Game
	MakeMove(move Move) (bool, PlayerColor)
	UndoLastMove()
	GetPossibleMoves() []Move
	GetCurrentPlayerColor() PlayerColor
	IsGameFinished() (bool, PlayerColor)
	DrawBoard()
	SerializeBoard(flipColors bool) string
	OneHotBoard() [][][]float32
}

const WHITE = PlayerColor(-1)
const BLACK = PlayerColor(1)
