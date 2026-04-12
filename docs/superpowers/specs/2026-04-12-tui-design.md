# TUI Design Specification

**Date:** 2026-04-12  
**Status:** Draft  
**Author:** Todd Wardzinski + Claude

## Overview

A polished terminal user interface (TUI) for the Red Hat Engagement Kit that provides full feature parity with the web GUI while running entirely from the command line with minimal dependencies.

## Goals

- Run from command line without requiring Node.js, Python, or other runtimes
- Support Red Hat architects on managed workstations and field engineers on varied systems (Linux/Mac)
- Bundle missing dependencies automatically; support air-gapped environments via pre-bundled binaries
- Full feature parity with web GUI: skills, agents, engagements, phases, artifacts, checklists

## Architecture

### Hybrid Approach: Bash Logic + Go Viewer

The TUI uses a two-layer architecture:

1. **Bash Logic Layer (`core/`)** — Handles all business operations: reading skills, managing engagements, invoking Claude CLI, parsing output
2. **Go Viewer Layer (`viewer/`)** — Handles all terminal rendering using Bubbletea: split-pane layout, streaming display, user input

The layers communicate via a simple JSON protocol over stdin/stdout.

```
┌─────────────────────────────────────────────────────────────────┐
│                         User Terminal                           │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Go Viewer (tui-viewer)                       │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────────┐ │
│  │ Split Panes │  │  Input      │  │  Event Stream Display   │ │
│  │ (Bubbletea) │  │  Handling   │  │  (Activity Log)         │ │
│  └─────────────┘  └─────────────┘  └─────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
          │ stdin/stdout (JSON protocol)          ▲
          ▼                                       │
┌─────────────────────────────────────────────────────────────────┐
│                    Bash Logic (tui-core)                        │
│  ┌──────────────┐  ┌──────────────┐  ┌───────────────────────┐ │
│  │ Skill Reader │  │ Engagement   │  │ Claude Executor       │ │
│  │              │  │ Manager      │  │ (spawns claude CLI)   │ │
│  └──────────────┘  └──────────────┘  └───────────────────────┘ │
│  ┌──────────────┐  ┌──────────────┐  ┌───────────────────────┐ │
│  │ Phase        │  │ Artifact     │  │ Agent Invoker         │ │
│  │ Detector     │  │ Browser      │  │                       │ │
│  └──────────────┘  └──────────────┘  └───────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
          │                                       │
          ▼                                       ▼
┌─────────────────────┐              ┌────────────────────────────┐
│  Filesystem         │              │  Claude CLI                │
│  (.claude/skills,   │              │  (claude --print           │
│   engagements/)     │              │   --output-format stream-json)
└─────────────────────┘              └────────────────────────────┘
```

### Why This Architecture

- **Bash for logic** — Readable and modifiable by architects without a toolchain; direct access to filesystem and Claude CLI
- **Go for rendering** — Bubbletea is the best TUI framework available; compiles to a single static binary
- **JSON protocol** — Simple, debuggable, language-agnostic communication

## Directory Structure

```
tui/
├── bin/                          # Bundled dependencies (gitignored, populated by setup)
│   ├── jq                        # JSON parsing
│   └── tui-viewer                # Compiled Go viewer (built or downloaded)
│
├── core/                         # Bash logic layer
│   ├── main.sh                   # Entry point, orchestrates everything
│   ├── lib/
│   │   ├── skills.sh             # List skills, read SKILL.md, parse frontmatter
│   │   ├── engagements.sh        # List/create/select engagements
│   │   ├── phase.sh              # Detect current phase from CONTEXT.md
│   │   ├── artifacts.sh          # List/read files in engagement directories
│   │   ├── checklists.sh         # Parse/update markdown checklists
│   │   ├── agents.sh             # Agent definitions and invocation
│   │   ├── claude.sh             # Spawn claude CLI, stream output
│   │   └── state.sh              # Persistent state management
│   └── protocol.sh               # JSON message encoding/decoding for viewer
│
├── viewer/                       # Go viewer layer
│   ├── main.go                   # Entry point
│   ├── go.mod
│   ├── go.sum
│   ├── ui/
│   │   ├── app.go                # Main Bubbletea app, layout management
│   │   ├── sidebar.go            # Context sidebar (engagement, phase, agent)
│   │   ├── activity.go           # Activity log for streaming events
│   │   ├── menu.go               # Menu view
│   │   ├── palette.go            # Command palette (fuzzy search)
│   │   ├── input.go              # User input during execution
│   │   └── styles.go             # Lipgloss styling (Red Hat brand colors)
│   └── protocol/
│       └── messages.go           # JSON protocol types
│
├── setup.sh                      # First-run setup: detect/download dependencies
├── tui.sh                        # User-facing entry point
└── README.md                     # Usage documentation
```

## Bash Logic Layer

### Core Modules

| Module | Responsibility |
|--------|----------------|
| `skills.sh` | `list_skills`, `get_skill <name>` — reads `.claude/skills/*/SKILL.md`, extracts frontmatter |
| `engagements.sh` | `list_engagements`, `create_engagement`, `get_engagement <slug>` — manages `engagements/` |
| `phase.sh` | `detect_phase <slug>` — parses CONTEXT.md to determine pre-engagement/live/leave-behind |
| `artifacts.sh` | `list_artifacts <slug>`, `read_artifact <path>` — traverses engagement subdirectories |
| `checklists.sh` | `list_checklists`, `get_checklist <name>`, `toggle_item <file> <line>` — parses markdown checkboxes |
| `agents.sh` | `list_agents`, `get_agent <name>` — reads `.claude/agents/`, builds invocation commands |
| `claude.sh` | `execute_skill <skill> <engagement>`, `execute_agent <agent> <prompt>` — spawns claude CLI, streams JSON |
| `state.sh` | `load_state`, `save_state`, `get_state <key>`, `set_state <key> <value>` — reads/writes `.tui-state.json` |

### Event Loop Pattern

```bash
# main.sh (simplified)
while IFS= read -r message; do
    cmd=$(echo "$message" | jq -r '.cmd')
    
    case "$cmd" in
        list_skills)     list_skills | send_response ;;
        execute_skill)   execute_skill "$(jq -r '.skill')" "$(jq -r '.engagement')" ;;
        user_input)      handle_user_input "$(jq -r '.text')" ;;
        # ... etc
    esac
done
```

### Claude Execution Flow

```bash
# claude.sh (simplified)
execute_skill() {
    local skill="$1" engagement="$2"
    
    claude --print \
           --output-format stream-json \
           --verbose \
           --permission-mode acceptEdits \
           "Run /${skill}. Use engagement at engagements/${engagement}/." \
    | while IFS= read -r line; do
        # Parse and forward to viewer as activity events
        parse_claude_event "$line" | send_to_viewer
    done
}
```

## Go Viewer Layer

### Component Hierarchy

```
App (root model)
├── Sidebar (left pane, fixed width)
│   ├── Logo/Title
│   ├── Engagement name & phase indicator
│   ├── Current agent (if active)
│   └── Quick stats (artifact counts)
│
├── Main (right pane, flexible)
│   ├── MenuView (when navigating)
│   ├── ActivityView (during execution)
│   ├── InputView (when Claude asks a question)
│   ├── ArtifactView (browsing files)
│   └── ChecklistView (viewing/editing checklists)
│
└── CommandPalette (overlay, activated by /)
    ├── Fuzzy search input
    └── Filtered results list
```

### Bubbletea Model Structure

```go
type App struct {
    sidebar     Sidebar
    view        View           // Current active view
    palette     CommandPalette
    showPalette bool
    
    bashIn      io.Writer      // Send commands to bash
    bashOut     <-chan Message // Receive responses from bash
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == "/" && !a.showPalette {
            a.showPalette = true
            return a, nil
        }
    case BashMessage:
        return a.handleBashMessage(msg)
    }
    return a.view.Update(msg)
}
```

### Styling (Red Hat Brand)

```go
var (
    RedHatRed    = lipgloss.Color("#EE0000")
    RedHatBlack  = lipgloss.Color("#151515")
    TextPrimary  = lipgloss.Color("#E8E8E8")
    TextMuted    = lipgloss.Color("#888888")
    
    SidebarStyle = lipgloss.NewStyle().
        Width(30).
        BorderRight(true).
        BorderStyle(lipgloss.NormalBorder()).
        BorderForeground(RedHatRed)
)
```

## Communication Protocol

### Message Format

```json
// Bash → Viewer (responses, events)
{
    "type": "response" | "event" | "error",
    "id": "optional-correlation-id",
    "payload": { ... }
}

// Viewer → Bash (commands)
{
    "cmd": "command_name",
    "id": "correlation-id",
    "args": { ... }
}
```

### Command Catalog

| Command | Args | Response |
|---------|------|----------|
| `init` | — | `{ engagements, skills, agents, state }` |
| `list_skills` | — | `{ skills: [{ name, description, path }] }` |
| `list_engagements` | — | `{ engagements: [{ slug, hasContext }] }` |
| `get_phase` | `{ engagement }` | `{ phase, artifactCounts }` |
| `list_artifacts` | `{ engagement }` | `{ tree: [...] }` |
| `read_artifact` | `{ path }` | `{ content }` |
| `list_checklists` | — | `{ checklists: [...] }` |
| `toggle_checklist` | `{ file, line }` | `{ success }` |
| `execute_skill` | `{ skill, engagement }` | Streams events until `{ type: "complete" }` |
| `execute_agent` | `{ agent, prompt, engagement }` | Streams events until `{ type: "complete" }` |
| `user_input` | `{ text }` | Sent during execution when user responds |
| `cancel` | — | Cancels running execution |
| `set_state` | `{ key, value }` | `{ success }` |

### Event Types (During Execution)

| Event Type | Payload | Display |
|------------|---------|---------|
| `assistant` | `{ text }` | Assistant message in activity log |
| `tool_use` | `{ tool, input }` | Tool invocation line |
| `tool_result` | `{ tool, output }` | Tool result (collapsed by default) |
| `question` | `{ text, options? }` | Prompts user input view |
| `cost` | `{ usd }` | Running cost indicator |
| `complete` | `{ status, totalCost }` | Execution finished |
| `error` | `{ message }` | Error display |

## Features & Views

### Main Menu Structure

```
Red Hat Engagement Kit
━━━━━━━━━━━━━━━━━━━━━━

[1] Skills
    ├── /setup
    ├── /discover-infrastructure
    ├── /assess-app-portfolio
    └── /build-deliverable-deck

[2] Agents
    ├── Architect (Opus)
    ├── Senior Developer (Sonnet)
    ├── QA Specialist (Opus)
    └── Documentation Specialist (Opus)

[3] Engagements
    ├── Switch engagement
    ├── Create new
    └── View CONTEXT.md

[4] Artifacts
    └── Browse discovery/, assessments/, deliverables/

[5] Checklists
    ├── Architect checklist
    ├── Consultant checklist
    └── PM checklist

[/] Command palette
[q] Quit
```

### View Flows

| Action | Flow |
|--------|------|
| Run a skill | Menu → Select skill → Confirm engagement → Activity view (streaming) → Complete → Return to menu |
| Invoke an agent | Menu → Select agent → Enter prompt → Activity view (streaming) → Complete |
| Respond to question | Activity view → Question appears → Input view overlay → Submit → Resume streaming |
| Browse artifacts | Menu → Artifact tree → Select file → Content preview |
| Edit checklist | Menu → Select checklist → View items → Toggle with Enter/Space → Auto-saves |
| Command palette | Press `/` anywhere → Fuzzy search → Select action → Execute |

### Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `/` | Open command palette |
| `Esc` | Close overlay / go back |
| `q` | Quit (with confirmation if execution running) |
| `Ctrl+C` | Cancel current execution |
| `?` | Show help overlay |
| `Tab` | Cycle focus between panes |

## Startup Behavior

**Smart detection:**
- If no engagements exist → prompt to create one
- If one engagement exists → load it automatically
- If multiple engagements exist → show picker

**Persistent state (`.tui-state.json`):**
- Last active engagement
- Preferred view/layout
- Execution history

## Dependency Management

### Required Dependencies

| Tool | Purpose | Size |
|------|---------|------|
| `jq` | JSON parsing in bash | ~1.5MB |
| `tui-viewer` | Go viewer binary | ~10MB |

### Setup Flow

```bash
#!/usr/bin/env bash
# setup.sh - Detects platform, downloads missing dependencies

PLATFORM=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
[[ "$ARCH" == "x86_64" ]] && ARCH="amd64"
[[ "$ARCH" == "aarch64" || "$ARCH" == "arm64" ]] && ARCH="arm64"

BIN_DIR="$(dirname "$0")/bin"
mkdir -p "$BIN_DIR"

# Download jq if missing
if ! command -v jq &>/dev/null && [[ ! -x "$BIN_DIR/jq" ]]; then
    curl -sL "https://github.com/jqlang/jq/releases/download/jq-1.7.1/jq-${PLATFORM}-${ARCH}" \
        -o "$BIN_DIR/jq"
    chmod +x "$BIN_DIR/jq"
fi

# Build or download tui-viewer
if [[ ! -x "$BIN_DIR/tui-viewer" ]]; then
    if command -v go &>/dev/null; then
        (cd viewer && go build -o "../$BIN_DIR/tui-viewer" .)
    else
        curl -sL "https://github.com/<repo>/releases/latest/download/tui-viewer-${PLATFORM}-${ARCH}" \
            -o "$BIN_DIR/tui-viewer"
        chmod +x "$BIN_DIR/tui-viewer"
    fi
fi
```

### Air-Gapped Support

- Pre-populate `bin/` with all platform binaries before deployment
- `setup.sh` skips downloads if binaries already exist
- Document manual download URLs in README

## Build & Distribution

### Development Workflow

```bash
# First time setup
cd tui && ./setup.sh

# Run TUI
./tui.sh

# Modify bash logic — changes take effect immediately
vim core/lib/skills.sh

# Modify Go viewer — rebuild required
cd viewer && go build -o ../bin/tui-viewer .
```

### Release Artifacts

```
tui-release-v1.0.0/
├── tui.sh
├── setup.sh
├── core/                    # Bash scripts
├── bin/
│   ├── jq-linux-amd64
│   ├── jq-linux-arm64
│   ├── jq-darwin-amd64
│   ├── jq-darwin-arm64
│   ├── tui-viewer-linux-amd64
│   ├── tui-viewer-linux-arm64
│   ├── tui-viewer-darwin-amd64
│   └── tui-viewer-darwin-arm64
└── README.md
```

### Entry Point

```bash
#!/usr/bin/env bash
# tui.sh
set -euo pipefail

TUI_DIR="$(cd "$(dirname "$0")" && pwd)"
BIN_DIR="$TUI_DIR/bin"
CORE_DIR="$TUI_DIR/core"

# Ensure setup has run
if [[ ! -x "$BIN_DIR/tui-viewer" ]]; then
    echo "First run detected. Running setup..."
    "$TUI_DIR/setup.sh"
fi

export PATH="$BIN_DIR:$PATH"

# Launch with bidirectional pipes
coproc CORE { bash "$CORE_DIR/main.sh"; }
"$BIN_DIR/tui-viewer" <&"${CORE[0]}" >&"${CORE[1]}"
```

## Out of Scope

- Windows support (Linux and Mac only)
- Web browser integration
- Multi-user/server mode
- Plugin system for custom views

## Open Questions

None — design is complete and approved.

## Appendix: Relation to Web GUI

This TUI reuses the architectural patterns from `feature/gui`:

| GUI Component | TUI Equivalent |
|---------------|----------------|
| Express routes | Bash command handlers in `main.sh` |
| Service layer (skillReader, etc.) | Bash modules in `core/lib/` |
| claudeProvider.ts | `core/lib/claude.sh` |
| React components | Go Bubbletea views in `viewer/ui/` |
| SSE streaming | JSON protocol over stdin/stdout |
| CSS modules | Lipgloss styles in `viewer/ui/styles.go` |

The core difference is that the TUI removes the HTTP layer entirely — the viewer and logic communicate directly via pipes.
