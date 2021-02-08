run: up go-run

go-run:
	go mod download & go run ./cmd/listener

go-build:
	go mod download & go build -v ./cmd/listener

up: docker-up
down: docker-down
restart: up down

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down --remove-orphans

.DEFAULT_GOAL := run