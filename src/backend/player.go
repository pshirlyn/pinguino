package pinguino

import (
	"fmt"
	"log"
	"pinguino/src/backend/labgob"
	"pinguino/src/backend/labrpc"
	"sync"
	"sync/atomic"
	"time"
)

type Player struct {
	mu   sync.Mutex
	dead int32

	coordinator *labrpc.ClientEnd
	servers     []*labrpc.ClientEnd
	username    string

	region      int // region connected to
	serverIndex int

	state     *PlayerState
	moveIndex int
}

func (pl *Player) Kill() {
	atomic.StoreInt32(&pl.dead, 1)
}

func (pl *Player) killed() bool {
	z := atomic.LoadInt32(&pl.dead)
	return z == 1
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
	args.Command = MoveCommand{move, pl.username, pl.region}

	reply := StableMoveReply{}

	for !pl.isAssigned() {
		time.Sleep(10 * time.Millisecond) // wait until player is assigned
	}

	log.Println("Sending stable move")

	ok := false
	for i := 0; i < 10; i++ {
		ok = pl.servers[pl.serverIndex].Call("Worker.StableMove", &args, &reply)
		if ok && reply.Success {
			break
		}
	}
}

func (pl *Player) sendFastMove(move interface{}) {
	args := FastMoveArgs{}
	args.Command = MoveCommand{move, pl.username, pl.region}

	reply := FastMoveReply{}

	for !pl.isAssigned() {
		time.Sleep(10 * time.Millisecond) // wait until player is assigned
	}

	log.Println("Sending fast move")

	ok := false
	// TOOD: handle result of call
	for i := 0; i < 10; i++ {
		ok = pl.servers[pl.serverIndex].Call("Worker.FastMove", &args, &reply)
		if ok && reply.Success {
			break
		}
	}
}

// Checks whether the |pl.region| and |pl.serverIndex| has been assigned. This status check verifies whether the player is allowed to start sending moves.
func (pl *Player) isAssigned() bool {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	return pl.region >= 0 && pl.serverIndex >= 0
}

func (pl *Player) getServerIndex() int {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	return pl.serverIndex
}

// RPC handle for a player to be notified of reassignment.
func (pl *Player) WorkerReassignment(args *WorkerReassignmentArgs, reply *WorkerReassignmentReply) {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	pl.serverIndex = args.Worker

	reply.Success = true
	fmt.Printf("NEW SERVER INDEX: %d\n", pl.serverIndex)
}

// Used upon player initialization for a player to request assignment to a reigon and worker
func (pl *Player) requestRegionAssignment() {
	pl.mu.Lock()
	username := pl.username
	pl.mu.Unlock()

	args := AssignPlayerToRegionArgs{Username: username}
	reply := AssignPlayerToRegionReply{}

	// TODO: only tries to contact coordinator for assignment 10 times. might want to see if we can repeatedly send the RPC instead.
	ok := false
	for i := 0; i < 10; i++ {
		ok = pl.coordinator.Call("Coordinator.AssignPlayerToRegion", &args, &reply)
		if ok && reply.Success {
			break
		}
	}

	// TODO: handle !ok and !reply.Success ?

	if ok && reply.Success {
		pl.mu.Lock()
		defer pl.mu.Unlock()
		pl.region = reply.Region
		pl.serverIndex = reply.Worker
	}

}

func MakePlayer(coordinator *labrpc.ClientEnd, servers []*labrpc.ClientEnd, username string) *Player {
	pl := &Player{}

	pl.mu.Lock()
	defer pl.mu.Unlock()

	pl.coordinator = coordinator
	pl.servers = servers
	pl.username = username

	// Temporary initialization
	pl.region = -1
	pl.serverIndex = -1
	pl.moveIndex = 0

	pl.initialize()

	go pl.requestRegionAssignment()

	return pl
}

/////
//
// This is the stuff a developer would need to implement on top of our framework.
//
//
/////

type Move struct {
	X        int    `json:"x"`
	Y        int    `json:"y"`
	Username string `json:"username"`
}

type ChatMessage struct {
	Message  string `json:"message"`
	Username string `json:"username"`
}

func newMove(x int, y int, username string) *Move {
	move := Move{}
	move.X = x
	move.Y = y
	move.Username = username
	return &move
}

func newChatMessage(msg string, username string) *ChatMessage {
	chatMessage := ChatMessage{}
	chatMessage.Message = msg
	chatMessage.Username = username
	return &chatMessage
}

func (pl *Player) ClientMovePlayer(x int, y int) {
	playerMove := newMove(x, y, pl.username)

	pl.sendFastMove(playerMove) // handle reply

	pl.mu.Lock()
	pl.state.x = x
	pl.state.y = y
	pl.mu.Unlock()
}

func (pl *Player) initialize() { // don't hold lock
	labgob.Register(Move{})
	labgob.Register(ChatMessage{})

	pl.state = &PlayerState{0, 0} // initialize to 0, 0
}

func (pl *Player) SendChatMessage(message string) {
	chatMessage := newChatMessage(message, pl.username)

	log.Println(pl.username, "sending chat message")

	pl.sendStableMove(chatMessage)
}
