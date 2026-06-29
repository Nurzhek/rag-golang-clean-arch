.PHONY: help tidy run build test vet fmt docker-build docker-up docker-down

BINARY := rag-server

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-14s\033[0m %s\n", $$1, $$2}'

tidy: ## Resolve and lock dependencies (also generates go.sum)
	go mod tidy

run: ## Run the server locally
	go run ./cmd/server

build: ## Build the server binary into ./bin
	go build -o bin/$(BINARY) ./cmd/server

test: ## Run unit tests
	go test ./...

vet: ## Run go vet
	go vet ./...

fmt: ## Format the code
	go fmt ./...

docker-build: ## Build the Docker image
	docker build -t $(BINARY):latest .

docker-up: ## Start via docker compose
	docker compose up --build

docker-down: ## Stop docker compose
	docker compose down
