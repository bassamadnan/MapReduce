package utils

import (
	c_utils "mapreduce/internal/common"
	"strings"
)

// place to keep the map and reduce functions for now

/*
map (k1,v1) → list(k2,v2)
reduce (k2,list(v2)) → list(v2)
*/

// could take in line number, value but then would have to do mini reduce
func Map(lines []string) []c_utils.KeyValue {
	wordCount := make(map[string]int)
	for _, line := range lines {
		words := strings.Fields(line)
		for _, word := range words {
			wordCount[word]++
		}
	}
	kv_pairs := make([]c_utils.KeyValue, 0, len(wordCount))
	for word, count := range wordCount {
		kv_pairs = append(kv_pairs, c_utils.KeyValue{
			Key:   word,
			Value: count,
		})
	}

	return kv_pairs
}
