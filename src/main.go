package main

import "fmt"

func main() {
	game := NewGame()
	game.DrawBoard()
	fmt.Println("")
	game.MakeMove(19)
	game.DrawBoard()
	fmt.Println(game.GetPossibleMoves())
}
