package w_utils

import (
	"context"
	"fmt"
	c_utils "mapreduce/internal/common"
	utils "mapreduce/pkg"
	wpb "mapreduce/pkg/proto/worker"
)

func (w *WorkerMachine) execTask(start int, end int, taskID int, filePath string) {
	lines := c_utils.GetLines(start, end, filePath)
	results := utils.Map(lines)
	fmt.Printf("results :%v\n", results)
	outFile := fmt.Sprintf("%v/task_%v.txt", w.OutputDirectory, taskID)
	fmt.Printf("writing to %v\n", outFile)
	c_utils.WriteMapResults(results, outFile)
	// notify master about this next
}

// SendTask(context.Context, *TaskDescription) (*Empty, error)
func (s *Server) SendTask(ctx context.Context, task *wpb.TaskDescription) (*wpb.Empty, error) {
	start, end, taskid := task.Start, task.End, task.TaskID
	fmt.Printf("Recieved task start: %d, %d\n", start, end)
	go s.WorkerMachineInstance.execTask(int(start), int(end), int(taskid), s.InputFile)
	return &wpb.Empty{}, nil
}

// master machine will ping the worker, this will be running on the go routine hosting the worker server
// the tasks executed by a worker will be running in the execTask go routine, so this will never be blocked by other operations
// i.e reply will be instant
func (s *Server) Ping(ctx context.Context, req *wpb.PingRequest) (*wpb.PingResponse, error) {
	return &wpb.PingResponse{
		Status: true,
	}, nil
}
