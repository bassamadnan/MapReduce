package m_utils

import (
	"context"
	"fmt"
	wpb "mapreduce/pkg/proto/worker"
	"time"
)

func SendTask(client wpb.WorkerServiceClient, task *Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := client.SendTask(ctx, &wpb.TaskDescription{
		Start:  int32(task.Start),
		End:    int32(task.End),
		TaskID: int32(task.TaskID),
	})
	if err != nil {
		fmt.Printf("error in  sendtask %v to client %v\n", err, client)
		return err
	}
	return nil

}

func (s *Server) AssignTasks(tasks []Task) {
	for i, worker := range s.Workers {
		if s.Workers[i].Status != IDLE {
			continue
		}
		err := SendTask(worker.Client, &tasks[i])
		if err != nil {
			s.Workers[i].Status = FAIL
		}
		// tasks[i] = ASSIGNED
	}
}
