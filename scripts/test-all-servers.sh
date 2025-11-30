#!/bin/bash

set -e

TIMEOUT=30
HEALTH_CHECK_WAIT=5

FRAMEWORKS=(
  "8081:Standard Library"
  "8082:Gin"
  "8083:Fiber"
  "8084:Echo"
  "8085:Chi"
  "8086:Gorilla Mux"
)

ORMS=("sql" "gorm" "sqlx" "pgx")

ENDPOINTS=(
  "/health"
  "/api/test/simple"
  "/api/test/json"
  "/api/info"
)

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

PASSED=0
FAILED=0

log_info() {
  echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
  echo -e "${GREEN}[PASS]${NC} $1"
  ((PASSED++))
}

log_error() {
  echo -e "${RED}[FAIL]${NC} $1"
  ((FAILED++))
}

log_warning() {
  echo -e "${YELLOW}[WARN]${NC} $1"
}

wait_for_servers() {
  log_info "Waiting ${HEALTH_CHECK_WAIT}s for servers to be ready..."
  sleep $HEALTH_CHECK_WAIT
}

test_endpoint() {
  local url=$1
  local expected_status=${2:-200}
  local description=$3

  response=$(curl -s -w "\n%{http_code}" -m 5 "$url" 2>/dev/null) || {
    log_error "$description - Connection failed"
    return 1
  }

  http_code=$(echo "$response" | tail -n1)
  body=$(echo "$response" | head -n-1)

  if [ "$http_code" == "$expected_status" ]; then
    log_success "$description"
    return 0
  else
    log_error "$description - Expected $expected_status, got $http_code"
    return 1
  fi
}

test_database_with_orm() {
  local port=$1
  local framework=$2
  local orm=$3

  url="http://localhost:${port}/api/test/database?orm=${orm}&limit=5"
  description="$framework - Database query with ORM: $orm"

  response=$(curl -s -w "\n%{http_code}" -m 5 "$url" 2>/dev/null) || {
    log_error "$description - Connection failed"
    return 1
  }

  http_code=$(echo "$response" | tail -n1)
  body=$(echo "$response" | head -n-1)

  if [ "$http_code" == "200" ]; then
    if echo "$body" | grep -q "\"orm\":\"$orm\""; then
      log_success "$description - ORM verified in response"
      return 0
    else
      log_error "$description - ORM not found in response"
      return 1
    fi
  else
    log_error "$description - HTTP $http_code"
    return 1
  fi
}

echo ""
log_info "========================================="
log_info "Bananas Multi-Server Test Suite"
log_info "========================================="
echo ""

wait_for_servers

log_info "Phase 1: Health Checks"
log_info "---------------------------------------"
for fw in "${FRAMEWORKS[@]}"; do
  port=$(echo $fw | cut -d: -f1)
  name=$(echo $fw | cut -d: -f2)
  test_endpoint "http://localhost:${port}/health" 200 "$name - Health check"
done

echo ""
log_info "Phase 2: Basic Endpoints (All Frameworks)"
log_info "---------------------------------------"
for fw in "${FRAMEWORKS[@]}"; do
  port=$(echo $fw | cut -d: -f1)
  name=$(echo $fw | cut -d: -f2)

  test_endpoint "http://localhost:${port}/api/test/simple" 200 "$name - Simple request"
  test_endpoint "http://localhost:${port}/api/test/json" 200 "$name - JSON response"
  test_endpoint "http://localhost:${port}/api/info" 200 "$name - Framework info"
done

echo ""
log_info "Phase 3: Database Queries (All ORMs × All Frameworks)"
log_info "---------------------------------------"
for fw in "${FRAMEWORKS[@]}"; do
  port=$(echo $fw | cut -d: -f1)
  name=$(echo $fw | cut -d: -f2)

  for orm in "${ORMS[@]}"; do
    test_database_with_orm "$port" "$name" "$orm"
  done
done

echo ""
log_info "========================================="
log_info "Test Results Summary"
log_info "========================================="
echo -e "${GREEN}Passed:${NC} $PASSED"
echo -e "${RED}Failed:${NC} $FAILED"
echo -e "Total:  $((PASSED + FAILED))"
echo ""

if [ $FAILED -eq 0 ]; then
  log_success "All tests passed! 🎉"
  exit 0
else
  log_error "$FAILED test(s) failed"
  exit 1
fi
