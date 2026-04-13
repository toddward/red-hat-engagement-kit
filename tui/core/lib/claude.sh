#!/usr/bin/env bash
# claude.sh - Claude CLI execution and streaming

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/../protocol.sh"

PROJECT_ROOT="${PROJECT_ROOT:-$(cd "$SCRIPT_DIR/../../.." && pwd)}"

# Track Claude session ID for resumption
CLAUDE_SESSION_ID=""
# Track whether a complete event was sent this execution
_SENT_COMPLETE=""
# Track whether we're waiting for user input (AskUserQuestion pending)
_WAITING_FOR_INPUT=""

_parse_claude_event() {
    local line="$1"
    local type
    type=$($JQ -r '.type // empty' <<< "$line" 2>/dev/null) || return

    case "$type" in
        system)
            # Capture session ID for later resumption
            local session_id
            session_id=$($JQ -r '.session_id // empty' <<< "$line" 2>/dev/null)
            if [[ -n "$session_id" ]]; then
                CLAUDE_SESSION_ID="$session_id"
            fi
            ;;
        assistant)
            local text
            text=$($JQ -r '.message.content[]? | select(.type == "text") | .text // empty' <<< "$line" 2>/dev/null)
            if [[ -n "$text" ]]; then
                send_event "assistant" "$($JQ -cn --arg text "$text" '{text: $text}')"
                # Check if text ends with a question mark
                if [[ "$text" =~ \?[[:space:]]*$ ]]; then
                    send_event "question" "$($JQ -cn --arg text "$text" '{text: $text}')"
                fi
            fi
            local tool_uses
            tool_uses=$($JQ -c '.message.content[]? | select(.type == "tool_use")' <<< "$line" 2>/dev/null)
            if [[ -n "$tool_uses" ]]; then
                while IFS= read -r tool_use; do
                    [[ -z "$tool_use" ]] && continue
                    local tool_name input
                    tool_name=$($JQ -r '.name' <<< "$tool_use")
                    input=$($JQ -c '.input' <<< "$tool_use")
                    # AskUserQuestion triggers input mode
                    if [[ "$tool_name" == "AskUserQuestion" ]]; then
                        local question options_json
                        # AskUserQuestion schema: .input.questions[0].question for text
                        question=$($JQ -r '.input.questions[0].question // .input.question // .input.text // empty' <<< "$tool_use")
                        # Extract option labels if present
                        options_json=$($JQ -c '[.input.questions[0].options[]?.label // empty] // []' <<< "$tool_use" 2>/dev/null)
                        [[ "$options_json" == "null" || -z "$options_json" ]] && options_json="[]"
                        send_event "question" "$($JQ -cn --arg text "$question" --argjson options "$options_json" '{text: $text, options: $options}')"
                        # Don't send complete while waiting for user input
                        _WAITING_FOR_INPUT="true"
                    fi
                    send_event "tool_use" "$($JQ -cn --arg tool "$tool_name" --arg input "$input" '{tool: $tool, input: $input}')"
                done <<< "$tool_uses"
            fi
            ;;
        result)
            local tool_name
            tool_name=$($JQ -r '.tool_name // empty' <<< "$line" 2>/dev/null)
            if [[ -n "$tool_name" ]]; then
                local output
                output=$($JQ -r '.message.content[]? | select(.type == "text") | .text // empty' <<< "$line" 2>/dev/null)
                send_event "tool_result" "$($JQ -cn --arg tool "$tool_name" --arg output "$output" '{tool: $tool, output: $output}')"
            fi
            local cost
            cost=$($JQ -r '.total_cost_usd // empty' <<< "$line" 2>/dev/null)
            if [[ -n "$cost" ]]; then
                local result_text
                result_text=$($JQ -r '.result // empty' <<< "$line" 2>/dev/null)
                # Check if result ends with question
                if [[ "$result_text" =~ \?[[:space:]]*$ ]]; then
                    send_event "question" "$($JQ -cn --arg text "$result_text" '{text: $text}')"
                else
                    _SENT_COMPLETE="true"
                    send_event "complete" "$($JQ -cn --arg status "success" --argjson cost "$cost" '{status: $status, totalCost: $cost}')"
                fi
            fi
            ;;
        *)
            ;;
    esac
}

execute_skill() {
    local skill="$1" engagement="$2"
    local prompt="Run /${skill}."
    [[ -n "$engagement" ]] && prompt="$prompt Use engagement at engagements/${engagement}/."
    _execute_claude "$prompt"
}

execute_agent() {
    local agent="$1" prompt="$2" engagement="$3"
    local full_prompt="Using the $agent agent: $prompt"
    [[ -n "$engagement" ]] && full_prompt="$full_prompt (Engagement: $engagement)"
    _execute_claude "$full_prompt"
}

_execute_claude() {
    local prompt="$1"

    # Use process substitution to capture Claude's output
    # Claude's stdout is streamed line-by-line to the parser
    while IFS= read -r line; do
        [[ -z "$line" ]] && continue
        _parse_claude_event "$line"
    done < <(claude --print \
           --output-format stream-json \
           --verbose \
           --permission-mode acceptEdits \
           "$prompt" 2>/dev/null)

    # If we got here without a complete event, send one (unless waiting for user input)
    if [[ -z "$_SENT_COMPLETE" && -z "$_WAITING_FOR_INPUT" ]]; then
        send_event "complete" '{"status": "success", "totalCost": 0}'
    fi
    _SENT_COMPLETE=""
}

cancel_execution() {
    # Send cancel signal to any running claude process
    pkill -f "claude --print" 2>/dev/null || true
    _WAITING_FOR_INPUT=""
    send_event "complete" '{"status": "cancelled", "totalCost": 0}'
}

send_user_input() {
    local text="$1"
    if [[ -z "$CLAUDE_SESSION_ID" ]]; then
        send_error "No active session to resume"
        return
    fi
    # Clear waiting flag - user has provided input
    _WAITING_FOR_INPUT=""
    # Resume the Claude session with the user's response
    _execute_claude_resume "$text"
}

_execute_claude_resume() {
    local message="$1"
    _SENT_COMPLETE=""
    while IFS= read -r line; do
        [[ -z "$line" ]] && continue
        _parse_claude_event "$line"
    done < <(claude --print \
           --output-format stream-json \
           --verbose \
           --permission-mode acceptEdits \
           --resume "$CLAUDE_SESSION_ID" \
           "$message" 2>/dev/null)

    # If we got here without a complete event, send one (unless waiting for more input)
    if [[ -z "$_SENT_COMPLETE" && -z "$_WAITING_FOR_INPUT" ]]; then
        send_event "complete" '{"status": "success", "totalCost": 0}'
    fi
    _SENT_COMPLETE=""
}
