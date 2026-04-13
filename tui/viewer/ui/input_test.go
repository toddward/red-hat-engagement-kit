package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewInput(t *testing.T) {
	i := NewInput()

	if i.prompt != "" {
		t.Errorf("expected empty prompt, got %q", i.prompt)
	}
	if len(i.options) != 0 {
		t.Errorf("expected 0 options, got %d", len(i.options))
	}
	if i.cursor != 0 {
		t.Errorf("expected cursor 0, got %d", i.cursor)
	}
}

func TestInput_SetPrompt_FreeText(t *testing.T) {
	i := NewInput()
	i.SetPrompt("Enter your name:", nil)

	if i.prompt != "Enter your name:" {
		t.Errorf("expected prompt 'Enter your name:', got %q", i.prompt)
	}
	if len(i.options) != 0 {
		t.Errorf("expected 0 options, got %d", len(i.options))
	}
	// textinput should be focused (we can check via View containing the input)
	view := i.View()
	if !strings.Contains(view, "Enter your name:") {
		t.Error("expected View() to contain prompt text")
	}
}

func TestInput_SetPrompt_WithOptions(t *testing.T) {
	i := NewInput()
	i.SetPrompt("Choose one:", []string{"Alpha", "Beta", "Gamma"})

	if len(i.options) != 3 {
		t.Fatalf("expected 3 options, got %d", len(i.options))
	}

	view := i.View()
	if !strings.Contains(view, "Alpha") {
		t.Error("expected View() to contain 'Alpha'")
	}
	if !strings.Contains(view, "Beta") {
		t.Error("expected View() to contain 'Beta'")
	}
	if !strings.Contains(view, "Gamma") {
		t.Error("expected View() to contain 'Gamma'")
	}
}

func TestInput_Value_Options(t *testing.T) {
	i := NewInput()
	i.SetPrompt("Choose:", []string{"Alpha", "Beta", "Gamma"})

	// Default cursor=0, should return first option
	if i.Value() != "Alpha" {
		t.Errorf("expected Value()='Alpha', got %q", i.Value())
	}

	// Move cursor down
	i, _ = i.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	if i.Value() != "Beta" {
		t.Errorf("expected Value()='Beta', got %q", i.Value())
	}
}

func TestInput_OptionNavigation(t *testing.T) {
	i := NewInput()
	i.SetPrompt("Choose:", []string{"A", "B", "C"})

	// Down with j
	i, _ = i.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	if i.cursor != 1 {
		t.Errorf("expected cursor 1, got %d", i.cursor)
	}

	// Down with arrow
	i, _ = i.Update(tea.KeyMsg{Type: tea.KeyDown})
	if i.cursor != 2 {
		t.Errorf("expected cursor 2, got %d", i.cursor)
	}

	// Bounds: down at bottom stays
	i, _ = i.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	if i.cursor != 2 {
		t.Errorf("expected cursor to stay at 2, got %d", i.cursor)
	}

	// Up with k
	i, _ = i.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	if i.cursor != 1 {
		t.Errorf("expected cursor 1, got %d", i.cursor)
	}

	// Up with arrow
	i, _ = i.Update(tea.KeyMsg{Type: tea.KeyUp})
	if i.cursor != 0 {
		t.Errorf("expected cursor 0, got %d", i.cursor)
	}

	// Bounds: up at top stays
	i, _ = i.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	if i.cursor != 0 {
		t.Errorf("expected cursor to stay at 0, got %d", i.cursor)
	}
}

func TestInput_SetPrompt_ResetsState(t *testing.T) {
	i := NewInput()
	i.SetPrompt("First:", []string{"A", "B", "C"})

	// Move cursor
	i, _ = i.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	i, _ = i.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	if i.cursor != 2 {
		t.Fatalf("expected cursor 2, got %d", i.cursor)
	}

	// SetPrompt again should reset
	i.SetPrompt("Second:", []string{"X", "Y"})
	if i.cursor != 0 {
		t.Errorf("expected cursor reset to 0, got %d", i.cursor)
	}
	if i.prompt != "Second:" {
		t.Errorf("expected prompt 'Second:', got %q", i.prompt)
	}
}

func TestInput_View_SelectedOption(t *testing.T) {
	i := NewInput()
	i.SetPrompt("Pick:", []string{"One", "Two", "Three"})

	view := i.View()
	// The selected item should have the cursor prefix
	if !strings.Contains(view, "❯") {
		t.Error("expected View() to contain cursor prefix '❯' for selected option")
	}
}
