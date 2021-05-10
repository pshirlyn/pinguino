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

func TestInitalizeNetworkAndWait(t *testing.T) {
	workers := 5
	regions := 5
	reliable := false

	cfg := make_config(t, workers, regions, reliable)
	defer cfg.cleanup()

	cfg.begin("TestInitalizeNetworkAndWait: After initializing, wait for a bit to see if heartbeats are ok")

	time.Sleep(3 * time.Second)

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
	cfg.start1(4, cfg.applier, false)
	cfg.end()
}

func TestAddMultipleNewWorkers(t *testing.T) {
	workers := 3
	regions := 3
	reliable := false

	cfg := make_config(t, workers, regions, reliable)
	defer cfg.cleanup()

	cfg.begin("TestAddMultipleNewWorkers: Multiple new server workers are added")

	// Add a new fourth worker
	cfg.start1(4, cfg.applier, false)

	// Add a new fifth worker
	cfg.start1(5, cfg.applier, false)
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
			player0.ClientMovePlayer(1, 1)

			time.Sleep(10 * time.Millisecond)

			if player0.state.x != 1 || player0.state.y != 1 {
				log.Println("Test basic send: wrong player state after 10ms")
				log.Printf("Expected (%d, %d), got (%d, %d)\n", 1, 1, player0.state.x, player0.state.y)
			}
			cfg.end()
			return
		}

		time.Sleep(10 * time.Millisecond)
	}
}

func TestStableSend(t *testing.T) {
	servers := 5
	regions := 5
	reliable := false

	cfg := make_config(t, servers, regions, reliable)
	defer cfg.cleanup()

	cfg.begin("TestStableSend: Can send messages")

	player0 := cfg.startPlayer("player0")

	t0 := time.Now()
	for time.Since(t0).Seconds() < 10 {
		assigned := player0.isAssigned()

		if assigned {
			player0.SendChatMessage("hello")

			time.Sleep(10 * time.Millisecond)

			// how to confirm that message was added to game??

			// msg := ChatMessage{"hello", "player0"}
			// if player0.servers[player0.serverIndex].logs[0] != msg {
			// 	log.Println("Test stable send: wrong player state after 10ms")
			// 	log.Printf("Expected (%d, %d), got (%d, %d)\n", 1, 1, player0.state.x, player0.state.y)
			// }
			cfg.end()
			return
		}

		time.Sleep(10 * time.Millisecond)
	}
}

func TestCoordinatorBackupAssigned(t *testing.T) {
	servers := 5
	regions := 5
	reliable := false

	cfg := make_config(t, servers, regions, reliable)
	defer cfg.cleanup()

	cfg.begin("TestCoordinatorBackupAssigned: Has backup")

	if cfg.coordinator.backup == -1 {
		log.Fatalf("Coordinator has no backup assigned")
	}

	// player0 := cfg.startPlayer("player0")

	// t0 := time.Now()
	// for time.Since(t0).Seconds() < 10 {
	// 	assigned := player0.isAssigned()

	// 	if assigned {
	// 		player0.SendChatMessage("hello")

	// 		time.Sleep(10 * time.Millisecond)

	// 		// how to confirm that message was added to game??

	// 		// msg := ChatMessage{"hello", "player0"}
	// 		// if player0.servers[player0.serverIndex].logs[0] != msg {
	// 		// 	log.Println("Test stable send: wrong player state after 10ms")
	// 		// 	log.Printf("Expected (%d, %d), got (%d, %d)\n", 1, 1, player0.state.x, player0.state.y)
	// 		// }
	// 		cfg.end()
	// 		return
	// 	}

	// 	time.Sleep(10 * time.Millisecond)
	// }
}
