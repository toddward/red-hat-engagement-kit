package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func HelpView() string {
	var b strings.Builder

	b.WriteString(TitleStyle.Render("Keyboard Shortcuts"))
	b.WriteString("\n\n")

	shortcuts := []struct {
		key  string
		desc string
	}{
		{"↑/↓ or j/k", "Navigate items"},
		{"Enter", "Select / confirm"},
		{"Esc", "Go back / close"},
		{"1-5", "Jump to menu item"},
		{"/", "Open command palette"},
		{"?", "Toggle this help"},
		{"Space", "Toggle checklist item"},
		{"q", "Quit (from main menu)"},
		{"Ctrl+C", "Cancel running execution"},
	}

	for _, s := range shortcuts {
		key := HelpKeyStyle.Width(16).Render(s.key)
		desc := HelpDescStyle.Render(s.desc)
		b.WriteString("  " + key + desc + "\n")
	}

	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(TextDim).Italic(true).
		Render("Press ? or Esc to close"))

	return PaletteStyle.Width(50).Render(b.String())
}
