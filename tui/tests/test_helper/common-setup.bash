# common-setup.bash — sourced by every .bats file via setup()
_common_setup() {
    load 'test_helper/bats-support/load'
    load 'test_helper/bats-assert/load'
    load 'test_helper/bats-file/load'

    TESTS_DIR="$(cd "$(dirname "$BATS_TEST_FILENAME")" && pwd)"
    TUI_DIR="$(cd "$TESTS_DIR/.." && pwd)"
    CORE_DIR="$TUI_DIR/core"

    export PROJECT_ROOT="$TESTS_DIR/fixtures"
}
