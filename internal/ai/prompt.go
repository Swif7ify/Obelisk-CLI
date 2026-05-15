package ai

import (
	"fmt"
	"strings"

	"github.com/Swif7ify/Obelisk-CLI/internal/scanner"
)

// BuildPrompt constructs the Gemini prompt from scan results.
func BuildPrompt(result *scanner.ScanResult) string {
	var sb strings.Builder

	sb.WriteString(`You are an expert software architect and security auditor. Analyze the following codebase scan results and produce a health report.

PROJECT INFO:
- Path: ` + result.ProjectPath + `
- Type: ` + result.ProjectType + `
- Files scanned: ` + fmt.Sprintf("%d", result.FileCount) + `
- Directories: ` + fmt.Sprintf("%d", result.DirCount) + `

FINDINGS SUMMARY:
- Critical: ` + fmt.Sprintf("%d", result.CountBySeverity(scanner.SeverityCritical)) + `
- Errors: ` + fmt.Sprintf("%d", result.CountBySeverity(scanner.SeverityError)) + `
- Warnings: ` + fmt.Sprintf("%d", result.CountBySeverity(scanner.SeverityWarning)) + `
- Info: ` + fmt.Sprintf("%d", result.CountBySeverity(scanner.SeverityInfo)) + `
`)

	// Add directory tree for Vibe Check
	if result.DirTree != "" {
		sb.WriteString(`
DIRECTORY TREE (for Vibe Check analysis):
Analyze this structure to determine if the project follows industry-standard patterns
(e.g., MVC, Atomic Design, feature-based, domain-driven). Grade the organizational quality.

` + result.DirTree + `
`)
	}

	sb.WriteString(`
DETAILED FINDINGS:
`)

	for i, f := range result.Findings {
		if i >= 50 { // Cap at 50 findings for prompt size
			sb.WriteString(fmt.Sprintf("\n... and %d more findings\n", len(result.Findings)-50))
			break
		}
		sb.WriteString(fmt.Sprintf("\n%d. [%s][%s] %s", i+1, f.Severity, f.Category, f.Title))
		if f.File != "" {
			sb.WriteString(fmt.Sprintf(" (file: %s", f.File))
			if f.Line > 0 {
				sb.WriteString(fmt.Sprintf(":%d", f.Line))
			}
			sb.WriteString(")")
		}
		sb.WriteString("\n   " + f.Description)
	}

	sb.WriteString(`

INSTRUCTIONS:
Respond with ONLY a JSON object (no markdown, no code fences) in this exact format:
{
  "grade": "A|B|C|D|F",
  "security_score": <0-100>,
  "architecture_score": <0-100>,
  "quality_score": <0-100>,
  "overall_score": <0-100>,
  "top_issues": [
    {"title": "...", "description": "...", "priority": "critical|high|medium|low"}
  ],
  "praise": ["List of things the project does well"],
  "summary": "A 2-3 sentence executive summary of the project health",
  "recommendations": ["Top 3 actionable recommendations"]
}

GRADING RUBRIC:
- A (90-100): Excellent. No critical issues, strong architecture, good practices.
- B (75-89): Good. Minor issues, mostly well-structured.
- C (60-74): Fair. Several issues need attention.
- D (40-59): Poor. Significant problems across multiple areas.
- F (0-39): Failing. Critical security/architectural issues.

Be constructive and specific. Provide actionable advice.`)

	return sb.String()
}
