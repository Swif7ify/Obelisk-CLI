<p align="center">
  <h1 align="center">🏛️ Obelisk CLI</h1>
  <p align="center">
    <strong>The AI-Powered "Automated Tech Lead" for Modern Codebases</strong>
  </p>
  <p align="center">
    <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="License: MIT"></a>
    <a href="https://golang.org"><img src="https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go&logoColor=white" alt="Go Version"></a>
    <a href="https://github.com/Swif7ify/Obelisk-CLI/releases"><img src="https://img.shields.io/badge/version-0.1.0-blue" alt="Version"></a>
    <a href="https://github.com/Swif7ify/Obelisk-CLI/actions"><img src="https://img.shields.io/badge/build-passing-brightgreen" alt="Build Status"></a>
  </p>
</p>

---

**Obelisk** is a high-performance CLI tool built in **Go** that acts as a final gatekeeper for your project. It doesn't just check for syntax errors — it evaluates project integrity, security, and architectural health using a combination of static analysis and LLM intelligence.

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
- **Interactive TUI (`obelisk`)** — Launch a premium, visual menu system.
- **Local Dashboard (`obelisk check`)** — Runs a visual health check with an animated spinner, displays an interactive scorecard, and automatically generates a Markdown report file in your project directory upon exiting.
- **Headless Mode (`obelisk scan`)** — Completely bypasses the UI and prints standard text directly to `stdout`. Designed for automated pipelines (GitHub Actions, Jenkins). Supports raw JSON (`--format json`) and strict mode (`--strict`) which fails the pipeline with an Exit Code 1 if critical issues are found.

### 🛡️ Security Shield
- **Secret Scanner** — Deep-scans files for hardcoded API keys, JWTs, AWS credentials, and private keys using regex + entropy analysis
- **Integrity Validator** — Ensures `.gitignore` and `.env` setups prevent sensitive files from reaching version control
- **Pre-Push Hook** — Integrates into Git to automatically block pushes with security vulnerabilities
- **Ignore Engine** — Automatically reads `.gitignore` and `.obelisk-ignore` to exclude files/folders across all scans.

### 🧹 Architectural & Code Quality Linting
- **Native Syntax Checking** — If ESLint is missing, Obelisk uses its own blazing-fast embedded `esbuild` engine to natively parse JS/TS files and catch syntax errors.
- **Cyclomatic Complexity Scanner** — Mathematically analyzes your code's branching density (`if/for/switch`) to detect and flag highly unmaintainable "Spaghetti Code".
- **Technical Debt Tracker** — Hunts down lingering `TODO`, `FIXME`, and `HACK` comments across your codebase and aggregates them into architectural debt warnings.
- **Naming Enforcer** — Validates file/folder naming conventions per framework (PascalCase for React, kebab-case for assets)
- **Dependency Audit** — Scans `package.json` for deprecated or vulnerable packages
- **Import Integrity** — Flags circular dependencies and enforces clean import patterns
- **Unused Dependency Detection** — Identifies and flags packages declared but never imported

### 🧠 AI "Senior Dev" Brain
- **Health Score (A–F)** — Composite grade based on Security, Architecture, and Code Quality rubrics
- **Vibe Check** — AI analyzes your project structure against industry-standard patterns
- **Technical Debt Summary** — Plain-English summary of critical issues + praise for good practices

### 📄 Automatic Report Generation
- **Markdown or Text** — Automatically save reports to `obelisk-report-<timestamp>.md` or configure it to save as `.txt`.
- **Custom Output Paths** — Use `--output` flag to save reports to any location
- **Persistent Configuration** — Manage global settings via `~/.obelisk/config.json`

---

## 🚀 Quick Start

### Installation

```bash
# Install from source
go install github.com/Swif7ify/Obelisk-CLI@latest

# Or clone and build
git clone https://github.com/Swif7ify/Obelisk-CLI.git
cd Obelisk-CLI
make build
```

### Usage

```bash
# 1. Interactive TUI (Highly Recommended)
# Launch the main menu to run scans and configure settings visually
obelisk

# 2. Local Dashboard (For Developers)
# Run a visual health check with an animated spinner
obelisk check
obelisk check "C:\path\to\your\project"

# 3. CI/CD Headless Mode (For Pipelines)
obelisk scan
obelisk scan --format json
obelisk scan --strict # Exit code 1 on errors

# 3. Report Generation
obelisk report --export=markdown
obelisk report --export=json

# 4. Git Hooks
obelisk protect --install

# 5. Configuration Management
obelisk config list
obelisk config set api-key YOUR_API_KEY
obelisk config set report-format txt
```

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

| Framework | Status |
|-----------|--------|
| JavaScript / TypeScript (React, Next.js) | ✅ Supported |
| Go (Golang) | 🔜 Coming Soon |
| Laravel (PHP) | 🔜 Coming Soon |
| Python (Django/Flask) | 🔜 Planned |

---

## 🛠️ Tech Stack

- **Core Logic:** [Go (Golang)](https://golang.org) — High-speed concurrency, tiny binaries, cross-platform
- **Terminal UI:** [Bubble Tea](https://github.com/charmbracelet/bubbletea) + [Lip Gloss](https://github.com/charmbracelet/lipgloss) — Premium interactive CLI experience
- **AI Engine:** [Google Gemini API](https://ai.google.dev/) — Synthesizes findings into architectural insights
- **CLI Framework:** [Cobra](https://github.com/spf13/cobra) — Industry-standard Go CLI toolkit

---

## 📁 Project Structure

```
Obelisk-CLI/
├── main.go                 # Entry point
├── cmd/                    # Cobra command definitions
│   ├── root.go             # Root command (Launches TUI)
│   ├── check.go            # Original health check command
│   ├── scan.go             # Headless scan mode
│   ├── config_cmd.go       # Config management
│   ├── report.go           # Export report command
│   ├── protect.go          # Git hook integration
│   └── version.go          # Version info
├── ui/                     # Bubble Tea TUI layer
│   ├── interactive.go      # Main menu state machine
│   ├── handlers.go         # Key press routing
│   ├── components.go       # Reusable views (Input, Spinner)
│   └── styles.go           # Lip Gloss styling
├── internal/               # Core logic
│   ├── config/             # Persistent JSON settings
│   ├── scanner/            # All scanning modules (.obelisk-ignore)
│   ├── detector/           # Framework detection
│   ├── ai/                 # Gemini API integration
│   └── engine/             # Orchestrator
└── adapters/               # Framework-specific rules
    ├── javascript.go       # JS/TS/React/Next.js
    ├── laravel.go          # PHP/Laravel
    └── golang.go           # Go
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
  Built using IBM BOB and Go for the hackathon
</p>
