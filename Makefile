APP_NAME=koalbot_api
BIN_DIR=bin

.PHONY: build run tidy test docker-build compose-up compose-down up down clean prod ui-tidy

build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/api ./cmd/api

run:
	go run ./cmd/api

tidy:
	go mod tidy

test:
	go test ./...

docker-build:
	docker build -t $(APP_NAME):latest .

up:
	docker compose up -d --build

down:
	docker compose down

clean:
	docker compose down -v

ui-tidy:
	cd koalbot-ui && npm ci

prod: tidy ui-tidy
	docker compose up -d --build
