package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/toddward/red-hat-engagement-kit/tui/viewer/protocol"
)

type ActivityEntry struct {
	Timestamp time.Time
	Event     protocol.EventType
	Text      string
	Tool      string
	Expanded  bool
}

type Activity struct {
	viewport viewport.Model
	entries  []ActivityEntry
	width    int
	height   int
	running  bool
	cost     float64
}

func NewActivity() Activity {
	vp := viewport.New(80, 20)
	return Activity{
		viewport: vp,
		entries:  make([]ActivityEntry, 0),
	}
}

func (a *Activity) SetSize(width, height int) {
	a.width = width
	a.height = height
	a.viewport.Width = width
	a.viewport.Height = height - 4
}

func (a *Activity) Clear() {
	a.entries = make([]ActivityEntry, 0)
	a.running = false
	a.cost = 0
	a.updateContent()
}

func (a *Activity) SetRunning(running bool) {
	a.running = running
}

func (a *Activity) AddEvent(event protocol.StreamEvent) {
	entry := ActivityEntry{
		Timestamp: time.Now(),
		Event:     event.Event,
		Text:      event.Text,
		Tool:      event.Tool,
	}

	switch event.Event {
	case protocol.EventCost:
		a.cost = event.Cost
		return
	case protocol.EventComplete:
		a.running = false
		a.cost = event.Cost
	case protocol.EventToolResult:
		entry.Text = event.Output
	}

	a.entries = append(a.entries, entry)
	a.updateContent()
	a.viewport.GotoBottom()
}

func (a *Activity) updateContent() {
	var b strings.Builder

	for _, entry := range a.entries {
		ts := ActivityTimestampStyle.Render(entry.Timestamp.Format("15:04:05"))

		var content string
		switch entry.Event {
		case protocol.EventAssistant:
			content = ActivityAssistantStyle.Render(entry.Text)
		case protocol.EventToolUse:
			content = ActivityToolStyle.Render(fmt.Sprintf("▶ %s", entry.Tool))
		case protocol.EventToolResult:
			text := entry.Text
			if len(text) > 200 {
				text = text[:200] + "..."
			}
			content = ActivityToolResultStyle.Render("  └─ " + strings.ReplaceAll(text, "\n", " "))
		case protocol.EventError:
			content = ActivityErrorStyle.Render("✗ " + entry.Text)
		case protocol.EventComplete:
			if entry.Text == "success" || entry.Text == "" {
				content = StatusCompleteStyle.Render("✓ Complete")
			} else {
				content = StatusErrorStyle.Render("✗ " + entry.Text)
			}
		default:
			content = entry.Text
		}

		b.WriteString(ts + " " + content + "\n")
	}

	a.viewport.SetContent(b.String())
}

func (a Activity) Update(msg tea.Msg) (Activity, tea.Cmd) {
	var cmd tea.Cmd
	a.viewport, cmd = a.viewport.Update(msg)
	return a, cmd
}

func (a Activity) View() string {
	var b strings.Builder

	title := "Execution Log"
	if a.running {
		title = StatusRunningStyle.Render("● Running")
	}
	b.WriteString(TitleStyle.Render(title))
	b.WriteString("\n")
	maxWidth := a.width - 4
	if maxWidth > 60 {
		maxWidth = 60
	}
	if maxWidth < 1 {
		maxWidth = 1
	}
	b.WriteString(lipgloss.NewStyle().
		Foreground(RedHatRed).
		Render(strings.Repeat("━", maxWidth)))
	b.WriteString("\n\n")

	b.WriteString(a.viewport.View())
	b.WriteString("\n")

	if a.cost > 0 {
		costStr := fmt.Sprintf("Cost: $%.4f", a.cost)
		b.WriteString(lipgloss.NewStyle().Foreground(TextMuted).Render(costStr))
	}

	b.WriteString("\n")
	b.WriteString(HelpKeyStyle.Render("↑/↓"))
	b.WriteString(HelpDescStyle.Render(" scroll  "))
	b.WriteString(HelpKeyStyle.Render("Ctrl+C"))
	b.WriteString(HelpDescStyle.Render(" cancel  "))
	b.WriteString(HelpKeyStyle.Render("Esc"))
	b.WriteString(HelpDescStyle.Render(" back"))

	return MainStyle.Render(b.String())
}
