#!/usr/bin/env bats

setup() {
    load 'test_helper/common-setup'
    _common_setup
    source "$CORE_DIR/lib/checklists.sh"
}

@test "list_checklists finds test-checklist" {
    run list_checklists
    assert_success
    assert_output --partial '"name":"test-checklist"'
}

@test "get_checklist parses sections" {
    run get_checklist "test-checklist"
    assert_success
    assert_output --partial '"title":"Section One"'
}

@test "get_checklist detects checked items" {
    run get_checklist "test-checklist"
    assert_output --partial '"checked":true'
}
