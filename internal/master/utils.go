package m_utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func GetMapTasks(job *Job) ([]Task, error) {
	// first line of input file will have number of lines
	file, err := os.Open(job.InputFileName)
	if err != nil {
		log.Fatalf("job input filename error %v\n", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	ok := scanner.Scan()

	if !ok {
		log.Printf("error totallinse %v\n", scanner.Err())
	}
	totalLines, err := strconv.Atoi(scanner.Text())
	if err != nil {
		log.Fatalf("total linse error %v\n", err)
	}
	fmt.Printf("totalLines %v, numworkers: %v\n", totalLines, job.NumWorkers)
	// split tasks equally among workers
	linesPerSplit := totalLines / job.Split
	extraLines := totalLines % job.Split
	tasks := make([]Task, job.Split)
	currentLine := 1
	for i := 0; i < job.Split; i++ {
		lines := linesPerSplit
		if extraLines > 0 {
			lines++
			extraLines--
		}

		tasks[i] = Task{
			TaskID:   i,
			Start:    currentLine,
			End:      currentLine + lines - 1,
			TaskType: MAP_TASK,
		}
		currentLine += lines
	}

	return tasks, nil

}
