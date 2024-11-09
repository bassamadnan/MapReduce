package w_utils

import (
	"context"
	"fmt"
	"log"
	c_utils "mapreduce/internal/common"
	mpb "mapreduce/pkg/proto/master"
	wpb "mapreduce/pkg/proto/worker"
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

func ExecuteReduceTask(partition int, addr string) {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	conn, err := grpc.NewClient(addr, opts...)
	if err != nil {
		log.Fatalf("conn failed %v", err)
	}
	defer conn.Close()

	client := wpb.NewWorkerServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.GetPartitionData(ctx, &wpb.Partition{Partition: int32(partition)})
	if err != nil {
		log.Fatalf("Error in get partition data client %v\n", err)
	}

	var kvList []c_utils.KeyValue
	for _, kv := range resp.Kv {
		kvList = append(kvList, c_utils.KeyValue{
			Key:   kv.Key,
			Value: int(kv.Value),
		})
	}

	reducedData := c_utils.Reduce(kvList)
	for _, kv := range reducedData {
		fmt.Printf("%s %d\n", kv.Key, kv.Value)
	}
}
