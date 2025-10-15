#!/bin/bash

# Vignette Authentication API Test Script

echo "🟣 Testing Vignette Authentication Service"
echo "=========================================="
echo ""

BASE_URL="http://localhost:8002/api/v1"

# Test 1: Health Check
echo "1️⃣  Testing Health Check..."
curl -s -X GET http://localhost:8002/health | jq '.'
echo ""
echo ""

# Test 2: Sign Up
echo "2️⃣  Testing User Signup..."
SIGNUP_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/signup" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john.doe@vignette.com",
    "full_name": "John Doe",
    "password": "SecurePassword123"
  }')

echo "$SIGNUP_RESPONSE" | jq '.'
ACCESS_TOKEN=$(echo "$SIGNUP_RESPONSE" | jq -r '.data.access_token')
echo ""
echo "✅ Access Token: ${ACCESS_TOKEN:0:50}..."
echo ""
echo ""

# Test 3: Get Current User
echo "3️⃣  Testing Get Current User..."
curl -s -X GET "$BASE_URL/auth/me" \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq '.'
echo ""
echo ""

# Test 4: Login
echo "4️⃣  Testing Login with existing user..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username_or_email": "johndoe",
    "password": "SecurePassword123"
  }')

echo "$LOGIN_RESPONSE" | jq '.'
echo ""
echo ""

# Test 5: Logout
echo "5️⃣  Testing Logout..."
curl -s -X POST "$BASE_URL/auth/logout" \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq '.'
echo ""
echo ""

echo "✅ Vignette Authentication Service Tests Complete!"
echo ""
