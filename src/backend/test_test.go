package backend

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func checkPlayerAssigned(pl *Player) bool {
	// Check if the player was assigned to a region successfully within 10 secs
	t0 := time.Now()
	for time.Since(t0).Seconds() < 10 {
		assigned := pl.isAssigned()

		if assigned {
			return true
		}

		time.Sleep(10 * time.Millisecond)
	}

	log.Fatalf("Player not assigned within 10 sec")
	return false
}

// ---------------------------------------------
// 				   	Test Suite
// ---------------------------------------------

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
	checkPlayerAssigned(player0)
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

func TestBasicDisconnectWorker(t *testing.T) {
	workers := 3
	regions := 3
	reliable := false

	cfg := make_config(t, workers, regions, reliable)
	defer cfg.cleanup()

	cfg.begin("TestBasicDisconnectWorker: Disconnect a worker")

	// Disconnect server worker 1
	cfg.disconnect(1)

	time.Sleep(100 * time.Millisecond)
	cfg.end()
}

func TestDisconnectWorkerAndReassignPlayer(t *testing.T) {
	workers := 3
	regions := 3
	reliable := false

	cfg := make_config(t, workers, regions, reliable)
	defer cfg.cleanup()

	cfg.begin("TestDisconnectWorkerAndReassignPlayer: Disconnect a worker and reassign its players to another worker")

	player0 := cfg.startPlayer("player0")

	// Ensure that the player had enough time to initialize and be assigned to a region and an corresponding worker
	checkPlayerAssigned(player0)

	// Add 1 because server index uses zero indexing, but we need one indexing to connect/disconnect with config
	assignedWorker := player0.getServerIndex() + 1

	fmt.Printf("Disconnecting worker %d\n", assignedWorker)
	cfg.disconnect(assignedWorker)

	// Check if the player was reassigned to another worker
	t0 := time.Now()
	for time.Since(t0).Seconds() < 10 {
		newAssignedWorker := player0.getServerIndex() + 1
		if assignedWorker != newAssignedWorker {
			cfg.end()
			return
		}

		time.Sleep(10 * time.Millisecond)
	}

	log.Fatalf("Player was not reassigned within 10 secs")
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
