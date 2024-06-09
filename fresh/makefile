all:
	protoc --go_out=./services/three --go_opt=paths=source_relative \
    --go-grpc_out=./services/three --go-grpc_opt=paths=source_relative \
    threes.proto

clean:
	rm -rf $(wildcard services/*)