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

const heartbeatTimeInterval = time.Duration(10 * time.Millisecond)

type Coordinator struct {
	mu sync.Mutex

	nRegions          int
	workers           []*labrpc.ClientEnd
	lastHeartbeats    []time.Time
	playerToRegionMap map[string]int
	regionToWorkerMap map[int]int
	backup            int

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

func (c *Coordinator) sendHeartbeatToWorker(workerIndex int, args *HeartbeatArgs, reply *HeartbeatReply) {
	ok := c.workers[workerIndex].Call("Worker.Heartbeat", args, reply)
	// fmt.Printf("Heartbeat: S%d\n", workerIndex+1)
	if !ok {
		log.Println("couldn't reach worker")
		// TODO: handle worker disconnect
		return
	}

	// Successfully sent heartbeat to worker, so update lastHearbeat
	c.lastHeartbeats[workerIndex] = time.Now()
}

func (c *Coordinator) maybeSendHeartbeats() {
	for i := 0; i < len(c.workers); i++ {
		args := HeartbeatArgs{}
		reply := HeartbeatReply{}
		// Check time since last hearbeat
		if time.Since(c.lastHeartbeats[i]) > heartbeatTimeInterval {
			go c.sendHeartbeatToWorker(i, &args, &reply)
		}
	}
}

// to implement:
// func (c *Coordinator) MovePlayer()

// func (c *Coordinator) NewWorkersAdded(args *NewWorkersAddedArgs, reply *NewWorkersAddedReply) {
func (c *Coordinator) NewWorkersAdded(workers []*labrpc.ClientEnd) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.workers = append(c.workers, workers...)

	// TODO: rather than just adding a new region, check if there's a worker that is handling multiple regions first
	for i := 0; i < len(workers); i++ {
		c.regionToWorkerMap[c.nRegions+i] = c.nRegions + 1
		c.lastHeartbeats = append(c.lastHeartbeats, time.Now())
	}
	c.nRegions += len(workers)
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

func (c *Coordinator) run() {
	// main loop
	for !c.killed {
		c.maybeSendHeartbeats()
		time.Sleep(heartbeatTimeInterval)
	}
}

func (c *Coordinator) SelectBackup() {
	// Selects backup out of existing servers
	idx := rand.Intn(len(c.workers))
	c.backup = idx
}

func MakeCoordinator(workers []*labrpc.ClientEnd, regions int) *Coordinator {
	c := &Coordinator{}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.nRegions = regions
	c.killed = false
	c.backup = -1

	c.workers = workers
	c.SelectBackup()
	c.lastHeartbeats = make([]time.Time, len(c.workers))
	// Initialize all last heartbeats to current time.
	for i := 0; i < len(c.workers); i++ {
		c.lastHeartbeats[i] = time.Now()
	}

	c.playerToRegionMap = make(map[string]int)
	c.regionToWorkerMap = make(map[int]int)

	// Assign main worker to each region.
	// TODO: handle cases where the number of regions != number of workers
	// TODO: perhaps randomize the assignment
	for i := 0; i < c.nRegions; i++ {
		c.regionToWorkerMap[i] = i

	}

	go c.run()
	return c
}
