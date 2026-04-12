#!/usr/bin/env bash
# state.sh - Persistent state management

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/../protocol.sh"

# Default state file location (can be overridden for testing)
STATE_FILE="${STATE_FILE:-$(cd "$SCRIPT_DIR/../.." && pwd)/.tui-state.json}"

# Initialize state file if it doesn't exist or is empty
init_state() {
    if [[ ! -f "$STATE_FILE" ]] || [[ ! -s "$STATE_FILE" ]]; then
        echo '{"lastEngagement":"","preferences":{}}' > "$STATE_FILE"
    fi
}

# Get a state value
# Usage: get_state <key>
get_state() {
    local key="$1"
    init_state
    $JQ -r ".$key // empty" "$STATE_FILE"
}

# Set a state value
# Usage: set_state <key> <value>
set_state() {
    local key="$1"
    local value="$2"
    init_state

    local tmp
    tmp=$(mktemp)
    $JQ --arg key "$key" --arg val "$value" '.[$key] = $val' "$STATE_FILE" > "$tmp"
    mv "$tmp" "$STATE_FILE"
}

# Get all state as JSON
get_all_state() {
    init_state
    cat "$STATE_FILE"
}

# Clear state
clear_state() {
    rm -f "$STATE_FILE"
    init_state
}
