package pinguino

import (
	"fmt"
	"log"
	"net/rpc"
	"pinguino/src/labrpc"
	"sync"
)

type Worker struct {
	mu sync.Mutex

	peers []*labrpc.ClientEnd
	me    int
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
}

func (wk *Worker) StableMove(args *StableMoveArgs, reply *StableMoveReply) {

}

func (wk *Worker) FastMove(args *FastMoveArgs, reply *FastMoveReply) {

}

func MakeWorker(coordinator *labrpc.ClientEnd, peers []*labrpc.ClientEnd, me int) *Worker {
	wk := &Worker{}

	wk.mu.Lock()
	defer wk.mu.Unlock()

	wk.me = me
	wk.peers = peers

	return wk
}
