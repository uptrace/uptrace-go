all: build up

build:
	protoc -I api --go_out=plugins=grpc,paths=source_relative:./api api/hello-service.proto
	docker-compose build
up:
	docker-compose up --detach grpc-server grpc-client
down:
	docker-compose down
logs:
	docker-compose logs grpc-server grpc-client

