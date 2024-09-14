# Setup
LOCALHOST := "127.0.0.1"
SERVER_PORT := "8080"
PROJECT_NAME := "tcp-network"

# Default command
default:
    @just --list --unsorted

# Run UDP server
run port:
    go run ./src/cmd/main.go --port={{port}}

# Execute unit tests
test:
    @echo "Running unit tests!"
    go clean -testcache
    go test -cover ./src/...

# Build Docker image
build tag="latest": test
	@echo "Building Docker image (tag={{ tag }})..."
	docker build -t {{ PROJECT_NAME }}:{{ tag }} -f ./Dockerfile .
	@echo "Docker image built successfully!"

# Sync Go modules
tidy:
    cd src && go mod tidy
    go work sync
    echo "Go workspace and modules synced successfully!"

# Start Compose with load-balancer
compose-up:
    docker compose -f docker-compose.yml up

# Stop all Compose containers and delete images created
compose-down:
    docker compose -f docker-compose.yml down
    docker rmi $(docker images | grep "socket" | awk "{print \$3}")
