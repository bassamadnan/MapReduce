package main

import (
	"context"
	"fmt"
	"log"
	mpb "mapreduce/pkg/proto/master"
	"net"
	"sync"

	"google.golang.org/grpc"
)

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
