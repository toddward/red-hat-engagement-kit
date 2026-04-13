#!/usr/bin/env bash
set -euo pipefail
TESTS_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
"$TESTS_DIR/bats/bin/bats" "$TESTS_DIR"/*.bats "$@"
