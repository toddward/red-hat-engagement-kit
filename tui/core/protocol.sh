#!/usr/bin/env bash
# protocol.sh - JSON protocol helpers for viewer communication

# Get jq path (prefer bundled, fall back to system)
_get_jq() {
    local script_dir
    script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
    if [[ -x "$script_dir/bin/jq" ]]; then
        echo "$script_dir/bin/jq"
    elif command -v jq &>/dev/null; then
        echo "jq"
    else
        echo "ERROR: jq not found" >&2
        exit 1
    fi
}

JQ=$(_get_jq)

# Send a response to the viewer
# Usage: send_response <id> <payload_json>
send_response() {
    local id="$1"
    local payload="$2"

    $JQ -cn \
        --arg id "$id" \
        --argjson payload "$payload" \
        '{type: "response", id: $id, payload: $payload}'
}

# Send an event to the viewer (no correlation ID)
# Usage: send_event <event_type> <payload_json>
send_event() {
    local event_type="$1"
    local payload="$2"

    $JQ -cn \
        --arg event "$event_type" \
        --argjson payload "$payload" \
        '{type: "event", payload: ({event: $event} + $payload)}'
}

# Send an error to the viewer
# Usage: send_error <message> [id]
send_error() {
    local message="$1"
    local id="${2:-}"

    if [[ -n "$id" ]]; then
        $JQ -cn \
            --arg id "$id" \
            --arg msg "$message" \
            '{type: "error", id: $id, payload: {message: $msg}}'
    else
        $JQ -cn \
            --arg msg "$message" \
            '{type: "error", payload: {message: $msg}}'
    fi
}

# Parse a field from a command JSON
# Usage: echo "$json" | parse_command_field <field>
parse_command_field() {
    local field="$1"
    $JQ -r ".$field // empty"
}

# Parse args from a command JSON
# Usage: echo "$json" | parse_command_args
parse_command_args() {
    $JQ -c ".args // {}"
}
