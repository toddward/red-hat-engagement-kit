package protocol

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestCommand_MarshalJSON(t *testing.T) {
	t.Run("full command with all fields", func(t *testing.T) {
		args := json.RawMessage(`{"key":"value"}`)
		cmd := Command{
			Cmd:  "run-skill",
			ID:   "abc-123",
			Args: args,
		}

		data, err := json.Marshal(cmd)
		if err != nil {
			t.Fatalf("unexpected marshal error: %v", err)
		}

		var got map[string]json.RawMessage
		if err := json.Unmarshal(data, &got); err != nil {
			t.Fatalf("unexpected unmarshal error: %v", err)
		}

		assertJSONString(t, got["cmd"], "run-skill")
		assertJSONString(t, got["id"], "abc-123")

		if string(got["args"]) != `{"key":"value"}` {
			t.Errorf("args: got %s, want %s", got["args"], `{"key":"value"}`)
		}
	})

	t.Run("nil Args omitted from output", func(t *testing.T) {
		cmd := Command{Cmd: "init"}

		data, err := json.Marshal(cmd)
		if err != nil {
			t.Fatalf("unexpected marshal error: %v", err)
		}

		var got map[string]json.RawMessage
		if err := json.Unmarshal(data, &got); err != nil {
			t.Fatalf("unexpected unmarshal error: %v", err)
		}

		if _, ok := got["args"]; ok {
			t.Error("args field should be omitted when nil")
		}
		if _, ok := got["id"]; ok {
			t.Error("id field should be omitted when empty")
		}
	})
}

func TestCommand_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  Command
	}{
		{
			name:  "minimal cmd only",
			input: `{"cmd":"init"}`,
			want:  Command{Cmd: "init"},
		},
		{
			name:  "cmd with id",
			input: `{"cmd":"run","id":"req-42"}`,
			want:  Command{Cmd: "run", ID: "req-42"},
		},
		{
			name:  "cmd with args",
			input: `{"cmd":"run-skill","id":"req-1","args":{"skill":"setup"}}`,
			want: Command{
				Cmd:  "run-skill",
				ID:   "req-1",
				Args: json.RawMessage(`{"skill":"setup"}`),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var got Command
			if err := json.Unmarshal([]byte(tc.input), &got); err != nil {
				t.Fatalf("unexpected unmarshal error: %v", err)
			}

			if got.Cmd != tc.want.Cmd {
				t.Errorf("Cmd: got %q, want %q", got.Cmd, tc.want.Cmd)
			}
			if got.ID != tc.want.ID {
				t.Errorf("ID: got %q, want %q", got.ID, tc.want.ID)
			}
			if tc.want.Args != nil {
				if string(got.Args) != string(tc.want.Args) {
					t.Errorf("Args: got %s, want %s", got.Args, tc.want.Args)
				}
			} else if got.Args != nil {
				t.Errorf("Args: got %s, want nil", got.Args)
			}
		})
	}
}

func TestResponse_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		resp    Response
		wantID  bool
	}{
		{
			name: "response type with id",
			resp: Response{
				Type:    "response",
				ID:      "resp-1",
				Payload: json.RawMessage(`{"status":"ok"}`),
			},
			wantID: true,
		},
		{
			name: "event type without id",
			resp: Response{
				Type:    "event",
				Payload: json.RawMessage(`{"event":"assistant","text":"hello"}`),
			},
			wantID: false,
		},
		{
			name: "error type with id",
			resp: Response{
				Type:    "error",
				ID:      "err-99",
				Payload: json.RawMessage(`{"message":"something went wrong"}`),
			},
			wantID: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			data, err := json.Marshal(tc.resp)
			if err != nil {
				t.Fatalf("unexpected marshal error: %v", err)
			}

			var got map[string]json.RawMessage
			if err := json.Unmarshal(data, &got); err != nil {
				t.Fatalf("unexpected unmarshal error: %v", err)
			}

			assertJSONString(t, got["type"], tc.resp.Type)

			if tc.wantID {
				assertJSONString(t, got["id"], tc.resp.ID)
			} else {
				if _, ok := got["id"]; ok {
					t.Error("id field should be omitted when empty")
				}
			}

			if string(got["payload"]) != string(tc.resp.Payload) {
				t.Errorf("payload: got %s, want %s", got["payload"], tc.resp.Payload)
			}
		})
	}
}

func TestStreamEvent_AllEventTypes(t *testing.T) {
	tests := []struct {
		name  string
		event StreamEvent
	}{
		{
			name:  "assistant event with text",
			event: StreamEvent{Event: EventAssistant, Text: "Here is the plan."},
		},
		{
			name:  "tool_use event with tool and input",
			event: StreamEvent{Event: EventToolUse, Tool: "Bash", Input: "ls -la"},
		},
		{
			name:  "tool_result event with output",
			event: StreamEvent{Event: EventToolResult, Tool: "Bash", Output: "file1\nfile2"},
		},
		{
			name:  "question event with options",
			event: StreamEvent{Event: EventQuestion, Text: "Which env?", Options: []string{"dev", "prod"}},
		},
		{
			name:  "cost event with totalCost",
			event: StreamEvent{Event: EventCost, Cost: 0.00245},
		},
		{
			name:  "complete event with status",
			event: StreamEvent{Event: EventComplete, Status: "success"},
		},
		{
			name:  "error event with text",
			event: StreamEvent{Event: EventError, Text: "timeout"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			data, err := json.Marshal(tc.event)
			if err != nil {
				t.Fatalf("marshal error: %v", err)
			}

			var got StreamEvent
			if err := json.Unmarshal(data, &got); err != nil {
				t.Fatalf("unmarshal error: %v", err)
			}

			if got.Event != tc.event.Event {
				t.Errorf("Event: got %q, want %q", got.Event, tc.event.Event)
			}
			if got.Text != tc.event.Text {
				t.Errorf("Text: got %q, want %q", got.Text, tc.event.Text)
			}
			if got.Tool != tc.event.Tool {
				t.Errorf("Tool: got %q, want %q", got.Tool, tc.event.Tool)
			}
			if got.Input != tc.event.Input {
				t.Errorf("Input: got %q, want %q", got.Input, tc.event.Input)
			}
			if got.Output != tc.event.Output {
				t.Errorf("Output: got %q, want %q", got.Output, tc.event.Output)
			}
			if !reflect.DeepEqual(got.Options, tc.event.Options) {
				t.Errorf("Options: got %v, want %v", got.Options, tc.event.Options)
			}
			if got.Cost != tc.event.Cost {
				t.Errorf("Cost: got %v, want %v", got.Cost, tc.event.Cost)
			}
			if got.Status != tc.event.Status {
				t.Errorf("Status: got %q, want %q", got.Status, tc.event.Status)
			}
		})
	}

	t.Run("omitempty fields absent when zero", func(t *testing.T) {
		event := StreamEvent{Event: EventComplete}
		data, err := json.Marshal(event)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var raw map[string]json.RawMessage
		if err := json.Unmarshal(data, &raw); err != nil {
			t.Fatalf("unmarshal error: %v", err)
		}

		for _, field := range []string{"text", "tool", "input", "output", "options", "totalCost", "status"} {
			if _, ok := raw[field]; ok {
				t.Errorf("field %q should be omitted when zero-valued", field)
			}
		}
	})
}

func TestInitResponse_RoundTrip(t *testing.T) {
	original := InitResponse{
		Skills: []Skill{
			{Name: "setup", Description: "Initialize engagement", Path: ".claude/skills/setup"},
			{Name: "discover-infrastructure", Description: "Discovery interview", Path: ".claude/skills/discover-infrastructure"},
		},
		Engagements: []Engagement{
			{Slug: "acme-federal", HasContext: true},
			{Slug: "beta-corp", HasContext: false},
		},
		Agents: []Agent{
			{Name: "Architect", Model: "claude-opus-4", Role: "architect", Description: "Team lead"},
			{Name: "Senior Developer", Model: "claude-sonnet-4", Role: "developer", Description: "Implementation"},
		},
		State: State{
			LastEngagement: "acme-federal",
			Preferences:    map[string]string{"theme": "dark", "lang": "en"},
		},
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var got InitResponse
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if !reflect.DeepEqual(original, got) {
		t.Errorf("round-trip mismatch\ngot:  %+v\nwant: %+v", got, original)
	}
}

func TestPhaseInfo_RoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		phase PhaseInfo
	}{
		{
			name: "pre-engagement with zero counts",
			phase: PhaseInfo{
				Phase:          PhasePreEngagement,
				ArtifactCounts: ArtifactCounts{},
			},
		},
		{
			name: "live with discovery and assessments",
			phase: PhaseInfo{
				Phase: PhaseLive,
				ArtifactCounts: ArtifactCounts{
					Discovery:   3,
					Assessments: 2,
				},
			},
		},
		{
			name: "leave-behind with full counts",
			phase: PhaseInfo{
				Phase: PhaseLeaveBehind,
				ArtifactCounts: ArtifactCounts{
					Discovery:    5,
					Assessments:  4,
					Deliverables: 2,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			data, err := json.Marshal(tc.phase)
			if err != nil {
				t.Fatalf("marshal error: %v", err)
			}

			var got PhaseInfo
			if err := json.Unmarshal(data, &got); err != nil {
				t.Fatalf("unmarshal error: %v", err)
			}

			if !reflect.DeepEqual(tc.phase, got) {
				t.Errorf("round-trip mismatch\ngot:  %+v\nwant: %+v", got, tc.phase)
			}
		})
	}
}

func TestChecklist_RoundTrip(t *testing.T) {
	original := Checklist{
		Name:     "OCP Readiness",
		FileName: "ocp-readiness.md",
		Sections: []ChecklistSection{
			{
				Title: "Infrastructure",
				Items: []ChecklistItem{
					{Text: "DNS configured", Checked: true, Line: 5},
					{Text: "Storage class set", Checked: false, Line: 6},
				},
			},
			{
				Title: "Security",
				Items: []ChecklistItem{
					{Text: "TLS certs in place", Checked: true, Line: 12},
					{Text: "RBAC policies defined", Checked: true, Line: 13},
					{Text: "Secrets management configured", Checked: false, Line: 14},
				},
			},
		},
		CompletionPct: 60,
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var got Checklist
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if !reflect.DeepEqual(original, got) {
		t.Errorf("round-trip mismatch\ngot:  %+v\nwant: %+v", got, original)
	}
}

func TestArtifactNode_NestedChildren(t *testing.T) {
	original := ArtifactNode{
		Name: "acme-federal",
		Path: "engagements/acme-federal",
		Type: "directory",
		Children: []ArtifactNode{
			{
				Name: "discovery",
				Path: "engagements/acme-federal/discovery",
				Type: "directory",
				Children: []ArtifactNode{
					{Name: "notes.md", Path: "engagements/acme-federal/discovery/notes.md", Type: "file"},
					{Name: "inventory.csv", Path: "engagements/acme-federal/discovery/inventory.csv", Type: "file"},
				},
			},
			{
				Name: "CONTEXT.md",
				Path: "engagements/acme-federal/CONTEXT.md",
				Type: "file",
			},
		},
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var got ArtifactNode
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if !reflect.DeepEqual(original, got) {
		t.Errorf("round-trip mismatch\ngot:  %+v\nwant: %+v", got, original)
	}

	t.Run("leaf node has no children field", func(t *testing.T) {
		leaf := ArtifactNode{Name: "notes.md", Path: "discovery/notes.md", Type: "file"}
		leafData, err := json.Marshal(leaf)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var raw map[string]json.RawMessage
		if err := json.Unmarshal(leafData, &raw); err != nil {
			t.Fatalf("unmarshal error: %v", err)
		}

		if _, ok := raw["children"]; ok {
			t.Error("children field should be omitted for leaf nodes")
		}
	})
}

// assertJSONString checks that a JSON-encoded value decodes to the expected string.
func assertJSONString(t *testing.T, raw json.RawMessage, want string) {
	t.Helper()
	var got string
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("assertJSONString: unmarshal error: %v", err)
	}
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
