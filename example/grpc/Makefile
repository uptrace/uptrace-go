all:
	go run main.go

protoc:
	protoc -I api --go_out=plugins=grpc,paths=source_relative:./api api/hello-service.proto
