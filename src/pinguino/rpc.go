package pinguino

import (
	"os"
	"strconv"
)

type StableMoveArgs struct {
	Command MoveCommand
}

type StableMoveReply struct {
	Success bool
}

type FastMoveArgs struct {
	Command MoveCommand
}

type FastMoveReply struct {
	Success bool
}

type HeartbeatArgs struct {
	AddedPlayers []string
}

type HeartbeatReply struct {
	// report players who left
	Success             bool
	DisconnectedPlayers []string
}

// Cook up a unique-ish UNIX-domain socket name
// in /var/tmp, for the coordinator.
// Can't use the current directory since
// Athena AFS doesn't support UNIX-domain sockets.
func coordinatorSock() string {
	s := "/var/tmp/pinguino-"
	s += strconv.Itoa(os.Getuid())
	return s
}

type AssignPlayerToRegionArgs struct {
	Username string
}

type AssignPlayerToRegionReply struct {
	Success bool
	Region  int
	Worker  int
}
