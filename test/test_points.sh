#!/bin/bash

# Test Points API endpoints
# This script tests all point-related API endpoints

BASE_URL="${BASE_URL:-http://localhost:8080}"
API_URL="${BASE_URL}/api"

echo "Testing Points API endpoints..."
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

# Note: These tests assume that test_members.sh and test_events.sh have been run
# and that member IDs 1, 2 and event IDs 1, 2, 3 exist

# Test 1: Allocate points to member for event (using default points)
echo -e "\n${YELLOW}Test 1: Allocate points using event default${NC}"
allocation_data='{
    "event_id": 1,
    "member_id": 1,
    "notes": "Helped with setup and teardown"
}'
test_endpoint "POST" "/points/allocate" "$allocation_data" "201" "POST /api/points/allocate (default)"

# Test 2: Allocate specific points
echo -e "\n${YELLOW}Test 2: Allocate specific points${NC}"
allocation_data2='{
    "event_id": 2,
    "member_id": 1,
    "points": 25,
    "notes": "Led the entire event"
}'
test_endpoint "POST" "/points/allocate" "$allocation_data2" "201" "POST /api/points/allocate (specific)"

# Test 3: Allocate points to second member
echo -e "\n${YELLOW}Test 3: Allocate points to second member${NC}"
allocation_data3='{
    "event_id": 1,
    "member_id": 2,
    "points": 15,
    "notes": "Assisted with registration"
}'
test_endpoint "POST" "/points/allocate" "$allocation_data3" "201" "POST /api/points/allocate (member 2)"

# Test 4: Check member 1 balance
echo -e "\n${YELLOW}Test 4: Check member 1 balance${NC}"
test_endpoint "GET" "/points/balance?member_id=1" "" "200" "GET /api/points/balance"

# Test 5: Check member 2 balance
echo -e "\n${YELLOW}Test 5: Check member 2 balance${NC}"
test_endpoint "GET" "/points/balance?member_id=2" "" "200" "GET /api/points/balance (member 2)"

# Test 6: Redeem points for member 1
echo -e "\n${YELLOW}Test 6: Redeem points${NC}"
redemption_data='{
    "member_id": 1,
    "points": 20,
    "notes": "Redeemed for club t-shirt"
}'
test_endpoint "POST" "/points/redeem" "$redemption_data" "201" "POST /api/points/redeem"

# Test 7: Check updated balance after redemption
echo -e "\n${YELLOW}Test 7: Check balance after redemption${NC}"
test_endpoint "GET" "/points/balance?member_id=1" "" "200" "GET /api/points/balance (after redemption)"

# Test 8: Try to redeem more points than available (should fail)
echo -e "\n${YELLOW}Test 8: Redeem excessive points (should fail)${NC}"
excessive_redemption='{
    "member_id": 2,
    "points": 1000,
    "notes": "Trying to redeem too many points"
}'
test_endpoint "POST" "/points/redeem" "$excessive_redemption" "400" "POST /api/points/redeem (excessive)"

# Test 9: Allocate invalid points (negative)
echo -e "\n${YELLOW}Test 9: Allocate negative points (should fail)${NC}"
invalid_allocation='{
    "event_id": 1,
    "member_id": 1,
    "points": -10,
    "notes": "Invalid negative points"
}'
test_endpoint "POST" "/points/allocate" "$invalid_allocation" "400" "POST /api/points/allocate (negative)"

# Test 10: Allocate excessive points
echo -e "\n${YELLOW}Test 10: Allocate excessive points (should fail)${NC}"
excessive_allocation='{
    "event_id": 1,
    "member_id": 1,
    "points": 150,
    "notes": "Too many points"
}'
test_endpoint "POST" "/points/allocate" "$excessive_allocation" "400" "POST /api/points/allocate (excessive)"

# Test 11: Get member points history via member endpoint
echo -e "\n${YELLOW}Test 11: Get complete points history${NC}"
test_endpoint "GET" "/members/points?call_sign=W1TEST" "" "200" "GET /api/members/points (history)"

# Test 12: Allocate points for non-existent event
echo -e "\n${YELLOW}Test 12: Allocate for non-existent event${NC}"
bad_event='{
    "event_id": 99999,
    "member_id": 1,
    "points": 10
}'
test_endpoint "POST" "/points/allocate" "$bad_event" "500" "POST /api/points/allocate (bad event)"

# Test 13: Allocate points for non-existent member
echo -e "\n${YELLOW}Test 13: Allocate for non-existent member${NC}"
bad_member='{
    "event_id": 1,
    "member_id": 99999,
    "points": 10
}'
test_endpoint "POST" "/points/allocate" "$bad_member" "500" "POST /api/points/allocate (bad member)"

# Test 14: Multiple allocations for same member/event
echo -e "\n${YELLOW}Test 14: Multiple allocations same member/event${NC}"
duplicate_allocation='{
    "event_id": 3,
    "member_id": 1,
    "points": 10,
    "notes": "First allocation for event 3"
}'
test_endpoint "POST" "/points/allocate" "$duplicate_allocation" "201" "POST /api/points/allocate (first)"

duplicate_allocation2='{
    "event_id": 3,
    "member_id": 1,
    "points": 5,
    "notes": "Second allocation for same event"
}'
test_endpoint "POST" "/points/allocate" "$duplicate_allocation2" "500" "POST /api/points/allocate (second)"

echo -e "\n================================"
echo "Points API tests completed!"