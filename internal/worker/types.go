package w_utils

import (
	utils "mapreduce/pkg"
	mpb "mapreduce/pkg/proto/master"
	wpb "mapreduce/pkg/proto/worker"
	"sync"

	"google.golang.org/grpc"
)

// server side functions for grpc service of the worker
type Server struct {
	wpb.UnimplementedWorkerServiceServer
	Mu                    sync.Mutex
	InputFile             string // should this be here?
	WorkerMachineInstance *WorkerMachine
	ReduceResults         map[string]int
}

// information about the worker
type WorkerMachine struct {
	MasterAddr      string
	NumVertices     int
	WorkerType      int                     // mapper/reducer
	ID              string                  // port number
	OutputDirectory string                  // store results of task
	Client          mpb.MasterServiceClient // to communicate with master machine
	Conn            *grpc.ClientConn
	ServerInstance  *Server
	AdjList         map[int][]utils.Edge
	MST             []utils.Edge
	DSU             *utils.DisjointSetUnion
}
