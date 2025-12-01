#!/bin/bash

echo "🍌 Starting Bananas Framework Testing Service!"

# Set environment
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=bananas_user
export DB_PASSWORD=bananas_pass
export DB_NAME=bananas_dev
export DB_SSL_MODE=disable

# Check if PostgreSQL container already exists
if [ "$(docker ps -aq -f name=bananas-postgres)" ]; then
    echo "🐘 PostgreSQL container already exists, starting it..."
    docker start bananas-postgres
else
    echo "🐘 Starting PostgreSQL..."
    docker run -d \
      --name bananas-postgres \
      -p 5432:5432 \
      -e POSTGRES_DB=bananas_dev \
      -e POSTGRES_USER=bananas_user \
      -e POSTGRES_PASSWORD=bananas_pass \
      postgres:18-alpine
fi

# Wait for database to be ready
echo "⏳ Waiting for database to be ready..."
until docker exec bananas-postgres pg_isready -U bananas_user -d bananas_dev; do
  echo "Waiting for postgres..."
  sleep 2
done

# Run migrations
echo "🔧 Running database migrations..."
cd server && go run cmd/migration/main.go up

# Seed database
echo "🌱 Seeding database..."
cd server && go run cmd/migration/main.go seed

echo ""
echo "✅ Setup complete!"
echo ""
echo "💡 To start all 6 frameworks:"
echo "   make run"
echo "   OR"
echo "   cd server && go run ./cmd/api"
echo ""
echo "🚀 All frameworks will run simultaneously on:"
echo "   • Standard Library: http://localhost:8081"
echo "   • Gin:             http://localhost:8082"
echo "   • Fiber:           http://localhost:8083"
echo "   • Echo:            http://localhost:8084"
echo "   • Chi:             http://localhost:8085"
echo "   • Gorilla Mux:     http://localhost:8086"
echo ""
echo "🧪 To test all endpoints:"
echo "   ./scripts/test-endpoints.sh"
echo "   OR"
echo "   ./scripts/test-all-servers.sh"
echo ""
echo "🛑 To stop PostgreSQL:"
echo "   docker stop bananas-postgres && docker rm bananas-postgres"