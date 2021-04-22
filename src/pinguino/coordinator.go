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
)

type Coordinator struct {
	mu sync.Mutex

	nRegions          int
	workers           []*labrpc.ClientEnd
	playerToRegionMap map[string]int
	regionToWorkerMap map[int]int
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

func MakeCoordinator(workers []*labrpc.ClientEnd, nregions int) *Coordinator {
	c := &Coordinator{}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.nRegions = nregions
	// TODO: add coordinator backup server reference here
	c.workers = workers

	c.playerToRegionMap = make(map[string]int)
	c.regionToWorkerMap = make(map[int]int)

	// Assign main worker to each region.
	// TODO: handle cases where the number of regions != number of workers
	// TODO: perhaps randomize the assignment
	for i := 0; i < c.nRegions; i++ {
		c.regionToWorkerMap[i] = i

	}

	return c
}
