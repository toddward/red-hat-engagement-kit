package ui

import (
	"strings"
	"testing"

	"github.com/toddward/red-hat-engagement-kit/tui/viewer/protocol"
)

func TestNewActivity(t *testing.T) {
	a := NewActivity()

	if len(a.entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(a.entries))
	}
	if a.running {
		t.Error("expected running=false on new activity")
	}
}

func TestActivity_AddEvent_Assistant(t *testing.T) {
	a := NewActivity()
	a.SetSize(80, 24)

	a.AddEvent(protocol.StreamEvent{
		Event: protocol.EventAssistant,
		Text:  "Hello from the assistant",
	})

	if len(a.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(a.entries))
	}
	if a.entries[0].Event != protocol.EventAssistant {
		t.Errorf("expected EventAssistant, got %v", a.entries[0].Event)
	}

	view := a.View()
	if !strings.Contains(view, "Hello from the assistant") {
		t.Error("expected View() to contain assistant text")
	}
}

func TestActivity_AddEvent_ToolUse(t *testing.T) {
	a := NewActivity()
	a.SetSize(80, 24)

	a.AddEvent(protocol.StreamEvent{
		Event: protocol.EventToolUse,
		Tool:  "Read",
		Input: "some input",
	})

	if len(a.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(a.entries))
	}

	view := a.View()
	if !strings.Contains(view, "Read") {
		t.Error("expected View() to contain tool name 'Read'")
	}
}

func TestActivity_AddEvent_ToolResult(t *testing.T) {
	a := NewActivity()
	a.SetSize(80, 24)

	a.AddEvent(protocol.StreamEvent{
		Event:  protocol.EventToolResult,
		Output: "The result output text",
	})

	if len(a.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(a.entries))
	}
	if a.entries[0].Text != "The result output text" {
		t.Errorf("expected entry text from Output field, got %q", a.entries[0].Text)
	}

	// Test truncation at 200 chars: the entry stores the full Output,
	// but updateContent truncates to 200 chars + "..." when rendering into the viewport.
	longOutput := strings.Repeat("x", 300)
	a2 := NewActivity()
	a2.SetSize(400, 24) // wide enough so the line is not wrapped/clipped
	a2.AddEvent(protocol.StreamEvent{
		Event:  protocol.EventToolResult,
		Output: longOutput,
	})

	// The entry text should be the full output (set from Output field)
	if a2.entries[0].Text != longOutput {
		t.Error("expected entry text to contain full output")
	}
	// The viewport content (rendered by updateContent) should have the truncated text
	vpContent := a2.viewport.View()
	if !strings.Contains(vpContent, "...") {
		// The viewport may apply its own width clipping. Check the view as a whole.
		view2 := a2.View()
		if !strings.Contains(view2, "...") {
			t.Errorf("expected truncated tool result to contain '...'; viewport content length=%d", len(vpContent))
		}
	}
}

func TestActivity_AddEvent_Cost(t *testing.T) {
	a := NewActivity()
	a.SetSize(80, 24)

	a.AddEvent(protocol.StreamEvent{
		Event: protocol.EventCost,
		Cost:  0.0123,
	})

	// Cost events should NOT add an entry
	if len(a.entries) != 0 {
		t.Errorf("expected 0 entries after cost event, got %d", len(a.entries))
	}
	if a.cost != 0.0123 {
		t.Errorf("expected cost 0.0123, got %f", a.cost)
	}
}

func TestActivity_AddEvent_Complete(t *testing.T) {
	a := NewActivity()
	a.SetSize(80, 24)
	a.SetRunning(true)

	a.AddEvent(protocol.StreamEvent{
		Event: protocol.EventComplete,
		Text:  "success",
	})

	if a.running {
		t.Error("expected running=false after complete event")
	}
}

func TestActivity_AddEvent_Error(t *testing.T) {
	a := NewActivity()
	a.SetSize(80, 24)

	a.AddEvent(protocol.StreamEvent{
		Event: protocol.EventError,
		Text:  "something went wrong",
	})

	if len(a.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(a.entries))
	}

	view := a.View()
	if !strings.Contains(view, "something went wrong") {
		t.Error("expected View() to contain error text")
	}
}

func TestActivity_Clear(t *testing.T) {
	a := NewActivity()
	a.SetSize(80, 24)
	a.SetRunning(true)
	a.AddEvent(protocol.StreamEvent{
		Event: protocol.EventAssistant,
		Text:  "hello",
	})
	a.cost = 1.5

	a.Clear()

	if len(a.entries) != 0 {
		t.Errorf("expected 0 entries after clear, got %d", len(a.entries))
	}
	if a.running {
		t.Error("expected running=false after clear")
	}
	if a.cost != 0 {
		t.Errorf("expected cost=0 after clear, got %f", a.cost)
	}
}

func TestActivity_SetRunning(t *testing.T) {
	a := NewActivity()
	a.SetSize(80, 24)

	a.SetRunning(true)
	view := a.View()
	if !strings.Contains(view, "Running") {
		t.Error("expected View() to contain 'Running' when running=true")
	}

	a.SetRunning(false)
	view = a.View()
	if !strings.Contains(view, "Execution Log") {
		t.Error("expected View() to contain 'Execution Log' when running=false")
	}
}

func TestActivity_CostDisplay(t *testing.T) {
	a := NewActivity()
	a.SetSize(80, 24)

	// No cost: should not show cost line
	view := a.View()
	if strings.Contains(view, "Cost:") {
		t.Error("expected View() to NOT contain 'Cost:' when cost is 0")
	}

	// With cost
	a.AddEvent(protocol.StreamEvent{
		Event: protocol.EventCost,
		Cost:  0.0456,
	})
	view = a.View()
	if !strings.Contains(view, "Cost:") {
		t.Error("expected View() to contain 'Cost:' when cost > 0")
	}
	if !strings.Contains(view, "0.0456") {
		t.Error("expected View() to contain formatted cost value")
	}
}
