# AGENTS.md - Bananas Framework Testing Service

This document helps AI agents work effectively with the Bananas Framework Testing Service codebase.

## Project Overview

Bananas is a comprehensive testing service for comparing performance and characteristics of different Go web frameworks and frontend frameworks. The project allows running 6 different Go web frameworks simultaneously with shared business logic for fair comparison.

## Project Structure

```
bananas/
├── server/                 # Go backend code
│   ├── cmd/               # Main entry points for each framework
│   │   ├── api/          # Framework entry points
│   │   │   ├── all-frameworks.go   # Runs all 6 frameworks together
│   │   │   ├── main.go             # Standard Library
│   │   │   ├── gin.go              # Gin
│   │   │   ├── fiber.go            # Fiber
│   │   │   ├── echo.go             # Echo
│   │   │   ├── chi.go              # Chi
│   │   │   └── gorilla.go          # Gorilla Mux
│   │   └── migration/     # Database migration tool
│   ├── internal/         # Internal application code
│   │   ├── app/         # Application structure
│   │   ├── config/      # Configuration management
│   │   ├── controllers/ # Shared controllers (base_controller.go)
│   │   ├── database/    # Database layer
│   │   ├── logger/      # Logging utilities
│   │   ├── models/      # Data models
│   │   ├── repositories/ # Data access layer
│   │   └── services/    # Business logic
│   ├── config/          # Configuration files
│   ├── pkg/            # Public packages
│   └── Dockerfile.dev  # Development Dockerfile
├── clients/            # Frontend frameworks (to be implemented)
│   ├── react/
│   ├── vue/
│   ├── svelte/
│   ├── solid/
│   ├── angular/
│   ├── htmx/
│   └── templ/
├── database/          # Database initialization scripts
├── scripts/           # Utility scripts
│   ├── setup-dev.sh     # Development setup script
│   └── test-endpoints.sh # Endpoint testing script
└── Makefile           # Build and run commands
```

## Development Commands

### Setting Up Development Environment

```bash
# Copy environment variables
cp .env.example .env

# Start development environment (sets up PostgreSQL, runs migrations, seeds data)
./scripts/setup-dev.sh

# Start all frameworks with hot reloading (via Tilt)
tilt up

# Start all frameworks locally (without Docker)
make run-all
```

### Running Frameworks

All frameworks can run simultaneously on different ports:

```bash
# All frameworks together
make run-all                    # Runs all 6 frameworks at once

# Individual frameworks
make run-standard    # Port 8081
make run-gin        # Port 8082
make run-fiber      # Port 8083
make run-echo       # Port 8084
make run-chi        # Port 8085
make run-gorilla    # Port 8086
```

### Docker Commands

```bash
# Start services with Docker
docker compose -f docker-compose.dev.yml up -d

# Stop services
docker compose -f docker-compose.dev.yml down

# View logs
docker compose -f docker-compose.dev.yml logs -f
```

### Database Operations

```bash
# Run migrations (via Make)
make migrate-up

# Rollback migrations
make migrate-down

# Seed database
make seed

# Via Tilt
tilt trigger migrate-up
tilt trigger migrate-down
tilt trigger migrate-seed
```

### Testing

```bash
# Run Go tests
make test

# Test all endpoints
./scripts/test-endpoints.sh

# Via Tilt
tilt trigger server-1-check-all
```

### Building

```bash
# Build all framework binaries
make build

# Build and clean
make clean
```

## Architecture Patterns

### Shared Application Structure

All backend frameworks share the same:
- Controllers (`internal/controllers/base_controller.go`)
- Services (`internal/services/service.go`)
- Repositories (`internal/repositories/repository.go`)
- Database layer (`internal/database/database.go`)
- Models and configuration

### Framework Entry Points

Each framework has its own entry point in `server/cmd/api/` that:
1. Sets up the framework-specific router
2. Applies framework middleware
3. Adds framework context to requests (`context.WithValue(r.Context(), "framework", "<name>")`)
4. Calls the shared controllers

### API Endpoints

All frameworks expose the same API endpoints:
- `GET /health` - Health check
- `GET /api/test/simple` - Simple request test
- `GET /api/test/database?limit=N` - Database query test
- `GET /api/test/json` - JSON response test
- `GET /api/info` - Framework information

### Application Initialization Pattern

The app follows dependency injection via the `internal/app/app.go`:
1. Load configuration
2. Initialize database connection
3. Create repositories
4. Create services
5. Create controllers
6. Wire everything together

### Logging Pattern

All logging uses the custom logger in `internal/logger/`:
```go
log := logger.New("component")
log.Info("message", arg1, arg2)
log.Er("error message", err)
```

## Code Conventions

### Naming
- Use clear, descriptive names
- Framework identifiers: "standard", "gin", "fiber", "echo", "chi", "gorilla"
- Controllers follow `XxxRequest` pattern for handlers
- Services use descriptive method names

### Error Handling
- Always return errors from functions
- Use structured logging with logger.Er() for errors
- Return appropriate HTTP status codes

### Context Usage
- Pass context through all layers
- Framework name is added to context as "framework" key
- Use context for database operations

### Testing
- Write unit tests for business logic
- Use the endpoint test script for integration testing
- Test each framework's unique behavior

## Development Workflow

1. **Code changes**: Edit shared code in `internal/` or framework-specific code in `cmd/api/`
2. **Hot reloading**: Air automatically restarts servers on file changes
3. **Testing**: Run the endpoint test script to verify all frameworks work
4. **Database changes**: Create migration scripts and run them via make/tilt commands

## Important Gotchas

### Framework Context
Each framework must add its name to the request context:
```go
context.WithValue(r.Context(), "framework", "<framework-name>")
```

### Fiber Compatibility
Fiber requires a special adapter (`fiberResponseWriter`) to work with standard Go HTTP handlers.

### Port Configuration
Frameworks use fixed ports (8081-8086) defined in `.env` file.

### Database Connection
Database uses PostgreSQL with connection details from environment variables.

### Docker Development
Docker file has been fixed to use the new Air package path (`github.com/air-verse/air`).

### Simultaneous Execution
All frameworks are designed to run simultaneously from a single Go application via `all-frameworks.go`.

## Future Development Plans

1. **Multiple ORM Support**: Adding GORM, SQLx, and PGX implementations
2. **Frontend Clients**: Implementing React, Vue, Svelte, Solid, Angular, HTMX, and Templ clients
3. **Advanced Testing**: Load testing, performance metrics, and monitoring
4. **Analytics**: Performance comparison dashboards and reporting

## Testing API Endpoints

Example cURL commands:
```bash
# Test health endpoint
curl http://localhost:8081/health

# Test simple request
curl http://localhost:8082/api/test/simple

# Test database query with limit
curl http://localhost:8083/api/test/database?limit=5

# Test JSON response
curl http://localhost:8084/api/test/json

# Get framework info
curl http://localhost:8085/api/info
```