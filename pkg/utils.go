package utils

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type DisjointSetUnion struct {
	Parent []int
	rank   []int
}

func NewDSU(n int) *DisjointSetUnion {
	Parent := make([]int, n)
	rank := make([]int, n)
	for i := range Parent {
		Parent[i] = i
	}
	return &DisjointSetUnion{Parent, rank}
}

func (dsu *DisjointSetUnion) Find(u int) int {
	if dsu.Parent[u] != u {
		dsu.Parent[u] = dsu.Find(dsu.Parent[u])
	}
	return dsu.Parent[u]
}

func (dsu *DisjointSetUnion) Union(u, v int) {
	rootU := dsu.Find(u)
	rootV := dsu.Find(v)

	if rootU != rootV {
		if dsu.rank[rootU] > dsu.rank[rootV] {
			dsu.Parent[rootV] = rootU
		} else if dsu.rank[rootU] < dsu.rank[rootV] {
			dsu.Parent[rootU] = rootV
		} else {
			dsu.Parent[rootV] = rootU
			dsu.rank[rootU]++
		}
	}
}

type Edge struct {
	u, v, w int
}

// everything starts from 0
func ReadMTXFile(filename string) (map[int][]Edge, int) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	const (
		numWorkers = 4
		chunkSize  = 1000
	)

	type MTXLine struct {
		u, v, w int
		valid   bool
	}
	mappedLines := make(chan MTXLine, numWorkers*chunkSize)
	done := make(chan bool)

	var chunks [][]string
	var currentChunk []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		currentChunk = append(currentChunk, scanner.Text())
		if len(currentChunk) == chunkSize {
			chunks = append(chunks, currentChunk)
			currentChunk = []string{}
		}
	}
	if len(currentChunk) > 0 {
		chunks = append(chunks, currentChunk)
	}

	var wg sync.WaitGroup
	for _, chunk := range chunks {
		wg.Add(1)
		go func(lines []string) {
			defer wg.Done()
			for _, line := range lines {
				if strings.HasPrefix(line, "%") || len(strings.TrimSpace(line)) == 0 {
					continue
				}

				fields := strings.Fields(line)
				if len(fields) != 3 {
					continue
				}

				u, err1 := strconv.Atoi(fields[0])
				v, err2 := strconv.Atoi(fields[1])
				w, err3 := strconv.Atoi(fields[2])
				if w < 0 {
					w = -w
				}

				if err1 == nil && err2 == nil && err3 == nil {
					mappedLines <- MTXLine{u, v, w, true}
				}

				if err1 == nil && err2 == nil && err3 == nil {
					mappedLines <- MTXLine{v, u, w, true}
				}
			}
		}(chunk)
	}

	adjList := make(map[int][]Edge)
	maxVertex := 0
	minVertex := math.MaxInt32
	var mu sync.Mutex

	go func() {
		for line := range mappedLines {
			if !line.valid {
				continue
			}

			mu.Lock()
			// Track minimum and maximum vertex numbers
			if line.u < minVertex {
				minVertex = line.u
			}
			if line.v < minVertex {
				minVertex = line.v
			}
			if line.u > maxVertex {
				maxVertex = line.u
			}
			if line.v > maxVertex {
				maxVertex = line.v
			}
			mu.Unlock()
		}
		done <- true
	}()

	wg.Wait()
	close(mappedLines)
	<-done

	file.Seek(0, 0)
	scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "%") || len(strings.TrimSpace(line)) == 0 {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) != 3 {
			continue
		}

		u, err1 := strconv.Atoi(fields[0])
		v, err2 := strconv.Atoi(fields[1])
		w, err3 := strconv.Atoi(fields[2])
		if w < 0 {
			w = -w
		}

		if err1 == nil && err2 == nil && err3 == nil {
			// Normalize vertices to start from 0
			normalizedU := u - minVertex
			normalizedV := v - minVertex
			adjList[normalizedU] = append(adjList[normalizedU],
				Edge{normalizedU, normalizedV, w})
			adjList[normalizedV] = append(adjList[normalizedV],
				Edge{normalizedV, normalizedU, w})
		}
	}

	// Adjust maxVertex to reflect normalized numbering
	maxVertex = maxVertex - minVertex

	return adjList, maxVertex
}
func PrintAdjList(adjList map[int][]Edge) {
	vertices := make([]int, 0, len(adjList))
	for v := range adjList {
		vertices = append(vertices, v)
	}
	sort.Ints(vertices)

	for _, v := range vertices {
		fmt.Printf("Vertex %d -> ", v)
		edges := adjList[v]
		for i, edge := range edges {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Printf("(%d, w:%d)", edge.v, edge.w)
		}
		fmt.Println()
	}
}

func GetWorkerComponents(dsu *DisjointSetUnion, numWorkers int) [][]int {
	componentSet := make(map[int]bool)
	for _, parent := range dsu.Parent {
		componentSet[dsu.Find(parent)] = true
	}

	components := make([]int, 0, len(componentSet))
	for comp := range componentSet {
		components = append(components, comp)
	}
	sort.Ints(components)

	totalComponents := len(components)
	workerComponents := make([][]int, numWorkers)
	componentsPerWorker := totalComponents / numWorkers

	for i := 0; i < numWorkers; i++ {
		start := i * componentsPerWorker
		end := start + componentsPerWorker
		if i == numWorkers-1 {
			end = totalComponents
		}
		workerComponents[i] = components[start:end]
	}

	return workerComponents
}

func GetComponentOutgoingEdges(workerComps []int, adjList map[int][]Edge, dsu *DisjointSetUnion) map[int][]Edge {
	compSet := make(map[int]bool)
	for _, comp := range workerComps {
		compSet[comp] = true
	}

	compOutgoing := make(map[int][]Edge)
	for comp := range compSet {
		compOutgoing[comp] = make([]Edge, 0)
	}

	for u, edges := range adjList {
		for _, edge := range edges {
			v := edge.v
			if dsu.Find(u) != dsu.Find(v) {
				if compSet[dsu.Find(u)] {
					compOutgoing[dsu.Find(u)] = append(compOutgoing[dsu.Find(u)], Edge{u, v, edge.w})
				}
				if compSet[dsu.Find(v)] {
					compOutgoing[dsu.Find(v)] = append(compOutgoing[dsu.Find(v)], Edge{v, u, edge.w})
				}
			}
		}
	}

	return compOutgoing
}
