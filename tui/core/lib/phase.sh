#!/usr/bin/env bash
# phase.sh - Phase detection
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/../protocol.sh"

PROJECT_ROOT="${PROJECT_ROOT:-$(cd "$SCRIPT_DIR/../../.." && pwd)}"
ENGAGEMENTS_DIR="$PROJECT_ROOT/engagements"

_count_files() {
    local dir="$1"
    if [[ -d "$dir" ]]; then
        find "$dir" -maxdepth 1 -type f | wc -l | tr -d ' '
    else
        echo "0"
    fi
}

detect_phase() {
    local slug="$1"
    local eng_dir="$ENGAGEMENTS_DIR/$slug"
    local phase="pre-engagement"
    local discovery_count=0 assessments_count=0 deliverables_count=0

    if [[ -d "$eng_dir" ]]; then
        discovery_count=$(_count_files "$eng_dir/discovery")
        assessments_count=$(_count_files "$eng_dir/assessments")
        deliverables_count=$(_count_files "$eng_dir/deliverables")

        if [[ -f "$eng_dir/CONTEXT.md" ]]; then
            if (( deliverables_count > 0 )); then
                phase="leave-behind"
            else
                phase="live"
            fi
        fi
    fi

    $JQ -cn \
        --arg phase "$phase" \
        --argjson discovery "$discovery_count" \
        --argjson assessments "$assessments_count" \
        --argjson deliverables "$deliverables_count" \
        '{phase: $phase, artifactCounts: {discovery: $discovery, assessments: $assessments, deliverables: $deliverables}}'
}
