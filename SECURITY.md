# Security Policy

## Data Privacy & Security

1. Obelisk communicates directly with the official Google Gemini API over HTTPS (TLS 1.2+).
2. Code is **never** used to train Google's models (as per standard Gemini API enterprise terms).
3. The API key is required but is **never stored in plaintext**. Obelisk integrates with your Operating System's native Credential Manager (Windows Credential Manager, macOS Keychain, Linux Secret Service). This means your key is encrypted locally by the OS and strictly tied to your authenticated login session. Even if `config.json` is stolen, the key remains safe.
4. If you prefer, you can completely avoid storing the key locally by passing it at runtime via the `GEMINI_API_KEY` or `GOOGLE_API_KEY` environment variables.

## Reporting a Vulnerability

If you discover a security vulnerability in Obelisk CLI, please report it responsibly.

**Do NOT open a public GitHub issue for security vulnerabilities.**

Instead, please email: **security@obelisk-cli.dev** (or open a private security advisory on GitHub).

### What to Include

- Description of the vulnerability
- Steps to reproduce
- Potential impact
- Suggested fix (if any)

### Response Timeline

- **Acknowledgment:** Within 48 hours
- **Assessment:** Within 1 week
- **Fix & Disclosure:** Coordinated with reporter

## Supported Versions

| Version | Supported |
|---------|-----------|
| 0.1.x   | ✅        |

## Security Best Practices

Obelisk CLI itself is a security tool. We practice what we preach:

- No hardcoded secrets in the codebase
- Dependencies are regularly audited
- API keys are never logged or stored
