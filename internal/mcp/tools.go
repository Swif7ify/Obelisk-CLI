package mcp

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Swif7ify/Obelisk-CLI/internal/engine"
	"github.com/Swif7ify/Obelisk-CLI/internal/scanner"
)

// ToolHandler handles tool execution
type ToolHandler struct {
	cache  *ResultCache
	apiKey string
	model  string
}

// NewToolHandler creates a new tool handler
func NewToolHandler(cache *ResultCache, apiKey, model string) *ToolHandler {
	return &ToolHandler{
		cache:  cache,
		apiKey: apiKey,
		model:  model,
	}
}

// GetToolsList returns the list of available tools
func (h *ToolHandler) GetToolsList() []Tool {
	return []Tool{
		{
			Name:        "scan_project",
			Description: "Performs a comprehensive health scan of a project including security, architecture, and code quality analysis",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]Property{
					"path": {
						Type:        "string",
						Description: "Project directory path (defaults to current directory)",
					},
					"skip_ai": {
						Type:        "boolean",
						Description: "Skip AI-powered analysis (faster, no API key required)",
						Default:     false,
					},
				},
			},
		},
		{
			Name:        "check_security",
			Description: "Scans for security vulnerabilities including secrets, exposed credentials, and .gitignore issues",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]Property{
					"path": {
						Type:        "string",
						Description: "Project directory path (defaults to current directory)",
					},
				},
			},
		},
		{
			Name:        "analyze_complexity",
			Description: "Analyzes cyclomatic complexity to identify spaghetti code and maintainability issues",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]Property{
					"path": {
						Type:        "string",
						Description: "Project directory path (defaults to current directory)",
					},
					"threshold": {
						Type:        "number",
						Description: "Complexity threshold (default: 10)",
						Default:     10,
					},
				},
			},
		},
		{
			Name:        "track_tech_debt",
			Description: "Tracks TODO, FIXME, HACK, and XXX comments across the codebase",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]Property{
					"path": {
						Type:        "string",
						Description: "Project directory path (defaults to current directory)",
					},
				},
			},
		},
		{
			Name:        "audit_dependencies",
			Description: "Audits package.json for deprecated, vulnerable, or unused dependencies",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]Property{
					"path": {
						Type:        "string",
						Description: "Project directory path (defaults to current directory)",
					},
				},
			},
		},
		{
			Name:        "get_health_report",
			Description: "Generates an AI-powered health report with grade (A-F) and recommendations. Requires GEMINI_API_KEY.",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]Property{
					"path": {
						Type:        "string",
						Description: "Project directory path (defaults to current directory)",
					},
				},
			},
		},
	}
}

// ExecuteTool executes a tool by name
func (h *ToolHandler) ExecuteTool(name string, args map[string]interface{}) (*CallToolResult, error) {
	switch name {
	case "scan_project":
		return h.scanProject(args)
	case "check_security":
		return h.checkSecurity(args)
	case "analyze_complexity":
		return h.analyzeComplexity(args)
	case "track_tech_debt":
		return h.trackTechDebt(args)
	case "audit_dependencies":
		return h.auditDependencies(args)
	case "get_health_report":
		return h.getHealthReport(args)
	default:
		return nil, fmt.Errorf("unknown tool: %s", name)
	}
}

// scanProject performs a full project scan
func (h *ToolHandler) scanProject(args map[string]interface{}) (*CallToolResult, error) {
	projectPath, err := h.getProjectPath(args)
	if err != nil {
		return h.errorResult(err.Error()), nil
	}

	skipAI := false
	if val, ok := args["skip_ai"].(bool); ok {
		skipAI = val
	}

	// Run the scan
	cfg := engine.Config{
		ProjectPath: projectPath,
		APIKey:      h.apiKey,
		Model:       h.model,
		SkipAI:      skipAI || h.apiKey == "",
		Verbose:     false,
	}

	result, err := engine.Run(cfg, nil)
	if err != nil {
		return h.errorResult(fmt.Sprintf("Scan failed: %v", err)), nil
	}

	// Cache the result
	h.cache.Set(result, projectPath)

	// Format output
	output := h.formatScanResult(result)
	return &CallToolResult{
		Content: []Content{
			{
				Type: "text",
				Text: output,
			},
		},
	}, nil
}

// checkSecurity performs security-focused scanning
func (h *ToolHandler) checkSecurity(args map[string]interface{}) (*CallToolResult, error) {
	projectPath, err := h.getProjectPath(args)
	if err != nil {
		return h.errorResult(err.Error()), nil
	}

	// Run full scan but filter for security findings
	cfg := engine.Config{
		ProjectPath: projectPath,
		APIKey:      h.apiKey,
		Model:       h.model,
		SkipAI:      true, // Skip AI for focused scans
		Verbose:     false,
	}

	result, err := engine.Run(cfg, nil)
	if err != nil {
		return h.errorResult(fmt.Sprintf("Security scan failed: %v", err)), nil
	}

	// Filter for security findings
	securityFindings := []scanner.Finding{}
	for _, f := range result.ScanResult.Findings {
		if f.Category == scanner.CategorySecurity {
			securityFindings = append(securityFindings, f)
		}
	}

	output := h.formatSecurityFindings(securityFindings, projectPath)
	return &CallToolResult{
		Content: []Content{
			{
				Type: "text",
				Text: output,
			},
		},
	}, nil
}

// analyzeComplexity analyzes code complexity
func (h *ToolHandler) analyzeComplexity(args map[string]interface{}) (*CallToolResult, error) {
	projectPath, err := h.getProjectPath(args)
	if err != nil {
		return h.errorResult(err.Error()), nil
	}

	// Run complexity scan
	complexityFindings, err := scanner.ScanComplexity(projectPath)
	if err != nil {
		return h.errorResult(fmt.Sprintf("Complexity analysis failed: %v", err)), nil
	}

	output := h.formatComplexityFindings(complexityFindings, projectPath)
	return &CallToolResult{
		Content: []Content{
			{
				Type: "text",
				Text: output,
			},
		},
	}, nil
}

// trackTechDebt tracks technical debt
func (h *ToolHandler) trackTechDebt(args map[string]interface{}) (*CallToolResult, error) {
	projectPath, err := h.getProjectPath(args)
	if err != nil {
		return h.errorResult(err.Error()), nil
	}

	// Run tech debt scan
	techDebtFindings, err := scanner.ScanTechDebt(projectPath)
	if err != nil {
		return h.errorResult(fmt.Sprintf("Tech debt tracking failed: %v", err)), nil
	}

	output := h.formatTechDebtFindings(techDebtFindings, projectPath)
	return &CallToolResult{
		Content: []Content{
			{
				Type: "text",
				Text: output,
			},
		},
	}, nil
}

// auditDependencies audits project dependencies
func (h *ToolHandler) auditDependencies(args map[string]interface{}) (*CallToolResult, error) {
	projectPath, err := h.getProjectPath(args)
	if err != nil {
		return h.errorResult(err.Error()), nil
	}

	// Run dependency scan
	depFindings, err := scanner.ScanDependencies(projectPath)
	if err != nil {
		return h.errorResult(fmt.Sprintf("Dependency audit failed: %v", err)), nil
	}

	output := h.formatDependencyFindings(depFindings, projectPath)
	return &CallToolResult{
		Content: []Content{
			{
				Type: "text",
				Text: output,
			},
		},
	}, nil
}

// getHealthReport generates AI health report
func (h *ToolHandler) getHealthReport(args map[string]interface{}) (*CallToolResult, error) {
	if h.apiKey == "" {
		return h.errorResult("GEMINI_API_KEY environment variable is required for AI health reports"), nil
	}

	projectPath, err := h.getProjectPath(args)
	if err != nil {
		return h.errorResult(err.Error()), nil
	}

	// Run full scan with AI
	cfg := engine.Config{
		ProjectPath: projectPath,
		APIKey:      h.apiKey,
		Model:       h.model,
		SkipAI:      false,
		Verbose:     false,
	}

	result, err := engine.Run(cfg, nil)
	if err != nil {
		return h.errorResult(fmt.Sprintf("Health report generation failed: %v", err)), nil
	}

	// Cache the result
	h.cache.Set(result, projectPath)

	output := h.formatHealthReport(result)
	return &CallToolResult{
		Content: []Content{
			{
				Type: "text",
				Text: output,
			},
		},
	}, nil
}

// Helper functions

func (h *ToolHandler) getProjectPath(args map[string]interface{}) (string, error) {
	var projectPath string
	if val, ok := args["path"].(string); ok && val != "" {
		projectPath = val
	} else {
		var err error
		projectPath, err = os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to get current directory: %w", err)
		}
	}

	// Resolve to absolute path
	absPath, err := filepath.Abs(projectPath)
	if err != nil {
		return "", fmt.Errorf("invalid project path: %w", err)
	}

	// Verify path exists
	if _, err := os.Stat(absPath); err != nil {
		return "", fmt.Errorf("project path does not exist: %s", absPath)
	}

	return absPath, nil
}

func (h *ToolHandler) errorResult(message string) *CallToolResult {
	return &CallToolResult{
		Content: []Content{
			{
				Type: "text",
				Text: message,
			},
		},
		IsError: true,
	}
}

func (h *ToolHandler) formatScanResult(result *engine.Result) string {
	data := map[string]interface{}{
		"project_path": result.ScanResult.ProjectPath,
		"project_type": result.ScanResult.ProjectType,
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
			"grade":            result.Report.Grade,
			"overall_score":    result.Report.OverallScore,
			"security_score":   result.Report.SecurityScore,
			"architecture_score": result.Report.ArchitectureScore,
			"quality_score":    result.Report.QualityScore,
			"summary":          result.Report.Summary,
			"praise":           result.Report.Praise,
			"recommendations":  result.Report.Recommendations,
			"top_issues":       result.Report.TopIssues,
		}
	}

	jsonData, _ := json.MarshalIndent(data, "", "  ")
	return string(jsonData)
}

func (h *ToolHandler) formatSecurityFindings(findings []scanner.Finding, projectPath string) string {
	data := map[string]interface{}{
		"project_path":      projectPath,
		"security_findings": findings,
		"total":             len(findings),
		"critical":          countBySeverity(findings, scanner.SeverityCritical),
		"error":             countBySeverity(findings, scanner.SeverityError),
		"warning":           countBySeverity(findings, scanner.SeverityWarning),
	}

	jsonData, _ := json.MarshalIndent(data, "", "  ")
	return string(jsonData)
}

func (h *ToolHandler) formatComplexityFindings(findings []scanner.Finding, projectPath string) string {
	data := map[string]interface{}{
		"project_path":        projectPath,
		"complexity_findings": findings,
		"total":               len(findings),
		"high_complexity":     countBySeverity(findings, scanner.SeverityError),
		"moderate":            countBySeverity(findings, scanner.SeverityWarning),
	}

	jsonData, _ := json.MarshalIndent(data, "", "  ")
	return string(jsonData)
}

func (h *ToolHandler) formatTechDebtFindings(findings []scanner.Finding, projectPath string) string {
	data := map[string]interface{}{
		"project_path":       projectPath,
		"tech_debt_findings": findings,
		"total":              len(findings),
	}

	jsonData, _ := json.MarshalIndent(data, "", "  ")
	return string(jsonData)
}

func (h *ToolHandler) formatDependencyFindings(findings []scanner.Finding, projectPath string) string {
	data := map[string]interface{}{
		"project_path":         projectPath,
		"dependency_findings":  findings,
		"total":                len(findings),
		"vulnerable":           countBySeverity(findings, scanner.SeverityError),
		"deprecated":           countBySeverity(findings, scanner.SeverityWarning),
	}

	jsonData, _ := json.MarshalIndent(data, "", "  ")
	return string(jsonData)
}

func (h *ToolHandler) formatHealthReport(result *engine.Result) string {
	data := map[string]interface{}{
		"project_path": result.ScanResult.ProjectPath,
		"project_type": result.ScanResult.ProjectType,
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
		data["findings_summary"] = map[string]interface{}{
			"critical": result.ScanResult.CountBySeverity(scanner.SeverityCritical),
			"error":    result.ScanResult.CountBySeverity(scanner.SeverityError),
			"warning":  result.ScanResult.CountBySeverity(scanner.SeverityWarning),
			"info":     result.ScanResult.CountBySeverity(scanner.SeverityInfo),
		}
	}

	jsonData, _ := json.MarshalIndent(data, "", "  ")
	return string(jsonData)
}

// Made with Bob
