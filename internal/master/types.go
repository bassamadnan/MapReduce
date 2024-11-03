package m_utils

import (
	wpb "mapreduce/pkg/proto/worker"
	"time"

	"google.golang.org/grpc"
)

// used for worker status
const (
	IDLE = iota
	WORKING
	FAIL
)

// used for ping
const (
	NIL     = iota // didnt ping yet
	SUCCESS        // sent back a response to the ping
	WAITING        // awaiting ping response
)

// used for task type
const (
	MAP_TASK = iota
	REDUCE_TASK
)

type Worker struct {
	Status          int
	LastPingTime    time.Time
	PingResponse    int
	AssignedTask    int
	OutputDirectory string                  // workers local disk directory
	Addr            string                  // localhost address of the worker machine
	Client          wpb.WorkerServiceClient // for the master to communicate with this worker
	Conn            *grpc.ClientConn
}

// job submitted to map reduce, to be split into tasks and assigned to workers
type Job struct {
	InputFileName string // data to be processed
	NumWorkers    int    // number of total workers (initially)
	Split         int    // number of splits
}

// tasks
type Task struct {
	TaskID   int
	Start    int // strat of the input file line no
	End      int // end of the input file line no
	TaskType int
}
