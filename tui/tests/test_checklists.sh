#!/usr/bin/env bash
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
export PROJECT_ROOT="$SCRIPT_DIR/fixtures"
source "$SCRIPT_DIR/../core/lib/checklists.sh"

test_list_checklists() {
    local result
    result=$(list_checklists)
    if ! echo "$result" | grep -q '"name":"test-checklist"'; then
        echo "FAIL: list_checklists did not find test-checklist"
        echo "Got: $result"
        return 1
    fi
    echo "PASS: list_checklists"
}

test_get_checklist() {
    local result
    result=$(get_checklist "test-checklist")
    if ! echo "$result" | grep -q '"title":"Section One"'; then
        echo "FAIL: get_checklist did not parse sections"
        echo "Got: $result"
        return 1
    fi
    if ! echo "$result" | grep -q '"checked":true'; then
        echo "FAIL: get_checklist did not detect checked item"
        return 1
    fi
    echo "PASS: get_checklist"
}

test_list_checklists
test_get_checklist
echo "All checklists tests passed"
