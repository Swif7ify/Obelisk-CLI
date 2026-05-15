package ui

import "github.com/charmbracelet/lipgloss"

// Color palette
var (
	ColorPrimary   = lipgloss.Color("#7C3AED") // Purple
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
)

// Grade colors
var GradeColors = map[string]lipgloss.Color{
	"A": lipgloss.Color("#10B981"),
	"B": lipgloss.Color("#34D399"),
	"C": lipgloss.Color("#F59E0B"),
	"D": lipgloss.Color("#F97316"),
	"F": lipgloss.Color("#EF4444"),
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
			MarginBottom(1)

	CardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPrimary).
			Padding(1, 2).
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
)
