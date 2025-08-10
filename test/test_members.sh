#!/bin/bash

# Test Members API endpoints
# This script tests all member-related API endpoints

BASE_URL="${BASE_URL:-http://localhost:8080}"
API_URL="${BASE_URL}/api"

echo "Testing Members API endpoints..."
echo "================================"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test function
test_endpoint() {
    local method=$1
    local endpoint=$2
    local data=$3
    local expected_status=$4
    local test_name=$5
    
    echo -n "Testing: $test_name... "
    
    if [ -z "$data" ]; then
        response=$(curl -s -w "\n%{http_code}" -X $method "$API_URL$endpoint")
    else
        response=$(curl -s -w "\n%{http_code}" -X $method \
            -H "Content-Type: application/json" \
            -d "$data" \
            "$API_URL$endpoint")
    fi
    
    status_code=$(echo "$response" | tail -n 1)
    body=$(echo "$response" | sed '$d')
    
    if [ "$status_code" == "$expected_status" ]; then
        echo -e "${GREEN}PASSED${NC} (Status: $status_code)"
        if [ ! -z "$body" ] && [ "$body" != "" ]; then
            echo "  Response: $(echo $body | jq -c '.' 2>/dev/null || echo $body | head -c 100)"
        fi
        return 0
    else
        echo -e "${RED}FAILED${NC} (Expected: $expected_status, Got: $status_code)"
        echo "  Response: $body"
        return 1
    fi
}

# Test 1: List all members
echo -e "\n${YELLOW}Test 1: List all members${NC}"
test_endpoint "GET" "/members" "" "200" "GET /api/members"

# Test 2: Create a new member
echo -e "\n${YELLOW}Test 2: Create a new member${NC}"
member_data='{
    "name": "John Doe",
    "call_sign": "W1TEST",
    "email": "john.doe@example.com",
    "phone": "555-1234",
    "mailing_address": "123 Main St, Anytown, USA",
    "frn": "0012345678",
    "license_class": "General"
}'
test_endpoint "POST" "/members" "$member_data" "201" "POST /api/members"

# Test 3: Create duplicate member (should fail)
echo -e "\n${YELLOW}Test 3: Create duplicate member (should fail)${NC}"
test_endpoint "POST" "/members" "$member_data" "409" "POST /api/members (duplicate)"

# Test 4: Create member with missing required fields (should fail)
echo -e "\n${YELLOW}Test 4: Create member with missing required fields${NC}"
invalid_data='{
    "name": "Jane Doe"
}'
test_endpoint "POST" "/members" "$invalid_data" "400" "POST /api/members (invalid)"

# Test 5: Update member
echo -e "\n${YELLOW}Test 5: Update member${NC}"
update_data='{
    "phone": "555-5678",
    "license_class": "Extra"
}'
test_endpoint "PUT" "/members/W1TEST" "$update_data" "200" "PUT /api/members/{call_sign}"

# Test 6: Update non-existent member
echo -e "\n${YELLOW}Test 6: Update non-existent member${NC}"
test_endpoint "PUT" "/members/NOTEXIST" "$update_data" "404" "PUT /api/members/{call_sign} (not found)"

# Test 7: Get member points history
echo -e "\n${YELLOW}Test 7: Get member points history${NC}"
test_endpoint "GET" "/members/W1TEST/points" "" "200" "GET /api/members/{call_sign}/points"

# Test 8: Get points for non-existent member
echo -e "\n${YELLOW}Test 8: Get points for non-existent member${NC}"
test_endpoint "GET" "/members/NOTEXIST/points" "" "404" "GET /api/members/{call_sign}/points (not found)"

# Test 9: Create another member for testing
echo -e "\n${YELLOW}Test 9: Create another member${NC}"
member2_data='{
    "name": "Jane Smith",
    "call_sign": "KD2ABC",
    "email": "jane.smith@example.com",
    "phone": "555-9999",
    "license_class": "Technician"
}'
test_endpoint "POST" "/members" "$member2_data" "201" "POST /api/members (member 2)"

echo -e "\n================================"
echo "Members API tests completed!"