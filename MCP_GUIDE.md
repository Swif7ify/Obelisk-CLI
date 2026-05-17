# 🤖 Obelisk MCP Server Guide

> Model Context Protocol (MCP) integration for AI assistants and IDEs

## 📋 Table of Contents

- [What is MCP?](#what-is-mcp)
- [Quick Start](#quick-start)
- [Local Setup](#local-setup)
- [Available Tools](#available-tools)
- [Available Resources](#available-resources)
- [Client Configuration](#client-configuration)
- [Examples](#examples)
- [Troubleshooting](#troubleshooting)

---

## What is MCP?

**Model Context Protocol (MCP)** is an open standard that enables AI assistants to securely connect to external tools and data sources. Obelisk uses MCP to expose its code analysis capabilities to clients such as Claude, Bob IDE, and other compatible tools.

### Benefits

- 🔌 **Seamless Integration** - AI assistants can analyze your code directly
- 🔒 **Secure** - Analysis runs locally or on your infrastructure
- 🚀 **Real-time** - Get code health insights during conversations
- 🎯 **Context-Aware** - AI can use project structure and findings as context

---

## Quick Start

### 1. Start the MCP Server

```bash
# Local mode (stdio transport)
obelisk mcp
```

### 2. Configure Your AI Assistant

See [Client Configuration](#client-configuration) for specific setup instructions.

---

## Local Setup

### Prerequisites

- Obelisk CLI installed ([Installation Guide](README.md#installation))
- Google Gemini API key ([Get one here](https://aistudio.google.com/app/apikey))

### Configuration

```bash
# Set your API key
obelisk config set api-key YOUR_GEMINI_API_KEY

# Or use environment variable
export GEMINI_API_KEY=your-api-key-here
```

### Start Server

```bash
# Default stdio mode for local IDE integration
obelisk mcp

# The server will output:
# MCP Server started (stdio transport)
# Available tools: scan_project, check_security, analyze_complexity, ...
```

## Available Tools

MCP clients can invoke these tools to analyze your code:

### `scan_project`

**Description:** Run a comprehensive project health scan

**Parameters:**

- `path` (string, required) - Project directory path
- `skip_ai` (boolean, optional) - Skip AI analysis for faster results

**Returns:** Complete scan results with security, quality, and architecture findings

**Example:**

```json
{
	"name": "scan_project",
	"arguments": {
		"path": "/path/to/project",
		"skip_ai": false
	}
}
```

---

### `check_security`

**Description:** Security-focused scan for vulnerabilities

**Parameters:**

- `path` (string, required) - Project directory path

**Returns:** Security findings including secrets, exposed files, and vulnerabilities

**Example:**

```json
{
	"name": "check_security",
	"arguments": {
		"path": "/path/to/project"
	}
}
```

---

### `analyze_complexity`

**Description:** Analyze code complexity and identify high-complexity code

**Parameters:**

- `path` (string, required) - Project directory path

**Returns:** Cyclomatic complexity metrics and problematic files

**Example:**

```json
{
	"name": "analyze_complexity",
	"arguments": {
		"path": "/path/to/project"
	}
}
```

---

### `track_tech_debt`

**Description:** Track technical debt markers (TODO, FIXME, HACK, XXX)

**Parameters:**

- `path` (string, required) - Project directory path

**Returns:** List of technical debt items with locations

**Example:**

```json
{
	"name": "track_tech_debt",
	"arguments": {
		"path": "/path/to/project"
	}
}
```

---

### `audit_dependencies`

**Description:** Audit project dependencies for issues

**Parameters:**

- `path` (string, required) - Project directory path

**Returns:** Dependency analysis including unused, deprecated, and vulnerable packages

**Example:**

```json
{
	"name": "audit_dependencies",
	"arguments": {
		"path": "/path/to/project"
	}
}
```

---

### `get_health_report`

**Description:** Get AI-powered health assessment with grade

**Parameters:**

- `path` (string, required) - Project directory path

**Returns:** Health score (A-F), AI insights, and recommendations

**Example:**

```json
{
	"name": "get_health_report",
	"arguments": {
		"path": "/path/to/project"
	}
}
```

---

## Available Resources

MCP clients can read these resources to get cached analysis results:

### `obelisk://scan/latest`

**Description:** Latest scan results

**MIME Type:** `application/json`

**Content:** Complete scan results from the most recent analysis

---

### `obelisk://health/score`

**Description:** Current project health score

**MIME Type:** `text/plain`

**Content:** Health grade (A-F) with brief summary

---

### `obelisk://findings/security`

**Description:** Security findings only

**MIME Type:** `application/json`

**Content:** Filtered list of security-related findings

---

### `obelisk://findings/quality`

**Description:** Code quality findings only

**MIME Type:** `application/json`

**Content:** Filtered list of quality-related findings

---

### `obelisk://findings/architecture`

**Description:** Architecture findings only

**MIME Type:** `application/json`

**Content:** Filtered list of architecture-related findings

---

## Client Configuration

### Bob IDE

Add to your Bob configuration:

```json
{
	"mcpServers": {
		"obelisk": {
			"command": "obelisk",
			"args": ["mcp"]
		}
	}
}
```

---

### Claude Desktop

Add to `claude_desktop_config.json`:

**Windows:** `%APPDATA%\Claude\claude_desktop_config.json`

**macOS:** `~/Library/Application Support/Claude/claude_desktop_config.json`

```json
{
	"mcpServers": {
		"obelisk": {
			"command": "obelisk",
			"args": ["mcp"]
		}
	}
}
```

---

### Cline (VS Code Extension)

Add to Cline settings:

```json
{
	"mcpServers": {
		"obelisk": {
			"command": "obelisk",
			"args": ["mcp"]
		}
	}
}
```

---

### Custom MCP Client

For custom implementations:

**Stdio Transport:**

```javascript
const { spawn } = require("child_process");

const server = spawn("obelisk", ["mcp"]);
// Communicate via stdin/stdout using JSON-RPC 2.0
```

---

## Examples

### Example 1: Security Scan via AI Assistant

**User:** "Check my project for security issues"

**AI Assistant:** _Invokes `check_security` tool_

**Response:**

```
Found 3 security issues:
1. Hardcoded API key in config.js (line 15)
2. .env file tracked in Git
3. Weak password validation in auth.js
```

---

### Example 2: Health Report

**User:** "What's the overall health of my codebase?"

**AI Assistant:** _Invokes `get_health_report` tool_

**Response:**

```
Health Score: B

Strengths:
- Good test coverage
- Clean architecture
- No critical security issues

Areas for Improvement:
- 5 files with high complexity
- 12 TODO items need attention
- Consider updating deprecated dependencies
```

---

### Example 3: Complexity Analysis

**User:** "Which files are too complex?"

**AI Assistant:** _Invokes `analyze_complexity` tool_

**Response:**

```
High Complexity Files:
1. src/utils/parser.js - Complexity: 45 (Refactor recommended)
2. src/api/handler.js - Complexity: 38 (Needs simplification)
3. src/core/engine.js - Complexity: 32 (Consider splitting)
```

---

## Troubleshooting

### Server Won't Start

**Issue:** `obelisk mcp` fails to start

**Solutions:**

```bash
# Check if Obelisk is installed
obelisk version

# Verify API key is set
obelisk config get api-key
```

---

### Client Can't Connect

**Issue:** AI assistant can't connect to MCP server

**Solutions:**

1. **Verify server is running:**

    ```bash
    # Check process
    ps aux | grep obelisk
    ```

2. **Check configuration:**
    - Ensure correct command path in client config
    - Verify args are correct: `["mcp"]`

3. **Test manually:**
    ```bash
    # Start server with verbose output
    obelisk mcp --verbose
    ```

---

### Tools Not Working

**Issue:** MCP tools return errors

**Solutions:**

1. **Check API key:**

    ```bash
    # Verify key is set
    echo $GEMINI_API_KEY
    ```

2. **Verify project path:**
    - Ensure path exists and is accessible
    - Use absolute paths for reliability

3. **Check logs:**
    ```bash
    # Enable debug logging
    obelisk mcp --verbose
    ```

---

## Security Considerations

### API Key Protection

- ✅ Never commit API keys to Git
- ✅ Use environment variables in production
- ✅ Rotate keys regularly
- ✅ Use OS keyring for local storage

### Network Security

- ✅ Implement rate limiting for public endpoints
- ✅ Consider authentication for team deployments
- ✅ Monitor access logs

### Data Privacy

- ✅ Code analysis runs on your infrastructure
- ✅ Only summaries sent to Gemini API (not full code)
- ✅ No data stored by Obelisk servers
- ✅ Comply with your organization's data policies

---

## Resources

- **MCP Specification:** https://modelcontextprotocol.io/
- **Obelisk Documentation:** [README.md](README.md)
- **GitHub Issues:** https://github.com/Swif7ify/Obelisk-CLI/issues

---

## Support

For help with MCP integration:

- 📧 Email: weareonedev@gmail.com
- 🐛 Issues: https://github.com/Swif7ify/Obelisk-CLI/issues
- 📖 Docs: https://github.com/Swif7ify/Obelisk-CLI

---

**Built with ❤️ by OneDev PH for the IBM Bob Hackathon 2026**
