package m_utils

import (
	"context"
	"fmt"
	mpb "mapreduce/pkg/proto/master"
	"sync"
)

// server side functions for grpc service of the master
type Server struct {
	mpb.UnimplementedMasterServiceServer
	NumWorkers      int
	Mu              sync.Mutex
	ServiceRegistry []string
}

// PingWorker(context.Context, *PingRequest) (*PingResponse, error)
func (s *Server) PingWorker(ctx context.Context, req *mpb.PingRequest) (*mpb.PingResponse, error) {
	fmt.Printf("ping from %v\n", req.ID)
	return &mpb.PingResponse{
		Response: true,
	}, nil
}
