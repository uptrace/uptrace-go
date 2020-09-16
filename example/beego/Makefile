all: build up test

build:
	docker-compose build
up:
	docker-compose up --detach beego-server
down:
	docker-compose down
logs:
	docker-compose logs beego-server
test:
	curl http://localhost:9999/profiles/admin
	curl http://localhost:9999/profiles/foo
