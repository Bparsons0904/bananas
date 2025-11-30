.PHONY: test build run-standard run-gin run-fiber run-echo run-chi run-gorilla migrate-up migrate-down seed docker-up docker-down

# Run tests
test:
	cd server && go test ./...

# Build the application
build:
	cd server && go build -o bin/standard ./cmd/api/main.go
	cd server && go build -o bin/gin ./cmd/api/gin.go
	cd server && go build -o bin/fiber ./cmd/api/fiber.go
	cd server && go build -o bin/echo ./cmd/api/echo.go
	cd server && go build -o bin/chi ./cmd/api/chi.go
	cd server && go build -o bin/gorilla ./cmd/api/gorilla.go

# Run all frameworks together
run-all:
	cd server && go run ./cmd/api/all-frameworks.go

# Run individual frameworks (for testing)
run-standard:
	cd server && go run ./cmd/api/main.go

run-gin:
	cd server && go run ./cmd/api/gin.go

run-fiber:
	cd server && go run ./cmd/api/fiber.go

run-echo:
	cd server && go run ./cmd/api/echo.go

run-chi:
	cd server && go run ./cmd/api/chi.go

run-gorilla:
	cd server && go run ./cmd/api/gorilla.go

# Database operations
migrate-up:
	cd server && go run cmd/migration/main.go up

migrate-down:
	cd server && go run cmd/migration/main.go down

seed:
	cd server && go run cmd/migration/main.go seed

# Docker operations
docker-up:
	docker compose -f docker-compose.dev.yml up -d

docker-down:
	docker compose -f docker-compose.dev.yml down

docker-logs:
	docker compose -f docker-compose.dev.yml logs -f

# Tilt
dev:
	tilt up

dev-down:
	tilt down

# Install dependencies
deps:
	cd server && go mod download
	cd server && go mod tidy

# Clean build artifacts
clean:
	cd server && rm -rf bin/ tmp/ *.log