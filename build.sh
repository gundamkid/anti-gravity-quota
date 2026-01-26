#!/bin/bash

# Build script for ag-quota
# Builds binaries for multiple platforms

set -e

VERSION="0.1.0"
APP_NAME="ag-quota"
BUILD_DIR="dist"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}Building ${APP_NAME} v${VERSION}${NC}"
echo ""

# Clean previous builds
rm -rf ${BUILD_DIR}
mkdir -p ${BUILD_DIR}

# Build for multiple platforms
platforms=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
)

for platform in "${platforms[@]}"; do
    IFS="/" read -r -a platform_split <<< "$platform"
    GOOS="${platform_split[0]}"
    GOARCH="${platform_split[1]}"

    output_name="${APP_NAME}-${GOOS}-${GOARCH}"

    if [ "$GOOS" = "windows" ]; then
        output_name+=".exe"
    fi

    echo -e "${BLUE}Building for ${GOOS}/${GOARCH}...${NC}"

    GOOS=$GOOS GOARCH=$GOARCH go build \
        -ldflags "-X main.version=${VERSION}" \
        -o "${BUILD_DIR}/${output_name}" \
        ./cmd/ag-quota

    echo -e "${GREEN}✓ Built ${output_name}${NC}"
done

echo ""
echo -e "${GREEN}Build complete! Binaries are in ${BUILD_DIR}/${NC}"
echo ""
ls -lh ${BUILD_DIR}
