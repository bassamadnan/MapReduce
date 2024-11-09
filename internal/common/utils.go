package c_utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

// do we sort?
func GetPartitionData(directory string) []KeyValue {
	var kvList []KeyValue

	files, err := os.ReadDir(directory)
	if err != nil {
		log.Printf("Error reading directory %s: %v", directory, err)
		return kvList
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(strings.ToLower(file.Name()), ".txt") {
			continue
		}

		data, err := os.ReadFile(filepath.Join(directory, file.Name()))
		if err != nil {
			log.Printf("Error reading file %s: %v", file.Name(), err)
			continue
		}

		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			parts := strings.Fields(line)
			if len(parts) != 2 {
				continue
			}

			value, err := strconv.Atoi(parts[1])
			if err != nil {
				continue
			}

			kv := KeyValue{
				Key:   parts[0],
				Value: value,
			}
			kvList = append(kvList, kv)
		}
	}

	return kvList
}
