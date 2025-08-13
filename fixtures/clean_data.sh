#!/bin/bash

# Script to clean all data from the database
# WARNING: This will delete all data!

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

echo "================================================"
echo "   RARA Membership System - Data Cleanup"
echo "================================================"
echo ""
echo "⚠️  WARNING: This will delete ALL data from the database!"
echo ""
read -p "Are you sure you want to continue? (yes/no): " confirmation

if [ "$confirmation" != "yes" ]; then
    echo "Cleanup cancelled."
    exit 0
fi

cd "$PROJECT_ROOT"

# Remove the database file
if [ -f "membership.db" ]; then
    echo "Removing database file..."
    rm membership.db
    echo "✅ Database removed."
else
    echo "No database file found."
fi

echo ""
echo "Database has been cleaned."
echo "Run './fixtures/run_all.sh' to reload test data."
echo ""