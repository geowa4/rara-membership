#!/bin/bash

# Script to populate events via API
# Assumes the server is running on http://localhost:8080

API_BASE="http://localhost:8080/api"

echo "📅 Creating events..."

# Event 1: Monthly Meeting
echo "Creating event: Monthly Club Meeting..."
curl -X POST "$API_BASE/events" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "RARA Monthly Meeting - January",
    "description": "Regular monthly meeting for all club members. Topics include new member introductions and technical presentations.",
    "date": "2025-01-15T19:00",
    "timezone": "America/New_York",
    "points": 5,
    "latitude": 43.1547,
    "longitude": -77.6158
  }'
echo ""

# Event 2: Field Day
echo "Creating event: Winter Field Day..."
curl -X POST "$API_BASE/events" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Winter Field Day 2025",
    "description": "Annual Winter Field Day event. Emergency preparedness exercise with portable radio operations.",
    "date": "2025-01-25T09:00",
    "timezone": "America/New_York",
    "points": 20,
    "latitude": 43.1389,
    "longitude": -77.5725
  }'
echo ""

# Event 3: Ham Radio License Class
echo "Creating event: Technician License Class..."
curl -X POST "$API_BASE/events" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Technician License Class - Session 1",
    "description": "First session of the Technician license preparation class. Open to all prospective hams.",
    "date": "2025-02-01T10:00",
    "timezone": "America/New_York",
    "points": 10,
    "latitude": 43.1610,
    "longitude": -77.6109
  }'
echo ""

# Event 4: Public Service Event
echo "Creating event: Rochester Marathon Communications..."
curl -X POST "$API_BASE/events" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Rochester Marathon Communications Support",
    "description": "Provide radio communications support for the Rochester Marathon. Multiple operator positions needed.",
    "date": "2025-05-18T06:00",
    "timezone": "America/New_York",
    "points": 15,
    "latitude": 43.1547,
    "longitude": -77.6158
  }'
echo ""

# Event 5: VE Testing Session
echo "Creating event: VE Testing Session..."
curl -X POST "$API_BASE/events" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Volunteer Examiner Testing Session",
    "description": "Monthly VE testing session for all amateur radio license classes.",
    "date": "2025-02-15T09:00",
    "timezone": "America/New_York",
    "points": 8,
    "latitude": 43.1389,
    "longitude": -77.5725
  }'
echo ""

# Event 6: Special Event Station
echo "Creating event: Museum Ships Weekend..."
curl -X POST "$API_BASE/events" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Museum Ships Weekend Special Event",
    "description": "Operating special event station K2RMW from the USS Edson. Contact other museum ships worldwide.",
    "date": "2025-06-07T10:00",
    "timezone": "America/New_York",
    "points": 12,
    "latitude": 43.2651,
    "longitude": -77.6647
  }'
echo ""

echo "✅ Events created successfully!"