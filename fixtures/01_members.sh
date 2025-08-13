#!/bin/bash

# Script to populate members via API
# Assumes the server is running on http://localhost:8080

API_BASE="http://localhost:8080/api"

echo "🔧 Creating members..."

# Member 1: Extra class, Regular type
echo "Creating member: John Smith (KD2ABC)..."
curl -X POST "$API_BASE/members" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Smith",
    "call_sign": "KD2ABC",
    "email": "john.smith@example.com",
    "phone": "(585) 555-1234",
    "mailing_address": "123 Main St, Rochester, NY 14623",
    "frn": "0012345678",
    "license_class": "Extra",
    "member_type": "Regular"
  }'
echo ""

# Member 2: General class, Senior type
echo "Creating member: Mary Johnson (W2XYZ)..."
curl -X POST "$API_BASE/members" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Mary Johnson",
    "call_sign": "W2XYZ",
    "email": "mary.johnson@example.com",
    "phone": "(585) 555-2345",
    "mailing_address": "456 Oak Ave, Rochester, NY 14620",
    "frn": "0023456789",
    "license_class": "General",
    "member_type": "Senior"
  }'
echo ""

# Member 3: Technician class, Student type
echo "Creating member: Bob Wilson (KC2DEF)..."
curl -X POST "$API_BASE/members" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Bob Wilson",
    "call_sign": "KC2DEF",
    "email": "bob.wilson@example.com",
    "phone": "(585) 555-3456",
    "mailing_address": "789 Pine St, Rochester, NY 14615",
    "frn": "0034567890",
    "license_class": "Technician",
    "member_type": "Student"
  }'
echo ""

# Member 4: Advanced class, Regular type
echo "Creating member: Alice Davis (N2GHI)..."
curl -X POST "$API_BASE/members" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Alice Davis",
    "call_sign": "N2GHI",
    "email": "alice.davis@example.com",
    "phone": "(585) 555-4567",
    "mailing_address": "321 Elm St, Rochester, NY 14607",
    "frn": "0045678901",
    "license_class": "Advanced",
    "member_type": "Regular"
  }'
echo ""

# Member 5: General class, Associate type
echo "Creating member: Charlie Brown (KB2JKL)..."
curl -X POST "$API_BASE/members" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Charlie Brown",
    "call_sign": "KB2JKL",
    "email": "charlie.brown@example.com",
    "phone": "(585) 555-5678",
    "mailing_address": "654 Maple Ave, Rochester, NY 14612",
    "frn": "0056789012",
    "license_class": "General",
    "member_type": "Associate"
  }'
echo ""

echo "✅ Members created successfully!"