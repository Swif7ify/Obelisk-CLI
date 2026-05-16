# Obelisk VS Code Extension - Usage Guide

## 🚀 Two Ways to Use Obelisk

The Obelisk VS Code extension now supports **two modes**:

### 1. **Local Mode** (Default)

Uses the Obelisk CLI installed on your machine.

**Requirements:**

- Obelisk CLI must be installed
- Available in your PATH or configured in settings

**Pros:**

- Works offline
- Full control over execution
- No network latency

### 2. **Cloud Mode** (New!)

Connects to the official Obelisk Cloud MCP server.

**Requirements:**

- Internet connection
- No local installation needed!

**Pros:**

- No installation required
- Always up-to-date
- Works on any machine
- Instant setup

---

## ⚙️ Configuration

### Switch to Cloud Mode

1. Open VS Code Settings (`Ctrl+,` or `Cmd+,`)
2. Search for "Obelisk"
3. Change **Obelisk: Mode** from `local` to `cloud`
4. Reload VS Code

Or edit `settings.json` directly:

```json
{
	"obelisk.mode": "cloud"
}
```

### Cloud Server URL

The default cloud server is:

```
https://mcp-obelisk.onedevph.online/sse
```

To use a custom server:

```json
{
	"obelisk.mode": "cloud",
	"obelisk.cloudServerUrl": "https://your-custom-server.com/sse"
}
```

### All Settings

```json
{
	// Mode: "local" or "cloud"
	"obelisk.mode": "local",

	// Cloud server URL (only used in cloud mode)
	"obelisk.cloudServerUrl": "https://mcp-obelisk.onedevph.online/sse",

	// Path to local CLI (only used in local mode)
	"obelisk.executablePath": "obelisk",

	// Auto-scan when workspace opens
	"obelisk.scanOnOpen": true,

	// Auto-scan when files are saved
	"obelisk.scanOnSave": false,

	// Delay before auto-scan after save (ms)
	"obelisk.scanOnSaveDelay": 3000,

	// Skip AI analysis (faster scans)
	"obelisk.skipAI": false
}
```

---

## 📊 Features

### Sidebar Panel

The Obelisk panel in the Activity Bar shows:

1. **Findings Tree** - All issues found in your project
    - Click any finding to jump to the file and line
    - Organized by severity (Critical, Error, Warning, Info)

2. **Health Summary** - Visual health report
    - Overall grade (A-F)
    - Security, Architecture, and Quality scores
    - Top issues and recommendations

### Commands

Access via Command Palette (`Ctrl+Shift+P` or `Cmd+Shift+P`):

- **Obelisk: Scan Project** - Run a full scan
- **Obelisk: Refresh Scan** - Re-run the scan
- **Obelisk: Clear Findings** - Clear all results
- **Obelisk: Stop Scan** - Stop a running scan (local mode only)

### Status Bar

Shows current scan status:

- 🏛️ Ready
- ⏳ Scanning...
- ✅ Grade: A (or B, C, D, F)
- ❌ Error

---

## 🎯 Usage Examples

### Example 1: Quick Security Check (Cloud Mode)

1. Set mode to `cloud` in settings
2. Open your project in VS Code
3. Extension auto-scans on open
4. View security findings in the sidebar
5. Click any finding to see the code

### Example 2: Deep Analysis (Local Mode)

1. Ensure Obelisk CLI is installed
2. Set mode to `local` in settings
3. Run "Obelisk: Scan Project" command
4. Get AI-powered health report
5. Review recommendations in Health Summary

### Example 3: CI/CD Integration

Use local mode with auto-scan on save:

```json
{
	"obelisk.mode": "local",
	"obelisk.scanOnSave": true,
	"obelisk.scanOnSaveDelay": 5000
}
```

---

## 🔧 Troubleshooting

### Cloud Mode Issues

**"Connection timeout"**

- Check your internet connection
- Verify the cloud server URL is correct
- The free tier may sleep after inactivity (first request takes ~30s)

**"Invalid response from MCP server"**

- Server might be updating
- Try again in a few moments
- Check server status at the URL

### Local Mode Issues

**"Could not find 'obelisk'"**

- Install Obelisk CLI: See [main README](../README.md)
- Or set `obelisk.executablePath` to the full path

**"Scan failed"**

- Check if Obelisk CLI is working: Run `obelisk --version` in terminal
- Verify you have a valid project open
- Check the Output panel for detailed errors

---

## 🆚 Mode Comparison

| Feature        | Local Mode      | Cloud Mode            |
| -------------- | --------------- | --------------------- |
| Installation   | Required        | Not required          |
| Internet       | Optional        | Required              |
| Speed          | Fast            | Network dependent     |
| Offline        | ✅ Yes          | ❌ No                 |
| Always updated | Manual          | ✅ Automatic          |
| Custom config  | ✅ Full control | Limited               |
| Privacy        | ✅ All local    | Results sent to cloud |

---

## 🔒 Privacy & Security

### Local Mode

- All analysis happens on your machine
- No data leaves your computer
- Full privacy

### Cloud Mode

- Only analysis results are sent to the cloud
- Your source code is NOT uploaded
- Results are not stored permanently
- Communication over HTTPS

---

## 📚 Learn More

- [Main Documentation](../README.md)
- [Deployment Guide](../DEPLOYMENT.md)
- [MCP Guide](../MCP_GUIDE.md)
- [Quick Start](../QUICK_START_MCP.md)

---

**Made with ❤️ by OneDev PH**
