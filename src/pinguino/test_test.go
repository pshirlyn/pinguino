package pinguino

import (
	"pinguino/src/labrpc"
	"testing"
	"time"
)

func TestInitalizeNetwork(t *testing.T) {
	servers := 5
	regions := 5
	reliable := false

	cfg := make_config(t, servers, regions, reliable)
	defer cfg.cleanup()

	cfg.begin("TestInitalizeNetwork: Basic test")

	cfg.end()
}

func TestInitalizePlayer(t *testing.T) {
	servers := 5
	regions := 5
	reliable := false

	cfg := make_config(t, servers, regions, reliable)
	defer cfg.cleanup()

	cfg.begin("TestInitalizePlayer: Player initialization assigned to region")

	endnames := make([]string, cfg.nservers)
	for j := 0; j < cfg.nservers; j++ {
		endnames[j] = randstring(20)
	}

	// a fresh set of ClientEnds.
	ends := make([]*labrpc.ClientEnd, cfg.nservers)
	for j := 0; j < cfg.nservers; j++ {
		ends[j] = cfg.net.MakeEnd(endnames[j])
		cfg.net.Connect(endnames[j], j)
		cfg.net.Enable(endnames[j], true)
	}

	coordinatorEnd := ends[0]
	workerEnds := ends[1:]

	player0 := MakePlayer(coordinatorEnd, workerEnds, "player0")

	t0 := time.Now()
	assigned := false
	for time.Since(t0).Seconds() < 10 {
		assigned = player0.isAssigned()
		if assigned {
			break
		}

		time.Sleep(10 * time.Millisecond)
	}

	if !assigned {
		t.Fatalf("Player not assigned within 10 sec")
	}

	cfg.end()

}
