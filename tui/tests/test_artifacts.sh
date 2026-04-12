#!/usr/bin/env bash
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
export PROJECT_ROOT="$SCRIPT_DIR/fixtures"
source "$SCRIPT_DIR/../core/lib/artifacts.sh"

test_list_artifacts() {
    local result
    result=$(list_artifacts "test-customer")
    if ! echo "$result" | grep -q '"name":"discovery"'; then
        echo "FAIL: list_artifacts should include discovery directory"
        echo "Got: $result"
        return 1
    fi
    echo "PASS: list_artifacts"
}

test_read_artifact() {
    local result
    result=$(read_artifact "engagements/test-customer/CONTEXT.md")
    if ! echo "$result" | grep -q "Test Customer"; then
        echo "FAIL: read_artifact did not return correct content"
        return 1
    fi
    echo "PASS: read_artifact"
}

test_list_artifacts
test_read_artifact
echo "All artifacts tests passed"
