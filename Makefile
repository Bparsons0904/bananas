.PHONY: test build run migrate-up migrate-down seed docker-up docker-down dev dev-down deps clean

# Run tests
test:
	cd server && go test ./...

# Build the application (all 6 frameworks in one binary)
build:
	cd server && go build -o bin/api ./cmd/api

# Run all 6 frameworks
run:
	cd server && go run ./cmd/api

# Database operations
create-db:
	cd server && go run cmd/migration/main.go create-db

migrate-up:
	cd server && go run cmd/migration/main.go up

migrate-down:
	cd server && go run cmd/migration/main.go down

seed:
	cd server && go run cmd/migration/main.go seed

# Complete database setup (for production/fresh installs)
db-setup: create-db migrate-up seed

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