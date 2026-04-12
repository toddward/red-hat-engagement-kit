#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
export PROJECT_ROOT="$SCRIPT_DIR/fixtures"
source "$SCRIPT_DIR/../core/lib/skills.sh"

test_list_skills() {
    local result
    result=$(list_skills)

    if ! echo "$result" | grep -q '"name":"test-skill"'; then
        echo "FAIL: list_skills did not find test-skill"
        echo "Got: $result"
        return 1
    fi
    echo "PASS: list_skills"
}

test_get_skill() {
    local result
    result=$(get_skill "test-skill")

    if ! echo "$result" | grep -q '"description"'; then
        echo "FAIL: get_skill description missing"
        echo "Got: $result"
        return 1
    fi
    echo "PASS: get_skill"
}

test_get_skill_not_found() {
    local result
    result=$(get_skill "nonexistent" 2>&1 || true)

    if ! echo "$result" | grep -q "null"; then
        echo "FAIL: get_skill should return null for missing skill"
        return 1
    fi
    echo "PASS: get_skill_not_found"
}

test_list_skills
test_get_skill
test_get_skill_not_found

echo "All skills tests passed"
