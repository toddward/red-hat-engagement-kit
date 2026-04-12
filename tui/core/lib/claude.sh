#!/usr/bin/env bash
# claude.sh - Claude CLI execution and streaming

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/../protocol.sh"

PROJECT_ROOT="${PROJECT_ROOT:-$(cd "$SCRIPT_DIR/../../.." && pwd)}"

CLAUDE_PID=""

_parse_claude_event() {
    local line="$1"
    local type
    type=$($JQ -r '.type // empty' <<< "$line" 2>/dev/null) || return

    case "$type" in
        assistant)
            local text
            text=$($JQ -r '.message.content[]? | select(.type == "text") | .text // empty' <<< "$line" 2>/dev/null)
            if [[ -n "$text" ]]; then
                send_event "assistant" "$($JQ -cn --arg text "$text" '{text: $text}')"
            fi
            local tool_uses
            tool_uses=$($JQ -c '.message.content[]? | select(.type == "tool_use")' <<< "$line" 2>/dev/null)
            if [[ -n "$tool_uses" ]]; then
                while IFS= read -r tool_use; do
                    local tool_name input
                    tool_name=$($JQ -r '.name' <<< "$tool_use")
                    input=$($JQ -c '.input' <<< "$tool_use")
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
                send_event "complete" "$($JQ -cn --arg status "success" --argjson cost "$cost" '{status: $status, totalCost: $cost}')"
            fi
            ;;
        system) ;;
        *) ;;
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
    claude --print \
           --output-format stream-json \
           --verbose \
           --permission-mode acceptEdits \
           "$prompt" 2>/dev/null &
    CLAUDE_PID=$!

    while IFS= read -r line; do
        [[ -z "$line" ]] && continue
        _parse_claude_event "$line"
    done < <(wait $CLAUDE_PID 2>/dev/null; true)
    CLAUDE_PID=""
}

cancel_execution() {
    if [[ -n "$CLAUDE_PID" ]] && kill -0 "$CLAUDE_PID" 2>/dev/null; then
        kill -TERM "$CLAUDE_PID" 2>/dev/null
        CLAUDE_PID=""
        send_event "complete" '{"status": "cancelled"}'
    fi
}

send_user_input() {
    local text="$1"
    send_error "User input during execution not yet implemented"
}
