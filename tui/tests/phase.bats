#!/usr/bin/env bats

setup() {
    load 'test_helper/common-setup'
    _common_setup
    source "$CORE_DIR/lib/phase.sh"
}

@test "detect_phase returns live for engagement with CONTEXT.md" {
    run detect_phase "test-customer"
    assert_success
    assert_output --partial '"phase":"live"'
}

@test "detect_phase returns pre-engagement for missing engagement" {
    run detect_phase "nonexistent"
    assert_success
    assert_output --partial '"phase":"pre-engagement"'
}
