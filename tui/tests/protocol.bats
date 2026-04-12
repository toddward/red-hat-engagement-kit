#!/usr/bin/env bats

setup() {
    load 'test_helper/common-setup'
    _common_setup
    source "$CORE_DIR/protocol.sh"
}

# --- send_response ---

@test "send_response produces correct JSON envelope" {
    run send_response "test-123" '{"foo":"bar"}'
    assert_success
    assert_output '{"type":"response","id":"test-123","payload":{"foo":"bar"}}'
}

# --- send_event ---

@test "send_event produces correct JSON envelope" {
    run send_event "assistant" '{"text":"hello"}'
    assert_success
    assert_output '{"type":"event","payload":{"event":"assistant","text":"hello"}}'
}

# --- send_error ---

@test "send_error with id includes correlation id" {
    run send_error "something went wrong" "err-456"
    assert_success
    assert_output '{"type":"error","id":"err-456","payload":{"message":"something went wrong"}}'
}

@test "send_error without id omits id field" {
    run send_error "something went wrong"
    assert_success
    assert_output '{"type":"error","payload":{"message":"something went wrong"}}'
    refute_output --partial '"id"'
}

# --- parse_command_field ---

@test "parse_command_field extracts cmd field" {
    local json='{"cmd":"list_skills","id":"abc-123","args":{}}'
    run bash -c "source '$CORE_DIR/protocol.sh' && echo '$json' | parse_command_field cmd"
    assert_success
    assert_output "list_skills"
}

@test "parse_command_field extracts id field" {
    local json='{"cmd":"list_skills","id":"abc-123","args":{}}'
    run bash -c "source '$CORE_DIR/protocol.sh' && echo '$json' | parse_command_field id"
    assert_success
    assert_output "abc-123"
}

# --- parse_command_args ---

@test "parse_command_args extracts args object" {
    local json='{"cmd":"run","id":"x","args":{"skill":"setup"}}'
    run bash -c "source '$CORE_DIR/protocol.sh' && echo '$json' | parse_command_args"
    assert_success
    assert_output '{"skill":"setup"}'
}

@test "parse_command_args with missing args returns empty object" {
    local json='{"cmd":"run","id":"x"}'
    run bash -c "source '$CORE_DIR/protocol.sh' && echo '$json' | parse_command_args"
    assert_success
    assert_output '{}'
}
