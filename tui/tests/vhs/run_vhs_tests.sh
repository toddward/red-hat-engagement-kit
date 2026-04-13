#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"
SCREENSHOTS_DIR="$SCRIPT_DIR/screenshots"
BASELINES_DIR="$SCRIPT_DIR/baselines"

# Check VHS is installed
if ! command -v vhs &>/dev/null; then
    echo "ERROR: VHS not installed. Install with: brew install charmbracelet/tap/vhs"
    exit 1
fi

mkdir -p "$SCREENSHOTS_DIR"

failed=0
total=0

for tape in "$SCRIPT_DIR"/*.tape; do
    name=$(basename "$tape" .tape)
    total=$((total + 1))
    echo "Running: $name..."

    # Run VHS with project root as working directory
    if ! (cd "$PROJECT_ROOT" && vhs "$tape") 2>/dev/null; then
        echo "  FAIL: VHS execution failed for $name"
        failed=$((failed + 1))
        continue
    fi

    # Compare screenshots against baselines (if baselines exist)
    for screenshot in "$SCREENSHOTS_DIR"/${name}_*.png "$SCREENSHOTS_DIR"/${name}.png; do
        [ -f "$screenshot" ] || continue
        baseline="$BASELINES_DIR/$(basename "$screenshot")"

        if [ ! -f "$baseline" ]; then
            echo "  NEW: $(basename "$screenshot") — no baseline yet"
            continue
        fi

        if command -v compare &>/dev/null; then
            diff_metric=$(compare -metric RMSE "$screenshot" "$baseline" /dev/null 2>&1 | awk -F'[( )]' '{print $1}')
            threshold="50"
            if (( $(echo "$diff_metric > $threshold" | bc -l 2>/dev/null || echo 0) )); then
                echo "  FAIL: $(basename "$screenshot") differs from baseline (RMSE: $diff_metric)"
                failed=$((failed + 1))
            else
                echo "  PASS: $(basename "$screenshot")"
            fi
        else
            echo "  SKIP: no 'compare' tool (install ImageMagick for image diff)"
        fi
    done
done

echo
echo "VHS tests: $total total, $failed failed"
if [ "$failed" -gt 0 ]; then
    exit 1
else
    echo "To update baselines: cp $SCREENSHOTS_DIR/*.png $BASELINES_DIR/"
fi
