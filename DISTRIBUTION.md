# 📦 Obelisk CLI Distribution Guide

Complete guide for building, packaging, and distributing Obelisk CLI.

## 🎯 Distribution Methods

### 1. Windows MSI Installer (Recommended)

**Best for:** End users who want a professional installation experience

#### Build the MSI

```powershell
# Install WiX Toolset first
# Download from: https://wixtoolset.org/releases/

# Build the installer
.\installer\build-installer.ps1

# Output: installer/ObeliskCLI-1.1.0-x64.msi
```

#### Features

- ✅ Professional installation wizard
- ✅ License agreement (EULA)
- ✅ Custom installation directory
- ✅ Automatic PATH configuration
- ✅ Start menu shortcuts
- ✅ Clean uninstallation
- ✅ Upgrade support

#### Distribution

```powershell
# 1. Upload to GitHub Releases
# 2. Users download and double-click to install
# 3. Installer handles everything automatically
```

---

### 2. PowerShell Installer Script

**Best for:** Developers who prefer command-line installation

#### Usage

```powershell
# Clone and build
git clone https://github.com/Swif7ify/Obelisk-CLI.git
cd Obelisk-CLI
make build

# Run installer
powershell -ExecutionPolicy Bypass -File install.ps1

# Restart terminal
obelisk version
```

#### Features

- ✅ Quick installation
- ✅ Automatic PATH setup
- ✅ User or system-wide install
- ✅ No admin required (user mode)

#### Uninstall

```powershell
powershell -ExecutionPolicy Bypass -File uninstall.ps1
```

---

### 3. Portable Executable

**Best for:** Users who don't want to install anything

#### Usage

```powershell
# Download obelisk.exe from GitHub Releases
# Run directly without installation
.\obelisk.exe version
```

#### Features

- ✅ No installation required
- ✅ Run from any directory
- ✅ Perfect for USB drives
- ❌ Not in PATH (must use full path)

---

### 4. Go Install (For Developers)

**Best for:** Go developers who want the latest version

#### Usage

```bash
# Install directly from source
go install github.com/Swif7ify/Obelisk-CLI@latest

# Or install specific version
go install github.com/Swif7ify/Obelisk-CLI@v1.1.0
```

#### Features

- ✅ Always latest version
- ✅ Automatic updates with `go install`
- ✅ Integrates with Go toolchain
- ❌ Requires Go installed

---

## 🚀 Release Process

### 1. Prepare Release

```bash
# Update version in files
# - main.go (Version constant)
# - installer/obelisk.wxs (Version attribute)
# - README.md (version badge)

# Update CHANGELOG.md
# Document all changes since last release

# Commit changes
git add .
git commit -m "chore: bump version to 1.1.0"
git push
```

### 2. Build Artifacts

```powershell
# Build for Windows x64
$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -o bin/obelisk-windows-amd64.exe -ldflags "-s -w -X main.Version=1.1.0" .

# Build for Linux x64
$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -o bin/obelisk-linux-amd64 -ldflags "-s -w -X main.Version=1.1.0" .

# Build for macOS x64
$env:GOOS = "darwin"
$env:GOARCH = "amd64"
go build -o bin/obelisk-darwin-amd64 -ldflags "-s -w -X main.Version=1.1.0" .

# Build for macOS ARM64 (Apple Silicon)
$env:GOOS = "darwin"
$env:GOARCH = "arm64"
go build -o bin/obelisk-darwin-arm64 -ldflags "-s -w -X main.Version=1.1.0" .

# Build MSI installer (Windows only)
.\installer\build-installer.ps1 -Version "1.1.0"
```

### 3. Create GitHub Release

```bash
# Create and push tag
git tag -a v1.1.0 -m "Release version 1.1.0"
git push origin v1.1.0

# Go to GitHub → Releases → Create new release
# - Tag: v1.1.0
# - Title: Obelisk CLI v1.1.0
# - Description: Copy from CHANGELOG.md
```

### 4. Upload Artifacts

Upload these files to the GitHub Release:

```
✅ ObeliskCLI-1.1.0-x64.msi          (Windows Installer)
✅ obelisk-windows-amd64.exe         (Windows Portable)
✅ obelisk-linux-amd64               (Linux Binary)
✅ obelisk-darwin-amd64              (macOS Intel)
✅ obelisk-darwin-arm64              (macOS Apple Silicon)
✅ checksums.txt                     (SHA256 hashes)
```

### 5. Generate Checksums

```powershell
# Windows
Get-FileHash bin/*.exe, installer/*.msi -Algorithm SHA256 |
    Format-Table Hash, Path |
    Out-File checksums.txt

# Linux/macOS
sha256sum bin/* > checksums.txt
```

---

## 📊 Distribution Channels

### GitHub Releases (Primary)

**Pros:**

- ✅ Free hosting
- ✅ Version control integration
- ✅ Automatic changelog
- ✅ Download statistics

**Setup:**

1. Create release on GitHub
2. Upload all artifacts
3. Users download from Releases page

---

### Chocolatey (Windows Package Manager)

**Pros:**

- ✅ Easy updates: `choco upgrade obelisk-cli`
- ✅ Trusted by Windows users
- ✅ Automatic dependency management

**Setup:**

1. Create `obelisk-cli.nuspec`:

```xml
<?xml version="1.0"?>
<package xmlns="http://schemas.microsoft.com/packaging/2015/06/nuspec.xsd">
  <metadata>
    <id>obelisk-cli</id>
    <version>1.1.0</version>
    <title>Obelisk CLI</title>
    <authors>OneDev PH</authors>
    <projectUrl>https://github.com/Swif7ify/Obelisk-CLI</projectUrl>
    <licenseUrl>https://github.com/Swif7ify/Obelisk-CLI/blob/main/LICENSE</licenseUrl>
    <requireLicenseAcceptance>false</requireLicenseAcceptance>
    <description>AI-Powered Automated Tech Lead for Modern Codebases</description>
    <tags>cli developer-tools ai code-quality security</tags>
  </metadata>
</package>
```

2. Submit to Chocolatey:

```powershell
choco pack
choco push obelisk-cli.1.1.0.nupkg --source https://push.chocolatey.org/
```

3. Users install:

```powershell
choco install obelisk-cli
```

---

### Winget (Windows Package Manager)

**Pros:**

- ✅ Built into Windows 11
- ✅ Official Microsoft package manager
- ✅ Growing user base

**Setup:**

1. Fork https://github.com/microsoft/winget-pkgs
2. Create manifest in `manifests/o/OneDev/ObeliskCLI/1.1.0/`
3. Submit PR

4. Users install:

```powershell
winget install OneDev.ObeliskCLI
```

---

### Homebrew (macOS/Linux)

**Pros:**

- ✅ Popular on macOS
- ✅ Easy updates: `brew upgrade obelisk-cli`
- ✅ Trusted by developers

**Setup:**

1. Create formula:

```ruby
class ObeliskCli < Formula
  desc "AI-Powered Automated Tech Lead"
  homepage "https://github.com/Swif7ify/Obelisk-CLI"
  url "https://github.com/Swif7ify/Obelisk-CLI/archive/v1.1.0.tar.gz"
  sha256 "..."
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args
  end

  test do
    system "#{bin}/obelisk", "version"
  end
end
```

2. Submit to Homebrew
3. Users install:

```bash
brew install obelisk-cli
```

---

## 📈 Analytics & Metrics

### Track Downloads

- GitHub Release download counts
- Chocolatey package statistics
- Winget telemetry (opt-in)
- Homebrew analytics

### Monitor Usage

```go
// Optional: Add telemetry (with user consent)
// Track:
// - Installation method
// - OS/Architecture
// - Feature usage
// - Error rates
```

---

## 🔒 Security

### Code Signing (Recommended)

**Windows:**

```powershell
# Sign the executable
signtool sign /f certificate.pfx /p password /t http://timestamp.digicert.com obelisk.exe

# Sign the MSI
signtool sign /f certificate.pfx /p password /t http://timestamp.digicert.com ObeliskCLI-1.1.0-x64.msi
```

**macOS:**

```bash
# Sign and notarize
codesign --sign "Developer ID" obelisk
xcrun notarytool submit obelisk.zip
```

### Verify Downloads

Provide checksums for all releases:

```bash
# Users verify downloads
sha256sum -c checksums.txt
```

---

## 📝 Documentation

Update these files for each release:

- ✅ `README.md` - Installation instructions
- ✅ `CHANGELOG.md` - What's new
- ✅ `DISTRIBUTION.md` - This file
- ✅ `installer/README.md` - Installer docs
- ✅ GitHub Release notes

---

## 🎯 Recommended Distribution Strategy

### For v1.1.0 (Initial Release)

1. **Primary:** GitHub Releases with MSI installer
2. **Secondary:** PowerShell install script
3. **Tertiary:** Portable executable

### For v1.0.0 (Stable Release)

1. **Primary:** GitHub Releases + MSI
2. **Secondary:** Chocolatey package
3. **Tertiary:** Winget package
4. **Optional:** Homebrew formula

---

## 🤝 Contributing

To improve distribution:

1. Test installation on clean systems
2. Document any issues
3. Submit PRs with improvements
4. Help maintain package managers

---

Built with ❤️ by OneDev PH for the IBM Bob Hackathon 2026
