package pinguino

import (
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

	player0 := cfg.startPlayer("player0")

	// Check if the player was assigned to a region successfully within 10 secs
	t0 := time.Now()
	for time.Since(t0).Seconds() < 10 {
		assigned := player0.isAssigned()

		if assigned {
			cfg.end()
			return
		}

		time.Sleep(10 * time.Millisecond)
	}

	t.Fatalf("Player not assigned within 10 sec")
	cfg.end()

}

func TestBasicSend(t *testing.T) {
	servers := 5
	regions := 5
	reliable := false

	cfg := make_config(t, servers, regions, reliable)
	defer cfg.cleanup()

	cfg.begin("TestBasicSend: Can send messages")

	player0 := cfg.startPlayer("player0")
	player0.Move(1, 1)

	cfg.end()
}
