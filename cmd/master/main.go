package main

import (
	"fmt"
	"log"
	m_utils "mapreduce/internal/master"
	mpb "mapreduce/pkg/proto/master"
	wpb "mapreduce/pkg/proto/worker"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

/*
Users submit jobs to a scheduling system. Each job
consists of a set of tasks, and is mapped by the scheduler
to a set of available machines within a cluster

Assumptions:
1. On failure, tasks are re-executed by a different worker
2. Number of tasks depends is predefined (DEFAULT_SPLIT) based on number of lines
3. The master does not have to send the input file to the workers, just the line numbers (can be changed later if needed)
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
	s := grpc.NewServer()
	mpb.RegisterMasterServiceServer(s, &m_utils.Server{
		NumWorkers: NUM_WORKERS,
	})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	// create map tasks
	tasks, _ := m_utils.GetMapTasks(&m_utils.Job{
		InputFileName: "data/input_1.txt",
		NumWorkers:    NUM_WORKERS,
		Split:         DEFAULT_SPLIT,
	})
	fmt.Printf(("%v\n"), tasks)
	// connect to all worker machines
	serviceRegistry := []string{"localhost:7070", "localhost:7071", "localhost:7072", "localhost:7073"}
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	conn, err := grpc.NewClient(BASE_SERVER_ADDR, opts...)
	if err != nil {
		log.Fatalf("conn failed %v", err)
	}
	defer conn.Close()
	// crdt.IsCRDT(&doc)
	client = m_utils.MasterClient{
		Client: wpb.NewWorkerServiceClient(conn),
	}

}
