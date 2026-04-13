package ui

import (
	"strings"
	"testing"

	"github.com/toddward/red-hat-engagement-kit/tui/viewer/protocol"
)

func TestNewSidebar(t *testing.T) {
	s := NewSidebar()

	if s.phase != protocol.PhasePreEngagement {
		t.Errorf("expected default phase PhasePreEngagement, got %q", s.phase)
	}
	if s.width != SidebarWidth {
		t.Errorf("expected width %d, got %d", SidebarWidth, s.width)
	}
}

func TestSidebar_SetEngagement(t *testing.T) {
	s := NewSidebar()
	s.SetEngagement("acme-corp")

	view := s.View()
	if !strings.Contains(view, "acme-corp") {
		t.Error("expected View() to contain engagement slug 'acme-corp'")
	}
}

func TestSidebar_NoEngagement(t *testing.T) {
	s := NewSidebar()

	view := s.View()
	if !strings.Contains(view, "None selected") {
		t.Error("expected View() to contain 'None selected' when no engagement set")
	}
	if strings.Contains(view, "ARTIFACTS") {
		t.Error("expected View() to NOT contain 'ARTIFACTS' when no engagement set")
	}
}

func TestSidebar_SetPhase_AllPhases(t *testing.T) {
	tests := []struct {
		phase    protocol.Phase
		expected string
	}{
		{protocol.PhasePreEngagement, "Pre-Engagement"},
		{protocol.PhaseLive, "Live"},
		{protocol.PhaseLeaveBehind, "Leave-Behind"},
	}

	for _, tt := range tests {
		t.Run(string(tt.phase), func(t *testing.T) {
			s := NewSidebar()
			s.SetEngagement("test-eng")
			s.SetPhase(protocol.PhaseInfo{Phase: tt.phase})

			view := s.View()
			if !strings.Contains(view, tt.expected) {
				t.Errorf("expected View() to contain %q for phase %q", tt.expected, tt.phase)
			}
		})
	}
}

func TestSidebar_ArtifactCounts(t *testing.T) {
	s := NewSidebar()
	s.SetEngagement("test-eng")
	s.SetPhase(protocol.PhaseInfo{
		Phase: protocol.PhaseLive,
		ArtifactCounts: protocol.ArtifactCounts{
			Discovery:    3,
			Assessments:  2,
			Deliverables: 1,
		},
	})

	view := s.View()
	if !strings.Contains(view, "ARTIFACTS") {
		t.Error("expected View() to contain 'ARTIFACTS' when engagement is set")
	}
	if !strings.Contains(view, "3") {
		t.Error("expected View() to contain discovery count '3'")
	}
	if !strings.Contains(view, "2") {
		t.Error("expected View() to contain assessments count '2'")
	}
	if !strings.Contains(view, "1") {
		t.Error("expected View() to contain deliverables count '1'")
	}
}

func TestSidebar_SetAgent(t *testing.T) {
	s := NewSidebar()
	s.SetAgent("Architect")

	view := s.View()
	if !strings.Contains(view, "AGENT") {
		t.Error("expected View() to contain 'AGENT'")
	}
	if !strings.Contains(view, "Architect") {
		t.Error("expected View() to contain agent name 'Architect'")
	}
}
