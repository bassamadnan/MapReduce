package m_utils

func GetMapTasks(workerComponents [][]int) ([]Task, error) {
	tasks := make([]Task, len(workerComponents))
	for i, components := range workerComponents {
		tasks[i] = Task{
			TaskID:     i,
			Components: components,
			TaskType:   MAP_TASK,
			TaskStatus: PENDING,
		}
	}
	return tasks, nil
}

func GetAvailableTask(tasks []Task) *Task {
	for _, task := range tasks {
		if task.TaskStatus == PENDING {
			return &task
		}
	}
	return nil
}

func GetWorkerID(port string) string {
	return "localhost:" + port
}
