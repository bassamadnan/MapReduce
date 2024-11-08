package w_utils

import (
	"context"
	mpb "mapreduce/pkg/proto/master"
	"time"
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
