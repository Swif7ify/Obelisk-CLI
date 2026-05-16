# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.0] - 2026-05-17

> 🚀 **Official Release** — The AI-Powered Automated Tech Lead.

    ### Added
    - Complete UI redesign for the interactive TUI with a premium aesthetic
    - Polished component architecture and scorecard rendering
    - Streamlined MCP server documentation focusing exclusively on local stdio integration

## [0.1.0] - 2026-05-17

> 🏆 Initial Pre-release — Built during the **IBM Bob Hackathon** (May 15–17, 2026) by **OneDev PH**.

### Added

#### CLI Modes
- Interactive TUI with Bubble Tea — premium visual menu system
- `obelisk check` — Visual health check dashboard with animated spinner
- `obelisk scan` — Headless CI/CD mode with `--format json` and `--strict` flags
- `obelisk report` — Export findings as Markdown or JSON
- `obelisk protect` — Git pre-push hook installation
- `obelisk config` — Persistent configuration management (list/get/set)
- `obelisk version` — Version info display
- Scrollable viewport for long scan results in the TUI

#### Security Scanners
- Secret scanner with regex + Shannon entropy analysis
- `.gitignore` and `.env` integrity validation
- Sensitive file tracking — detects files tracked by Git that shouldn't be (with `.example`/`.template` exemptions)
- Hardcoded credential detection (AWS, GitHub, Slack, Stripe, Google, JWT, SSH keys)

#### Code Quality Scanners
- Native esbuild-powered JS/TS syntax checker (no Node.js required)
- ESLint orchestration with automatic native fallback
- Cyclomatic Complexity (Spaghetti Code) scanner
- Technical Debt (TODO/FIXME/HACK/XXX) tracker
- File/folder naming convention enforcer
- Dependency audit for `package.json`
- Circular dependency and import integrity scanner
- Unused dependency detection

#### AI Integration
- Google Gemini API integration for architectural health scoring
- Health Score grading (A–F) with strict hard-ceiling enforcement
- Directory structure "Vibe Check" analysis
- Constructive praise and actionable recommendations

#### Security Architecture
- OS-level API key encryption via `go-keyring` (Windows Credential Manager, macOS Keychain, Linux Secret Service)
- Automatic plaintext key migration from legacy `config.json`
- Zero-credential config files (`json:"-"` struct tags)
- Strict file permissions (`0700` dirs, `0600` files)

#### Framework Support
- JavaScript / TypeScript (React, Next.js) — ✅ Supported
- Merged `.gitignore` + `.obelisk-ignore` ignore engine

#### Documentation
- README with full feature documentation
- SECURITY.md with credential architecture details
- CONTRIBUTING.md with development setup guide
- CODE_OF_CONDUCT.md
- MIT License
