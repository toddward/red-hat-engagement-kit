#!/usr/bin/env bats

setup() {
    load 'test_helper/common-setup'
    _common_setup
    source "$CORE_DIR/lib/claude.sh"
}

@test "system event captures session_id" {
    run _parse_claude_event '{"type":"system","session_id":"sess-abc123"}'
    assert_success
    # No output expected — just captures session ID
}

@test "assistant text emits assistant event" {
    run _parse_claude_event '{"type":"assistant","message":{"content":[{"type":"text","text":"Hello world"}]}}'
    assert_success
    assert_output --partial '"event":"assistant"'
    assert_output --partial '"text":"Hello world"'
}

@test "assistant text ending with ? also emits question event" {
    run _parse_claude_event '{"type":"assistant","message":{"content":[{"type":"text","text":"Which file?"}]}}'
    assert_success
    assert_output --partial '"event":"assistant"'
    assert_output --partial '"event":"question"'
}

@test "assistant tool_use emits tool_use event" {
    run _parse_claude_event '{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Read","input":{"file":"foo.txt"}}]}}'
    assert_success
    assert_output --partial '"event":"tool_use"'
    assert_output --partial '"tool":"Read"'
}

@test "AskUserQuestion tool emits question event" {
    run _parse_claude_event '{"type":"assistant","message":{"content":[{"type":"tool_use","name":"AskUserQuestion","input":{"question":"Pick one"}}]}}'
    assert_success
    assert_output --partial '"event":"question"'
    assert_output --partial '"text":"Pick one"'
}

@test "result with cost emits complete event" {
    run _parse_claude_event '{"type":"result","total_cost_usd":0.05,"result":"Done."}'
    assert_success
    assert_output --partial '"event":"complete"'
    assert_output --partial '"status":"success"'
}

@test "result with question mark emits question instead of complete" {
    run _parse_claude_event '{"type":"result","total_cost_usd":0.01,"result":"Which environment?"}'
    assert_success
    assert_output --partial '"event":"question"'
}

@test "result with tool_name emits tool_result event" {
    run _parse_claude_event '{"type":"result","tool_name":"Read","message":{"content":[{"type":"text","text":"file contents"}]}}'
    assert_success
    assert_output --partial '"event":"tool_result"'
    assert_output --partial '"tool":"Read"'
}
