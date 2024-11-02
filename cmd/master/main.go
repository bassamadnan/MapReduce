package main

import (
	"context"
	"fmt"
	"log"
	m_utils "mapreduce/internal/master"
	mpb "mapreduce/pkg/proto/master"
	"net"
	"sync"

	"google.golang.org/grpc"
)

/*
Users submit jobs to a scheduling system. Each job
consists of a set of tasks, and is mapped by the scheduler
to a set of available machines within a cluster

Assumptions:
1. On failure, tasks are re-executed by a different worker
2. Number of tasks depends on number of workers (the size of a task is equally divided among workers)
this does not matter really, since we will test with killing workers, just a starting point to divide tasks
*/
const (
	NUM_WORKERS       = 5
	BASE_WORKERS_ADDR = "localhost:7070"
)

type Server struct {
	mpb.UnimplementedMasterServiceServer
	num_workers int
	mu          sync.Mutex
}

// PingWorker(context.Context, *PingRequest) (*PingResponse, error)
func (s *Server) PingWorker(ctx context.Context, req *mpb.PingRequest) (*mpb.PingResponse, error) {
	fmt.Printf("ping from %v\n", req.ID)
	return &mpb.PingResponse{
		Response: true,
	}, nil
}

func main() {
	tasks, _ := m_utils.GetMapTasks(&m_utils.Job{
		InputFileName:   "data/input_1.txt",
		OutputDirectory: ".",
		NumWorkers:      NUM_WORKERS,
	})
	fmt.Printf(("%v\n"), tasks)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 5050))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	mpb.RegisterMasterServiceServer(s, &Server{
		num_workers: NUM_WORKERS,
	})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
