# Bananas Setup Guide

Quick setup guide for the Bananas framework testing environment.

## Initial Setup

### 1. Environment Configuration

Copy the example environment file:

```bash
cp .env.example .env
```

**Important Port Note**: The default `.env.example` uses port `5433` for PostgreSQL to avoid conflicts with system PostgreSQL installations that typically use port `5432`.

If you need to customize ports or settings for your local environment:

```bash
# Create a local override file (gitignored)
cp .env .env.local

# Edit .env.local with your custom settings
# .env.local takes precedence over .env
```

### 2. Install Dependencies

**Backend (Go)**:
```bash
cd server
go mod download
```

**Frontend**:
```bash
cd frontend
make install
```

### 3. Start Services

**Option A: Use Tilt (Recommended)**
```bash
# Start everything with hot reloading
tilt up

# Access:
# - Main Frontend: http://localhost:5172
# - Tilt Dashboard: http://localhost:10350
# - Backend APIs: http://localhost:8081-8086
# - PostgreSQL: localhost:5433
```

**Option B: Manual Start**
```bash
# Terminal 1: Database
docker compose -f docker-compose.dev.yml up postgres

# Terminal 2: Backend
cd server
make run-all

# Terminal 3: Frontend
cd frontend
make dev
```

## Common Issues

### Port 5432 Already in Use

**Error**: `Bind for 0.0.0.0:5432 failed: port is already allocated`

**Solution**: You have PostgreSQL running on port 5432. Either:

1. **Use a different port** (default in `.env.example`):
   ```bash
   # In .env or .env.local
   DB_PORT=5433
   ```

2. **Stop system PostgreSQL**:
   ```bash
   sudo systemctl stop postgresql
   # or
   brew services stop postgresql
   ```

3. **Create .env.local override**:
   ```bash
   echo "DB_PORT=5433" > .env.local
   ```

### Frontend Not Loading

Make sure all frontend dependencies are installed:
```bash
cd frontend
make install
```

### Database Connection Issues

Check that the database is running and the port matches:
```bash
# Check running containers
docker ps | grep postgres

# Check .env file
cat .env | grep DB_PORT

# Test connection
psql -h localhost -p 5433 -U bananas_user -d bananas_dev
```

## Environment Variables Reference

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_PORT` | 5433 | PostgreSQL port |
| `DB_HOST` | localhost | Database host |
| `DB_USER` | bananas_user | Database username |
| `DB_PASSWORD` | bananas_pass | Database password |
| `DB_NAME` | bananas_dev | Database name |
| `TILT_PORT` | 10350 | Tilt dashboard port |
| `STANDARD_PORT` | 8081 | Standard Library server port |
| `GIN_PORT` | 8082 | Gin server port |
| `FIBER_PORT` | 8083 | Fiber server port |
| `ECHO_PORT` | 8084 | Echo server port |
| `CHI_PORT` | 8085 | Chi server port |
| `GORILLA_PORT` | 8086 | Gorilla Mux server port |

## Next Steps

Once everything is running:

1. Visit **http://localhost:5172** for the main testing interface
2. Select a frontend framework (React, Solid, or Angular)
3. Choose a backend framework and ORM to test
4. Run tests and compare performance!

## Development Workflow

- **Frontend changes**: Hot reload automatically via Vite/Angular
- **Backend changes**: Hot reload via Air (Go)
- **Database migrations**: `tilt trigger migrate-up`
- **Run tests**: `tilt trigger server-1-check-all`

See the main [README.md](README.md) for more details.
