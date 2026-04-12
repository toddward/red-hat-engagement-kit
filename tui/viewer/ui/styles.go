package ui

import "github.com/charmbracelet/lipgloss"

// Red Hat brand colors
var (
	RedHatRed     = lipgloss.Color("#EE0000")
	RedHatRedDark = lipgloss.Color("#A30000")
	RedHatBlack   = lipgloss.Color("#151515")
	Surface       = lipgloss.Color("#1A1A1A")
	SurfaceLight  = lipgloss.Color("#2E2E2E")
	Border        = lipgloss.Color("#3A3A3A")
	TextPrimary   = lipgloss.Color("#E8E8E8")
	TextMuted     = lipgloss.Color("#888888")
	TextDim       = lipgloss.Color("#555555")
	Green         = lipgloss.Color("#3E8635")
	Yellow        = lipgloss.Color("#F0AB00")
	Blue          = lipgloss.Color("#0066CC")
)

var (
	SidebarWidth = 32

	SidebarStyle = lipgloss.NewStyle().
			Width(SidebarWidth).
			Padding(1, 2).
			BorderRight(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(RedHatRed)

	MainStyle = lipgloss.NewStyle().
			Padding(1, 2)

	TitleStyle = lipgloss.NewStyle().
			Foreground(RedHatRed).
			Bold(true).
			MarginBottom(1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(TextMuted).
			Italic(true)
)

var (
	MenuItemStyle = lipgloss.NewStyle().
			Foreground(TextPrimary).
			PaddingLeft(2)

	MenuItemSelectedStyle = lipgloss.NewStyle().
				Foreground(RedHatRed).
				Bold(true).
				PaddingLeft(2)

	MenuHeaderStyle = lipgloss.NewStyle().
			Foreground(TextMuted).
			Bold(true).
			MarginTop(1).
			MarginBottom(1)
)

var (
	ActivityTimestampStyle = lipgloss.NewStyle().
				Foreground(TextDim).
				Width(10)

	ActivityAssistantStyle = lipgloss.NewStyle().
				Foreground(TextPrimary)

	ActivityToolStyle = lipgloss.NewStyle().
				Foreground(Blue)

	ActivityToolResultStyle = lipgloss.NewStyle().
				Foreground(TextMuted)

	ActivityErrorStyle = lipgloss.NewStyle().
				Foreground(RedHatRed).
				Bold(true)
)

var (
	InputStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(RedHatRed).
			Padding(0, 1)

	InputLabelStyle = lipgloss.NewStyle().
			Foreground(TextMuted).
			MarginBottom(1)
)

var (
	PaletteStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(RedHatRed).
			Padding(1, 2).
			Width(60)

	PaletteInputStyle = lipgloss.NewStyle().
				MarginBottom(1)

	PaletteItemStyle = lipgloss.NewStyle().
				Foreground(TextPrimary)

	PaletteItemSelectedStyle = lipgloss.NewStyle().
					Foreground(RedHatRed).
					Bold(true)

	PaletteItemDescStyle = lipgloss.NewStyle().
				Foreground(TextMuted).
				MarginLeft(2)
)

var (
	StatusRunningStyle = lipgloss.NewStyle().
				Foreground(Yellow).
				Bold(true)

	StatusCompleteStyle = lipgloss.NewStyle().
				Foreground(Green).
				Bold(true)

	StatusErrorStyle = lipgloss.NewStyle().
				Foreground(RedHatRed).
				Bold(true)

	PhasePreEngagementStyle = lipgloss.NewStyle().
				Foreground(TextMuted)

	PhaseLiveStyle = lipgloss.NewStyle().
			Foreground(Green).
			Bold(true)

	PhaseLeaveBehindStyle = lipgloss.NewStyle().
				Foreground(Blue)
)

var (
	HelpKeyStyle = lipgloss.NewStyle().
			Foreground(RedHatRed)

	HelpDescStyle = lipgloss.NewStyle().
			Foreground(TextMuted)
)

var (
	ArtifactDirStyle  = lipgloss.NewStyle().Foreground(Blue).Bold(true)
	ArtifactFileStyle = lipgloss.NewStyle().Foreground(TextPrimary)
)

var (
	ChecklistCheckedStyle   = lipgloss.NewStyle().Foreground(Green)
	ChecklistUncheckedStyle = lipgloss.NewStyle().Foreground(TextMuted)
	ChecklistProgressStyle  = lipgloss.NewStyle().Foreground(Yellow)
)
