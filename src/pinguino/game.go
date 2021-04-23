package pinguino

import (
	"log"
	"sync"
)

type PlayerState struct {
	x int
	y int
}

type Game struct {
	mu sync.Mutex

	moveCh      chan MoveCommand // channel where it receives all the moves
	playerState map[string]*PlayerState
}

func (g *Game) handleMessage(move MoveCommand) {
	// TODO: handle types of commands
	log.Printf("Received move %s from %s in region %d", move.Command, move.Username, move.Region)
	g.playerState[move.Username].x = move.Command.X
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

	go game.run()

	return game
}
