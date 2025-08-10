#!/bin/bash

# Test Events API endpoints
# This script tests all event-related API endpoints

BASE_URL="${BASE_URL:-http://localhost:8080}"
API_URL="${BASE_URL}/api"

echo "Testing Events API endpoints..."
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
            # Extract event ID if created
            if [ "$method" == "POST" ] && [ "$status_code" == "201" ]; then
                EVENT_ID=$(echo $body | jq -r '.id' 2>/dev/null)
                if [ ! -z "$EVENT_ID" ] && [ "$EVENT_ID" != "null" ]; then
                    echo "  Created Event ID: $EVENT_ID"
                    export EVENT_ID
                fi
            fi
            echo "  Response: $(echo $body | jq -c '.' 2>/dev/null || echo $body | head -c 100)"
        fi
        return 0
    else
        echo -e "${RED}FAILED${NC} (Expected: $expected_status, Got: $status_code)"
        echo "  Response: $body"
        return 1
    fi
}

# Test 1: List all events
echo -e "\n${YELLOW}Test 1: List all events${NC}"
test_endpoint "GET" "/events" "" "200" "GET /api/events"

# Test 2: Create a new event
echo -e "\n${YELLOW}Test 2: Create a new event${NC}"
event_data='{
    "name": "Winter Field Day 2024",
    "description": "Annual winter field day event for amateur radio operators",
    "default_points": 10
}'
test_endpoint "POST" "/events/create" "$event_data" "201" "POST /api/events/create"

# Test 3: Create another event
echo -e "\n${YELLOW}Test 3: Create another event${NC}"
event2_data='{
    "name": "Summer Hamfest 2024",
    "description": "Summer hamfest and swap meet",
    "default_points": 5
}'
test_endpoint "POST" "/events/create" "$event2_data" "201" "POST /api/events/create"

# Test 4: Create event with high points value
echo -e "\n${YELLOW}Test 4: Create event with high points${NC}"
event3_data='{
    "name": "Emergency Communications Exercise",
    "description": "ARES emergency communications training",
    "default_points": 20
}'
test_endpoint "POST" "/events/create" "$event3_data" "201" "POST /api/events/create"

# Test 5: Update event (using first created event ID)
echo -e "\n${YELLOW}Test 5: Update event${NC}"
if [ ! -z "$EVENT_ID" ]; then
    update_data='{
        "name": "Winter Field Day 2024 - Updated",
        "default_points": 15
    }'
    test_endpoint "PUT" "/events/update?id=$EVENT_ID" "$update_data" "200" "PUT /api/events/update"
else
    echo -e "${YELLOW}Skipping - No event ID available${NC}"
fi

# Test 6: Update with partial data
echo -e "\n${YELLOW}Test 6: Update event with partial data${NC}"
if [ ! -z "$EVENT_ID" ]; then
    partial_update='{
        "description": "Updated description for winter field day"
    }'
    test_endpoint "PATCH" "/events/update?id=$EVENT_ID" "$partial_update" "200" "PATCH /api/events/update"
else
    echo -e "${YELLOW}Skipping - No event ID available${NC}"
fi

# Test 7: Update non-existent event
echo -e "\n${YELLOW}Test 7: Update non-existent event${NC}"
test_endpoint "PUT" "/events/update?id=99999" '{"name": "Test"}' "500" "PUT /api/events/update (not found)"

# Test 8: List events again to see updates
echo -e "\n${YELLOW}Test 8: List events after updates${NC}"
test_endpoint "GET" "/events" "" "200" "GET /api/events (after updates)"

# Test 9: Create event with invalid data
echo -e "\n${YELLOW}Test 9: Create event with invalid data${NC}"
invalid_data='{}'
test_endpoint "POST" "/events/create" "$invalid_data" "500" "POST /api/events/create (invalid)"

echo -e "\n================================"
echo "Events API tests completed!"

# Export event IDs for use in other tests
echo -e "\n${YELLOW}Note:${NC} Event IDs have been created for use in point allocation tests"