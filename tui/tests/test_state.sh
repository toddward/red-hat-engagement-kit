#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/../core/lib/state.sh"

# Use temp file for testing
STATE_FILE=$(mktemp)
export STATE_FILE

cleanup() {
    rm -f "$STATE_FILE"
}
trap cleanup EXIT

# Test initial state
test_init_state() {
    init_state

    if [[ ! -f "$STATE_FILE" ]]; then
        echo "FAIL: init_state did not create file"
        return 1
    fi
    echo "PASS: init_state"
}

# Test get/set state
test_get_set_state() {
    set_state "lastEngagement" "acme-corp"
    local result
    result=$(get_state "lastEngagement")

    if [[ "$result" != "acme-corp" ]]; then
        echo "FAIL: get_state returned '$result' instead of 'acme-corp'"
        return 1
    fi
    echo "PASS: get_set_state"
}

# Test get_all_state
test_get_all_state() {
    set_state "foo" "bar"
    local all
    all=$(get_all_state)

    if ! echo "$all" | grep -q '"foo"'; then
        echo "FAIL: get_all_state missing foo"
        return 1
    fi
    echo "PASS: get_all_state"
}

# Run tests
test_init_state
test_get_set_state
test_get_all_state

echo "All state tests passed"
