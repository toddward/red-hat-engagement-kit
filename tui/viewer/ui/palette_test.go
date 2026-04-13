package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func testPaletteItems() []PaletteItem {
	return []PaletteItem{
		{Name: "Setup", Description: "Init engagement", Action: "setup", Category: "Skills"},
		{Name: "Discover", Description: "Run discovery", Action: "discover", Category: "Skills"},
		{Name: "Assess", Description: "Portfolio assessment", Action: "assess", Category: "Skills"},
		{Name: "Build Deck", Description: "Generate presentation", Action: "deck", Category: "Deliverables"},
	}
}

func TestNewPalette(t *testing.T) {
	p := NewPalette()

	if len(p.items) != 0 {
		t.Errorf("expected 0 items, got %d", len(p.items))
	}
	if len(p.filtered) != 0 {
		t.Errorf("expected 0 filtered items, got %d", len(p.filtered))
	}
}

func TestPalette_SetItems(t *testing.T) {
	p := NewPalette()
	p.SetItems(testPaletteItems())

	if len(p.items) != 4 {
		t.Errorf("expected 4 items, got %d", len(p.items))
	}
	// With no filter, all items should be in filtered
	if len(p.filtered) != 4 {
		t.Errorf("expected 4 filtered items, got %d", len(p.filtered))
	}
}

func TestPalette_Open(t *testing.T) {
	p := NewPalette()
	p.SetItems(testPaletteItems())

	p.Open()

	// After Open, cursor should be 0 and filter cleared
	if p.cursor != 0 {
		t.Errorf("expected cursor 0, got %d", p.cursor)
	}
	if p.input.Value() != "" {
		t.Errorf("expected empty input value, got %q", p.input.Value())
	}
	if !p.input.Focused() {
		t.Error("expected input to be focused after Open()")
	}
}

func TestPalette_Close(t *testing.T) {
	p := NewPalette()
	p.SetItems(testPaletteItems())
	p.Open()

	p.Close()

	if p.input.Focused() {
		t.Error("expected input to be blurred after Close()")
	}
}

func TestPalette_Selected(t *testing.T) {
	p := NewPalette()
	p.SetItems(testPaletteItems())
	p.Open()

	sel := p.Selected()
	if sel == nil {
		t.Fatal("expected non-nil Selected()")
	}
	if sel.Name != "Setup" {
		t.Errorf("expected 'Setup', got %q", sel.Name)
	}
}

func TestPalette_Selected_Empty(t *testing.T) {
	p := NewPalette()

	sel := p.Selected()
	if sel != nil {
		t.Errorf("expected nil Selected() for empty palette, got %v", sel)
	}
}

func TestPalette_NavigateDownUp(t *testing.T) {
	p := NewPalette()
	p.SetItems(testPaletteItems())
	p.Open()

	// Down with ctrl+n
	p, _ = p.Update(tea.KeyMsg{Type: tea.KeyCtrlN})
	if p.cursor != 1 {
		t.Errorf("expected cursor 1 after ctrl+n, got %d", p.cursor)
	}

	// Down with down key
	p, _ = p.Update(tea.KeyMsg{Type: tea.KeyDown})
	if p.cursor != 2 {
		t.Errorf("expected cursor 2 after down, got %d", p.cursor)
	}

	// Up with ctrl+p
	p, _ = p.Update(tea.KeyMsg{Type: tea.KeyCtrlP})
	if p.cursor != 1 {
		t.Errorf("expected cursor 1 after ctrl+p, got %d", p.cursor)
	}

	// Up with up key
	p, _ = p.Update(tea.KeyMsg{Type: tea.KeyUp})
	if p.cursor != 0 {
		t.Errorf("expected cursor 0 after up, got %d", p.cursor)
	}
}

func TestPalette_NavigateBounds(t *testing.T) {
	p := NewPalette()
	p.SetItems(testPaletteItems())
	p.Open()

	// Up at top stays at 0
	p, _ = p.Update(tea.KeyMsg{Type: tea.KeyUp})
	if p.cursor != 0 {
		t.Errorf("expected cursor to stay at 0, got %d", p.cursor)
	}

	// Move to bottom
	p.cursor = 3
	// Down at bottom stays
	p, _ = p.Update(tea.KeyMsg{Type: tea.KeyDown})
	if p.cursor != 3 {
		t.Errorf("expected cursor to stay at 3, got %d", p.cursor)
	}
}

func TestPalette_Filter_CursorClamp(t *testing.T) {
	p := NewPalette()
	p.SetItems(testPaletteItems())
	p.Open()

	// Move cursor to the end
	p.cursor = 3

	// Type a filter character that reduces the list.
	// "Deck" should match only "Build Deck"
	for _, r := range "Deck" {
		p, _ = p.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}

	if p.cursor >= len(p.filtered) {
		t.Errorf("expected cursor clamped within filtered range, cursor=%d, filtered=%d", p.cursor, len(p.filtered))
	}
	if len(p.filtered) != 1 {
		t.Errorf("expected 1 filtered item for 'Deck', got %d", len(p.filtered))
	}
}

func TestPalette_View_MaxItems(t *testing.T) {
	// Create more than 10 items
	items := make([]PaletteItem, 15)
	for i := range items {
		items[i] = PaletteItem{
			Name:        strings.Repeat("Item", 1) + string(rune('A'+i)),
			Description: "desc",
			Category:    "cat",
		}
	}

	p := NewPalette()
	p.SetItems(items)
	p.Open()

	view := p.View()
	// Items beyond 10 should not appear. Item index 10 is 'K' (A=0, ..., K=10)
	lastVisibleName := "Item" + string(rune('A'+9))  // ItemJ
	firstHiddenName := "Item" + string(rune('A'+10)) // ItemK

	if !strings.Contains(view, lastVisibleName) {
		t.Errorf("expected View() to contain %q (10th item)", lastVisibleName)
	}
	if strings.Contains(view, firstHiddenName) {
		t.Errorf("expected View() to NOT contain %q (11th item)", firstHiddenName)
	}
}

func TestPalette_View_NoMatches(t *testing.T) {
	p := NewPalette()
	p.SetItems(testPaletteItems())
	p.Open()

	// Type something that matches nothing
	for _, r := range "zzzzz" {
		p, _ = p.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}

	view := p.View()
	if !strings.Contains(view, "No matches") {
		t.Error("expected View() to contain 'No matches'")
	}
}
