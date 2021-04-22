package pinguino

import (
	"pinguino/src/labrpc"
	"sync"
)

type Player struct {
	mu sync.Mutex

	coordinator *labrpc.ClientEnd
	servers     []*labrpc.ClientEnd
	username    string

	region      int // region connected to
	serverIndex int
}

func (pl *Player) Kill() {
}

func (pl *Player) SetWorkers(servers []*labrpc.ClientEnd, region int, serverIndex int) {
	// In case workers change, we update our list of worker ClientEnds
	pl.mu.Lock()
	pl.servers = servers
	pl.region = region
	pl.serverIndex = serverIndex
	pl.mu.Unlock()
}

func (pl *Player) sendStableMove(move interface{}) {
	args := StableMoveArgs{}
	args.Move = move

	reply := StableMoveReply{}

	pl.servers[pl.serverIndex].Call("Worker.StableMove", &args, &reply)
	// TOOD: handle result of call
}

func (pl *Player) sendFastMove(move interface{}) {
	args := FastMoveArgs{}
	args.Move = move

	reply := FastMoveReply{}

	pl.servers[pl.serverIndex].Call("Worker.FastMove", &args, &reply)
	// TOOD: handle result of call
}

func MakePlayer(coordinator *labrpc.ClientEnd, servers []*labrpc.ClientEnd, username string) *Player {
	pl := &Player{}

	pl.mu.Lock()
	defer pl.mu.Unlock()

	pl.coordinator = coordinator
	pl.servers = servers
	pl.username = username

	return pl
}

/////
//
// This is the stuff a developer would need to implement on top of our framework.
//
//
/////

type Move struct {
	X        int
	Y        int
	Username string
}

func NewMove(x int, y int, username string) *Move {
	move := Move{}
	move.X = x
	move.Y = y
	move.Username = username
	return &move
}

func (pl *Player) move(x int, y int) {
	playerMove := NewMove(x, y, pl.username)

	pl.sendFastMove(playerMove)
}
