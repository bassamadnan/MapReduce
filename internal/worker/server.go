package w_utils

import (
	wpb "mapreduce/pkg/proto/worker"
	"sync"
)

// server side functions for grpc service of the worker
type Server struct {
	wpb.UnimplementedWorkerServiceServer
	ID int
	Mu sync.Mutex
}
