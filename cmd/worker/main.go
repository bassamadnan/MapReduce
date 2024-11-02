package main

import (
	"context"
	"fmt"
	"log"
	mpb "mapreduce/pkg/proto/master"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	Client mpb.MasterServiceClient
	ID     string
}

func (c *Client) PingMaster() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := c.Client.PingWorker(ctx, &mpb.PingRequest{ID: c.ID})
	if err != nil {
		log.Fatalf("error in getting servers %v\n", err)
	}
	fmt.Printf("status: %v\n", resp.Response)
}
func main() {
	BASE_SERVER_ADDR := "localhost:5050"
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	conn, err := grpc.NewClient(BASE_SERVER_ADDR, opts...)
	if err != nil {
		log.Fatalf("conn failed %v", err)
	}
	defer conn.Close()
	client := Client{
		Client: mpb.NewMasterServiceClient(conn),
		ID:     "4",
	}
	client.PingMaster()

}
