# Obelisk VS Code Extension - Usage Guide

## 🚀 How to Use Obelisk

The Obelisk VS Code extension uses the Obelisk CLI installed on your machine.

**Requirements:**

- Obelisk CLI must be installed
- Available in your PATH or configured in settings

**Pros:**

- Works offline
- Full control over execution
- No network latency

---

## ⚙️ Configuration

### All Settings

```json
{
	// Path to local CLI
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
- **Obelisk: Stop Scan** - Stop a running scan

### Status Bar

Shows current scan status:

- 🏛️ Ready
- ⏳ Scanning...
- ✅ Grade: A (or B, C, D, F)
- ❌ Error

---

## 🎯 Usage Examples

### Example 1: Deep Analysis

1. Ensure Obelisk CLI is installed
2. Run "Obelisk: Scan Project" command
3. Get AI-powered health report
4. Review recommendations in Health Summary

### Example 2: CI/CD Integration

Use auto-scan on save:

```json
{
	"obelisk.scanOnSave": true,
	"obelisk.scanOnSaveDelay": 5000
}
```

---

## 🔧 Troubleshooting

**"Could not find 'obelisk'"**

- Install Obelisk CLI: See [main README](../README.md)
- Or set `obelisk.executablePath` to the full path

**"Scan failed"**

- Check if Obelisk CLI is working: Run `obelisk --version` in terminal
- Verify you have a valid project open
- Check the Output panel for detailed errors

---

## 🔒 Privacy & Security

- All analysis happens on your machine
- No data leaves your computer
- Full privacy

---

## 📚 Learn More

- [Main Documentation](../README.md)
- [Deployment Guide](../DEPLOYMENT.md)

---

**Made with ❤️ by OneDev PH**
