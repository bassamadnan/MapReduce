package m_utils

import (
	"context"
	"fmt"
	wpb "mapreduce/pkg/proto/worker"
	"sync"
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

func SendPing(client wpb.WorkerServiceClient, worker_id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // timeout for ping
	defer cancel()
	fmt.Printf("\033[34mPing reqiest to client :%v\033[0m\n", worker_id) // blue color for request
	resp, err := client.Ping(ctx, &wpb.PingRequest{
		Id: worker_id,
	})
	if err != nil {
		fmt.Printf("\033[31merror in  ping %v to client %v\033[0m\n", err, client) // red color for errors
		// disable this worker machine connection
		return err
	}

	fmt.Printf("\033[32mPing response from client :%v status: %v\033[0m \n", worker_id, resp.Status) // green color for success
	// handle fail situation, colorize the print statement above
	return nil
}

func (s *Server) StartPing() {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		var wg sync.WaitGroup
		for _, worker := range s.Workers {
			wg.Add(1)
			go func(w *Worker) {
				defer wg.Add(-1)
				SendPing(w.Client, w.Addr)
			}(worker)
		}

		wg.Wait()
		<-ticker.C
	}
}
