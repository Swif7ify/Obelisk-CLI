# 🏗️ Obelisk CLI Architecture

> Deep dive into the internal architecture and design patterns

## 📋 Table of Contents

- [Overview](#overview)
- [Design Principles](#design-principles)
- [System Architecture](#system-architecture)
- [Core Components](#core-components)
- [Data Flow](#data-flow)
- [Framework Adapter Pattern](#framework-adapter-pattern)
- [Security Architecture](#security-architecture)
- [Performance Considerations](#performance-considerations)
- [Extension Points](#extension-points)

---

## Overview

Obelisk CLI is built with a modular, extensible architecture that separates concerns and enables easy addition of new features and framework support.

### Key Characteristics

- **Modular Design** - Independent, composable components
- **Adapter Pattern** - Framework-agnostic core with pluggable adapters
- **Concurrent Scanning** - Parallel execution of independent scanners
- **Secure by Default** - OS-level credential encryption, no plaintext secrets
- **CLI-First** - Optimized for both interactive and headless usage

---

## Design Principles

### 1. Separation of Concerns

Each package has a single, well-defined responsibility:

```
cmd/          → CLI interface and command routing
ui/           → Terminal UI and user interaction
internal/     → Core business logic (private)
  ├── engine/     → Scan orchestration
  ├── scanner/    → Individual scanners
  ├── detector/   → Framework detection
  ├── ai/         → AI integration
  ├── config/     → Configuration management
  ├── report/     → Report generation
  └── mcp/        → MCP server implementation
adapters/     → Framework-specific rules
```

### 2. Dependency Inversion

High-level modules don't depend on low-level modules. Both depend on abstractions:

```go
// Adapter interface (abstraction)
type Adapter interface {
    GetNamingRules() []NamingRule
    GetIgnorePatterns() []string
}

// Engine depends on interface, not concrete implementations
type Engine struct {
    adapter Adapter
}
```

### 3. Fail-Safe Defaults

- Missing config → Use sensible defaults
- Missing API key → Skip AI analysis, continue with local scans
- Scanner error → Log and continue with other scanners
- Network error → Graceful degradation

### 4. Immutability

Configuration and scan results are immutable after creation:

```go
type ScanResult struct {
    // All fields are read-only after creation
    ProjectPath  string
    Findings     []Finding  // Slice is copied, not modified
    Timestamp    time.Time
}
```

---

## System Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                         CLI Layer                            │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐   │
│  │   Root   │  │  Check   │  │   Scan   │  │   MCP    │   │
│  │ (TUI)    │  │(Dashboard)│  │(Headless)│  │ (Server) │   │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘   │
└───────┼─────────────┼─────────────┼─────────────┼──────────┘
        │             │             │             │
        └─────────────┴─────────────┴─────────────┘
                      │
        ┌─────────────▼─────────────┐
        │      Engine Layer          │
        │  ┌──────────────────────┐  │
        │  │  Scan Orchestrator   │  │
        │  └──────────┬───────────┘  │
        └─────────────┼───────────────┘
                      │
        ┌─────────────▼─────────────┐
        │     Scanner Layer          │
        │  ┌────────┐  ┌──────────┐ │
        │  │Security│  │ Quality  │ │
        │  │Scanners│  │ Scanners │ │
        │  └────────┘  └──────────┘ │
        └────────────────────────────┘
                      │
        ┌─────────────▼─────────────┐
        │    Integration Layer       │
        │  ┌────────┐  ┌──────────┐ │
        │  │   AI   │  │ Adapters │ │
        │  │(Gemini)│  │(Framework)│ │
        │  └────────┘  └──────────┘ │
        └────────────────────────────┘
```

---

## Core Components

### 1. Command Layer (`cmd/`)

**Responsibility:** CLI interface, argument parsing, command routing

**Key Files:**

- `root.go` - Root command, launches TUI
- `check.go` - Visual health check dashboard
- `scan.go` - Headless CI/CD mode
- `mcp.go` - MCP server mode
- `config_cmd.go` - Configuration management
- `protect.go` - Git hook management

**Pattern:** Command pattern with Cobra framework

```go
var scanCmd = &cobra.Command{
    Use:   "scan [path]",
    Short: "Run headless scan",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Parse args, load config, run scan
        return runScan(args)
    },
}
```

---

### 2. UI Layer (`ui/`)

**Responsibility:** Terminal user interface, user interaction

**Key Files:**

- `interactive.go` - Main TUI state machine
- `dashboard.go` - Health check dashboard
- `components.go` - Reusable UI components
- `styles.go` - Visual styling with Lip Gloss

**Pattern:** Model-View-Update (Elm Architecture) via Bubble Tea

```go
type Model struct {
    state    string
    findings []Finding
    // ... other state
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd)
func (m Model) View() string
```

---

### 3. Engine Layer (`internal/engine/`)

**Responsibility:** Scan orchestration, result aggregation

**Key Concepts:**

1. **Concurrent Execution** - Scanners run in parallel
2. **Error Isolation** - One scanner failure doesn't stop others
3. **Result Aggregation** - Combines findings from all scanners

```go
type Engine struct {
    config  *config.Config
    adapter adapters.Adapter
}

func (e *Engine) Scan(path string) (*scanner.ScanResult, error) {
    // 1. Detect framework
    framework := detector.Detect(path)

    // 2. Load adapter
    adapter := adapters.Get(framework)

    // 3. Run scanners concurrently
    results := e.runScanners(path, adapter)

    // 4. Aggregate findings
    return e.aggregate(results), nil
}
```

---

### 4. Scanner Layer (`internal/scanner/`)

**Responsibility:** Individual code analysis scanners

**Scanner Types:**

**Security Scanners:**

- `secrets.go` - Hardcoded credentials detection
- `gitignore.go` - .gitignore validation

**Quality Scanners:**

- `complexity.go` - Cyclomatic complexity analysis
- `techdebt.go` - TODO/FIXME tracking
- `naming.go` - Naming convention enforcement

**Dependency Scanners:**

- `dependencies.go` - Package.json audit
- `imports.go` - Import analysis

**Pattern:** Strategy pattern - each scanner implements common interface

```go
type Scanner interface {
    Scan(path string) ([]Finding, error)
}

type Finding struct {
    Type       string
    Severity   string
    Category   string
    Message    string
    File       string
    Line       int
    Suggestion string
}
```

---

### 5. Detector Layer (`internal/detector/`)

**Responsibility:** Framework detection from project structure

**Detection Strategy:**

1. Check for signature files (package.json, go.mod, composer.json)
2. Analyze file structure patterns
3. Return framework identifier

```go
func Detect(path string) string {
    if hasFile(path, "package.json") {
        return "javascript"
    }
    if hasFile(path, "go.mod") {
        return "golang"
    }
    if hasFile(path, "composer.json") {
        return "laravel"
    }
    return "unknown"
}
```

---

### 6. Adapter Layer (`adapters/`)

**Responsibility:** Framework-specific rules and patterns

**Adapter Interface:**

```go
type Adapter interface {
    GetNamingRules() []NamingRule
    GetIgnorePatterns() []string
    GetComplexityThreshold() int
}

type NamingRule struct {
    Pattern     string
    Convention  string  // "PascalCase", "camelCase", "kebab-case"
    FileType    string  // "component", "utility", "config"
}
```

**Example: JavaScript Adapter**

```go
type JavaScriptAdapter struct{}

func (a *JavaScriptAdapter) GetNamingRules() []NamingRule {
    return []NamingRule{
        {Pattern: "*.jsx", Convention: "PascalCase", FileType: "component"},
        {Pattern: "*.js", Convention: "camelCase", FileType: "utility"},
        {Pattern: "*.config.js", Convention: "kebab-case", FileType: "config"},
    }
}
```

---

### 7. AI Layer (`internal/ai/`)

**Responsibility:** AI integration for health assessment

**Components:**

1. **Client** - Gemini API communication
2. **Prompt Builder** - Constructs AI prompts with context
3. **Report Parser** - Extracts structured data from AI response

```go
type Client struct {
    apiKey string
    model  string
}

func (c *Client) GenerateHealthReport(result *ScanResult) (string, error) {
    // 1. Build prompt with findings summary
    prompt := buildPrompt(result)

    // 2. Call Gemini API
    response := c.callAPI(prompt)

    // 3. Parse and validate response
    return parseReport(response), nil
}
```

---

### 8. Config Layer (`internal/config/`)

**Responsibility:** Configuration management, credential security

**Security Features:**

1. **OS Keyring Integration** - API keys stored in OS credential manager
2. **Automatic Migration** - Detects and migrates legacy plaintext keys
3. **Zero-Credential Config** - Config file contains no secrets

```go
type Config struct {
    APIKey       string `json:"-"`  // Never serialized
    Model        string `json:"model"`
    ReportFormat string `json:"report_format"`
}

func (c *Config) SetAPIKey(key string) error {
    // Store in OS keyring, not config file
    return keyring.Set("obelisk-cli", "gemini-api-key", key)
}

func (c *Config) GetAPIKey() (string, error) {
    // Retrieve from OS keyring
    return keyring.Get("obelisk-cli", "gemini-api-key")
}
```

---

### 9. MCP Layer (`internal/mcp/`)

**Responsibility:** Model Context Protocol server implementation

**Components:**

1. **Server** - MCP protocol handler
2. **Tools** - Exposed analysis functions
3. **Resources** - Cached scan results
4. **HTTP Server** - HTTP/SSE transport for cloud deployment

```go
type Server struct {
    tools     map[string]Tool
    resources map[string]Resource
    cache     *Cache
}

type Tool struct {
    Name        string
    Description string
    Handler     func(params map[string]interface{}) (interface{}, error)
}
```

---

## Data Flow

### Scan Flow

```
1. User Input
   ↓
2. Command Parser (Cobra)
   ↓
3. Config Loader
   ↓
4. Framework Detector
   ↓
5. Adapter Loader
   ↓
6. Engine.Scan()
   ↓
7. Parallel Scanner Execution
   ├─→ Security Scanners
   ├─→ Quality Scanners
   └─→ Dependency Scanners
   ↓
8. Result Aggregation
   ↓
9. AI Analysis (optional)
   ↓
10. Report Generation
    ↓
11. Output (TUI/JSON/Markdown)
```

### Configuration Flow

```
1. Load Config File
   ↓
2. Check for Legacy API Key
   ↓
3. If found → Migrate to Keyring
   ↓
4. Load API Key from Keyring
   ↓
5. Check Environment Variables
   ↓
6. Merge with Command Flags
   ↓
7. Return Final Config
```

---

## Framework Adapter Pattern

### Why Adapters?

Different frameworks have different conventions:

- **React:** PascalCase components, camelCase utilities
- **Laravel:** kebab-case views, PascalCase models
- **Go:** PascalCase exports, camelCase private

### Adapter Registration

```go
var adapters = map[string]Adapter{
    "javascript": &JavaScriptAdapter{},
    "golang":     &GolangAdapter{},
    "laravel":    &LaravelAdapter{},
}

func Get(framework string) Adapter {
    if adapter, ok := adapters[framework]; ok {
        return adapter
    }
    return &DefaultAdapter{}
}
```

### Adding New Adapter

1. Create new file in `adapters/`
2. Implement `Adapter` interface
3. Register in adapter map
4. Add detection logic in `detector/`

---

## Security Architecture

### Credential Storage

```
┌─────────────────────────────────────┐
│         User Input                   │
│    (obelisk config set api-key)     │
└──────────────┬──────────────────────┘
               ↓
┌──────────────▼──────────────────────┐
│      Config Layer                    │
│   (validates and sanitizes)         │
└──────────────┬──────────────────────┘
               ↓
┌──────────────▼──────────────────────┐
│      go-keyring Library              │
│   (OS-native encryption)            │
└──────────────┬──────────────────────┘
               ↓
┌──────────────▼──────────────────────┐
│   OS Credential Manager              │
│  • Windows: Credential Manager      │
│  • macOS: Keychain                  │
│  • Linux: Secret Service            │
└─────────────────────────────────────┘
```

### API Communication

```
┌─────────────────────────────────────┐
│      Obelisk CLI                     │
└──────────────┬──────────────────────┘
               ↓ HTTPS (TLS 1.2+)
┌──────────────▼──────────────────────┐
│      Google Gemini API               │
│   (receives only summaries)         │
└─────────────────────────────────────┘
```

---

## Performance Considerations

### Concurrent Scanning

Scanners run in parallel using goroutines:

```go
func (e *Engine) runScanners(path string) []Finding {
    var wg sync.WaitGroup
    results := make(chan []Finding, len(scanners))

    for _, scanner := range scanners {
        wg.Add(1)
        go func(s Scanner) {
            defer wg.Done()
            findings, _ := s.Scan(path)
            results <- findings
        }(scanner)
    }

    wg.Wait()
    close(results)

    return aggregate(results)
}
```

### Caching

MCP server caches scan results:

```go
type Cache struct {
    mu      sync.RWMutex
    results map[string]*ScanResult
    ttl     time.Duration
}

func (c *Cache) Get(key string) (*ScanResult, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()

    result, ok := c.results[key]
    if !ok || time.Since(result.Timestamp) > c.ttl {
        return nil, false
    }
    return result, true
}
```

### File Scanning Optimization

- **Gitignore Respect** - Skip ignored files early
- **Binary Detection** - Skip binary files
- **Size Limits** - Skip files > 1MB
- **Parallel Processing** - Process multiple files concurrently

---

## Extension Points

### Adding a New Scanner

1. Create scanner file in `internal/scanner/`
2. Implement scanner logic
3. Register in engine
4. Add tests

```go
// internal/scanner/myscan.go
func ScanMyFeature(path string) ([]Finding, error) {
    var findings []Finding
    // Scan logic here
    return findings, nil
}

// internal/engine/engine.go
func (e *Engine) Scan(path string) (*ScanResult, error) {
    // Add to scanner list
    myFindings, _ := scanner.ScanMyFeature(path)
    allFindings = append(allFindings, myFindings...)
}
```

### Adding a New Command

1. Create command file in `cmd/`
2. Define Cobra command
3. Register in `init()`

```go
// cmd/mycommand.go
var myCmd = &cobra.Command{
    Use:   "mycommand",
    Short: "My new command",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Command logic
        return nil
    },
}

func init() {
    rootCmd.AddCommand(myCmd)
}
```

### Adding a New MCP Tool

1. Define tool in `internal/mcp/tools.go`
2. Register in server initialization

```go
func (s *Server) registerTools() {
    s.tools["my_tool"] = Tool{
        Name:        "my_tool",
        Description: "My custom tool",
        Handler: func(params map[string]interface{}) (interface{}, error) {
            // Tool logic
            return result, nil
        },
    }
}
```

---

## Testing Strategy

### Unit Tests

Each package has its own test file:

```
internal/scanner/
  ├── secrets.go
  ├── secrets_test.go
  ├── complexity.go
  └── complexity_test.go
```

### Integration Tests

Test complete scan flow:

```go
func TestFullScan(t *testing.T) {
    cfg := config.Default()
    eng := engine.New(cfg)
    result, err := eng.Scan("./testdata/sample-project")

    assert.NoError(t, err)
    assert.NotEmpty(t, result.Findings)
}
```

### Test Coverage

Run with coverage:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## Future Architecture Improvements

### Planned Enhancements

1. **Plugin System** - Load external scanners dynamically
2. **Distributed Scanning** - Scan large projects across multiple machines
3. **Real-time Monitoring** - Watch mode for continuous scanning
4. **Custom Rules Engine** - User-defined scanning rules
5. **Database Backend** - Store scan history and trends

---

## Resources

- **Go Best Practices:** https://golang.org/doc/effective_go
- **Cobra Documentation:** https://github.com/spf13/cobra
- **Bubble Tea Guide:** https://github.com/charmbracelet/bubbletea
- **MCP Specification:** https://modelcontextprotocol.io/

---

**Built with ❤️ by OneDev PH for the IBM Bob Hackathon 2026**
