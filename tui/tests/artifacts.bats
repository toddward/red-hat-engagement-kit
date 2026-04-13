#!/usr/bin/env bats

setup() {
    load 'test_helper/common-setup'
    _common_setup
    source "$CORE_DIR/lib/artifacts.sh"
}

@test "list_artifacts includes discovery directory" {
    run list_artifacts "test-customer"
    assert_success
    assert_output --partial '"name":"discovery"'
}

@test "read_artifact returns file content" {
    run read_artifact "engagements/test-customer/CONTEXT.md"
    assert_success
    assert_output --partial "Test Customer"
}
