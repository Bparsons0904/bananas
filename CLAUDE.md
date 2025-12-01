# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Overview

Bananas is a framework comparison testing service that runs **6 different Go web frameworks simultaneously** with **shared business logic**. This unique architecture allows fair performance comparison by ensuring all frameworks use identical controllers, services, repositories, and database operations.

**Key Design Principle**: Framework-specific code is isolated to routing/middleware (`cmd/api/*.go`), while all business logic is shared (`internal/`).

## Common Commands

### Development Environment

```bash
# Start all services with hot reloading (Tilt + Docker)
tilt up

# Start all frameworks locally (without Docker)
make run-all

# Stop Tilt services
tilt down
```

### Running Individual Frameworks

All 6 frameworks can run simultaneously from a single Go application:

```bash
make run-all        # All frameworks (ports 8081-8086)
make run-standard   # Standard Library only (port 8081)
make run-gin        # Gin only (port 8082)
make run-fiber      # Fiber only (port 8083)
make run-echo       # Echo only (port 8084)
make run-chi        # Chi only (port 8085)
make run-gorilla    # Gorilla Mux only (port 8086)
```

### Testing

```bash
# Run Go tests
make test
cd server && go test ./...

# Run via Tilt
tilt trigger server-1-check-all

# Test all API endpoints
./scripts/test-endpoints.sh
```

### Database Operations

```bash
# Via Make
make migrate-up     # Run migrations
make migrate-down   # Rollback migrations
make seed           # Seed database

# Via Tilt (recommended for dev)
tilt trigger migrate-up
tilt trigger migrate-down
tilt trigger migrate-seed
tilt trigger postgres-info
```

### Build Commands

```bash
make build    # Build all framework binaries
make deps     # Install/update dependencies
make clean    # Clean build artifacts
```

## Architecture

### Multi-Framework Design

The application uses a **shared application structure** with **multi-ORM support**:

```
┌─────────────────────────────────────────────────────────┐
│  6 Framework Entry Points (cmd/api/*.go)               │
│  Standard│Gin│Fiber│Echo│Chi│Gorilla                  │
└─────────────────────────────────────────────────────────┘
                         ↓
┌─────────────────────────────────────────────────────────┐
│  Shared Application Layer (internal/)                  │
│  Controllers → Services → Repository Manager           │
└─────────────────────────────────────────────────────────┘
                         ↓
┌─────────────────────────────────────────────────────────┐
│  Multi-ORM Repository Layer                            │
│  database/sql │ GORM │ SQLx │ PGX                      │
└─────────────────────────────────────────────────────────┘
                         ↓
┌─────────────────────────────────────────────────────────┐
│  PostgreSQL Database                                    │
└─────────────────────────────────────────────────────────┘
```

### Application Initialization Flow

The `internal/app/app.go` follows a dependency injection pattern:

1. Load configuration from environment variables
2. Initialize database connection (`internal/database/database.go`)
3. Create repository manager with all 4 ORMs (`internal/repositories/manager.go`)
4. Create services with repository manager (`internal/services/service.go`)
5. Create controllers with services (`internal/controllers/base_controller.go`)

**All frameworks share the same App instance** to ensure identical behavior.

### Multi-ORM Repository System

The `repositories.Manager` (`internal/repositories/manager.go`) provides access to 4 different ORM implementations:

- `sql` - Standard library `database/sql`
- `gorm` - GORM ORM
- `sqlx` - jmoiron/sqlx
- `pgx` - jackc/pgx/v5

Services can switch between ORMs dynamically for performance testing.

### Framework Entry Point Pattern

Each framework entry point (`cmd/api/*.go`) must:

1. Initialize the shared App
2. Set up framework-specific router
3. **Add framework name to context**: `context.WithValue(r.Context(), "framework", "<name>")`
4. Register shared controller methods as handlers
5. Start server on assigned port (8081-8086)

**Framework identifiers**: `"standard"`, `"gin"`, `"fiber"`, `"echo"`, `"chi"`, `"gorilla"`

### API Endpoints

All frameworks expose identical REST endpoints:

- `GET /health` - Health check
- `GET /api/test/simple` - Simple request test
- `GET /api/test/database?limit=N` - Database query test (supports `?orm=<type>` parameter)
- `GET /api/test/json` - JSON response test
- `GET /api/info` - Framework information

### Important Code Patterns

**Logging** - Use the custom logger (`internal/logger/logger.go`):
```go
log := logger.New("component")
log.Info("message", arg1, arg2)
log.Er("error message", err)
```

**Context Propagation** - Always pass context through all layers for framework identification and database operations.

**Error Handling** - Return errors from functions and log them with structured logging.

**Fiber Compatibility** - Fiber requires a special adapter (`fiberResponseWriter`) to work with standard `http.Handler` interfaces. This is already implemented in `cmd/api/fiber.go`.

## Development Workflow

1. **Make changes** to shared code (`internal/`) or framework-specific code (`cmd/api/`)
2. **Hot reloading** via Air automatically restarts affected services
3. **Test endpoints** with curl or the test script:
   ```bash
   curl http://localhost:8081/api/test/simple
   curl http://localhost:8082/api/test/database?limit=5&orm=gorm
   ```
4. **Run tests** via `make test` or `tilt trigger server-1-check-all`

## Project Structure

```
server/
├── cmd/
│   ├── api/
│   │   ├── all-frameworks.go   # Runs all 6 frameworks simultaneously
│   │   ├── main.go             # Standard Library
│   │   ├── gin.go              # Gin
│   │   ├── fiber.go            # Fiber
│   │   ├── echo.go             # Echo
│   │   ├── chi.go              # Chi
│   │   └── gorilla.go          # Gorilla Mux
│   └── migration/              # Database migration tool
└── internal/                    # Shared application code
    ├── app/                    # Application initialization
    ├── config/                 # Configuration management
    ├── controllers/            # Shared controllers
    ├── database/               # Database connection
    ├── logger/                 # Logging utilities
    ├── models/                 # Data models
    ├── repositories/           # Multi-ORM data access layer
    │   ├── manager.go          # ORM manager
    │   ├── repository.go       # database/sql implementation
    │   ├── gorm_repository.go  # GORM implementation
    │   ├── sqlx_repository.go  # SQLx implementation
    │   └── pgx_repository.go   # PGX implementation
    └── services/               # Business logic
```

## Environment Configuration

All frameworks use fixed ports defined in `.env`:

- **Standard Library**: 8081
- **Gin**: 8082
- **Fiber**: 8083
- **Echo**: 8084
- **Chi**: 8085
- **Gorilla Mux**: 8086

Database connection uses PostgreSQL with credentials from environment variables.

## Testing ORM Performance

To test different ORMs, use the `orm` query parameter:

```bash
curl "http://localhost:8081/api/test/database?limit=10&orm=sql"
curl "http://localhost:8082/api/test/database?limit=10&orm=gorm"
curl "http://localhost:8083/api/test/database?limit=10&orm=sqlx"
curl "http://localhost:8084/api/test/database?limit=10&orm=pgx"
```

## Current Development Status

- ✅ **Phase 1 Complete**: All 6 Go frameworks running with shared business logic
- ✅ **Phase 2 Complete**: Multi-ORM support (database/sql, GORM, SQLx, PGX)
- 📋 **Phase 3 Planned**: Frontend clients (React, Vue, Svelte, Solid, Angular, HTMX, Templ)
- 📋 **Phase 4 Planned**: Load testing and performance analytics

See `PROJECT_PLAN.md` for full roadmap.
