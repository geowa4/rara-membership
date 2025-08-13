#!/bin/bash

# Script to populate point allocations via API
# Assumes the server is running on http://localhost:8080

API_BASE="http://localhost:8080/api"

echo "🎯 Creating point allocations..."

# Note: These allocations assume events are created with IDs 1-6
# and members exist with the call signs from 01_members.sh

# KD2ABC attended monthly meeting (event 1)
echo "Allocating points to KD2ABC for Monthly Meeting..."
curl -X POST "$API_BASE/members/KD2ABC/points" \
  -H "Content-Type: application/json" \
  -d '{
    "event_id": 1,
    "notes": "Attended January monthly meeting and gave presentation on antenna modeling"
  }'
echo ""

# W2XYZ helped with Winter Field Day (event 2)
echo "Allocating points to W2XYZ for Winter Field Day..."
curl -X POST "$API_BASE/members/W2XYZ/points" \
  -H "Content-Type: application/json" \
  -d '{
    "event_id": 2,
    "points": 25,
    "notes": "Led Winter Field Day setup and operated for 12 hours. Extra points for coordination role."
  }'
echo ""

# KC2DEF attended Tech class (event 3)
echo "Allocating points to KC2DEF for Tech License Class..."
curl -X POST "$API_BASE/members/KC2DEF/points" \
  -H "Content-Type: application/json" \
  -d '{
    "event_id": 3,
    "notes": "Attended as student helper, assisting new hams with practice questions"
  }'
echo ""

# N2GHI multiple events
echo "Allocating points to N2GHI for Monthly Meeting..."
curl -X POST "$API_BASE/members/N2GHI/points" \
  -H "Content-Type: application/json" \
  -d '{
    "event_id": 1,
    "notes": "Attended meeting and volunteered for refreshments"
  }'
echo ""

echo "Allocating points to N2GHI for Winter Field Day..."
curl -X POST "$API_BASE/members/N2GHI/points" \
  -H "Content-Type: application/json" \
  -d '{
    "event_id": 2,
    "notes": "Operated digital modes station for 6 hours"
  }'
echo ""

# KB2JKL teaching Tech class
echo "Allocating points to KB2JKL for Tech License Class..."
curl -X POST "$API_BASE/members/KB2JKL/points" \
  -H "Content-Type: application/json" \
  -d '{
    "event_id": 3,
    "points": 15,
    "notes": "Instructor for Session 1 - covered FCC rules and regulations"
  }'
echo ""

# Multiple members at monthly meeting
echo "Allocating points to KB2JKL for Monthly Meeting..."
curl -X POST "$API_BASE/members/KB2JKL/points" \
  -H "Content-Type: application/json" \
  -d '{
    "event_id": 1,
    "notes": "Attended and participated in technical discussion"
  }'
echo ""

echo "Allocating points to KC2DEF for Monthly Meeting..."
curl -X POST "$API_BASE/members/KC2DEF/points" \
  -H "Content-Type: application/json" \
  -d '{
    "event_id": 1,
    "notes": "First meeting as new member"
  }'
echo ""

# Bonus points for special contributions
echo "Allocating bonus points to KD2ABC..."
curl -X POST "$API_BASE/members/KD2ABC/points" \
  -H "Content-Type: application/json" \
  -d '{
    "event_id": 2,
    "points": 30,
    "notes": "Provided generator and emergency power equipment for Winter Field Day"
  }'
echo ""

echo "Allocating points to W2XYZ for Tech License Class..."
curl -X POST "$API_BASE/members/W2XYZ/points" \
  -H "Content-Type: application/json" \
  -d '{
    "event_id": 3,
    "points": 12,
    "notes": "Guest speaker on propagation and antennas"
  }'
echo ""

echo "✅ Point allocations created successfully!"