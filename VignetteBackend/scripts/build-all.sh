#!/bin/bash

set -e

echo "🏗️  Building all Vignette services..."
echo ""

FAILED_BUILDS=()
SUCCESS_BUILDS=()

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

build_go_service() {
    local service_name=$1
    local service_path=$2
    
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "Building: $service_name"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    
    cd "$service_path"
    
    # Fix module path
    if grep -q "github.com/entativa" go.mod 2>/dev/null; then
        sed -i 's|github.com/entativa/vignette/[^/]*|vignette/'"$(basename $service_path)"'|g' go.mod
    fi
    
    # Fix imports in all Go files
    find . -name "*.go" -type f -exec sed -i 's|github.com/entativa/vignette/[^/]*/|vignette/'"$(basename $service_path)"'/|g' {} \;
    
    # Build
    if go mod tidy && go build -o /tmp/$(basename $service_path) ./cmd/api > /tmp/build.log 2>&1; then
        echo -e "${GREEN}✅ $service_name built successfully${NC}"
        SUCCESS_BUILDS+=("$service_name")
    else
        echo -e "${RED}❌ $service_name build failed${NC}"
        tail -20 /tmp/build.log
        FAILED_BUILDS+=("$service_name")
    fi
    
    cd - > /dev/null
    echo ""
}

# Build Go services
build_go_service "User Service" "/workspace/VignetteBackend/services/user-service"
build_go_service "Community Service" "/workspace/VignetteBackend/services/community-service"
build_go_service "Streaming Service" "/workspace/VignetteBackend/services/live-streaming-service"
build_go_service "Creator Service" "/workspace/VignetteBackend/services/creator-service"
build_go_service "Settings Service" "/workspace/VignetteBackend/services/settings-service"
build_go_service "Event Service" "/workspace/SocialinkBackend/services/event-service"

# Summary
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "BUILD SUMMARY"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo -e "${GREEN}✅ Successful builds: ${#SUCCESS_BUILDS[@]}${NC}"
for service in "${SUCCESS_BUILDS[@]}"; do
    echo "   - $service"
done

if [ ${#FAILED_BUILDS[@]} -gt 0 ]; then
    echo ""
    echo -e "${RED}❌ Failed builds: ${#FAILED_BUILDS[@]}${NC}"
    for service in "${FAILED_BUILDS[@]}"; do
        echo "   - $service"
    done
    exit 1
fi

echo ""
echo -e "${GREEN}🎉 All services built successfully!${NC}"
