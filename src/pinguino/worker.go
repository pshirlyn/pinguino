package pinguino

import (
	"fmt"
	"labrpc"
	"log"
	"net/rpc"
	"sync"
)

type Worker struct {
	mu sync.Mutex
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

func SetPeers(peers []*labrpc.ClientEnd) {
	worker.mu.Lock()
	worker.peers = peers
	worker.mu.Unlock()
}

func MakeWorker(coordinator *labrpc.ClientEnd, me int) *Worker {
	worker := &Worker{}

	worker.mu.Lock()
	defer worker.mu.Unlock()

	worker.me = me

	return worker
}
