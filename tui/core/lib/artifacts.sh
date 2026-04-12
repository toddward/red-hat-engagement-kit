#!/usr/bin/env bash
# artifacts.sh - Artifact browsing
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/../protocol.sh"

PROJECT_ROOT="${PROJECT_ROOT:-$(cd "$SCRIPT_DIR/../../.." && pwd)}"
ENGAGEMENTS_DIR="$PROJECT_ROOT/engagements"

_build_tree() {
    local dir="$1" rel_path="$2"
    local result="[]"

    for entry in "$dir"/*; do
        [[ -e "$entry" ]] || continue
        local name
        name=$(basename "$entry")
        local entry_path="$rel_path/$name"

        if [[ -d "$entry" ]]; then
            local children
            children=$(_build_tree "$entry" "$entry_path")
            result=$($JQ -c --arg name "$name" --arg path "$entry_path" --argjson children "$children" \
                '. + [{name: $name, path: $path, type: "directory", children: $children}]' <<< "$result")
        else
            result=$($JQ -c --arg name "$name" --arg path "$entry_path" \
                '. + [{name: $name, path: $path, type: "file"}]' <<< "$result")
        fi
    done
    echo "$result"
}

list_artifacts() {
    local slug="$1"
    local eng_dir="$ENGAGEMENTS_DIR/$slug"
    [[ ! -d "$eng_dir" ]] && { echo "[]"; return; }
    _build_tree "$eng_dir" "engagements/$slug"
}

read_artifact() {
    local path="$1"
    local full_path="$PROJECT_ROOT/$path"
    if [[ ! -f "$full_path" ]]; then
        $JQ -cn --arg path "$path" '{error: "File not found", path: $path}'
        return 1
    fi
    local content
    content=$(cat "$full_path")
    $JQ -cn --arg content "$content" '{content: $content}'
}
