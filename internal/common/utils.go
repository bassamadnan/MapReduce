package c_utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

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

func WriteMapResults(kvPairs []KeyValue, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, kv := range kvPairs {
		_, err := fmt.Fprintf(writer, "%s %d\n", kv.Key, kv.Value)
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}
