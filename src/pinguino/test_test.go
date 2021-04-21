package pinguino

import "testing"

func TestInitalizeNetwork(t *testing.T) {
	servers := 5
	reliable := false

	cfg := make_config(t, servers, reliable)
	defer cfg.cleanup()

	cfg.end()
}
