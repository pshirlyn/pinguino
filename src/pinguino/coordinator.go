package pinguino

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"pinguino/src/labrpc"
	"sync"
)

type Coordinator struct {
	mu sync.Mutex

	workers []*labrpc.ClientEnd
}

//
// start a thread that listens for RPCs from worker.go
//
func (c *Coordinator) server() {
	rpc.Register(c)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := coordinatorSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

func (c *Coordinator) Kill() {
}

func (c *Coordinator) SetWorkers(workers []*labrpc.ClientEnd) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// TODO: add coordinator backup server reference here

	c.workers = workers
}

func MakeCoordinator() *Coordinator {
	cr := &Coordinator{}

	return cr
}
