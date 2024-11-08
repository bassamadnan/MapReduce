package w_utils

import (
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
}

// information about the worker
type WorkerMachine struct {
	MasterAddr      string
	WorkerType      int                     // mapper/reducer
	ID              string                  // port number
	OutputDirectory string                  // store results of task
	Client          mpb.MasterServiceClient // to communicate with master machine
	Conn            *grpc.ClientConn
	ServerInstance  *Server
}
