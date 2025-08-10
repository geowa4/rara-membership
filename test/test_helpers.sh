#!/bin/bash

# Helper functions for API testing
# Source this file in other test scripts for reusable functions

# Colors for output
export GREEN='\033[0;32m'
export RED='\033[0;31m'
export YELLOW='\033[1;33m'
export BLUE='\033[0;34m'
export NC='\033[0m' # No Color

# Default configuration
export BASE_URL="${BASE_URL:-http://localhost:8080}"
export API_URL="${BASE_URL}/api"

# Generic test function
test_api() {
    local method=$1
    local endpoint=$2
    local data=$3
    local expected_status=$4
    local test_name=$5
    
    echo -n "Testing: $test_name... "
    
    local curl_args="-s -w \n%{http_code} -X $method"
    
    if [ ! -z "$data" ]; then
        curl_args="$curl_args -H Content-Type: application/json -d $data"
    fi
    
    response=$(curl $curl_args "$API_URL$endpoint")
    status_code=$(echo "$response" | tail -n 1)
    body=$(echo "$response" | sed '$d')
    
    if [ "$status_code" == "$expected_status" ]; then
        echo -e "${GREEN}PASSED${NC} (Status: $status_code)"
        return 0
    else
        echo -e "${RED}FAILED${NC} (Expected: $expected_status, Got: $status_code)"
        echo "  Response: $body"
        return 1
    fi
}

# Function to extract JSON field
get_json_field() {
    local json=$1
    local field=$2
    echo "$json" | jq -r ".$field" 2>/dev/null
}

# Function to create a test member
create_test_member() {
    local name=${1:-"Test User"}
    local call_sign=${2:-"TEST123"}
    local email=${3:-"test@example.com"}
    
    local data="{
        \"name\": \"$name\",
        \"call_sign\": \"$call_sign\",
        \"email\": \"$email\"
    }"
    
    response=$(curl -s -X POST \
        -H "Content-Type: application/json" \
        -d "$data" \
        "$API_URL/members/create")
    
    echo "$response"
}

# Function to create a test event
create_test_event() {
    local name=${1:-"Test Event"}
    local description=${2:-"Test event description"}
    local points=${3:-10}
    
    local data="{
        \"name\": \"$name\",
        \"description\": \"$description\",
        \"default_points\": $points
    }"
    
    response=$(curl -s -X POST \
        -H "Content-Type: application/json" \
        -d "$data" \
        "$API_URL/events/create")
    
    echo "$response"
}

# Function to allocate points
allocate_points() {
    local event_id=$1
    local member_id=$2
    local points=$3
    local notes=${4:-""}
    
    local data="{
        \"event_id\": $event_id,
        \"member_id\": $member_id"
    
    if [ ! -z "$points" ]; then
        data="$data, \"points\": $points"
    fi
    
    if [ ! -z "$notes" ]; then
        data="$data, \"notes\": \"$notes\""
    fi
    
    data="$data }"
    
    response=$(curl -s -X POST \
        -H "Content-Type: application/json" \
        -d "$data" \
        "$API_URL/points/allocate")
    
    echo "$response"
}

# Function to check member balance
check_balance() {
    local member_id=$1
    
    response=$(curl -s "$API_URL/points/balance?member_id=$member_id")
    balance=$(echo "$response" | jq -r '.balance' 2>/dev/null)
    
    if [ ! -z "$balance" ] && [ "$balance" != "null" ]; then
        echo "$balance"
    else
        echo "0"
    fi
}

# Function to wait for server
wait_for_server() {
    local max_attempts=${1:-30}
    local url="${2:-$BASE_URL/api/members}"
    
    echo -n "Waiting for server"
    for i in $(seq 1 $max_attempts); do
        if curl -s "$url" > /dev/null 2>&1; then
            echo -e " ${GREEN}OK${NC}"
            return 0
        fi
        echo -n "."
        sleep 1
    done
    
    echo -e " ${RED}TIMEOUT${NC}"
    return 1
}

# Function to print test section header
print_section() {
    local title=$1
    echo ""
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}$title${NC}"
    echo -e "${BLUE}========================================${NC}"
}

# Function to print test result summary
print_summary() {
    local passed=$1
    local failed=$2
    local total=$((passed + failed))
    
    echo ""
    echo -e "${BLUE}Test Results:${NC}"
    echo -e "  Total:  $total"
    echo -e "  Passed: ${GREEN}$passed${NC}"
    echo -e "  Failed: ${RED}$failed${NC}"
    
    if [ $failed -eq 0 ]; then
        echo -e "\n${GREEN}✓ All tests passed!${NC}"
        return 0
    else
        echo -e "\n${RED}✗ Some tests failed${NC}"
        return 1
    fi
}