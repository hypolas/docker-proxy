#!/bin/bash
# Test multi-platform Docker build locally
# This simulates what GitHub Actions will do

set -e

echo "🐳 Testing multi-platform Docker build locally"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
IMAGE_NAME="hypolas/proxy-docker"
VERSION="${1:-test}"
PLATFORMS="linux/amd64,linux/arm64,linux/arm/v7"

echo -e "${YELLOW}📦 Version:${NC} $VERSION"
echo -e "${YELLOW}🏗️  Platforms:${NC} $PLATFORMS"
echo ""

# Check if docker buildx is available
if ! docker buildx version > /dev/null 2>&1; then
    echo -e "${RED}❌ Error: docker buildx is not available${NC}"
    echo "Install it with: docker buildx install"
    exit 1
fi

# Create buildx builder if it doesn't exist
if ! docker buildx inspect multiarch > /dev/null 2>&1; then
    echo -e "${YELLOW}🔧 Creating buildx builder...${NC}"
    docker buildx create --name multiarch --use
    docker buildx inspect --bootstrap
fi

# Build for all platforms (without push)
echo -e "${GREEN}🚀 Building for multiple platforms...${NC}"
echo ""

docker buildx build \
    --platform "$PLATFORMS" \
    --tag "$IMAGE_NAME:$VERSION" \
    --tag "$IMAGE_NAME:latest" \
    --build-arg VERSION="$VERSION" \
    --build-arg BUILD_DATE="$(date -u +'%Y-%m-%dT%H:%M:%SZ')" \
    --build-arg VCS_REF="$(git rev-parse --short HEAD)" \
    --load \
    .

echo ""
echo -e "${GREEN}✅ Build successful!${NC}"
echo ""
echo "📋 Image details:"
docker images "$IMAGE_NAME" | head -2

echo ""
echo "🧪 Test the image:"
echo "  docker run --rm $IMAGE_NAME:$VERSION --version"
echo ""
echo "🚀 To push to Docker Hub:"
echo "  docker buildx build --platform $PLATFORMS --push -t $IMAGE_NAME:$VERSION ."
