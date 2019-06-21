package game

type PlayerColor = int8
type Field = int8
type Move = Field

type Game interface {
	Copy() Game
	MakeMove(move Move) (bool, PlayerColor)
	GetPossibleMoves() []Move
	GetCurrentPlayerColor() PlayerColor
	IsGameFinished() (bool, PlayerColor)
	DrawBoard()
	SerializeBoard(flipColors bool) string
	OneHotBoard() [][][]float32
}
