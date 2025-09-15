#!/bin/bash

set -e

# Configuration - UPDATE THESE VALUES
DOCKER_USERNAME="${DOCKER_USERNAME:-your-dockerhub-username}"
DOCKER_REPO="${DOCKER_REPO:-${DOCKER_USERNAME}/ollama}"
PLATFORM="linux/amd64"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}🐳 Ollama Private Docker Build & Push Script${NC}"
echo "=================================================="

# Check if Docker username is set
if [ "$DOCKER_USERNAME" = "your-dockerhub-username" ]; then
    echo -e "${RED}❌ Please set your Docker Hub username:${NC}"
    echo "   export DOCKER_USERNAME=your-actual-username"
    echo "   Or edit this script and replace 'your-dockerhub-username'"
    exit 1
fi

# Source the environment variables
echo -e "${YELLOW}📋 Loading environment variables...${NC}"
. $(dirname $0)/env.sh

# Override the image repository
export FINAL_IMAGE_REPO="$DOCKER_REPO"
export PLATFORM="$PLATFORM"

echo "Building for: $PLATFORM"
echo "Image repository: $FINAL_IMAGE_REPO"
echo "Version: $VERSION"

# Check if user is logged into Docker
echo -e "${YELLOW}🔐 Checking Docker login status...${NC}"
if ! docker info | grep -q "Username"; then
    echo -e "${YELLOW}⚠️  Not logged into Docker Hub. Attempting login...${NC}"
    
    # Check if access token is provided via environment variable
    if [ -n "${DOCKER_ACCESS_TOKEN:-}" ]; then
        echo "Using access token from environment variable..."
        echo "$DOCKER_ACCESS_TOKEN" | docker login --username "$DOCKER_USERNAME" --password-stdin
    else
        echo "Please enter your Docker Hub credentials:"
        echo "💡 Tip: Use your access token as the password for better security"
        echo "   You can create one at: https://hub.docker.com/settings/security"
        docker login
    fi
    
    if [ $? -ne 0 ]; then
        echo -e "${RED}❌ Docker login failed${NC}"
        exit 1
    fi
fi

echo -e "${GREEN}✅ Docker login confirmed${NC}"

# Build and push the main image
echo -e "${YELLOW}🔨 Building and pushing main Ollama image...${NC}"
docker buildx build \
    --push \
    --platform=$PLATFORM \
    ${OLLAMA_COMMON_BUILD_ARGS} \
    -f Dockerfile \
    -t ${FINAL_IMAGE_REPO}:$VERSION \
    -t ${FINAL_IMAGE_REPO}:latest \
    .

if [ $? -ne 0 ]; then
    echo -e "${RED}❌ Main image build failed${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Main image built and pushed successfully${NC}"

# Build and push ROCm variant (optional, comment out if not needed)
echo -e "${YELLOW}🔨 Building and pushing ROCm variant...${NC}"
docker buildx build \
    --push \
    --platform=$PLATFORM \
    ${OLLAMA_COMMON_BUILD_ARGS} \
    --build-arg FLAVOR=rocm \
    -f Dockerfile \
    -t ${FINAL_IMAGE_REPO}:$VERSION-rocm \
    -t ${FINAL_IMAGE_REPO}:latest-rocm \
    .

if [ $? -ne 0 ]; then
    echo -e "${YELLOW}⚠️  ROCm image build failed (this is optional)${NC}"
else
    echo -e "${GREEN}✅ ROCm image built and pushed successfully${NC}"
fi

echo ""
echo -e "${GREEN}🎉 Build and push completed!${NC}"
echo "=================================================="
echo "Your images are available at:"
echo "  📦 ${FINAL_IMAGE_REPO}:latest"
echo "  📦 ${FINAL_IMAGE_REPO}:$VERSION"
echo "  📦 ${FINAL_IMAGE_REPO}:latest-rocm"
echo "  📦 ${FINAL_IMAGE_REPO}:$VERSION-rocm"
echo ""
echo -e "${YELLOW}To run your private image:${NC}"
echo "  docker run -d -v ollama:/root/.ollama -p 11434:11434 --gpus=all ${FINAL_IMAGE_REPO}:latest"
echo ""
echo -e "${YELLOW}For NVIDIA GPU support, make sure you have:${NC}"
echo "  - NVIDIA Container Toolkit installed"
echo "  - Docker configured with nvidia runtime"