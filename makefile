all:
	protoc --go_out=./service --go_opt=paths=source_relative \
    --go-grpc_out=./service --go-grpc_opt=paths=source_relative \
    threes.proto

clean:
	rm -rf $(wildcard services/*)