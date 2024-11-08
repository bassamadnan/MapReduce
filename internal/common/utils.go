package c_utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func GetPartition(key string) int {
	if len(key) == 0 {
		return 0
	}
	return int(key[0]) % 2 // num_reducers
}

func GetLines(start int, end int, filePath string) []string {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to filepath %v", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var lines []string
	line := 0
	for scanner.Scan() {
		if line > end {
			break
		}
		if line >= start {
			lines = append(lines, scanner.Text())
		}
		line++
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("error reading file: %v", err)
	}
	return lines
}

func WriteMapResults(kvPairs []KeyValue, outputDirectory string, taskID int) ([]int, error) {
	partitions := make(map[int]bool) // use map to track unique partitions

	partitioned := make(map[int][]KeyValue)
	for _, kv := range kvPairs {
		partition := GetPartition(kv.Key)
		partitioned[partition] = append(partitioned[partition], kv)
		partitions[partition] = true
	}

	for partition, pairs := range partitioned {
		dirPath := fmt.Sprintf("%v/%v", outputDirectory, partition)
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			return nil, err
		}

		filename := fmt.Sprintf("%v/task_%v.txt", dirPath, taskID)
		file, err := os.Create(filename)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		writer := bufio.NewWriter(file)
		for _, kv := range pairs {
			_, err := fmt.Fprintf(writer, "%s %d\n", kv.Key, kv.Value)
			if err != nil {
				return nil, err
			}
		}
		err = writer.Flush()
		if err != nil {
			return nil, err
		}
	}

	result := make([]int, 0, len(partitions))
	for p := range partitions {
		result = append(result, p)
	}

	return result, nil
}
