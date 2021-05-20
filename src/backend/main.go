package pinguino

type Framework struct {
	cfg         *config
	players     []*Player
	playerCount int // count of players that have already been handed to client
}

func Setup() (*Framework, *Player) {
	workers := 5
	regions := 5
	reliable := false

	fw := &Framework{}
	cfg := make_config(nil, workers, regions, reliable)
	fw.cfg = cfg

	// Create 3 players
	fw.players = append(fw.players, cfg.startPlayer("player0"))
	fw.players = append(fw.players, cfg.startPlayer("player1"))
	fw.players = append(fw.players, cfg.startPlayer("player2"))

	fw.playerCount = 1

	return fw, fw.players[0]
}

func (fr *Framework) GetNextPlayer() *Player {
	nextPlayer := fr.players[fr.playerCount]
	fr.playerCount += 1
	return nextPlayer
}
