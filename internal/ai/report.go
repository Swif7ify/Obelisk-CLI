package ai

import (
	"encoding/json"
	"strings"
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
func FallbackReport(criticals, errors, warnings, infos int, aiErr error) *HealthReport {
	total := criticals + errors + warnings + infos
	score := 100
	
	score -= criticals * 40
	score -= errors * 25
	score -= warnings * 5

	if score < 0 {
		score = 0
	}

	// Hard ceilings for major issues
	if criticals > 0 && score > 39 {
		score = 39 // F
	} else if errors > 0 && score > 59 {
		score = 59 // F
	}

	grade := "A"
	switch {
	case score >= 90:
		grade = "A"
	case score >= 75:
		grade = "B"
	case score >= 60:
		grade = "C"
	case score >= 40:
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
		SecurityScore:     score,
		ArchitectureScore: score,
		QualityScore:      score,
		OverallScore:      score,
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
