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
obelisk check
```

One command. Full project health report. AI-graded score.

---

## ✨ Features

### 🛡️ Security Shield
- **Secret Scanner** — Deep-scans files for hardcoded API keys, JWTs, AWS credentials, and private keys using regex + entropy analysis
- **Integrity Validator** — Ensures `.gitignore` and `.env` setups prevent sensitive files from reaching version control
- **Pre-Push Hook** — Integrates into Git to automatically block pushes with security vulnerabilities

### 🧹 Architectural Linting
- **Naming Enforcer** — Validates file/folder naming conventions per framework (PascalCase for React, kebab-case for assets)
- **Dependency Audit** — Scans `package.json` for deprecated or vulnerable packages
- **Import Integrity** — Flags circular dependencies and enforces clean import patterns
- **Unused Dependency Detection** — Identifies and flags packages declared but never imported

### 🧠 AI "Senior Dev" Brain
- **Health Score (A–F)** — Composite grade based on Security, Architecture, and Code Quality rubrics
- **Vibe Check** — AI analyzes your project structure against industry-standard patterns
- **Technical Debt Summary** — Plain-English summary of critical issues + praise for good practices

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

### Setup

```bash
# Set your Gemini API key
export GOOGLE_API_KEY="your-api-key-here"
```

### Usage

```bash
# Run a full health check with interactive TUI
obelisk check

# Run with an API key flag
obelisk check --api-key="your-key"

# Export a health report as Markdown
obelisk report --export=markdown

# Export as JSON
obelisk report --export=json

# Install as a Git pre-push hook
obelisk protect --install

# Run in strict mode (non-zero exit on failures)
obelisk protect --strict
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
| Laravel (PHP) | 🔜 Coming Soon |
| Go | 🔜 Coming Soon |
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
│   ├── root.go             # Root command & global flags
│   ├── check.go            # Health check command (TUI)
│   ├── report.go           # Export report command
│   ├── protect.go          # Git hook integration
│   └── version.go          # Version info
├── ui/                     # Bubble Tea TUI layer
│   ├── dashboard.go        # Main dashboard model
│   ├── spinner.go          # Animated scan progress
│   ├── styles.go           # Lip Gloss styling
│   └── components.go       # Reusable view components
├── internal/               # Core logic
│   ├── scanner/            # All scanning modules
│   │   ├── secrets.go      # Secret/credential detection
│   │   ├── gitignore.go    # .gitignore validation
│   │   ├── dependencies.go # Dependency auditing
│   │   ├── naming.go       # Naming convention checks
│   │   ├── imports.go      # Import/circular dep analysis
│   │   └── types.go        # Shared types
│   ├── detector/           # Framework detection
│   ├── ai/                 # Gemini API integration
│   └── engine/             # Orchestrator
└── adapters/               # Framework-specific rules
    ├── javascript.go       # JS/TS/React/Next.js
    ├── laravel.go          # PHP/Laravel (stub)
    └── golang.go           # Go (stub)
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
