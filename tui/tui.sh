#!/usr/bin/env bash
set -euo pipefail

TUI_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BIN_DIR="$TUI_DIR/bin"
CORE_DIR="$TUI_DIR/core"

# Check if setup has been run
if [[ ! -x "$BIN_DIR/tui-viewer" ]]; then
    echo "TUI not set up. Running setup..."
    "$TUI_DIR/setup.sh"
    echo
fi

# Add bin to PATH
export PATH="$BIN_DIR:$PATH"

# Change to project root (parent of tui/)
cd "$TUI_DIR/.."

# Create named pipes for bidirectional communication
# Works on bash 3.2+ (macOS default) unlike coproc which requires bash 4+
FIFO_DIR=$(mktemp -d)
FIFO_TO_CORE="$FIFO_DIR/to-core"
FIFO_FROM_CORE="$FIFO_DIR/from-core"
mkfifo "$FIFO_TO_CORE" "$FIFO_FROM_CORE"

cleanup() {
    kill "$CORE_PID" 2>/dev/null || true
    rm -rf "$FIFO_DIR"
}
trap cleanup EXIT

# Launch bash core: reads commands from FIFO, writes responses to FIFO
bash "$CORE_DIR/main.sh" < "$FIFO_TO_CORE" > "$FIFO_FROM_CORE" &
CORE_PID=$!

# Launch viewer: reads responses from core, writes commands to core
# Viewer opens /dev/tty directly for terminal rendering
"$BIN_DIR/tui-viewer" < "$FIFO_FROM_CORE" > "$FIFO_TO_CORE"
