#!/bin/bash

# Script to populate point redemptions via API
# Assumes the server is running on http://localhost:8080

API_BASE="http://localhost:8080/api"

echo "🎁 Creating point redemptions..."

# KD2ABC redeems points for club merchandise
echo "Creating redemption for KD2ABC..."
curl -X POST "$API_BASE/members/KD2ABC/redemptions" \
  -H "Content-Type: application/json" \
  -d '{
    "points": 10,
    "notes": "Redeemed for RARA club hat"
  }'
echo ""

# W2XYZ redeems points for hamfest ticket
echo "Creating redemption for W2XYZ..."
curl -X POST "$API_BASE/members/W2XYZ/redemptions" \
  -H "Content-Type: application/json" \
  -d '{
    "points": 15,
    "notes": "Redeemed for Rochester Hamfest admission ticket"
  }'
echo ""

# N2GHI redeems points for club dinner
echo "Creating redemption for N2GHI..."
curl -X POST "$API_BASE/members/N2GHI/redemptions" \
  -H "Content-Type: application/json" \
  -d '{
    "points": 20,
    "notes": "Redeemed for Annual Club Dinner ticket"
  }'
echo ""

# KB2JKL small redemption for QSL cards
echo "Creating redemption for KB2JKL..."
curl -X POST "$API_BASE/members/KB2JKL/redemptions" \
  -H "Content-Type: application/json" \
  -d '{
    "points": 5,
    "notes": "Redeemed for pack of club QSL cards"
  }'
echo ""

# KD2ABC another redemption
echo "Creating second redemption for KD2ABC..."
curl -X POST "$API_BASE/members/KD2ABC/redemptions" \
  -H "Content-Type: application/json" \
  -d '{
    "points": 8,
    "notes": "Redeemed for RARA coffee mug"
  }'
echo ""

echo "✅ Point redemptions created successfully!"