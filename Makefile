compile:
	protoc -I=. --go_out=plugins=grpc,paths=source_relative:. api/v1/log.proto

test:
	go test -race ./...