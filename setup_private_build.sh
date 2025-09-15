#!/bin/bash

# Quick setup script for building private Ollama Docker images
# This script helps you set up the environment and provides instructions

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${GREEN}🚀 Ollama Private Docker Build Setup${NC}"
echo "======================================"

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo -e "${YELLOW}⚠️  Docker is not installed. Please install Docker first.${NC}"
    exit 1
fi

# Check if Docker buildx is available
if ! docker buildx version &> /dev/null; then
    echo -e "${YELLOW}⚠️  Docker buildx is not available. Please update Docker to a newer version.${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Docker and buildx are available${NC}"

# Prompt for Docker Hub username
read -p "Enter your Docker Hub username: " DOCKER_USERNAME

if [ -z "$DOCKER_USERNAME" ]; then
    echo "Username cannot be empty"
    exit 1
fi

# Create or update .env file
echo "# Ollama Private Build Configuration" > .env
echo "export DOCKER_USERNAME=\"$DOCKER_USERNAME\"" >> .env
echo "export DOCKER_REPO=\"$DOCKER_USERNAME/ollama\"" >> .env

echo -e "${GREEN}✅ Configuration saved to .env${NC}"

echo ""
echo -e "${BLUE}📋 Next steps:${NC}"
echo "1. Source the configuration:"
echo "   ${YELLOW}source .env${NC}"
echo ""
echo "2. Login to Docker Hub:"
echo "   ${YELLOW}docker login${NC}"
echo ""
echo "3. Run the build script:"
echo "   ${YELLOW}./scripts/build_and_push_private.sh${NC}"
echo ""
echo -e "${BLUE}💡 Or run everything in one go:${NC}"
echo "   ${YELLOW}source .env && ./scripts/build_and_push_private.sh${NC}"