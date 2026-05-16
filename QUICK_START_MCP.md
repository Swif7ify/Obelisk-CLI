# 🚀 Quick Start - Use Obelisk MCP Server

Get started with Obelisk's AI-powered code analysis in under 2 minutes!

## ⚡ Instant Access - No Installation Required

Use our **official public MCP server** - no deployment, no setup, just connect and go!

```
https://mcp-obelisk.onedevph.online/sse
```

---

## 📱 Connect Your AI Client

### For Bob IDE

1. Open your Bob configuration file (`.bob/config.json` or `~/.bob/config.json`)
2. Add this configuration:

```json
{
	"mcpServers": {
		"obelisk": {
			"url": "https://mcp-obelisk.onedevph.online/sse"
		}
	}
}
```

3. Restart Bob IDE
4. Done! Ask Bob to analyze your code using Obelisk tools

### For Claude Desktop

1. Open Claude Desktop config:
    - **Mac**: `~/Library/Application Support/Claude/claude_desktop_config.json`
    - **Windows**: `%APPDATA%/Claude/claude_desktop_config.json`

2. Add this configuration:

```json
{
	"mcpServers": {
		"obelisk": {
			"url": "https://mcp-obelisk.onedevph.online/sse"
		}
	}
}
```

3. Restart Claude Desktop
4. Done! Claude can now use Obelisk to analyze your projects

### For Cline (VS Code Extension)

1. Open VS Code Settings
2. Search for "Cline MCP"
3. Add this configuration:

```json
{
	"mcpServers": {
		"obelisk": {
			"url": "https://mcp-obelisk.onedevph.online/sse"
		}
	}
}
```

4. Reload VS Code
5. Done! Cline can now use Obelisk tools

---

## 🎯 What Can You Do?

Once connected, ask your AI assistant to:

### Security Analysis

```
"Use Obelisk to check my project for security vulnerabilities"
"Scan for exposed API keys and secrets in my code"
```

### Code Quality

```
"Analyze the complexity of my codebase"
"Find all TODO and FIXME comments in my project"
```

### Full Health Check

```
"Run a complete Obelisk scan on my project"
"Generate a health report with recommendations"
```

### Dependency Audit

```
"Check my dependencies for vulnerabilities"
"Audit my package.json for issues"
```

---

## 🛠️ Available Tools

Your AI assistant can now use these Obelisk tools:

| Tool                 | Description                              |
| -------------------- | ---------------------------------------- |
| `scan_project`       | Full project health scan with AI grading |
| `check_security`     | Security vulnerability scan              |
| `analyze_complexity` | Code complexity analysis                 |
| `track_tech_debt`    | Find TODO/FIXME/HACK comments            |
| `audit_dependencies` | Dependency vulnerability check           |
| `get_health_report`  | AI-powered health assessment             |

---

## 📊 Available Resources

Access cached analysis results:

| Resource                          | Description           |
| --------------------------------- | --------------------- |
| `obelisk://scan/latest`           | Latest scan results   |
| `obelisk://health/score`          | Project health score  |
| `obelisk://findings/security`     | Security findings     |
| `obelisk://findings/quality`      | Code quality findings |
| `obelisk://findings/architecture` | Architecture findings |

---

## 💡 Example Conversations

### Example 1: Quick Security Check

**You:** "Use Obelisk to check my project for security issues"

**AI:** _Uses `check_security` tool and reports findings_

### Example 2: Full Analysis

**You:** "Run a complete Obelisk health check on my codebase"

**AI:** _Uses `scan_project` tool and provides detailed report with grade_

### Example 3: Tech Debt Tracking

**You:** "Find all the TODOs in my project using Obelisk"

**AI:** _Uses `track_tech_debt` tool and lists all technical debt markers_

---

## 🔒 Privacy & Security

- **Your code stays local** - Only analysis results are sent to the AI
- **No code storage** - We don't store your source code
- **Secure connection** - All communication over HTTPS
- **Open source** - Audit the code yourself on GitHub

---

## 🆘 Troubleshooting

### Connection Issues

- Verify the URL is correct: `https://mcp-obelisk.onedevph.online/sse`
- Check your internet connection
- Restart your AI client

### Tool Not Working

- Make sure you're asking the AI to "use Obelisk" or mention the tool name
- Try: "Use the Obelisk scan_project tool on my current directory"

### Need Help?

- Check [DEPLOYMENT.md](DEPLOYMENT.md) for detailed documentation
- Open an issue on [GitHub](https://github.com/Swif7ify/Obelisk-CLI/issues)

---

## 🚀 Want Your Own Instance?

See [DEPLOYMENT.md](DEPLOYMENT.md) to deploy your own Obelisk MCP server on:

- Render (free tier available)
- Fly.io
- Railway
- Heroku
- Any Docker-compatible platform

---

## 📚 Learn More

- [Full Documentation](README.md)
- [Deployment Guide](DEPLOYMENT.md)
- [MCP Guide](MCP_GUIDE.md)
- [API Reference](API.md)

---

**Made with ❤️ by OneDev PH**
