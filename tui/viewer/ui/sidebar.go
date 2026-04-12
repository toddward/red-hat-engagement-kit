package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/toddward/red-hat-engagement-kit/tui/viewer/protocol"
)

type Sidebar struct {
	width      int
	height     int
	engagement string
	phase      protocol.Phase
	agent      string
	counts     protocol.ArtifactCounts
}

func NewSidebar() Sidebar {
	return Sidebar{
		width: SidebarWidth,
		phase: protocol.PhasePreEngagement,
	}
}

func (s *Sidebar) SetSize(width, height int) {
	s.width = width
	s.height = height
}

func (s *Sidebar) SetEngagement(slug string) {
	s.engagement = slug
}

func (s *Sidebar) SetPhase(info protocol.PhaseInfo) {
	s.phase = info.Phase
	s.counts = info.ArtifactCounts
}

func (s *Sidebar) SetAgent(name string) {
	s.agent = name
}

func (s Sidebar) View() string {
	var b strings.Builder

	// Brand header
	b.WriteString(lipgloss.NewStyle().Foreground(RedHatRed).Bold(true).Render("Red Hat"))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(TextDim).Render("Engagement Kit"))
	b.WriteString("\n")

	// Thin divider
	divWidth := s.width - 4
	if divWidth < 1 {
		divWidth = 1
	}
	b.WriteString(DividerStyle.Render(strings.Repeat("─", divWidth)))
	b.WriteString("\n\n")

	// Engagement section
	b.WriteString(lipgloss.NewStyle().Foreground(TextDim).Render("ENGAGEMENT"))
	b.WriteString("\n")
	if s.engagement != "" {
		b.WriteString(lipgloss.NewStyle().Foreground(TextPrimary).Render(s.engagement))
		b.WriteString("\n")

		phaseStyle := PhasePreEngagementStyle
		phaseText := "Pre-Engagement"
		switch s.phase {
		case protocol.PhaseLive:
			phaseStyle = PhaseLiveStyle
			phaseText = "● Live"
		case protocol.PhaseLeaveBehind:
			phaseStyle = PhaseLeaveBehindStyle
			phaseText = "Leave-Behind"
		}
		b.WriteString(phaseStyle.Render(phaseText))
		b.WriteString("\n")
	} else {
		b.WriteString(lipgloss.NewStyle().Foreground(TextDim).Italic(true).Render("None selected"))
		b.WriteString("\n")
	}
	b.WriteString("\n")

	if s.agent != "" {
		b.WriteString(DividerStyle.Render(strings.Repeat("─", divWidth)))
		b.WriteString("\n\n")
		b.WriteString(lipgloss.NewStyle().Foreground(TextDim).Render("AGENT"))
		b.WriteString("\n")
		b.WriteString(StatusRunningStyle.Render("● " + s.agent))
		b.WriteString("\n\n")
	}

	if s.engagement != "" {
		b.WriteString(DividerStyle.Render(strings.Repeat("─", divWidth)))
		b.WriteString("\n\n")
		b.WriteString(lipgloss.NewStyle().Foreground(TextDim).Render("ARTIFACTS"))
		b.WriteString("\n")
		b.WriteString(lipgloss.NewStyle().Foreground(TextMuted).Render("Discovery    ") +
			lipgloss.NewStyle().Foreground(TextPrimary).Render(fmt.Sprintf("%d", s.counts.Discovery)))
		b.WriteString("\n")
		b.WriteString(lipgloss.NewStyle().Foreground(TextMuted).Render("Assessments  ") +
			lipgloss.NewStyle().Foreground(TextPrimary).Render(fmt.Sprintf("%d", s.counts.Assessments)))
		b.WriteString("\n")
		b.WriteString(lipgloss.NewStyle().Foreground(TextMuted).Render("Deliverables ") +
			lipgloss.NewStyle().Foreground(TextPrimary).Render(fmt.Sprintf("%d", s.counts.Deliverables)))
		b.WriteString("\n")
	}

	return SidebarStyle.Height(s.height).Render(b.String())
}
