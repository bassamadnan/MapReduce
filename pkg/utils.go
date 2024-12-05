package utils

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type DisjointSetUnion struct {
	Parent []int
	Rank   []int
}

func NewDSU(n int) *DisjointSetUnion {
	Parent := make([]int, n)
	Rank := make([]int, n)
	for i := range Parent {
		Parent[i] = i
	}
	return &DisjointSetUnion{Parent, Rank}
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
		if dsu.Rank[rootU] > dsu.Rank[rootV] {
			dsu.Parent[rootV] = rootU
		} else if dsu.Rank[rootU] < dsu.Rank[rootV] {
			dsu.Parent[rootU] = rootV
		} else {
			dsu.Parent[rootV] = rootU
			dsu.Rank[rootU]++
		}
	}
}

type Edge struct {
	U, V, W int
}

func ReadMTXFile(filename string) (map[int][]Edge, int) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	const (
		NumWorkers = 4
		ChunkSize  = 1000
	)

	type MTXLine struct {
		U, V, W int
		Valid   bool
	}
	MappedLines := make(chan MTXLine, NumWorkers*ChunkSize)
	Done := make(chan bool)

	var Chunks [][]string
	var CurrentChunk []string
	Scanner := bufio.NewScanner(file)

	for Scanner.Scan() {
		CurrentChunk = append(CurrentChunk, Scanner.Text())
		if len(CurrentChunk) == ChunkSize {
			Chunks = append(Chunks, CurrentChunk)
			CurrentChunk = []string{}
		}
	}
	if len(CurrentChunk) > 0 {
		Chunks = append(Chunks, CurrentChunk)
	}

	var Wg sync.WaitGroup
	for _, Chunk := range Chunks {
		Wg.Add(1)
		go func(Lines []string) {
			defer Wg.Done()
			for _, Line := range Lines {
				if strings.HasPrefix(Line, "%") || len(strings.TrimSpace(Line)) == 0 {
					continue
				}

				Fields := strings.Fields(Line)
				if len(Fields) != 3 {
					continue
				}

				U, Err1 := strconv.Atoi(Fields[0])
				V, Err2 := strconv.Atoi(Fields[1])
				W, Err3 := strconv.Atoi(Fields[2])
				if W < 0 {
					W = -W
				}

				if Err1 == nil && Err2 == nil && Err3 == nil {
					MappedLines <- MTXLine{U, V, W, true}
				}

				if Err1 == nil && Err2 == nil && Err3 == nil {
					MappedLines <- MTXLine{V, U, W, true}
				}
			}
		}(Chunk)
	}

	AdjList := make(map[int][]Edge)
	MaxVertex := 0
	MinVertex := math.MaxInt32
	var Mu sync.Mutex

	go func() {
		for Line := range MappedLines {
			if !Line.Valid {
				continue
			}

			Mu.Lock()
			if Line.U < MinVertex {
				MinVertex = Line.U
			}
			if Line.V < MinVertex {
				MinVertex = Line.V
			}
			if Line.U > MaxVertex {
				MaxVertex = Line.U
			}
			if Line.V > MaxVertex {
				MaxVertex = Line.V
			}
			Mu.Unlock()
		}
		Done <- true
	}()

	Wg.Wait()
	close(MappedLines)
	<-Done

	file.Seek(0, 0)
	Scanner = bufio.NewScanner(file)
	for Scanner.Scan() {
		Line := Scanner.Text()
		if strings.HasPrefix(Line, "%") || len(strings.TrimSpace(Line)) == 0 {
			continue
		}

		Fields := strings.Fields(Line)
		if len(Fields) != 3 {
			continue
		}

		U, Err1 := strconv.Atoi(Fields[0])
		V, Err2 := strconv.Atoi(Fields[1])
		W, Err3 := strconv.Atoi(Fields[2])
		if W < 0 {
			W = -W
		}

		if Err1 == nil && Err2 == nil && Err3 == nil {
			NormalizedU := U - MinVertex
			NormalizedV := V - MinVertex
			AdjList[NormalizedU] = append(AdjList[NormalizedU],
				Edge{NormalizedU, NormalizedV, W})
			AdjList[NormalizedV] = append(AdjList[NormalizedV],
				Edge{NormalizedV, NormalizedU, W})
		}
	}

	MaxVertex = MaxVertex - MinVertex

	return AdjList, MaxVertex
}

func PrintAdjList(AdjList map[int][]Edge) {
	Vertices := make([]int, 0, len(AdjList))
	for V := range AdjList {
		Vertices = append(Vertices, V)
	}
	sort.Ints(Vertices)

	for _, V := range Vertices {
		fmt.Printf("Vertex %d -> ", V)
		Edges := AdjList[V]
		for I, Edge := range Edges {
			if I > 0 {
				fmt.Print(", ")
			}
			fmt.Printf("(%d, W:%d)", Edge.V, Edge.W)
		}
		fmt.Println()
	}
}

func GetWorkerComponents(Dsu *DisjointSetUnion, NumWorkers int) [][]int {
	ComponentSet := make(map[int]bool)
	for _, Parent := range Dsu.Parent {
		ComponentSet[Dsu.Find(Parent)] = true
	}

	Components := make([]int, 0, len(ComponentSet))
	for Comp := range ComponentSet {
		Components = append(Components, Comp)
	}
	sort.Ints(Components)

	TotalComponents := len(Components)
	WorkerComponents := make([][]int, NumWorkers)
	ComponentsPerWorker := TotalComponents / NumWorkers

	for I := 0; I < NumWorkers; I++ {
		Start := I * ComponentsPerWorker
		End := Start + ComponentsPerWorker
		if I == NumWorkers-1 {
			End = TotalComponents
		}
		WorkerComponents[I] = Components[Start:End]
	}

	return WorkerComponents
}

func GetComponentOutgoingEdges(WorkerComps []int, AdjList map[int][]Edge, Dsu *DisjointSetUnion) map[int][]Edge {
	CompSet := make(map[int]bool)
	for _, Comp := range WorkerComps {
		CompSet[Comp] = true
	}

	CompOutgoing := make(map[int][]Edge)
	for Comp := range CompSet {
		CompOutgoing[Comp] = make([]Edge, 0)
	}

	for U, Edges := range AdjList {
		for _, e := range Edges {
			V := e.V
			if Dsu.Find(U) != Dsu.Find(V) {
				if CompSet[Dsu.Find(U)] {
					CompOutgoing[Dsu.Find(U)] = append(CompOutgoing[Dsu.Find(U)], Edge{U, V, e.W})
				}
				if CompSet[Dsu.Find(V)] {
					CompOutgoing[Dsu.Find(V)] = append(CompOutgoing[Dsu.Find(V)], Edge{V, U, e.W})
				}
			}
		}
	}
	// remove duplicates
	for Comp, edges := range CompOutgoing {
		uniqueEdges := make(map[Edge]bool)
		for _, edge := range edges {
			uniqueEdges[edge] = true
		}

		uniqueEdgeList := make([]Edge, 0, len(uniqueEdges))
		for edge := range uniqueEdges {
			uniqueEdgeList = append(uniqueEdgeList, edge)
		}

		CompOutgoing[Comp] = uniqueEdgeList
	}

	return CompOutgoing
}

func GetPartitionID(ComponentID int, NumReducers int) int {
	return ComponentID % NumReducers
}

func WriteMapResults(CompOutgoing map[int][]Edge, OutputDir string, TaskID int, NumReducers int) ([]int, error) {
	PartitionsUsed := make(map[int]bool)

	for Comp, Edges := range CompOutgoing {
		PartitionID := GetPartitionID(Comp, NumReducers)
		PartitionsUsed[PartitionID] = true

		PartitionDir := fmt.Sprintf("%s/%d", OutputDir, PartitionID)
		if Err := os.MkdirAll(PartitionDir, 0755); Err != nil {
			return nil, Err
		}

		OutFile := fmt.Sprintf("%s/task_%d.txt", PartitionDir, TaskID)
		File, Err := os.Create(OutFile)
		if Err != nil {
			return nil, Err
		}
		defer File.Close()

		// Use map to deduplicate edges for this component
		type EdgeKey struct {
			u, v int
		}
		uniqueEdges := make(map[EdgeKey]Edge)
		for _, Edge := range Edges {
			// Create canonical edge key (smaller vertex first)
			u, v := Edge.U, Edge.V
			if u > v {
				u, v = v, u
			}
			key := EdgeKey{u: u, v: v}

			// keep edge if we haven't seen it, or if it has lower weight
			if existing, exists := uniqueEdges[key]; !exists || Edge.W < existing.W {
				uniqueEdges[key] = Edge
			}
		}

		// unique edges
		for _, Edge := range uniqueEdges {
			// Always write with smaller vertex ID first
			u, v := Edge.U, Edge.V
			if u > v {
				u, v = v, u
			}
			_, Err := fmt.Fprintf(File, "%d %d %d %d\n", Comp, u, v, Edge.W)
			if Err != nil {
				return nil, Err
			}
		}
	}

	Partitions := make([]int, 0, len(PartitionsUsed))
	for P := range PartitionsUsed {
		Partitions = append(Partitions, P)
	}
	sort.Ints(Partitions)

	return Partitions, nil
}

func ReadDirectoryEdges(directory string) (map[int][]Edge, error) {
	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %v", err)
	}

	compEdges := make(map[int][]Edge)
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(directory, file.Name())
		file, err := os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("error opening file %s: %v", filePath, err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			fields := strings.Fields(line)
			if len(fields) != 4 { // Now expecting 4 fields
				continue
			}

			comp, err1 := strconv.Atoi(fields[0])
			u, err2 := strconv.Atoi(fields[1])
			v, err3 := strconv.Atoi(fields[2])
			w, err4 := strconv.Atoi(fields[3])

			if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
				continue
			}

			edge := Edge{U: u, V: v, W: w}
			compEdges[comp] = append(compEdges[comp], edge)
		}
	}

	return compEdges, nil
}

func GetComponents(dsu *DisjointSetUnion) []int {
	componentSet := make(map[int]bool)
	for i := range dsu.Parent {
		componentSet[dsu.Find(i)] = true
	}

	components := make([]int, 0, len(componentSet))
	for comp := range componentSet {
		components = append(components, comp)
	}
	return components
}
