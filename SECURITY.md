# Security Policy

> 🏆 **Obelisk CLI** — Created during the IBM Bob Hackathon (May 15–17, 2026) by OneDev PH.

Obelisk CLI is a security-first tool. We take the security of both our users' codebases and our own application seriously.

---

## 🔐 Credential Security Architecture

### OS-Level Keyring Encryption

Obelisk **never** stores API keys in plaintext. Instead, it integrates directly with your Operating System's native Credential Manager:

| Platform    | Backend                                                |
| ----------- | ------------------------------------------------------ |
| **Windows** | Windows Credential Manager (`wincred`)                 |
| **macOS**   | macOS Keychain (`security`)                            |
| **Linux**   | Secret Service API (GNOME Keyring / KWallet via D-Bus) |

When you run `obelisk config set api-key <YOUR_KEY>` or set it via the interactive TUI, the following happens:

1. The key is passed to the `go-keyring` library.
2. The OS encrypts the key using its native cryptographic APIs.
3. The decryption rights are bound to your authenticated user session.
4. The key is **never** written to `~/.obelisk/config.json` or any other file on disk.

### Automatic Migration

If you are upgrading from an older version that stored the key in plaintext inside `config.json`, Obelisk will:

1. Detect the legacy plaintext `api_key` field on the next run.
2. Securely vault it into the OS keyring.
3. Permanently remove the `api_key` field from `config.json` and overwrite the file.

This migration is silent and automatic — no user action required.

### Environment Variable Support

For CI/CD pipelines where no interactive keyring is available, pass the API key securely at runtime:

```bash
# Using environment variables (recommended for CI/CD)
export GEMINI_API_KEY="your-key-here"
obelisk scan --strict

# Or
export GOOGLE_API_KEY="your-key-here"
obelisk scan --format json
```

### Config File Security

The configuration file (`~/.obelisk/config.json`) contains **only** non-sensitive settings:

```json
{
	"model": "gemini-2.5-flash",
	"default_path": "",
	"report_format": "md",
	"no_color": false
}
```

- The config directory is created with `0700` permissions (owner-only access).
- The config file is written with `0600` permissions (owner-only read/write).
- No credentials, tokens, or secrets are ever written to this file.

---

## 🛡️ Data Privacy & API Communication

1. **HTTPS Only** — Obelisk communicates with the Google Gemini API exclusively over HTTPS (TLS 1.2+).
2. **No Training** — Your source code is **never** used to train Google's models (per standard Gemini API terms).
3. **Minimal Data** — Obelisk sends only a structured summary of scan findings and directory structure to the AI. Full source code is never uploaded.
4. **Local-First** — All scanners (secret detection, complexity analysis, tech debt tracking, native syntax checking) run entirely locally. The AI is optional and used only for the final grading synthesis.

---

## 🚨 Reporting a Vulnerability

If you discover a security vulnerability in Obelisk CLI, please report it responsibly.

**Do NOT open a public GitHub issue for security vulnerabilities.**

Instead, please email: **weareonedev@gmail.com** (or open a private security advisory on GitHub).

### What to Include

- Description of the vulnerability
- Steps to reproduce
- Potential impact
- Suggested fix (if any)

### Response Timeline

- **Acknowledgment:** Within 48 hours
- **Assessment:** Within 1 week
- **Fix & Disclosure:** Coordinated with reporter

---

## Supported Versions

| Version | Supported |
| ------- | --------- |
| 0.1.x   | ✅        |

---

## ✅ Security Best Practices We Follow

Obelisk CLI is a security tool. We practice what we preach:

- **No hardcoded secrets** in the codebase
- **Dependencies are regularly audited** using our own dependency scanner
- **API keys are never logged, printed, or stored** in plaintext
- **OS-native credential encryption** via `go-keyring`
- **File permissions are strict** (`0700` dirs, `0600` files)
- **All sensitive config fields** use `json:"-"` struct tags to prevent accidental serialization
