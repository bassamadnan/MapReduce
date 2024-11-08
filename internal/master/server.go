package m_utils

import (
	"context"
	"fmt"
	"log"

	mpb "mapreduce/pkg/proto/master"
	wpb "mapreduce/pkg/proto/worker"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type WorkerInfo struct {
	Addr string
	Role int
}

// server side functions for grpc service of the master
type Server struct {
	mpb.UnimplementedMasterServiceServer
	NumWorkers      int
	Mu              sync.Mutex
	Tasks           []Task             // list of tasks, to be cleared after map phase
	ServiceRegistry []string           // list of all worker machine address (indxes for the map blow)
	Workers         map[string]*Worker // worker[x] -> worker having id localhost:xxxx (xxxx : port number)
}

// function to initialize the worker list
// takes in a list of addresses , set's up a connection with each of them
// the id's are by default assumed to start from 0
func (s *Server) SetupWorkerClients(serviceRegistry []WorkerInfo) {
	s.Workers = make(map[string]*Worker)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	for _, worker := range serviceRegistry {
		conn, err := grpc.NewClient(worker.Addr, opts...)
		if err != nil {
			s.CloseAllConnections()
			log.Fatalf("conn failed %v", err)
		}
		s.Workers[worker.Addr] = &Worker{
			Status:       IDLE,
			Addr:         worker.Addr,
			WorkerType:   worker.Role, // 0 for map, 1 for reduce
			Client:       wpb.NewWorkerServiceClient(conn),
			Conn:         conn,
			AssignedTask: -1,
		}
	}

	addresses := make([]string, len(serviceRegistry))
	for i, worker := range serviceRegistry {
		addresses[i] = worker.Addr
	}
	s.ServiceRegistry = addresses
}

func (s *Server) CloseAllConnections() {
	for _, worker := range s.Workers {
		if worker != nil && worker.Conn != nil {
			worker.Conn.Close()
		}
	}
}

// CompleteTask(context.Context, *TaskStatus) (*Empty, error)

func (s *Server) CompleteTask(ctx context.Context, req *mpb.TaskStatus) (*mpb.Empty, error) {
	fmt.Printf("Task %v completd by worker %v saved on partitions: %v\n", req.TaskId, req.WorkerId, req.Partitions)
	s.Mu.Lock()
	workerID := GetWorkerID(req.WorkerId)
	if req.Status {
		s.Workers[workerID].Status = IDLE
		s.Workers[workerID].AssignedTask = -1
		s.Tasks[req.TaskId].TaskStatus = COMPLETED
	} else {
		s.Workers[workerID].Status = FAIL
		s.Workers[workerID].AssignedTask = -1
		s.Tasks[req.TaskId].TaskStatus = PENDING
	}
	s.Mu.Unlock()
	go s.AssignMapTasks() // assign remaining tasks
	return &mpb.Empty{}, nil
}
