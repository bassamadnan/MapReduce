package w_utils

import (
	"context"
	"fmt"
	wpb "mapreduce/pkg/proto/worker"
	"sync"
)

// server side functions for grpc service of the worker
type Server struct {
	wpb.UnimplementedWorkerServiceServer
	ID int
	Mu sync.Mutex
}

// SendTask(context.Context, *TaskDescription) (*Empty, error)
func (s *Server) SendTask(ctx context.Context, task *wpb.TaskDescription) (*wpb.Empty, error) {
	start, end := task.Start, task.End
	fmt.Printf("Worker machine: %d Recieved task start: %d, %d\n", s.ID, start, end)
	return &wpb.Empty{}, nil
}
