package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func testMenuItems() []MenuItem {
	return []MenuItem{
		{Key: "1", Label: "Setup", Description: "Initialize engagement", Action: "setup"},
		{Key: "2", Label: "Discover", Description: "Run discovery", Action: "discover"},
		{Key: "3", Label: "Assess", Description: "Assess portfolio", Action: "assess"},
	}
}

func keyMsg(r rune) tea.KeyMsg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}}
}

func specialKeyMsg(t tea.KeyType) tea.KeyMsg {
	return tea.KeyMsg{Type: t}
}

func TestNewMenu(t *testing.T) {
	items := testMenuItems()
	m := NewMenu("Main Menu", items)

	if m.title != "Main Menu" {
		t.Errorf("expected title 'Main Menu', got %q", m.title)
	}
	if m.cursor != 0 {
		t.Errorf("expected cursor 0, got %d", m.cursor)
	}
	if len(m.items) != 3 {
		t.Errorf("expected 3 items, got %d", len(m.items))
	}
}

func TestMenu_NavigateDown(t *testing.T) {
	m := NewMenu("Test", testMenuItems())

	m, _ = m.Update(keyMsg('j'))
	if m.cursor != 1 {
		t.Errorf("expected cursor 1 after j, got %d", m.cursor)
	}

	m, _ = m.Update(specialKeyMsg(tea.KeyDown))
	if m.cursor != 2 {
		t.Errorf("expected cursor 2 after down, got %d", m.cursor)
	}
}

func TestMenu_NavigateUp(t *testing.T) {
	m := NewMenu("Test", testMenuItems())
	m.cursor = 2

	m, _ = m.Update(keyMsg('k'))
	if m.cursor != 1 {
		t.Errorf("expected cursor 1 after k, got %d", m.cursor)
	}

	m, _ = m.Update(specialKeyMsg(tea.KeyUp))
	if m.cursor != 0 {
		t.Errorf("expected cursor 0 after up, got %d", m.cursor)
	}
}

func TestMenu_NavigateBounds(t *testing.T) {
	m := NewMenu("Test", testMenuItems())

	// At top, pressing up should stay at 0
	m, _ = m.Update(keyMsg('k'))
	if m.cursor != 0 {
		t.Errorf("expected cursor to stay at 0, got %d", m.cursor)
	}

	// At bottom, pressing down should stay at len-1
	m.cursor = 2
	m, _ = m.Update(keyMsg('j'))
	if m.cursor != 2 {
		t.Errorf("expected cursor to stay at 2, got %d", m.cursor)
	}
}

func TestMenu_HomeEnd(t *testing.T) {
	m := NewMenu("Test", testMenuItems())
	m.cursor = 1

	m, _ = m.Update(specialKeyMsg(tea.KeyEnd))
	if m.cursor != 2 {
		t.Errorf("expected cursor 2 after end, got %d", m.cursor)
	}

	m, _ = m.Update(specialKeyMsg(tea.KeyHome))
	if m.cursor != 0 {
		t.Errorf("expected cursor 0 after home, got %d", m.cursor)
	}
}

func TestMenu_Selected(t *testing.T) {
	m := NewMenu("Test", testMenuItems())
	m.cursor = 1

	sel := m.Selected()
	if sel == nil {
		t.Fatal("expected non-nil Selected()")
	}
	if sel.Label != "Discover" {
		t.Errorf("expected 'Discover', got %q", sel.Label)
	}
}

func TestMenu_Selected_EmptyMenu(t *testing.T) {
	m := NewMenu("Test", []MenuItem{})

	sel := m.Selected()
	if sel != nil {
		t.Errorf("expected nil Selected() for empty menu, got %v", sel)
	}
}

func TestMenu_SetItems(t *testing.T) {
	m := NewMenu("Test", testMenuItems())
	m.cursor = 2

	// Shrink the list: cursor should clamp
	m.SetItems([]MenuItem{
		{Key: "1", Label: "Only Item"},
	})
	if m.cursor != 0 {
		t.Errorf("expected cursor clamped to 0, got %d", m.cursor)
	}
	if len(m.items) != 1 {
		t.Errorf("expected 1 item, got %d", len(m.items))
	}
}

func TestMenu_PushPopMenu(t *testing.T) {
	m := NewMenu("Root", testMenuItems())
	m.cursor = 1

	subItems := []MenuItem{{Key: "a", Label: "Sub Item"}}
	m.PushMenu("Sub Menu", subItems)

	if m.title != "Sub Menu" {
		t.Errorf("expected title 'Sub Menu', got %q", m.title)
	}
	if m.cursor != 0 {
		t.Errorf("expected cursor 0 after push, got %d", m.cursor)
	}
	if len(m.items) != 1 {
		t.Errorf("expected 1 item after push, got %d", len(m.items))
	}
	if m.AtRoot() {
		t.Error("expected not at root after push")
	}

	ok := m.PopMenu()
	if !ok {
		t.Error("expected PopMenu to return true")
	}
	if m.title != "Root" {
		t.Errorf("expected title 'Root' after pop, got %q", m.title)
	}
	if m.cursor != 1 {
		t.Errorf("expected cursor 1 after pop, got %d", m.cursor)
	}
	if len(m.items) != 3 {
		t.Errorf("expected 3 items after pop, got %d", len(m.items))
	}
	if !m.AtRoot() {
		t.Error("expected at root after pop")
	}

	// Pop on empty history returns false
	ok = m.PopMenu()
	if ok {
		t.Error("expected PopMenu to return false on empty history")
	}
}

func TestMenu_View_ContainsTitle(t *testing.T) {
	m := NewMenu("Main Menu", testMenuItems())
	m.SetSize(80, 24)

	view := m.View()
	if !strings.Contains(view, "Main Menu") {
		t.Error("expected View() to contain title 'Main Menu'")
	}
}

func TestMenu_View_SelectedItemHighlighted(t *testing.T) {
	m := NewMenu("Test", testMenuItems())
	m.SetSize(80, 24)

	view := m.View()
	// The selected item (cursor=0, "Setup") should have the cursor prefix
	if !strings.Contains(view, "❯") {
		t.Error("expected View() to contain cursor prefix '❯'")
	}
	// The selected item label should appear
	if !strings.Contains(view, "Setup") {
		t.Error("expected View() to contain 'Setup'")
	}
}
