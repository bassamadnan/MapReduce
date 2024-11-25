package w_utils

import (
	"context"
	"fmt"
	"log"
	mpb "mapreduce/pkg/proto/master"
	wpb "mapreduce/pkg/proto/worker"
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

func (s *Server) ExecuteReduceTask(partition int, addr string, wg *sync.WaitGroup) {
	defer wg.Done()
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

	dataMap := make(map[string][]int)
	for _, kv := range resp.Kv {
		values := make([]int, len(kv.Value))
		for i, v := range kv.Value {
			values[i] = int(v)
		}
		dataMap[kv.Key] = values
	}

	for key, values := range dataMap {
		sum := 0
		for _, v := range values {
			sum += v
		}
		s.Mu.Lock()
		s.ReduceResults[key] += sum
		s.Mu.Unlock()
	}
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
