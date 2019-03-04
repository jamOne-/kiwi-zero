package main

import (
	"math/rand"
	"time"
)

type RandomPlayer struct{}

func NewRandomPlayer() *RandomPlayer {
	rand.Seed(time.Now().UnixNano())
	return &RandomPlayer{}
}

func (player *RandomPlayer) SelectMove(game *Game) Move {
	possibleMoves := game.GetPossibleMoves()
	return possibleMoves[rand.Intn(len(possibleMoves))]
}
