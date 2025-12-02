# CRUSH.md - Bananas Framework Testing Service

## Development Commands

### Build & Test
```bash
# Build all frameworks
make build

# Run tests (currently no test files exist)
make test
cd server && go test ./...

# Run single test (when test files exist)
cd server && go test -run TestSpecificFunction ./path/to/package

# Install/update dependencies
make deps

# Clean build artifacts
make clean
```

### Running Applications
```bash
# All frameworks simultaneously (ports 8081-8086)
make run
cd server && go run ./cmd/api

# Individual frameworks (alternative entry points)
# Standard Library (8081), Gin (8082), Fiber (8083), Echo (8084), Chi (8085), Gorilla Mux (8086)

# Development with hot reload
tilt up
```

### Database Operations
```bash
make migrate-up    # Run migrations
make migrate-down  # Rollback migrations
make seed          # Seed database
make db-setup      # Complete setup (create + migrate + seed)
```

## Code Style Guidelines

### Import Organization
- Group imports: standard library, third-party, internal modules
- Use blank lines between groups
- Order alphabetically within groups

### Naming Conventions
- Packages: lowercase, single word when possible (`controllers`, `services`)
- Functions: PascalCase for exported, camelCase for unexported
- Variables: camelCase, meaningful names
- Constants: UPPER_SNAKE_CASE for exported
- Framework identifiers: "standard", "gin", "fiber", "echo", "chi", "gorilla"

### Error Handling
- Always return errors from functions
- Use structured logging with `logger.Er("message", err)` for errors
- Log errors at the appropriate level
- Return appropriate HTTP status codes from controllers

### Context Usage
- Pass context through all layers for request tracing
- Add framework name to context: `context.WithValue(r.Context(), "framework", "<name>")`
- Use context for database operations and request-scoped data

### Logging
- Use custom logger: `log := logger.New("component")`
- `log.Info("message", args...)` for informational messages
- `log.Er("error message", err)` for errors
- Avoid logging sensitive data

### Database & ORM
- Repository pattern with multi-ORM support (sql, gorm, sqlx, pgx)
- Use `repositories.Manager` for ORM switching
- Always use transactions for multi-step operations

### Framework Integration
- Each framework entry point (`cmd/api/*.go`) follows same pattern
- Apply framework middleware after adding context
- Fiber requires `fiberResponseWriter` adapter
- All frameworks expose identical API endpoints

### Code Organization
- Shared logic in `internal/` (controllers, services, repositories)
- Framework-specific routing only in `cmd/api/`
- Models in `internal/models/` with JSON and DB tags
- Configuration from environment variables