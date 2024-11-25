package main

import (
	"fmt"
	"log"
	m_utils "mapreduce/internal/master"
	utils "mapreduce/pkg"
	mpb "mapreduce/pkg/proto/master"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

/*
Users submit jobs to a scheduling system. Each job
consists of a set of tasks, and is mapped by the scheduler
to a set of available machines within a cluster

Assumptions:
1. On failure, tasks are re-executed by a different worker
2. Number of tasks depends is predefined (DEFAULT_SPLIT) based on number of lines
3. The master does not have to send the input file to the workers, just the line numbers (can be changed later if needed)
4. We dont expect worker machines to come back up, and no new worker machines will be added during execution
*/
const (
	NUM_WORKERS       = 4
	BASE_WORKERS_ADDR = "localhost:7070"
	DEFAULT_SPLIT     = 7
)

func GetUniqueComponents(dsu *utils.DisjointSetUnion) int {
	componentSet := make(map[int]bool)
	for i := range dsu.Parent {
		componentSet[dsu.Find(i)] = true
	}
	return len(componentSet)
}

func main() {
	adjList, _ := utils.ReadMTXFile("data/input/sample.mtx")
	utils.PrintAdjList(adjList)
	NumMappers, NumReducers := 2, 2
	// setup server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 5050))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	masterServer := &m_utils.Server{
		NumWorkers:     NUM_WORKERS,
		NumMappers:     NumMappers,
		NumReducers:    NumReducers,
		NumVertices:    len(adjList),
		DSU:            utils.NewDSU(len(adjList)),
		AdjList:        adjList,
		MST:            make([]utils.Edge, 0),
		ComponentEdges: make(map[int][]utils.Edge),
	}
	fmt.Println(masterServer.DSU.Parent)
	s := grpc.NewServer()
	mpb.RegisterMasterServiceServer(s, masterServer)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	WorkerComponents := utils.GetWorkerComponents(masterServer.DSU, masterServer.NumMappers)
	for i, comp := range WorkerComponents {
		fmt.Printf("Worker %d components: %v\n", i, comp)
	}
	// create map tasks, done by a scheduler in the paper, out of scope for this project
	tasks, _ := m_utils.GetMapTasks(WorkerComponents)
	fmt.Printf(("%v\n"), tasks)
	for {
		if masterServer.NumWorkersReady == NUM_WORKERS {
			break
		}
	}
	fmt.Print("All workers ready!\n")
	// connect to all worker machines

	serviceRegistry := []m_utils.WorkerInfo{
		{Addr: "localhost:7070", Role: 0}, // map worker
		{Addr: "localhost:7071", Role: 0}, // map worker
		{Addr: "localhost:7072", Role: 1}, // reduce worker
		{Addr: "localhost:7073", Role: 1}, // reduce worker
	}
	masterServer.SetupWorkerClients(serviceRegistry)
	go masterServer.StartPing() // start pinging the machines periodically on background
	for {
		numComponents := GetUniqueComponents(masterServer.DSU)
		if numComponents <= 1 {
			sum := 0
			fmt.Println("MST construction complete!")
			fmt.Println("\nFinal MST edges:")

			for _, edge := range masterServer.MST {
				fmt.Printf("%d -> %d (weight: %d)\n", edge.U, edge.V, edge.W)
				sum += edge.W
			}
			fmt.Printf("Final sum: %v\n", sum)
			break
		}
		components := utils.GetComponents(masterServer.DSU)

		fmt.Printf("\nStarting new iteration with %d components\n", len(components))

		// Divide components among mappers
		WorkerComponents := utils.GetWorkerComponents(masterServer.DSU, masterServer.NumMappers)
		for i, comp := range WorkerComponents {
			fmt.Printf("Worker %d components: %v\n", i, comp)
		}

		// Create and assign map tasks
		tasks, _ := m_utils.GetMapTasks(WorkerComponents)
		masterServer.Tasks = tasks
		masterServer.AssignMapTasks()

		// Wait for map phase completion
		for {
			count := 0
			masterServer.Mu.Lock()
			for _, task := range masterServer.Tasks {
				if task.TaskStatus == m_utils.COMPLETED {
					count++
				}
			}
			masterServer.Mu.Unlock()
			if count == len(masterServer.Tasks) {
				break
			}
		}

		// Start reduce phase
		masterServer.NumReducersResponse = 0                     // Reset counter
		masterServer.ComponentEdges = make(map[int][]utils.Edge) // Clear previous edges
		masterServer.AssignReducerTasks()

		fmt.Println("WAITING FOR REDUCERS")
		for {
			masterServer.Mu.Lock()
			if masterServer.NumReducersResponse == masterServer.NumReducers {
				masterServer.Mu.Unlock()
				break
			}
			masterServer.Mu.Unlock()
		}
		fmt.Println(" REDUCERS WAIT OVER")
		// Process edges and update MST
		masterServer.ProcessMinEdges()

		fmt.Println("\nCurrent MST edges:")
		for _, edge := range masterServer.MST {
			fmt.Printf("%d -> %d (weight: %d)\n", edge.U, edge.V, edge.W)
		}
	}
	// exit program
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	s.GracefulStop()

}
