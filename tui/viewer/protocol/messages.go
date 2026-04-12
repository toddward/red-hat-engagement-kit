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
	Cost    float64   `json:"totalCost,omitempty"`
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
	Name          string             `json:"name"`
	FileName      string             `json:"fileName"`
	Sections      []ChecklistSection `json:"sections"`
	CompletionPct int                `json:"completionPercent"`
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
