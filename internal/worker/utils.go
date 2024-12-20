package w_utils

import (
	"fmt"
	"log"
	utils "mapreduce/pkg"
	mpb "mapreduce/pkg/proto/master"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (w *WorkerMachine) Initialize(master_addr string, port string, graphFilePath string) {
	w.MasterAddr = master_addr
	w.ID = port
	outputPath := fmt.Sprintf("data/output/localhost%v", port)
	os.RemoveAll(outputPath)
	os.MkdirAll(outputPath, 0777)
	w.ServerInstance = &Server{InputFile: "data/input/input_1.txt", WorkerMachineInstance: w}
	w.OutputDirectory = outputPath
	w.SetupWorkerMachine(w.MasterAddr)
	w.ServerInstance.ReduceResults = make(map[string]int)
	adjList, _ := utils.ReadMTXFile(graphFilePath)
	w.AdjList = adjList
	w.DSU = utils.NewDSU(len(w.AdjList))
	w.NumVertices = len(w.AdjList)
	utils.PrintAdjList(w.AdjList)
	numReducers, _ := PingReady(w.Client, port)
	w.NumReducers = numReducers
}

func (w *WorkerMachine) SetupWorkerMachine(master_addr string) {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	conn, err := grpc.NewClient(master_addr, opts...)
	if err != nil {
		log.Fatalf("Couldnt setup client workermachine %v\n", err)
	}
	w.Client = mpb.NewMasterServiceClient(conn)
	w.Conn = conn
}

func (w *WorkerMachine) CloseConnection() {
	if w.Conn != nil {
		w.Conn.Close()
	}

}
