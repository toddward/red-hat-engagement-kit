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

	title := TitleStyle.Render("Red Hat")
	subtitle := SubtitleStyle.Render("Engagement Kit")
	b.WriteString(title + "\n")
	b.WriteString(subtitle + "\n\n")

	b.WriteString(MenuHeaderStyle.Render("ENGAGEMENT"))
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
		b.WriteString(MenuHeaderStyle.Render("AGENT"))
		b.WriteString("\n")
		b.WriteString(StatusRunningStyle.Render("● " + s.agent))
		b.WriteString("\n\n")
	}

	if s.engagement != "" {
		b.WriteString(MenuHeaderStyle.Render("ARTIFACTS"))
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("Discovery:    %d\n", s.counts.Discovery))
		b.WriteString(fmt.Sprintf("Assessments:  %d\n", s.counts.Assessments))
		b.WriteString(fmt.Sprintf("Deliverables: %d\n", s.counts.Deliverables))
	}

	return SidebarStyle.Height(s.height).Render(b.String())
}
