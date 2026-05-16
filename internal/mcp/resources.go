package mcp

import (
	"encoding/json"
	"fmt"

	"github.com/Swif7ify/Obelisk-CLI/internal/scanner"
)

// ResourceHandler handles resource access
type ResourceHandler struct {
	cache *ResultCache
}

// NewResourceHandler creates a new resource handler
func NewResourceHandler(cache *ResultCache) *ResourceHandler {
	return &ResourceHandler{
		cache: cache,
	}
}

// GetResourcesList returns the list of available resources
func (h *ResourceHandler) GetResourcesList() []Resource {
	return []Resource{
		{
			URI:         "obelisk://scan/latest",
			Name:        "Latest Scan Results",
			Description: "Access the most recent project scan results including all findings",
			MimeType:    "application/json",
		},
		{
			URI:         "obelisk://health/score",
			Name:        "Project Health Score",
			Description: "Current project health grade (A-F) and AI-generated recommendations",
			MimeType:    "application/json",
		},
		{
			URI:         "obelisk://findings/security",
			Name:        "Security Findings",
			Description: "All security-related findings from the latest scan",
			MimeType:    "application/json",
		},
		{
			URI:         "obelisk://findings/quality",
			Name:        "Code Quality Findings",
			Description: "All code quality findings from the latest scan",
			MimeType:    "application/json",
		},
		{
			URI:         "obelisk://findings/architecture",
			Name:        "Architecture Findings",
			Description: "All architecture-related findings from the latest scan",
			MimeType:    "application/json",
		},
	}
}

// ReadResource reads a resource by URI
func (h *ResourceHandler) ReadResource(uri string) (*ReadResourceResult, error) {
	switch uri {
	case "obelisk://scan/latest":
		return h.getLatestScan()
	case "obelisk://health/score":
		return h.getHealthScore()
	case "obelisk://findings/security":
		return h.getSecurityFindings()
	case "obelisk://findings/quality":
		return h.getQualityFindings()
	case "obelisk://findings/architecture":
		return h.getArchitectureFindings()
	default:
		return nil, fmt.Errorf("unknown resource URI: %s", uri)
	}
}

// getLatestScan returns the complete latest scan results
func (h *ResourceHandler) getLatestScan() (*ReadResourceResult, error) {
	result, scanTime, projectPath, ok := h.cache.Get()
	if !ok {
		return nil, fmt.Errorf("no scan results available. Run a scan first using the scan_project tool")
	}

	data := map[string]interface{}{
		"project_path": projectPath,
		"project_type": result.ScanResult.ProjectType,
		"scanned_at":   scanTime.Format("2006-01-02 15:04:05"),
		"file_count":   result.ScanResult.FileCount,
		"dir_count":    result.ScanResult.DirCount,
		"findings": map[string]interface{}{
			"total":    len(result.ScanResult.Findings),
			"critical": result.ScanResult.CountBySeverity(scanner.SeverityCritical),
			"error":    result.ScanResult.CountBySeverity(scanner.SeverityError),
			"warning":  result.ScanResult.CountBySeverity(scanner.SeverityWarning),
			"info":     result.ScanResult.CountBySeverity(scanner.SeverityInfo),
		},
		"findings_list": result.ScanResult.Findings,
	}

	if result.Report != nil {
		data["health_report"] = map[string]interface{}{
			"grade":              result.Report.Grade,
			"overall_score":      result.Report.OverallScore,
			"security_score":     result.Report.SecurityScore,
			"architecture_score": result.Report.ArchitectureScore,
			"quality_score":      result.Report.QualityScore,
			"summary":            result.Report.Summary,
			"praise":             result.Report.Praise,
			"recommendations":    result.Report.Recommendations,
			"top_issues":         result.Report.TopIssues,
		}
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal scan results: %w", err)
	}

	return &ReadResourceResult{
		Contents: []ResourceContent{
			{
				URI:      "obelisk://scan/latest",
				MimeType: "application/json",
				Text:     string(jsonData),
			},
		},
	}, nil
}

// getHealthScore returns just the health score and report
func (h *ResourceHandler) getHealthScore() (*ReadResourceResult, error) {
	report, ok := h.cache.GetHealthScore()
	if !ok {
		return nil, fmt.Errorf("no health report available. Run get_health_report tool first")
	}

	data := map[string]interface{}{
		"grade":              report.Grade,
		"overall_score":      report.OverallScore,
		"security_score":     report.SecurityScore,
		"architecture_score": report.ArchitectureScore,
		"quality_score":      report.QualityScore,
		"summary":            report.Summary,
		"praise":             report.Praise,
		"recommendations":    report.Recommendations,
		"top_issues":         report.TopIssues,
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal health score: %w", err)
	}

	return &ReadResourceResult{
		Contents: []ResourceContent{
			{
				URI:      "obelisk://health/score",
				MimeType: "application/json",
				Text:     string(jsonData),
			},
		},
	}, nil
}

// getSecurityFindings returns security-related findings
func (h *ResourceHandler) getSecurityFindings() (*ReadResourceResult, error) {
	result, _, projectPath, ok := h.cache.Get()
	if !ok {
		return nil, fmt.Errorf("no scan results available. Run a scan first")
	}

	securityFindings := []scanner.Finding{}
	for _, f := range result.ScanResult.Findings {
		if f.Category == scanner.CategorySecurity {
			securityFindings = append(securityFindings, f)
		}
	}

	data := map[string]interface{}{
		"project_path":      projectPath,
		"security_findings": securityFindings,
		"total":             len(securityFindings),
		"critical":          countBySeverity(securityFindings, scanner.SeverityCritical),
		"error":             countBySeverity(securityFindings, scanner.SeverityError),
		"warning":           countBySeverity(securityFindings, scanner.SeverityWarning),
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal security findings: %w", err)
	}

	return &ReadResourceResult{
		Contents: []ResourceContent{
			{
				URI:      "obelisk://findings/security",
				MimeType: "application/json",
				Text:     string(jsonData),
			},
		},
	}, nil
}

// getQualityFindings returns code quality findings
func (h *ResourceHandler) getQualityFindings() (*ReadResourceResult, error) {
	result, _, projectPath, ok := h.cache.Get()
	if !ok {
		return nil, fmt.Errorf("no scan results available. Run a scan first")
	}

	qualityFindings := []scanner.Finding{}
	for _, f := range result.ScanResult.Findings {
		if f.Category == scanner.CategoryQuality {
			qualityFindings = append(qualityFindings, f)
		}
	}

	data := map[string]interface{}{
		"project_path":     projectPath,
		"quality_findings": qualityFindings,
		"total":            len(qualityFindings),
		"error":            countBySeverity(qualityFindings, scanner.SeverityError),
		"warning":          countBySeverity(qualityFindings, scanner.SeverityWarning),
		"info":             countBySeverity(qualityFindings, scanner.SeverityInfo),
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal quality findings: %w", err)
	}

	return &ReadResourceResult{
		Contents: []ResourceContent{
			{
				URI:      "obelisk://findings/quality",
				MimeType: "application/json",
				Text:     string(jsonData),
			},
		},
	}, nil
}

// getArchitectureFindings returns architecture-related findings
func (h *ResourceHandler) getArchitectureFindings() (*ReadResourceResult, error) {
	result, _, projectPath, ok := h.cache.Get()
	if !ok {
		return nil, fmt.Errorf("no scan results available. Run a scan first")
	}

	archFindings := []scanner.Finding{}
	for _, f := range result.ScanResult.Findings {
		if f.Category == scanner.CategoryArchitecture {
			archFindings = append(archFindings, f)
		}
	}

	data := map[string]interface{}{
		"project_path":          projectPath,
		"architecture_findings": archFindings,
		"total":                 len(archFindings),
		"error":                 countBySeverity(archFindings, scanner.SeverityError),
		"warning":               countBySeverity(archFindings, scanner.SeverityWarning),
		"info":                  countBySeverity(archFindings, scanner.SeverityInfo),
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal architecture findings: %w", err)
	}

	return &ReadResourceResult{
		Contents: []ResourceContent{
			{
				URI:      "obelisk://findings/architecture",
				MimeType: "application/json",
				Text:     string(jsonData),
			},
		},
	}, nil
}

// Helper function to count findings by severity
func countBySeverity(findings []scanner.Finding, severity scanner.Severity) int {
	count := 0
	for _, f := range findings {
		if f.Severity == severity {
			count++
		}
	}
	return count
}

// Made with Bob
