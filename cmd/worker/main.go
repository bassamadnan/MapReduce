package main

import (
	"flag"
	"fmt"
	"log"
	w_utils "mapreduce/internal/worker"
	mpb "mapreduce/pkg/proto/master"
	wpb "mapreduce/pkg/proto/worker"
	"net"

	"google.golang.org/grpc"
)

type Client struct {
	Client mpb.MasterServiceClient
	ID     string
}

func main() {
	port := flag.Int("port", 7070, "")
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	wpb.RegisterWorkerServiceServer(s, &w_utils.Server{
		ID: *port,
	})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
