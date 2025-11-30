# Bananas Framework Testing Service

A comprehensive testing service for comparing performance and characteristics of different Go web frameworks and frontend frameworks.

## Architecture

The project consists of:
- **Backend**: Multiple Go web frameworks (Standard Library, Gin, Fiber, Echo, Chi, Gorilla Mux) sharing the same controllers and business logic
- **Frontend**: Multiple client frameworks (React, Vue, Svelte, Solid, Angular, HTMX, Templ) for testing different approaches
- **Database**: PostgreSQL with support for multiple ORMs
- **Development**: Tilt for local development with hot reloading

## Quick Start

1. Copy environment variables:
   ```bash
   cp .env.example .env
   ```

2. Start the development environment:
   ```bash
   ./scripts/setup-dev.sh
   ```

3. Run all frameworks simultaneously:
   ```bash
   make run-all
   ```

## Backend Frameworks

All backend frameworks share:
- Same controllers (`internal/controllers/base_controller.go`)
- Same services (`internal/services/service.go`)
- Same repositories (`internal/repositories/repository.go`)
- Same database layer (`internal/database/database.go`)

### Available Endpoints

Each framework exposes the same API endpoints:

- `GET /health` - Health check
- `GET /api/test/simple` - Simple request test
- `GET /api/test/database` - Database query test
- `GET /api/test/json` - JSON response test
- `GET /api/info` - Framework information

### Framework Ports (All Running Simultaneously)

All frameworks can run simultaneously from a single Go application:

- **Standard Library**: http://localhost:8081
- **Gin**: http://localhost:8082
- **Fiber**: http://localhost:8083
- **Echo**: http://localhost:8084
- **Chi**: http://localhost:8085
- **Gorilla Mux**: http://localhost:8086

### Running Frameworks

**All Frameworks Together:**
```bash
make run-all  # Starts all 6 frameworks simultaneously
```

**Individual Frameworks (for testing):**
```bash
make run-standard   # Standard Library only (port 8081)
make run-gin        # Gin only (port 8082)
make run-fiber      # Fiber only (port 8083)
make run-echo       # Echo only (port 8084)
make run-chi        # Chi only (port 8085)
make run-gorilla    # Gorilla Mux only (port 8086)
```

**Docker Compose:**
```bash
# Runs all frameworks in one container
docker-compose -f docker-compose.dev.yml up server
```

## Frontend Frameworks

Frontend clients will be implemented in the `clients/` directory:
- React
- Vue
- Svelte
- Solid
- Angular
- HTMX
- Templ

## Project Structure

```
bananas/
├── server/                 # Go backend code
│   ├── cmd/               # Main entry points for each framework
│   │   └── api/
│   │       ├── main.go    # Standard library
│   │       ├── gin.go     # Gin framework
│   │       ├── fiber.go   # Fiber framework
│   │       ├── echo.go    # Echo framework
│   │       ├── chi.go     # Chi framework
│   │       └── gorilla.go # Gorilla Mux framework
│   ├── internal/          # Internal application code
│   │   ├── app/          # Application structure
│   │   ├── config/       # Configuration
│   │   ├── controllers/  # Shared controllers
│   │   ├── database/     # Database layer
│   │   ├── models/       # Data models
│   │   ├── repositories/ # Data access layer
│   │   ├── services/     # Business logic
│   │   └── logger/       # Logging utilities
│   └── cmd/migration/    # Database migration tool
├── clients/              # Frontend frameworks
│   ├── react/
│   ├── vue/
│   ├── svelte/
│   ├── solid/
│   ├── angular/
│   ├── htmx/
│   └── templ/
├── database/             # Database configuration
└── scripts/             # Utility scripts
```

## Development

### Running Tests

```bash
# Run all server tests
tilt trigger server-2-tests

# Run all server checks (tests + linting when configured)
tilt trigger server-1-check-all
```

### Database Operations

```bash
# Run migrations
tilt trigger migrate-up

# Rollback migrations
tilt trigger migrate-down

# Seed database
tilt trigger migrate-seed

# Check database info
tilt trigger postgres-info
```

### Hot Reloading

All backend services have hot reloading enabled via Air. Changes to Go code will automatically rebuild and restart the respective service.

## Contributing

1. Follow the existing code structure
2. Each framework should maintain the same API contract
3. All shared logic should be in the `internal/` directory
4. Framework-specific code should be in `cmd/api/` files

## Testing API Endpoints

You can test the endpoints directly:

```bash
# Test simple request
curl http://localhost:8081/api/test/simple

# Test database query
curl http://localhost:8082/api/test/database?limit=5

# Test JSON response
curl http://localhost:8083/api/test/json

# Get framework info
curl http://localhost:8084/api/info
```