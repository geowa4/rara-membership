#!/bin/bash

# Main script to run all fixtures and display the created data
# This script assumes the server is already running on port 8080

set -e  # Exit on any error

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

echo "================================================"
echo "   RARA Membership System - Fixture Loader"
echo "================================================"
echo ""

# Check if server is running
echo "🔍 Checking if server is running on localhost:8080..."
if ! curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/api/members | grep -q "200"; then
    echo "❌ Server is not running!"
    echo "Please start the server first with: go run main.go serve"
    exit 1
fi
echo "✅ Server is running!"
echo ""

# Make all scripts executable
chmod +x "$SCRIPT_DIR"/*.sh

# Run fixture scripts in order
echo "================================================"
echo "Loading fixture data..."
echo "================================================"
echo ""

# Create members
echo "Step 1: Creating Members"
echo "------------------------"
bash "$SCRIPT_DIR/01_members.sh"
echo ""
sleep 1

# Create events
echo "Step 2: Creating Events"
echo "-----------------------"
bash "$SCRIPT_DIR/02_events.sh"
echo ""
sleep 1

# Create point allocations
echo "Step 3: Creating Point Allocations"
echo "-----------------------------------"
bash "$SCRIPT_DIR/03_point_allocations.sh"
echo ""
sleep 1

# Create point redemptions
echo "Step 4: Creating Point Redemptions"
echo "-----------------------------------"
bash "$SCRIPT_DIR/04_point_redemptions.sh"
echo ""
sleep 1

echo ""
echo "================================================"
echo "   Viewing Created Data with CLI"
echo "================================================"
echo ""

# Change to project root to run CLI commands
cd "$PROJECT_ROOT"

# Build the CLI if needed
if [ ! -f "./rara-membership" ]; then
    echo "Building CLI..."
    go build -o rara-membership .
fi

# Display members
echo "📻 Active Members:"
echo "=================="
./rara-membership member list
echo ""

# Display events
echo "📅 Upcoming Events:"
echo "==================="
./rara-membership event list
echo ""

# Display points for each member using member points command
echo "🎯 Member Points:"
echo "================="
for call_sign in KD2ABC W2XYZ KC2DEF N2GHI KB2JKL; do
    echo ""
    echo "Points for $call_sign:"
    echo "-------------------"
    ./rara-membership member points "$call_sign" 2>/dev/null || echo "Unable to view points for $call_sign"
done

echo ""
echo "================================================"
echo "   Fixture Loading Complete!"
echo "================================================"
echo ""
echo "✅ All test data has been created successfully!"
echo ""
echo "You can now:"
echo "  - View the web UI at http://localhost:8080"
echo "  - Use the CLI: ./rara-membership --help"
echo "  - Check the database: membership.db"
echo ""