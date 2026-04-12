# TUI -- Terminal Interface for Red Hat Engagement Kit

A polished terminal UI that wraps the engagement kit's skills, agents, and artifacts into a single navigable interface. Built for Red Hat architects who want a visual workflow without leaving the terminal.

## Quick Start

```bash
cd tui/
./setup.sh   # Downloads jq, builds the Go viewer
./tui.sh     # Launches the TUI
```

`setup.sh` is idempotent -- it skips anything already present in `bin/`.

## Requirements

| Dependency | Version | When Needed |
|------------|---------|-------------|
| Bash       | 4+      | Always      |
| jq         | 1.7+    | Always (auto-downloaded by `setup.sh`) |
| Go         | 1.21+   | Building `tui-viewer` from source only  |
| Claude CLI | Latest  | Skill/agent execution |

Go is not required at runtime. If you have a pre-built `tui-viewer` binary in `bin/`, Go is never invoked.

## Architecture

The TUI is a hybrid of bash and Go, connected by a JSON protocol over stdin/stdout.

```
 ┌──────────────────────────────────┐
 │  tui.sh (launcher)               │
 │    Starts bash core as coproc     │
 │    Connects viewer via pipes      │
 └───────┬──────────────┬───────────┘
         │ stdout        │ stdin
         ▼              ▼
 ┌───────────────┐  ┌───────────────┐
 │  core/main.sh │  │  tui-viewer   │
 │  Bash logic   │◄─┤  Go / Bubble  │
 │  layer        │──►  Tea renderer │
 └───────────────┘  └───────────────┘
   Reads skills,      Renders UI to
   engagements,       /dev/tty, sends
   runs Claude CLI    commands as JSON
```

**Why two processes?** Bash handles file I/O, skill discovery, and Claude CLI orchestration -- things shell scripts do well. Go (via Bubble Tea) handles terminal rendering, keyboard input, and layout -- things TUI frameworks do well. Neither process manages the other's domain.

### Protocol

Every message is a single JSON line. The viewer sends commands:

```json
{"cmd": "list_skills", "id": "cmd-1", "args": {}}
```

The core responds with one of three message types:

| Type       | Purpose                              |
|------------|--------------------------------------|
| `response` | Direct reply to a command (has `id`) |
| `event`    | Streaming output during execution    |
| `error`    | Error with optional correlation `id` |

## Keyboard Shortcuts

| Key        | Action                                |
|------------|---------------------------------------|
| `/`        | Open command palette (fuzzy search)   |
| Up / Down  | Navigate menu or palette              |
| Enter      | Select item or confirm input          |
| Esc        | Close palette or return to main menu  |
| q          | Quit (from main menu)                 |
| Ctrl+C     | Cancel running skill, or quit         |
| 1-5        | Jump to menu item by number           |

The command palette supports type-ahead filtering across skill names, agent names, and descriptions.

## Features

- **Skill execution** -- Run `/setup`, `/discover-infrastructure`, `/assess-app-portfolio`, and `/build-deliverable-deck` with streaming output
- **Agent invocation** -- Prompt Architect, Senior Developer, QA, or Documentation agents with free-text input
- **Engagement management** -- Browse, select, and create engagements; auto-selects when only one exists
- **Artifact browsing** -- Navigate the engagement directory tree and read files inline
- **Checklist tracking** -- View and toggle checklist items from `knowledge/checklists/`
- **Phase detection** -- Sidebar shows current engagement phase (pre-engagement, live, leave-behind) with artifact counts

## Directory Structure

```
tui/
├── tui.sh              # Entry point -- launches core + viewer
├── setup.sh            # One-time setup (jq + viewer binary)
├── bin/                # Runtime binaries (gitignored)
│   ├── jq              # JSON processor
│   └── tui-viewer      # Compiled Go viewer
├── core/               # Bash logic layer
│   ├── main.sh         # Event loop -- reads commands, dispatches handlers
│   ├── protocol.sh     # JSON message helpers (send_response, send_event, send_error)
│   └── lib/            # Domain modules
│       ├── skills.sh       # Skill discovery and metadata
│       ├── engagements.sh  # Engagement listing and selection
│       ├── agents.sh       # Agent discovery
│       ├── artifacts.sh    # File tree and content reading
│       ├── checklists.sh   # Checklist parsing and toggling
│       ├── phase.sh        # Engagement phase detection
│       ├── state.sh        # Persistent state (JSON file)
│       └── claude.sh       # Claude CLI execution and streaming
├── viewer/             # Go terminal renderer
│   ├── main.go         # Entry point -- wires pipes to Bubble Tea
│   ├── go.mod
│   ├── protocol/       # Shared message types
│   │   └── messages.go
│   └── ui/             # Bubble Tea components
│       ├── app.go          # Root model, routing, bash communication
│       ├── sidebar.go      # Engagement info and phase display
│       ├── menu.go         # Navigable menu lists
│       ├── activity.go     # Streaming execution log
│       ├── input.go        # Text input with prompt
│       ├── palette.go      # Fuzzy command palette overlay
│       └── styles.go       # Red Hat brand colors and layout
└── tests/              # Bash unit tests for core modules
    ├── test_protocol.sh
    ├── test_skills.sh
    ├── test_engagements.sh
    └── ...
```

## Development

**Bash scripts (core/)** -- Edit and re-run. Changes take effect on next `./tui.sh` launch. No build step.

**Go viewer (viewer/)** -- Requires rebuild after changes:

```bash
cd viewer/
go build -o ../bin/tui-viewer .
```

**Tests** -- Run individual test files directly:

```bash
bash tests/test_protocol.sh
bash tests/test_skills.sh
```

## Air-Gapped Environments

For disconnected use, pre-populate `bin/` before deployment:

```bash
# On a connected machine, build for the target platform
cd viewer/
GOOS=linux GOARCH=amd64 go build -o ../bin/tui-viewer .

# Download jq for the target
curl -sL https://github.com/jqlang/jq/releases/download/jq-1.7.1/jq-linux-amd64 -o ../bin/jq
chmod +x ../bin/jq
```

Transfer the entire `tui/` directory (including `bin/`) to the air-gapped host. `setup.sh` will detect existing binaries and skip downloads.
