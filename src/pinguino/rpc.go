package pinguino

import (
	"os"
	"strconv"
)

type StableMoveArgs struct {
	Move interface{}
}

type StableMoveReply struct {
}

type FastMoveArgs struct {
	Move interface{}
}

type FastMoveReply struct {
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
