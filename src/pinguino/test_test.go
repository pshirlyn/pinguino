package pinguino

import (
	"log"
	"testing"
	"time"
)

func TestInitalizeNetwork(t *testing.T) {
	workers := 5
	regions := 5
	reliable := false

	cfg := make_config(t, workers, regions, reliable)
	defer cfg.cleanup()

	cfg.begin("TestInitalizeNetwork: Basic test")

	cfg.end()
}

func TestInitalizePlayer(t *testing.T) {
	workers := 5
	regions := 5
	reliable := false

	cfg := make_config(t, workers, regions, reliable)
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

func TestAddNewWorker(t *testing.T) {
	workers := 3
	regions := 3
	reliable := false

	cfg := make_config(t, workers, regions, reliable)
	defer cfg.cleanup()

	cfg.begin("TestAddNewWorker: New server worker is added")

	// Add a new fourth worker
	cfg.start1(4, cfg.applier)
	cfg.end()
}

func TestBasicSend(t *testing.T) {
	workers := 5
	regions := 5
	reliable := false

	cfg := make_config(t, workers, regions, reliable)
	defer cfg.cleanup()

	cfg.begin("TestBasicSend: Can send messages")

	player0 := cfg.startPlayer("player0")

	t0 := time.Now()
	for time.Since(t0).Seconds() < 10 {
		assigned := player0.isAssigned()

		if assigned {
			player0.MovePlayer(1, 1)

			time.Sleep(10 * time.Millisecond)

			if player0.state.x != 1 || player0.state.y != 1 {
				log.Println("Test basic send: wrong player state")
				log.Printf("Expected (%d, %d), got (%d, %d)\n", 1, 1, player0.state.x, player0.state.y)
			}
			cfg.end()
			return
		}

		time.Sleep(10 * time.Millisecond)
	}
}
