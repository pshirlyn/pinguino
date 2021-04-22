package pinguino

import (
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"pinguino/src/labrpc"
	"sync"
	"time"
)

type Coordinator struct {
	mu sync.Mutex

	nRegions          int
	workers           []*labrpc.ClientEnd
	playerToRegionMap map[string]int
	regionToWorkerMap map[int]int

	killed bool
}

// Used to pick a region to assign to players.
// For now, pick the region randomly. In the future, consider load balancing.
func (c *Coordinator) pickRegion() int {
	return rand.Intn(c.nRegions)
}

func (c *Coordinator) AssignPlayerToRegion(args *AssignPlayerToRegionArgs, reply *AssignPlayerToRegionReply) {
	c.mu.Lock()
	defer c.mu.Unlock()

	username := args.Username
	region := c.pickRegion()
	c.playerToRegionMap[username] = region

	reply.Success = true
	reply.Region = region
	reply.Worker = c.regionToWorkerMap[region]
}

// to implement:
// func (c *Coordinator) MovePlayer()

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

	// TODO: assign workers to region and set c.regionToWorkerMap
	// Remember that workers[0] is set to be nil
	c.workers = workers
}

func (c *Coordinator) sendHeartbeatToWorker(worker *labrpc.ClientEnd, args *HeartbeatArgs, reply *HeartbeatReply) {
	ok := worker.Call("Worker.Heartbeat", &args, &reply)

	if !ok {
		log.Println("couldn't reach worker")
		// TODO: handle worker disconnect
	}
}

func (c *Coordinator) SendHeartbeats() {
	for _, worker := range c.workers {
		args := HeartbeatArgs{}
		reply := HeartbeatReply{}
		go c.sendHeartbeatToWorker(worker, &args, &reply)
	}
}

func (c *Coordinator) run() {
	for !c.killed {
		c.SendHeartbeats()
		time.Sleep(10 * time.Millisecond)
	}
}

func MakeCoordinator(regions int) *Coordinator {
	c := &Coordinator{}
	c.nRegions = regions
	c.playerToRegionMap = make(map[string]int)
	c.regionToWorkerMap = make(map[int]int)
	c.killed = false

	go c.run()
	return c
}
