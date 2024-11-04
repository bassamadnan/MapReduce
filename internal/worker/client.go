package w_utils

import (
	"context"
	"fmt"
	mpb "mapreduce/pkg/proto/master"
	"time"
)

func CompleteTask(client mpb.MasterServiceClient, worker_id string, task_id int, status bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // timeout for ping
	defer cancel()
	_, err := client.CompleteTask(ctx, &mpb.TaskStatus{
		TaskId:   int32(task_id),
		Status:   status,
		WorkerId: worker_id,
	})
	if err != nil {
		fmt.Printf("Complete task error %v\n", err)
		return err
	}

	return nil
}
