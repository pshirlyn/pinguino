package pinguino

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"pinguino/labrpc"
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

func MakeCoordinator(workers []*labrpc.ClientEnd) *Coordinator {
	cr := &Coordinator{}

	cr.mu.Lock()
	defer cr.mu.Unlock()

	// TODO: add coordinator backup server reference here

	cr.workers = workers

	return cr
}
