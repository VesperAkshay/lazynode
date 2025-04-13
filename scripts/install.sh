#!/bin/bash
set -e

# LazyNode installation script
# Automatically detects OS and architecture and installs the appropriate binary

# Detect OS and architecture
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

# Map architecture to Go arch
if [ "$ARCH" = "x86_64" ]; then
  ARCH="amd64"
elif [ "$ARCH" = "aarch64" ] || [ "$ARCH" = "arm64" ]; then
  ARCH="arm64"
else
  echo "Unsupported architecture: $ARCH"
  exit 1
fi

# Get the latest version
VERSION=$(curl -s https://api.github.com/repos/VesperAkshay/lazynode/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/' | sed 's/v//')

if [ -z "$VERSION" ]; then
  echo "Could not determine the latest version. Please check your internet connection."
  exit 1
fi

# Determine download URL
DOWNLOAD_URL="https://github.com/VesperAkshay/lazynode/releases/download/v${VERSION}/lazynode_${VERSION}_${OS}_${ARCH}.tar.gz"
echo "Downloading LazyNode v${VERSION} for ${OS}/${ARCH}..."
echo "URL: $DOWNLOAD_URL"

# Create a temporary directory for the download
TMP_DIR=$(mktemp -d)
cd "$TMP_DIR"

# Download and extract
if command -v curl > /dev/null 2>&1; then
  curl -L -o lazynode.tar.gz "$DOWNLOAD_URL"
else
  if command -v wget > /dev/null 2>&1; then
    wget -O lazynode.tar.gz "$DOWNLOAD_URL"
  else
    echo "Error: Neither curl nor wget found. Please install one of them and try again."
    exit 1
  fi
fi

tar -xzf lazynode.tar.gz

# Install
echo "Installing LazyNode to /usr/local/bin (may require sudo)..."
sudo mv lazynode /usr/local/bin/
sudo chmod +x /usr/local/bin/lazynode

# Clean up
cd -
rm -rf "$TMP_DIR"

echo "LazyNode v${VERSION} has been installed successfully!"
echo "You can now run 'lazynode' to start using it." 