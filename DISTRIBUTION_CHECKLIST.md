# Distribution Checklist for Obelisk CLI

Complete this checklist to prepare Obelisk CLI for distribution.

## Phase 1: Pre-Distribution Setup (REQUIRED)

### 1.1 Install Required Tools

- [ ] **Install WiX Toolset 3.11+**
    ```powershell
    choco install wixtoolset
    # OR download from: https://wixtoolset.org/releases/
    ```
- [ ] **Verify WiX Installation**
    ```powershell
    $env:WIX  # Should output: C:\Program Files (x86)\WiX Toolset v3.11\
    ```

### 1.2 Build and Test Locally

- [ ] **Build the executable**
    ```powershell
    go build -o bin/obelisk.exe .
    ```
- [ ] **Test the executable**
    ```powershell
    .\bin\obelisk.exe version
    .\bin\obelisk.exe --help
    ```
- [ ] **Run tests**
    ```powershell
    go test ./...
    ```

### 1.3 Build MSI Installer

- [ ] **Run the installer build script**
    ```powershell
    .\installer\build-installer.ps1
    ```
- [ ] **Verify MSI was created**
    ```powershell
    # Check: installer/ObeliskCLI-0.1.0-x64.msi exists
    ```

### 1.4 Test MSI Installer

- [ ] **Install MSI on test machine**
    - Double-click `ObeliskCLI-0.1.0-x64.msi`
    - Follow installation wizard
    - Accept license agreement
    - Choose installation directory
    - Complete installation
- [ ] **Verify installation**

    ```powershell
    # Open NEW terminal window
    obelisk version
    obelisk --help
    ```

- [ ] **Test uninstallation**
    - Go to Windows Settings > Apps
    - Find "Obelisk CLI"
    - Click Uninstall
    - Verify it's removed

---

## Phase 2: GitHub Release (REQUIRED)

### 2.1 Prepare Release

- [ ] **Update version numbers**
    - [ ] `main.go` - Version constant
    - [ ] `installer/obelisk.wxs` - Version attribute
    - [ ] `README.md` - Version badge
    - [ ] `CHANGELOG.md` - Add release notes

- [ ] **Commit version changes**
    ```bash
    git add .
    git commit -m "chore: bump version to 0.1.0"
    git push
    ```

### 2.2 Build All Platform Binaries

- [ ] **Build for Windows**

    ```powershell
    $env:GOOS = "windows"; $env:GOARCH = "amd64"
    go build -o bin/obelisk-windows-amd64.exe -ldflags "-s -w -X main.Version=0.1.0" .
    ```

- [ ] **Build for Linux**

    ```powershell
    $env:GOOS = "linux"; $env:GOARCH = "amd64"
    go build -o bin/obelisk-linux-amd64 -ldflags "-s -w -X main.Version=0.1.0" .
    ```

- [ ] **Build for macOS Intel**

    ```powershell
    $env:GOOS = "darwin"; $env:GOARCH = "amd64"
    go build -o bin/obelisk-darwin-amd64 -ldflags "-s -w -X main.Version=0.1.0" .
    ```

- [ ] **Build for macOS Apple Silicon**
    ```powershell
    $env:GOOS = "darwin"; $env:GOARCH = "arm64"
    go build -o bin/obelisk-darwin-arm64 -ldflags "-s -w -X main.Version=0.1.0" .
    ```

### 2.3 Generate Checksums

- [ ] **Calculate SHA256 hashes**
    ```powershell
    Get-FileHash bin/*.exe, bin/obelisk-*, installer/*.msi -Algorithm SHA256 |
        Format-Table Hash, @{Label="File";Expression={Split-Path $_.Path -Leaf}} |
        Out-File checksums.txt
    ```

### 2.4 Create GitHub Release

- [ ] **Create and push tag**

    ```bash
    git tag -a v0.1.0 -m "Release version 0.1.0"
    git push origin v0.1.0
    ```

- [ ] **Create release on GitHub**
    1. Go to: https://github.com/Swif7ify/Obelisk-CLI/releases
    2. Click "Create new release"
    3. Select tag: `v0.1.0`
    4. Release title: `Obelisk CLI v0.1.0`
    5. Description: Copy from CHANGELOG.md

- [ ] **Upload release files**
    - [ ] `ObeliskCLI-0.1.0-x64.msi` (Windows Installer - PRIMARY)
    - [ ] `obelisk-windows-amd64.exe` (Windows Portable)
    - [ ] `obelisk-linux-amd64` (Linux Binary)
    - [ ] `obelisk-darwin-amd64` (macOS Intel)
    - [ ] `obelisk-darwin-arm64` (macOS Apple Silicon)
    - [ ] `checksums.txt` (SHA256 Hashes)

- [ ] **Publish release**

---

## Phase 3: Package Managers (OPTIONAL - Can do later)

### 3.1 Chocolatey (Windows)

- [ ] **Update checksums in package**

    ```powershell
    # Calculate checksum
    $hash = (Get-FileHash bin/obelisk-windows-amd64.exe -Algorithm SHA256).Hash

    # Update packages/chocolatey/tools/chocolateyinstall.ps1
    # Replace: checksum64 = '' with checksum64 = '$hash'
    ```

- [ ] **Build Chocolatey package**

    ```powershell
    cd packages/chocolatey
    choco pack
    ```

- [ ] **Test locally**

    ```powershell
    choco install obelisk-cli -source . -force
    obelisk version
    choco uninstall obelisk-cli
    ```

- [ ] **Submit to Chocolatey**
    ```powershell
    choco apikey --key YOUR_API_KEY --source https://push.chocolatey.org/
    choco push obelisk-cli.0.1.0.nupkg --source https://push.chocolatey.org/
    ```

### 3.2 Winget (Windows)

- [ ] **Calculate MSI checksum**

    ```powershell
    $hash = (Get-FileHash installer/ObeliskCLI-0.1.0-x64.msi -Algorithm SHA256).Hash
    # Update packages/winget/OneDev.ObeliskCLI.installer.yaml
    ```

- [ ] **Extract Product Code**

    ```powershell
    # See packages/winget/README.md for instructions
    # Update packages/winget/OneDev.ObeliskCLI.installer.yaml
    ```

- [ ] **Use wingetcreate (easier)**

    ```powershell
    winget install wingetcreate
    wingetcreate update OneDev.ObeliskCLI --version 0.1.0 --urls https://github.com/Swif7ify/Obelisk-CLI/releases/download/v0.1.0/ObeliskCLI-0.1.0-x64.msi
    ```

- [ ] **Submit to winget-pkgs**
    - Fork: https://github.com/microsoft/winget-pkgs
    - Create PR with manifest files

### 3.3 Homebrew (macOS/Linux)

- [ ] **Calculate tarball checksum**

    ```bash
    curl -L https://github.com/Swif7ify/Obelisk-CLI/archive/refs/tags/v0.1.0.tar.gz | shasum -a 256
    # Update packages/homebrew/obelisk-cli.rb
    ```

- [ ] **Create Homebrew tap**
    - Create repo: `homebrew-obelisk`
    - Add formula to `Formula/obelisk-cli.rb`
    - Users install: `brew tap swif7ify/obelisk && brew install obelisk-cli`

---

## Phase 4: Optional Enhancements

### 4.1 Code Signing (Recommended)

- [ ] **Purchase code signing certificate**
    - DigiCert: $474/year
    - Sectigo: $199/year
    - SSL.com: $199/year

- [ ] **Sign executable**

    ```powershell
    signtool sign /f certificate.pfx /p password /fd SHA256 /tr http://timestamp.digicert.com /td SHA256 bin/obelisk.exe
    ```

- [ ] **Sign MSI**
    ```powershell
    signtool sign /f certificate.pfx /p password /fd SHA256 /tr http://timestamp.digicert.com /td SHA256 installer/ObeliskCLI-0.1.0-x64.msi
    ```

### 4.2 Installer Assets

- [ ] **Create icon.ico** (256x256)
    - See: `installer/ASSETS.md`
    - Place in: `installer/icon.ico`

- [ ] **Create banner.bmp** (493x58)
    - See: `installer/ASSETS.md`
    - Place in: `installer/banner.bmp`

- [ ] **Create dialog.bmp** (493x312)
    - See: `installer/ASSETS.md`
    - Place in: `installer/dialog.bmp`

### 4.3 Documentation

- [ ] **Update README.md**
    - Installation instructions
    - Usage examples
    - Screenshots

- [ ] **Create video demo**
    - Record usage demo
    - Upload to YouTube
    - Add link to README

---

## Phase 5: Post-Release

### 5.1 Announce Release

- [ ] **Social media**
    - Twitter/X
    - LinkedIn
    - Reddit (r/golang, r/programming)

- [ ] **Dev.to article**
    - Write about the project
    - Share lessons learned

- [ ] **Hacker News**
    - Submit Show HN post

### 5.2 Monitor

- [ ] **Watch GitHub issues**
- [ ] **Monitor download stats**
- [ ] **Respond to feedback**

---

## Quick Start (Minimum Viable Distribution)

If you just want to get started quickly, do these steps:

1. [ ] Install WiX Toolset
2. [ ] Build MSI: `.\installer\build-installer.ps1`
3. [ ] Test MSI locally
4. [ ] Create GitHub Release
5. [ ] Upload MSI to release
6. [ ] Done! Users can download and install

Everything else can be added later as your project grows!

---

## Current Status

**Ready for Distribution:** YES / NO

**Blockers:**

- [ ] None - Ready to go!
- [ ] Need to install WiX Toolset
- [ ] Need to test MSI installer
- [ ] Need to create GitHub release
- [ ] Other: ******\_\_\_******

**Next Step:** ******\_\_\_******

---

## Support

- **MSI Installer:** `installer/README.md`
- **Chocolatey:** `packages/chocolatey/README.md`
- **Winget:** `packages/winget/README.md`
- **Homebrew:** `packages/homebrew/README.md`
- **Distribution Guide:** `DISTRIBUTION.md`
- **Code Signing:** `installer/CODE_SIGNING.md`
