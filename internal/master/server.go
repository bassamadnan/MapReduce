package m_utils

import (
	"log"
	mpb "mapreduce/pkg/proto/master"
	wpb "mapreduce/pkg/proto/worker"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// server side functions for grpc service of the master
type Server struct {
	mpb.UnimplementedMasterServiceServer
	NumWorkers      int
	Mu              sync.Mutex
	ServiceRegistry []string // list of all worker machine address
	Workers         []*Worker
}

// function to initialize the worker list
// takes in a list of addresses , set's up a connection with each of them
// the id's are by default assumed to start from 0
func (s *Server) SetupWorkerClients(serviceRegistry []string) {
	s.Workers = make([]*Worker, len(serviceRegistry))
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	for i, addr := range serviceRegistry {

		conn, err := grpc.NewClient(serviceRegistry[i], opts...)
		if err != nil {
			s.CloseAllConnections()
			log.Fatalf("conn failed %v", err)
		}
		s.Workers[i] = &Worker{
			Status: IDLE,
			Addr:   addr,
			Client: wpb.NewWorkerServiceClient(conn),
			Conn:   conn,
		}
	}
	s.ServiceRegistry = serviceRegistry
}

func (s *Server) CloseAllConnections() {
	for _, worker := range s.Workers {
		if worker != nil && worker.Conn != nil {
			worker.Conn.Close()
		}
	}
}
