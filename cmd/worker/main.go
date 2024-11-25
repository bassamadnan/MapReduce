package main

import (
	"flag"
	"fmt"
	"log"
	w_utils "mapreduce/internal/worker"
	mpb "mapreduce/pkg/proto/master"
	wpb "mapreduce/pkg/proto/worker"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"google.golang.org/grpc"
)

type Client struct {
	Client          mpb.MasterServiceClient
	ID              string
	OutputDirectory string
}

func main() {
	port := flag.Int("port", 7070, "")
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	w := w_utils.WorkerMachine{}
	w.Initialize("localhost:5050", strconv.Itoa(*port), "data/input/19.mtx")
	s := grpc.NewServer()
	wpb.RegisterWorkerServiceServer(s, w.ServerInstance)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// exit program
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	s.GracefulStop()
}
