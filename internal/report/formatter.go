package report

import (
	"fmt"
	"strings"
	"time"

	"github.com/Swif7ify/Obelisk-CLI/internal/ai"
	"github.com/Swif7ify/Obelisk-CLI/internal/scanner"
)

// FormatPlainText generates a plain text report from scan results and AI analysis.
func FormatPlainText(result *scanner.ScanResult, report *ai.HealthReport) string {
	var sb strings.Builder

	// Header with timestamp
	sb.WriteString("================================================================================\n")
	sb.WriteString("                           OBELISK HEALTH REPORT\n")
	sb.WriteString("                    The AI-Powered Automated Tech Lead\n")
	sb.WriteString("================================================================================\n\n")
	sb.WriteString(fmt.Sprintf("Generated: %s\n", time.Now().Format("2006-01-02 15:04:05")))
	sb.WriteString(fmt.Sprintf("Project Path: %s\n", result.ProjectPath))
	sb.WriteString(fmt.Sprintf("Project Type: %s\n", result.ProjectType))
	sb.WriteString("\n")

	// Statistics
	sb.WriteString("--------------------------------------------------------------------------------\n")
	sb.WriteString("SCAN STATISTICS\n")
	sb.WriteString("--------------------------------------------------------------------------------\n")
	criticals := result.CountBySeverity(scanner.SeverityCritical)
	errors := result.CountBySeverity(scanner.SeverityError)
	warnings := result.CountBySeverity(scanner.SeverityWarning)
	infos := result.CountBySeverity(scanner.SeverityInfo)

	sb.WriteString(fmt.Sprintf("Files Scanned:     %d\n", result.FileCount))
	sb.WriteString(fmt.Sprintf("Directories:       %d\n", result.DirCount))
	sb.WriteString(fmt.Sprintf("Critical Issues:   %d\n", criticals))
	sb.WriteString(fmt.Sprintf("Errors:            %d\n", errors))
	sb.WriteString(fmt.Sprintf("Warnings:          %d\n", warnings))
	sb.WriteString(fmt.Sprintf("Info:              %d\n", infos))
	sb.WriteString("\n")

	// Health Score Card
	if report != nil {
		sb.WriteString("--------------------------------------------------------------------------------\n")
		sb.WriteString("HEALTH SCORE CARD\n")
		sb.WriteString("--------------------------------------------------------------------------------\n")
		sb.WriteString(fmt.Sprintf("Overall Grade:     %s\n", report.Grade))
		sb.WriteString(fmt.Sprintf("Overall Score:     %d/100\n\n", report.OverallScore))
		sb.WriteString(fmt.Sprintf("Security Score:    %d/100  %s\n", report.SecurityScore, renderTextBar(report.SecurityScore)))
		sb.WriteString(fmt.Sprintf("Architecture:      %d/100  %s\n", report.ArchitectureScore, renderTextBar(report.ArchitectureScore)))
		sb.WriteString(fmt.Sprintf("Quality Score:     %d/100  %s\n", report.QualityScore, renderTextBar(report.QualityScore)))
		sb.WriteString("\n")
	}

	// Findings
	sb.WriteString("--------------------------------------------------------------------------------\n")
	sb.WriteString("FINDINGS\n")
	sb.WriteString("--------------------------------------------------------------------------------\n")
	if len(result.Findings) == 0 {
		sb.WriteString("✓ No issues found — excellent!\n\n")
	} else {
		for i, f := range result.Findings {
			icon := getSeverityIcon(f.Severity)
			sb.WriteString(fmt.Sprintf("%d. %s [%s] %s\n", i+1, icon, f.Severity, f.Title))
			if f.File != "" {
				sb.WriteString(fmt.Sprintf("   File: %s", f.File))
				if f.Line > 0 {
					sb.WriteString(fmt.Sprintf(":%d", f.Line))
				}
				sb.WriteString("\n")
			}
			if f.Description != "" {
				sb.WriteString(fmt.Sprintf("   Description: %s\n", f.Description))
			}
			if f.Suggestion != "" {
				sb.WriteString(fmt.Sprintf("   Suggestion: %s\n", f.Suggestion))
			}
			sb.WriteString("\n")
		}
	}

	// AI Summary
	if report != nil {
		sb.WriteString("--------------------------------------------------------------------------------\n")
		sb.WriteString("AI SUMMARY\n")
		sb.WriteString("--------------------------------------------------------------------------------\n")
		sb.WriteString(report.Summary)
		sb.WriteString("\n\n")

		// Praise
		if len(report.Praise) > 0 {
			sb.WriteString("✨ PRAISE:\n")
			for _, p := range report.Praise {
				sb.WriteString(fmt.Sprintf("  ✓ %s\n", p))
			}
			sb.WriteString("\n")
		}

		// Recommendations
		if len(report.Recommendations) > 0 {
			sb.WriteString("💡 RECOMMENDATIONS:\n")
			for i, r := range report.Recommendations {
				sb.WriteString(fmt.Sprintf("  %d. %s\n", i+1, r))
			}
			sb.WriteString("\n")
		}

		// Top Issues
		if len(report.TopIssues) > 0 {
			sb.WriteString("🔴 TOP PRIORITY ISSUES:\n")
			for i, issue := range report.TopIssues {
				sb.WriteString(fmt.Sprintf("  %d. [%s] %s\n", i+1, issue.Priority, issue.Title))
				if issue.Description != "" {
					sb.WriteString(fmt.Sprintf("     %s\n", issue.Description))
				}
			}
			sb.WriteString("\n")
		}
	}

	sb.WriteString("================================================================================\n")
	sb.WriteString("End of Report\n")
	sb.WriteString("================================================================================\n")

	return sb.String()
}

// renderTextBar creates a simple ASCII progress bar for scores.
func renderTextBar(score int) string {
	barWidth := 30
	filled := score * barWidth / 100
	if filled > barWidth {
		filled = barWidth
	}
	empty := barWidth - filled

	bar := "[" + strings.Repeat("█", filled) + strings.Repeat("░", empty) + "]"
	return bar
}

// getSeverityIcon returns a text icon for the severity level.
func getSeverityIcon(s scanner.Severity) string {
	switch s {
	case scanner.SeverityCritical:
		return "🔴 CRITICAL"
	case scanner.SeverityError:
		return "🟠 ERROR"
	case scanner.SeverityWarning:
		return "🟡 WARNING"
	case scanner.SeverityInfo:
		return "🔵 INFO"
	default:
		return "● UNKNOWN"
	}
}

// Made with Bob
