package backend

import (
	"bytes"
	"fmt"
	"log"
	"net/rpc"
	"pinguino/src/backend/labgob"
	"pinguino/src/backend/labrpc"
	"sync"
	"sync/atomic"
)

type MoveCommand struct {
	Command  interface{}
	Username string
	Region   int
}

type Replica struct {
	Data      []byte
	Moves     []*MoveCommand
	MoveIndex int
}

type Worker struct {
	mu sync.Mutex

	peers         []*labrpc.ClientEnd
	me            int
	replicas      []int
	localReplicas map[int]*Replica

	moveIndex int
	log       []*MoveCommand
	dead      int32

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

func (wk *Worker) NewWorkersAdded(workers []*labrpc.ClientEnd) {
	wk.mu.Lock()
	defer wk.mu.Unlock()

	wk.peers = append(wk.peers, workers...)
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
	atomic.StoreInt32(&wk.dead, 1)
}

func (wk *Worker) killed() bool {
	z := atomic.LoadInt32(&wk.dead)
	return z == 1
}

func (wk *Worker) sendMoveToReplica(replica int, move *MoveCommand) {
	args := SendReplicaArgs{}
	args.Worker = wk.me
	args.MoveIndex = wk.moveIndex
	args.Move = *move

	reply := SendReplicaReply{}

	ok := false
	for i := 0; i < 10; i++ {
		wk.peers[replica].Call("Worker.SendReplica", &args, &reply)
		if ok && reply.Success {
			break
		}
	}

}

func (wk *Worker) SendReplica(args *SendReplicaArgs, reply *SendReplicaReply) {
	wk.mu.Lock()
	defer wk.mu.Unlock()

	if oldReplica, ok := wk.localReplicas[args.Worker]; ok {
		if args.MoveIndex < oldReplica.MoveIndex {
			return // outdated replica sent
		}

		if args.Replica != nil {
			oldReplica.Data = args.Replica
		}

		if args.Move.Command != nil {
			oldReplica.Moves = append(oldReplica.Moves, &args.Move)
		}
	}

	replica := Replica{}

	if args.Replica != nil {
		replica.Data = args.Replica
	}

	if args.Move.Command != nil {
		replica.Moves = append(replica.Moves, &args.Move)
	}

	wk.localReplicas[args.Worker] = &replica

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

	wk.mu.Lock()
	wk.log = make([]*MoveCommand, 0)
	wk.mu.Unlock()
}

func (wk *Worker) SetReplicas(args *SetReplicasArgs, reply *SetReplicasReply) {
	wk.mu.Lock()
	replicas := args.Replicas
	wk.replicas = replicas
	wk.mu.Unlock()
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
