#!/bin/bash
# build.sh - Cross-platform build script for LazyNode

set -e

# Set application name and version
APP_NAME="lazynode"
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# Create dist directory
mkdir -p dist

# Build flags
LDFLAGS="-s -w -X 'main.Version=$VERSION' -X 'main.Commit=$COMMIT' -X 'main.BuildDate=$BUILD_DATE'"

echo "Building LazyNode $VERSION (commit: $COMMIT, date: $BUILD_DATE)"

# Build for Linux
echo "Building for Linux..."
GOOS=linux GOARCH=amd64 go build -ldflags "$LDFLAGS" -o "dist/${APP_NAME}_${VERSION}_linux_amd64" cmd/lazynode/main.go
GOOS=linux GOARCH=arm64 go build -ldflags "$LDFLAGS" -o "dist/${APP_NAME}_${VERSION}_linux_arm64" cmd/lazynode/main.go

# Build for macOS
echo "Building for macOS..."
GOOS=darwin GOARCH=amd64 go build -ldflags "$LDFLAGS" -o "dist/${APP_NAME}_${VERSION}_darwin_amd64" cmd/lazynode/main.go
GOOS=darwin GOARCH=arm64 go build -ldflags "$LDFLAGS" -o "dist/${APP_NAME}_${VERSION}_darwin_arm64" cmd/lazynode/main.go

# Build for Windows
echo "Building for Windows..."
GOOS=windows GOARCH=amd64 go build -ldflags "$LDFLAGS" -o "dist/${APP_NAME}_${VERSION}_windows_amd64.exe" cmd/lazynode/main.go

# Create zip archives
echo "Creating distribution archives..."
cd dist

# Linux archives
tar -czvf "${APP_NAME}_${VERSION}_linux_amd64.tar.gz" "${APP_NAME}_${VERSION}_linux_amd64"
tar -czvf "${APP_NAME}_${VERSION}_linux_arm64.tar.gz" "${APP_NAME}_${VERSION}_linux_arm64"

# macOS archives
tar -czvf "${APP_NAME}_${VERSION}_darwin_amd64.tar.gz" "${APP_NAME}_${VERSION}_darwin_amd64"
tar -czvf "${APP_NAME}_${VERSION}_darwin_arm64.tar.gz" "${APP_NAME}_${VERSION}_darwin_arm64"

# Windows archive
zip "${APP_NAME}_${VERSION}_windows_amd64.zip" "${APP_NAME}_${VERSION}_windows_amd64.exe"

# Generate SHA256 checksums
echo "Generating checksums..."
shasum -a 256 "${APP_NAME}_${VERSION}_linux_amd64.tar.gz" > "${APP_NAME}_${VERSION}_checksums.txt"
shasum -a 256 "${APP_NAME}_${VERSION}_linux_arm64.tar.gz" >> "${APP_NAME}_${VERSION}_checksums.txt"
shasum -a 256 "${APP_NAME}_${VERSION}_darwin_amd64.tar.gz" >> "${APP_NAME}_${VERSION}_checksums.txt"
shasum -a 256 "${APP_NAME}_${VERSION}_darwin_arm64.tar.gz" >> "${APP_NAME}_${VERSION}_checksums.txt"
shasum -a 256 "${APP_NAME}_${VERSION}_windows_amd64.zip" >> "${APP_NAME}_${VERSION}_checksums.txt"

cd ..

echo "Build complete! Distribution files are in the dist directory." 