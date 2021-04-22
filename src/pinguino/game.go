package pinguino

import (
	"fmt"
	"sync"
)

type Game struct {
	mu sync.Mutex

	moveCh chan MoveCommand // channel where it receives all the moves
}

func (g *Game) handleMessage(move MoveCommand) {
	fmt.Println("Received move %s from %s in region", move.Command, move.Username, move.Region)
}

func (g *Game) run() {
	for {
		move := <-g.moveCh
		go g.handleMessage(move)
	}
}

func MakeGame(moveCh chan MoveCommand) *Game {
	game := &Game{}

	game.mu.Lock()
	game.moveCh = moveCh

	return game
}
