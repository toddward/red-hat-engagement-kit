#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/../core/protocol.sh"

# Test send_response
test_send_response() {
    local output
    output=$(send_response "test-123" '{"foo":"bar"}')

    local expected='{"type":"response","id":"test-123","payload":{"foo":"bar"}}'
    if [[ "$output" != "$expected" ]]; then
        echo "FAIL: send_response"
        echo "  Expected: $expected"
        echo "  Got: $output"
        return 1
    fi
    echo "PASS: send_response"
}

# Test send_event
test_send_event() {
    local output
    output=$(send_event "assistant" '{"text":"hello"}')

    local expected='{"type":"event","payload":{"event":"assistant","text":"hello"}}'
    if [[ "$output" != "$expected" ]]; then
        echo "FAIL: send_event"
        echo "  Expected: $expected"
        echo "  Got: $output"
        return 1
    fi
    echo "PASS: send_event"
}

# Test parse_command
test_parse_command() {
    local cmd='{"cmd":"list_skills","id":"abc-123","args":{}}'

    local parsed_cmd parsed_id
    parsed_cmd=$(echo "$cmd" | parse_command_field "cmd")
    parsed_id=$(echo "$cmd" | parse_command_field "id")

    if [[ "$parsed_cmd" != "list_skills" ]]; then
        echo "FAIL: parse_command_field cmd"
        return 1
    fi
    if [[ "$parsed_id" != "abc-123" ]]; then
        echo "FAIL: parse_command_field id"
        return 1
    fi
    echo "PASS: parse_command"
}

# Run tests
test_send_response
test_send_event
test_parse_command

echo "All protocol tests passed"
