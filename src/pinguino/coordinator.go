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
	"sync/atomic"
	"time"
)

const heartbeatTimeInterval = time.Duration(10 * time.Millisecond)

type Coordinator struct {
	mu sync.Mutex

	nRegions          int
	players           map[string]*labrpc.ClientEnd
	workers           []*labrpc.ClientEnd
	workerReplicas    map[int][]int // index by worker, list of replicas for workers
	isConnectedWorker map[int]bool

	lastHeartbeats    []time.Time
	playerToRegionMap map[string]int
	regionToWorkerMap map[int]int
	backup            int

	dead int32
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

	c.mu.Lock()
	defer c.mu.Unlock()

	if !ok && !c.killed() {
		if c.isConnectedWorker[workerIndex] {
			// Add one because config and test references 1 indexing for workers
			log.Printf("couldn't reach worker %d\n", workerIndex+1)

			c.WorkerCrashed(workerIndex)
		}
		return
	}

	// TODO: handle worker that reconnects. need to redistribute region

	// Successfully sent heartbeat to worker, so update lastHearbeat
	c.lastHeartbeats[workerIndex] = time.Now()
	c.isConnectedWorker[workerIndex] = true
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

func (c *Coordinator) NewPlayerAdded(username string, end *labrpc.ClientEnd) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.players[username] = end
}

func (c *Coordinator) NewWorkersAdded(workers []*labrpc.ClientEnd) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// TODO: rather than just adding a new region, check if there's a worker that is handling multiple regions first
	for i := 0; i < len(workers); i++ {
		c.regionToWorkerMap[c.nRegions+i] = c.nRegions + 1
		c.lastHeartbeats = append(c.lastHeartbeats, time.Now())
		c.isConnectedWorker[len(c.workers)+i] = true
	}
	c.nRegions += len(workers)

	c.workers = append(c.workers, workers...)
}

// Generate a random connected worker
func (c *Coordinator) randomConnectedWorker() int {
	for i := 0; i < len(c.workers); i++ {
		w := rand.Intn(len(c.workers))
		if c.isConnectedWorker[w] {
			return w
		}
	}
	log.Println("Could not get a random connected worker")
	return -1
}

func (c *Coordinator) WorkerCrashed(i int) {
	c.isConnectedWorker[i] = false

	regionsReassigned := make(map[int]bool, c.nRegions)
	// Reassign regions that were assigned to worker i
	for region, worker := range c.regionToWorkerMap {
		if worker == i {
			newWorker := c.randomConnectedWorker()
			c.regionToWorkerMap[region] = newWorker
			regionsReassigned[region] = true
		}
	}

	// Reassign players that were assigned to reassigned regions
	for username, region := range c.playerToRegionMap {
		if regionsReassigned[region] {
			newWorker := c.regionToWorkerMap[region]
			args := WorkerReassignmentArgs{Worker: newWorker}
			reply := WorkerReassignmentReply{}
			go c.players[username].Call("Player.WorkerReassignment", &args, &reply)
		}
	}
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
	atomic.StoreInt32(&c.dead, 1)
}

func (c *Coordinator) killed() bool {
	z := atomic.LoadInt32(&c.dead)
	return z == 1
}

func (c *Coordinator) run() {
	// main loop
	for !c.killed() {
		c.maybeSendHeartbeats()
		time.Sleep(heartbeatTimeInterval)
	}
}

func (c *Coordinator) SelectBackup() {
	// Selects backup out of existing servers
	idx := rand.Intn(len(c.workers))
	c.backup = idx
}

func (c *Coordinator) SelectWorkerReplicas() {
	connectedWorkers := make([]int, 0)
	replicas := make(map[int][]int)
	for i, connected := range c.isConnectedWorker {
		if connected {
			connectedWorkers = append(connectedWorkers, i)
		}
	}

	for i, worker := range connectedWorkers {
		replicas[worker] = make([]int, 0)
		right := (i + 1) % len(connectedWorkers)
		left := (i + len(connectedWorkers) - 1) % len(connectedWorkers)
		replicas[worker] = append(replicas[worker], right, left)

		args := SetReplicasArgs{replicas[worker]}
		reply := SetReplicasReply{}
		c.workers[worker].Call("Worker.Heartbeat", args, reply)
	}

	c.workerReplicas = replicas
}

func MakeCoordinator(workers []*labrpc.ClientEnd, regions int) *Coordinator {
	c := &Coordinator{}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.nRegions = regions
	c.backup = -1

	c.players = make(map[string]*labrpc.ClientEnd)
	c.workers = workers
	c.SelectBackup()
	c.SelectWorkerReplicas()
	c.lastHeartbeats = make([]time.Time, len(c.workers))
	c.isConnectedWorker = make(map[int]bool, len(c.workers))

	// Initialize all workers to be connected and set last heartbeats to current time.
	for i := 0; i < len(c.workers); i++ {
		c.isConnectedWorker[i] = true
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
