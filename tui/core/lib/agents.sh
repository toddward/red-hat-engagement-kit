#!/usr/bin/env bash
# agents.sh - Agent discovery and invocation
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/../protocol.sh"

PROJECT_ROOT="${PROJECT_ROOT:-$(cd "$SCRIPT_DIR/../../.." && pwd)}"
AGENTS_DIR="$PROJECT_ROOT/.claude/agents"

_extract_frontmatter_field() {
    local file="$1" field="$2"
    local frontmatter
    frontmatter=$(sed -n '/^---$/,/^---$/p' "$file")

    local first_line
    first_line=$(echo "$frontmatter" | grep -E "^${field}:" | sed "s/^${field}:[[:space:]]*//" | tr -d '\n')

    if [[ "$first_line" == ">" || "$first_line" == "|" || -z "$first_line" ]]; then
        echo "$frontmatter" | sed -n "/^${field}:/,/^[^ ]/p" | tail -n +2 | grep -E '^  ' | sed 's/^  //' | tr '\n' ' ' | sed 's/[[:space:]]*$//'
    else
        echo "$first_line"
    fi
}

list_agents() {
    local agents="[]"
    [[ ! -d "$AGENTS_DIR" ]] && { echo "$agents"; return; }

    # Check both patterns: agents/*.md and agents/*/AGENT.md
    for agent_dir in "$AGENTS_DIR"/*/; do
        [[ -d "$agent_dir" ]] || continue
        local file="$agent_dir/AGENT.md"
        [[ -f "$file" ]] || continue

        local name
        name=$(_extract_frontmatter_field "$file" "name")
        [[ -z "$name" ]] && name=$(basename "$agent_dir")
        local model
        model=$(_extract_frontmatter_field "$file" "model")
        local role
        role=$(_extract_frontmatter_field "$file" "role")
        local description
        description=$(_extract_frontmatter_field "$file" "description")

        agents=$($JQ -c --arg name "$name" --arg model "$model" --arg role "$role" --arg desc "$description" \
            '. + [{name: $name, model: $model, role: $role, description: $desc}]' <<< "$agents")
    done
    echo "$agents"
}

get_agent() {
    local name="$1"
    # Check subdirectory pattern first (agents/<name>/AGENT.md), then flat (agents/<name>.md)
    local file="$AGENTS_DIR/$name/AGENT.md"
    [[ ! -f "$file" ]] && file="$AGENTS_DIR/$name.md"
    [[ ! -f "$file" ]] && { echo "null"; return; }

    local model role description content
    model=$(_extract_frontmatter_field "$file" "model")
    role=$(_extract_frontmatter_field "$file" "role")
    description=$(_extract_frontmatter_field "$file" "description")
    content=$(cat "$file")

    $JQ -cn --arg name "$name" --arg model "$model" --arg role "$role" --arg desc "$description" --arg content "$content" \
        '{name: $name, model: $model, role: $role, description: $desc, content: $content}'
}
