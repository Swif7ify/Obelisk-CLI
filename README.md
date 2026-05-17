<p align="center">
  <h1 align="center">🏛️ Obelisk CLI</h1>
  <p align="center">
    <strong>AI-powered code health and analysis for modern codebases</strong>
  </p>
  <p align="center">
    <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="License: MIT"></a>
    <a href="https://golang.org"><img src="https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go&logoColor=white" alt="Go Version"></a>
    <a href="https://github.com/Swif7ify/Obelisk-CLI/releases"><img src="https://img.shields.io/badge/version-0.1.0-blue" alt="Version"></a>
    <a href="https://github.com/Swif7ify/Obelisk-CLI/actions"><img src="https://img.shields.io/badge/build-passing-brightgreen" alt="Build Status"></a>
  </p>
  <p align="center">
    <em>🏆 Created during the <strong>IBM Bob Hackathon</strong> (May 15–17, 2026)</em><br/>
    <em>By <strong>OneDev PH</strong> — built with the help of <strong>IBM Bob IDE</strong></em>
  </p>
</p>

---

**Obelisk** is a tool ecosystem consisting of a Go-based CLI, a VS Code extension, and an integrated Model Context Protocol (MCP) server. It evaluates project integrity, security, and architectural health using a combination of static analysis, native code parsing, and LLM-assisted review in your terminal, IDE, or AI agents.

## ⚡ Why Obelisk?

In fast-paced development, "tool fatigue" leads to missed security leaks and messy architectures. Developers forget to configure `.gitignore`, leave API keys in code, or ignore framework best practices.

**Obelisk** solves this with a single command:

```bash
obelisk
```

One command. Full project health report. AI-graded score.

---

## ✨ Features

### 💻 CLI Modes

Obelisk is designed for both humans and machines:

- **Interactive TUI (`obelisk`)** — Launch an interactive menu system with scrollable results.
- **Local Dashboard (`obelisk check`)** — Runs a visual health check with an animated spinner, displays an interactive scorecard, and automatically generates a report file in your project directory.
- **Headless Mode (`obelisk scan`)** — Completely bypasses the UI and prints standard text directly to `stdout`. Designed for automated pipelines (GitHub Actions, Jenkins). Supports raw JSON (`--format json`) and strict mode (`--strict`) which fails the pipeline with an Exit Code 1 if critical issues are found.
- **MCP Server Mode (`obelisk mcp`)** — Runs as a Model Context Protocol server, exposing code analysis capabilities to AI assistants and IDEs. Supports stdio (local)

### 🛡️ Security Shield

- **Secret Scanner** — Deep-scans files for hardcoded API keys, JWTs, AWS credentials, and private keys using regex + Shannon entropy analysis
- **Integrity Validator** — Ensures `.gitignore` and `.env` setups prevent sensitive files from reaching version control. Intelligently exempts safe files like `.env.example`.
- **Pre-Push Hook** — Integrates into Git to automatically block pushes with security vulnerabilities
- **Ignore Engine** — Automatically reads `.gitignore` and `.obelisk-ignore` to exclude files/folders across all scans

### 🧹 Architectural & Code Quality Linting

- **Native Syntax Checking (esbuild)** — If ESLint is missing, Obelisk uses its own blazing-fast embedded `esbuild` engine to natively parse `.js`, `.ts`, `.jsx`, and `.tsx` files and catch syntax errors — no Node.js required!
- **Cyclomatic Complexity Scanner** — Analyzes branching density (`if/for/switch/&&/||`) to identify maintainability risks
- **Technical Debt Tracker** — Hunts down lingering `TODO`, `FIXME`, `HACK`, and `XXX` comments across your codebase and aggregates them into architectural debt warnings
- **Naming Enforcer** — Validates file/folder naming conventions per framework (PascalCase for React, kebab-case for assets)
- **Dependency Audit** — Scans `package.json` for deprecated or vulnerable packages
- **Import Integrity** — Flags circular dependencies and enforces clean import patterns
- **Unused Dependency Detection** — Identifies and flags packages declared but never imported

### 🔐 Security Architecture

- **OS-Level Keyring Encryption** — Your Gemini API key is securely encrypted and stored natively in your Operating System's Credential Manager (Windows Credential Manager / macOS Keychain / Linux Secret Service). It is **never** saved in plaintext to disk.
- **Automatic Security Migration** — If upgrading from an older version that stored the key in `config.json`, Obelisk automatically detects it, migrates it into the secure OS keychain, and permanently scrubs the plaintext from the configuration file.
- **Zero Trust Config** — The `config.json` file contains only non-sensitive settings (model, report format, default path). Credentials are strictly isolated in the OS vault.
- **Environment Variable Support** — For CI/CD pipelines, pass keys securely at runtime via `GEMINI_API_KEY` or `GOOGLE_API_KEY` environment variables.

### 🧠 AI-Assisted Analysis

- **Health Score (A–F)** — Composite grade based on Security, Architecture, and Code Quality rubrics with strict ceiling enforcement
- **Architecture Review** — AI analyzes your project structure against common patterns such as MVC, Atomic Design, and feature-based organization
- **Technical Debt Summary** — Plain-English summary of critical issues and strong practices

### 📄 Automatic Report Generation

- **Markdown or Text** — Automatically save reports to `obelisk-report-<timestamp>.md` or configure it to save as `.txt`
- **Custom Output Paths** — Use `--output` flag to save reports to any location
- **Persistent Configuration** — Manage global settings via `~/.obelisk/config.json`

---

## 🚀 Quick Start

### Installation

#### Option 1: MSI Installer (Windows - Recommended) 🎯

**For End Users - Professional Installation Experience**

1. **Download the installer** from [GitHub Releases](https://github.com/Swif7ify/Obelisk-CLI/releases)
2. **Double-click** `ObeliskCLI-0.1.0-x64.msi`
3. **Follow the wizard:**
    - Accept the license agreement
    - Choose installation directory
    - Click Install
4. **Done!** Open a new terminal and run:
    ```powershell
    obelisk version
    ```

**Features:**

- ✅ Professional installation wizard with EULA
- ✅ Automatic PATH configuration
- ✅ Custom installation directory
- ✅ Start menu shortcuts
- ✅ Clean uninstallation via Windows Settings

**To build the MSI yourself:**

```powershell
# Requires WiX Toolset (https://wixtoolset.org/releases/)
git clone https://github.com/Swif7ify/Obelisk-CLI.git
cd Obelisk-CLI
.\installer\build-installer.ps1
```

---

#### Option 2: PowerShell Quick Install

**For Developers - Fast Command-Line Installation**

```powershell
# 1. Clone and build
git clone https://github.com/Swif7ify/Obelisk-CLI.git
cd Obelisk-CLI
make build

# 2. Run installer (adds to PATH automatically)
powershell -ExecutionPolicy Bypass -File install.ps1

# 3. Restart terminal
obelisk version
```

**To uninstall:**

```powershell
powershell -ExecutionPolicy Bypass -File uninstall.ps1
```

---

#### Option 3: Go Install

**For Go Developers**

```bash
go install github.com/Swif7ify/Obelisk-CLI@latest
```

---

#### Option 4: Portable Executable

**No Installation Required**

1. Download `obelisk.exe` from [GitHub Releases](https://github.com/Swif7ify/Obelisk-CLI/releases)
2. Run directly: `.\obelisk.exe version`

---

**📦 See [DISTRIBUTION.md](DISTRIBUTION.md) for complete distribution guide including Chocolatey, Winget, and Homebrew.**

### Usage

```bash
# 1. Interactive TUI (Highly Recommended)
# Launch the main menu to run scans and configure settings visually
obelisk

# 2. Local Dashboard (For Developers)
# Run a visual health check with an animated spinner
obelisk check
obelisk check "C:\path\to\your\project"
obelisk check --skip-ai

# 3. CI/CD Headless Mode (For Pipelines)
obelisk scan
obelisk scan --format json
obelisk scan --strict # Exit code 1 on errors
obelisk scan --skip-ai

# 4. Report Generation
obelisk report --export=markdown
obelisk report --export=json

# 5. Git Hooks
obelisk protect install
obelisk protect uninstall

# 6. Configuration Management
obelisk config list
obelisk config set api-key YOUR_API_KEY
obelisk config set report-format txt
obelisk config set auto-save false

# 7. MCP Server Mode (For AI Assistants)
obelisk mcp
```

**🤖 MCP Server Mode:** Run Obelisk as a Model Context Protocol server to expose code analysis capabilities to AI assistants. See [MCP_GUIDE.md](MCP_GUIDE.md) for complete setup and usage instructions.

**📚 Additional Documentation:**

- [API Reference](API.md) - Complete API documentation for developers
- [MCP Guide](MCP_GUIDE.md) - Model Context Protocol integration guide
- [Commands Reference](COMMANDS.md) - Detailed command documentation
- [Deployment Guide](DEPLOYMENT.md) - Cloud deployment instructions
- [Distribution Guide](DISTRIBUTION.md) - Package and release management

### Ignored Files

To completely skip scanning specific files or folders, create a `.obelisk-ignore` file at the root of your project. Obelisk natively respects both `.gitignore` and `.obelisk-ignore`.

```text
# .obelisk-ignore
*.log
vendor/
legacy_code.js
```

---

## 🏗️ Architecture

Obelisk uses a modular **Adapter Pattern** to support multiple frameworks:

```
┌─────────────┐     ┌──────────────┐     ┌──────────────┐     ┌─────────────┐
│  Detector   │────▶│   Adapter    │────▶│  Aggregator  │────▶│ Synthesizer │
│             │     │              │     │              │     │  (Gemini)   │
│ Finds what  │     │ Loads rules  │     │ Runs scans   │     │ AI grading  │
│ framework   │     │ for that     │     │ & collects   │     │ & report    │
│ you use     │     │ framework    │     │ findings     │     │ generation  │
└─────────────┘     └──────────────┘     └──────────────┘     └─────────────┘
```

### Supported Frameworks

| Framework                                | Status       |
| ---------------------------------------- | ------------ |
| JavaScript / TypeScript (React, Next.js) | ✅ Supported |
| Go (Golang)                              | Planned      |
| Laravel (PHP)                            | Planned      |
| Python (Django/Flask)                    | 🔜 Planned   |

---

## 🔐 Security Model

Obelisk takes security extremely seriously — both as a scanning tool and in its own implementation:

| Layer                 | Protection                                                                                                                                              |
| --------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **API Key Storage**   | OS-native Credential Manager encryption (Windows Credential Manager, macOS Keychain, Linux Secret Service) via `go-keyring`. Never stored in plaintext. |
| **API Communication** | HTTPS (TLS 1.2+) exclusively to Google Gemini API                                                                                                       |
| **Config File**       | `~/.obelisk/config.json` contains zero credentials — only non-sensitive preferences                                                                     |
| **Migration Safety**  | Automatic detection and secure migration of legacy plaintext keys                                                                                       |
| **CI/CD**             | Supports environment variable injection (`GEMINI_API_KEY`) for zero-persistence pipelines                                                               |
| **File Permissions**  | Config directory created with `0700`, config file with `0600` (owner-only read/write)                                                                   |

For full details, see the [Security Policy](SECURITY.md).

---

## 🛠️ Tech Stack

- **Core Logic:** [Go (Golang)](https://golang.org) — High-speed concurrency, tiny binaries, cross-platform
- **Terminal UI:** [Bubble Tea](https://github.com/charmbracelet/bubbletea) + [Lip Gloss](https://github.com/charmbracelet/lipgloss) — Interactive CLI experience
- **Native Parser:** [esbuild](https://esbuild.github.io/) — Embedded Go-native JS/TS syntax checker (no Node.js required)
- **AI Engine:** [Google Gemini API](https://ai.google.dev/) — Synthesizes findings into architectural insights
- **Credential Security:** [go-keyring](https://github.com/zalando/go-keyring) — OS-native encrypted credential storage
- **CLI Framework:** [Cobra](https://github.com/spf13/cobra) — Industry-standard Go CLI toolkit
- **Scrollable Output:** [Bubbles Viewport](https://github.com/charmbracelet/bubbles) — Scrollable TUI results panel

---

## 📁 Project Structure

```
Obelisk-CLI/
├── main.go                        # Entry point
├── cmd/                           # Cobra command definitions
│   ├── root.go                    # Root command (launches interactive TUI)
│   ├── check.go                   # Visual health check dashboard
│   ├── scan.go                    # Headless CI/CD scan mode
│   ├── config_cmd.go              # Configuration management (get/set/list)
│   ├── report.go                  # Export report command
│   ├── protect.go                 # Git pre-push hook integration
│   └── version.go                 # Version info
├── ui/                            # Bubble Tea TUI layer
│   ├── interactive.go             # Main menu state machine + viewport
│   ├── handlers.go                # Key press routing & scan orchestration
│   ├── components.go              # Reusable UI components (ScoreCard, Findings)
│   ├── views.go                   # Settings, API Key, Help, Protect views
│   ├── dashboard.go               # Standalone dashboard model
│   ├── input.go                   # Text input widget
│   ├── menu.go                    # Menu navigation widget
│   ├── spinner.go                 # Animated scanning spinner
│   └── styles.go                  # Lip Gloss style definitions
├── internal/                      # Core logic (private packages)
│   ├── config/                    # Persistent config + OS keyring integration
│   │   └── config.go              # Load, Save, SetAPIKey (keyring), GetAPIKey
│   ├── scanner/                   # All scanning modules
│   │   ├── secrets.go             # Secret scanner (regex + entropy)
│   │   ├── gitignore.go           # .gitignore validator & sensitive file tracker
│   │   ├── gitignore_matcher.go   # Merged .gitignore + .obelisk-ignore matcher
│   │   ├── naming.go              # File/folder naming convention enforcer
│   │   ├── dependencies.go        # package.json dependency auditor
│   │   ├── imports.go             # Circular dependency & import scanner
│   │   ├── complexity.go          # Cyclomatic complexity scanner
│   │   ├── techdebt.go            # Technical Debt (TODO/FIXME) tracker
│   │   └── types.go               # Finding, ScanResult type definitions
│   ├── linter/                    # Orchestrated linting
│   │   ├── eslint.go              # ESLint integration with auto-fallback
│   │   └── native.go              # Native esbuild-powered syntax checker
│   ├── detector/                  # Framework detection
│   │   └── detector.go            # Identifies project type from signature files
│   ├── ai/                        # Gemini API integration
│   │   ├── client.go              # Gemini API client
│   │   ├── prompt.go              # AI prompt builder with grading rubric
│   │   └── report.go              # Health report parser + fallback scoring
│   ├── engine/                    # Scan orchestrator
│   │   └── engine.go              # Coordinates all scanners + AI pipeline
│   └── report/                    # Report generation
│       ├── formatter.go           # Markdown/Text report formatters
│       └── writer.go              # File output writer
└── adapters/                      # Framework-specific naming rules
    ├── javascript.go              # JS/TS/React/Next.js rules
    ├── laravel.go                 # PHP/Laravel rules (planned)
    └── golang.go                  # Go rules (planned)
```

---

## 🤝 Contributing

Contributions are welcome! Please read our [Contributing Guide](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## 📜 License

This project is licensed under the MIT License — see the [LICENSE](LICENSE) file for details.

## 🔒 Security

For security concerns, please see our [Security Policy](SECURITY.md).

---

<p align="center">
  🏆 Built during the <strong>IBM Bob Hackathon</strong> (May 15–17, 2026)<br/>
  By <strong>OneDev PH</strong> — with the help of <strong>IBM Bob IDE</strong><br/><br/>
  <em>Powered by Go, Google Gemini AI & IBM Bob</em>
</p>
