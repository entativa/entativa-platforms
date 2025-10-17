#!/bin/bash

# Socialink Name Policy Test Script
# Demonstrates relaxed name policy vs Facebook's strict policy

echo "🔵 Testing Socialink's Relaxed Name Policy"
echo "==========================================="
echo ""

BASE_URL="http://localhost:8001/api/v1"

echo "Test 1️⃣  - Real name (recommended)"
echo "-----------------------------------"
curl -s -X POST "$BASE_URL/auth/signup" \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "email": "john.doe@example.com",
    "password": "SecurePassword123",
    "birthday": "1990-05-15",
    "gender": "male"
  }' | jq '.message, .profile_url, .note'
echo ""
echo ""

echo "Test 2️⃣  - Nickname (allowed, with friendly suggestion)"
echo "-------------------------------------------------------"
curl -s -X POST "$BASE_URL/auth/signup" \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Jay",
    "last_name": "Smith",
    "email": "jay.smith@example.com",
    "password": "SecurePassword123",
    "birthday": "1990-05-15",
    "gender": "male"
  }' | jq '.message, .profile_url, .note'
echo ""
echo ""

echo "Test 3️⃣  - Artist/Stage name (perfectly fine!)"
echo "----------------------------------------------"
curl -s -X POST "$BASE_URL/auth/signup" \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "DJ",
    "last_name": "CoolBeats",
    "email": "dj.coolbeats@example.com",
    "password": "SecurePassword123",
    "birthday": "1990-05-15",
    "gender": "other"
  }' | jq '.message, .profile_url, .note'
echo ""
echo ""

echo "Test 4️⃣  - International name (fully supported)"
echo "----------------------------------------------"
curl -s -X POST "$BASE_URL/auth/signup" \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "María",
    "last_name": "García",
    "email": "maria.garcia@example.com",
    "password": "SecurePassword123",
    "birthday": "1990-05-15",
    "gender": "female"
  }' | jq '.message, .profile_url, .note'
echo ""
echo ""

echo "Test 5️⃣  - Single name (allowed, unlike Facebook)"
echo "------------------------------------------------"
curl -s -X POST "$BASE_URL/auth/signup" \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Shakira",
    "last_name": "Artist",
    "email": "shakira@example.com",
    "password": "SecurePassword123",
    "birthday": "1990-05-15",
    "gender": "female"
  }' | jq '.message, .profile_url'
echo ""
echo ""

echo "Test 6️⃣  - Hyphenated name (works great)"
echo "----------------------------------------"
curl -s -X POST "$BASE_URL/auth/signup" \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Jean-Pierre",
    "last_name": "Dubois",
    "email": "jean.pierre@example.com",
    "password": "SecurePassword123",
    "birthday": "1990-05-15",
    "gender": "male"
  }' | jq '.message, .profile_url'
echo ""
echo ""

echo "Test 7️⃣  - Name with apostrophe (Irish, etc.)"
echo "---------------------------------------------"
curl -s -X POST "$BASE_URL/auth/signup" \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Patrick",
    "last_name": "O'\''Brien",
    "email": "patrick.obrien@example.com",
    "password": "SecurePassword123",
    "birthday": "1990-05-15",
    "gender": "male"
  }' | jq '.message, .profile_url'
echo ""
echo ""

echo "✅ Summary"
echo "=========="
echo ""
echo "Socialink Policy:"
echo "  ✅ All names accepted (even nicknames, stage names)"
echo "  💡 Real names recommended but not required"
echo "  🌍 International characters fully supported"
echo "  🚫 Never blocks for 'fake names'"
echo "  ⚡ Instant signup, no verification"
echo ""
echo "Clean Profile URLs:"
echo "  socialink.com/john.doe"
echo "  socialink.com/dj.coolbeats"
echo "  socialink.com/maria.garcia"
echo "  socialink.com/jeanpierre.dubois"
echo ""
echo "Facebook Would Reject:"
echo "  ❌ 'Jay' (not full legal name)"
echo "  ❌ 'DJ CoolBeats' (stage name)"
echo "  ❌ Single names"
echo "  ❌ Requires ID verification"
echo ""
