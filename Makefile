APP_NAME := sub-hub

.PHONY: tidy test run build lint generate docker-up docker-down

tidy:
	go mod tidy

test:
	go test ./...

run:
	go run ./cmd/$(APP_NAME)

build:
	go build -o bin/$(APP_NAME) ./cmd/$(APP_NAME)

lint:
	golangci-lint run ./...

generate:
	./scripts/generate.sh

docker-up:
	docker compose -f deploy/compose/docker-compose.yml up -d

docker-down:
	docker compose -f deploy/compose/docker-compose.yml down

