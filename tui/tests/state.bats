#!/usr/bin/env bats

setup() {
    load 'test_helper/common-setup'
    _common_setup
    source "$CORE_DIR/lib/state.sh"
    export STATE_FILE="$(mktemp)"
}

teardown() {
    rm -f "$STATE_FILE"
}

@test "init_state creates file" {
    init_state
    assert_file_exists "$STATE_FILE"
}

@test "get/set state round-trips a value" {
    init_state
    set_state "lastEngagement" "acme-corp"
    run get_state "lastEngagement"
    assert_success
    assert_output "acme-corp"
}

@test "get_all_state returns all state" {
    init_state
    set_state "foo" "bar"
    run get_all_state
    assert_success
    assert_output --partial '"foo"'
}
