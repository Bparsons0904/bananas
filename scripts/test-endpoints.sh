#!/bin/bash

echo "🍌 Testing Bananas Framework Endpoints"
echo "====================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test function
test_endpoint() {
    local name=$1
    local port=$2
    local endpoint=$3
    local expected_status=$4
    
    echo -e "\n${YELLOW}Testing $name (Port $port):$NC"
    echo "GET http://localhost:$port$endpoint"
    
    response=$(curl -s -w "%{http_code}" -o /dev/null "http://localhost:$port$endpoint" 2>/dev/null)
    
    if [ "$response" = "$expected_status" ]; then
        echo -e "${GREEN}✅ $name: OK${NC}"
        return 0
    else
        echo -e "${RED}❌ $name: Failed (Status: $response)${NC}"
        return 1
    fi
}

# Wait a moment for servers to start
echo "⏳ Waiting 3 seconds for servers to start..."
sleep 3

# Framework names and ports
declare -A frameworks=(
    ["Standard Library"]="8081"
    ["Gin"]="8082"
    ["Fiber"]="8083"
    ["Echo"]="8084"
    ["Chi"]="8085"
    ["Gorilla Mux"]="8086"
)

# Test endpoints
failed=0

echo "Testing health endpoints..."
for name in "${!frameworks[@]}"; do
    port=${frameworks[$name]}
    if ! test_endpoint "$name" "$port" "/health" "200"; then
        ((failed++))
    fi
done

echo -e "\n${YELLOW}Testing API endpoints...$NC"

# Test a few API endpoints
if ! test_endpoint "Standard Library" "8081" "/api/test/simple" "200"; then
    ((failed++))
fi

if ! test_endpoint "Gin" "8082" "/api/test/database" "200"; then
    ((failed++))
fi

if ! test_endpoint "Fiber" "8083" "/api/test/json" "200"; then
    ((failed++))
fi

if ! test_endpoint "Echo" "8084" "/api/info" "200"; then
    ((failed++))
fi

# Summary
echo -e "\n====================================="
if [ $failed -eq 0 ]; then
    echo -e "${GREEN}🎉 All tests passed! All frameworks are running correctly.${NC}"
else
    echo -e "${RED}❌ $failed tests failed. Please check the server logs.${NC}"
fi

echo -e "\n📊 Available endpoints for all frameworks:"
echo "• http://localhost:8081-8086/health"
echo "• http://localhost:8081-8086/api/test/simple"
echo "• http://localhost:8081-8086/api/test/database"  
echo "• http://localhost:8081-8086/api/test/json"
echo "• http://localhost:8081-8086/api/info"