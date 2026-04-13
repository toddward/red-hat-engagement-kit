#!/usr/bin/env bats

setup() {
    load 'test_helper/common-setup'
    _common_setup
    export STATE_FILE="$(mktemp)"
    source "$CORE_DIR/main.sh"
}

teardown() {
    rm -f "$STATE_FILE"
}

@test "init command returns skills, engagements, agents, state" {
    init_state
    run handle_command '{"cmd":"init","id":"cmd-1","args":{}}'
    assert_success
    assert_output --partial '"type":"response"'
    assert_output --partial '"id":"cmd-1"'
    assert_output --partial '"skills"'
    assert_output --partial '"engagements"'
}

@test "list_skills returns skills" {
    run handle_command '{"cmd":"list_skills","id":"cmd-2","args":{}}'
    assert_success
    assert_output --partial '"test-skill"'
}

@test "get_phase returns live for test-customer" {
    run handle_command '{"cmd":"get_phase","id":"cmd-3","args":{"engagement":"test-customer"}}'
    assert_success
    assert_output --partial '"phase":"live"'
}

@test "get_phase returns pre-engagement for missing" {
    run handle_command '{"cmd":"get_phase","id":"cmd-4","args":{"engagement":"nonexistent"}}'
    assert_success
    assert_output --partial '"phase":"pre-engagement"'
}

@test "unknown command returns error" {
    run handle_command '{"cmd":"bogus","id":"cmd-5","args":{}}'
    assert_success
    assert_output --partial '"type":"error"'
    assert_output --partial '"Unknown command: bogus"'
}

@test "set_state returns success" {
    init_state
    run handle_command '{"cmd":"set_state","id":"cmd-6","args":{"key":"test","value":"hello"}}'
    assert_success
    assert_output --partial '"success":true'
}

@test "list_engagements returns engagements" {
    run handle_command '{"cmd":"list_engagements","id":"cmd-7","args":{}}'
    assert_success
    assert_output --partial '"test-customer"'
}

@test "list_agents returns agents" {
    run handle_command '{"cmd":"list_agents","id":"cmd-8","args":{}}'
    assert_success
    assert_output --partial '"test-agent"'
}

@test "list_checklists returns checklists" {
    run handle_command '{"cmd":"list_checklists","id":"cmd-9","args":{}}'
    assert_success
    assert_output --partial '"test-checklist"'
}

@test "read_artifact returns content" {
    run handle_command '{"cmd":"read_artifact","id":"cmd-10","args":{"path":"engagements/test-customer/CONTEXT.md"}}'
    assert_success
    assert_output --partial 'Test Customer'
}
