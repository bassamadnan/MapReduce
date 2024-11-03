package m_utils

import (
	"context"
	"log"
	wpb "mapreduce/pkg/proto/worker"
	"time"
)

func SendTask(client wpb.WorkerServiceClient, task *Task) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := client.SendTask(ctx, &wpb.TaskDescription{
		Start: int32(task.Start),
		End:   int32(task.End),
	})
	if err != nil {
		log.Fatalf("error in getting sendtask %v\n", err)
	}
}

func (s *Server) AssignTasks(tasks []Task) {
	for i, worker := range s.Workers {
		SendTask(worker.Client, &tasks[i])
	}
}
