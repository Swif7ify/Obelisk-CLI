package ai

import (
	"encoding/json"
	"strings"

	"github.com/Swif7ify/Obelisk-CLI/internal/scanner"
)

// HealthReport represents the AI-generated health assessment.
type HealthReport struct {
	Grade             string      `json:"grade"`
	SecurityScore     int         `json:"security_score"`
	ArchitectureScore int         `json:"architecture_score"`
	QualityScore      int         `json:"quality_score"`
	OverallScore      int         `json:"overall_score"`
	TopIssues         []Issue     `json:"top_issues"`
	Praise            []string    `json:"praise"`
	Summary           string      `json:"summary"`
	Recommendations   []string    `json:"recommendations"`
	Error             string      `json:"error,omitempty"`
}

// Issue represents a prioritized issue from the AI analysis.
type Issue struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
}

// ParseReport parses the AI response into a HealthReport.
func ParseReport(response string) (*HealthReport, error) {
	// Clean the response — strip markdown code fences if present
	cleaned := response
	cleaned = strings.TrimSpace(cleaned)
	if strings.HasPrefix(cleaned, "```json") {
		cleaned = strings.TrimPrefix(cleaned, "```json")
	}
	if strings.HasPrefix(cleaned, "```") {
		cleaned = strings.TrimPrefix(cleaned, "```")
	}
	if strings.HasSuffix(cleaned, "```") {
		cleaned = strings.TrimSuffix(cleaned, "```")
	}
	cleaned = strings.TrimSpace(cleaned)

	var report HealthReport
	if err := json.Unmarshal([]byte(cleaned), &report); err != nil {
		return nil, err
	}

	// Validate grade
	validGrades := map[string]bool{"A": true, "B": true, "C": true, "D": true, "F": true}
	if !validGrades[report.Grade] {
		report.Grade = "C" // Default fallback
	}

	// Clamp scores
	report.SecurityScore = clamp(report.SecurityScore, 0, 100)
	report.ArchitectureScore = clamp(report.ArchitectureScore, 0, 100)
	report.QualityScore = clamp(report.QualityScore, 0, 100)
	report.OverallScore = clamp(report.OverallScore, 0, 100)

	return &report, nil
}

// FallbackReport creates a report when AI is unavailable.
func FallbackReport(result *scanner.ScanResult, aiErr error) *HealthReport {
	secScore := 100
	archScore := 100
	qualScore := 100

	criticals := result.CountBySeverity(scanner.SeverityCritical)
	errors := result.CountBySeverity(scanner.SeverityError)
	warnings := result.CountBySeverity(scanner.SeverityWarning)
	infos := result.CountBySeverity(scanner.SeverityInfo)
	total := criticals + errors + warnings + infos

	secWarnPenalty := 0
	archWarnPenalty := 0
	qualWarnPenalty := 0

	for _, f := range result.Findings {
		penalty := 0
		switch f.Severity {
		case scanner.SeverityCritical:
			penalty = 40
		case scanner.SeverityError:
			penalty = 25
		case scanner.SeverityWarning:
			penalty = 5
		}

		switch f.Category {
		case scanner.CategorySecurity:
			if f.Severity == scanner.SeverityWarning {
				secWarnPenalty += penalty
			} else {
				secScore -= penalty
			}
		case scanner.CategoryArchitecture:
			if f.Severity == scanner.SeverityWarning {
				archWarnPenalty += penalty
			} else {
				archScore -= penalty
			}
		default:
			if f.Severity == scanner.SeverityWarning {
				qualWarnPenalty += penalty
			} else {
				qualScore -= penalty
			}
		}
	}

	// Cap the penalty from warnings so they can't completely tank a category
	if secWarnPenalty > 45 { secWarnPenalty = 45 }
	secScore -= secWarnPenalty

	if archWarnPenalty > 45 { archWarnPenalty = 45 }
	archScore -= archWarnPenalty

	if qualWarnPenalty > 45 { qualWarnPenalty = 45 }
	qualScore -= qualWarnPenalty

	secScore = clamp(secScore, 0, 100)
	archScore = clamp(archScore, 0, 100)
	qualScore = clamp(qualScore, 0, 100)

	overallScore := (secScore + archScore + qualScore) / 3

	grade := "A"
	switch {
	case overallScore >= 90:
		grade = "A"
	case overallScore >= 75:
		grade = "B"
	case overallScore >= 60:
		grade = "C"
	case overallScore >= 40:
		grade = "D"
	default:
		grade = "F"
	}

	summary := "Project analysis completed without AI. "
	if total == 0 {
		summary += "No issues detected — great job!"
	} else {
		summary += "Issues were found that should be addressed."
	}

	errMsg := ""
	if aiErr != nil {
		errMsg = aiErr.Error()
	}

	return &HealthReport{
		Grade:             grade,
		SecurityScore:     secScore,
		ArchitectureScore: archScore,
		QualityScore:      qualScore,
		OverallScore:      overallScore,
		Summary:           summary,
		Praise:            []string{"Scan completed successfully"},
		Recommendations:   []string{"Run with an API key for AI-powered insights"},
		Error:             errMsg,
	}
}

func clamp(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}
