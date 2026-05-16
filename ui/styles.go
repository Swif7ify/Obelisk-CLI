package ui

import "github.com/charmbracelet/lipgloss"

// Color palette
var (
	ColorPrimary   = lipgloss.Color("#8B5CF6") // Vibrant Purple
	ColorSecondary = lipgloss.Color("#06B6D4") // Cyan
	ColorSuccess   = lipgloss.Color("#10B981") // Green
	ColorWarning   = lipgloss.Color("#F59E0B") // Amber
	ColorError     = lipgloss.Color("#EF4444") // Red
	ColorCritical  = lipgloss.Color("#DC2626") // Dark Red
	ColorInfo      = lipgloss.Color("#3B82F6") // Blue
	ColorMuted     = lipgloss.Color("#6B7280") // Gray
	ColorBg        = lipgloss.Color("#1E1B2E") // Dark purple bg
	ColorCardBg    = lipgloss.Color("#2D2640") // Card bg
	ColorText      = lipgloss.Color("#E5E7EB") // Light gray
	ColorTextBold  = lipgloss.Color("#F9FAFB") // White
	ColorHighlight = lipgloss.Color("#C4B5FD") // Light Purple
)

// Grade colors
var GradeColors = map[string]lipgloss.Color{
	"A": lipgloss.Color("#10B981"), // Green
	"B": lipgloss.Color("#34D399"), // Light Green
	"C": lipgloss.Color("#F59E0B"), // Amber
	"D": lipgloss.Color("#F97316"), // Orange
	"F": lipgloss.Color("#EF4444"), // Red
}

// Styles
var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorTextBold).
			Background(ColorPrimary).
			Padding(0, 2).
			MarginBottom(1)

	SubtitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorSecondary).
			Border(lipgloss.NormalBorder(), false, false, true, false). // bottom border
			BorderForeground(ColorSecondary).
			MarginBottom(1)

	CardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPrimary).
			Padding(1, 3).
			MarginBottom(1)

	ActiveCardStyle = lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			BorderForeground(ColorSecondary).
			Padding(1, 3).
			MarginBottom(1)

	GradeStyle = lipgloss.NewStyle().
			Bold(true).
			Padding(1, 4).
			Align(lipgloss.Center)

	ScoreBarFilled = lipgloss.NewStyle().
			Background(ColorSuccess)

	ScoreBarEmpty = lipgloss.NewStyle().
			Background(lipgloss.Color("#374151"))

	FindingCritical = lipgloss.NewStyle().
			Foreground(ColorCritical).
			Bold(true)

	FindingError = lipgloss.NewStyle().
			Foreground(ColorError)

	FindingWarning = lipgloss.NewStyle().
			Foreground(ColorWarning)

	FindingInfo = lipgloss.NewStyle().
			Foreground(ColorInfo)

	MutedStyle = lipgloss.NewStyle().
			Foreground(ColorMuted)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(ColorSuccess).
			Bold(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ColorError).
			Bold(true)

	BannerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary)

	PraiseStyle = lipgloss.NewStyle().
			Foreground(ColorSuccess).
			PaddingLeft(2)

	HighlightStyle = lipgloss.NewStyle().
			Foreground(ColorHighlight).
			Bold(true)
)
