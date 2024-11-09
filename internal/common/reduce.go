package c_utils

import "sort"

func Reduce(kvList []KeyValue) []KeyValue {
    reduced := make(map[string]int)

    for _, kv := range kvList {
        reduced[kv.Key] += kv.Value
    }

    var result []KeyValue
    for key, value := range reduced {
        result = append(result, KeyValue{
            Key:   key,
            Value: value,
        })
    }

    sort.Slice(result, func(i, j int) bool {
        return result[i].Key < result[j].Key
    })

    return result
}
