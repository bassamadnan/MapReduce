package m_utils

import (
	"context"
	"log"
	wpb "mapreduce/pkg/proto/worker"
	"sync"
	"time"
)

// client side functions for grpc service of the master, to communicate with worker service
type MasterClient struct {
	Client          wpb.WorkerServiceClient
	Mu              sync.Mutex
	ServiceRegistry []string
}

func (c *MasterClient) SendTask(task *Task) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := c.Client.SendTask(ctx, &wpb.TaskDescription{
		Start: int32(task.Start),
		End:   int32(task.End),
	})
	if err != nil {
		log.Fatalf("error in getting sendtask %v\n", err)
	}
}
