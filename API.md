# 🔌 Obelisk CLI API Reference

> Complete API documentation for developers integrating with Obelisk

## 📋 Table of Contents

- [Overview](#overview)
- [Command Line Interface](#command-line-interface)
- [Go Package API](#go-package-api)
- [MCP Server API](#mcp-server-api)
- [Configuration API](#configuration-api)
- [Exit Codes](#exit-codes)
- [Error Handling](#error-handling)

---

## Overview

Obelisk provides multiple APIs for different integration scenarios:

1. **CLI API** - Command-line interface for terminal usage
2. **Go Package API** - Import Obelisk as a Go library
3. **MCP Server API** - Model Context Protocol for AI assistants
4. **Configuration API** - Programmatic configuration management

---

## Command Line Interface

### Global Flags

Available on all commands:

```bash
--api-key string      Gemini API key (overrides env var)
--model string        AI model to use (default: "gemini-2.0-flash-exp")
--no-color           Disable colored output
-v, --verbose        Enable verbose output
-h, --help           Show help
--version            Show version
```

### Commands

#### `obelisk` (Interactive Mode)

Launch interactive TUI.

```bash
obelisk [flags]
```

**Flags:** Global flags only

**Exit Codes:**

- `0` - Success
- `1` - Error

---

#### `obelisk check`

Run visual health check dashboard.

```bash
obelisk check [path] [flags]
```

**Arguments:**

- `path` (optional) - Project directory (default: current directory)

**Flags:**

- `--skip-ai` - Skip AI analysis
- `--output string` - Save report to file

**Exit Codes:**

- `0` - Success
- `1` - Error

**Example:**

```bash
obelisk check /path/to/project --skip-ai
```

---

#### `obelisk scan`

Headless scan for CI/CD pipelines.

```bash
obelisk scan [path] [flags]
```

**Arguments:**

- `path` (optional) - Project directory (default: current directory)

**Flags:**

- `--format string` - Output format: text, json, markdown (default: "text")
- `--strict` - Exit with code 1 if critical issues found
- `--skip-ai` - Skip AI analysis
- `--output string` - Save report to file

**Exit Codes:**

- `0` - No critical issues
- `1` - Critical issues found (with --strict)
- `2` - Scan error

**Example:**

```bash
obelisk scan . --format json --strict > results.json
```

---

#### `obelisk report`

Generate and export detailed report.

```bash
obelisk report [path] [flags]
```

**Arguments:**

- `path` (optional) - Project directory (default: current directory)

**Flags:**

- `--export string` - Export format: markdown, json, text (default: "markdown")
- `--skip-ai` - Skip AI analysis

**Exit Codes:**

- `0` - Success
- `1` - Error

**Example:**

```bash
obelisk report . --export=markdown
```

---

#### `obelisk config`

Manage configuration settings.

```bash
obelisk config <subcommand> [flags]
```

**Subcommands:**

- `list` - List all settings
- `get <key>` - Get setting value
- `set <key> <value>` - Set setting value

**Available Keys:**

- `api-key` - Gemini API key
- `model` - AI model name
- `report-format` - Default report format (md, txt, json)
- `default-path` - Default project path
- `no-color` - Disable colors (true/false)

**Example:**

```bash
obelisk config set api-key YOUR_API_KEY
obelisk config get model
obelisk config list
```

---

#### `obelisk protect`

Manage Git pre-push hooks.

```bash
obelisk protect <subcommand> [flags]
```

**Subcommands:**

- `install` - Install pre-push hook
- `uninstall` - Remove pre-push hook
- `status` - Check hook status

**Example:**

```bash
obelisk protect install
```

---

#### `obelisk mcp`

Start MCP server for AI assistants.

```bash
obelisk mcp [flags]
```

**Flags:**

- `--http` - Use HTTP/SSE transport (default: stdio)
- `--port int` - HTTP server port (default: 8080)

**Exit Codes:**

- `0` - Server stopped gracefully
- `1` - Server error

**Example:**

```bash
# Local stdio mode
obelisk mcp

# HTTP mode for cloud deployment
obelisk mcp --http --port 8080
```

---

#### `obelisk version`

Display version information.

```bash
obelisk version
```

**Output:**

```
Obelisk CLI
  Version:    0.1.0
  Commit:     abc1234
  Built:      2026-05-16T08:00:00Z
  Go Version: go1.26.3
```

---

## Go Package API

### Installation

```bash
go get github.com/Swif7ify/Obelisk-CLI
```

### Basic Usage

```go
package main

import (
    "fmt"
    "github.com/Swif7ify/Obelisk-CLI/internal/engine"
    "github.com/Swif7ify/Obelisk-CLI/internal/config"
)

func main() {
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        panic(err)
    }

    // Create scan engine
    eng := engine.New(cfg)

    // Run scan
    result, err := eng.Scan("/path/to/project")
    if err != nil {
        panic(err)
    }

    // Process results
    fmt.Printf("Health Score: %s\n", result.HealthScore)
    fmt.Printf("Total Findings: %d\n", len(result.Findings))
}
```

---

### Core Packages

#### `internal/engine`

Main scan orchestration engine.

```go
type Engine struct {
    Config *config.Config
}

func New(cfg *config.Config) *Engine

func (e *Engine) Scan(path string) (*scanner.ScanResult, error)
```

---

#### `internal/scanner`

Individual scanner implementations.

```go
// Scan result structure
type ScanResult struct {
    ProjectPath  string
    Framework    string
    Findings     []Finding
    HealthScore  string
    AIInsights   string
    Timestamp    time.Time
}

// Finding structure
type Finding struct {
    Type        string   // "security", "quality", "architecture"
    Severity    string   // "critical", "high", "medium", "low", "info"
    Category    string   // "secrets", "complexity", "naming", etc.
    Message     string
    File        string
    Line        int
    Suggestion  string
}
```

**Available Scanners:**

```go
// Security scanners
func ScanSecrets(path string) ([]Finding, error)
func ScanGitignore(path string) ([]Finding, error)

// Quality scanners
func ScanComplexity(path string) ([]Finding, error)
func ScanTechDebt(path string) ([]Finding, error)
func ScanNaming(path string, framework string) ([]Finding, error)

// Dependency scanners
func ScanDependencies(path string) ([]Finding, error)
func ScanImports(path string) ([]Finding, error)
```

---

#### `internal/config`

Configuration management.

```go
type Config struct {
    Model         string `json:"model"`
    DefaultPath   string `json:"default_path"`
    ReportFormat  string `json:"report_format"`
    NoColor       bool   `json:"no_color"`
}

func Load() (*Config, error)
func (c *Config) Save() error
func SetAPIKey(key string) error
func GetAPIKey() (string, error)
```

---

#### `internal/ai`

AI integration for health assessment.

```go
type Client struct {
    APIKey string
    Model  string
}

func NewClient(apiKey, model string) *Client

func (c *Client) GenerateHealthReport(result *scanner.ScanResult) (string, error)
```

---

#### `internal/report`

Report generation and formatting.

```go
func FormatMarkdown(result *scanner.ScanResult) string
func FormatJSON(result *scanner.ScanResult) ([]byte, error)
func FormatText(result *scanner.ScanResult) string

func WriteToFile(content string, path string) error
```

---

### Example: Custom Scanner

```go
package main

import (
    "fmt"
    "github.com/Swif7ify/Obelisk-CLI/internal/scanner"
)

func main() {
    // Run individual scanners
    secrets, err := scanner.ScanSecrets("./myproject")
    if err != nil {
        panic(err)
    }

    complexity, err := scanner.ScanComplexity("./myproject")
    if err != nil {
        panic(err)
    }

    // Process findings
    for _, finding := range secrets {
        if finding.Severity == "critical" {
            fmt.Printf("CRITICAL: %s in %s:%d\n",
                finding.Message, finding.File, finding.Line)
        }
    }
}
```

---

### Example: Custom Report Generator

```go
package main

import (
    "github.com/Swif7ify/Obelisk-CLI/internal/engine"
    "github.com/Swif7ify/Obelisk-CLI/internal/report"
    "github.com/Swif7ify/Obelisk-CLI/internal/config"
)

func main() {
    cfg, _ := config.Load()
    eng := engine.New(cfg)

    result, err := eng.Scan("./myproject")
    if err != nil {
        panic(err)
    }

    // Generate custom report
    markdown := report.FormatMarkdown(result)
    report.WriteToFile(markdown, "custom-report.md")

    // Or JSON
    jsonData, _ := report.FormatJSON(result)
    report.WriteToFile(string(jsonData), "report.json")
}
```

---

## MCP Server API

See [MCP_GUIDE.md](MCP_GUIDE.md) for complete MCP API documentation.

### Quick Reference

**Transport Modes:**

- `stdio` - JSON-RPC over stdin/stdout (local)
- `http` - HTTP/SSE (cloud deployment)

**Endpoints (HTTP mode):**

- `GET /health` - Health check
- `GET /sse` - SSE endpoint for MCP communication
- `GET /` - Alias for /sse

**Available Tools:**

- `scan_project` - Full project scan
- `check_security` - Security scan
- `analyze_complexity` - Complexity analysis
- `track_tech_debt` - Tech debt tracking
- `audit_dependencies` - Dependency audit
- `get_health_report` - AI health assessment

**Available Resources:**

- `obelisk://scan/latest` - Latest scan results
- `obelisk://health/score` - Health score
- `obelisk://findings/security` - Security findings
- `obelisk://findings/quality` - Quality findings
- `obelisk://findings/architecture` - Architecture findings

---

## Configuration API

### Configuration File

**Location:**

- Windows: `%USERPROFILE%\.obelisk\config.json`
- Linux/macOS: `~/.obelisk/config.json`

**Structure:**

```json
{
	"model": "gemini-2.0-flash-exp",
	"default_path": "",
	"report_format": "md",
	"no_color": false
}
```

**Note:** API keys are stored securely in OS keyring, not in config file.

---

### Environment Variables

```bash
# API Keys (priority order)
GEMINI_API_KEY=your-key-here
GOOGLE_API_KEY=your-key-here

# Configuration overrides
OBELISK_MODEL=gemini-2.0-flash-exp
OBELISK_CONFIG=/custom/path/config.json
NO_COLOR=1

# MCP Server (HTTP mode)
PORT=8080
```

---

### Programmatic Configuration

```go
import "github.com/Swif7ify/Obelisk-CLI/internal/config"

// Load config
cfg, err := config.Load()

// Modify settings
cfg.Model = "gemini-2.0-flash-exp"
cfg.ReportFormat = "json"

// Save changes
err = cfg.Save()

// Manage API key (uses OS keyring)
err = config.SetAPIKey("your-api-key")
key, err := config.GetAPIKey()
```

---

## Exit Codes

| Code | Meaning      | Usage                                                  |
| ---- | ------------ | ------------------------------------------------------ |
| `0`  | Success      | Normal operation completed                             |
| `1`  | Error        | General error or critical issues found (with --strict) |
| `2`  | Scan Error   | Failed to complete scan                                |
| `3`  | Config Error | Configuration issue                                    |
| `4`  | API Error    | AI API communication failed                            |

---

## Error Handling

### CLI Error Messages

Errors are printed to stderr with appropriate formatting:

```bash
# Error format
Error: <error message>

# With suggestions
Error: API key not found
Suggestion: Run 'obelisk config set api-key YOUR_KEY'
```

---

### Go Package Errors

All functions return standard Go errors:

```go
result, err := engine.Scan("/path")
if err != nil {
    // Handle error
    log.Fatalf("Scan failed: %v", err)
}
```

**Common Error Types:**

```go
// Configuration errors
var ErrConfigNotFound = errors.New("config file not found")
var ErrInvalidConfig = errors.New("invalid configuration")

// Scan errors
var ErrInvalidPath = errors.New("invalid project path")
var ErrNoFramework = errors.New("framework not detected")

// API errors
var ErrAPIKeyMissing = errors.New("API key not configured")
var ErrAPIFailed = errors.New("API request failed")
```

---

### MCP Error Responses

MCP errors follow JSON-RPC 2.0 format:

```json
{
	"jsonrpc": "2.0",
	"id": 1,
	"error": {
		"code": -32600,
		"message": "Invalid Request",
		"data": {
			"details": "Missing required parameter: path"
		}
	}
}
```

**Error Codes:**

- `-32700` - Parse error
- `-32600` - Invalid request
- `-32601` - Method not found
- `-32602` - Invalid params
- `-32603` - Internal error

---

## Best Practices

### CLI Integration

```bash
# Always check exit codes
obelisk scan --strict
if [ $? -ne 0 ]; then
    echo "Scan failed or found critical issues"
    exit 1
fi

# Use JSON output for parsing
obelisk scan --format json > results.json
jq '.findings[] | select(.severity=="critical")' results.json
```

---

### Go Package Integration

```go
// Always handle errors
result, err := engine.Scan(path)
if err != nil {
    return fmt.Errorf("scan failed: %w", err)
}

// Check for critical findings
hasCritical := false
for _, finding := range result.Findings {
    if finding.Severity == "critical" {
        hasCritical = true
        break
    }
}

if hasCritical {
    return errors.New("critical issues found")
}
```

---

### MCP Integration

```javascript
// Always validate responses
const response = await mcpClient.callTool("scan_project", {
	path: "/path/to/project",
});

if (response.error) {
	console.error("Scan failed:", response.error.message);
	return;
}

// Process results
const findings = response.result.findings;
```

---

## Rate Limits

### Gemini API

- **Free tier:** 15 requests per minute
- **Paid tier:** 60 requests per minute

**Recommendation:** Use `--skip-ai` flag for frequent scans, enable AI only for final reports.

---

## Support

For API questions and issues:

- 📧 Email: weareonedev@gmail.com
- 🐛 Issues: https://github.com/Swif7ify/Obelisk-CLI/issues
- 📖 Docs: https://github.com/Swif7ify/Obelisk-CLI

---

**Built by OneDev PH for the IBM Bob Hackathon 2026**
