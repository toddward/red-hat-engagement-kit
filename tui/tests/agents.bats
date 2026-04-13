#!/usr/bin/env bats

setup() {
    load 'test_helper/common-setup'
    _common_setup
    source "$CORE_DIR/lib/agents.sh"
}

@test "list_agents finds test-agent" {
    run list_agents
    assert_success
    assert_output --partial '"name":"test-agent"'
}

@test "list_agents parses model" {
    run list_agents
    assert_output --partial '"model":"sonnet"'
}
