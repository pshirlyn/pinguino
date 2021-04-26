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
	g.mu.Lock()
	defer g.mu.Unlock()
	// TODO: handle types of commands
	log.Printf("Received move %s from %s in region %d", move.Command, move.Username, move.Region)
	command := move.Command.(Move)
	user := move.Username
	g.playerState[user] = &PlayerState{command.X, command.Y}
	log.Printf("Player %s moved to (%d, %d)", move.Username, command.X, command.Y)
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
	game.playerState = make(map[string]*PlayerState)
	game.mu.Unlock()

	go game.run()

	return game
}
