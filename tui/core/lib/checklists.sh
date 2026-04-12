#!/usr/bin/env bash
# checklists.sh - Checklist parsing and toggling
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/../protocol.sh"

PROJECT_ROOT="${PROJECT_ROOT:-$(cd "$SCRIPT_DIR/../../.." && pwd)}"
CHECKLISTS_DIR="$PROJECT_ROOT/knowledge/checklists"

list_checklists() {
    local checklists="[]"
    [[ ! -d "$CHECKLISTS_DIR" ]] && { echo "$checklists"; return; }

    for file in "$CHECKLISTS_DIR"/*.md; do
        [[ -f "$file" ]] || continue
        local filename
        filename=$(basename "$file")
        local name="${filename%.md}"
        checklists=$($JQ -c --arg name "$name" --arg fileName "$filename" \
            '. + [{name: $name, fileName: $fileName}]' <<< "$checklists")
    done
    echo "$checklists"
}

get_checklist() {
    local name="$1"
    local file="$CHECKLISTS_DIR/$name.md"
    [[ ! -f "$file" ]] && { echo "null"; return; }

    local sections="[]" current_section="" current_items="[]"
    local line_num=0 total_items=0 checked_items=0

    while IFS= read -r line || [[ -n "$line" ]]; do
        line_num=$(( line_num + 1 ))
        if [[ "$line" =~ ^##[[:space:]]+(.+)$ ]]; then
            if [[ -n "$current_section" ]]; then
                sections=$($JQ -c --arg title "$current_section" --argjson items "$current_items" \
                    '. + [{title: $title, items: $items}]' <<< "$sections")
            fi
            current_section="${BASH_REMATCH[1]}"
            current_items="[]"
            continue
        fi
        if [[ "$line" =~ ^-[[:space:]]+\[([xX[:space:]])\][[:space:]]+(.+)$ ]]; then
            local checkbox_state="${BASH_REMATCH[1]}"
            local text="${BASH_REMATCH[2]}"
            local checked="false"
            [[ "$checkbox_state" =~ [xX] ]] && checked="true" || true
            total_items=$(( total_items + 1 ))
            [[ "$checked" == "true" ]] && checked_items=$(( checked_items + 1 )) || true
            current_items=$($JQ -c --arg text "$text" --argjson checked "$checked" --argjson line "$line_num" \
                '. + [{text: $text, checked: $checked, line: $line}]' <<< "$current_items")
        fi
    done < "$file"

    if [[ -n "$current_section" ]]; then
        sections=$($JQ -c --arg title "$current_section" --argjson items "$current_items" \
            '. + [{title: $title, items: $items}]' <<< "$sections")
    fi

    local completion_pct=0
    if [[ "$total_items" -gt 0 ]]; then
        completion_pct=$(( (checked_items * 100) / total_items ))
    fi

    $JQ -cn --arg name "$name" --arg fileName "$name.md" --argjson sections "$sections" --argjson pct "$completion_pct" \
        '{name: $name, fileName: $fileName, sections: $sections, completionPercent: $pct}'
}

toggle_checklist_item() {
    local name="$1" line_num="$2"
    local file="$CHECKLISTS_DIR/$name.md"
    [[ ! -f "$file" ]] && { echo '{"success": false, "error": "Checklist not found"}'; return 1; }

    local line
    line=$(sed -n "${line_num}p" "$file")
    if [[ "$line" =~ \[\ \] ]]; then
        sed -i.bak "${line_num}s/\[ \]/[x]/" "$file"
        rm -f "$file.bak"
    elif [[ "$line" =~ \[[xX]\] ]]; then
        sed -i.bak "${line_num}s/\[[xX]\]/[ ]/" "$file"
        rm -f "$file.bak"
    else
        echo '{"success": false, "error": "Line is not a checkbox"}'
        return 1
    fi
    echo '{"success": true}'
}
