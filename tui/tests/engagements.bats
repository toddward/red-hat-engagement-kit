#!/usr/bin/env bats

setup() {
    load 'test_helper/common-setup'
    _common_setup
    source "$CORE_DIR/lib/engagements.sh"
}

@test "list_engagements finds test-customer" {
    run list_engagements
    assert_success
    assert_output --partial '"slug":"test-customer"'
}

@test "list_engagements includes hasContext" {
    run list_engagements
    assert_output --partial '"hasContext":true'
}

@test "get_engagement returns correct slug" {
    run get_engagement "test-customer"
    assert_success
    assert_output --partial '"slug":"test-customer"'
}
