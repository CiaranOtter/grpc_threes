all:
	protoc --go_out=./gameservice --go_opt=paths=source_relative \
    --go-grpc_out=./gameservice --go-grpc_opt=paths=source_relative \
    threes.proto

clean:
	rm -rf $(wildcard services/*)