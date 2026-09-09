# .PHONY: build run dev lint migrate-up migrate-down help docker-up docker-down

# help:
# 	@echo "Makefile commands:"
# 	@echo "  build          Build the Go application"
# 	@echo "  run            Run the Go application"
# 	@echo "  dev            Run golangci-lint for code linting"
# 	@echo "  migrate-up     Apply database migrations"
# 	@echo "  migrate-down   Rollback database migrations"

# build:
# 	go build -o bin/api cmd/api/main.go

# dev:
# 	go run ./cmd/api

# lint:
# 	golangci-lint run ./...

# migrate-up:
# 	migrate -path db/migrations -database  "postgres://postgres:postgres@localhost:5432/ecommerce_shop?sslmode=disable" up

# migrate-down:
# 	migrate -path db/migrations -database  "postgres://postgres:postgres@localhost:5432/ecommerce_shop?sslmode=disable" down

# docker-up:
# 	docker compose -f docker/docker-compose.yml up -d	

# docker-down:
# 	docker compose -f docker/docker-compose.yml down

.PHONY: build run dev lint migrate-up migrate-down docker-up docker-down help

help:
	@echo "Makefile commands:"
	@echo "  build          Build the Go application"
	@echo "  run            Run the Go application"
	@echo "  dev            Run the Go application in development mode"
	@echo "  lint           Run golangci-lint for code linting"
	@echo "  migrate-up     Apply database migrations"
	@echo "  migrate-down   Rollback database migrations"
	@echo "  docker-up      Start Docker containers"
	@echo "  docker-down    Stop Docker containers"

build:
	go build -o bin/api cmd/api/main.go

run:
	go run cmd/api/main.go

dev:
	go run ./cmd/api

lint: format
	@command -v golangci-lint >/dev/null 2>&1 || { echo "golangci-lint not installed. Installing..."; curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.54.2; }
	golangci-lint run ./...

format:
	@gofmt -s -w .	
	@goimports -w .	

migrate-up:
	migrate -path db/migrations -database "postgres://postgres:password@localhost:5432/ecommerce_shop?sslmode=disable" up

migrate-down:
	migrate -path db/migrations -database "postgres://postgres:password@localhost:5432/ecommerce_shop?sslmode=disable" down

docker-up:
	docker compose -f docker/docker-compose.yml up -d

docker-down:
	docker compose -f docker/docker-compose.yml down