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

# Launch with bidirectional communication using coproc
# Bash core reads commands from stdin, writes responses to stdout
# Viewer reads responses from bash's stdout, writes commands to bash's stdin
# Viewer opens /dev/tty directly for terminal rendering
coproc CORE {
    bash "$CORE_DIR/main.sh"
}

# Viewer: stdin connected to bash stdout (reads responses),
#          stdout connected to bash stdin (sends commands)
"$BIN_DIR/tui-viewer" <&"${CORE[0]}" >&"${CORE[1]}"

# Cleanup
kill "${CORE_PID}" 2>/dev/null || true
