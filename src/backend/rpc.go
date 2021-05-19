package backend

import (
	"os"
	"strconv"
)

type SendReplicaArgs struct {
	Replica   []byte
	Move      MoveCommand
	Worker    int
	MoveIndex int
}

type SendReplicaReply struct {
	Success bool
}

type SetReplicasArgs struct {
	Replicas []int
}

type SetReplicasReply struct {
	Success bool
}

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

type WorkerReassignmentArgs struct {
	Worker int
}

type WorkerReassignmentReply struct {
	Success bool
}

type AssignPlayerToRegionArgs struct {
	Username string
}

type AssignPlayerToRegionReply struct {
	Success bool
	Region  int
	Worker  int
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
