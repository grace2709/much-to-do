#!/bin/bash
set -euo pipefail

IMAGE_NAME="muchtodo-api"
IMAGE_TAG="${1:-latest}"

echo "======================================"
echo "  Building Docker image: $IMAGE_NAME:$IMAGE_TAG"
echo "======================================"

docker build \
  --tag "$IMAGE_NAME:$IMAGE_TAG" \
  --file Dockerfile \
  .

echo ""
echo "✅ Build complete: $IMAGE_NAME:$IMAGE_TAG"
echo ""
docker images "$IMAGE_NAME"
