# 🤖 Obelisk MCP Server Guide

> Model Context Protocol (MCP) integration for AI assistants and IDEs

## 📋 Table of Contents

- [What is MCP?](#what-is-mcp)
- [Quick Start](#quick-start)
- [Local Setup](#local-setup)
- [Cloud Deployment](#cloud-deployment)
- [Available Tools](#available-tools)
- [Available Resources](#available-resources)
- [Client Configuration](#client-configuration)
- [Examples](#examples)
- [Troubleshooting](#troubleshooting)

---

## What is MCP?

**Model Context Protocol (MCP)** is an open standard that enables AI assistants to securely connect to external tools and data sources. Obelisk implements MCP to expose its code analysis capabilities to AI assistants like Claude, Bob IDE, and other MCP-compatible clients.

### Benefits

- 🔌 **Seamless Integration** - AI assistants can analyze your code directly
- 🔒 **Secure** - All analysis runs locally or on your infrastructure
- 🚀 **Real-time** - Get instant code health insights during conversations
- 🎯 **Context-Aware** - AI understands your project structure and issues

---

## Quick Start

### 1. Start MCP Server

```bash
# Local mode (stdio transport)
obelisk mcp

# HTTP mode (for remote access)
obelisk mcp --http --port 8080
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
# Default stdio mode (for local IDE integration)
obelisk mcp

# The server will output:
# MCP Server started (stdio transport)
# Available tools: scan_project, check_security, analyze_complexity, ...
```

---

## Cloud Deployment

For remote access and team collaboration, deploy Obelisk MCP Server to the cloud.

### Deploy to Render (Recommended)

See [DEPLOYMENT.md](DEPLOYMENT.md) for complete cloud deployment guide.

**Quick Deploy:**

1. Push code to GitHub
2. Connect to Render
3. Set `GEMINI_API_KEY` environment variable
4. Deploy with `obelisk mcp --http`

Your server will be available at: `https://your-app.onrender.com/sse`

---

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

**Description:** Analyze code complexity and identify spaghetti code

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

For cloud deployment:

```json
{
	"mcpServers": {
		"obelisk-cloud": {
			"url": "https://your-app.onrender.com/sse"
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

**HTTP/SSE Transport:**

```javascript
const eventSource = new EventSource("https://your-server.com/sse");

eventSource.onmessage = (event) => {
	const message = JSON.parse(event.data);
	// Handle MCP messages
};
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

# Check for port conflicts (HTTP mode)
netstat -ano | findstr :8080
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

### HTTP Mode Issues

**Issue:** Can't access HTTP endpoint

**Solutions:**

1. **Check firewall:**

    ```bash
    # Windows
    netsh advfirewall firewall add rule name="Obelisk MCP" dir=in action=allow protocol=TCP localport=8080
    ```

2. **Verify port:**

    ```bash
    # Test endpoint
    curl http://localhost:8080/health
    ```

3. **Check CORS:**
    - Server includes CORS headers by default
    - Verify client origin is allowed

---

## Advanced Usage

### Custom Port

```bash
# Use custom port
obelisk mcp --http --port 3000
```

### Environment Variables

```bash
# Set via environment
export GEMINI_API_KEY=your-key
export OBELISK_MODEL=gemini-2.0-flash-exp
export PORT=8080

obelisk mcp --http
```

### Docker Deployment

```bash
# Build image
docker build -t obelisk-mcp .

# Run container
docker run -p 8080:8080 \
  -e GEMINI_API_KEY=your-key \
  obelisk-mcp
```

---

## Security Considerations

### API Key Protection

- ✅ Never commit API keys to Git
- ✅ Use environment variables in production
- ✅ Rotate keys regularly
- ✅ Use OS keyring for local storage

### Network Security

- ✅ Use HTTPS in production (Render provides this)
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
- **Cloud Deployment:** [DEPLOYMENT.md](DEPLOYMENT.md)
- **GitHub Issues:** https://github.com/Swif7ify/Obelisk-CLI/issues

---

## Support

For help with MCP integration:

- 📧 Email: weareonedev@gmail.com
- 🐛 Issues: https://github.com/Swif7ify/Obelisk-CLI/issues
- 📖 Docs: https://github.com/Swif7ify/Obelisk-CLI

---

**Built with ❤️ by OneDev PH for the IBM Bob Hackathon 2026**
