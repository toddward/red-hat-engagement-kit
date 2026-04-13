#!/usr/bin/env bats
# protocol_integration.bats - Full round-trip integration tests for main.sh
#
# These tests start main.sh as a background process connected via named pipes,
# then send JSON commands and assert on the JSON responses. This exercises the
# full bash-side command dispatch without mocking any internals.
#
# Named pipes (mkfifo) are used instead of coproc because the system bash on
# macOS is version 3.2, which does not support the coproc builtin.

setup() {
    load 'test_helper/common-setup'
    _common_setup

    # Isolated state file for each test
    export STATE_FILE
    STATE_FILE="$(mktemp)"

    # Named pipes for bidirectional communication with main.sh
    _IN_PIPE="$(mktemp -u)"
    _OUT_PIPE="$(mktemp -u)"
    mkfifo "$_IN_PIPE"
    mkfifo "$_OUT_PIPE"

    # Launch main.sh as a background process connected to both pipes.
    # STATE_FILE and PROJECT_ROOT are already exported so main.sh picks them up.
    bash "$CORE_DIR/main.sh" <"$_IN_PIPE" >"$_OUT_PIPE" &
    _MAIN_PID=$!

    # Open persistent file descriptors so the pipes stay open across
    # multiple reads/writes within a single test.
    # Use fixed FD numbers (7 and 8) because bash 3.2 on macOS does not
    # support the brace-based dynamic allocation syntax (exec {var}>file).
    exec 8>"$_IN_PIPE"   # write commands to main.sh
    exec 7<"$_OUT_PIPE"  # read responses from main.sh
}

teardown() {
    # Close the write end first — this signals EOF to main.sh so it exits
    # its read loop cleanly.
    exec 8>&- 2>/dev/null || true
    exec 7<&- 2>/dev/null || true
    wait "$_MAIN_PID" 2>/dev/null || true
    rm -f "$_IN_PIPE" "$_OUT_PIPE" "$STATE_FILE"
}

# _send_recv <json_command>
#   Sends one JSON command to main.sh and reads one line of response.
#   Times out after 5 seconds to prevent hanging tests.
_send_recv() {
    local cmd="$1"
    local response
    echo "$cmd" >&8
    IFS= read -r -t 5 response <&7
    echo "$response"
}

# ---------------------------------------------------------------------------
# Tests
# ---------------------------------------------------------------------------

@test "round-trip: init returns valid JSON with expected fields" {
    local response
    response=$(_send_recv '{"cmd":"init","id":"rt-1","args":{}}')

    echo "$response" | jq -e '.type == "response"'
    echo "$response" | jq -e '.id == "rt-1"'
    echo "$response" | jq -e '.payload.skills | type == "array"'
    echo "$response" | jq -e '.payload.engagements | type == "array"'
    echo "$response" | jq -e '.payload.agents | type == "array"'
    echo "$response" | jq -e '.payload.state | type == "object"'
}

@test "round-trip: list_skills returns skills array" {
    local response
    response=$(_send_recv '{"cmd":"list_skills","id":"rt-2","args":{}}')

    echo "$response" | jq -e '.payload.skills | length > 0'
    echo "$response" | jq -e '.payload.skills[0].name == "test-skill"'
}

@test "round-trip: sequential commands maintain state" {
    _send_recv '{"cmd":"set_state","id":"rt-3a","args":{"key":"lastEngagement","value":"test-customer"}}' > /dev/null

    local response
    response=$(_send_recv '{"cmd":"init","id":"rt-3b","args":{}}')
    echo "$response" | jq -e '.payload.state.lastEngagement == "test-customer"'
}

@test "round-trip: unknown command returns error" {
    local response
    response=$(_send_recv '{"cmd":"nonexistent","id":"rt-5","args":{}}')

    echo "$response" | jq -e '.type == "error"'
    echo "$response" | jq -e '.id == "rt-5"'
}

@test "round-trip: get_phase returns phase info" {
    local response
    response=$(_send_recv '{"cmd":"get_phase","id":"rt-6","args":{"engagement":"test-customer"}}')

    echo "$response" | jq -e '.payload.phase == "live"'
}

@test "round-trip: multiple commands get individual responses" {
    # Send three commands back-to-back without waiting for responses
    echo '{"cmd":"list_skills","id":"multi-1","args":{}}' >&8
    echo '{"cmd":"list_agents","id":"multi-2","args":{}}' >&8
    echo '{"cmd":"list_checklists","id":"multi-3","args":{}}' >&8

    local r1 r2 r3
    IFS= read -r -t 5 r1 <&7
    IFS= read -r -t 5 r2 <&7
    IFS= read -r -t 5 r3 <&7

    echo "$r1" | jq -e '.id == "multi-1"'
    echo "$r2" | jq -e '.id == "multi-2"'
    echo "$r3" | jq -e '.id == "multi-3"'
}
