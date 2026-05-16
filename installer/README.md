# 🏛️ Obelisk CLI Installer

This directory contains the Windows MSI installer configuration for Obelisk CLI.

## 📦 What's Included

### Core Files

- **obelisk.wxs** - WiX Toolset installer definition (XML)
- **License.rtf** - End User License Agreement (EULA) in RTF format
- **build-installer.ps1** - Automated build script for creating the MSI

### Optional Assets (Recommended)

- **icon.ico** - Application icon (256x256 recommended)
- **banner.bmp** - Installer banner (493x58 pixels)
- **dialog.bmp** - Installer dialog background (493x312 pixels)

## 🚀 Building the Installer

### Prerequisites

1. **Install WiX Toolset 3.11+**

    ```powershell
    # Download from: https://wixtoolset.org/releases/
    # Or use Chocolatey:
    choco install wixtoolset
    ```

2. **Verify Installation**
    ```powershell
    # Check if WiX is in PATH
    $env:WIX
    # Should output something like: C:\Program Files (x86)\WiX Toolset v3.11\
    ```

### Build Steps

```powershell
# From the project root directory:

# 1. Build the MSI installer (includes Go build)
.\installer\build-installer.ps1

# 2. Or skip Go build if already built
.\installer\build-installer.ps1 -SkipBuild

# 3. Specify custom version
.\installer\build-installer.ps1 -Version "0.2.0"
```

### Output

The script will create:

```
installer/ObeliskCLI-0.1.0-x64.msi
```

## 📋 Installer Features

### ✅ What the Installer Does

1. **License Agreement** - Shows EULA that users must accept
2. **Installation Directory** - Allows users to choose install location (default: `C:\Program Files\Obelisk CLI`)
3. **PATH Configuration** - Automatically adds Obelisk to system PATH
4. **Start Menu Shortcuts** - Creates uninstall shortcut
5. **Documentation** - Installs README, LICENSE, and CHANGELOG
6. **Upgrade Support** - Handles upgrades from previous versions
7. **Clean Uninstall** - Removes all files and PATH entries

### 🎨 Customization

#### Change Product Information

Edit `installer/obelisk.wxs`:

```xml
<Product Id="*"
         Name="Obelisk CLI"
         Version="0.1.0"
         Manufacturer="OneDev PH">
```

#### Add Custom Icons

Place these files in the `installer/` directory:

- `icon.ico` - 256x256 ICO file
- `banner.bmp` - 493x58 BMP (top banner)
- `dialog.bmp` - 493x312 BMP (left side image)

#### Modify License Text

Edit `installer/License.rtf` with WordPad or any RTF editor.

## 📤 Distribution

### GitHub Releases (Recommended)

1. **Create a Release**

    ```bash
    git tag v0.1.0
    git push origin v0.1.0
    ```

2. **Upload MSI**
    - Go to GitHub Releases
    - Create new release from tag
    - Upload `ObeliskCLI-0.1.0-x64.msi`
    - Add release notes

3. **Users Install**
    ```powershell
    # Download from GitHub Releases
    # Double-click ObeliskCLI-0.1.0-x64.msi
    # Follow installer wizard
    ```

### Alternative: Chocolatey Package

Create a Chocolatey package for easier distribution:

```powershell
choco pack
choco push obelisk-cli.0.1.0.nupkg --source https://push.chocolatey.org/
```

### Alternative: Winget

Submit to Windows Package Manager:

```yaml
# Create winget manifest
winget install OneDev.ObeliskCLI
```

## 🔧 Troubleshooting

### "WiX Toolset not found"

- Install WiX from https://wixtoolset.org/releases/
- Restart PowerShell after installation
- Verify `$env:WIX` is set

### "Build failed"

- Ensure `main.go` exists in project root
- Run `go mod tidy` to fix dependencies
- Check Go version: `go version` (requires 1.21+)

### "MSI linking failed"

- Check all file paths in `obelisk.wxs` are correct
- Ensure `bin/obelisk.exe` exists
- Verify `License.rtf` is valid RTF format

### "Installer won't run"

- Right-click MSI → Properties → Unblock
- Run as Administrator if needed
- Check Windows Defender/Antivirus

## 📝 Technical Details

### MSI Properties

- **Product Code**: Auto-generated GUID (changes per version)
- **Upgrade Code**: Fixed GUID (same across versions)
- **Install Scope**: Per-machine (requires admin)
- **Platform**: x64 Windows

### Registry Keys

```
HKLM\Software\OneDev\Obelisk
HKCU\Software\OneDev\Obelisk
```

### Installed Files

```
C:\Program Files\Obelisk CLI\
├── obelisk.exe
├── README.md
├── LICENSE
└── CHANGELOG.md
```

### PATH Entry

```
C:\Program Files\Obelisk CLI
```

## 🤝 Contributing

To improve the installer:

1. Edit `obelisk.wxs` for installer logic
2. Update `License.rtf` for legal changes
3. Modify `build-installer.ps1` for build process
4. Test on clean Windows VM
5. Submit PR with changes

## 📚 Resources

- [WiX Toolset Documentation](https://wixtoolset.org/documentation/)
- [WiX Tutorial](https://www.firegiant.com/wix/tutorial/)
- [MSI Best Practices](https://docs.microsoft.com/en-us/windows/win32/msi/windows-installer-best-practices)

---

Built by OneDev PH for the IBM Bob Hackathon 2026
