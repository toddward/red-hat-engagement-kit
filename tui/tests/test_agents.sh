#!/usr/bin/env bash
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
export PROJECT_ROOT="$SCRIPT_DIR/fixtures"
source "$SCRIPT_DIR/../core/lib/agents.sh"

test_list_agents() {
    local result
    result=$(list_agents)
    if ! echo "$result" | grep -q '"name":"test-agent"'; then
        echo "FAIL: list_agents did not find test-agent"
        echo "Got: $result"
        return 1
    fi
    if ! echo "$result" | grep -q '"model":"sonnet"'; then
        echo "FAIL: list_agents did not parse model"
        return 1
    fi
    echo "PASS: list_agents"
}

test_list_agents
echo "All agents tests passed"
