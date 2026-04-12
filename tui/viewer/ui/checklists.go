package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/toddward/red-hat-engagement-kit/tui/viewer/protocol"
)

type checklistFlatItem struct {
	// isHeader is true for section title rows.
	isHeader bool
	// For header rows:
	title string
	// For item rows:
	item *protocol.ChecklistItem
}

// ChecklistBrowser is a two-mode view: list of checklists and detail view.
type ChecklistBrowser struct {
	names      []string
	listCursor int
	checklist  *protocol.Checklist
	items      []checklistFlatItem
	cursor     int
	viewing    bool
	width      int
	height     int
}

func NewChecklistBrowser() ChecklistBrowser {
	return ChecklistBrowser{}
}

func (c *ChecklistBrowser) SetSize(width, height int) {
	c.width = width
	c.height = height
}

// SetNames populates the list mode.
func (c *ChecklistBrowser) SetNames(names []string) {
	c.names = names
	c.listCursor = 0
}

// SetChecklist flattens sections into items and switches to detail mode.
func (c *ChecklistBrowser) SetChecklist(cl protocol.Checklist) {
	c.checklist = &cl
	c.items = nil
	c.cursor = 0
	for i := range cl.Sections {
		sec := &cl.Sections[i]
		c.items = append(c.items, checklistFlatItem{isHeader: true, title: sec.Title})
		for j := range sec.Items {
			item := &cl.Sections[i].Items[j]
			c.items = append(c.items, checklistFlatItem{isHeader: false, item: item})
		}
	}
	c.viewing = true
}

// CloseDetail returns to list mode.
func (c *ChecklistBrowser) CloseDetail() {
	c.viewing = false
	c.checklist = nil
	c.items = nil
}

// IsViewing reports whether detail mode is active.
func (c *ChecklistBrowser) IsViewing() bool {
	return c.viewing
}

// SelectedName returns the name of the currently highlighted checklist.
func (c *ChecklistBrowser) SelectedName() string {
	if len(c.names) == 0 {
		return ""
	}
	return c.names[c.listCursor]
}

// SelectedItem returns the currently selected checklist item (skips headers).
func (c *ChecklistBrowser) SelectedItem() *protocol.ChecklistItem {
	if c.cursor >= len(c.items) {
		return nil
	}
	fi := c.items[c.cursor]
	if fi.isHeader {
		return nil
	}
	return fi.item
}

// ChecklistName returns the name of the currently viewed checklist.
func (c *ChecklistBrowser) ChecklistName() string {
	if c.checklist == nil {
		return ""
	}
	return c.checklist.Name
}

func (c ChecklistBrowser) Update(msg tea.Msg) (ChecklistBrowser, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if c.viewing {
			switch keyMsg.String() {
			case "up", "k":
				c.moveCursorUp()
			case "down", "j":
				c.moveCursorDown()
			}
		} else {
			switch keyMsg.String() {
			case "up", "k":
				if c.listCursor > 0 {
					c.listCursor--
				}
			case "down", "j":
				if c.listCursor < len(c.names)-1 {
					c.listCursor++
				}
			}
		}
	}
	return c, nil
}

// moveCursorUp moves cursor up, skipping header rows.
func (c *ChecklistBrowser) moveCursorUp() {
	for i := c.cursor - 1; i >= 0; i-- {
		if !c.items[i].isHeader {
			c.cursor = i
			return
		}
	}
}

// moveCursorDown moves cursor down, skipping header rows.
func (c *ChecklistBrowser) moveCursorDown() {
	for i := c.cursor + 1; i < len(c.items); i++ {
		if !c.items[i].isHeader {
			c.cursor = i
			return
		}
	}
}

func (c ChecklistBrowser) View() string {
	var sb strings.Builder

	sb.WriteString(TitleStyle.Render("Checklists"))
	sb.WriteString("\n")

	maxWidth := c.width - 4
	if maxWidth < 1 {
		maxWidth = 1
	}
	if maxWidth > 60 {
		maxWidth = 60
	}
	sb.WriteString(DividerStyle.Render(strings.Repeat("─", maxWidth)))
	sb.WriteString("\n\n")

	if c.viewing {
		c.renderDetail(&sb)
	} else {
		c.renderList(&sb)
	}

	return MainStyle.Render(sb.String())
}

func (c *ChecklistBrowser) renderList(sb *strings.Builder) {
	if len(c.names) == 0 {
		sb.WriteString(lipgloss.NewStyle().Foreground(TextMuted).Render("No checklists found."))
		sb.WriteString("\n")
	} else {
		for i, name := range c.names {
			if i == c.listCursor {
				cursor := lipgloss.NewStyle().Foreground(RedHatRed).Bold(true).Render("❯ ")
				sb.WriteString(cursor + lipgloss.NewStyle().Foreground(TextPrimary).Bold(true).Render(name))
			} else {
				sb.WriteString("  " + lipgloss.NewStyle().Foreground(TextPrimary).Render(name))
			}
			sb.WriteString("\n")
		}
	}
	sb.WriteString("\n")
	sb.WriteString(HelpKeyStyle.Render("↑/↓"))
	sb.WriteString(HelpDescStyle.Render(" navigate  "))
	sb.WriteString(HelpKeyStyle.Render("Enter"))
	sb.WriteString(HelpDescStyle.Render(" open  "))
	sb.WriteString(HelpKeyStyle.Render("Esc"))
	sb.WriteString(HelpDescStyle.Render(" back"))
}

func (c *ChecklistBrowser) renderDetail(sb *strings.Builder) {
	if c.checklist == nil {
		return
	}

	pct := c.checklist.CompletionPct
	progress := ChecklistProgressStyle.Render(fmt.Sprintf("Progress: %d%%", pct))
	sb.WriteString(progress)
	sb.WriteString("\n\n")

	for i, fi := range c.items {
		if fi.isHeader {
			sb.WriteString(lipgloss.NewStyle().Foreground(TextMuted).Bold(true).Render(fi.title))
			sb.WriteString("\n")
			continue
		}

		item := fi.item
		checkbox := "[ ]"
		if item.Checked {
			checkbox = "[x]"
		}

		var itemStyle lipgloss.Style
		if item.Checked {
			itemStyle = ChecklistCheckedStyle
		} else {
			itemStyle = ChecklistUncheckedStyle
		}

		text := itemStyle.Render(checkbox + " " + item.Text)

		if i == c.cursor {
			cursor := lipgloss.NewStyle().Foreground(RedHatRed).Bold(true).Render("❯ ")
			sb.WriteString("  " + cursor + text)
		} else {
			sb.WriteString("    " + text)
		}
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	sb.WriteString(HelpKeyStyle.Render("↑/↓"))
	sb.WriteString(HelpDescStyle.Render(" navigate  "))
	sb.WriteString(HelpKeyStyle.Render("Space"))
	sb.WriteString(HelpDescStyle.Render(" toggle  "))
	sb.WriteString(HelpKeyStyle.Render("Esc"))
	sb.WriteString(HelpDescStyle.Render(" back"))
}
