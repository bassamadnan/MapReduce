build_proto:
	protoc --go_out=. --go-grpc_out=. pkg/proto/master.proto
	protoc --go_out=. --go-grpc_out=. pkg/proto/worker.proto
