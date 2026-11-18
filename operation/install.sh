#!/bin/bash

set -e

# Configuration
REPO="montralis/flight-tracker"
INSTALL_DIR="$HOME/flight-tracker"
SERVICE_NAME="flight-tracker"

echo "Flight Tracker Setup"

# Check dependencies
for cmd in gh systemctl; do
    if ! command -v "$cmd" &> /dev/null; then
        echo "Error: $cmd is not installed."
        exit 1
    fi
done

# Detect architecture
ARCH=$(uname -m)

case "$ARCH" in
    aarch64|arm64)
        ARCH="ARM64"
        ;;
    x86_64|amd64)
        ARCH="AMD64"
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

echo "Detected architecture: $ARCH"

# Get latest release tag
TAG=$(gh release view \
    --repo "$REPO" \
    --json tagName \
    --jq '.tagName')

if [ -z "$TAG" ]; then
    echo "Could not determine latest release tag."
    exit 1
fi

echo "Latest release: $TAG"

# Download binary
ASSET="Flight_Tracker_${ARCH}"
BINARY="Flight_Tracker_${ARCH}_${TAG}"

TEMP_DIR=$(mktemp -d)
trap 'rm -rf "$TEMP_DIR"' EXIT

echo "Downloading $ASSET..."

gh release download "$TAG" \
    --repo "$REPO" \
    --pattern "$ASSET" \
    --dir "$TEMP_DIR"

# Check download
if [ ! -f "$TEMP_DIR/$ASSET" ]; then
    echo "Binary not found in release."
    exit 1
fi

# Prepare installation directory
mkdir -p "$INSTALL_DIR"

# Install versioned binary
cp "$TEMP_DIR/$ASSET" "$INSTALL_DIR/$BINARY"
chmod +x "$INSTALL_DIR/$BINARY"

# Create stable symlink for systemd
ln -sfn "$BINARY" "$INSTALL_DIR/$ASSET"

echo "Installed: $INSTALL_DIR/$BINARY"

# Install systemd service
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

sudo cp "$SCRIPT_DIR/flight-tracker@.service" \
    /etc/systemd/system/

# Reload systemd
sudo systemctl daemon-reload

# Enable and start service
echo "Starting Flight Tracker..."

sudo systemctl enable --now "$SERVICE_NAME@$USER"

echo "Installation completed!"

sudo systemctl status "$SERVICE_NAME@$USER" --no-pager