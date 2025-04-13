#!/bin/bash
# install.sh - LazyNode installation script for Linux and macOS

set -e

# Determine OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Map architecture to Go architecture naming
if [ "$ARCH" = "x86_64" ]; then
    ARCH="amd64"
elif [ "$ARCH" = "aarch64" ] || [ "$ARCH" = "arm64" ]; then
    ARCH="arm64"
else
    echo "Unsupported architecture: $ARCH"
    exit 1
fi

# Get the latest version from GitHub
VERSION=$(curl -s https://api.github.com/repos/yourusername/lazynode/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
VERSION=${VERSION#v} # Remove 'v' prefix if present

echo "Installing LazyNode $VERSION for $OS/$ARCH..."

# Create temporary directory
TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

# Download the appropriate binary
ARCHIVE="lazynode_${VERSION}_${OS}_${ARCH}.tar.gz"
DOWNLOAD_URL="https://github.com/yourusername/lazynode/releases/download/v${VERSION}/${ARCHIVE}"

echo "Downloading $DOWNLOAD_URL..."
curl -L "$DOWNLOAD_URL" -o "$TMP_DIR/$ARCHIVE"

# Extract the archive
echo "Extracting..."
tar -xzf "$TMP_DIR/$ARCHIVE" -C "$TMP_DIR"

# Install the binary
BINARY="lazynode_${VERSION}_${OS}_${ARCH}"
INSTALL_DIR="/usr/local/bin"

echo "Installing to $INSTALL_DIR/lazynode..."
if [ ! -d "$INSTALL_DIR" ]; then
    mkdir -p "$INSTALL_DIR"
fi

if [ -w "$INSTALL_DIR" ]; then
    # User has write permission
    mv "$TMP_DIR/$BINARY" "$INSTALL_DIR/lazynode"
    chmod +x "$INSTALL_DIR/lazynode"
else
    # Need sudo
    echo "Elevated permissions required to install to $INSTALL_DIR"
    sudo mv "$TMP_DIR/$BINARY" "$INSTALL_DIR/lazynode"
    sudo chmod +x "$INSTALL_DIR/lazynode"
fi

echo "LazyNode $VERSION has been installed to $INSTALL_DIR/lazynode"
echo "Run 'lazynode --version' to verify the installation." 