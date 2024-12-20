package c_utils

import "strings"

const NUM_REDUCERS = 2 // number of partitions, the R of hash(key)modR

// place to keep the map and reduce functions for now

/*
map (k1,v1) → list(k2,v2)
reduce (k2,list(v2)) → list(v2)
*/

// could take in line number, value but then would have to do mini reduce
func Map(lines []string) []KeyValue {
	wordCount := make(map[string]int)
	for _, line := range lines {
		words := strings.Fields(line)
		for _, word := range words {
			wordCount[word]++
		}
	}
	kv_pairs := make([]KeyValue, 0, len(wordCount))
	for word, count := range wordCount {
		kv_pairs = append(kv_pairs, KeyValue{
			Key:   word,
			Value: count,
		})
	}

	return kv_pairs
}
