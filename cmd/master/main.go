package main

import (
	"fmt"
	"log"
	m_utils "mapreduce/internal/master"
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

func main() {
	// setup server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 5050))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	masterServer := &m_utils.Server{
		NumWorkers: NUM_WORKERS,
	}
	s := grpc.NewServer()
	mpb.RegisterMasterServiceServer(s, masterServer)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	// create map tasks
	tasks, _ := m_utils.GetMapTasks(&m_utils.Job{
		InputFileName: "data/input/input_1.txt",
		NumWorkers:    NUM_WORKERS,
		Split:         DEFAULT_SPLIT,
	})
	fmt.Printf(("%v\n"), tasks)
	// connect to all worker machines
	serviceRegistry := []string{"localhost:7070", "localhost:7071", "localhost:7072", "localhost:7073"}
	masterServer.SetupWorkerClients(serviceRegistry)
	go masterServer.StartPing() // start pinging the machines periodically on background
	masterServer.AssignTasks(tasks)

	// exit program
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	s.GracefulStop()

}
