#!/usr/bin/env bash
# skills.sh - Skill discovery and reading

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/../protocol.sh"

PROJECT_ROOT="${PROJECT_ROOT:-$(cd "$SCRIPT_DIR/../../.." && pwd)}"
SKILLS_DIR="$PROJECT_ROOT/.claude/skills"

_extract_description() {
    local file="$1"
    local frontmatter
    frontmatter=$(sed -n '/^---$/,/^---$/p' "$file")

    # Check if description is single-line or multi-line (YAML > or |)
    local first_line
    first_line=$(echo "$frontmatter" | grep -E '^description:' | sed 's/^description:[[:space:]]*//')

    if [[ "$first_line" == ">" || "$first_line" == "|" || -z "$first_line" ]]; then
        # Multi-line: grab indented lines after 'description:'
        echo "$frontmatter" | sed -n '/^description:/,/^[^ ]/p' | tail -n +2 | grep -E '^  ' | sed 's/^  //' | tr '\n' ' ' | sed 's/[[:space:]]*$//'
    else
        echo "$first_line"
    fi
}

list_skills() {
    local skills="[]"

    if [[ ! -d "$SKILLS_DIR" ]]; then
        echo "$skills"
        return
    fi

    for skill_dir in "$SKILLS_DIR"/*/; do
        [[ -d "$skill_dir" ]] || continue

        local skill_file="$skill_dir/SKILL.md"
        [[ -f "$skill_file" ]] || continue

        local name
        name=$(basename "$skill_dir")
        local description
        description=$(_extract_description "$skill_file")
        local path
        path=".claude/skills/$name/SKILL.md"

        skills=$($JQ -c \
            --arg name "$name" \
            --arg desc "$description" \
            --arg path "$path" \
            '. + [{name: $name, description: $desc, path: $path}]' <<< "$skills")
    done

    echo "$skills"
}

get_skill() {
    local name="$1"
    local skill_file="$SKILLS_DIR/$name/SKILL.md"

    if [[ ! -f "$skill_file" ]]; then
        echo "null"
        return
    fi

    local description
    description=$(_extract_description "$skill_file")
    local content
    content=$(cat "$skill_file")
    local path=".claude/skills/$name/SKILL.md"

    $JQ -cn \
        --arg name "$name" \
        --arg desc "$description" \
        --arg path "$path" \
        --arg content "$content" \
        '{name: $name, description: $desc, path: $path, content: $content}'
}
