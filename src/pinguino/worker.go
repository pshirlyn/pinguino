package pinguino

import (
	"bytes"
	"fmt"
	"log"
	"net/rpc"
	"pinguino/src/labgob"
	"pinguino/src/labrpc"
	"sync"
)

type MoveCommand struct {
	Command  interface{}
	Username string
	Region   int
}

type Replica struct {
	Data []byte
}

type Worker struct {
	mu sync.Mutex

	peers         []*labrpc.ClientEnd
	me            int
	replicas      []int
	localReplicas map[int]*Replica
	// replicaStates []*[]*PlayerState

	moveIndex int
	log       []*MoveCommand
	killed    bool

	gameChannel chan MoveCommand
	game        *Game
}

func (wk *Worker) StableMove(args *StableMoveArgs, reply *StableMoveReply) {
	wk.mu.Lock()
	wk.log = append(wk.log, &args.Command)
	reply.Success = true
	wk.moveIndex++
	wk.mu.Unlock()

	wk.gameChannel <- args.Command
}

func (wk *Worker) FastMove(args *FastMoveArgs, reply *FastMoveReply) {
	wk.mu.Lock()
	wk.log = append(wk.log, &args.Command)
	reply.Success = true
	wk.moveIndex++
	wk.mu.Unlock()

	wk.gameChannel <- args.Command
}

func (wk *Worker) Heartbeat(args *HeartbeatArgs, reply *HeartbeatReply) {
	reply.Success = true
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

func (wk *Worker) SendReplica(args *SendReplicaArgs, reply *SendReplicaReply) {
	wk.mu.Lock()
	wk.localReplicas[args.Worker] = &Replica{args.Replica}
	wk.mu.Unlock()
}

func (wk *Worker) getSnapshot() []byte {
	w := new(bytes.Buffer)
	e := labgob.NewEncoder(w)
	e.Encode(wk.log)
	e.Encode(wk.game.chatLog)
	data := w.Bytes()
	return data
}

func (wk *Worker) sendUpdateReplica(replica int) {
	// TODO: smart diff
	// for now just send everything
	args := SendReplicaArgs{}
	args.Replica = wk.getSnapshot()
	args.Worker = wk.me
	args.MoveIndex = wk.moveIndex
	reply := SendReplicaReply{}
	wk.peers[replica].Call("Worker.SendReplica", &args, &reply)

}

func (wk *Worker) sendToReplicas() {
	for replica := range wk.replicas {
		go wk.sendUpdateReplica(replica)
	}
}

func (wk *Worker) SetReplicas(replicas []int) {
	wk.replicas = replicas
	wk.sendToReplicas()
}

func MakeWorker(coordinator *labrpc.ClientEnd, peers []*labrpc.ClientEnd, me int) *Worker {
	wk := &Worker{}

	wk.mu.Lock()
	defer wk.mu.Unlock()

	wk.me = me
	wk.peers = peers
	wk.gameChannel = make(chan MoveCommand)
	wk.game = MakeGame(wk.gameChannel)
	wk.replicas = make([]int, 2)

	return wk
}
