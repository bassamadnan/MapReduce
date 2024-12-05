package w_utils

import (
	"context"
	"fmt"
	utils "mapreduce/pkg"
	wpb "mapreduce/pkg/proto/worker"
)

func (w *WorkerMachine) execTask(components []int32, dsuParent []int32, taskID int) {
	// Convert int32 slices back to int for our utilities
	comps := make([]int, len(components))
	for i, c := range components {
		comps[i] = int(c)
	}

	dsu := &utils.DisjointSetUnion{
		Parent: make([]int, len(dsuParent)),
	}
	for i, p := range dsuParent {
		dsu.Parent[i] = int(p)
	}

	results := utils.GetComponentOutgoingEdges(comps, w.AdjList, dsu)
	fmt.Printf("Task %d Results:\n", taskID)
	for comp, edges := range results {
		fmt.Printf("Component %d outgoing edges: %v\n", comp, edges)
	}
	partitions, err := utils.WriteMapResults(results, w.OutputDirectory, taskID, w.NumReducers)
	if err != nil {
		fmt.Printf("Error writing map results: %v\n", err)
		return
	}
	fmt.Println("Partitions written:", partitions)
	// verify
	// for _, partition := range partitions {
	// 	dirPath := fmt.Sprintf("%s/%s", w.OutputDirectory, partition)
	// 	edges, err := utils.ReadDirectoryEdges(dirPath)
	// 	if err != nil {
	// 		fmt.Printf("Error reading partition %s: %v\n", partition, err)
	// 		continue
	// 	}
	// 	fmt.Printf("\nRead from partition %s:\n", partition)
	// 	for _, edge := range edges {
	// 		fmt.Printf("Edge: %d -> %d (weight: %d)\n", edge.U, edge.V, edge.W)
	// 	}
	// }
	CompleteTask(w.Client, w.ID, taskID, true, partitions)
}

// SendTask(context.Context, *TaskDescription) (*Empty, error)
func (s *Server) SendMapTask(ctx context.Context, task *wpb.MapTaskDescription) (*wpb.Empty, error) {
	fmt.Printf("Received task ID: %d\n", task.TaskID)
	fmt.Printf("Worker Components: %v\n", task.WorkerComponent)
	fmt.Printf("DSU length: %d\n", len(task.DSU))
	// go s.WorkerMachineInstance.execTask(int(start), int(end), int(taskid), s.InputFile)

	go s.WorkerMachineInstance.execTask(task.WorkerComponent, task.DSU, int(task.TaskID))

	return &wpb.Empty{}, nil
}

// master machine will ping the worker, this will be running on the go routine hosting the worker server
// the tasks executed by a worker will be running in the execTask go routine, so this will never be blocked by other operations
// i.e reply will be instant
func (s *Server) Ping(ctx context.Context, req *wpb.PingRequest) (*wpb.PingResponse, error) {
	return &wpb.PingResponse{
		Status: true,
	}, nil
}

func (s *Server) SendReduceTask(ctx context.Context, req *wpb.ReduceTaskDescription) (*wpb.Empty, error) {
	fmt.Printf("Recived partition:%v from addresses:%v\n", req.Partition, req.Addr)
	go s.StartReduceTask(int(req.Partition), req.Addr)
	return &wpb.Empty{}, nil
}

func (s *Server) GetPartitionData(ctx context.Context, req *wpb.Partition) (*wpb.Data, error) {
	directory := fmt.Sprintf("%v/%v", s.WorkerMachineInstance.OutputDirectory, req.Partition)
	compEdges, err := utils.ReadDirectoryEdges(directory)
	if err != nil {
		return nil, fmt.Errorf("error reading partition %d: %v", req.Partition, err)
	}

	protoCompEdges := make([]*wpb.ComponentEdges, 0, len(compEdges))
	for comp, edges := range compEdges {
		protoEdges := make([]*wpb.Edge, len(edges))
		for i, edge := range edges {
			protoEdges[i] = &wpb.Edge{
				U: int32(edge.U),
				V: int32(edge.V),
				W: int32(edge.W),
			}
		}
		protoCompEdges = append(protoCompEdges, &wpb.ComponentEdges{
			ComponentId: int32(comp),
			Edges:       protoEdges,
		})
	}

	return &wpb.Data{
		CompEdges: protoCompEdges,
	}, nil
}
