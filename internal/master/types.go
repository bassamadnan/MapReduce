package m_utils

import (
	wpb "mapreduce/pkg/proto/worker"

	"google.golang.org/grpc"
)

// used for worker status
const (
	IDLE = iota
	WORKING
	FAIL
)

// used for worker status
const (
	MAPPER = iota
	REDUCER
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

// used for task type
const (
	PENDING  = iota // task is available for the worker
	ASSIGNED        // worker is working on this task
	COMPLETED
)

type Worker struct {
	Status       int
	PingResponse int
	AssignedTask int
	WorkerType   int                     //mapper or reducer
	Addr         string                  // localhost address of the worker machine, also servers as worker id
	Client       wpb.WorkerServiceClient // for the master to communicate with this worker
	Conn         *grpc.ClientConn
}

// job submitted to map reduce, to be split into tasks and assigned to workers
type Job struct {
	InputFileName string // data to be processed
	NumWorkers    int    // number of total workers (initially)
	Split         int    // number of splits
}

// tasks, do we need client here?
type Task struct {
	TaskID           int
	Start            int // start of the input file line no
	End              int // end of the input file line no
	TaskType         int
	TaskStatus       int
	Components       []int            // list of components assigned to mapper
	OutputPartitions map[string][]int // destination of all the partitions  output after its completed
}
