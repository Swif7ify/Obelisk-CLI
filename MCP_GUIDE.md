# 🏛️ Obelisk MCP Server Guide

## Overview

Obelisk can now run as a **Model Context Protocol (MCP) server**, exposing its powerful code analysis capabilities to AI assistants and IDEs that support MCP. This allows AI assistants to scan your projects, check security, analyze code quality, and provide architectural insights programmatically.

## What is MCP?

The Model Context Protocol (MCP) is an open protocol that enables AI assistants to securely access tools and data sources. By running Obelisk as an MCP server, you give AI assistants the ability to:

- Scan projects for security vulnerabilities
- Analyze code complexity and quality
- Track technical debt
- Audit dependencies
- Generate AI-powered health reports

## Quick Start

### 1. Build Obelisk with MCP Support

```bash
go build -o obelisk.exe
```

### 2. Configure MCP Server

The MCP server has been automatically configured in your Bob IDE settings at:
`C:\Users\nikon\.bob\settings\mcp_settings.json`

Current configuration:
```json
{
  "mcpServers": {
    "obelisk": {
      "command": "C:\\Users\\nikon\\Desktop\\Obelisk-CLI\\obelisk.exe",
      "args": ["mcp"],
      "env": {
        "GEMINI_API_KEY": ""
      },
      "disabled": false,
      "alwaysAllow": [],
      "disabledTools": []
    }
  }
}
```

### 3. Add Your Gemini API Key (Optional)

To enable AI-powered health reports, add your Gemini API key:

1. Get a free API key from: https://aistudio.google.com/app/apikey
2. Edit `mcp_settings.json` and replace the empty `GEMINI_API_KEY` value
3. Restart Bob IDE to apply changes

**Note:** The MCP server works without an API key, but AI-powered features like `get_health_report` will be unavailable.

### 4. Restart Bob IDE

After configuration, restart Bob IDE to load the Obelisk MCP server.

## Available Tools

### 1. `scan_project` - Comprehensive Project Analysis

Performs a full health scan including security, architecture, and code quality analysis.

**Parameters:**
- `path` (string, optional): Project directory path (defaults to current directory)
- `skip_ai` (boolean, optional): Skip AI-powered analysis (default: false)

**Example Usage:**
```
"Scan the current project for issues"
"Run a full scan on C:\path\to\project"
"Scan this project without AI analysis"
```

**Returns:**
- Project metadata (type, file count, directory count)
- All findings categorized by severity (critical, error, warning, info)
- Health report with grade (A-F) if AI is enabled

---

### 2. `check_security` - Security-Focused Scan

Scans specifically for security vulnerabilities including secrets, exposed credentials, and .gitignore issues.

**Parameters:**
- `path` (string, optional): Project directory path

**Example Usage:**
```
"Check security vulnerabilities in this project"
"Scan for secrets and exposed credentials"
```

**Returns:**
- Security findings only
- Count by severity level
- Specific file locations and line numbers

---

### 3. `analyze_complexity` - Code Complexity Analysis

Analyzes cyclomatic complexity to identify spaghetti code and maintainability issues.

**Parameters:**
- `path` (string, optional): Project directory path
- `threshold` (number, optional): Complexity threshold (default: 10)

**Example Usage:**
```
"Analyze code complexity in this project"
"Find functions with high cyclomatic complexity"
```

**Returns:**
- Complexity findings
- Functions exceeding threshold
- Maintainability warnings

---

### 4. `track_tech_debt` - Technical Debt Tracking

Tracks TODO, FIXME, HACK, and XXX comments across the codebase.

**Parameters:**
- `path` (string, optional): Project directory path

**Example Usage:**
```
"Track technical debt in this project"
"Find all TODO and FIXME comments"
```

**Returns:**
- All technical debt markers
- File locations and line numbers
- Total count

---

### 5. `audit_dependencies` - Dependency Audit

Audits package.json for deprecated, vulnerable, or unused dependencies.

**Parameters:**
- `path` (string, optional): Project directory path

**Example Usage:**
```
"Audit dependencies in this project"
"Check for vulnerable or deprecated packages"
```

**Returns:**
- Dependency findings
- Vulnerable packages
- Deprecated packages
- Unused dependencies

---

### 6. `get_health_report` - AI Health Assessment

Generates an AI-powered health report with grade (A-F) and recommendations.

**Parameters:**
- `path` (string, optional): Project directory path

**Requirements:**
- Requires `GEMINI_API_KEY` environment variable

**Example Usage:**
```
"Generate a health report for this project"
"What's the overall health score of this codebase?"
```

**Returns:**
- Health grade (A-F)
- Overall score and category scores (security, architecture, quality)
- AI-generated summary
- Praise for good practices
- Actionable recommendations
- Top priority issues

---

## Available Resources

Resources provide quick access to cached scan results without re-running scans.

### 1. `obelisk://scan/latest`
Access the most recent project scan results including all findings.

### 2. `obelisk://health/score`
Current project health grade (A-F) and AI-generated recommendations.

### 3. `obelisk://findings/security`
All security-related findings from the latest scan.

### 4. `obelisk://findings/quality`
All code quality findings from the latest scan.

### 5. `obelisk://findings/architecture`
All architecture-related findings from the latest scan.

**Example Usage:**
```
"Show me the latest scan results"
"What's the current health score?"
"List all security findings"
```

---

## Environment Variables

### `GEMINI_API_KEY` (or `GOOGLE_API_KEY`)
Google Gemini API key for AI-powered analysis. Required for `get_health_report` tool.

**How to set:**
```json
{
  "mcpServers": {
    "obelisk": {
      "env": {
        "GEMINI_API_KEY": "your-api-key-here"
      }
    }
  }
}
```

### `OBELISK_MODEL`
AI model to use (default: `gemini-2.0-flash-exp`)

**Supported models:**
- `gemini-2.0-flash-exp` (default, fastest)
- `gemini-1.5-pro`
- `gemini-1.5-flash`

---

## Supported Project Types

Currently, Obelisk MCP server supports:

- ✅ JavaScript / TypeScript
- ✅ React
- ✅ Next.js
- 🔜 Go (Golang) - Coming Soon
- 🔜 Laravel (PHP) - Coming Soon
- 🔜 Python (Django/Flask) - Planned

---

## Testing the MCP Server

### Manual Test (Command Line)

You can test the MCP server manually using stdio:

```bash
# Start the server
.\obelisk.exe mcp

# In another terminal, send a JSON-RPC request:
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}' | .\obelisk.exe mcp
```

### Test with Bob IDE

1. Restart Bob IDE after configuration
2. Ask the AI assistant: "What MCP tools are available?"
3. Try a scan: "Scan the current project for issues"
4. Check results: "Show me the latest scan results"

---

## Troubleshooting

### Server Not Starting

**Issue:** MCP server doesn't appear in Bob IDE

**Solutions:**
1. Verify the executable path in `mcp_settings.json` is correct
2. Ensure `disabled` is set to `false`
3. Restart Bob IDE completely
4. Check Bob IDE logs for error messages

### Tools Not Working

**Issue:** Tools return errors or don't execute

**Solutions:**
1. Verify the project path exists
2. Check that the project is a supported type (JS/TS/React/Next.js)
3. For AI features, ensure `GEMINI_API_KEY` is set
4. Check stderr output for detailed error messages

### API Key Issues

**Issue:** `get_health_report` fails with API key error

**Solutions:**
1. Verify API key is correctly set in `mcp_settings.json`
2. Test API key at: https://aistudio.google.com/app/apikey
3. Ensure no extra spaces or quotes in the key value
4. Restart Bob IDE after adding the key

---

## Advanced Configuration

### Custom Installation Path

If you install Obelisk to a different location, update the `command` path:

```json
{
  "mcpServers": {
    "obelisk": {
      "command": "C:\\Program Files\\Obelisk\\obelisk.exe",
      "args": ["mcp"]
    }
  }
}
```

### Always Allow Specific Tools

To skip confirmation prompts for certain tools:

```json
{
  "mcpServers": {
    "obelisk": {
      "alwaysAllow": ["scan_project", "check_security"]
    }
  }
}
```

### Disable Specific Tools

To prevent certain tools from being used:

```json
{
  "mcpServers": {
    "obelisk": {
      "disabledTools": ["get_health_report"]
    }
  }
}
```

---

## Security Considerations

### API Key Storage

- API keys are stored in the MCP settings file
- The file is located in your user directory: `~/.bob/settings/mcp_settings.json`
- Ensure this file has appropriate permissions (readable only by you)
- Never commit this file to version control

### Project Access

- The MCP server can only access projects you explicitly specify
- It respects `.gitignore` and `.obelisk-ignore` files
- No data is sent to external services except the Gemini API (when enabled)

### Network Communication

- The MCP server communicates via stdio (no network ports)
- Only the Gemini API is contacted (when AI features are used)
- All communication uses HTTPS (TLS 1.2+)

---

## Example Workflows

### Daily Code Review

```
1. "Scan the current project"
2. "Show me all critical and error findings"
3. "What's the security status?"
4. "Generate a health report"
```

### Pre-Commit Check

```
1. "Check security vulnerabilities"
2. "Track technical debt"
3. "Analyze code complexity"
```

### Dependency Maintenance

```
1. "Audit dependencies"
2. "Show me vulnerable packages"
3. "List unused dependencies"
```

---

## Integration with CI/CD

While the MCP server is designed for interactive use with AI assistants, you can still use Obelisk in CI/CD pipelines with the standard CLI commands:

```bash
# Headless scan for CI/CD
obelisk scan --format json --strict

# This will exit with code 1 if critical issues are found
```

---

## Support and Feedback

- **GitHub Issues:** https://github.com/Swif7ify/Obelisk-CLI/issues
- **Documentation:** https://github.com/Swif7ify/Obelisk-CLI
- **Security Policy:** See SECURITY.md

---

## Version Information

- **Obelisk Version:** 0.1.0
- **MCP Protocol Version:** 2024-11-05
- **Supported Transports:** stdio

---

**Built during the IBM Bob Hackathon (May 15–17, 2026)**  
**By OneDev PH — with the help of IBM Bob IDE**