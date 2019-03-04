package main

import "fmt"

type HumanPlayer struct{}

func NewHumanPlayer() *HumanPlayer {
	return &HumanPlayer{}
}

func (player *HumanPlayer) SelectMove(game *Game) Move {
	var move int

	fmt.Println(game.GetPossibleMoves())
	fmt.Scan(&move)

	return int8(move)
}
