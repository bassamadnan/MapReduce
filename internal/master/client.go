package m_utils

import (
	"context"
	"fmt"
	wpb "mapreduce/pkg/proto/worker"
	"sync"
	"time"
)

func SendMapTask(client wpb.WorkerServiceClient, task *Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := client.SendMapTask(ctx, &wpb.TaskDescription{
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

func (s *Server) AssignMapTasks() {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	for i, worker := range s.Workers {
		if s.Workers[i].Status != IDLE || s.Workers[i].WorkerType != MAPPER {
			continue
		}
		task := GetAvailableTask(s.Tasks)
		if task == nil {
			fmt.Print("all tasks over\n")
			return
		}
		err := SendMapTask(worker.Client, task)
		if err != nil {
			s.Workers[i].Status = FAIL
			s.Tasks[task.TaskID].TaskStatus = PENDING
		} else {
			s.Tasks[task.TaskID].TaskStatus = ASSIGNED
		}
		fmt.Printf("assigned task %v to worker %v, worker: %v\n", task.TaskID, i, worker)
	}
}

func SendPing(client wpb.WorkerServiceClient, worker_id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // timeout for ping
	defer cancel()
	fmt.Printf("\033[34mPing request to client :%v\033[0m\n", worker_id) // blue color for request
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
		for i, worker := range s.Workers {
			if s.Workers[i].Status == FAIL {
				continue
			}
			wg.Add(1)
			go func(w *Worker) {
				defer wg.Add(-1)
				err := SendPing(w.Client, w.Addr)
				if err != nil {
					s.Mu.Lock()
					s.Workers[i].Status = FAIL
					// re assign task assigned to it
					if s.Workers[i].AssignedTask != -1 {
						s.Tasks[s.Workers[i].AssignedTask].TaskStatus = PENDING
					}
					s.Mu.Unlock()
				}
			}(worker)
		}

		wg.Wait()
		<-ticker.C
	}
}

func SendReducerTask(client wpb.WorkerServiceClient, addr []string, partition int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := client.SendReduceTask(ctx, &wpb.ReduceTaskDescription{
		Addr:      addr,
		Partition: int32(partition),
	})
	if err != nil {
		fmt.Printf("error in  send reduce task %v to client %v\n", err, client)
		return err
	}
	return nil
}

/*
Tasks contain a map string, partitions
string -> localhost address of map worker machine
partition -> list of numbers

We have to assign reduce worker machine tasks now,
*/
func (s *Server) AssignReducerTasks() {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	reducerInputs := make(map[int]map[string]bool)

	for _, task := range s.Tasks {
		for addr, partitions := range task.OutputPartitions {
			for _, partition := range partitions {
				if reducerInputs[partition] == nil {
					reducerInputs[partition] = make(map[string]bool)
				}
				reducerInputs[partition][addr] = true
			}
		}
	}

	reducerLists := make(map[int][]string)
	for partition, addrSet := range reducerInputs {
		addrs := make([]string, 0, len(addrSet))
		for addr := range addrSet {
			addrs = append(addrs, addr)
		}
		reducerLists[partition] = addrs
	}
	for partition, addr := range reducerLists {
		fmt.Print(partition, addr)
	}
	partition := 0
	for i, worker := range s.Workers {
		if s.Workers[i].Status != IDLE || s.Workers[i].WorkerType != REDUCER {
			continue
		}
		SendReducerTask(worker.Client, reducerLists[partition], partition)
		partition++
	}
}
