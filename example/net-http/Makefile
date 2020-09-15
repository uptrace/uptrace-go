all: build up test

build:
	docker-compose build
up:
	docker-compose up --detach http-server
down:
	docker-compose down
logs:
	docker-compose logs http-server
test:
	curl http://localhost:9999/profiles/admin
	curl http://localhost:9999/profiles/foo
