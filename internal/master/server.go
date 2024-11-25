package m_utils

import (
	"context"
	"fmt"
	"log"

	utils "mapreduce/pkg"
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
	NumWorkers          int
	NumWorkersReady     int
	NumVertices         int
	NumMappers          int
	NumReducers         int
	NumReducersResponse int
	Mu                  sync.Mutex
	Tasks               []Task             // list of tasks, to be cleared after map phase
	ServiceRegistry     []string           // list of all worker machine address (indxes for the map blow)
	Workers             map[string]*Worker // worker[x] -> worker having id localhost:xxxx (xxxx : port number)
	AdjList             map[int][]utils.Edge
	MST                 []utils.Edge
	DSU                 *utils.DisjointSetUnion
	ComponentEdges      map[int][]utils.Edge // map[componentID][]Edge - stores all edges received for each component

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

func (s *Server) CompleteTask(ctx context.Context, req *mpb.TaskStatus) (*mpb.Empty, error) {
	fmt.Printf("Task %v completd by worker %v saved on partitions: %v\n", req.TaskId, req.WorkerId, req.Partitions)
	s.Mu.Lock()
	partitions := make([]int, 0, len(req.Partitions))
	for _, p := range req.Partitions {
		partitions = append(partitions, int(p))
	}
	workerID := GetWorkerID(req.WorkerId)
	if req.Status {
		s.Workers[workerID].Status = IDLE
		s.Workers[workerID].AssignedTask = -1
		s.Tasks[req.TaskId].TaskStatus = COMPLETED
		if s.Tasks[req.TaskId].OutputPartitions == nil {
			s.Tasks[req.TaskId].OutputPartitions = make(map[string][]int)
		}
		s.Tasks[req.TaskId].OutputPartitions[fmt.Sprintf("localhost:%v", req.WorkerId)] = partitions
	} else {
		s.Workers[workerID].Status = FAIL
		s.Workers[workerID].AssignedTask = -1
		s.Tasks[req.TaskId].TaskStatus = PENDING
	}
	s.Mu.Unlock()
	go s.AssignMapTasks() // assign remaining tasks
	return &mpb.Empty{}, nil
}

func (s *Server) Ready(ctx context.Context, req *mpb.WorkerStatus) (*mpb.ReadyResponse, error) {
	id := req.WorkerId
	fmt.Printf("Worker wiht %v ready! \n", id)
	s.Mu.Lock()
	s.NumWorkersReady += 1
	s.Mu.Unlock()

	return &mpb.ReadyResponse{NumReducers: int32(s.NumReducers)}, nil
}

func (s *Server) SendMinEdge(ctx context.Context, info *mpb.EdgeInfo) (*mpb.Empty, error) {
	comp := int(info.Component)
	edge := utils.Edge{U: int(info.U), V: int(info.V), W: int(info.W)}

	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.ComponentEdges[comp] = append(s.ComponentEdges[comp], edge)
	// fmt.Println("Current responses: ", s.NumReducersResponse)
	return &mpb.Empty{}, nil
}

func (s *Server) Complete(ctx context.Context, req *mpb.Empty) (*mpb.Empty, error) {
	s.Mu.Lock()
	s.NumReducersResponse++
	s.Mu.Unlock()
	return &mpb.Empty{}, nil
}
func (s *Server) ProcessMinEdges() {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	// Map to track already added edges (using both orientations)
	edgeSet := make(map[string]bool)
	for _, existingEdge := range s.MST {
		// Add both orientations to prevent duplicates
		edgeStr1 := fmt.Sprintf("%d-%d", existingEdge.U, existingEdge.V)
		edgeStr2 := fmt.Sprintf("%d-%d", existingEdge.V, existingEdge.U)
		edgeSet[edgeStr1] = true
		edgeSet[edgeStr2] = true
	}

	// Process new edges
	for comp, edges := range s.ComponentEdges {
		if len(edges) == 0 {
			continue
		}

		// Find minimum edge for this component
		minEdge := edges[0]
		for _, edge := range edges {
			if edge.W < minEdge.W {
				minEdge = edge
			}
		}

		// Check if we've already added this edge (in either orientation)
		edgeStr1 := fmt.Sprintf("%d-%d", minEdge.U, minEdge.V)
		edgeStr2 := fmt.Sprintf("%d-%d", minEdge.V, minEdge.U)

		if !edgeSet[edgeStr1] && !edgeSet[edgeStr2] {
			// Add edge to MST
			s.MST = append(s.MST, minEdge)
			// Mark edge as added
			edgeSet[edgeStr1] = true
			edgeSet[edgeStr2] = true
			// Merge components
			s.DSU.Union(minEdge.U, minEdge.V)
			fmt.Printf("Added MST edge for component %d: %d -> %d (weight: %d)\n",
				comp, minEdge.U, minEdge.V, minEdge.W)
		}
	}

	// Clear for next iteration
	s.ComponentEdges = make(map[int][]utils.Edge)

	// Print current state
	fmt.Println("\nCurrent MST edges:")
	for _, edge := range s.MST {
		fmt.Printf("%d -> %d (weight: %d)\n", edge.U, edge.V, edge.W)
	}
}
