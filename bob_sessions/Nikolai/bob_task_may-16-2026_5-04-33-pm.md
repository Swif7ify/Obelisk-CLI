**User:**

<task>
What MCP tools are available?
</task>
<environment_details>
# VSCode Visible Files
README.md

# VSCode Open Tabs
internal/mcp/types.go,internal/mcp/cache.go,internal/mcp/server.go,cmd/mcp.go,internal/mcp/tools.go,internal/mcp/resources.go,../../.bob/settings/mcp_settings.json,MCP_GUIDE.md,README.md

# Current Time
Current time in ISO 8601 UTC format: 2026-05-16T09:04:16.398Z
User time zone: Asia/Singapore, UTC+8:00

# Current Cost
$0.00

# Review Exclusions
The following file patterns are excluded from code reviews:
- .vscode/**

# Current Mode
<slug>code</slug>
<name>💻 Code</name>


# Current Workspace Directory (c:/Users/nikon/Desktop/Obelisk-CLI) Files
.gitignore
CHANGELOG.md
CODE_OF_CONDUCT.md
CONTRIBUTING.md
DISTRIBUTION_CHECKLIST.md
DISTRIBUTION.md
go.mod
go.sum
install.ps1
LICENSE
main.go
Makefile
MCP_GUIDE.md
README.md
SECURITY.md
TESTING.md
uninstall.ps1
.github/
adapters/
adapters/adapter.go
adapters/golang.go
adapters/javascript.go
adapters/laravel.go
cmd/
cmd/check.go
cmd/config_cmd.go
cmd/mcp.go
cmd/protect.go
cmd/report.go
cmd/root.go
cmd/scan.go
cmd/update_test.go
cmd/update.go
cmd/version.go
installer/
installer/ASSETS.md
installer/build-installer.ps1
installer/CODE_SIGNING.md
installer/icon.ico
installer/License.rtf
installer/obelisk.wxs
installer/README.md
internal/
internal/ai/
internal/ai/client.go
internal/ai/prompt.go
internal/ai/report.go
internal/config/
internal/config/config.go
internal/detector/
internal/detector/detector.go
internal/engine/
internal/engine/engine.go
internal/linter/
internal/linter/eslint.go
internal/linter/native.go
internal/mcp/
internal/mcp/cache.go
internal/mcp/resources.go
internal/mcp/server.go
internal/mcp/tools.go
internal/mcp/types.go
internal/report/
internal/report/formatter.go
internal/report/writer.go
internal/scanner/
internal/scanner/complexity.go
internal/scanner/dependencies.go
internal/scanner/gitignore_matcher.go
internal/scanner/gitignore.go
internal/scanner/imports.go
internal/scanner/naming.go
internal/scanner/secrets.go
internal/scanner/techdebt.go
internal/scanner/types.go
packages/
packages/chocolatey/
packages/chocolatey/obelisk-cli.nuspec
packages/chocolatey/README.md
packages/chocolatey/tools/
packages/chocolatey/tools/chocolateyinstall.ps1
packages/chocolatey/tools/chocolateyuninstall.ps1
packages/homebrew/
packages/homebrew/obelisk-cli.rb
packages/homebrew/README.md
packages/winget/
packages/winget/OneDev.ObeliskCLI.installer.yaml
packages/winget/OneDev.ObeliskCLI.locale.en-US.yaml
packages/winget/OneDev.ObeliskCLI.yaml
packages/winget/README.md
ui/
ui/components.go
ui/dashboard.go
ui/handlers.go
ui/input.go
ui/interactive.go
ui/menu.go
ui/messages.go
ui/spinner.go
ui/styles.go
ui/views.go
You have not created a todo list yet. Create one with `update_todo_list` if your task is complicated or involves multiple steps.
</environment_details>

---

**Assistant:**

I'll check what MCP tools are available by reading the MCP-related files in the project.

<read_file>
<args>
<file>
<path>internal/mcp/tools.go</path>
</file>
<file>
<path>internal/mcp/types.go</path>
</file>
<file>
<path>internal/mcp/server.go</path>
</file>
<file>
<path>MCP_GUIDE.md</path>
</file>
</args>
</read_file>

---

**User:**

[read_file for 4 files] Result:
<files>
<file><path>internal/mcp/tools.go</path>
<content lines="1-490">
  1 | package mcp
  2 | 
  3 | import (
  4 | 	"encoding/json"
  5 | 	"fmt"
  6 | 	"os"
  7 | 	"path/filepath"
  8 | 
  9 | 	"github.com/Swif7ify/Obelisk-CLI/internal/engine"
 10 | 	"github.com/Swif7ify/Obelisk-CLI/internal/scanner"
 11 | )
 12 | 
 13 | // ToolHandler handles tool execution
 14 | type ToolHandler struct {
 15 | 	cache  *ResultCache
 16 | 	apiKey string
 17 | 	model  string
 18 | }
 19 | 
 20 | // NewToolHandler creates a new tool handler
 21 | func NewToolHandler(cache *ResultCache, apiKey, model string) *ToolHandler {
 22 | 	return &ToolHandler{
 23 | 		cache:  cache,
 24 | 		apiKey: apiKey,
 25 | 		model:  model,
 26 | 	}
 27 | }
 28 | 
 29 | // GetToolsList returns the list of available tools
 30 | func (h *ToolHandler) GetToolsList() []Tool {
 31 | 	return []Tool{
 32 | 		{
 33 | 			Name:        "scan_project",
 34 | 			Description: "Performs a comprehensive health scan of a project including security, architecture, and code quality analysis",
 35 | 			InputSchema: InputSchema{
 36 | 				Type: "object",
 37 | 				Properties: map[string]Property{
 38 | 					"path": {
 39 | 						Type:        "string",
 40 | 						Description: "Project directory path (defaults to current directory)",
 41 | 					},
 42 | 					"skip_ai": {
 43 | 						Type:        "boolean",
 44 | 						Description: "Skip AI-powered analysis (faster, no API key required)",
 45 | 						Default:     false,
 46 | 					},
 47 | 				},
 48 | 			},
 49 | 		},
 50 | 		{
 51 | 			Name:        "check_security",
 52 | 			Description: "Scans for security vulnerabilities including secrets, exposed credentials, and .gitignore issues",
 53 | 			InputSchema: InputSchema{
 54 | 				Type: "object",
 55 | 				Properties: map[string]Property{
 56 | 					"path": {
 57 | 						Type:        "string",
 58 | 						Description: "Project directory path (defaults to current directory)",
 59 | 					},
 60 | 				},
 61 | 			},
 62 | 		},
 63 | 		{
 64 | 			Name:        "analyze_complexity",
 65 | 			Description: "Analyzes cyclomatic complexity to identify spaghetti code and maintainability issues",
 66 | 			InputSchema: InputSchema{
 67 | 				Type: "object",
 68 | 				Properties: map[string]Property{
 69 | 					"path": {
 70 | 						Type:        "string",
 71 | 						Description: "Project directory path (defaults to current directory)",
 72 | 					},
 73 | 					"threshold": {
 74 | 						Type:        "number",
 75 | 						Description: "Complexity threshold (default: 10)",
 76 | 						Default:     10,
 77 | 					},
 78 | 				},
 79 | 			},
 80 | 		},
 81 | 		{
 82 | 			Name:        "track_tech_debt",
 83 | 			Description: "Tracks TODO, FIXME, HACK, and XXX comments across the codebase",
 84 | 			InputSchema: InputSchema{
 85 | 				Type: "object",
 86 | 				Properties: map[string]Property{
 87 | 					"path": {
 88 | 						Type:        "string",
 89 | 						Description: "Project directory path (defaults to current directory)",
 90 | 					},
 91 | 				},
 92 | 			},
 93 | 		},
 94 | 		{
 95 | 			Name:        "audit_dependencies",
 96 | 			Description: "Audits package.json for deprecated, vulnerable, or unused dependencies",
 97 | 			InputSchema: InputSchema{
 98 | 				Type: "object",
 99 | 				Properties: map[string]Property{
100 | 					"path": {
101 | 						Type:        "string",
102 | 						Description: "Project directory path (defaults to current directory)",
103 | 					},
104 | 				},
105 | 			},
106 | 		},
107 | 		{
108 | 			Name:        "get_health_report",
109 | 			Description: "Generates an AI-powered health report with grade (A-F) and recommendations. Requires GEMINI_API_KEY.",
110 | 			InputSchema: InputSchema{
111 | 				Type: "object",
112 | 				Properties: map[string]Property{
113 | 					"path": {
114 | 						Type:        "string",
115 | 						Description: "Project directory path (defaults to current directory)",
116 | 					},
117 | 				},
118 | 			},
119 | 		},
120 | 	}
121 | }
122 | 
123 | // ExecuteTool executes a tool by name
124 | func (h *ToolHandler) ExecuteTool(name string, args map[string]interface{}) (*CallToolResult, error) {
125 | 	switch name {
126 | 	case "scan_project":
127 | 		return h.scanProject(args)
128 | 	case "check_security":
129 | 		return h.checkSecurity(args)
130 | 	case "analyze_complexity":
131 | 		return h.analyzeComplexity(args)
132 | 	case "track_tech_debt":
133 | 		return h.trackTechDebt(args)
134 | 	case "audit_dependencies":
135 | 		return h.auditDependencies(args)
136 | 	case "get_health_report":
137 | 		return h.getHealthReport(args)
138 | 	default:
139 | 		return nil, fmt.Errorf("unknown tool: %s", name)
140 | 	}
141 | }
142 | 
143 | // scanProject performs a full project scan
144 | func (h *ToolHandler) scanProject(args map[string]interface{}) (*CallToolResult, error) {
145 | 	projectPath, err := h.getProjectPath(args)
146 | 	if err != nil {
147 | 		return h.errorResult(err.Error()), nil
148 | 	}
149 | 
150 | 	skipAI := false
151 | 	if val, ok := args["skip_ai"].(bool); ok {
152 | 		skipAI = val
153 | 	}
154 | 
155 | 	// Run the scan
156 | 	cfg := engine.Config{
157 | 		ProjectPath: projectPath,
158 | 		APIKey:      h.apiKey,
159 | 		Model:       h.model,
160 | 		SkipAI:      skipAI || h.apiKey == "",
161 | 		Verbose:     false,
162 | 	}
163 | 
164 | 	result, err := engine.Run(cfg, nil)
165 | 	if err != nil {
166 | 		return h.errorResult(fmt.Sprintf("Scan failed: %v", err)), nil
167 | 	}
168 | 
169 | 	// Cache the result
170 | 	h.cache.Set(result, projectPath)
171 | 
172 | 	// Format output
173 | 	output := h.formatScanResult(result)
174 | 	return &CallToolResult{
175 | 		Content: []Content{
176 | 			{
177 | 				Type: "text",
178 | 				Text: output,
179 | 			},
180 | 		},
181 | 	}, nil
182 | }
183 | 
184 | // checkSecurity performs security-focused scanning
185 | func (h *ToolHandler) checkSecurity(args map[string]interface{}) (*CallToolResult, error) {
186 | 	projectPath, err := h.getProjectPath(args)
187 | 	if err != nil {
188 | 		return h.errorResult(err.Error()), nil
189 | 	}
190 | 
191 | 	// Run full scan but filter for security findings
192 | 	cfg := engine.Config{
193 | 		ProjectPath: projectPath,
194 | 		APIKey:      h.apiKey,
195 | 		Model:       h.model,
196 | 		SkipAI:      true, // Skip AI for focused scans
197 | 		Verbose:     false,
198 | 	}
199 | 
200 | 	result, err := engine.Run(cfg, nil)
201 | 	if err != nil {
202 | 		return h.errorResult(fmt.Sprintf("Security scan failed: %v", err)), nil
203 | 	}
204 | 
205 | 	// Filter for security findings
206 | 	securityFindings := []scanner.Finding{}
207 | 	for _, f := range result.ScanResult.Findings {
208 | 		if f.Category == scanner.CategorySecurity {
209 | 			securityFindings = append(securityFindings, f)
210 | 		}
211 | 	}
212 | 
213 | 	output := h.formatSecurityFindings(securityFindings, projectPath)
214 | 	return &CallToolResult{
215 | 		Content: []Content{
216 | 			{
217 | 				Type: "text",
218 | 				Text: output,
219 | 			},
220 | 		},
221 | 	}, nil
222 | }
223 | 
224 | // analyzeComplexity analyzes code complexity
225 | func (h *ToolHandler) analyzeComplexity(args map[string]interface{}) (*CallToolResult, error) {
226 | 	projectPath, err := h.getProjectPath(args)
227 | 	if err != nil {
228 | 		return h.errorResult(err.Error()), nil
229 | 	}
230 | 
231 | 	// Run complexity scan
232 | 	complexityFindings, err := scanner.ScanComplexity(projectPath)
233 | 	if err != nil {
234 | 		return h.errorResult(fmt.Sprintf("Complexity analysis failed: %v", err)), nil
235 | 	}
236 | 
237 | 	output := h.formatComplexityFindings(complexityFindings, projectPath)
238 | 	return &CallToolResult{
239 | 		Content: []Content{
240 | 			{
241 | 				Type: "text",
242 | 				Text: output,
243 | 			},
244 | 		},
245 | 	}, nil
246 | }
247 | 
248 | // trackTechDebt tracks technical debt
249 | func (h *ToolHandler) trackTechDebt(args map[string]interface{}) (*CallToolResult, error) {
250 | 	projectPath, err := h.getProjectPath(args)
251 | 	if err != nil {
252 | 		return h.errorResult(err.Error()), nil
253 | 	}
254 | 
255 | 	// Run tech debt scan
256 | 	techDebtFindings, err := scanner.ScanTechDebt(projectPath)
257 | 	if err != nil {
258 | 		return h.errorResult(fmt.Sprintf("Tech debt tracking failed: %v", err)), nil
259 | 	}
260 | 
261 | 	output := h.formatTechDebtFindings(techDebtFindings, projectPath)
262 | 	return &CallToolResult{
263 | 		Content: []Content{
264 | 			{
265 | 				Type: "text",
266 | 				Text: output,
267 | 			},
268 | 		},
269 | 	}, nil
270 | }
271 | 
272 | // auditDependencies audits project dependencies
273 | func (h *ToolHandler) auditDependencies(args map[string]interface{}) (*CallToolResult, error) {
274 | 	projectPath, err := h.getProjectPath(args)
275 | 	if err != nil {
276 | 		return h.errorResult(err.Error()), nil
277 | 	}
278 | 
279 | 	// Run dependency scan
280 | 	depFindings, err := scanner.ScanDependencies(projectPath)
281 | 	if err != nil {
282 | 		return h.errorResult(fmt.Sprintf("Dependency audit failed: %v", err)), nil
283 | 	}
284 | 
285 | 	output := h.formatDependencyFindings(depFindings, projectPath)
286 | 	return &CallToolResult{
287 | 		Content: []Content{
288 | 			{
289 | 				Type: "text",
290 | 				Text: output,
291 | 			},
292 | 		},
293 | 	}, nil
294 | }
295 | 
296 | // getHealthReport generates AI health report
297 | func (h *ToolHandler) getHealthReport(args map[string]interface{}) (*CallToolResult, error) {
298 | 	if h.apiKey == "" {
299 | 		return h.errorResult("GEMINI_API_KEY environment variable is required for AI health reports"), nil
300 | 	}
301 | 
302 | 	projectPath, err := h.getProjectPath(args)
303 | 	if err != nil {
304 | 		return h.errorResult(err.Error()), nil
305 | 	}
306 | 
307 | 	// Run full scan with AI
308 | 	cfg := engine.Config{
309 | 		ProjectPath: projectPath,
310 | 		APIKey:      h.apiKey,
311 | 		Model:       h.model,
312 | 		SkipAI:      false,
313 | 		Verbose:     false,
314 | 	}
315 | 
316 | 	result, err := engine.Run(cfg, nil)
317 | 	if err != nil {
318 | 		return h.errorResult(fmt.Sprintf("Health report generation failed: %v", err)), nil
319 | 	}
320 | 
321 | 	// Cache the result
322 | 	h.cache.Set(result, projectPath)
323 | 
324 | 	output := h.formatHealthReport(result)
325 | 	return &CallToolResult{
326 | 		Content: []Content{
327 | 			{
328 | 				Type: "text",
329 | 				Text: output,
330 | 			},
331 | 		},
332 | 	}, nil
333 | }
334 | 
335 | // Helper functions
336 | 
337 | func (h *ToolHandler) getProjectPath(args map[string]interface{}) (string, error) {
338 | 	var projectPath string
339 | 	if val, ok := args["path"].(string); ok && val != "" {
340 | 		projectPath = val
341 | 	} else {
342 | 		var err error
343 | 		projectPath, err = os.Getwd()
344 | 		if err != nil {
345 | 			return "", fmt.Errorf("failed to get current directory: %w", err)
346 | 		}
347 | 	}
348 | 
349 | 	// Resolve to absolute path
350 | 	absPath, err := filepath.Abs(projectPath)
351 | 	if err != nil {
352 | 		return "", fmt.Errorf("invalid project path: %w", err)
353 | 	}
354 | 
355 | 	// Verify path exists
356 | 	if _, err := os.Stat(absPath); err != nil {
357 | 		return "", fmt.Errorf("project path does not exist: %s", absPath)
358 | 	}
359 | 
360 | 	return absPath, nil
361 | }
362 | 
363 | func (h *ToolHandler) errorResult(message string) *CallToolResult {
364 | 	return &CallToolResult{
365 | 		Content: []Content{
366 | 			{
367 | 				Type: "text",
368 | 				Text: message,
369 | 			},
370 | 		},
371 | 		IsError: true,
372 | 	}
373 | }
374 | 
375 | func (h *ToolHandler) formatScanResult(result *engine.Result) string {
376 | 	data := map[string]interface{}{
377 | 		"project_path": result.ScanResult.ProjectPath,
378 | 		"project_type": result.ScanResult.ProjectType,
379 | 		"file_count":   result.ScanResult.FileCount,
380 | 		"dir_count":    result.ScanResult.DirCount,
381 | 		"findings": map[string]interface{}{
382 | 			"total":    len(result.ScanResult.Findings),
383 | 			"critical": result.ScanResult.CountBySeverity(scanner.SeverityCritical),
384 | 			"error":    result.ScanResult.CountBySeverity(scanner.SeverityError),
385 | 			"warning":  result.ScanResult.CountBySeverity(scanner.SeverityWarning),
386 | 			"info":     result.ScanResult.CountBySeverity(scanner.SeverityInfo),
387 | 		},
388 | 		"findings_list": result.ScanResult.Findings,
389 | 	}
390 | 
391 | 	if result.Report != nil {
392 | 		data["health_report"] = map[string]interface{}{
393 | 			"grade":            result.Report.Grade,
394 | 			"overall_score":    result.Report.OverallScore,
395 | 			"security_score":   result.Report.SecurityScore,
396 | 			"architecture_score": result.Report.ArchitectureScore,
397 | 			"quality_score":    result.Report.QualityScore,
398 | 			"summary":          result.Report.Summary,
399 | 			"praise":           result.Report.Praise,
400 | 			"recommendations":  result.Report.Recommendations,
401 | 			"top_issues":       result.Report.TopIssues,
402 | 		}
403 | 	}
404 | 
405 | 	jsonData, _ := json.MarshalIndent(data, "", "  ")
406 | 	return string(jsonData)
407 | }
408 | 
409 | func (h *ToolHandler) formatSecurityFindings(findings []scanner.Finding, projectPath string) string {
410 | 	data := map[string]interface{}{
411 | 		"project_path":      projectPath,
412 | 		"security_findings": findings,
413 | 		"total":             len(findings),
414 | 		"critical":          countBySeverity(findings, scanner.SeverityCritical),
415 | 		"error":             countBySeverity(findings, scanner.SeverityError),
416 | 		"warning":           countBySeverity(findings, scanner.SeverityWarning),
417 | 	}
418 | 
419 | 	jsonData, _ := json.MarshalIndent(data, "", "  ")
420 | 	return string(jsonData)
421 | }
422 | 
423 | func (h *ToolHandler) formatComplexityFindings(findings []scanner.Finding, projectPath string) string {
424 | 	data := map[string]interface{}{
425 | 		"project_path":        projectPath,
426 | 		"complexity_findings": findings,
427 | 		"total":               len(findings),
428 | 		"high_complexity":     countBySeverity(findings, scanner.SeverityError),
429 | 		"moderate":            countBySeverity(findings, scanner.SeverityWarning),
430 | 	}
431 | 
432 | 	jsonData, _ := json.MarshalIndent(data, "", "  ")
433 | 	return string(jsonData)
434 | }
435 | 
436 | func (h *ToolHandler) formatTechDebtFindings(findings []scanner.Finding, projectPath string) string {
437 | 	data := map[string]interface{}{
438 | 		"project_path":       projectPath,
439 | 		"tech_debt_findings": findings,
440 | 		"total":              len(findings),
441 | 	}
442 | 
443 | 	jsonData, _ := json.MarshalIndent(data, "", "  ")
444 | 	return string(jsonData)
445 | }
446 | 
447 | func (h *ToolHandler) formatDependencyFindings(findings []scanner.Finding, projectPath string) string {
448 | 	data := map[string]interface{}{
449 | 		"project_path":         projectPath,
450 | 		"dependency_findings":  findings,
451 | 		"total":                len(findings),
452 | 		"vulnerable":           countBySeverity(findings, scanner.SeverityError),
453 | 		"deprecated":           countBySeverity(findings, scanner.SeverityWarning),
454 | 	}
455 | 
456 | 	jsonData, _ := json.MarshalIndent(data, "", "  ")
457 | 	return string(jsonData)
458 | }
459 | 
460 | func (h *ToolHandler) formatHealthReport(result *engine.Result) string {
461 | 	data := map[string]interface{}{
462 | 		"project_path": result.ScanResult.ProjectPath,
463 | 		"project_type": result.ScanResult.ProjectType,
464 | 	}
465 | 
466 | 	if result.Report != nil {
467 | 		data["health_report"] = map[string]interface{}{
468 | 			"grade":              result.Report.Grade,
469 | 			"overall_score":      result.Report.OverallScore,
470 | 			"security_score":     result.Report.SecurityScore,
471 | 			"architecture_score": result.Report.ArchitectureScore,
472 | 			"quality_score":      result.Report.QualityScore,
473 | 			"summary":            result.Report.Summary,
474 | 			"praise":             result.Report.Praise,
475 | 			"recommendations":    result.Report.Recommendations,
476 | 			"top_issues":         result.Report.TopIssues,
477 | 		}
478 | 		data["findings_summary"] = map[string]interface{}{
479 | 			"critical": result.ScanResult.CountBySeverity(scanner.SeverityCritical),
480 | 			"error":    result.ScanResult.CountBySeverity(scanner.SeverityError),
481 | 			"warning":  result.ScanResult.CountBySeverity(scanner.SeverityWarning),
482 | 			"info":     result.ScanResult.CountBySeverity(scanner.SeverityInfo),
483 | 		}
484 | 	}
485 | 
486 | 	jsonData, _ := json.MarshalIndent(data, "", "  ")
487 | 	return string(jsonData)
488 | }
489 | 
490 | // Made with Bob
</content>
</file>
<file><path>internal/mcp/types.go</path>
<content lines="1-186">
  1 | package mcp
  2 | 
  3 | // JSON-RPC 2.0 message types for MCP protocol
  4 | 
  5 | // JSONRPCRequest represents a JSON-RPC 2.0 request
  6 | type JSONRPCRequest struct {
  7 | 	JSONRPC string                 `json:"jsonrpc"`
  8 | 	ID      interface{}            `json:"id,omitempty"`
  9 | 	Method  string                 `json:"method"`
 10 | 	Params  map[string]interface{} `json:"params,omitempty"`
 11 | }
 12 | 
 13 | // JSONRPCResponse represents a JSON-RPC 2.0 response
 14 | type JSONRPCResponse struct {
 15 | 	JSONRPC string      `json:"jsonrpc"`
 16 | 	ID      interface{} `json:"id,omitempty"`
 17 | 	Result  interface{} `json:"result,omitempty"`
 18 | 	Error   *RPCError   `json:"error,omitempty"`
 19 | }
 20 | 
 21 | // RPCError represents a JSON-RPC 2.0 error
 22 | type RPCError struct {
 23 | 	Code    int         `json:"code"`
 24 | 	Message string      `json:"message"`
 25 | 	Data    interface{} `json:"data,omitempty"`
 26 | }
 27 | 
 28 | // Standard JSON-RPC error codes
 29 | const (
 30 | 	ParseError     = -32700
 31 | 	InvalidRequest = -32600
 32 | 	MethodNotFound = -32601
 33 | 	InvalidParams  = -32602
 34 | 	InternalError  = -32603
 35 | )
 36 | 
 37 | // MCP Protocol Types
 38 | 
 39 | // InitializeRequest represents the initialize method params
 40 | type InitializeRequest struct {
 41 | 	ProtocolVersion string                 `json:"protocolVersion"`
 42 | 	Capabilities    ClientCapabilities     `json:"capabilities"`
 43 | 	ClientInfo      ClientInfo             `json:"clientInfo"`
 44 | 	Meta            map[string]interface{} `json:"_meta,omitempty"`
 45 | }
 46 | 
 47 | // ClientCapabilities describes what the client supports
 48 | type ClientCapabilities struct {
 49 | 	Roots       *RootsCapability       `json:"roots,omitempty"`
 50 | 	Sampling    map[string]interface{} `json:"sampling,omitempty"`
 51 | 	Experimental map[string]interface{} `json:"experimental,omitempty"`
 52 | }
 53 | 
 54 | // RootsCapability indicates if client supports roots
 55 | type RootsCapability struct {
 56 | 	ListChanged bool `json:"listChanged,omitempty"`
 57 | }
 58 | 
 59 | // ClientInfo contains information about the client
 60 | type ClientInfo struct {
 61 | 	Name    string `json:"name"`
 62 | 	Version string `json:"version"`
 63 | }
 64 | 
 65 | // InitializeResult is the response to initialize
 66 | type InitializeResult struct {
 67 | 	ProtocolVersion string             `json:"protocolVersion"`
 68 | 	Capabilities    ServerCapabilities `json:"capabilities"`
 69 | 	ServerInfo      ServerInfo         `json:"serverInfo"`
 70 | }
 71 | 
 72 | // ServerCapabilities describes what the server supports
 73 | type ServerCapabilities struct {
 74 | 	Tools     *ToolsCapability     `json:"tools,omitempty"`
 75 | 	Resources *ResourcesCapability `json:"resources,omitempty"`
 76 | 	Prompts   *PromptsCapability   `json:"prompts,omitempty"`
 77 | 	Logging   map[string]interface{} `json:"logging,omitempty"`
 78 | }
 79 | 
 80 | // ToolsCapability indicates tool support
 81 | type ToolsCapability struct {
 82 | 	ListChanged bool `json:"listChanged,omitempty"`
 83 | }
 84 | 
 85 | // ResourcesCapability indicates resource support
 86 | type ResourcesCapability struct {
 87 | 	Subscribe   bool `json:"subscribe,omitempty"`
 88 | 	ListChanged bool `json:"listChanged,omitempty"`
 89 | }
 90 | 
 91 | // PromptsCapability indicates prompt support
 92 | type PromptsCapability struct {
 93 | 	ListChanged bool `json:"listChanged,omitempty"`
 94 | }
 95 | 
 96 | // ServerInfo contains information about the server
 97 | type ServerInfo struct {
 98 | 	Name    string `json:"name"`
 99 | 	Version string `json:"version"`
100 | }
101 | 
102 | // Tool represents an MCP tool definition
103 | type Tool struct {
104 | 	Name        string      `json:"name"`
105 | 	Description string      `json:"description"`
106 | 	InputSchema InputSchema `json:"inputSchema"`
107 | }
108 | 
109 | // InputSchema defines the JSON schema for tool parameters
110 | type InputSchema struct {
111 | 	Type       string                 `json:"type"`
112 | 	Properties map[string]Property    `json:"properties,omitempty"`
113 | 	Required   []string               `json:"required,omitempty"`
114 | 	Additional map[string]interface{} `json:"additionalProperties,omitempty"`
115 | }
116 | 
117 | // Property defines a parameter property
118 | type Property struct {
119 | 	Type        string        `json:"type"`
120 | 	Description string        `json:"description,omitempty"`
121 | 	Enum        []string      `json:"enum,omitempty"`
122 | 	Default     interface{}   `json:"default,omitempty"`
123 | 	Items       *PropertyType `json:"items,omitempty"`
124 | }
125 | 
126 | // PropertyType for array items
127 | type PropertyType struct {
128 | 	Type string `json:"type"`
129 | }
130 | 
131 | // ToolsListResult is the response to tools/list
132 | type ToolsListResult struct {
133 | 	Tools []Tool `json:"tools"`
134 | }
135 | 
136 | // CallToolRequest represents a tools/call request params
137 | type CallToolRequest struct {
138 | 	Name      string                 `json:"name"`
139 | 	Arguments map[string]interface{} `json:"arguments,omitempty"`
140 | }
141 | 
142 | // CallToolResult is the response to tools/call
143 | type CallToolResult struct {
144 | 	Content []Content `json:"content"`
145 | 	IsError bool      `json:"isError,omitempty"`
146 | }
147 | 
148 | // Content represents tool output content
149 | type Content struct {
150 | 	Type     string `json:"type"`
151 | 	Text     string `json:"text,omitempty"`
152 | 	MimeType string `json:"mimeType,omitempty"`
153 | }
154 | 
155 | // Resource represents an MCP resource definition
156 | type Resource struct {
157 | 	URI         string `json:"uri"`
158 | 	Name        string `json:"name"`
159 | 	Description string `json:"description,omitempty"`
160 | 	MimeType    string `json:"mimeType,omitempty"`
161 | }
162 | 
163 | // ResourcesListResult is the response to resources/list
164 | type ResourcesListResult struct {
165 | 	Resources []Resource `json:"resources"`
166 | }
167 | 
168 | // ReadResourceRequest represents a resources/read request params
169 | type ReadResourceRequest struct {
170 | 	URI string `json:"uri"`
171 | }
172 | 
173 | // ReadResourceResult is the response to resources/read
174 | type ReadResourceResult struct {
175 | 	Contents []ResourceContent `json:"contents"`
176 | }
177 | 
178 | // ResourceContent represents resource data
179 | type ResourceContent struct {
180 | 	URI      string `json:"uri"`
181 | 	MimeType string `json:"mimeType"`
182 | 	Text     string `json:"text,omitempty"`
183 | 	Blob     string `json:"blob,omitempty"`
184 | }
185 | 
186 | // Made with Bob
</content>
</file>
<file><path>internal/mcp/server.go</path>
<content lines="1-276">
  1 | package mcp
  2 | 
  3 | import (
  4 | 	"bufio"
  5 | 	"encoding/json"
  6 | 	"fmt"
  7 | 	"io"
  8 | 	"os"
  9 | )
 10 | 
 11 | // Server represents the MCP server
 12 | type Server struct {
 13 | 	toolHandler     *ToolHandler
 14 | 	resourceHandler *ResourceHandler
 15 | 	cache           *ResultCache
 16 | 	version         string
 17 | 	initialized     bool
 18 | }
 19 | 
 20 | // NewServer creates a new MCP server
 21 | func NewServer(version string) *Server {
 22 | 	cache := NewResultCache()
 23 | 	
 24 | 	// Get API key from environment
 25 | 	apiKey := os.Getenv("GEMINI_API_KEY")
 26 | 	if apiKey == "" {
 27 | 		apiKey = os.Getenv("GOOGLE_API_KEY")
 28 | 	}
 29 | 	
 30 | 	// Get model from environment or use default
 31 | 	model := os.Getenv("OBELISK_MODEL")
 32 | 	if model == "" {
 33 | 		model = "gemini-2.0-flash-exp"
 34 | 	}
 35 | 	
 36 | 	return &Server{
 37 | 		toolHandler:     NewToolHandler(cache, apiKey, model),
 38 | 		resourceHandler: NewResourceHandler(cache),
 39 | 		cache:           cache,
 40 | 		version:         version,
 41 | 		initialized:     false,
 42 | 	}
 43 | }
 44 | 
 45 | // Run starts the MCP server with stdio transport
 46 | func (s *Server) Run() error {
 47 | 	reader := bufio.NewReader(os.Stdin)
 48 | 	
 49 | 	// Log to stderr (stdout is reserved for JSON-RPC)
 50 | 	fmt.Fprintln(os.Stderr, "🏛️ Obelisk MCP Server starting...")
 51 | 	fmt.Fprintf(os.Stderr, "Version: %s\n", s.version)
 52 | 	fmt.Fprintln(os.Stderr, "Listening on stdio for JSON-RPC 2.0 messages...")
 53 | 	
 54 | 	for {
 55 | 		// Read a line from stdin
 56 | 		line, err := reader.ReadBytes('\n')
 57 | 		if err != nil {
 58 | 			if err == io.EOF {
 59 | 				fmt.Fprintln(os.Stderr, "Client disconnected")
 60 | 				return nil
 61 | 			}
 62 | 			return fmt.Errorf("error reading from stdin: %w", err)
 63 | 		}
 64 | 		
 65 | 		// Parse JSON-RPC request
 66 | 		var req JSONRPCRequest
 67 | 		if err := json.Unmarshal(line, &req); err != nil {
 68 | 			s.sendError(nil, ParseError, "Invalid JSON", err.Error())
 69 | 			continue
 70 | 		}
 71 | 		
 72 | 		// Handle the request
 73 | 		s.handleRequest(&req)
 74 | 	}
 75 | }
 76 | 
 77 | // handleRequest processes a JSON-RPC request
 78 | func (s *Server) handleRequest(req *JSONRPCRequest) {
 79 | 	// Validate JSON-RPC version
 80 | 	if req.JSONRPC != "2.0" {
 81 | 		s.sendError(req.ID, InvalidRequest, "Invalid JSON-RPC version", "Expected 2.0")
 82 | 		return
 83 | 	}
 84 | 	
 85 | 	// Route to appropriate handler
 86 | 	switch req.Method {
 87 | 	case "initialize":
 88 | 		s.handleInitialize(req)
 89 | 	case "tools/list":
 90 | 		s.handleToolsList(req)
 91 | 	case "tools/call":
 92 | 		s.handleToolsCall(req)
 93 | 	case "resources/list":
 94 | 		s.handleResourcesList(req)
 95 | 	case "resources/read":
 96 | 		s.handleResourcesRead(req)
 97 | 	case "ping":
 98 | 		s.handlePing(req)
 99 | 	default:
100 | 		s.sendError(req.ID, MethodNotFound, fmt.Sprintf("Method not found: %s", req.Method), nil)
101 | 	}
102 | }
103 | 
104 | // handleInitialize handles the initialize method
105 | func (s *Server) handleInitialize(req *JSONRPCRequest) {
106 | 	// Parse initialize params
107 | 	var initReq InitializeRequest
108 | 	if req.Params != nil {
109 | 		paramsJSON, _ := json.Marshal(req.Params)
110 | 		if err := json.Unmarshal(paramsJSON, &initReq); err != nil {
111 | 			s.sendError(req.ID, InvalidParams, "Invalid initialize parameters", err.Error())
112 | 			return
113 | 		}
114 | 	}
115 | 	
116 | 	// Log client info
117 | 	fmt.Fprintf(os.Stderr, "Client connected: %s v%s\n", initReq.ClientInfo.Name, initReq.ClientInfo.Version)
118 | 	
119 | 	// Create initialize result
120 | 	result := InitializeResult{
121 | 		ProtocolVersion: "2024-11-05",
122 | 		Capabilities: ServerCapabilities{
123 | 			Tools: &ToolsCapability{
124 | 				ListChanged: false,
125 | 			},
126 | 			Resources: &ResourcesCapability{
127 | 				Subscribe:   false,
128 | 				ListChanged: false,
129 | 			},
130 | 		},
131 | 		ServerInfo: ServerInfo{
132 | 			Name:    "obelisk-mcp-server",
133 | 			Version: s.version,
134 | 		},
135 | 	}
136 | 	
137 | 	s.initialized = true
138 | 	s.sendResult(req.ID, result)
139 | }
140 | 
141 | // handleToolsList handles the tools/list method
142 | func (s *Server) handleToolsList(req *JSONRPCRequest) {
143 | 	if !s.initialized {
144 | 		s.sendError(req.ID, InvalidRequest, "Server not initialized", "Call initialize first")
145 | 		return
146 | 	}
147 | 	
148 | 	tools := s.toolHandler.GetToolsList()
149 | 	result := ToolsListResult{
150 | 		Tools: tools,
151 | 	}
152 | 	
153 | 	s.sendResult(req.ID, result)
154 | }
155 | 
156 | // handleToolsCall handles the tools/call method
157 | func (s *Server) handleToolsCall(req *JSONRPCRequest) {
158 | 	if !s.initialized {
159 | 		s.sendError(req.ID, InvalidRequest, "Server not initialized", "Call initialize first")
160 | 		return
161 | 	}
162 | 	
163 | 	// Parse call tool params
164 | 	var callReq CallToolRequest
165 | 	if req.Params != nil {
166 | 		paramsJSON, _ := json.Marshal(req.Params)
167 | 		if err := json.Unmarshal(paramsJSON, &callReq); err != nil {
168 | 			s.sendError(req.ID, InvalidParams, "Invalid tool call parameters", err.Error())
169 | 			return
170 | 		}
171 | 	}
172 | 	
173 | 	fmt.Fprintf(os.Stderr, "Executing tool: %s\n", callReq.Name)
174 | 	
175 | 	// Execute the tool
176 | 	result, err := s.toolHandler.ExecuteTool(callReq.Name, callReq.Arguments)
177 | 	if err != nil {
178 | 		s.sendError(req.ID, InternalError, fmt.Sprintf("Tool execution failed: %s", callReq.Name), err.Error())
179 | 		return
180 | 	}
181 | 	
182 | 	s.sendResult(req.ID, result)
183 | }
184 | 
185 | // handleResourcesList handles the resources/list method
186 | func (s *Server) handleResourcesList(req *JSONRPCRequest) {
187 | 	if !s.initialized {
188 | 		s.sendError(req.ID, InvalidRequest, "Server not initialized", "Call initialize first")
189 | 		return
190 | 	}
191 | 	
192 | 	resources := s.resourceHandler.GetResourcesList()
193 | 	result := ResourcesListResult{
194 | 		Resources: resources,
195 | 	}
196 | 	
197 | 	s.sendResult(req.ID, result)
198 | }
199 | 
200 | // handleResourcesRead handles the resources/read method
201 | func (s *Server) handleResourcesRead(req *JSONRPCRequest) {
202 | 	if !s.initialized {
203 | 		s.sendError(req.ID, InvalidRequest, "Server not initialized", "Call initialize first")
204 | 		return
205 | 	}
206 | 	
207 | 	// Parse read resource params
208 | 	var readReq ReadResourceRequest
209 | 	if req.Params != nil {
210 | 		paramsJSON, _ := json.Marshal(req.Params)
211 | 		if err := json.Unmarshal(paramsJSON, &readReq); err != nil {
212 | 			s.sendError(req.ID, InvalidParams, "Invalid resource read parameters", err.Error())
213 | 			return
214 | 		}
215 | 	}
216 | 	
217 | 	fmt.Fprintf(os.Stderr, "Reading resource: %s\n", readReq.URI)
218 | 	
219 | 	// Read the resource
220 | 	result, err := s.resourceHandler.ReadResource(readReq.URI)
221 | 	if err != nil {
222 | 		s.sendError(req.ID, InternalError, fmt.Sprintf("Resource read failed: %s", readReq.URI), err.Error())
223 | 		return
224 | 	}
225 | 	
226 | 	s.sendResult(req.ID, result)
227 | }
228 | 
229 | // handlePing handles the ping method (for testing)
230 | func (s *Server) handlePing(req *JSONRPCRequest) {
231 | 	result := map[string]string{
232 | 		"status": "ok",
233 | 		"server": "obelisk-mcp-server",
234 | 	}
235 | 	s.sendResult(req.ID, result)
236 | }
237 | 
238 | // sendResult sends a successful JSON-RPC response
239 | func (s *Server) sendResult(id interface{}, result interface{}) {
240 | 	response := JSONRPCResponse{
241 | 		JSONRPC: "2.0",
242 | 		ID:      id,
243 | 		Result:  result,
244 | 	}
245 | 	
246 | 	s.sendResponse(response)
247 | }
248 | 
249 | // sendError sends an error JSON-RPC response
250 | func (s *Server) sendError(id interface{}, code int, message string, data interface{}) {
251 | 	response := JSONRPCResponse{
252 | 		JSONRPC: "2.0",
253 | 		ID:      id,
254 | 		Error: &RPCError{
255 | 			Code:    code,
256 | 			Message: message,
257 | 			Data:    data,
258 | 		},
259 | 	}
260 | 	
261 | 	s.sendResponse(response)
262 | }
263 | 
264 | // sendResponse writes a JSON-RPC response to stdout
265 | func (s *Server) sendResponse(response JSONRPCResponse) {
266 | 	data, err := json.Marshal(response)
267 | 	if err != nil {
268 | 		fmt.Fprintf(os.Stderr, "Error marshaling response: %v\n", err)
269 | 		return
270 | 	}
271 | 	
272 | 	// Write to stdout with newline
273 | 	fmt.Println(string(data))
274 | }
275 | 
276 | // Made with Bob
</content>
</file>
<file><path>MCP_GUIDE.md</path>
<content lines="1-449">
  1 | # 🏛️ Obelisk MCP Server Guide
  2 | 
  3 | ## Overview
  4 | 
  5 | Obelisk can now run as a **Model Context Protocol (MCP) server**, exposing its powerful code analysis capabilities to AI assistants and IDEs that support MCP. This allows AI assistants to scan your projects, check security, analyze code quality, and provide architectural insights programmatically.
  6 | 
  7 | ## What is MCP?
  8 | 
  9 | The Model Context Protocol (MCP) is an open protocol that enables AI assistants to securely access tools and data sources. By running Obelisk as an MCP server, you give AI assistants the ability to:
 10 | 
 11 | - Scan projects for security vulnerabilities
 12 | - Analyze code complexity and quality
 13 | - Track technical debt
 14 | - Audit dependencies
 15 | - Generate AI-powered health reports
 16 | 
 17 | ## Quick Start
 18 | 
 19 | ### 1. Build Obelisk with MCP Support
 20 | 
 21 | ```bash
 22 | go build -o obelisk.exe
 23 | ```
 24 | 
 25 | ### 2. Configure MCP Server
 26 | 
 27 | The MCP server has been automatically configured in your Bob IDE settings at:
 28 | `C:\Users\nikon\.bob\settings\mcp_settings.json`
 29 | 
 30 | Current configuration:
 31 | ```json
 32 | {
 33 |   "mcpServers": {
 34 |     "obelisk": {
 35 |       "command": "C:\\Users\\nikon\\Desktop\\Obelisk-CLI\\obelisk.exe",
 36 |       "args": ["mcp"],
 37 |       "env": {
 38 |         "GEMINI_API_KEY": ""
 39 |       },
 40 |       "disabled": false,
 41 |       "alwaysAllow": [],
 42 |       "disabledTools": []
 43 |     }
 44 |   }
 45 | }
 46 | ```
 47 | 
 48 | ### 3. Add Your Gemini API Key (Optional)
 49 | 
 50 | To enable AI-powered health reports, add your Gemini API key:
 51 | 
 52 | 1. Get a free API key from: https://aistudio.google.com/app/apikey
 53 | 2. Edit `mcp_settings.json` and replace the empty `GEMINI_API_KEY` value
 54 | 3. Restart Bob IDE to apply changes
 55 | 
 56 | **Note:** The MCP server works without an API key, but AI-powered features like `get_health_report` will be unavailable.
 57 | 
 58 | ### 4. Restart Bob IDE
 59 | 
 60 | After configuration, restart Bob IDE to load the Obelisk MCP server.
 61 | 
 62 | ## Available Tools
 63 | 
 64 | ### 1. `scan_project` - Comprehensive Project Analysis
 65 | 
 66 | Performs a full health scan including security, architecture, and code quality analysis.
 67 | 
 68 | **Parameters:**
 69 | - `path` (string, optional): Project directory path (defaults to current directory)
 70 | - `skip_ai` (boolean, optional): Skip AI-powered analysis (default: false)
 71 | 
 72 | **Example Usage:**
 73 | ```
 74 | "Scan the current project for issues"
 75 | "Run a full scan on C:\path\to\project"
 76 | "Scan this project without AI analysis"
 77 | ```
 78 | 
 79 | **Returns:**
 80 | - Project metadata (type, file count, directory count)
 81 | - All findings categorized by severity (critical, error, warning, info)
 82 | - Health report with grade (A-F) if AI is enabled
 83 | 
 84 | ---
 85 | 
 86 | ### 2. `check_security` - Security-Focused Scan
 87 | 
 88 | Scans specifically for security vulnerabilities including secrets, exposed credentials, and .gitignore issues.
 89 | 
 90 | **Parameters:**
 91 | - `path` (string, optional): Project directory path
 92 | 
 93 | **Example Usage:**
 94 | ```
 95 | "Check security vulnerabilities in this project"
 96 | "Scan for secrets and exposed credentials"
 97 | ```
 98 | 
 99 | **Returns:**
100 | - Security findings only
101 | - Count by severity level
102 | - Specific file locations and line numbers
103 | 
104 | ---
105 | 
106 | ### 3. `analyze_complexity` - Code Complexity Analysis
107 | 
108 | Analyzes cyclomatic complexity to identify spaghetti code and maintainability issues.
109 | 
110 | **Parameters:**
111 | - `path` (string, optional): Project directory path
112 | - `threshold` (number, optional): Complexity threshold (default: 10)
113 | 
114 | **Example Usage:**
115 | ```
116 | "Analyze code complexity in this project"
117 | "Find functions with high cyclomatic complexity"
118 | ```
119 | 
120 | **Returns:**
121 | - Complexity findings
122 | - Functions exceeding threshold
123 | - Maintainability warnings
124 | 
125 | ---
126 | 
127 | ### 4. `track_tech_debt` - Technical Debt Tracking
128 | 
129 | Tracks TODO, FIXME, HACK, and XXX comments across the codebase.
130 | 
131 | **Parameters:**
132 | - `path` (string, optional): Project directory path
133 | 
134 | **Example Usage:**
135 | ```
136 | "Track technical debt in this project"
137 | "Find all TODO and FIXME comments"
138 | ```
139 | 
140 | **Returns:**
141 | - All technical debt markers
142 | - File locations and line numbers
143 | - Total count
144 | 
145 | ---
146 | 
147 | ### 5. `audit_dependencies` - Dependency Audit
148 | 
149 | Audits package.json for deprecated, vulnerable, or unused dependencies.
150 | 
151 | **Parameters:**
152 | - `path` (string, optional): Project directory path
153 | 
154 | **Example Usage:**
155 | ```
156 | "Audit dependencies in this project"
157 | "Check for vulnerable or deprecated packages"
158 | ```
159 | 
160 | **Returns:**
161 | - Dependency findings
162 | - Vulnerable packages
163 | - Deprecated packages
164 | - Unused dependencies
165 | 
166 | ---
167 | 
168 | ### 6. `get_health_report` - AI Health Assessment
169 | 
170 | Generates an AI-powered health report with grade (A-F) and recommendations.
171 | 
172 | **Parameters:**
173 | - `path` (string, optional): Project directory path
174 | 
175 | **Requirements:**
176 | - Requires `GEMINI_API_KEY` environment variable
177 | 
178 | **Example Usage:**
179 | ```
180 | "Generate a health report for this project"
181 | "What's the overall health score of this codebase?"
182 | ```
183 | 
184 | **Returns:**
185 | - Health grade (A-F)
186 | - Overall score and category scores (security, architecture, quality)
187 | - AI-generated summary
188 | - Praise for good practices
189 | - Actionable recommendations
190 | - Top priority issues
191 | 
192 | ---
193 | 
194 | ## Available Resources
195 | 
196 | Resources provide quick access to cached scan results without re-running scans.
197 | 
198 | ### 1. `obelisk://scan/latest`
199 | Access the most recent project scan results including all findings.
200 | 
201 | ### 2. `obelisk://health/score`
202 | Current project health grade (A-F) and AI-generated recommendations.
203 | 
204 | ### 3. `obelisk://findings/security`
205 | All security-related findings from the latest scan.
206 | 
207 | ### 4. `obelisk://findings/quality`
208 | All code quality findings from the latest scan.
209 | 
210 | ### 5. `obelisk://findings/architecture`
211 | All architecture-related findings from the latest scan.
212 | 
213 | **Example Usage:**
214 | ```
215 | "Show me the latest scan results"
216 | "What's the current health score?"
217 | "List all security findings"
218 | ```
219 | 
220 | ---
221 | 
222 | ## Environment Variables
223 | 
224 | ### `GEMINI_API_KEY` (or `GOOGLE_API_KEY`)
225 | Google Gemini API key for AI-powered analysis. Required for `get_health_report` tool.
226 | 
227 | **How to set:**
228 | ```json
229 | {
230 |   "mcpServers": {
231 |     "obelisk": {
232 |       "env": {
233 |         "GEMINI_API_KEY": "your-api-key-here"
234 |       }
235 |     }
236 |   }
237 | }
238 | ```
239 | 
240 | ### `OBELISK_MODEL`
241 | AI model to use (default: `gemini-2.0-flash-exp`)
242 | 
243 | **Supported models:**
244 | - `gemini-2.0-flash-exp` (default, fastest)
245 | - `gemini-1.5-pro`
246 | - `gemini-1.5-flash`
247 | 
248 | ---
249 | 
250 | ## Supported Project Types
251 | 
252 | Currently, Obelisk MCP server supports:
253 | 
254 | - ✅ JavaScript / TypeScript
255 | - ✅ React
256 | - ✅ Next.js
257 | - 🔜 Go (Golang) - Coming Soon
258 | - 🔜 Laravel (PHP) - Coming Soon
259 | - 🔜 Python (Django/Flask) - Planned
260 | 
261 | ---
262 | 
263 | ## Testing the MCP Server
264 | 
265 | ### Manual Test (Command Line)
266 | 
267 | You can test the MCP server manually using stdio:
268 | 
269 | ```bash
270 | # Start the server
271 | .\obelisk.exe mcp
272 | 
273 | # In another terminal, send a JSON-RPC request:
274 | echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}' | .\obelisk.exe mcp
275 | ```
276 | 
277 | ### Test with Bob IDE
278 | 
279 | 1. Restart Bob IDE after configuration
280 | 2. Ask the AI assistant: "What MCP tools are available?"
281 | 3. Try a scan: "Scan the current project for issues"
282 | 4. Check results: "Show me the latest scan results"
283 | 
284 | ---
285 | 
286 | ## Troubleshooting
287 | 
288 | ### Server Not Starting
289 | 
290 | **Issue:** MCP server doesn't appear in Bob IDE
291 | 
292 | **Solutions:**
293 | 1. Verify the executable path in `mcp_settings.json` is correct
294 | 2. Ensure `disabled` is set to `false`
295 | 3. Restart Bob IDE completely
296 | 4. Check Bob IDE logs for error messages
297 | 
298 | ### Tools Not Working
299 | 
300 | **Issue:** Tools return errors or don't execute
301 | 
302 | **Solutions:**
303 | 1. Verify the project path exists
304 | 2. Check that the project is a supported type (JS/TS/React/Next.js)
305 | 3. For AI features, ensure `GEMINI_API_KEY` is set
306 | 4. Check stderr output for detailed error messages
307 | 
308 | ### API Key Issues
309 | 
310 | **Issue:** `get_health_report` fails with API key error
311 | 
312 | **Solutions:**
313 | 1. Verify API key is correctly set in `mcp_settings.json`
314 | 2. Test API key at: https://aistudio.google.com/app/apikey
315 | 3. Ensure no extra spaces or quotes in the key value
316 | 4. Restart Bob IDE after adding the key
317 | 
318 | ---
319 | 
320 | ## Advanced Configuration
321 | 
322 | ### Custom Installation Path
323 | 
324 | If you install Obelisk to a different location, update the `command` path:
325 | 
326 | ```json
327 | {
328 |   "mcpServers": {
329 |     "obelisk": {
330 |       "command": "C:\\Program Files\\Obelisk\\obelisk.exe",
331 |       "args": ["mcp"]
332 |     }
333 |   }
334 | }
335 | ```
336 | 
337 | ### Always Allow Specific Tools
338 | 
339 | To skip confirmation prompts for certain tools:
340 | 
341 | ```json
342 | {
343 |   "mcpServers": {
344 |     "obelisk": {
345 |       "alwaysAllow": ["scan_project", "check_security"]
346 |     }
347 |   }
348 | }
349 | ```
350 | 
351 | ### Disable Specific Tools
352 | 
353 | To prevent certain tools from being used:
354 | 
355 | ```json
356 | {
357 |   "mcpServers": {
358 |     "obelisk": {
359 |       "disabledTools": ["get_health_report"]
360 |     }
361 |   }
362 | }
363 | ```
364 | 
365 | ---
366 | 
367 | ## Security Considerations
368 | 
369 | ### API Key Storage
370 | 
371 | - API keys are stored in the MCP settings file
372 | - The file is located in your user directory: `~/.bob/settings/mcp_settings.json`
373 | - Ensure this file has appropriate permissions (readable only by you)
374 | - Never commit this file to version control
375 | 
376 | ### Project Access
377 | 
378 | - The MCP server can only access projects you explicitly specify
379 | - It respects `.gitignore` and `.obelisk-ignore` files
380 | - No data is sent to external services except the Gemini API (when enabled)
381 | 
382 | ### Network Communication
383 | 
384 | - The MCP server communicates via stdio (no network ports)
385 | - Only the Gemini API is contacted (when AI features are used)
386 | - All communication uses HTTPS (TLS 1.2+)
387 | 
388 | ---
389 | 
390 | ## Example Workflows
391 | 
392 | ### Daily Code Review
393 | 
394 | ```
395 | 1. "Scan the current project"
396 | 2. "Show me all critical and error findings"
397 | 3. "What's the security status?"
398 | 4. "Generate a health report"
399 | ```
400 | 
401 | ### Pre-Commit Check
402 | 
403 | ```
404 | 1. "Check security vulnerabilities"
405 | 2. "Track technical debt"
406 | 3. "Analyze code complexity"
407 | ```
408 | 
409 | ### Dependency Maintenance
410 | 
411 | ```
412 | 1. "Audit dependencies"
413 | 2. "Show me vulnerable packages"
414 | 3. "List unused dependencies"
415 | ```
416 | 
417 | ---
418 | 
419 | ## Integration with CI/CD
420 | 
421 | While the MCP server is designed for interactive use with AI assistants, you can still use Obelisk in CI/CD pipelines with the standard CLI commands:
422 | 
423 | ```bash
424 | # Headless scan for CI/CD
425 | obelisk scan --format json --strict
426 | 
427 | # This will exit with code 1 if critical issues are found
428 | ```
429 | 
430 | ---
431 | 
432 | ## Support and Feedback
433 | 
434 | - **GitHub Issues:** https://github.com/Swif7ify/Obelisk-CLI/issues
435 | - **Documentation:** https://github.com/Swif7ify/Obelisk-CLI
436 | - **Security Policy:** See SECURITY.md
437 | 
438 | ---
439 | 
440 | ## Version Information
441 | 
442 | - **Obelisk Version:** 0.1.0
443 | - **MCP Protocol Version:** 2024-11-05
444 | - **Supported Transports:** stdio
445 | 
446 | ---
447 | 
448 | **Built during the IBM Bob Hackathon (May 15–17, 2026)**  
449 | **By OneDev PH — with the help of IBM Bob IDE**
</content>
</file>
</files>
<environment_details>
# VSCode Visible Files
README.md

# VSCode Open Tabs
internal/mcp/types.go,internal/mcp/cache.go,internal/mcp/server.go,cmd/mcp.go,internal/mcp/tools.go,internal/mcp/resources.go,../../.bob/settings/mcp_settings.json,MCP_GUIDE.md,README.md

# Current Time
Current time in ISO 8601 UTC format: 2026-05-16T09:04:22.857Z
User time zone: Asia/Singapore, UTC+8:00

# Current Cost
$0.01

# Review Exclusions
The following file patterns are excluded from code reviews:
- .vscode/**

# Current Mode
<slug>code</slug>
<name>💻 Code</name>
You have not created a todo list yet. Create one with `update_todo_list` if your task is complicated or involves multiple steps.
</environment_details>

---

**Assistant:**

<attempt_completion>
<result>
Based on the code analysis, Obelisk MCP Server provides **6 powerful tools** for code analysis:

## Available MCP Tools

### 1. **scan_project** - Comprehensive Project Analysis
- Performs full health scan (security, architecture, code quality)
- Parameters: `path` (optional), `skip_ai` (optional, default: false)
- Returns: Project metadata, all findings by severity, health report with grade (A-F) if AI enabled

### 2. **check_security** - Security-Focused Scan
- Scans for vulnerabilities, secrets, exposed credentials, .gitignore issues
- Parameters: `path` (optional)
- Returns: Security findings only with severity counts and file locations

### 3. **analyze_complexity** - Code Complexity Analysis
- Analyzes cyclomatic complexity to identify spaghetti code
- Parameters: `path` (optional), `threshold` (optional, default: 10)
- Returns: Complexity findings and functions exceeding threshold

### 4. **track_tech_debt** - Technical Debt Tracking
- Tracks TODO, FIXME, HACK, and XXX comments
- Parameters: `path` (optional)
- Returns: All technical debt markers with locations

### 5. **audit_dependencies** - Dependency Audit
- Audits package.json for deprecated, vulnerable, or unused dependencies
- Parameters: `path` (optional)
- Returns: Dependency findings with vulnerability and deprecation info

### 6. **get_health_report** - AI Health Assessment
- Generates AI-powered health report with grade and recommendations
- Parameters: `path` (optional)
- **Requires**: GEMINI_API_KEY environment variable
- Returns: Health grade (A-F), scores, summary, praise, recommendations, top issues

## Additional Resources
The server also provides 5 cached resources for quick access:
- `obelisk://scan/latest` - Latest scan results
- `obelisk://health/score` - Current health grade
- `obelisk://findings/security` - Security findings
- `obelisk://findings/quality` - Quality findings
- `obelisk://findings/architecture` - Architecture findings

All tools support JavaScript, TypeScript, React, and Next.js projects.
</result>
</attempt_completion>

