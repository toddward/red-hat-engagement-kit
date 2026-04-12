#!/usr/bin/env bash
set -euo pipefail

echo "Red Hat Engagement Kit TUI Setup"
echo "================================"
echo

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BIN_DIR="$SCRIPT_DIR/bin"

mkdir -p "$BIN_DIR"

# Detect platform
PLATFORM=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
[[ "$ARCH" == "x86_64" ]] && ARCH="amd64"
[[ "$ARCH" == "aarch64" || "$ARCH" == "arm64" ]] && ARCH="arm64"

echo "Platform: $PLATFORM-$ARCH"
echo

# Check/install jq
echo "Checking jq..."
if command -v jq &>/dev/null; then
    echo "  ✓ jq found in PATH"
elif [[ -x "$BIN_DIR/jq" ]]; then
    echo "  ✓ jq found in bin/"
else
    echo "  ↓ Downloading jq..."
    JQ_URL="https://github.com/jqlang/jq/releases/download/jq-1.7.1/jq-${PLATFORM}-${ARCH}"
    if curl -sL "$JQ_URL" -o "$BIN_DIR/jq" 2>/dev/null; then
        chmod +x "$BIN_DIR/jq"
        echo "  ✓ jq installed to bin/"
    else
        echo "  ✗ Failed to download jq"
        echo "    Manual install: https://jqlang.github.io/jq/download/"
        exit 1
    fi
fi

# Check/build tui-viewer
echo "Checking tui-viewer..."
if [[ -x "$BIN_DIR/tui-viewer" ]]; then
    echo "  ✓ tui-viewer found in bin/"
else
    if command -v go &>/dev/null; then
        echo "  → Building tui-viewer from source..."
        (cd "$SCRIPT_DIR/viewer" && go build -o "$BIN_DIR/tui-viewer" .)
        echo "  ✓ tui-viewer built"
    else
        echo "  ✗ Go not found, cannot build tui-viewer"
        echo "    Install Go: https://go.dev/dl/"
        echo "    Or download pre-built binary to bin/tui-viewer"
        exit 1
    fi
fi

echo
echo "Setup complete!"
echo "Run: ./tui.sh"
