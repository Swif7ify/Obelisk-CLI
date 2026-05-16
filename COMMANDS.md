# Obelisk CLI Commands Reference

## Quick Reference

```bash
# Show version
obelisk --version
obelisk version

# Show help
obelisk --help
obelisk [command] --help

# Interactive mode (TUI)
obelisk

# Run health check
obelisk check /path/to/project

# Headless scan (CI/CD)
obelisk scan /path/to/project

# Generate report
obelisk report /path/to/project --output report.json

# Configure settings
obelisk config set api-key YOUR_API_KEY
obelisk config get api-key
obelisk config list

# Git pre-push hook
obelisk protect install
obelisk protect uninstall

# Check for updates
obelisk update
```

---

## Available Commands

### 1. **Interactive Mode** (Default)

Launch the interactive Terminal UI when no subcommand is provided.

```bash
obelisk
```

**Features:**

- Visual dashboard with real-time scanning
- Interactive project selection
- Color-coded health metrics
- Progress indicators

**Flags:**

- `--api-key string` - Gemini API key (overrides GOOGLE_API_KEY env var)
- `--model string` - AI model to use (default: "gemini-2.5-flash")
- `--no-color` - Disable colored output
- `-v, --verbose` - Enable verbose output

---

### 2. **check** - Full Health Check

Run a comprehensive health check on your project with detailed output.

```bash
obelisk check --path /path/to/project
```

**Flags:**

- `--path string` - Path to project directory (required)
- `--skip-ai` - Skip AI analysis (faster, local-only checks)
- `--api-key string` - Gemini API key
- `--model string` - AI model to use

**Output:**

- Security vulnerabilities
- Code quality issues
- Architectural problems
- Dependency analysis
- Overall health score

---

### 3. **scan** - Headless Check (CI/CD)

Run a headless health check optimized for CI/CD pipelines.

```bash
obelisk scan --path /path/to/project
```

**Flags:**

- `--path string` - Path to project directory (required)
- `--skip-ai` - Skip AI analysis
- `--fail-on-critical` - Exit with code 1 if critical issues found
- `--output string` - Output format: text, json, or markdown

**Exit Codes:**

- `0` - No critical issues
- `1` - Critical issues found (with --fail-on-critical)
- `2` - Scan error

**Example CI/CD Usage:**

```yaml
# GitHub Actions
- name: Run Obelisk Scan
  run: |
      obelisk scan --path . --fail-on-critical --output json > scan-results.json
```

---

### 4. **report** - Generate Report

Generate and export a detailed health report.

```bash
obelisk report --path /path/to/project --output report.json
```

**Flags:**

- `--path string` - Path to project directory (required)
- `--output string` - Output file path (required)
- `--format string` - Report format: json, markdown, html (default: json)
- `--skip-ai` - Skip AI analysis

**Supported Formats:**

- **JSON** - Machine-readable, ideal for automation
- **Markdown** - Human-readable, great for documentation
- **HTML** - Rich formatting, perfect for sharing

---

### 5. **config** - Manage Configuration

Manage Obelisk configuration settings.

```bash
# Set API key
obelisk config set api-key YOUR_API_KEY

# Get API key
obelisk config get api-key

# List all settings
obelisk config list

# Set AI model
obelisk config set model gemini-2.5-flash

# Enable/disable AI
obelisk config set skip-ai true
```

**Available Settings:**

- `api-key` - Gemini API key
- `model` - AI model to use
- `skip-ai` - Skip AI analysis (true/false)
- `no-color` - Disable colored output (true/false)

**Config File Location:**

- Windows: `%USERPROFILE%\.obelisk\config.yaml`
- Linux/macOS: `~/.obelisk/config.yaml`

---

### 6. **protect** - Git Pre-Push Hook

Install/uninstall Git pre-push hook for automatic checks.

```bash
# Install hook
obelisk protect install

# Uninstall hook
obelisk protect uninstall

# Check hook status
obelisk protect status
```

**What it does:**

- Runs `obelisk scan` before every `git push`
- Blocks push if critical issues are found
- Provides immediate feedback on code quality

**Hook Location:**
`.git/hooks/pre-push`

---

### 7. **update** - Check for Updates

Check for and install Obelisk CLI updates.

```bash
# Check for updates
obelisk update

# Check without installing
obelisk update --check-only
```

**Features:**

- Automatic version checking
- Secure download from GitHub Releases
- SHA256 checksum verification
- Backup of current version

---

### 8. **version** - Show Version

Display version information.

```bash
obelisk version
# or
obelisk --version
```

**Output:**

```
Obelisk CLI
  Version:    0.1.0
  Commit:     abc1234
  Built:      2026-05-16T08:00:00Z
```

---

### 9. **completion** - Shell Completion

Generate shell completion scripts.

```bash
# Bash
obelisk completion bash > /etc/bash_completion.d/obelisk

# Zsh
obelisk completion zsh > "${fpath[1]}/_obelisk"

# Fish
obelisk completion fish > ~/.config/fish/completions/obelisk.fish

# PowerShell
obelisk completion powershell > obelisk.ps1
```

---

### 10. **help** - Command Help

Get help for any command.

```bash
# General help
obelisk help

# Command-specific help
obelisk help check
obelisk check --help
```

---

## Global Flags

These flags work with all commands:

- `--api-key string` - Gemini API key (overrides GOOGLE_API_KEY env var)
- `--model string` - AI model to use (default: "gemini-2.5-flash")
- `--no-color` - Disable colored output
- `-v, --verbose` - Enable verbose output
- `-h, --help` - Show help
- `--version` - Show version

---

## Environment Variables

- `GOOGLE_API_KEY` - Gemini API key (can be overridden by --api-key flag)
- `NO_COLOR` - Disable colored output (set to any value)
- `OBELISK_CONFIG` - Custom config file path

---

## Examples

### Basic Usage

```bash
# Quick check
obelisk check --path .

# With AI analysis
obelisk check --path . --api-key YOUR_KEY

# Skip AI (faster)
obelisk check --path . --skip-ai
```

### CI/CD Integration

```bash
# GitHub Actions
obelisk scan --path . --fail-on-critical --output json

# GitLab CI
obelisk scan --path $CI_PROJECT_DIR --fail-on-critical

# Jenkins
obelisk scan --path ${WORKSPACE} --fail-on-critical
```

### Report Generation

```bash
# JSON report
obelisk report --path . --output report.json --format json

# Markdown report
obelisk report --path . --output REPORT.md --format markdown

# HTML report
obelisk report --path . --output report.html --format html
```

### Configuration

```bash
# One-time setup
obelisk config set api-key YOUR_API_KEY
obelisk config set model gemini-2.5-flash

# Then use without flags
obelisk check --path .
```

---

## Tips

1. **Use Interactive Mode** for development - provides best UX
2. **Use `scan` for CI/CD** - optimized for automation
3. **Configure API key once** - saves typing
4. **Install Git hook** - catch issues before pushing
5. **Generate reports** - track progress over time
6. **Use `--skip-ai`** - when you need speed over depth

---

## Getting Help

- **Documentation**: See README.md
- **Issues**: https://github.com/Swif7ify/Obelisk-CLI/issues
- **Command Help**: `obelisk [command] --help`
- **Version Info**: `obelisk --version`
