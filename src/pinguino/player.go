package pinguino

import (
	"pinguino/src/labrpc"
	"sync"
)

type Player struct {
	mu sync.Mutex

	coordinator *labrpc.ClientEnd
	servers     []*labrpc.ClientEnd
	username    string

	region      int // region connected to
	serverIndex int
}

func (pl *Player) Kill() {
}

func (pl *Player) SetWorkers(servers []*labrpc.ClientEnd, region int, serverIndex int) {
	// In case workers change, we update our list of worker ClientEnds
	pl.mu.Lock()
	pl.servers = servers
	pl.region = region
	pl.serverIndex = serverIndex
	pl.mu.Unlock()
}

func MakePlayer(coordinator *labrpc.ClientEnd, servers []*labrpc.ClientEnd, username string) *Player {
	pl := &Player{}

	pl.mu.Lock()
	defer pl.mu.Unlock()

	pl.coordinator = coordinator
	pl.servers = servers
	pl.username = username

	return pl
}
