# Obelisk CLI — VS Code Extension

> AI-Powered Automated Tech Lead for your codebase, right in your IDE.

## Features

### Side Panel with Findings Tree

A dedicated **Obelisk** panel in the activity bar shows all scan findings organized by category:

- **Security** — Hardcoded secrets, exposed API keys, missing .gitignore rules
- **Architecture** — High cyclomatic complexity, circular imports
- **Quality** — TODO/FIXME tech debt, code smells
- **Dependency** — Deprecated or vulnerable packages
- **Naming** — Convention violations (PascalCase, kebab-case, etc.)

Click any finding to **jump directly to the file and line** where the issue was found.

### Editor Integration

- **Squiggly underlines** on problematic lines (red for critical/error, yellow for warnings)
- Issues appear in the built-in **Problems panel**
- Source shown as `Obelisk` for easy filtering

### Health Summary Panel

A rich visual summary showing:
- **Letter grade** (A-F) with color coding
- **Score breakdown** — Security, Architecture, Quality bars
- **AI-generated summary** and recommendations
- **Praise notes** for good coding practices

### Status Bar

Bottom status bar shows the current health grade at a glance. Click to re-scan.

## Requirements

- **Obelisk CLI** must be installed and available on your PATH
- Install via: `powershell -ExecutionPolicy Bypass -File install.ps1`
- Or download from [GitHub Releases](https://github.com/Swif7ify/Obelisk-CLI/releases)

## Extension Settings

| Setting | Default | Description |
|---|---|---|
| `obelisk.executablePath` | `obelisk` | Path to the Obelisk CLI executable |
| `obelisk.scanOnOpen` | `true` | Auto-scan when a workspace opens |
| `obelisk.scanOnSave` | `false` | Auto re-scan on file save (debounced) |
| `obelisk.scanOnSaveDelay` | `3000` | Debounce delay in ms for save-triggered scans |
| `obelisk.skipAI` | `false` | Skip AI analysis (faster, no API key needed) |

## Commands

| Command | Description |
|---|---|
| `Obelisk: Scan Project` | Run a full health scan |
| `Obelisk: Refresh Scan` | Re-run the scan |
| `Obelisk: Clear Findings` | Clear all findings and diagnostics |
| `Obelisk: Stop Scan` | Cancel a running scan |

## How It Works

The extension runs `obelisk scan --format json` on your workspace and parses the structured output. Findings are mapped to:

1. **TreeView items** in the Obelisk side panel (grouped by category)
2. **VS Code Diagnostics** for squiggly underlines and the Problems panel
3. **Status bar badge** showing the health grade
4. **Webview summary** with scores and AI insights

---

Built by **OneDev PH** for the IBM Bob Hackathon 2026.
