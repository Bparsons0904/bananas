# Tiltfile for Bananas development environment

# Load environment variables from .env
load('ext://dotenv', 'dotenv')
load('ext://restart_process', 'docker_build_with_restart')

dotenv('./.env')
# Load .env.local for local overrides (if it exists)
if os.path.exists('./.env.local'):
    dotenv('./.env.local')

# Configuration - use environment variables with defaults
DB_PORT = os.getenv('DB_PORT', '5432')
DOCKER_ENV = os.getenv('DOCKER_ENV', 'dev')
TILT_PORT = os.getenv('TILT_PORT', '10350')

# Development mode toggle
DEV_MODE = True

# Go Server with hot reloading
docker_build(
    'bananas-server-' + DOCKER_ENV,
    context='./server',
    dockerfile='./server/Dockerfile.dev',
    target='development',
    ignore=[
        'tmp/', 
        '*.log', 
        'main',
        '.git/',
        'Dockerfile*',
        '.dockerignore',
        'data/',
        '*.db',
        '*.db-journal',
    ]
)

# PostgreSQL database service
docker_build(
    'bananas-postgres-' + DOCKER_ENV,
    context='.',
    dockerfile_contents="""
FROM postgres:18
# Minimal Dockerfile - PostgreSQL image already configured
"""
)

# Use docker-compose for orchestration
docker_compose('./docker-compose.' + DOCKER_ENV + '.yml')

# ==========================================
# CORE SERVICES
# ==========================================

dc_resource('postgres',
    labels=['1-services'],
    resource_deps=[],
)

dc_resource('server',
    labels=['1-services'],
    resource_deps=['postgres'],
)

# Development utilities
if DEV_MODE:
    # ==========================================
    # SERVER/BACKEND QUALITY CHECKS
    # ==========================================
    
    # Server full check - runs tests
    local_resource(
        'server-1-check-all',
        cmd='cd server && go test ./...',
        deps=['./server'],
        ignore=['./server/tmp', './server/*.log', './server/main'],
        labels=['2-server'],
        auto_init=False,
        trigger_mode=TRIGGER_MODE_MANUAL
    )

    # Build tests for individual frameworks
    local_resource(
        'server-2-build-all',
        cmd='''
        cd server && echo "Building all framework servers..."
        cd server && go build -o tmp/bananas ./cmd/api/main.go
        echo "✅ All frameworks built successfully!"
        ''',
        deps=['./server'],
        ignore=['./server/tmp', './server/*.log', './server/main'],
        labels=['2-server'],
        auto_init=False,
        trigger_mode=TRIGGER_MODE_MANUAL
    )

    # ==========================================
    # DATABASE UTILITIES
    # ==========================================

    # Database migration commands
    local_resource(
        'migrate-up',
        cmd='cd server && go run cmd/migration/main.go up',
        deps=['./server/cmd/migration', './server/internal', './server/config'],
        ignore=['./server/tmp', './server/*.log', './server/main'],
        labels=['3-database'],
        auto_init=False,
        trigger_mode=TRIGGER_MODE_MANUAL,
        resource_deps=['server']
    )

    local_resource(
        'migrate-down',
        cmd='cd server && go run cmd/migration/main.go down',
        deps=['./server/cmd/migration', './server/internal', './server/config'],
        ignore=['./server/tmp', './server/*.log', './server/main'],
        labels=['3-database'],
        auto_init=False,
        trigger_mode=TRIGGER_MODE_MANUAL,
        resource_deps=['server']
    )

    local_resource(
        'migrate-seed',
        cmd='cd server && go run cmd/migration/main.go seed',
        deps=['./server/cmd/migration', './server/internal', './server/config'],
        ignore=['./server/tmp', './server/*.log', './server/main'],
        labels=['3-database'],
        auto_init=False,
        trigger_mode=TRIGGER_MODE_MANUAL,
        resource_deps=['server']
    )

    # PostgreSQL utilities
    local_resource(
        'postgres-info',
        cmd='docker compose -f docker-compose.' + DOCKER_ENV + '.yml exec postgres psql -U bananas_user -d bananas_dev -c "\\l"',
        labels=['3-database'],
        auto_init=False,
        trigger_mode=TRIGGER_MODE_MANUAL,
        resource_deps=['postgres']
    )

    # ==========================================
    # FRONTEND SERVICES
    # ==========================================
    # Note: Frontend services have NO dependencies on backend/DB
    # They will start in parallel with everything else

    # Main frontend launcher (framework selector)
    local_resource(
        'frontend-main',
        serve_cmd='cd frontend && npx serve . -p 5172 --cors',
        deps=['./frontend/index.html', './frontend/package.json'],
        labels=['4-frontend'],
        links=['http://localhost:5172'],
    )

    # React frontend
    local_resource(
        'frontend-react',
        serve_cmd='cd frontend/react && npm run dev',
        deps=['./frontend/react/src', './frontend/react/package.json'],
        labels=['4-frontend'],
        links=['http://localhost:5173'],
    )

    # Solid frontend
    local_resource(
        'frontend-solid',
        serve_cmd='cd frontend/solid && npm run dev',
        deps=['./frontend/solid/src', './frontend/solid/package.json'],
        labels=['4-frontend'],
        links=['http://localhost:5174'],
    )

    # Angular frontend
    local_resource(
        'frontend-angular',
        serve_cmd='cd frontend/angular && npm start',
        deps=['./frontend/angular/src', './frontend/angular/package.json'],
        labels=['4-frontend'],
        links=['http://localhost:5175'],
    )

print("🚀 Bananas Development Environment (Environment: %s)" % DOCKER_ENV)
print("📊 Tilt Dashboard: http://localhost:%s" % TILT_PORT)
print("🔧 Server APIs: http://localhost:8081-8086")
print("🐘 PostgreSQL: localhost:%s" % DB_PORT)
print("💡 Hot reloading enabled for all services!")
print("🧪 Manual test/migration resources available in Tilt UI")
print("\n📋 Backend - All Frameworks Running Simultaneously:")
print("• Standard Library: http://localhost:8081")
print("• Gin:             http://localhost:8082")
print("• Fiber:           http://localhost:8083")
print("• Echo:            http://localhost:8084")
print("• Chi:             http://localhost:8085")
print("• Gorilla Mux:     http://localhost:8086")

print("\n📋 Frontend - Testing Clients:")
print("• 🌟 MAIN:         http://localhost:5172  (Framework Switcher)")
print("• React:           http://localhost:5173")
print("• Solid:           http://localhost:5174")
print("• Angular:         http://localhost:5175")

print("\n📋 Quick Commands:")
print("\n🔧 SERVER (Backend):")
print("• tilt trigger server-1-check-all     - Run server tests")
print("• tilt trigger server-2-build-all     - Build all frameworks")
print("\n💾 DATABASE:")
print("• tilt trigger migrate-up             - Run database migrations")
print("• tilt trigger migrate-down           - Rollback migrations")
print("• tilt trigger migrate-seed           - Seed database")
print("• tilt trigger postgres-info          - Show PostgreSQL info")
print("\n⚡ GENERAL:")
print("• tilt down                           - Stop all services")
print("• tilt up --stream                    - Start with streaming logs")
print("\n🚀 LOCAL DEVELOPMENT:")
print("• make run-all                        - Start all frameworks locally")
print("• make run-standard                    - Test Standard Library only")
print("• make run-gin                        - Test Gin only")
print("• make run-fiber                      - Test Fiber only")
print("• make run-echo                       - Test Echo only")
print("• make run-chi                        - Test Chi only")
print("• make run-gorilla                    - Test Gorilla Mux only")