# 🍫 Chocolatey Package for Obelisk CLI

This directory contains the Chocolatey package definition for Obelisk CLI.

## 📦 What is Chocolatey?

Chocolatey is a package manager for Windows, similar to apt-get or Homebrew. It allows users to install software with a single command:

```powershell
choco install obelisk-cli
```

## 🚀 Building the Package

### Prerequisites

1. **Install Chocolatey:**

    ```powershell
    Set-ExecutionPolicy Bypass -Scope Process -Force
    [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072
    iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))
    ```

2. **Verify Installation:**
    ```powershell
    choco --version
    ```

### Build Steps

1. **Update Version:**
   Edit `obelisk-cli.nuspec` and update the version number:

    ```xml
    <version>0.1.0</version>
    ```

2. **Update Download URL:**
   Edit `tools/chocolateyinstall.ps1` and update the URL:

    ```powershell
    $url64 = 'https://github.com/Swif7ify/Obelisk-CLI/releases/download/v0.1.0/obelisk-windows-amd64.exe'
    ```

3. **Calculate Checksum:**

    ```powershell
    # Download the release file
    $url = "https://github.com/Swif7ify/Obelisk-CLI/releases/download/v0.1.0/obelisk-windows-amd64.exe"
    Invoke-WebRequest -Uri $url -OutFile "obelisk.exe"

    # Calculate SHA256
    $checksum = (Get-FileHash "obelisk.exe" -Algorithm SHA256).Hash
    Write-Host "Checksum: $checksum"
    ```

4. **Update Checksum:**
   Edit `tools/chocolateyinstall.ps1`:

    ```powershell
    checksum64 = 'YOUR_CHECKSUM_HERE'
    ```

5. **Build Package:**

    ```powershell
    cd packages/chocolatey
    choco pack
    ```

    This creates: `obelisk-cli.0.1.0.nupkg`

## 🧪 Testing Locally

```powershell
# Install from local package
choco install obelisk-cli -source . -force

# Test the installation
obelisk version

# Uninstall
choco uninstall obelisk-cli
```

## 📤 Publishing to Chocolatey

### 1. Create Chocolatey Account

1. Go to https://community.chocolatey.org/
2. Sign up for an account
3. Get your API key from your profile

### 2. Set API Key

```powershell
choco apikey --key YOUR_API_KEY --source https://push.chocolatey.org/
```

### 3. Push Package

```powershell
choco push obelisk-cli.0.1.0.nupkg --source https://push.chocolatey.org/
```

### 4. Wait for Approval

- Chocolatey moderators will review your package
- Usually takes 1-3 days for first submission
- You'll receive email notifications
- Once approved, users can install with `choco install obelisk-cli`

## 📋 Package Validation

Before submitting, validate your package:

```powershell
# Install Chocolatey Package Validator
choco install chocolatey-package-validator

# Validate package
choco-package-validator obelisk-cli.0.1.0.nupkg
```

## 🔄 Updating the Package

For new versions:

1. Update version in `obelisk-cli.nuspec`
2. Update URL in `tools/chocolateyinstall.ps1`
3. Calculate new checksum
4. Update checksum in `tools/chocolateyinstall.ps1`
5. Build: `choco pack`
6. Test locally
7. Push: `choco push obelisk-cli.X.Y.Z.nupkg`

## 📝 Package Structure

```
packages/chocolatey/
├── obelisk-cli.nuspec              # Package metadata
├── tools/
│   ├── chocolateyinstall.ps1       # Installation script
│   └── chocolateyuninstall.ps1     # Uninstallation script
└── README.md                        # This file
```

## 🎯 Best Practices

### Version Numbers

- Follow semantic versioning (X.Y.Z)
- Match GitHub release versions
- Update for every release

### Checksums

- Always include SHA256 checksums
- Verify checksums before publishing
- Never skip checksum validation

### Testing

- Test on clean Windows VM
- Test both install and uninstall
- Verify PATH is set correctly
- Check executable works after install

### Documentation

- Keep README.md updated
- Document breaking changes
- Include usage examples

## 🆘 Troubleshooting

### "Package already exists"

```powershell
# Use --force to overwrite during testing
choco install obelisk-cli -source . -force
```

### "Checksum mismatch"

```powershell
# Recalculate checksum
Get-FileHash obelisk.exe -Algorithm SHA256
```

### "Package validation failed"

```powershell
# Run validator for detailed errors
choco-package-validator obelisk-cli.0.1.0.nupkg
```

### "API key not set"

```powershell
# Set API key again
choco apikey --key YOUR_API_KEY --source https://push.chocolatey.org/
```

## 📚 Resources

- [Chocolatey Documentation](https://docs.chocolatey.org/)
- [Package Creation Guide](https://docs.chocolatey.org/en-us/create/create-packages)
- [Package Validator](https://github.com/chocolatey/package-validator)
- [Community Repository](https://community.chocolatey.org/)

## 🤝 Maintenance

### Automated Updates

Consider setting up automatic package updates:

1. **AU (Automatic Updater):**

    ```powershell
    # Install AU
    choco install au

    # Create update script
    # See: https://github.com/majkinetor/au
    ```

2. **GitHub Actions:**
    - Trigger on new releases
    - Auto-update checksums
    - Auto-push to Chocolatey

## 📊 Statistics

Once published, track your package:

- https://community.chocolatey.org/packages/obelisk-cli
- View download counts
- Monitor user feedback
- Track version adoption

---

**Note:** First-time package submissions require manual approval. Subsequent updates are usually automatic if you maintain good standing.
