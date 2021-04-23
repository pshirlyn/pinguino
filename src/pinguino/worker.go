package pinguino

import (
	"fmt"
	"log"
	"net/rpc"
	"pinguino/src/labrpc"
	"sync"
)

type MoveCommand struct {
	Command  interface{}
	Username string
	Region   int
}

type Worker struct {
	mu sync.Mutex

	peers []*labrpc.ClientEnd
	me    int

	log    []*MoveCommand
	killed bool

	gameChannel chan MoveCommand
	game        *Game
}

func (wk *Worker) StableMove(args *StableMoveArgs, reply *StableMoveReply) {
	// TODO
}

func (wk *Worker) FastMove(args *FastMoveArgs, reply *FastMoveReply) {
	wk.mu.Lock()
	wk.log = append(wk.log, &args.Command)
	wk.mu.Unlock()

	wk.gameChannel <- args.Command
}

func (wk *Worker) Heartbeat(args *HeartbeatArgs, reply *HeartbeatReply) {

}

func call(rpcname string, args interface{}, reply interface{}) bool {
	// c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	sockname := coordinatorSock()
	c, err := rpc.DialHTTP("unix", sockname)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}

func (wk *Worker) Kill() {
	wk.killed = true // change to atomic write
}

func MakeWorker(coordinator *labrpc.ClientEnd, peers []*labrpc.ClientEnd, me int) *Worker {
	wk := &Worker{}

	wk.mu.Lock()
	defer wk.mu.Unlock()

	wk.me = me
	wk.peers = peers
	wk.gameChannel = make(chan MoveCommand)
	wk.game = MakeGame(wk.gameChannel)

	return wk
}
