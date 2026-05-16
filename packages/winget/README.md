# 📦 Winget Package for Obelisk CLI

This directory contains the Windows Package Manager (winget) manifest for Obelisk CLI.

## 🎯 What is Winget?

Winget is Microsoft's official package manager for Windows 10/11. It's built into Windows and allows users to install software with:

```powershell
winget install OneDev.ObeliskCLI
```

## 📋 Manifest Files

Winget requires three YAML files:

1. **OneDev.ObeliskCLI.yaml** - Version manifest
2. **OneDev.ObeliskCLI.installer.yaml** - Installer details
3. **OneDev.ObeliskCLI.locale.en-US.yaml** - Package metadata

## 🚀 Preparing for Submission

### Step 1: Build and Upload MSI

First, create a GitHub release with your MSI:

```powershell
# Build the MSI
.\installer\build-installer.ps1

# Create GitHub release and upload:
# - ObeliskCLI-0.1.0-x64.msi
```

### Step 2: Calculate SHA256

```powershell
# Calculate checksum of the MSI
$msiPath = "installer\ObeliskCLI-0.1.0-x64.msi"
$hash = (Get-FileHash $msiPath -Algorithm SHA256).Hash
Write-Host "SHA256: $hash"
```

### Step 3: Get Product Code

```powershell
# Extract Product Code from MSI
$windowsInstaller = New-Object -ComObject WindowsInstaller.Installer
$database = $windowsInstaller.GetType().InvokeMember("OpenDatabase", "InvokeMethod", $null, $windowsInstaller, @("$msiPath", 0))
$view = $database.GetType().InvokeMember("OpenView", "InvokeMethod", $null, $database, ("SELECT Value FROM Property WHERE Property='ProductCode'"))
$view.GetType().InvokeMember("Execute", "InvokeMethod", $null, $view, $null)
$record = $view.GetType().InvokeMember("Fetch", "InvokeMethod", $null, $view, $null)
$productCode = $record.GetType().InvokeMember("StringData", "GetProperty", $null, $record, 1)
Write-Host "Product Code: $productCode"
```

### Step 4: Update Manifest Files

Edit `OneDev.ObeliskCLI.installer.yaml`:

```yaml
Installers:
    - Architecture: x64
      InstallerUrl: https://github.com/Swif7ify/Obelisk-CLI/releases/download/v0.1.0/ObeliskCLI-0.1.0-x64.msi
      InstallerSha256: YOUR_SHA256_HERE # From Step 2
      ProductCode: "{YOUR_PRODUCT_CODE_HERE}" # From Step 3
```

## 📤 Submitting to Winget

### Method 1: Using Winget-Create (Recommended)

```powershell
# Install winget-create
winget install wingetcreate

# Create/update manifest
wingetcreate update OneDev.ObeliskCLI --version 0.1.0 --urls https://github.com/Swif7ify/Obelisk-CLI/releases/download/v0.1.0/ObeliskCLI-0.1.0-x64.msi

# This will:
# 1. Download the MSI
# 2. Calculate SHA256 automatically
# 3. Extract Product Code
# 4. Generate manifests
# 5. Create PR to winget-pkgs repo
```

### Method 2: Manual Submission

1. **Fork the winget-pkgs repository:**

    ```bash
    # Go to: https://github.com/microsoft/winget-pkgs
    # Click "Fork"
    ```

2. **Clone your fork:**

    ```bash
    git clone https://github.com/YOUR_USERNAME/winget-pkgs.git
    cd winget-pkgs
    ```

3. **Create package directory:**

    ```bash
    mkdir -p manifests/o/OneDev/ObeliskCLI/0.1.0
    ```

4. **Copy manifest files:**

    ```bash
    cp packages/winget/*.yaml manifests/o/OneDev/ObeliskCLI/0.1.0/
    ```

5. **Validate manifests:**

    ```powershell
    # Install validation tool
    winget install winget-pkgs-validation

    # Validate
    winget validate --manifest manifests/o/OneDev/ObeliskCLI/0.1.0
    ```

6. **Create PR:**

    ```bash
    git checkout -b obelisk-cli-0.1.0
    git add manifests/o/OneDev/ObeliskCLI/0.1.0/
    git commit -m "New package: OneDev.ObeliskCLI version 0.1.0"
    git push origin obelisk-cli-0.1.0

    # Go to GitHub and create Pull Request
    ```

## ✅ Validation Checklist

Before submitting:

- [ ] MSI is uploaded to GitHub Releases
- [ ] SHA256 checksum is correct
- [ ] Product Code is extracted from MSI
- [ ] All three YAML files are present
- [ ] Package identifier follows format: `Publisher.PackageName`
- [ ] Version matches GitHub release tag
- [ ] URLs are accessible
- [ ] Manifests pass validation: `winget validate`

## 🔄 Updating for New Versions

For version 0.2.0:

1. **Create new directory:**

    ```bash
    mkdir -p manifests/o/OneDev/ObeliskCLI/0.2.0
    ```

2. **Copy and update manifests:**
    - Update version numbers
    - Update InstallerUrl
    - Calculate new SHA256
    - Extract new Product Code

3. **Submit PR:**
    ```bash
    git checkout -b obelisk-cli-0.2.0
    git add manifests/o/OneDev/ObeliskCLI/0.2.0/
    git commit -m "Update: OneDev.ObeliskCLI version 0.2.0"
    git push origin obelisk-cli-0.2.0
    ```

## 🧪 Testing

After PR is merged:

```powershell
# Search for package
winget search ObeliskCLI

# Show package info
winget show OneDev.ObeliskCLI

# Install
winget install OneDev.ObeliskCLI

# Test
obelisk version

# Uninstall
winget uninstall OneDev.ObeliskCLI
```

## 📊 Review Process

1. **Automated Checks:**
    - Manifest validation
    - URL accessibility
    - SHA256 verification
    - Installer testing

2. **Manual Review:**
    - Microsoft team reviews PR
    - Usually takes 1-3 days
    - May request changes

3. **Approval:**
    - PR is merged
    - Package appears in winget within hours
    - Users can install immediately

## 🆘 Troubleshooting

### "Manifest validation failed"

```powershell
# Run validator for details
winget validate --manifest manifests/o/OneDev/ObeliskCLI/0.1.0
```

### "SHA256 mismatch"

```powershell
# Recalculate checksum
Get-FileHash installer\ObeliskCLI-0.1.0-x64.msi -Algorithm SHA256
```

### "Product Code not found"

```powershell
# Use Orca MSI editor (part of Windows SDK)
# Or use the PowerShell script above
```

### "URL not accessible"

- Ensure GitHub release is public
- Verify URL is correct
- Check file name matches exactly

## 📚 Resources

- [Winget Documentation](https://docs.microsoft.com/en-us/windows/package-manager/)
- [Manifest Schema](https://github.com/microsoft/winget-pkgs/tree/master/doc/manifest)
- [Submission Guidelines](https://github.com/microsoft/winget-pkgs/blob/master/CONTRIBUTING.md)
- [Winget-Create Tool](https://github.com/microsoft/winget-create)

## 🎯 Benefits of Winget

- ✅ Built into Windows 10/11
- ✅ Official Microsoft package manager
- ✅ No installation required
- ✅ Automatic updates
- ✅ Trusted by enterprises
- ✅ Growing user base

## 📝 Notes

- First submission requires manual review
- Subsequent updates are usually automatic
- Package identifier cannot be changed
- Maintain consistent versioning
- Keep manifests in sync with releases

---

**Tip:** Use `wingetcreate` tool for easiest submission process!
