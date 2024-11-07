package w_utils

import (
	"context"
	mpb "mapreduce/pkg/proto/master"
	"time"
)

func CompleteTask(client mpb.MasterServiceClient, worker_id string, task_id int, status bool, output_path string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // timeout for ping
	defer cancel()
	client.CompleteTask(ctx, &mpb.TaskStatus{
		TaskId:     int32(task_id),
		Status:     status,
		WorkerId:   worker_id,
		OutputPath: output_path,
	})
	return nil
}
