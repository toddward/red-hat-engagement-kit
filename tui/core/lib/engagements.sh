#!/usr/bin/env bash
# engagements.sh - Engagement management

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/../protocol.sh"

PROJECT_ROOT="${PROJECT_ROOT:-$(cd "$SCRIPT_DIR/../../.." && pwd)}"
ENGAGEMENTS_DIR="$PROJECT_ROOT/engagements"

list_engagements() {
    local engagements="[]"

    if [[ ! -d "$ENGAGEMENTS_DIR" ]]; then
        echo "$engagements"
        return
    fi

    for eng_dir in "$ENGAGEMENTS_DIR"/*/; do
        [[ -d "$eng_dir" ]] || continue

        local slug
        slug=$(basename "$eng_dir")
        [[ "$slug" == ".template" ]] && continue

        local has_context="false"
        [[ -f "$eng_dir/CONTEXT.md" ]] && has_context="true"

        engagements=$($JQ -c \
            --arg slug "$slug" \
            --argjson hasContext "$has_context" \
            '. + [{slug: $slug, hasContext: $hasContext}]' <<< "$engagements")
    done

    echo "$engagements"
}

get_engagement() {
    local slug="$1"
    local eng_dir="$ENGAGEMENTS_DIR/$slug"

    if [[ ! -d "$eng_dir" ]]; then
        echo "null"
        return
    fi

    local has_context="false"
    [[ -f "$eng_dir/CONTEXT.md" ]] && has_context="true"

    $JQ -cn \
        --arg slug "$slug" \
        --argjson hasContext "$has_context" \
        '{slug: $slug, hasContext: $hasContext}'
}

create_engagement() {
    local slug="$1"
    local eng_dir="$ENGAGEMENTS_DIR/$slug"

    if [[ -d "$eng_dir" ]]; then
        echo '{"error": "Engagement already exists"}'
        return 1
    fi

    mkdir -p "$eng_dir"/{discovery,assessments,deliverables}
    get_engagement "$slug"
}
