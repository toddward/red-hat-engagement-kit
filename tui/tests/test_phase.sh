#!/usr/bin/env bash
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
export PROJECT_ROOT="$SCRIPT_DIR/fixtures"
source "$SCRIPT_DIR/../core/lib/phase.sh"

test_detect_phase() {
    local result
    result=$(detect_phase "test-customer")
    if ! echo "$result" | grep -q '"phase":"live"'; then
        echo "FAIL: detect_phase should return 'live' for engagement with CONTEXT.md"
        echo "Got: $result"
        return 1
    fi
    echo "PASS: detect_phase"
}

test_detect_phase_missing() {
    local result
    result=$(detect_phase "nonexistent")
    if ! echo "$result" | grep -q '"phase":"pre-engagement"'; then
        echo "FAIL: detect_phase should return 'pre-engagement' for missing engagement"
        return 1
    fi
    echo "PASS: detect_phase_missing"
}

test_detect_phase
test_detect_phase_missing
echo "All phase tests passed"
