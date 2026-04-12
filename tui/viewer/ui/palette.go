package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type PaletteItem struct {
	Name        string
	Description string
	Action      string
	Category    string
}

type Palette struct {
	input    textinput.Model
	items    []PaletteItem
	filtered []PaletteItem
	cursor   int
	width    int
	height   int
}

func NewPalette() Palette {
	ti := textinput.New()
	ti.Placeholder = "Type to search..."
	ti.CharLimit = 100
	return Palette{input: ti, items: make([]PaletteItem, 0), filtered: make([]PaletteItem, 0)}
}

func (p *Palette) SetItems(items []PaletteItem) {
	p.items = items
	p.filter()
}

func (p *Palette) Open() {
	p.input.SetValue("")
	p.input.Focus()
	p.cursor = 0
	p.filter()
}

func (p *Palette) Close() { p.input.Blur() }

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
		if strings.Contains(name, query) || strings.Contains(desc, query) || strings.Contains(cat, query) {
			p.filtered = append(p.filtered, item)
		}
	}
	if p.cursor >= len(p.filtered) {
		p.cursor = len(p.filtered) - 1
	}
	if p.cursor < 0 {
		p.cursor = 0
	}
}

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

func (p Palette) View() string {
	var b strings.Builder

	b.WriteString(PaletteInputStyle.Render(p.input.View()))
	b.WriteString("\n")

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
