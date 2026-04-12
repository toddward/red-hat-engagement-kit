#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
export PROJECT_ROOT="$SCRIPT_DIR/fixtures"
source "$SCRIPT_DIR/../core/lib/engagements.sh"

test_list_engagements() {
    local result
    result=$(list_engagements)

    if ! echo "$result" | grep -q '"slug":"test-customer"'; then
        echo "FAIL: list_engagements did not find test-customer"
        echo "Got: $result"
        return 1
    fi
    if ! echo "$result" | grep -q '"hasContext":true'; then
        echo "FAIL: test-customer should have hasContext=true"
        return 1
    fi
    echo "PASS: list_engagements"
}

test_get_engagement() {
    local result
    result=$(get_engagement "test-customer")

    if ! echo "$result" | grep -q '"slug":"test-customer"'; then
        echo "FAIL: get_engagement did not return correct slug"
        return 1
    fi
    echo "PASS: get_engagement"
}

test_list_engagements
test_get_engagement

echo "All engagements tests passed"
