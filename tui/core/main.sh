#!/usr/bin/env bash
# main.sh - Main event loop for TUI core
set -euo pipefail

CORE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

source "$CORE_DIR/protocol.sh"
source "$CORE_DIR/lib/state.sh"
source "$CORE_DIR/lib/skills.sh"
source "$CORE_DIR/lib/engagements.sh"
source "$CORE_DIR/lib/phase.sh"
source "$CORE_DIR/lib/artifacts.sh"
source "$CORE_DIR/lib/checklists.sh"
source "$CORE_DIR/lib/agents.sh"
source "$CORE_DIR/lib/claude.sh"

handle_command() {
    local message="$1"
    local cmd id args
    cmd=$(echo "$message" | parse_command_field "cmd")
    id=$(echo "$message" | parse_command_field "id")
    args=$(echo "$message" | parse_command_args)

    case "$cmd" in
        init)
            local skills engagements agents state
            skills=$(list_skills)
            engagements=$(list_engagements)
            agents=$(list_agents)
            state=$(get_all_state)
            send_response "$id" "$($JQ -cn \
                --argjson skills "$skills" \
                --argjson engagements "$engagements" \
                --argjson agents "$agents" \
                --argjson state "$state" \
                '{skills: $skills, engagements: $engagements, agents: $agents, state: $state}')"
            ;;
        list_skills)
            send_response "$id" "{\"skills\": $(list_skills)}"
            ;;
        list_engagements)
            send_response "$id" "{\"engagements\": $(list_engagements)}"
            ;;
        get_phase)
            local engagement
            engagement=$(echo "$args" | $JQ -r '.engagement')
            send_response "$id" "$(detect_phase "$engagement")"
            ;;
        list_artifacts)
            local engagement
            engagement=$(echo "$args" | $JQ -r '.engagement')
            send_response "$id" "{\"tree\": $(list_artifacts "$engagement")}"
            ;;
        read_artifact)
            local path
            path=$(echo "$args" | $JQ -r '.path')
            send_response "$id" "$(read_artifact "$path")"
            ;;
        list_checklists)
            send_response "$id" "{\"checklists\": $(list_checklists)}"
            ;;
        get_checklist)
            local name
            name=$(echo "$args" | $JQ -r '.name')
            send_response "$id" "$(get_checklist "$name")"
            ;;
        toggle_checklist)
            local name line
            name=$(echo "$args" | $JQ -r '.name')
            line=$(echo "$args" | $JQ -r '.line')
            send_response "$id" "$(toggle_checklist_item "$name" "$line")"
            ;;
        list_agents)
            send_response "$id" "{\"agents\": $(list_agents)}"
            ;;
        execute_skill)
            local skill engagement
            skill=$(echo "$args" | $JQ -r '.skill')
            engagement=$(echo "$args" | $JQ -r '.engagement // empty')
            execute_skill "$skill" "$engagement"
            ;;
        execute_agent)
            local agent prompt engagement
            agent=$(echo "$args" | $JQ -r '.agent')
            prompt=$(echo "$args" | $JQ -r '.prompt')
            engagement=$(echo "$args" | $JQ -r '.engagement // empty')
            execute_agent "$agent" "$prompt" "$engagement"
            ;;
        user_input)
            local text
            text=$(echo "$args" | $JQ -r '.text')
            send_user_input "$text"
            ;;
        cancel)
            cancel_execution
            ;;
        set_state)
            local key value
            key=$(echo "$args" | $JQ -r '.key')
            value=$(echo "$args" | $JQ -r '.value')
            set_state "$key" "$value"
            send_response "$id" '{"success": true}'
            ;;
        *)
            send_error "Unknown command: $cmd" "$id"
            ;;
    esac
}

main() {
    init_state
    while IFS= read -r message; do
        [[ -z "$message" ]] && continue
        handle_command "$message"
    done
}

if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main
fi
