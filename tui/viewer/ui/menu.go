package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MenuItem struct {
	Key         string
	Label       string
	Description string
	Action      string
	Children    []MenuItem
}

type Menu struct {
	items  []MenuItem
	cursor int
	parent *Menu
	title  string
	width  int
	height int
}

func NewMenu(title string, items []MenuItem) Menu {
	return Menu{items: items, title: title}
}

func (m *Menu) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m *Menu) SetItems(items []MenuItem) {
	m.items = items
	if m.cursor >= len(items) {
		m.cursor = len(items) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
}

func (m Menu) Selected() *MenuItem {
	if m.cursor < 0 || m.cursor >= len(m.items) {
		return nil
	}
	return &m.items[m.cursor]
}

func (m Menu) Update(msg tea.Msg) (Menu, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case "home":
			m.cursor = 0
		case "end":
			m.cursor = len(m.items) - 1
		}
	}
	return m, nil
}

func (m Menu) View() string {
	var b strings.Builder

	if m.title != "" {
		b.WriteString(TitleStyle.Render(m.title))
		b.WriteString("\n")
		maxWidth := m.width - 4
		if maxWidth > 40 {
			maxWidth = 40
		}
		b.WriteString(lipgloss.NewStyle().
			Foreground(RedHatRed).
			Render(strings.Repeat("━", maxWidth)))
		b.WriteString("\n\n")
	}

	for i, item := range m.items {
		cursor := "  "
		style := MenuItemStyle
		if i == m.cursor {
			cursor = "> "
			style = MenuItemSelectedStyle
		}

		line := cursor + item.Label
		if item.Key != "" {
			line = cursor + "[" + item.Key + "] " + item.Label
		}

		b.WriteString(style.Render(line))

		if item.Description != "" && i == m.cursor {
			b.WriteString("\n")
			b.WriteString(lipgloss.NewStyle().
				Foreground(TextMuted).
				PaddingLeft(4).
				Render(item.Description))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(HelpKeyStyle.Render("↑/↓"))
	b.WriteString(HelpDescStyle.Render(" navigate  "))
	b.WriteString(HelpKeyStyle.Render("enter"))
	b.WriteString(HelpDescStyle.Render(" select  "))
	b.WriteString(HelpKeyStyle.Render("/"))
	b.WriteString(HelpDescStyle.Render(" search  "))
	b.WriteString(HelpKeyStyle.Render("q"))
	b.WriteString(HelpDescStyle.Render(" quit"))

	return MainStyle.Render(b.String())
}
