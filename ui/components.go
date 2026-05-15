package ui

import (
	"fmt"
	"strings"

	"github.com/Swif7ify/Obelisk-CLI/internal/ai"
	"github.com/Swif7ify/Obelisk-CLI/internal/scanner"
	"github.com/charmbracelet/lipgloss"
)

// RenderBanner returns the Obelisk ASCII banner.
func RenderBanner() string {
	banner := `
   ____  __         ___      __  
  / __ \/ /_  ___  / (_)____/ /__
 / / / / __ \/ _ \/ / / ___/ //_/
/ /_/ / /_/ /  __/ / (__  ) ,<   
\____/_.___/\___/_/_/____/_/|_|  
`
	return BannerStyle.Render(banner) + "\n" +
		MutedStyle.Render("  The AI-Powered Automated Tech Lead") + "\n"
}

// RenderScoreCard renders the grade and scores card.
func RenderScoreCard(report *ai.HealthReport) string {
	gradeColor, ok := GradeColors[report.Grade]
	if !ok {
		gradeColor = ColorMuted
	}

	grade := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(gradeColor).
		Padding(1, 6).
		Align(lipgloss.Center).
		Render(fmt.Sprintf("  %s  ", report.Grade))

	scores := fmt.Sprintf(
		"%s %s\n%s %s\n%s %s\n%s %s",
		SubtitleStyle.Render("Security:"),
		renderScoreBar(report.SecurityScore),
		SubtitleStyle.Render("Architecture:"),
		renderScoreBar(report.ArchitectureScore),
		SubtitleStyle.Render("Quality:"),
		renderScoreBar(report.QualityScore),
		lipgloss.NewStyle().Bold(true).Foreground(ColorTextBold).Render("Overall:"),
		renderScoreBar(report.OverallScore),
	)

	left := lipgloss.NewStyle().Width(20).Align(lipgloss.Center).Render(grade)
	right := lipgloss.NewStyle().Width(50).Render(scores)

	content := lipgloss.JoinHorizontal(lipgloss.Top, left, "  ", right)
	return CardStyle.Width(76).Render(content)
}

// renderScoreBar creates a visual progress bar for a score.
func renderScoreBar(score int) string {
	barWidth := 30
	filled := score * barWidth / 100
	if filled > barWidth {
		filled = barWidth
	}

	color := ColorSuccess
	switch {
	case score < 40:
		color = ColorError
	case score < 60:
		color = ColorWarning
	case score < 75:
		color = lipgloss.Color("#F59E0B")
	}

	bar := lipgloss.NewStyle().Background(color).Render(strings.Repeat(" ", filled))
	empty := lipgloss.NewStyle().Background(lipgloss.Color("#374151")).Render(strings.Repeat(" ", barWidth-filled))

	return fmt.Sprintf("%s%s %s",
		bar, empty,
		lipgloss.NewStyle().Bold(true).Foreground(color).Render(fmt.Sprintf("%d%%", score)),
	)
}

// RenderFindings renders the findings list.
func RenderFindings(findings []scanner.Finding) string {
	if len(findings) == 0 {
		return CardStyle.Width(76).Render(
			SuccessStyle.Render("✓ No issues found — excellent!"),
		)
	}

	var sb strings.Builder
	sb.WriteString(SubtitleStyle.Render("📋 Findings") + "\n\n")

	maxShow := 15
	if len(findings) < maxShow {
		maxShow = len(findings)
	}

	for i := 0; i < maxShow; i++ {
		f := findings[i]
		var style lipgloss.Style
		icon := "●"

		switch f.Severity {
		case scanner.SeverityCritical:
			style = FindingCritical
			icon = "🔴"
		case scanner.SeverityError:
			style = FindingError
			icon = "🟠"
		case scanner.SeverityWarning:
			style = FindingWarning
			icon = "🟡"
		default:
			style = FindingInfo
			icon = "🔵"
		}

		line := fmt.Sprintf("%s %s", icon, style.Render(f.Title))
		if f.File != "" {
			line += MutedStyle.Render(fmt.Sprintf(" (%s)", f.File))
		}
		sb.WriteString(line + "\n")
	}

	if len(findings) > maxShow {
		sb.WriteString(MutedStyle.Render(fmt.Sprintf("\n  ... and %d more findings", len(findings)-maxShow)))
	}

	return CardStyle.Width(76).Render(sb.String())
}

// RenderSummary renders the AI summary and recommendations.
func RenderSummary(report *ai.HealthReport) string {
	var sb strings.Builder

	sb.WriteString(SubtitleStyle.Render("📊 Summary") + "\n\n")
	sb.WriteString(lipgloss.NewStyle().Foreground(ColorText).Render(report.Summary) + "\n")

	if len(report.Praise) > 0 {
		sb.WriteString("\n" + SuccessStyle.Render("✨ Praise:") + "\n")
		for _, p := range report.Praise {
			sb.WriteString(PraiseStyle.Render("  ✓ "+p) + "\n")
		}
	}

	if len(report.Recommendations) > 0 {
		sb.WriteString("\n" + SubtitleStyle.Render("💡 Recommendations:") + "\n")
		for i, r := range report.Recommendations {
			sb.WriteString(fmt.Sprintf("  %d. %s\n", i+1, r))
		}
	}

	return CardStyle.Width(76).Render(sb.String())
}

// RenderStats renders scan statistics.
func RenderStats(result *scanner.ScanResult) string {
	criticals := result.CountBySeverity(scanner.SeverityCritical)
	errors := result.CountBySeverity(scanner.SeverityError)
	warnings := result.CountBySeverity(scanner.SeverityWarning)
	infos := result.CountBySeverity(scanner.SeverityInfo)

	stats := fmt.Sprintf(
		"%s  %s  %s  %s  %s  %s",
		MutedStyle.Render(fmt.Sprintf("📁 %d files", result.FileCount)),
		MutedStyle.Render(fmt.Sprintf("📂 %d dirs", result.DirCount)),
		FindingCritical.Render(fmt.Sprintf("🔴 %d critical", criticals)),
		FindingError.Render(fmt.Sprintf("🟠 %d errors", errors)),
		FindingWarning.Render(fmt.Sprintf("🟡 %d warnings", warnings)),
		FindingInfo.Render(fmt.Sprintf("🔵 %d info", infos)),
	)
	return lipgloss.NewStyle().Padding(0, 1).Render(stats)
}
