#!/bin/bash

# Release script to build binaries for multiple platforms

set -e

# Output directory
OUTPUT_DIR="bin"
mkdir -p "$OUTPUT_DIR"

# Build function
build() {
    local GOOS=$1
    local GOARCH=$2
    local OUTPUT_NAME="q-${GOOS}-${GOARCH}"
    
    if [ "$GOOS" = "windows" ]; then
        OUTPUT_NAME="${OUTPUT_NAME}.exe"
    fi
    
    echo "Building for $GOOS/$GOARCH..."
    GOOS=$GOOS GOARCH=$GOARCH go build -o "$OUTPUT_DIR/$OUTPUT_NAME" ./cmd/q
    echo "Built $OUTPUT_DIR/$OUTPUT_NAME"
}

# Build for each platform
build linux amd64
build linux arm64
build darwin amd64
build darwin arm64
build windows amd64
build windows arm64

echo "All binaries built in $OUTPUT_DIR/"
