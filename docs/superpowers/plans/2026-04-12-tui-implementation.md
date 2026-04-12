# TUI Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a polished terminal UI for the Red Hat Engagement Kit with bash logic layer and Go viewer layer communicating via JSON protocol.

**Architecture:** Hybrid approach — Bash scripts handle business logic (reading skills, engagements, invoking Claude CLI), while a Go binary (Bubbletea) handles terminal rendering. They communicate via JSON messages over stdin/stdout pipes.

**Tech Stack:** Bash 4+, Go 1.21+, Bubbletea, Lipgloss, jq

---

## File Structure

```
tui/
├── bin/                          # Bundled binaries (gitignored)
├── core/
│   ├── main.sh                   # Event loop entry point
│   ├── protocol.sh               # JSON encode/decode helpers
│   └── lib/
│       ├── skills.sh             # Skill discovery and reading
│       ├── engagements.sh        # Engagement management
│       ├── phase.sh              # Phase detection
│       ├── artifacts.sh          # File tree and reading
│       ├── checklists.sh         # Checklist parsing/toggling
│       ├── agents.sh             # Agent definitions
│       ├── claude.sh             # Claude CLI execution
│       └── state.sh              # Persistent state
├── viewer/
│   ├── go.mod
│   ├── main.go
│   ├── protocol/
│   │   └── messages.go           # Protocol types
│   └── ui/
│       ├── styles.go             # Red Hat brand colors
│       ├── sidebar.go            # Sidebar component
│       ├── menu.go               # Menu view
│       ├── activity.go           # Activity log view
│       ├── input.go              # User input view
│       ├── artifacts.go          # Artifact browser view
│       ├── checklists.go         # Checklist view
│       ├── palette.go            # Command palette overlay
│       └── app.go                # Root app model
├── tests/
│   ├── test_skills.sh
│   ├── test_engagements.sh
│   ├── test_protocol.sh
│   └── fixtures/                 # Test data
├── setup.sh
├── tui.sh
└── README.md
```

---

## Phase 1: Foundation & Protocol

### Task 1: Create Directory Structure

**Files:**
- Create: `tui/` directory tree
- Create: `tui/.gitignore`

- [ ] **Step 1: Create directory structure**

```bash
cd "/Users/toddwardzinski/Desktop/devX preso, AI native  installer/red-hat-engagement-kit"
mkdir -p tui/{bin,core/lib,viewer/{protocol,ui},tests/fixtures}
```

- [ ] **Step 2: Create .gitignore for tui/**

Create file `tui/.gitignore`:

```gitignore
# Bundled binaries
bin/

# Go build artifacts
viewer/tui-viewer

# State file (contains engagement-specific data)
.tui-state.json
```

- [ ] **Step 3: Commit**

```bash
git add tui/.gitignore
git commit -m "feat(tui): initialize directory structure"
```

---

### Task 2: Define Protocol Types (Go)

**Files:**
- Create: `tui/viewer/go.mod`
- Create: `tui/viewer/protocol/messages.go`

- [ ] **Step 1: Initialize Go module**

```bash
cd tui/viewer
go mod init github.com/toddward/red-hat-engagement-kit/tui/viewer
```

- [ ] **Step 2: Create protocol message types**

Create file `tui/viewer/protocol/messages.go`:

```go
package protocol

import "encoding/json"

// Command is sent from viewer to bash
type Command struct {
	Cmd  string          `json:"cmd"`
	ID   string          `json:"id,omitempty"`
	Args json.RawMessage `json:"args,omitempty"`
}

// Response is sent from bash to viewer
type Response struct {
	Type    string          `json:"type"` // "response", "event", "error"
	ID      string          `json:"id,omitempty"`
	Payload json.RawMessage `json:"payload"`
}

// Event types for streaming execution
type EventType string

const (
	EventAssistant  EventType = "assistant"
	EventToolUse    EventType = "tool_use"
	EventToolResult EventType = "tool_result"
	EventQuestion   EventType = "question"
	EventCost       EventType = "cost"
	EventComplete   EventType = "complete"
	EventError      EventType = "error"
)

// StreamEvent represents a Claude execution event
type StreamEvent struct {
	Event   EventType `json:"event"`
	Text    string    `json:"text,omitempty"`
	Tool    string    `json:"tool,omitempty"`
	Input   string    `json:"input,omitempty"`
	Output  string    `json:"output,omitempty"`
	Options []string  `json:"options,omitempty"`
	Cost    float64   `json:"cost,omitempty"`
	Status  string    `json:"status,omitempty"`
}

// Skill represents a discovered skill
type Skill struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Path        string `json:"path"`
}

// Engagement represents a customer engagement
type Engagement struct {
	Slug       string `json:"slug"`
	HasContext bool   `json:"hasContext"`
}

// Agent represents a team agent
type Agent struct {
	Name        string `json:"name"`
	Model       string `json:"model"`
	Role        string `json:"role"`
	Description string `json:"description"`
}

// Phase represents engagement phase
type Phase string

const (
	PhasePreEngagement Phase = "pre-engagement"
	PhaseLive          Phase = "live"
	PhaseLeaveBehind   Phase = "leave-behind"
)

// PhaseInfo contains phase detection results
type PhaseInfo struct {
	Phase          Phase          `json:"phase"`
	ArtifactCounts ArtifactCounts `json:"artifactCounts"`
}

// ArtifactCounts tracks files per category
type ArtifactCounts struct {
	Discovery    int `json:"discovery"`
	Assessments  int `json:"assessments"`
	Deliverables int `json:"deliverables"`
}

// ArtifactNode represents a file/directory in the tree
type ArtifactNode struct {
	Name     string         `json:"name"`
	Path     string         `json:"path"`
	Type     string         `json:"type"` // "file" or "directory"
	Children []ArtifactNode `json:"children,omitempty"`
}

// ChecklistItem represents a single checkbox
type ChecklistItem struct {
	Text    string `json:"text"`
	Checked bool   `json:"checked"`
	Line    int    `json:"line"`
}

// ChecklistSection groups items under a heading
type ChecklistSection struct {
	Title string          `json:"title"`
	Items []ChecklistItem `json:"items"`
}

// Checklist represents a complete checklist file
type Checklist struct {
	Name            string             `json:"name"`
	FileName        string             `json:"fileName"`
	Sections        []ChecklistSection `json:"sections"`
	CompletionPct   int                `json:"completionPercent"`
}

// InitResponse is the payload for the init command
type InitResponse struct {
	Skills      []Skill      `json:"skills"`
	Engagements []Engagement `json:"engagements"`
	Agents      []Agent      `json:"agents"`
	State       State        `json:"state"`
}

// State represents persistent TUI state
type State struct {
	LastEngagement string            `json:"lastEngagement,omitempty"`
	Preferences    map[string]string `json:"preferences,omitempty"`
}
```

- [ ] **Step 3: Verify Go compiles**

```bash
cd tui/viewer
go build ./protocol/...
```

Expected: No errors

- [ ] **Step 4: Commit**

```bash
git add tui/viewer/go.mod tui/viewer/protocol/messages.go
git commit -m "feat(tui): add protocol message types for bash-go communication"
```

---

### Task 3: Create Protocol Helpers (Bash)

**Files:**
- Create: `tui/core/protocol.sh`
- Create: `tui/tests/test_protocol.sh`

- [ ] **Step 1: Write protocol test**

Create file `tui/tests/test_protocol.sh`:

```bash
#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/../core/protocol.sh"

# Test send_response
test_send_response() {
    local output
    output=$(send_response "test-123" '{"foo":"bar"}')
    
    local expected='{"type":"response","id":"test-123","payload":{"foo":"bar"}}'
    if [[ "$output" != "$expected" ]]; then
        echo "FAIL: send_response"
        echo "  Expected: $expected"
        echo "  Got: $output"
        return 1
    fi
    echo "PASS: send_response"
}

# Test send_event
test_send_event() {
    local output
    output=$(send_event "assistant" '{"text":"hello"}')
    
    local expected='{"type":"event","payload":{"event":"assistant","text":"hello"}}'
    if [[ "$output" != "$expected" ]]; then
        echo "FAIL: send_event"
        echo "  Expected: $expected"
        echo "  Got: $output"
        return 1
    fi
    echo "PASS: send_event"
}

# Test parse_command
test_parse_command() {
    local cmd='{"cmd":"list_skills","id":"abc-123","args":{}}'
    
    local parsed_cmd parsed_id
    parsed_cmd=$(echo "$cmd" | parse_command_field "cmd")
    parsed_id=$(echo "$cmd" | parse_command_field "id")
    
    if [[ "$parsed_cmd" != "list_skills" ]]; then
        echo "FAIL: parse_command_field cmd"
        return 1
    fi
    if [[ "$parsed_id" != "abc-123" ]]; then
        echo "FAIL: parse_command_field id"
        return 1
    fi
    echo "PASS: parse_command"
}

# Run tests
test_send_response
test_send_event
test_parse_command

echo "All protocol tests passed"
```

- [ ] **Step 2: Run test to verify it fails**

```bash
chmod +x tui/tests/test_protocol.sh
./tui/tests/test_protocol.sh
```

Expected: FAIL with "source: no such file" (protocol.sh doesn't exist)

- [ ] **Step 3: Implement protocol.sh**

Create file `tui/core/protocol.sh`:

```bash
#!/usr/bin/env bash
# protocol.sh - JSON protocol helpers for viewer communication

# Requires jq in PATH or in bin/

# Get jq path (prefer bundled, fall back to system)
_get_jq() {
    local script_dir
    script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
    if [[ -x "$script_dir/bin/jq" ]]; then
        echo "$script_dir/bin/jq"
    elif command -v jq &>/dev/null; then
        echo "jq"
    else
        echo "ERROR: jq not found" >&2
        exit 1
    fi
}

JQ=$(_get_jq)

# Send a response to the viewer
# Usage: send_response <id> <payload_json>
send_response() {
    local id="$1"
    local payload="$2"
    
    $JQ -cn \
        --arg id "$id" \
        --argjson payload "$payload" \
        '{type: "response", id: $id, payload: $payload}'
}

# Send an event to the viewer (no correlation ID)
# Usage: send_event <event_type> <payload_json>
send_event() {
    local event_type="$1"
    local payload="$2"
    
    $JQ -cn \
        --arg event "$event_type" \
        --argjson payload "$payload" \
        '{type: "event", payload: ($payload + {event: $event})}'
}

# Send an error to the viewer
# Usage: send_error <message> [id]
send_error() {
    local message="$1"
    local id="${2:-}"
    
    if [[ -n "$id" ]]; then
        $JQ -cn \
            --arg id "$id" \
            --arg msg "$message" \
            '{type: "error", id: $id, payload: {message: $msg}}'
    else
        $JQ -cn \
            --arg msg "$message" \
            '{type: "error", payload: {message: $msg}}'
    fi
}

# Parse a field from a command JSON
# Usage: echo "$json" | parse_command_field <field>
parse_command_field() {
    local field="$1"
    $JQ -r ".$field // empty"
}

# Parse args from a command JSON
# Usage: echo "$json" | parse_command_args
parse_command_args() {
    $JQ -c ".args // {}"
}
```

- [ ] **Step 4: Run test to verify it passes**

```bash
./tui/tests/test_protocol.sh
```

Expected: "All protocol tests passed"

- [ ] **Step 5: Commit**

```bash
git add tui/core/protocol.sh tui/tests/test_protocol.sh
git commit -m "feat(tui): add bash protocol helpers for JSON communication"
```

---

## Phase 2: Bash Logic Layer

### Task 4: Implement State Management

**Files:**
- Create: `tui/core/lib/state.sh`
- Create: `tui/tests/test_state.sh`

- [ ] **Step 1: Write state test**

Create file `tui/tests/test_state.sh`:

```bash
#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/../core/lib/state.sh"

# Use temp file for testing
STATE_FILE=$(mktemp)
export STATE_FILE

cleanup() {
    rm -f "$STATE_FILE"
}
trap cleanup EXIT

# Test initial state
test_init_state() {
    init_state
    
    if [[ ! -f "$STATE_FILE" ]]; then
        echo "FAIL: init_state did not create file"
        return 1
    fi
    echo "PASS: init_state"
}

# Test get/set state
test_get_set_state() {
    set_state "lastEngagement" "acme-corp"
    local result
    result=$(get_state "lastEngagement")
    
    if [[ "$result" != "acme-corp" ]]; then
        echo "FAIL: get_state returned '$result' instead of 'acme-corp'"
        return 1
    fi
    echo "PASS: get_set_state"
}

# Test get_all_state
test_get_all_state() {
    set_state "foo" "bar"
    local all
    all=$(get_all_state)
    
    if ! echo "$all" | grep -q '"foo"'; then
        echo "FAIL: get_all_state missing foo"
        return 1
    fi
    echo "PASS: get_all_state"
}

# Run tests
test_init_state
test_get_set_state
test_get_all_state

echo "All state tests passed"
```

- [ ] **Step 2: Run test to verify it fails**

```bash
chmod +x tui/tests/test_state.sh
./tui/tests/test_state.sh
```

Expected: FAIL with "source: no such file"

- [ ] **Step 3: Implement state.sh**

Create file `tui/core/lib/state.sh`:

```bash
#!/usr/bin/env bash
# state.sh - Persistent state management

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/../protocol.sh"

# Default state file location (can be overridden for testing)
STATE_FILE="${STATE_FILE:-$(cd "$SCRIPT_DIR/../.." && pwd)/.tui-state.json}"

# Initialize state file if it doesn't exist
init_state() {
    if [[ ! -f "$STATE_FILE" ]]; then
        echo '{"lastEngagement":"","preferences":{}}' > "$STATE_FILE"
    fi
}

# Get a state value
# Usage: get_state <key>
get_state() {
    local key="$1"
    init_state
    $JQ -r ".$key // empty" "$STATE_FILE"
}

# Set a state value
# Usage: set_state <key> <value>
set_state() {
    local key="$1"
    local value="$2"
    init_state
    
    local tmp
    tmp=$(mktemp)
    $JQ --arg key "$key" --arg val "$value" '.[$key] = $val' "$STATE_FILE" > "$tmp"
    mv "$tmp" "$STATE_FILE"
}

# Get all state as JSON
get_all_state() {
    init_state
    cat "$STATE_FILE"
}

# Clear state
clear_state() {
    rm -f "$STATE_FILE"
    init_state
}
```

- [ ] **Step 4: Run test to verify it passes**

```bash
./tui/tests/test_state.sh
```

Expected: "All state tests passed"

- [ ] **Step 5: Commit**

```bash
git add tui/core/lib/state.sh tui/tests/test_state.sh
git commit -m "feat(tui): add persistent state management"
```

---

### Task 5: Implement Skills Module

**Files:**
- Create: `tui/core/lib/skills.sh`
- Create: `tui/tests/test_skills.sh`
- Create: `tui/tests/fixtures/.claude/skills/test-skill/SKILL.md`

- [ ] **Step 1: Create test fixture**

```bash
mkdir -p "tui/tests/fixtures/.claude/skills/test-skill"
```

Create file `tui/tests/fixtures/.claude/skills/test-skill/SKILL.md`:

```markdown
---
name: test-skill
description: A test skill for unit testing
---

# Test Skill

This is a test skill.
```

- [ ] **Step 2: Write skills test**

Create file `tui/tests/test_skills.sh`:

```bash
#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
export PROJECT_ROOT="$SCRIPT_DIR/fixtures"
source "$SCRIPT_DIR/../core/lib/skills.sh"

# Test list_skills
test_list_skills() {
    local result
    result=$(list_skills)
    
    if ! echo "$result" | grep -q '"name":"test-skill"'; then
        echo "FAIL: list_skills did not find test-skill"
        echo "Got: $result"
        return 1
    fi
    echo "PASS: list_skills"
}

# Test get_skill
test_get_skill() {
    local result
    result=$(get_skill "test-skill")
    
    if ! echo "$result" | grep -q '"description":"A test skill for unit testing"'; then
        echo "FAIL: get_skill description mismatch"
        echo "Got: $result"
        return 1
    fi
    echo "PASS: get_skill"
}

# Test get_skill not found
test_get_skill_not_found() {
    local result
    result=$(get_skill "nonexistent" 2>&1 || true)
    
    if ! echo "$result" | grep -q "null\|not found"; then
        echo "FAIL: get_skill should return null for missing skill"
        return 1
    fi
    echo "PASS: get_skill_not_found"
}

# Run tests
test_list_skills
test_get_skill
test_get_skill_not_found

echo "All skills tests passed"
```

- [ ] **Step 3: Run test to verify it fails**

```bash
chmod +x tui/tests/test_skills.sh
./tui/tests/test_skills.sh
```

Expected: FAIL with "source: no such file"

- [ ] **Step 4: Implement skills.sh**

Create file `tui/core/lib/skills.sh`:

```bash
#!/usr/bin/env bash
# skills.sh - Skill discovery and reading

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/../protocol.sh"

# Project root (can be overridden for testing)
PROJECT_ROOT="${PROJECT_ROOT:-$(cd "$SCRIPT_DIR/../../.." && pwd)}"
SKILLS_DIR="$PROJECT_ROOT/.claude/skills"

# Extract description from SKILL.md frontmatter
_extract_description() {
    local file="$1"
    # Parse YAML frontmatter for description field
    sed -n '/^---$/,/^---$/p' "$file" | grep -E '^description:' | sed 's/^description:[[:space:]]*//' | tr -d '\n'
}

# List all available skills
# Returns JSON array of skill objects
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

# Get a single skill by name
# Returns skill object or null
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
```

- [ ] **Step 5: Run test to verify it passes**

```bash
./tui/tests/test_skills.sh
```

Expected: "All skills tests passed"

- [ ] **Step 6: Commit**

```bash
git add tui/core/lib/skills.sh tui/tests/test_skills.sh tui/tests/fixtures/
git commit -m "feat(tui): add skill discovery and reading"
```

---

### Task 6: Implement Engagements Module

**Files:**
- Create: `tui/core/lib/engagements.sh`
- Create: `tui/tests/test_engagements.sh`
- Create: `tui/tests/fixtures/engagements/test-customer/CONTEXT.md`

- [ ] **Step 1: Create test fixture**

```bash
mkdir -p "tui/tests/fixtures/engagements/test-customer"
```

Create file `tui/tests/fixtures/engagements/test-customer/CONTEXT.md`:

```markdown
# Engagement Context: Test Customer

## Engagement Metadata
- **Customer:** Test Customer
- **Type:** Platform Assessment
```

- [ ] **Step 2: Write engagements test**

Create file `tui/tests/test_engagements.sh`:

```bash
#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
export PROJECT_ROOT="$SCRIPT_DIR/fixtures"
source "$SCRIPT_DIR/../core/lib/engagements.sh"

# Test list_engagements
test_list_engagements() {
    local result
    result=$(list_engagements)
    
    if ! echo "$result" | grep -q '"slug":"test-customer"'; then
        echo "FAIL: list_engagements did not find test-customer"
        echo "Got: $result"
        return 1
    fi
    if ! echo "$result" | grep -q '"hasContext":true'; then
        echo "FAIL: test-customer should have hasContext=true"
        return 1
    fi
    echo "PASS: list_engagements"
}

# Test get_engagement
test_get_engagement() {
    local result
    result=$(get_engagement "test-customer")
    
    if ! echo "$result" | grep -q '"slug":"test-customer"'; then
        echo "FAIL: get_engagement did not return correct slug"
        return 1
    fi
    echo "PASS: get_engagement"
}

# Run tests
test_list_engagements
test_get_engagement

echo "All engagements tests passed"
```

- [ ] **Step 3: Run test to verify it fails**

```bash
chmod +x tui/tests/test_engagements.sh
./tui/tests/test_engagements.sh
```

Expected: FAIL with "source: no such file"

- [ ] **Step 4: Implement engagements.sh**

Create file `tui/core/lib/engagements.sh`:

```bash
#!/usr/bin/env bash
# engagements.sh - Engagement management

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/../protocol.sh"

PROJECT_ROOT="${PROJECT_ROOT:-$(cd "$SCRIPT_DIR/../../.." && pwd)}"
ENGAGEMENTS_DIR="$PROJECT_ROOT/engagements"

# List all engagements
# Returns JSON array of engagement objects
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
        
        # Skip template directory
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

# Get a single engagement by slug
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

# Create a new engagement
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
```

- [ ] **Step 5: Run test to verify it passes**

```bash
./tui/tests/test_engagements.sh
```

Expected: "All engagements tests passed"

- [ ] **Step 6: Commit**

```bash
git add tui/core/lib/engagements.sh tui/tests/test_engagements.sh tui/tests/fixtures/engagements/
git commit -m "feat(tui): add engagement management"
```

---

### Task 7: Implement Phase Detection

**Files:**
- Create: `tui/core/lib/phase.sh`
- Create: `tui/tests/test_phase.sh`

- [ ] **Step 1: Write phase test**

Create file `tui/tests/test_phase.sh`:

```bash
#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
export PROJECT_ROOT="$SCRIPT_DIR/fixtures"
source "$SCRIPT_DIR/../core/lib/phase.sh"

# Test detect_phase for engagement with context
test_detect_phase() {
    local result
    result=$(detect_phase "test-customer")
    
    if ! echo "$result" | grep -q '"phase":"live"'; then
        echo "FAIL: detect_phase should return 'live' for engagement with CONTEXT.md"
        echo "Got: $result"
        return 1
    fi
    echo "PASS: detect_phase"
}

# Test detect_phase for nonexistent engagement
test_detect_phase_missing() {
    local result
    result=$(detect_phase "nonexistent")
    
    if ! echo "$result" | grep -q '"phase":"pre-engagement"'; then
        echo "FAIL: detect_phase should return 'pre-engagement' for missing engagement"
        return 1
    fi
    echo "PASS: detect_phase_missing"
}

# Run tests
test_detect_phase
test_detect_phase_missing

echo "All phase tests passed"
```

- [ ] **Step 2: Run test to verify it fails**

```bash
chmod +x tui/tests/test_phase.sh
./tui/tests/test_phase.sh
```

Expected: FAIL

- [ ] **Step 3: Implement phase.sh**

Create file `tui/core/lib/phase.sh`:

```bash
#!/usr/bin/env bash
# phase.sh - Phase detection

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/../protocol.sh"

PROJECT_ROOT="${PROJECT_ROOT:-$(cd "$SCRIPT_DIR/../../.." && pwd)}"
ENGAGEMENTS_DIR="$PROJECT_ROOT/engagements"

# Count files in a directory (non-recursive, files only)
_count_files() {
    local dir="$1"
    if [[ -d "$dir" ]]; then
        find "$dir" -maxdepth 1 -type f | wc -l | tr -d ' '
    else
        echo "0"
    fi
}

# Detect phase for an engagement
# Returns PhaseInfo JSON
detect_phase() {
    local slug="$1"
    local eng_dir="$ENGAGEMENTS_DIR/$slug"
    
    local phase="pre-engagement"
    local discovery_count=0
    local assessments_count=0
    local deliverables_count=0
    
    if [[ -d "$eng_dir" ]]; then
        discovery_count=$(_count_files "$eng_dir/discovery")
        assessments_count=$(_count_files "$eng_dir/assessments")
        deliverables_count=$(_count_files "$eng_dir/deliverables")
        
        # Determine phase based on artifacts
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
        '{
            phase: $phase,
            artifactCounts: {
                discovery: $discovery,
                assessments: $assessments,
                deliverables: $deliverables
            }
        }'
}
```

- [ ] **Step 4: Run test to verify it passes**

```bash
./tui/tests/test_phase.sh
```

Expected: "All phase tests passed"

- [ ] **Step 5: Commit**

```bash
git add tui/core/lib/phase.sh tui/tests/test_phase.sh
git commit -m "feat(tui): add phase detection"
```

---

### Task 8: Implement Artifacts Module

**Files:**
- Create: `tui/core/lib/artifacts.sh`
- Create: `tui/tests/test_artifacts.sh`

- [ ] **Step 1: Create additional test fixtures**

```bash
echo "# Discovery Notes" > "tui/tests/fixtures/engagements/test-customer/discovery/notes.md"
```

- [ ] **Step 2: Write artifacts test**

Create file `tui/tests/test_artifacts.sh`:

```bash
#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
export PROJECT_ROOT="$SCRIPT_DIR/fixtures"
source "$SCRIPT_DIR/../core/lib/artifacts.sh"

# Test list_artifacts
test_list_artifacts() {
    local result
    result=$(list_artifacts "test-customer")
    
    if ! echo "$result" | grep -q '"name":"discovery"'; then
        echo "FAIL: list_artifacts should include discovery directory"
        echo "Got: $result"
        return 1
    fi
    echo "PASS: list_artifacts"
}

# Test read_artifact
test_read_artifact() {
    local result
    result=$(read_artifact "engagements/test-customer/CONTEXT.md")
    
    if ! echo "$result" | grep -q "Test Customer"; then
        echo "FAIL: read_artifact did not return correct content"
        return 1
    fi
    echo "PASS: read_artifact"
}

# Run tests
test_list_artifacts
test_read_artifact

echo "All artifacts tests passed"
```

- [ ] **Step 3: Run test to verify it fails**

```bash
chmod +x tui/tests/test_artifacts.sh
./tui/tests/test_artifacts.sh
```

Expected: FAIL

- [ ] **Step 4: Implement artifacts.sh**

Create file `tui/core/lib/artifacts.sh`:

```bash
#!/usr/bin/env bash
# artifacts.sh - Artifact browsing

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/../protocol.sh"

PROJECT_ROOT="${PROJECT_ROOT:-$(cd "$SCRIPT_DIR/../../.." && pwd)}"
ENGAGEMENTS_DIR="$PROJECT_ROOT/engagements"

# Build tree structure for a directory
_build_tree() {
    local dir="$1"
    local rel_path="$2"
    local result="[]"
    
    for entry in "$dir"/*; do
        [[ -e "$entry" ]] || continue
        
        local name
        name=$(basename "$entry")
        local entry_path="$rel_path/$name"
        
        if [[ -d "$entry" ]]; then
            local children
            children=$(_build_tree "$entry" "$entry_path")
            result=$($JQ -c \
                --arg name "$name" \
                --arg path "$entry_path" \
                --argjson children "$children" \
                '. + [{name: $name, path: $path, type: "directory", children: $children}]' <<< "$result")
        else
            result=$($JQ -c \
                --arg name "$name" \
                --arg path "$entry_path" \
                '. + [{name: $name, path: $path, type: "file"}]' <<< "$result")
        fi
    done
    
    echo "$result"
}

# List artifacts for an engagement as a tree
list_artifacts() {
    local slug="$1"
    local eng_dir="$ENGAGEMENTS_DIR/$slug"
    
    if [[ ! -d "$eng_dir" ]]; then
        echo "[]"
        return
    fi
    
    _build_tree "$eng_dir" "engagements/$slug"
}

# Read artifact content
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
```

- [ ] **Step 5: Run test to verify it passes**

```bash
./tui/tests/test_artifacts.sh
```

Expected: "All artifacts tests passed"

- [ ] **Step 6: Commit**

```bash
git add tui/core/lib/artifacts.sh tui/tests/test_artifacts.sh tui/tests/fixtures/engagements/test-customer/discovery/
git commit -m "feat(tui): add artifact browsing"
```

---

### Task 9: Implement Checklists Module

**Files:**
- Create: `tui/core/lib/checklists.sh`
- Create: `tui/tests/test_checklists.sh`
- Create: `tui/tests/fixtures/knowledge/checklists/test-checklist.md`

- [ ] **Step 1: Create test fixture**

```bash
mkdir -p "tui/tests/fixtures/knowledge/checklists"
```

Create file `tui/tests/fixtures/knowledge/checklists/test-checklist.md`:

```markdown
# Test Checklist

## Section One

- [ ] First item
- [x] Second item (checked)
- [ ] Third item

## Section Two

- [ ] Fourth item
```

- [ ] **Step 2: Write checklists test**

Create file `tui/tests/test_checklists.sh`:

```bash
#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
export PROJECT_ROOT="$SCRIPT_DIR/fixtures"
source "$SCRIPT_DIR/../core/lib/checklists.sh"

# Test list_checklists
test_list_checklists() {
    local result
    result=$(list_checklists)
    
    if ! echo "$result" | grep -q '"name":"test-checklist"'; then
        echo "FAIL: list_checklists did not find test-checklist"
        echo "Got: $result"
        return 1
    fi
    echo "PASS: list_checklists"
}

# Test get_checklist
test_get_checklist() {
    local result
    result=$(get_checklist "test-checklist")
    
    if ! echo "$result" | grep -q '"title":"Section One"'; then
        echo "FAIL: get_checklist did not parse sections"
        echo "Got: $result"
        return 1
    fi
    if ! echo "$result" | grep -q '"checked":true'; then
        echo "FAIL: get_checklist did not detect checked item"
        return 1
    fi
    echo "PASS: get_checklist"
}

# Run tests
test_list_checklists
test_get_checklist

echo "All checklists tests passed"
```

- [ ] **Step 3: Run test to verify it fails**

```bash
chmod +x tui/tests/test_checklists.sh
./tui/tests/test_checklists.sh
```

Expected: FAIL

- [ ] **Step 4: Implement checklists.sh**

Create file `tui/core/lib/checklists.sh`:

```bash
#!/usr/bin/env bash
# checklists.sh - Checklist parsing and toggling

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/../protocol.sh"

PROJECT_ROOT="${PROJECT_ROOT:-$(cd "$SCRIPT_DIR/../../.." && pwd)}"
CHECKLISTS_DIR="$PROJECT_ROOT/knowledge/checklists"

# List all available checklists
list_checklists() {
    local checklists="[]"
    
    if [[ ! -d "$CHECKLISTS_DIR" ]]; then
        echo "$checklists"
        return
    fi
    
    for file in "$CHECKLISTS_DIR"/*.md; do
        [[ -f "$file" ]] || continue
        
        local filename
        filename=$(basename "$file")
        local name="${filename%.md}"
        
        checklists=$($JQ -c \
            --arg name "$name" \
            --arg fileName "$filename" \
            '. + [{name: $name, fileName: $fileName}]' <<< "$checklists")
    done
    
    echo "$checklists"
}

# Parse a checklist file into structured JSON
get_checklist() {
    local name="$1"
    local file="$CHECKLISTS_DIR/$name.md"
    
    if [[ ! -f "$file" ]]; then
        echo "null"
        return
    fi
    
    local sections="[]"
    local current_section=""
    local current_items="[]"
    local line_num=0
    local total_items=0
    local checked_items=0
    
    while IFS= read -r line || [[ -n "$line" ]]; do
        ((line_num++))
        
        # Check for section header
        if [[ "$line" =~ ^##[[:space:]]+(.+)$ ]]; then
            # Save previous section if exists
            if [[ -n "$current_section" ]]; then
                sections=$($JQ -c \
                    --arg title "$current_section" \
                    --argjson items "$current_items" \
                    '. + [{title: $title, items: $items}]' <<< "$sections")
            fi
            current_section="${BASH_REMATCH[1]}"
            current_items="[]"
            continue
        fi
        
        # Check for checkbox item
        if [[ "$line" =~ ^-[[:space:]]+\[([xX[:space:]])\][[:space:]]+(.+)$ ]]; then
            local checked="false"
            [[ "${BASH_REMATCH[1]}" =~ [xX] ]] && checked="true"
            local text="${BASH_REMATCH[2]}"
            
            ((total_items++))
            [[ "$checked" == "true" ]] && ((checked_items++))
            
            current_items=$($JQ -c \
                --arg text "$text" \
                --argjson checked "$checked" \
                --argjson line "$line_num" \
                '. + [{text: $text, checked: $checked, line: $line}]' <<< "$current_items")
        fi
    done < "$file"
    
    # Save last section
    if [[ -n "$current_section" ]]; then
        sections=$($JQ -c \
            --arg title "$current_section" \
            --argjson items "$current_items" \
            '. + [{title: $title, items: $items}]' <<< "$sections")
    fi
    
    # Calculate completion percentage
    local completion_pct=0
    if (( total_items > 0 )); then
        completion_pct=$(( (checked_items * 100) / total_items ))
    fi
    
    $JQ -cn \
        --arg name "$name" \
        --arg fileName "$name.md" \
        --argjson sections "$sections" \
        --argjson pct "$completion_pct" \
        '{name: $name, fileName: $fileName, sections: $sections, completionPercent: $pct}'
}

# Toggle a checklist item
toggle_checklist_item() {
    local name="$1"
    local line_num="$2"
    local file="$CHECKLISTS_DIR/$name.md"
    
    if [[ ! -f "$file" ]]; then
        echo '{"success": false, "error": "Checklist not found"}'
        return 1
    fi
    
    # Read line and toggle checkbox
    local line
    line=$(sed -n "${line_num}p" "$file")
    
    if [[ "$line" =~ \[\ \] ]]; then
        # Unchecked -> checked
        sed -i.bak "${line_num}s/\[ \]/[x]/" "$file"
        rm -f "$file.bak"
    elif [[ "$line" =~ \[[xX]\] ]]; then
        # Checked -> unchecked
        sed -i.bak "${line_num}s/\[[xX]\]/[ ]/" "$file"
        rm -f "$file.bak"
    else
        echo '{"success": false, "error": "Line is not a checkbox"}'
        return 1
    fi
    
    echo '{"success": true}'
}
```

- [ ] **Step 5: Run test to verify it passes**

```bash
./tui/tests/test_checklists.sh
```

Expected: "All checklists tests passed"

- [ ] **Step 6: Commit**

```bash
git add tui/core/lib/checklists.sh tui/tests/test_checklists.sh tui/tests/fixtures/knowledge/
git commit -m "feat(tui): add checklist parsing and toggling"
```

---

### Task 10: Implement Agents Module

**Files:**
- Create: `tui/core/lib/agents.sh`
- Create: `tui/tests/test_agents.sh`
- Create: `tui/tests/fixtures/.claude/agents/test-agent.md`

- [ ] **Step 1: Create test fixture**

```bash
mkdir -p "tui/tests/fixtures/.claude/agents"
```

Create file `tui/tests/fixtures/.claude/agents/test-agent.md`:

```markdown
---
name: test-agent
model: sonnet
role: Testing Agent
description: An agent for unit testing
---

# Test Agent

This is a test agent.
```

- [ ] **Step 2: Write agents test**

Create file `tui/tests/test_agents.sh`:

```bash
#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
export PROJECT_ROOT="$SCRIPT_DIR/fixtures"
source "$SCRIPT_DIR/../core/lib/agents.sh"

# Test list_agents
test_list_agents() {
    local result
    result=$(list_agents)
    
    if ! echo "$result" | grep -q '"name":"test-agent"'; then
        echo "FAIL: list_agents did not find test-agent"
        echo "Got: $result"
        return 1
    fi
    if ! echo "$result" | grep -q '"model":"sonnet"'; then
        echo "FAIL: list_agents did not parse model"
        return 1
    fi
    echo "PASS: list_agents"
}

# Run tests
test_list_agents

echo "All agents tests passed"
```

- [ ] **Step 3: Run test to verify it fails**

```bash
chmod +x tui/tests/test_agents.sh
./tui/tests/test_agents.sh
```

Expected: FAIL

- [ ] **Step 4: Implement agents.sh**

Create file `tui/core/lib/agents.sh`:

```bash
#!/usr/bin/env bash
# agents.sh - Agent discovery and invocation

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/../protocol.sh"

PROJECT_ROOT="${PROJECT_ROOT:-$(cd "$SCRIPT_DIR/../../.." && pwd)}"
AGENTS_DIR="$PROJECT_ROOT/.claude/agents"

# Extract a field from markdown frontmatter
_extract_frontmatter_field() {
    local file="$1"
    local field="$2"
    sed -n '/^---$/,/^---$/p' "$file" | grep -E "^${field}:" | sed "s/^${field}:[[:space:]]*//" | tr -d '\n'
}

# List all available agents
list_agents() {
    local agents="[]"
    
    if [[ ! -d "$AGENTS_DIR" ]]; then
        echo "$agents"
        return
    fi
    
    for file in "$AGENTS_DIR"/*.md; do
        [[ -f "$file" ]] || continue
        
        local filename
        filename=$(basename "$file")
        
        # Skip README
        [[ "$filename" == "README.md" ]] && continue
        
        local name
        name=$(_extract_frontmatter_field "$file" "name")
        [[ -z "$name" ]] && name="${filename%.md}"
        
        local model
        model=$(_extract_frontmatter_field "$file" "model")
        local role
        role=$(_extract_frontmatter_field "$file" "role")
        local description
        description=$(_extract_frontmatter_field "$file" "description")
        
        agents=$($JQ -c \
            --arg name "$name" \
            --arg model "$model" \
            --arg role "$role" \
            --arg desc "$description" \
            '. + [{name: $name, model: $model, role: $role, description: $desc}]' <<< "$agents")
    done
    
    echo "$agents"
}

# Get a single agent by name
get_agent() {
    local name="$1"
    local file="$AGENTS_DIR/$name.md"
    
    if [[ ! -f "$file" ]]; then
        echo "null"
        return
    fi
    
    local model
    model=$(_extract_frontmatter_field "$file" "model")
    local role
    role=$(_extract_frontmatter_field "$file" "role")
    local description
    description=$(_extract_frontmatter_field "$file" "description")
    local content
    content=$(cat "$file")
    
    $JQ -cn \
        --arg name "$name" \
        --arg model "$model" \
        --arg role "$role" \
        --arg desc "$description" \
        --arg content "$content" \
        '{name: $name, model: $model, role: $role, description: $desc, content: $content}'
}
```

- [ ] **Step 5: Run test to verify it passes**

```bash
./tui/tests/test_agents.sh
```

Expected: "All agents tests passed"

- [ ] **Step 6: Commit**

```bash
git add tui/core/lib/agents.sh tui/tests/test_agents.sh tui/tests/fixtures/.claude/agents/
git commit -m "feat(tui): add agent discovery"
```

---

### Task 11: Implement Claude Execution

**Files:**
- Create: `tui/core/lib/claude.sh`

- [ ] **Step 1: Implement claude.sh**

Create file `tui/core/lib/claude.sh`:

```bash
#!/usr/bin/env bash
# claude.sh - Claude CLI execution and streaming

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/../protocol.sh"

PROJECT_ROOT="${PROJECT_ROOT:-$(cd "$SCRIPT_DIR/../../.." && pwd)}"

# Current execution PID (for cancellation)
CLAUDE_PID=""

# Parse Claude stream-json output and forward as events
_parse_claude_event() {
    local line="$1"
    
    # Try to parse as JSON
    local type
    type=$($JQ -r '.type // empty' <<< "$line" 2>/dev/null) || return
    
    case "$type" in
        assistant)
            # Extract text from content blocks
            local text
            text=$($JQ -r '.message.content[]? | select(.type == "text") | .text // empty' <<< "$line" 2>/dev/null)
            if [[ -n "$text" ]]; then
                send_event "assistant" "$($JQ -cn --arg text "$text" '{text: $text}')"
            fi
            
            # Check for tool use
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
            # Check for tool result
            local tool_name output
            tool_name=$($JQ -r '.tool_name // empty' <<< "$line" 2>/dev/null)
            if [[ -n "$tool_name" ]]; then
                output=$($JQ -r '.message.content[]? | select(.type == "text") | .text // empty' <<< "$line" 2>/dev/null)
                send_event "tool_result" "$($JQ -cn --arg tool "$tool_name" --arg output "$output" '{tool: $tool, output: $output}')"
            fi
            
            # Check for final result with cost
            local cost
            cost=$($JQ -r '.total_cost_usd // empty' <<< "$line" 2>/dev/null)
            if [[ -n "$cost" ]]; then
                send_event "complete" "$($JQ -cn --arg status "success" --argjson cost "$cost" '{status: $status, totalCost: $cost}')"
            fi
            ;;
        
        system)
            # Session started, ignore for now
            ;;
        
        *)
            # Unknown event type, log for debugging
            ;;
    esac
}

# Execute a skill
execute_skill() {
    local skill="$1"
    local engagement="$2"
    
    local prompt="Run /${skill}."
    [[ -n "$engagement" ]] && prompt="$prompt Use engagement at engagements/${engagement}/."
    
    _execute_claude "$prompt"
}

# Execute an agent with a prompt
execute_agent() {
    local agent="$1"
    local prompt="$2"
    local engagement="$3"
    
    # Build agent invocation command
    local full_prompt="Using the $agent agent: $prompt"
    [[ -n "$engagement" ]] && full_prompt="$full_prompt (Engagement: $engagement)"
    
    _execute_claude "$full_prompt"
}

# Core Claude execution
_execute_claude() {
    local prompt="$1"
    
    # Start Claude in background
    claude --print \
           --output-format stream-json \
           --verbose \
           --permission-mode acceptEdits \
           "$prompt" 2>&1 &
    
    CLAUDE_PID=$!
    
    # Read output line by line
    while IFS= read -r line; do
        [[ -z "$line" ]] && continue
        _parse_claude_event "$line"
    done < <(wait $CLAUDE_PID 2>&1; echo '{"type":"_done"}')
    
    CLAUDE_PID=""
}

# Cancel running execution
cancel_execution() {
    if [[ -n "$CLAUDE_PID" ]] && kill -0 "$CLAUDE_PID" 2>/dev/null; then
        kill -TERM "$CLAUDE_PID" 2>/dev/null
        CLAUDE_PID=""
        send_event "complete" '{"status": "cancelled"}'
    fi
}

# Send user input to running Claude (for questions)
send_user_input() {
    local text="$1"
    # For now, this requires session resumption which is complex
    # Placeholder for future implementation
    send_error "User input during execution not yet implemented"
}
```

- [ ] **Step 2: Verify syntax**

```bash
bash -n tui/core/lib/claude.sh
```

Expected: No output (no syntax errors)

- [ ] **Step 3: Commit**

```bash
git add tui/core/lib/claude.sh
git commit -m "feat(tui): add Claude CLI execution and streaming"
```

---

### Task 12: Implement Main Event Loop

**Files:**
- Create: `tui/core/main.sh`

- [ ] **Step 1: Implement main.sh**

Create file `tui/core/main.sh`:

```bash
#!/usr/bin/env bash
# main.sh - Main event loop for TUI core

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Source all modules
source "$SCRIPT_DIR/protocol.sh"
source "$SCRIPT_DIR/lib/state.sh"
source "$SCRIPT_DIR/lib/skills.sh"
source "$SCRIPT_DIR/lib/engagements.sh"
source "$SCRIPT_DIR/lib/phase.sh"
source "$SCRIPT_DIR/lib/artifacts.sh"
source "$SCRIPT_DIR/lib/checklists.sh"
source "$SCRIPT_DIR/lib/agents.sh"
source "$SCRIPT_DIR/lib/claude.sh"

# Handle a command from the viewer
handle_command() {
    local message="$1"
    
    local cmd id args
    cmd=$(echo "$message" | parse_command_field "cmd")
    id=$(echo "$message" | parse_command_field "id")
    args=$(echo "$message" | parse_command_args)
    
    case "$cmd" in
        init)
            # Return initial state
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

# Main event loop
main() {
    # Initialize state
    init_state
    
    # Read commands from stdin, respond to stdout
    while IFS= read -r message; do
        [[ -z "$message" ]] && continue
        handle_command "$message"
    done
}

# Run if executed directly (not sourced)
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main
fi
```

- [ ] **Step 2: Make executable and verify syntax**

```bash
chmod +x tui/core/main.sh
bash -n tui/core/main.sh
```

Expected: No output (no syntax errors)

- [ ] **Step 3: Commit**

```bash
git add tui/core/main.sh
git commit -m "feat(tui): add main event loop"
```

---

## Phase 3: Go Viewer Layer

### Task 13: Implement Styles

**Files:**
- Create: `tui/viewer/ui/styles.go`

- [ ] **Step 1: Add Lipgloss dependency**

```bash
cd tui/viewer
go get github.com/charmbracelet/lipgloss
```

- [ ] **Step 2: Implement styles.go**

Create file `tui/viewer/ui/styles.go`:

```go
package ui

import "github.com/charmbracelet/lipgloss"

// Red Hat brand colors
var (
	RedHatRed      = lipgloss.Color("#EE0000")
	RedHatRedDark  = lipgloss.Color("#A30000")
	RedHatBlack    = lipgloss.Color("#151515")
	Surface        = lipgloss.Color("#1A1A1A")
	SurfaceLight   = lipgloss.Color("#2E2E2E")
	Border         = lipgloss.Color("#3A3A3A")
	TextPrimary    = lipgloss.Color("#E8E8E8")
	TextMuted      = lipgloss.Color("#888888")
	TextDim        = lipgloss.Color("#555555")
	Green          = lipgloss.Color("#3E8635")
	Yellow         = lipgloss.Color("#F0AB00")
	Blue           = lipgloss.Color("#0066CC")
)

// Layout styles
var (
	SidebarWidth = 32

	SidebarStyle = lipgloss.NewStyle().
			Width(SidebarWidth).
			Padding(1, 2).
			BorderRight(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(RedHatRed)

	MainStyle = lipgloss.NewStyle().
			Padding(1, 2)

	TitleStyle = lipgloss.NewStyle().
			Foreground(RedHatRed).
			Bold(true).
			MarginBottom(1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(TextMuted).
			Italic(true)
)

// Menu styles
var (
	MenuItemStyle = lipgloss.NewStyle().
			Foreground(TextPrimary).
			PaddingLeft(2)

	MenuItemSelectedStyle = lipgloss.NewStyle().
				Foreground(RedHatRed).
				Bold(true).
				PaddingLeft(2).
				SetString("> ")

	MenuHeaderStyle = lipgloss.NewStyle().
			Foreground(TextMuted).
			Bold(true).
			MarginTop(1).
			MarginBottom(1)
)

// Activity log styles
var (
	ActivityTimestampStyle = lipgloss.NewStyle().
				Foreground(TextDim).
				Width(10)

	ActivityAssistantStyle = lipgloss.NewStyle().
				Foreground(TextPrimary)

	ActivityToolStyle = lipgloss.NewStyle().
				Foreground(Blue)

	ActivityToolResultStyle = lipgloss.NewStyle().
				Foreground(TextMuted)

	ActivityErrorStyle = lipgloss.NewStyle().
				Foreground(RedHatRed).
				Bold(true)
)

// Input styles
var (
	InputStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(RedHatRed).
			Padding(0, 1)

	InputLabelStyle = lipgloss.NewStyle().
			Foreground(TextMuted).
			MarginBottom(1)
)

// Command palette styles
var (
	PaletteStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(RedHatRed).
			Padding(1, 2).
			Width(60)

	PaletteInputStyle = lipgloss.NewStyle().
				MarginBottom(1)

	PaletteItemStyle = lipgloss.NewStyle().
				Foreground(TextPrimary)

	PaletteItemSelectedStyle = lipgloss.NewStyle().
					Foreground(RedHatRed).
					Bold(true)

	PaletteItemDescStyle = lipgloss.NewStyle().
				Foreground(TextMuted).
				MarginLeft(2)
)

// Status indicators
var (
	StatusRunningStyle = lipgloss.NewStyle().
				Foreground(Yellow).
				Bold(true)

	StatusCompleteStyle = lipgloss.NewStyle().
				Foreground(Green).
				Bold(true)

	StatusErrorStyle = lipgloss.NewStyle().
				Foreground(RedHatRed).
				Bold(true)

	PhasePreEngagementStyle = lipgloss.NewStyle().
				Foreground(TextMuted)

	PhaseLiveStyle = lipgloss.NewStyle().
			Foreground(Green).
			Bold(true)

	PhaseLeaveBehindStyle = lipgloss.NewStyle().
				Foreground(Blue)
)

// Help styles
var (
	HelpKeyStyle = lipgloss.NewStyle().
			Foreground(RedHatRed)

	HelpDescStyle = lipgloss.NewStyle().
			Foreground(TextMuted)
)
```

- [ ] **Step 3: Verify Go compiles**

```bash
cd tui/viewer
go build ./ui/...
```

Expected: No errors

- [ ] **Step 4: Commit**

```bash
git add tui/viewer/go.mod tui/viewer/go.sum tui/viewer/ui/styles.go
git commit -m "feat(tui): add Red Hat brand styles"
```

---

### Task 14: Implement Sidebar Component

**Files:**
- Create: `tui/viewer/ui/sidebar.go`

- [ ] **Step 1: Implement sidebar.go**

Create file `tui/viewer/ui/sidebar.go`:

```go
package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/toddward/red-hat-engagement-kit/tui/viewer/protocol"
)

// Sidebar displays engagement context
type Sidebar struct {
	width       int
	height      int
	engagement  string
	phase       protocol.Phase
	agent       string
	counts      protocol.ArtifactCounts
}

// NewSidebar creates a new sidebar
func NewSidebar() Sidebar {
	return Sidebar{
		width: SidebarWidth,
		phase: protocol.PhasePreEngagement,
	}
}

// SetSize sets the sidebar dimensions
func (s *Sidebar) SetSize(width, height int) {
	s.width = width
	s.height = height
}

// SetEngagement sets the current engagement
func (s *Sidebar) SetEngagement(slug string) {
	s.engagement = slug
}

// SetPhase sets the current phase info
func (s *Sidebar) SetPhase(info protocol.PhaseInfo) {
	s.phase = info.Phase
	s.counts = info.ArtifactCounts
}

// SetAgent sets the current active agent
func (s *Sidebar) SetAgent(name string) {
	s.agent = name
}

// View renders the sidebar
func (s Sidebar) View() string {
	var b strings.Builder

	// Logo/Title
	title := TitleStyle.Render("Red Hat")
	subtitle := SubtitleStyle.Render("Engagement Kit")
	b.WriteString(title + "\n")
	b.WriteString(subtitle + "\n\n")

	// Engagement section
	b.WriteString(MenuHeaderStyle.Render("ENGAGEMENT"))
	b.WriteString("\n")
	if s.engagement != "" {
		b.WriteString(lipgloss.NewStyle().Foreground(TextPrimary).Render(s.engagement))
		b.WriteString("\n")
		
		// Phase indicator
		phaseStyle := PhasePreEngagementStyle
		phaseText := "Pre-Engagement"
		switch s.phase {
		case protocol.PhaseLive:
			phaseStyle = PhaseLiveStyle
			phaseText = "● Live"
		case protocol.PhaseLeaveBehind:
			phaseStyle = PhaseLeaveBehindStyle
			phaseText = "Leave-Behind"
		}
		b.WriteString(phaseStyle.Render(phaseText))
		b.WriteString("\n")
	} else {
		b.WriteString(lipgloss.NewStyle().Foreground(TextDim).Italic(true).Render("None selected"))
		b.WriteString("\n")
	}
	b.WriteString("\n")

	// Agent section (if active)
	if s.agent != "" {
		b.WriteString(MenuHeaderStyle.Render("AGENT"))
		b.WriteString("\n")
		b.WriteString(StatusRunningStyle.Render("● " + s.agent))
		b.WriteString("\n\n")
	}

	// Artifact counts
	if s.engagement != "" {
		b.WriteString(MenuHeaderStyle.Render("ARTIFACTS"))
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("Discovery:    %d\n", s.counts.Discovery))
		b.WriteString(fmt.Sprintf("Assessments:  %d\n", s.counts.Assessments))
		b.WriteString(fmt.Sprintf("Deliverables: %d\n", s.counts.Deliverables))
	}

	content := b.String()
	return SidebarStyle.Height(s.height).Render(content)
}
```

- [ ] **Step 2: Add Bubbletea dependency**

```bash
cd tui/viewer
go get github.com/charmbracelet/bubbletea
```

- [ ] **Step 3: Verify Go compiles**

```bash
cd tui/viewer
go build ./ui/...
```

Expected: No errors

- [ ] **Step 4: Commit**

```bash
git add tui/viewer/go.mod tui/viewer/go.sum tui/viewer/ui/sidebar.go
git commit -m "feat(tui): add sidebar component"
```

---

### Task 15: Implement Menu Component

**Files:**
- Create: `tui/viewer/ui/menu.go`

- [ ] **Step 1: Implement menu.go**

Create file `tui/viewer/ui/menu.go`:

```go
package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// MenuItem represents a menu option
type MenuItem struct {
	Key         string
	Label       string
	Description string
	Action      string // Command to send to bash
	Children    []MenuItem
}

// Menu handles menu navigation
type Menu struct {
	items    []MenuItem
	cursor   int
	parent   *Menu
	title    string
	width    int
	height   int
}

// NewMenu creates a new menu
func NewMenu(title string, items []MenuItem) Menu {
	return Menu{
		items: items,
		title: title,
	}
}

// SetSize sets the menu dimensions
func (m *Menu) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// SetItems updates the menu items
func (m *Menu) SetItems(items []MenuItem) {
	m.items = items
	if m.cursor >= len(items) {
		m.cursor = len(items) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
}

// Selected returns the currently selected item
func (m Menu) Selected() *MenuItem {
	if m.cursor < 0 || m.cursor >= len(m.items) {
		return nil
	}
	return &m.items[m.cursor]
}

// Update handles key events
func (m Menu) Update(msg tea.Msg) (Menu, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case "home":
			m.cursor = 0
		case "end":
			m.cursor = len(m.items) - 1
		}
	}
	return m, nil
}

// View renders the menu
func (m Menu) View() string {
	var b strings.Builder

	// Title
	if m.title != "" {
		b.WriteString(TitleStyle.Render(m.title))
		b.WriteString("\n")
		b.WriteString(lipgloss.NewStyle().
			Foreground(RedHatRed).
			Render(strings.Repeat("━", min(m.width-4, 40))))
		b.WriteString("\n\n")
	}

	// Items
	for i, item := range m.items {
		cursor := "  "
		style := MenuItemStyle
		if i == m.cursor {
			cursor = "> "
			style = MenuItemSelectedStyle
		}

		line := cursor + item.Label
		if item.Key != "" {
			line = cursor + "[" + item.Key + "] " + item.Label
		}

		b.WriteString(style.Render(line))
		
		if item.Description != "" && i == m.cursor {
			b.WriteString("\n")
			b.WriteString(lipgloss.NewStyle().
				Foreground(TextMuted).
				PaddingLeft(4).
				Render(item.Description))
		}
		b.WriteString("\n")
	}

	// Help
	b.WriteString("\n")
	b.WriteString(HelpKeyStyle.Render("↑/↓"))
	b.WriteString(HelpDescStyle.Render(" navigate  "))
	b.WriteString(HelpKeyStyle.Render("enter"))
	b.WriteString(HelpDescStyle.Render(" select  "))
	b.WriteString(HelpKeyStyle.Render("/"))
	b.WriteString(HelpDescStyle.Render(" search  "))
	b.WriteString(HelpKeyStyle.Render("q"))
	b.WriteString(HelpDescStyle.Render(" quit"))

	return MainStyle.Render(b.String())
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
```

- [ ] **Step 2: Verify Go compiles**

```bash
cd tui/viewer
go build ./ui/...
```

Expected: No errors

- [ ] **Step 3: Commit**

```bash
git add tui/viewer/ui/menu.go
git commit -m "feat(tui): add menu component"
```

---

### Task 16: Implement Activity Log Component

**Files:**
- Create: `tui/viewer/ui/activity.go`

- [ ] **Step 1: Add bubbles viewport dependency**

```bash
cd tui/viewer
go get github.com/charmbracelet/bubbles/viewport
```

- [ ] **Step 2: Implement activity.go**

Create file `tui/viewer/ui/activity.go`:

```go
package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/toddward/red-hat-engagement-kit/tui/viewer/protocol"
)

// ActivityEntry represents a single log entry
type ActivityEntry struct {
	Timestamp time.Time
	Event     protocol.EventType
	Text      string
	Tool      string
	Expanded  bool
}

// Activity displays streaming execution events
type Activity struct {
	viewport viewport.Model
	entries  []ActivityEntry
	width    int
	height   int
	running  bool
	cost     float64
}

// NewActivity creates a new activity log
func NewActivity() Activity {
	vp := viewport.New(80, 20)
	return Activity{
		viewport: vp,
		entries:  make([]ActivityEntry, 0),
	}
}

// SetSize sets the activity log dimensions
func (a *Activity) SetSize(width, height int) {
	a.width = width
	a.height = height
	a.viewport.Width = width
	a.viewport.Height = height - 4 // Reserve space for header/footer
}

// Clear clears the activity log
func (a *Activity) Clear() {
	a.entries = make([]ActivityEntry, 0)
	a.running = false
	a.cost = 0
	a.updateContent()
}

// SetRunning sets whether execution is active
func (a *Activity) SetRunning(running bool) {
	a.running = running
}

// AddEvent adds a new event to the log
func (a *Activity) AddEvent(event protocol.StreamEvent) {
	entry := ActivityEntry{
		Timestamp: time.Now(),
		Event:     event.Event,
		Text:      event.Text,
		Tool:      event.Tool,
	}

	// Handle special events
	switch event.Event {
	case protocol.EventCost:
		a.cost = event.Cost
		return // Don't add to log
	case protocol.EventComplete:
		a.running = false
		a.cost = event.Cost
	case protocol.EventToolResult:
		entry.Text = event.Output
	}

	a.entries = append(a.entries, entry)
	a.updateContent()
	a.viewport.GotoBottom()
}

func (a *Activity) updateContent() {
	var b strings.Builder

	for _, entry := range a.entries {
		ts := ActivityTimestampStyle.Render(entry.Timestamp.Format("15:04:05"))

		var content string
		switch entry.Event {
		case protocol.EventAssistant:
			content = ActivityAssistantStyle.Render(entry.Text)
		case protocol.EventToolUse:
			content = ActivityToolStyle.Render(fmt.Sprintf("▶ %s", entry.Tool))
		case protocol.EventToolResult:
			// Truncate long outputs
			text := entry.Text
			if len(text) > 200 {
				text = text[:200] + "..."
			}
			content = ActivityToolResultStyle.Render("  └─ " + strings.ReplaceAll(text, "\n", " "))
		case protocol.EventError:
			content = ActivityErrorStyle.Render("✗ " + entry.Text)
		case protocol.EventComplete:
			if entry.Text == "success" || entry.Text == "" {
				content = StatusCompleteStyle.Render("✓ Complete")
			} else {
				content = StatusErrorStyle.Render("✗ " + entry.Text)
			}
		default:
			content = entry.Text
		}

		b.WriteString(ts + " " + content + "\n")
	}

	a.viewport.SetContent(b.String())
}

// Update handles events
func (a Activity) Update(msg tea.Msg) (Activity, tea.Cmd) {
	var cmd tea.Cmd
	a.viewport, cmd = a.viewport.Update(msg)
	return a, cmd
}

// View renders the activity log
func (a Activity) View() string {
	var b strings.Builder

	// Header
	title := "Execution Log"
	if a.running {
		title = StatusRunningStyle.Render("● Running")
	}
	b.WriteString(TitleStyle.Render(title))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().
		Foreground(RedHatRed).
		Render(strings.Repeat("━", min(a.width-4, 60))))
	b.WriteString("\n\n")

	// Viewport
	b.WriteString(a.viewport.View())
	b.WriteString("\n")

	// Footer with cost
	if a.cost > 0 {
		costStr := fmt.Sprintf("Cost: $%.4f", a.cost)
		b.WriteString(lipgloss.NewStyle().Foreground(TextMuted).Render(costStr))
	}

	// Help
	b.WriteString("\n")
	b.WriteString(HelpKeyStyle.Render("↑/↓"))
	b.WriteString(HelpDescStyle.Render(" scroll  "))
	b.WriteString(HelpKeyStyle.Render("Ctrl+C"))
	b.WriteString(HelpDescStyle.Render(" cancel  "))
	b.WriteString(HelpKeyStyle.Render("Esc"))
	b.WriteString(HelpDescStyle.Render(" back"))

	return MainStyle.Render(b.String())
}
```

- [ ] **Step 3: Verify Go compiles**

```bash
cd tui/viewer
go build ./ui/...
```

Expected: No errors

- [ ] **Step 4: Commit**

```bash
git add tui/viewer/go.mod tui/viewer/go.sum tui/viewer/ui/activity.go
git commit -m "feat(tui): add activity log component"
```

---

### Task 17: Implement Input Component

**Files:**
- Create: `tui/viewer/ui/input.go`

- [ ] **Step 1: Add textinput dependency**

```bash
cd tui/viewer
go get github.com/charmbracelet/bubbles/textinput
```

- [ ] **Step 2: Implement input.go**

Create file `tui/viewer/ui/input.go`:

```go
package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Input handles user text input
type Input struct {
	textInput textinput.Model
	prompt    string
	options   []string
	cursor    int
	width     int
	height    int
}

// NewInput creates a new input component
func NewInput() Input {
	ti := textinput.New()
	ti.Placeholder = "Type your response..."
	ti.CharLimit = 500
	ti.Width = 60

	return Input{
		textInput: ti,
	}
}

// SetSize sets the input dimensions
func (i *Input) SetSize(width, height int) {
	i.width = width
	i.height = height
	i.textInput.Width = width - 10
}

// SetPrompt sets the question prompt
func (i *Input) SetPrompt(prompt string, options []string) {
	i.prompt = prompt
	i.options = options
	i.cursor = 0
	i.textInput.SetValue("")
	
	if len(options) == 0 {
		i.textInput.Focus()
	}
}

// Focus focuses the input
func (i *Input) Focus() {
	i.textInput.Focus()
}

// Blur unfocuses the input
func (i *Input) Blur() {
	i.textInput.Blur()
}

// Value returns the current input value
func (i Input) Value() string {
	if len(i.options) > 0 && i.cursor < len(i.options) {
		return i.options[i.cursor]
	}
	return i.textInput.Value()
}

// Update handles key events
func (i Input) Update(msg tea.Msg) (Input, tea.Cmd) {
	if len(i.options) > 0 {
		// Option selection mode
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "up", "k":
				if i.cursor > 0 {
					i.cursor--
				}
			case "down", "j":
				if i.cursor < len(i.options)-1 {
					i.cursor++
				}
			}
		}
		return i, nil
	}

	// Text input mode
	var cmd tea.Cmd
	i.textInput, cmd = i.textInput.Update(msg)
	return i, cmd
}

// View renders the input
func (i Input) View() string {
	var b strings.Builder

	// Prompt
	b.WriteString(InputLabelStyle.Render(i.prompt))
	b.WriteString("\n\n")

	if len(i.options) > 0 {
		// Render options
		for idx, opt := range i.options {
			cursor := "  "
			style := MenuItemStyle
			if idx == i.cursor {
				cursor = "> "
				style = MenuItemSelectedStyle
			}
			b.WriteString(style.Render(cursor + opt))
			b.WriteString("\n")
		}
	} else {
		// Render text input
		b.WriteString(InputStyle.Render(i.textInput.View()))
	}

	b.WriteString("\n\n")
	b.WriteString(HelpKeyStyle.Render("enter"))
	b.WriteString(HelpDescStyle.Render(" submit  "))
	b.WriteString(HelpKeyStyle.Render("Esc"))
	b.WriteString(HelpDescStyle.Render(" cancel"))

	return MainStyle.Render(b.String())
}
```

- [ ] **Step 3: Verify Go compiles**

```bash
cd tui/viewer
go build ./ui/...
```

Expected: No errors

- [ ] **Step 4: Commit**

```bash
git add tui/viewer/go.mod tui/viewer/go.sum tui/viewer/ui/input.go
git commit -m "feat(tui): add input component"
```

---

### Task 18: Implement Command Palette

**Files:**
- Create: `tui/viewer/ui/palette.go`

- [ ] **Step 1: Implement palette.go**

Create file `tui/viewer/ui/palette.go`:

```go
package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// PaletteItem represents a searchable action
type PaletteItem struct {
	Name        string
	Description string
	Action      string
	Category    string
}

// Palette is a fuzzy-search command palette
type Palette struct {
	input     textinput.Model
	items     []PaletteItem
	filtered  []PaletteItem
	cursor    int
	width     int
	height    int
}

// NewPalette creates a new command palette
func NewPalette() Palette {
	ti := textinput.New()
	ti.Placeholder = "Type to search..."
	ti.CharLimit = 100

	return Palette{
		input:    ti,
		items:    make([]PaletteItem, 0),
		filtered: make([]PaletteItem, 0),
	}
}

// SetItems sets the available items
func (p *Palette) SetItems(items []PaletteItem) {
	p.items = items
	p.filter()
}

// Open shows the palette
func (p *Palette) Open() {
	p.input.SetValue("")
	p.input.Focus()
	p.cursor = 0
	p.filter()
}

// Close hides the palette
func (p *Palette) Close() {
	p.input.Blur()
}

// Selected returns the selected item
func (p Palette) Selected() *PaletteItem {
	if p.cursor < 0 || p.cursor >= len(p.filtered) {
		return nil
	}
	return &p.filtered[p.cursor]
}

func (p *Palette) filter() {
	query := strings.ToLower(p.input.Value())
	
	if query == "" {
		p.filtered = p.items
		return
	}

	p.filtered = make([]PaletteItem, 0)
	for _, item := range p.items {
		name := strings.ToLower(item.Name)
		desc := strings.ToLower(item.Description)
		cat := strings.ToLower(item.Category)
		
		if strings.Contains(name, query) || 
		   strings.Contains(desc, query) ||
		   strings.Contains(cat, query) {
			p.filtered = append(p.filtered, item)
		}
	}

	// Reset cursor if out of bounds
	if p.cursor >= len(p.filtered) {
		p.cursor = len(p.filtered) - 1
	}
	if p.cursor < 0 {
		p.cursor = 0
	}
}

// Update handles events
func (p Palette) Update(msg tea.Msg) (Palette, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "ctrl+p":
			if p.cursor > 0 {
				p.cursor--
			}
			return p, nil
		case "down", "ctrl+n":
			if p.cursor < len(p.filtered)-1 {
				p.cursor++
			}
			return p, nil
		}
	}

	var cmd tea.Cmd
	p.input, cmd = p.input.Update(msg)
	p.filter()
	return p, cmd
}

// View renders the palette
func (p Palette) View() string {
	var b strings.Builder

	// Input
	b.WriteString(PaletteInputStyle.Render(p.input.View()))
	b.WriteString("\n")

	// Results (max 10)
	maxItems := 10
	if len(p.filtered) < maxItems {
		maxItems = len(p.filtered)
	}

	for i := 0; i < maxItems; i++ {
		item := p.filtered[i]
		
		style := PaletteItemStyle
		if i == p.cursor {
			style = PaletteItemSelectedStyle
		}

		line := style.Render(item.Name)
		if item.Category != "" {
			line += " " + lipgloss.NewStyle().Foreground(TextDim).Render("["+item.Category+"]")
		}
		if item.Description != "" {
			line += PaletteItemDescStyle.Render(" — " + item.Description)
		}

		b.WriteString(line)
		b.WriteString("\n")
	}

	if len(p.filtered) == 0 {
		b.WriteString(lipgloss.NewStyle().Foreground(TextMuted).Italic(true).Render("No matches"))
	}

	return PaletteStyle.Render(b.String())
}
```

- [ ] **Step 2: Verify Go compiles**

```bash
cd tui/viewer
go build ./ui/...
```

Expected: No errors

- [ ] **Step 3: Commit**

```bash
git add tui/viewer/ui/palette.go
git commit -m "feat(tui): add command palette component"
```

---

### Task 19: Implement App Shell

**Files:**
- Create: `tui/viewer/ui/app.go`

- [ ] **Step 1: Implement app.go**

Create file `tui/viewer/ui/app.go`:

```go
package ui

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/toddward/red-hat-engagement-kit/tui/viewer/protocol"
)

// ViewType identifies the current view
type ViewType int

const (
	ViewMenu ViewType = iota
	ViewActivity
	ViewInput
	ViewArtifacts
	ViewChecklists
)

// BashMsg wraps a message from the bash layer
type BashMsg struct {
	Response protocol.Response
}

// App is the root model
type App struct {
	// Layout
	width  int
	height int

	// Components
	sidebar  Sidebar
	menu     Menu
	activity Activity
	input    Input
	palette  Palette

	// State
	currentView    ViewType
	showPalette    bool
	engagement     string
	skills         []protocol.Skill
	engagements    []protocol.Engagement
	agents         []protocol.Agent

	// Communication
	bashIn  io.Writer
	bashOut *bufio.Scanner
	cmdID   int
	pending map[string]chan protocol.Response
	mu      sync.Mutex
}

// NewApp creates a new app
func NewApp(bashIn io.Writer, bashOut io.Reader) App {
	menu := NewMenu("Main Menu", []MenuItem{
		{Key: "1", Label: "Skills", Action: "show_skills"},
		{Key: "2", Label: "Agents", Action: "show_agents"},
		{Key: "3", Label: "Engagements", Action: "show_engagements"},
		{Key: "4", Label: "Artifacts", Action: "show_artifacts"},
		{Key: "5", Label: "Checklists", Action: "show_checklists"},
	})

	return App{
		sidebar:     NewSidebar(),
		menu:        menu,
		activity:    NewActivity(),
		input:       NewInput(),
		palette:     NewPalette(),
		currentView: ViewMenu,
		bashIn:      bashIn,
		bashOut:     bufio.NewScanner(bashOut),
		pending:     make(map[string]chan protocol.Response),
	}
}

// Init initializes the app
func (a App) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		a.sendCommand("init", nil),
		a.listenToBash(),
	)
}

// sendCommand sends a command to bash and returns a Cmd that waits for response
func (a *App) sendCommand(cmd string, args interface{}) tea.Cmd {
	return func() tea.Msg {
		a.mu.Lock()
		a.cmdID++
		id := fmt.Sprintf("cmd-%d", a.cmdID)
		a.mu.Unlock()

		command := protocol.Command{
			Cmd: cmd,
			ID:  id,
		}
		if args != nil {
			argsJSON, _ := json.Marshal(args)
			command.Args = argsJSON
		}

		data, _ := json.Marshal(command)
		fmt.Fprintln(a.bashIn, string(data))

		return nil
	}
}

// listenToBash reads messages from bash
func (a *App) listenToBash() tea.Cmd {
	return func() tea.Msg {
		if a.bashOut.Scan() {
			line := a.bashOut.Text()
			var resp protocol.Response
			if err := json.Unmarshal([]byte(line), &resp); err == nil {
				return BashMsg{Response: resp}
			}
		}
		return nil
	}
}

// Update handles all messages
func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.sidebar.SetSize(SidebarWidth, msg.Height)
		mainWidth := msg.Width - SidebarWidth - 2
		a.menu.SetSize(mainWidth, msg.Height)
		a.activity.SetSize(mainWidth, msg.Height)
		a.input.SetSize(mainWidth, msg.Height)

	case tea.KeyMsg:
		// Global keys
		switch msg.String() {
		case "ctrl+c", "q":
			if a.currentView == ViewMenu && !a.showPalette {
				return a, tea.Quit
			}
		case "/":
			if !a.showPalette {
				a.showPalette = true
				a.palette.Open()
				return a, nil
			}
		case "esc":
			if a.showPalette {
				a.showPalette = false
				a.palette.Close()
				return a, nil
			}
			if a.currentView != ViewMenu {
				a.currentView = ViewMenu
				return a, nil
			}
		case "enter":
			if a.showPalette {
				if item := a.palette.Selected(); item != nil {
					a.showPalette = false
					a.palette.Close()
					return a.handleAction(item.Action)
				}
			}
			if a.currentView == ViewMenu {
				if item := a.menu.Selected(); item != nil {
					return a.handleAction(item.Action)
				}
			}
		}

		// Delegate to active component
		if a.showPalette {
			var cmd tea.Cmd
			a.palette, cmd = a.palette.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			switch a.currentView {
			case ViewMenu:
				var cmd tea.Cmd
				a.menu, cmd = a.menu.Update(msg)
				cmds = append(cmds, cmd)
			case ViewActivity:
				var cmd tea.Cmd
				a.activity, cmd = a.activity.Update(msg)
				cmds = append(cmds, cmd)
			case ViewInput:
				var cmd tea.Cmd
				a.input, cmd = a.input.Update(msg)
				cmds = append(cmds, cmd)
			}
		}

	case BashMsg:
		cmds = append(cmds, a.handleBashResponse(msg.Response))
		// Continue listening
		cmds = append(cmds, a.listenToBash())
	}

	return a, tea.Batch(cmds...)
}

func (a *App) handleAction(action string) (tea.Model, tea.Cmd) {
	switch action {
	case "show_skills":
		items := make([]MenuItem, len(a.skills))
		for i, s := range a.skills {
			items[i] = MenuItem{
				Label:       "/" + s.Name,
				Description: s.Description,
				Action:      "run_skill:" + s.Name,
			}
		}
		a.menu = NewMenu("Skills", items)
	case "show_agents":
		items := make([]MenuItem, len(a.agents))
		for i, ag := range a.agents {
			items[i] = MenuItem{
				Label:       ag.Name,
				Description: ag.Role + " (" + ag.Model + ")",
				Action:      "run_agent:" + ag.Name,
			}
		}
		a.menu = NewMenu("Agents", items)
	case "show_engagements":
		items := make([]MenuItem, len(a.engagements))
		for i, e := range a.engagements {
			label := e.Slug
			if !e.HasContext {
				label += " (no context)"
			}
			items[i] = MenuItem{
				Label:  label,
				Action: "select_engagement:" + e.Slug,
			}
		}
		a.menu = NewMenu("Engagements", items)
	}
	return a, nil
}

func (a *App) handleBashResponse(resp protocol.Response) tea.Cmd {
	switch resp.Type {
	case "response":
		// Handle init response
		var initResp protocol.InitResponse
		if err := json.Unmarshal(resp.Payload, &initResp); err == nil {
			a.skills = initResp.Skills
			a.engagements = initResp.Engagements
			a.agents = initResp.Agents
			
			// Smart engagement detection
			if len(a.engagements) == 1 {
				a.engagement = a.engagements[0].Slug
				a.sidebar.SetEngagement(a.engagement)
			}

			// Build palette items
			items := make([]PaletteItem, 0)
			for _, s := range a.skills {
				items = append(items, PaletteItem{
					Name:        "/" + s.Name,
					Description: s.Description,
					Action:      "run_skill:" + s.Name,
					Category:    "Skills",
				})
			}
			for _, ag := range a.agents {
				items = append(items, PaletteItem{
					Name:        ag.Name,
					Description: ag.Role,
					Action:      "run_agent:" + ag.Name,
					Category:    "Agents",
				})
			}
			a.palette.SetItems(items)
		}

	case "event":
		var event protocol.StreamEvent
		if err := json.Unmarshal(resp.Payload, &event); err == nil {
			a.activity.AddEvent(event)
		}
	}

	return nil
}

// View renders the app
func (a App) View() string {
	// Sidebar
	sidebar := a.sidebar.View()

	// Main content
	var main string
	switch a.currentView {
	case ViewMenu:
		main = a.menu.View()
	case ViewActivity:
		main = a.activity.View()
	case ViewInput:
		main = a.input.View()
	}

	// Combine sidebar and main
	content := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, main)

	// Overlay palette if visible
	if a.showPalette {
		paletteView := a.palette.View()
		// Center the palette
		x := (a.width - 60) / 2
		y := 5
		content = placeOverlay(x, y, paletteView, content)
	}

	return content
}

// placeOverlay places an overlay on top of content
func placeOverlay(x, y int, overlay, background string) string {
	bgLines := strings.Split(background, "\n")
	ovLines := strings.Split(overlay, "\n")

	for i, ovLine := range ovLines {
		bgY := y + i
		if bgY >= 0 && bgY < len(bgLines) {
			bgLine := bgLines[bgY]
			// Simple overlay - just replace characters
			if x >= 0 && x < len(bgLine) {
				bgLines[bgY] = bgLine[:x] + ovLine
			}
		}
	}

	return strings.Join(bgLines, "\n")
}

import "strings"
```

- [ ] **Step 2: Fix import and verify compile**

The strings import is duplicated. Remove the last line and verify:

```bash
cd tui/viewer
go build ./ui/...
```

- [ ] **Step 3: Commit**

```bash
git add tui/viewer/ui/app.go
git commit -m "feat(tui): add main app shell"
```

---

### Task 20: Implement Main Entry Point

**Files:**
- Create: `tui/viewer/main.go`

- [ ] **Step 1: Implement main.go**

Create file `tui/viewer/main.go`:

```go
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/toddward/red-hat-engagement-kit/tui/viewer/ui"
)

func main() {
	// stdin/stdout connect to bash core
	app := ui.NewApp(os.Stdout, os.Stdin)

	p := tea.NewProgram(app, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
```

- [ ] **Step 2: Build viewer**

```bash
cd tui/viewer
go build -o ../bin/tui-viewer .
```

Expected: Binary created at `tui/bin/tui-viewer`

- [ ] **Step 3: Commit**

```bash
git add tui/viewer/main.go
git commit -m "feat(tui): add viewer main entry point"
```

---

## Phase 4: Integration

### Task 21: Create Setup Script

**Files:**
- Create: `tui/setup.sh`

- [ ] **Step 1: Implement setup.sh**

Create file `tui/setup.sh`:

```bash
#!/usr/bin/env bash
set -euo pipefail

echo "Red Hat Engagement Kit TUI Setup"
echo "================================"
echo

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BIN_DIR="$SCRIPT_DIR/bin"

mkdir -p "$BIN_DIR"

# Detect platform
PLATFORM=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
[[ "$ARCH" == "x86_64" ]] && ARCH="amd64"
[[ "$ARCH" == "aarch64" || "$ARCH" == "arm64" ]] && ARCH="arm64"

echo "Platform: $PLATFORM-$ARCH"
echo

# Check/install jq
echo "Checking jq..."
if command -v jq &>/dev/null; then
    echo "  ✓ jq found in PATH"
elif [[ -x "$BIN_DIR/jq" ]]; then
    echo "  ✓ jq found in bin/"
else
    echo "  ↓ Downloading jq..."
    JQ_URL="https://github.com/jqlang/jq/releases/download/jq-1.7.1/jq-${PLATFORM}-${ARCH}"
    if curl -sL "$JQ_URL" -o "$BIN_DIR/jq" 2>/dev/null; then
        chmod +x "$BIN_DIR/jq"
        echo "  ✓ jq installed to bin/"
    else
        echo "  ✗ Failed to download jq"
        echo "    Manual install: https://jqlang.github.io/jq/download/"
        exit 1
    fi
fi

# Check/build tui-viewer
echo "Checking tui-viewer..."
if [[ -x "$BIN_DIR/tui-viewer" ]]; then
    echo "  ✓ tui-viewer found in bin/"
else
    if command -v go &>/dev/null; then
        echo "  → Building tui-viewer from source..."
        (cd "$SCRIPT_DIR/viewer" && go build -o "$BIN_DIR/tui-viewer" .)
        echo "  ✓ tui-viewer built"
    else
        echo "  ✗ Go not found, cannot build tui-viewer"
        echo "    Install Go: https://go.dev/dl/"
        echo "    Or download pre-built binary to bin/tui-viewer"
        exit 1
    fi
fi

echo
echo "Setup complete!"
echo "Run: ./tui.sh"
```

- [ ] **Step 2: Make executable**

```bash
chmod +x tui/setup.sh
```

- [ ] **Step 3: Commit**

```bash
git add tui/setup.sh
git commit -m "feat(tui): add setup script"
```

---

### Task 22: Create Entry Point Script

**Files:**
- Create: `tui/tui.sh`

- [ ] **Step 1: Implement tui.sh**

Create file `tui/tui.sh`:

```bash
#!/usr/bin/env bash
set -euo pipefail

TUI_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BIN_DIR="$TUI_DIR/bin"
CORE_DIR="$TUI_DIR/core"

# Check if setup has been run
if [[ ! -x "$BIN_DIR/tui-viewer" ]]; then
    echo "TUI not set up. Running setup..."
    "$TUI_DIR/setup.sh"
    echo
fi

# Add bin to PATH
export PATH="$BIN_DIR:$PATH"

# Change to project root (parent of tui/)
cd "$TUI_DIR/.."

# Launch with bidirectional communication
# Using coproc for bash-to-viewer communication
coproc CORE {
    bash "$CORE_DIR/main.sh"
}

# Connect viewer to core via pipes
"$BIN_DIR/tui-viewer" <&"${CORE[0]}" >&"${CORE[1]}"

# Cleanup
kill "${CORE_PID}" 2>/dev/null || true
```

- [ ] **Step 2: Make executable**

```bash
chmod +x tui/tui.sh
```

- [ ] **Step 3: Commit**

```bash
git add tui/tui.sh
git commit -m "feat(tui): add main entry point script"
```

---

### Task 23: Create README

**Files:**
- Create: `tui/README.md`

- [ ] **Step 1: Write README**

Create file `tui/README.md`:

```markdown
# Red Hat Engagement Kit TUI

A polished terminal user interface for running engagements.

## Quick Start

```bash
# First time setup
./setup.sh

# Run TUI
./tui.sh
```

## Requirements

- Bash 4+
- Go 1.21+ (for building, not required if using pre-built binaries)

## Architecture

The TUI uses a hybrid architecture:

- **Bash Logic Layer** (`core/`) — Handles business operations: reading skills, managing engagements, invoking Claude CLI
- **Go Viewer Layer** (`viewer/`) — Handles terminal rendering using Bubbletea

They communicate via JSON messages over stdin/stdout.

## Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `/` | Open command palette |
| `↑/↓` | Navigate |
| `Enter` | Select |
| `Esc` | Go back / Close |
| `q` | Quit |
| `Ctrl+C` | Cancel execution |

## Development

```bash
# Modify bash logic (changes take effect immediately)
vim core/lib/skills.sh

# Modify Go viewer (rebuild required)
cd viewer
go build -o ../bin/tui-viewer .
```

## Air-Gapped Environments

Pre-populate `bin/` with:
- `jq` - Download from https://jqlang.github.io/jq/download/
- `tui-viewer` - Build on a connected machine with `go build`

The `setup.sh` script will skip downloads if these exist.
```

- [ ] **Step 2: Commit**

```bash
git add tui/README.md
git commit -m "docs(tui): add README"
```

---

### Task 24: Final Integration Test

**Files:** None (manual testing)

- [ ] **Step 1: Run all bash tests**

```bash
cd tui
./tests/test_protocol.sh
./tests/test_state.sh
./tests/test_skills.sh
./tests/test_engagements.sh
./tests/test_phase.sh
./tests/test_artifacts.sh
./tests/test_checklists.sh
./tests/test_agents.sh
```

Expected: All tests pass

- [ ] **Step 2: Build Go viewer**

```bash
cd tui/viewer
go build -o ../bin/tui-viewer .
```

Expected: No errors

- [ ] **Step 3: Run setup**

```bash
cd tui
./setup.sh
```

Expected: Setup completes successfully

- [ ] **Step 4: Manual smoke test**

```bash
./tui.sh
```

Expected: TUI launches, shows main menu

- [ ] **Step 5: Final commit**

```bash
git add -A
git commit -m "feat(tui): complete TUI implementation

- Bash logic layer with full skill/engagement/agent support
- Go viewer with split-pane layout and command palette
- JSON protocol for bash-go communication
- Setup script with dependency bundling
- Full test coverage for bash modules"
```

---

## Summary

This plan implements the TUI in 24 tasks across 4 phases:

1. **Foundation** (Tasks 1-3): Directory structure, protocol definitions
2. **Bash Logic** (Tasks 4-12): All business logic modules
3. **Go Viewer** (Tasks 13-20): UI components and app shell
4. **Integration** (Tasks 21-24): Scripts, docs, testing

Each task follows TDD: write test → verify fail → implement → verify pass → commit.

Total estimated commits: 24
