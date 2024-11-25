package w_utils

import (
	"context"
	"fmt"
	"log"
	utils "mapreduce/pkg"
	mpb "mapreduce/pkg/proto/master"
	wpb "mapreduce/pkg/proto/worker"
	"math"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func CompleteTask(client mpb.MasterServiceClient, worker_id string, task_id int, status bool, partitions []int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // timeout for ping
	defer cancel()
	partition32 := make([]int32, 0, len(partitions))
	for _, p := range partitions {
		partition32 = append(partition32, int32(p))
	}
	client.CompleteTask(ctx, &mpb.TaskStatus{
		TaskId:     int32(task_id),
		Status:     status,
		WorkerId:   worker_id,
		Partitions: partition32,
	})
	return nil
}

func PingReady(client mpb.MasterServiceClient, worker_id string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // timeout for ping
	defer cancel()
	numReducers, _ := client.Ready(ctx, &mpb.WorkerStatus{
		WorkerId: worker_id,
	})
	return int(numReducers.NumReducers), nil
}

func SendEdge(client mpb.MasterServiceClient, component int, edge utils.Edge) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // timeout for ping
	defer cancel()
	client.SendMinEdge(ctx, &mpb.EdgeInfo{
		Component: int32(component),
		U:         int32(edge.U),
		V:         int32(edge.V),
		W:         int32(edge.W),
	})
	return nil
}

func (s *Server) ExecuteReduceTask(partition int, addr string, wg *sync.WaitGroup) {
	defer wg.Done()

	// Setup connection
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	conn, err := grpc.NewClient(addr, opts...)
	if err != nil {
		log.Printf("Failed to connect to %s: %v", addr, err)
		return
	}
	defer conn.Close()

	client := wpb.NewWorkerServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get edges from mapper output
	resp, err := client.GetPartitionData(ctx, &wpb.Partition{Partition: int32(partition)})
	if err != nil {
		log.Printf("Error getting partition data from %s: %v", addr, err)
		return
	}

	// Process each component's edges
	for _, compEdges := range resp.CompEdges {
		comp := int(compEdges.ComponentId)
		minEdge := utils.Edge{W: math.MaxInt32}

		// Find minimum edge for this component
		for _, e := range compEdges.Edges {
			edge := utils.Edge{
				U: int(e.U),
				V: int(e.V),
				W: int(e.W),
			}

			if edge.W < minEdge.W {
				minEdge = edge
			}
		}

		if minEdge.W != math.MaxInt32 {
			s.Mu.Lock()
			if s.MinOutgoingEdges == nil {
				s.MinOutgoingEdges = make(map[int]utils.Edge)
			}
			if currEdge, exists := s.MinOutgoingEdges[comp]; !exists || minEdge.W < currEdge.W {
				s.MinOutgoingEdges[comp] = minEdge
			}
			s.Mu.Unlock()
		}
	}
	s.Mu.Lock()
	fmt.Printf("\nCurrent minimum edges after processing partition %d from %s:\n", partition, addr)
	for comp, edge := range s.MinOutgoingEdges {
		fmt.Printf("Component %d: %d -> %d (weight: %d)\n", comp, edge.U, edge.V, edge.W)
		SendEdge(s.WorkerMachineInstance.Client, comp, edge)
	}
	s.Mu.Unlock()
}
func (s *Server) StartReduceTask(partition int, addresses []string) {
	var wg sync.WaitGroup
	wg.Add(len(addresses))
	for _, addr := range addresses {
		go s.ExecuteReduceTask(partition, addr, &wg)
	}
	wg.Wait()
	fmt.Println("All reduce tasks completed. Results:")
	s.Mu.Lock()
	for key, sum := range s.ReduceResults {
		fmt.Printf("%s: %d\n", key, sum)
	}
	s.Mu.Unlock()
}
