# Publishing Obelisk CLI Extension

## Overview

This guide covers how to publish the Obelisk CLI extension to:

1. **VS Code Marketplace** (for VS Code users)
2. **Open VSX Registry** (for IBM Bob IDE and other VS Code alternatives)
3. **Manual Distribution** (VSIX file)

---

## Prerequisites

### 1. Install Required Tools

```bash
# Install vsce (VS Code Extension Manager)
npm install -g @vscode/vsce

# Install ovsx (Open VSX CLI)
npm install -g ovsx
```

### 2. Create Publisher Accounts

#### VS Code Marketplace

1. Go to https://marketplace.visualstudio.com/manage
2. Sign in with Microsoft account
3. Create a publisher (e.g., `onedev-ph`)
4. Generate a Personal Access Token (PAT):
    - Go to https://dev.azure.com/
    - User Settings → Personal Access Tokens
    - Create new token with **Marketplace (Manage)** scope
    - Save the token securely

#### Open VSX Registry (for IBM Bob IDE)

1. Go to https://open-vsx.org/
2. Sign in with GitHub
3. Generate an access token:
    - Settings → Access Tokens
    - Create new token
    - Save the token securely

---

## Building the Extension

### 1. Install Dependencies

```bash
cd extension
npm install
```

### 2. Compile TypeScript

```bash
npm run compile
```

### 3. Package the Extension

```bash
# Create .vsix file
npm run package

# Or manually:
vsce package
```

This creates `obelisk-cli-0.1.0.vsix` in the extension directory.

---

## Publishing to VS Code Marketplace

### 1. Login to VS Code Marketplace

```bash
vsce login onedev-ph
# Enter your Personal Access Token when prompted
```

### 2. Publish the Extension

```bash
# Publish current version
vsce publish

# Or publish with version bump
vsce publish patch  # 0.1.0 → 0.1.1
vsce publish minor  # 0.1.0 → 0.2.0
vsce publish major  # 0.1.0 → 1.0.0
```

### 3. Verify Publication

- Visit: https://marketplace.visualstudio.com/items?itemName=onedev-ph.obelisk-cli
- Users can now install via: `ext install onedev-ph.obelisk-cli`

---

## Publishing to Open VSX (IBM Bob IDE)

### 1. Login to Open VSX

```bash
ovsx login
# Enter your Open VSX access token when prompted
```

### 2. Publish the Extension

```bash
# Publish the .vsix file
ovsx publish obelisk-cli-0.1.0.vsix -p YOUR_ACCESS_TOKEN
```

### 3. Verify Publication

- Visit: https://open-vsx.org/extension/onedev-ph/obelisk-cli
- IBM Bob IDE users can now install it

---

## Manual Distribution (VSIX File)

### 1. Build VSIX Package

```bash
cd extension
vsce package
```

### 2. Distribute the VSIX

Upload `obelisk-cli-0.1.0.vsix` to:

- GitHub Releases
- Your website
- Direct download link

### 3. Users Install Manually

**VS Code:**

```bash
code --install-extension obelisk-cli-0.1.0.vsix
```

**Or via UI:**

1. Open VS Code
2. Extensions panel (Ctrl+Shift+X)
3. Click "..." menu → "Install from VSIX..."
4. Select the .vsix file

**IBM Bob IDE:**
Same process as VS Code (Bob is based on VS Code/Eclipse Theia)

---

## IBM Bob IDE Specific Instructions

### What is IBM Bob IDE?

IBM Bob IDE is based on **Eclipse Theia**, which is compatible with VS Code extensions. Your extension will work in Bob IDE if:

1. ✅ It's published to **Open VSX Registry** (Bob's default marketplace)
2. ✅ It doesn't use VS Code-specific APIs (yours doesn't)
3. ✅ It's packaged as a standard .vsix file

### Publishing to Bob IDE

**Option 1: Open VSX Registry (Recommended)**

```bash
# Publish to Open VSX
ovsx publish obelisk-cli-0.1.0.vsix -p YOUR_TOKEN
```

Bob IDE users can then install via:

- Extensions panel → Search "Obelisk CLI"
- Or: Settings → Extensions → Install from Open VSX

**Option 2: Manual VSIX Installation**

1. Build the .vsix file
2. Share it with Bob IDE users
3. They install via: Extensions → Install from VSIX

**Option 3: Private Registry**
If IBM provides a private extension registry for the hackathon:

1. Contact IBM for registry URL and credentials
2. Configure your extension to publish there
3. Follow their specific publishing process

---

## Pre-Publication Checklist

Before publishing, ensure:

### Required Files

- [x] `package.json` - Extension manifest
- [x] `README.md` - User documentation
- [x] `LICENSE` - MIT license
- [x] `CHANGELOG.md` - Version history
- [x] `resources/obelisk-icon.png` - Extension icon (128x128)

### Package.json Validation

- [x] `name` - Lowercase, no spaces
- [x] `displayName` - User-friendly name
- [x] `description` - Clear, concise description
- [x] `version` - Semantic versioning (0.1.0)
- [x] `publisher` - Your publisher ID
- [x] `repository` - GitHub URL
- [x] `icon` - Path to icon file
- [x] `categories` - Appropriate categories
- [x] `keywords` - Searchable keywords

### Testing

- [ ] Test in VS Code
- [ ] Test in IBM Bob IDE (if available)
- [ ] Test all commands work
- [ ] Test with and without Obelisk CLI installed
- [ ] Test error handling

### Documentation

- [ ] Update README with installation instructions
- [ ] Add screenshots/GIFs
- [ ] Document all settings
- [ ] Document all commands
- [ ] Add troubleshooting section

---

## Updating the Extension

### 1. Update Version

Edit `package.json`:

```json
{
	"version": "0.2.0"
}
```

### 2. Update CHANGELOG

Add new version entry to `CHANGELOG.md`:

```markdown
## [0.2.0] - 2026-05-16

### Added

- New feature X

### Fixed

- Bug Y
```

### 3. Rebuild and Republish

```bash
# Compile
npm run compile

# Publish to VS Code Marketplace
vsce publish

# Publish to Open VSX
vsce package
ovsx publish obelisk-cli-0.2.0.vsix -p YOUR_TOKEN
```

---

## Quick Start Commands

```bash
# One-time setup
cd extension
npm install
npm run compile

# Build VSIX
npm run package

# Publish to VS Code Marketplace
vsce login onedev-ph
vsce publish

# Publish to Open VSX (IBM Bob IDE)
ovsx publish obelisk-cli-0.1.0.vsix -p YOUR_TOKEN

# Test locally
code --install-extension obelisk-cli-0.1.0.vsix
```

---

## Troubleshooting

### "Publisher not found"

- Create publisher at https://marketplace.visualstudio.com/manage
- Update `publisher` field in package.json

### "Extension validation failed"

- Run: `vsce package` to see validation errors
- Fix issues in package.json
- Ensure all required files exist

### "Icon not found"

- Verify `resources/obelisk-icon.png` exists
- Icon must be 128x128 pixels
- Use PNG format

### "Extension doesn't work in Bob IDE"

- Ensure published to Open VSX Registry
- Check Bob IDE's extension marketplace settings
- Try manual VSIX installation

---

## Resources

- **VS Code Extension API**: https://code.visualstudio.com/api
- **Publishing Extensions**: https://code.visualstudio.com/api/working-with-extensions/publishing-extension
- **Open VSX Registry**: https://open-vsx.org/
- **Eclipse Theia**: https://theia-ide.org/
- **IBM Bob IDE**: Contact IBM for specific documentation

---

## Support

For issues or questions:

- GitHub Issues: https://github.com/Swif7ify/Obelisk-CLI/issues
- Email: support@onedev.ph
- IBM Bob Hackathon: Contact hackathon organizers

---

## License

MIT License - See LICENSE file for details
