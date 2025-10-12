#!/bin/bash

# Script to build obsidian-index with custom version information
# Usage: ./scripts/build-with-version.sh [version]

set -e

VERSION=${1:-"1.0.0"}
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

echo "Building obsidian-index with version: $VERSION"
echo "Git commit: $GIT_COMMIT"
echo "Build date: $BUILD_DATE"

# Build with version information
go build -ldflags "-X github.com/nzb3/obsidian-index/internal/version.Version=$VERSION -X github.com/nzb3/obsidian-index/internal/version.GitCommit=$GIT_COMMIT -X github.com/nzb3/obsidian-index/internal/version.BuildDate=$BUILD_DATE" -o build/obsidian-index ./cmd

echo "Build completed successfully!"
echo "Test version information:"
./build/obsidian-index --version
