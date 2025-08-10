#!/bin/bash

# Master test script for RARA Membership API
# This script runs all API tests with a fresh database

set -e  # Exit on error

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
BASE_URL="${BASE_URL:-http://localhost:8080}"
TEST_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$TEST_DIR")"
# Use temporary directory for test database
TEST_DB_DIR=$(mktemp -d)
DB_FILE="$TEST_DB_DIR/membership.db"
BACKUP_DB_FILE="$TEST_DB_DIR/membership.db.backup"
SERVER_PID=""

echo -e "${BLUE}================================================${NC}"
echo -e "${BLUE}    RARA Membership API Test Suite${NC}"
echo -e "${BLUE}================================================${NC}"
echo ""

# Function to cleanup on exit
cleanup() {
    echo -e "\n${YELLOW}Cleaning up...${NC}"
    
    # Stop the server if it's running
    if [ ! -z "$SERVER_PID" ]; then
        echo "Stopping server (PID: $SERVER_PID)..."
        kill $SERVER_PID 2>/dev/null || true
        wait $SERVER_PID 2>/dev/null || true
    fi
    
    # Clean up temporary database directory
    if [ -d "$TEST_DB_DIR" ]; then
        echo "Removing temporary database directory..."
        rm -rf "$TEST_DB_DIR"
    fi
    
    echo -e "${GREEN}Cleanup complete${NC}"
}

# Set up trap to cleanup on exit
trap cleanup EXIT INT TERM

# Function to start the server
start_server() {
    echo -e "${YELLOW}Starting server...${NC}"
    
    # Build the application first
    echo "Building application..."
    cd "$PROJECT_ROOT"
    go build -o rara-membership . || {
        echo -e "${RED}Failed to build application${NC}"
        exit 1
    }
    
    # Start the server in the background with test database path
    RARA_DB_PATH="$DB_FILE" ./rara-membership serve > /tmp/rara-test-server.log 2>&1 &
    SERVER_PID=$!
    
    # Wait for server to start
    echo -n "Waiting for server to start"
    for i in {1..30}; do
        if curl -s "$BASE_URL/api/members" > /dev/null 2>&1; then
            echo -e " ${GREEN}OK${NC}"
            echo "Server started successfully (PID: $SERVER_PID)"
            return 0
        fi
        echo -n "."
        sleep 1
    done
    
    echo -e " ${RED}FAILED${NC}"
    echo "Server failed to start. Check /tmp/rara-test-server.log for details"
    tail -20 /tmp/rara-test-server.log
    exit 1
}

# Function to reset database
reset_database() {
    echo -e "${YELLOW}Setting up test database...${NC}"
    
    # Backup existing database if it exists
    if [ -f "$DB_FILE" ]; then
        echo "Backing up existing database..."
        cp "$DB_FILE" "$BACKUP_DB_FILE"
    fi
    
    # Remove existing database to start fresh
    rm -f "$DB_FILE"
    echo "Created fresh database"
}

# Function to run a test script
run_test() {
    local test_script=$1
    local test_name=$2
    
    echo ""
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}Running: $test_name${NC}"
    echo -e "${BLUE}========================================${NC}"
    
    if [ -f "$test_script" ]; then
        bash "$test_script"
        local exit_code=$?
        if [ $exit_code -eq 0 ]; then
            echo -e "${GREEN}✓ $test_name completed successfully${NC}"
        else
            echo -e "${RED}✗ $test_name failed with exit code $exit_code${NC}"
            return $exit_code
        fi
    else
        echo -e "${RED}Test script not found: $test_script${NC}"
        return 1
    fi
}

# Function to check dependencies
check_dependencies() {
    echo -e "${YELLOW}Checking dependencies...${NC}"
    
    # Check for required commands
    local missing_deps=()
    
    command -v curl >/dev/null 2>&1 || missing_deps+=("curl")
    command -v jq >/dev/null 2>&1 || missing_deps+=("jq")
    command -v go >/dev/null 2>&1 || missing_deps+=("go")
    
    if [ ${#missing_deps[@]} -gt 0 ]; then
        echo -e "${RED}Missing required dependencies: ${missing_deps[*]}${NC}"
        echo "Please install the missing dependencies and try again."
        exit 1
    fi
    
    echo -e "${GREEN}All dependencies satisfied${NC}"
}

# Function to show test summary
show_summary() {
    echo ""
    echo -e "${BLUE}================================================${NC}"
    echo -e "${BLUE}           Test Summary${NC}"
    echo -e "${BLUE}================================================${NC}"
    
    # Check server log for any errors
    if [ -f /tmp/rara-test-server.log ]; then
        local error_count=$(grep -c "ERROR\|FATAL\|panic" /tmp/rara-test-server.log 2>/dev/null || echo "0")
        error_count=$(echo "$error_count" | tr -d '\n\r' | head -c 10)
        if [ "$error_count" -gt 0 ]; then
            echo -e "${YELLOW}Warning: Found $error_count error(s) in server log${NC}"
            echo "Last 10 lines of server log:"
            tail -10 /tmp/rara-test-server.log
        fi
    fi
    
    # Show database statistics
    if [ -f "$DB_FILE" ]; then
        local db_size=$(du -h "$DB_FILE" | cut -f1)
        echo -e "Database size: ${db_size}"
    fi
    
    echo ""
    echo -e "${GREEN}All tests completed!${NC}"
}

# Main execution
main() {
    echo "Test configuration:"
    echo "  Base URL: $BASE_URL"
    echo "  Database: $DB_FILE"
    echo "  Test Dir: $TEST_DIR"
    echo ""
    
    # Check dependencies
    check_dependencies
    
    # Reset database
    reset_database
    
    # Start the server
    start_server
    
    # Make test scripts executable
    chmod +x "$TEST_DIR"/*.sh 2>/dev/null || true
    
    # Run tests in order
    echo -e "\n${YELLOW}Starting API tests...${NC}"
    
    # Run member tests first to create test members
    run_test "$TEST_DIR/test_members.sh" "Members API Tests" || exit 1
    
    # Run event tests to create test events
    run_test "$TEST_DIR/test_events.sh" "Events API Tests" || exit 1
    
    # Run points tests (depends on members and events)
    run_test "$TEST_DIR/test_points.sh" "Points API Tests" || exit 1
    
    # Show summary
    show_summary
}

# Run if not sourced
if [ "${BASH_SOURCE[0]}" = "${0}" ]; then
    main "$@"
fi