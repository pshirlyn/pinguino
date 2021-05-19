package pinguino

import (
	"fmt"
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
	chatLog     []*ChatMessage
}

func (g *Game) switchMessageType(msg interface{}, username string) {
	switch command := msg.(type) {
	case Move:
		g.processMove(command, username)
	case ChatMessage:
		g.processChatMessage(command, username)
	default:
		fmt.Println("Unknown type!")
	}
}

func (g *Game) processMove(command Move, username string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.playerState[username] = &PlayerState{command.X, command.Y}
	log.Printf("Player %s moved to (%d, %d)", username, command.X, command.Y)
}

func (g *Game) processChatMessage(msg ChatMessage, username string) {
	g.mu.Lock()
	g.chatLog = append(g.chatLog, &msg)
	g.mu.Unlock()
	log.Printf("Player %s typed: %s", msg.Username, msg.Message)
}

func (g *Game) handleMessage(move MoveCommand) {
	g.switchMessageType(move.Command, move.Username)
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
