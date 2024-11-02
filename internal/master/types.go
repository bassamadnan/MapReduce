package m_utils

import "time"

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
	Status       int
	LastPingTime time.Time
	PingResponse int
	AssignedTask int
}

// job submitted to map reduce, to be split into tasks and assigned to workers
type Job struct {
	InputFileName   string // data to be processed
	OutputDirectory string // map output before reduce phase, to be stored on local disk of the worker
	NumWorkers      int    // number of total workers (initially)
}

// tasks
type Task struct {
	TaskID   int
	Start    int // strat of the input file line no
	End      int // end of the input file line no
	TaskType int
}
