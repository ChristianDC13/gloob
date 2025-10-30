#!/bin/bash

# Gloob Installation Script
# Installs the latest version of Gloob for macOS

set -e

REPO="ChristianDC13/gloob"
BINARY_NAME="gloob"
INSTALL_DIR="/usr/local/bin"

echo "🫧 Installing Gloob Programming Language..."
echo ""

# Check if running on macOS
if [[ "$OSTYPE" != "darwin"* ]]; then
    echo "❌ Error: This installer only supports macOS."
    echo "   Please build from source: https://github.com/$REPO#build-from-source"
    exit 1
fi

# Get the latest release binary URL
echo "📦 Fetching latest release..."
DOWNLOAD_URL="https://github.com/$REPO/releases/latest/download/gloob-macos"

# Create temporary directory
TMP_DIR=$(mktemp -d)
TMP_FILE="$TMP_DIR/$BINARY_NAME"

# Download the binary
echo "⬇️  Downloading Gloob..."
if command -v curl &> /dev/null; then
    curl -fsSL "$DOWNLOAD_URL" -o "$TMP_FILE"
elif command -v wget &> /dev/null; then
    wget -q "$DOWNLOAD_URL" -O "$TMP_FILE"
else
    echo "❌ Error: curl or wget is required"
    exit 1
fi

# Make it executable
chmod +x "$TMP_FILE"

# Install to /usr/local/bin
echo "📥 Installing to $INSTALL_DIR..."
if [ -w "$INSTALL_DIR" ]; then
    mv "$TMP_FILE" "$INSTALL_DIR/$BINARY_NAME"
else
    echo "   (requires sudo permission)"
    sudo mv "$TMP_FILE" "$INSTALL_DIR/$BINARY_NAME"
fi

# Clean up
rm -rf "$TMP_DIR"

# Verify installation
if command -v gloob &> /dev/null; then
    echo ""
    echo "✅ Gloob installed successfully!"
    echo ""
    echo "Get started:"
    echo "  • Run 'gloob' to start the REPL"
    echo "  • Run 'gloob yourfile.gloob' to execute a file"
    echo ""
    echo "Learn more: https://github.com/$REPO"
else
    echo "⚠️  Installation complete, but 'gloob' not found in PATH"
    echo "   You may need to restart your terminal or add $INSTALL_DIR to your PATH"
fi

