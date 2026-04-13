#!/usr/bin/env bats

setup() {
    load 'test_helper/common-setup'
    _common_setup
    source "$CORE_DIR/lib/skills.sh"
}

@test "list_skills finds test-skill" {
    run list_skills
    assert_success
    assert_output --partial '"name":"test-skill"'
}

@test "get_skill returns description" {
    run get_skill "test-skill"
    assert_success
    assert_output --partial '"description"'
}

@test "get_skill returns null for nonexistent skill" {
    run get_skill "nonexistent"
    assert_output --partial "null"
}
