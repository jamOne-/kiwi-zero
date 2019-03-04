package main

type Player interface {
	SelectMove(game *Game) Move
}
