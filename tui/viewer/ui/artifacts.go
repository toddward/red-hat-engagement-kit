package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/toddward/red-hat-engagement-kit/tui/viewer/protocol"
)

type flatNode struct {
	name  string
	path  string
	isDir bool
	depth int
}

// ArtifactBrowser is a file tree browser with content preview.
type ArtifactBrowser struct {
	tree     []protocol.ArtifactNode
	flat     []flatNode
	cursor   int
	viewport viewport.Model
	content  string
	viewing  bool
	width    int
	height   int
}

func NewArtifactBrowser() ArtifactBrowser {
	vp := viewport.New(80, 20)
	return ArtifactBrowser{
		viewport: vp,
	}
}

func (b *ArtifactBrowser) SetSize(width, height int) {
	b.width = width
	b.height = height
	b.viewport.Width = width - 4
	b.viewport.Height = height - 6
}

// SetTree flattens the artifact tree for cursor navigation.
func (b *ArtifactBrowser) SetTree(nodes []protocol.ArtifactNode) {
	b.tree = nodes
	b.flat = nil
	b.cursor = 0
	flattenNodes(nodes, 0, &b.flat)
}

func flattenNodes(nodes []protocol.ArtifactNode, depth int, out *[]flatNode) {
	for _, n := range nodes {
		*out = append(*out, flatNode{
			name:  n.Name,
			path:  n.Path,
			isDir: n.Type == "directory",
			depth: depth,
		})
		if len(n.Children) > 0 {
			flattenNodes(n.Children, depth+1, out)
		}
	}
}

// SetContent switches to content viewing mode.
func (b *ArtifactBrowser) SetContent(content string) {
	b.content = content
	b.viewport.SetContent(content)
	b.viewport.GotoTop()
	b.viewing = true
}

// CloseContent returns to tree browsing mode.
func (b *ArtifactBrowser) CloseContent() {
	b.viewing = false
	b.content = ""
}

// IsViewing reports whether content viewing mode is active.
func (b *ArtifactBrowser) IsViewing() bool {
	return b.viewing
}

// SelectedPath returns the path of the currently highlighted node.
func (b *ArtifactBrowser) SelectedPath() string {
	if len(b.flat) == 0 {
		return ""
	}
	return b.flat[b.cursor].path
}

// SelectedIsFile reports whether the selected node is a file.
func (b *ArtifactBrowser) SelectedIsFile() bool {
	if len(b.flat) == 0 {
		return false
	}
	return !b.flat[b.cursor].isDir
}

func (b ArtifactBrowser) Update(msg tea.Msg) (ArtifactBrowser, tea.Cmd) {
	if b.viewing {
		var cmd tea.Cmd
		b.viewport, cmd = b.viewport.Update(msg)
		return b, cmd
	}

	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "up", "k":
			if b.cursor > 0 {
				b.cursor--
			}
		case "down", "j":
			if b.cursor < len(b.flat)-1 {
				b.cursor++
			}
		}
	}
	return b, nil
}

func (b ArtifactBrowser) View() string {
	var sb strings.Builder

	sb.WriteString(TitleStyle.Render("Artifacts"))
	sb.WriteString("\n")

	maxWidth := b.width - 4
	if maxWidth < 1 {
		maxWidth = 1
	}
	if maxWidth > 60 {
		maxWidth = 60
	}
	sb.WriteString(DividerStyle.Render(strings.Repeat("─", maxWidth)))
	sb.WriteString("\n\n")

	if b.viewing {
		sb.WriteString(b.viewport.View())
		sb.WriteString("\n\n")
		sb.WriteString(HelpKeyStyle.Render("↑/↓"))
		sb.WriteString(HelpDescStyle.Render(" scroll  "))
		sb.WriteString(HelpKeyStyle.Render("Esc"))
		sb.WriteString(HelpDescStyle.Render(" back to tree"))
	} else {
		if len(b.flat) == 0 {
			sb.WriteString(lipgloss.NewStyle().Foreground(TextMuted).Render("No artifacts found."))
		} else {
			for i, node := range b.flat {
				indent := strings.Repeat("  ", node.depth)
				var label string
				if node.isDir {
					label = ArtifactDirStyle.Render("📁 " + node.name)
				} else {
					label = ArtifactFileStyle.Render(node.name)
				}

				var line string
				if i == b.cursor {
					cursor := lipgloss.NewStyle().Foreground(RedHatRed).Bold(true).Render("❯ ")
					line = indent + cursor + label
				} else {
					line = indent + "  " + label
				}
				sb.WriteString(line + "\n")
			}
		}
		sb.WriteString("\n")
		sb.WriteString(HelpKeyStyle.Render("↑/↓"))
		sb.WriteString(HelpDescStyle.Render(" navigate  "))
		sb.WriteString(HelpKeyStyle.Render("Enter"))
		sb.WriteString(HelpDescStyle.Render(" open file  "))
		sb.WriteString(HelpKeyStyle.Render("Esc"))
		sb.WriteString(HelpDescStyle.Render(" back"))
	}

	return MainStyle.Render(sb.String())
}
