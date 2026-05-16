# 🔐 Code Signing Guide for Obelisk CLI

Code signing prevents Windows SmartScreen warnings and builds user trust. This guide covers the complete process.

## 🎯 Why Code Sign?

### Without Code Signing:

```
⚠️ Windows protected your PC
   Microsoft Defender SmartScreen prevented an unrecognized app from starting.
   Running this app might put your PC at risk.

   [Don't run]  [More info]
```

### With Code Signing:

```
✅ Verified publisher: OneDev PH
   This app has been verified by the publisher.

   [Install]  [Cancel]
```

**Benefits:**

- ✅ No SmartScreen warnings
- ✅ Builds user trust
- ✅ Professional appearance
- ✅ Required for some distribution channels
- ✅ Proves authenticity

---

## 📋 Prerequisites

### 1. Get a Code Signing Certificate

#### Option A: Commercial Certificate Authority (Recommended)

**Trusted CAs:**

- **DigiCert** - $474/year (most trusted)
- **Sectigo** - $199/year (good value)
- **GlobalSign** - $249/year
- **SSL.com** - $199/year

**What You Need:**

- Business registration documents
- Tax ID / EIN
- Business address verification
- Phone verification
- Email verification

**Process:**

1. Choose a CA and certificate type
2. Submit business verification documents
3. Wait 1-7 days for verification
4. Receive certificate (PFX file + password)
5. Install on signing machine

#### Option B: Self-Signed Certificate (Testing Only)

**⚠️ Warning:** Self-signed certificates still trigger SmartScreen warnings. Only use for testing!

```powershell
# Create self-signed certificate (testing only)
$cert = New-SelfSignedCertificate `
    -Type CodeSigningCert `
    -Subject "CN=OneDev PH, O=OneDev PH, C=PH" `
    -KeyAlgorithm RSA `
    -KeyLength 2048 `
    -Provider "Microsoft Enhanced RSA and AES Cryptographic Provider" `
    -CertStoreLocation "Cert:\CurrentUser\My" `
    -NotAfter (Get-Date).AddYears(3)

# Export to PFX
$password = ConvertTo-SecureString -String "YourPassword123!" -Force -AsPlainText
Export-PfxCertificate -Cert $cert -FilePath "test-cert.pfx" -Password $password

# ⚠️ This is for TESTING ONLY - users will still see warnings!
```

---

## 🔧 Signing Process

### 1. Install Certificate

```powershell
# Import PFX certificate
$password = Read-Host -AsSecureString "Enter certificate password"
Import-PfxCertificate -FilePath "certificate.pfx" -CertStoreLocation Cert:\CurrentUser\My -Password $password

# Verify installation
Get-ChildItem Cert:\CurrentUser\My -CodeSigningCert
```

### 2. Sign the Executable

```powershell
# Sign obelisk.exe
$cert = Get-ChildItem Cert:\CurrentUser\My -CodeSigningCert | Select-Object -First 1

Set-AuthenticodeSignature `
    -FilePath "bin\obelisk.exe" `
    -Certificate $cert `
    -TimestampServer "http://timestamp.digicert.com" `
    -HashAlgorithm SHA256

# Verify signature
Get-AuthenticodeSignature "bin\obelisk.exe" | Format-List
```

### 3. Sign the MSI Installer

```powershell
# Using signtool.exe (part of Windows SDK)
# Download from: https://developer.microsoft.com/en-us/windows/downloads/windows-sdk/

# Sign MSI
signtool sign /f certificate.pfx /p YourPassword /fd SHA256 /tr http://timestamp.digicert.com /td SHA256 "installer\ObeliskCLI-0.1.0-x64.msi"

# Verify signature
signtool verify /pa "installer\ObeliskCLI-0.1.0-x64.msi"
```

---

## 🤖 Automated Signing Script

Create `installer/sign-release.ps1`:

```powershell
# Automated Code Signing Script
param(
    [Parameter(Mandatory=$true)]
    [string]$CertPath,

    [Parameter(Mandatory=$true)]
    [SecureString]$CertPassword,

    [string]$TimestampServer = "http://timestamp.digicert.com"
)

$ErrorActionPreference = "Stop"

Write-Host "🔐 Obelisk CLI Code Signing" -ForegroundColor Cyan
Write-Host ""

# Check if signtool exists
$signtool = "C:\Program Files (x86)\Windows Kits\10\bin\10.0.22621.0\x64\signtool.exe"
if (-not (Test-Path $signtool)) {
    Write-Host "❌ signtool.exe not found!" -ForegroundColor Red
    Write-Host "   Install Windows SDK: https://developer.microsoft.com/windows/downloads/windows-sdk/" -ForegroundColor Yellow
    exit 1
}

# Convert SecureString to plain text for signtool
$BSTR = [System.Runtime.InteropServices.Marshal]::SecureStringToBSTR($CertPassword)
$password = [System.Runtime.InteropServices.Marshal]::PtrToStringAuto($BSTR)

# Sign executable
Write-Host "Signing obelisk.exe..." -ForegroundColor Yellow
& $signtool sign /f $CertPath /p $password /fd SHA256 /tr $TimestampServer /td SHA256 "bin\obelisk.exe"

if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Failed to sign executable!" -ForegroundColor Red
    exit 1
}
Write-Host "✓ Executable signed" -ForegroundColor Green

# Sign MSI
Write-Host "Signing MSI installer..." -ForegroundColor Yellow
& $signtool sign /f $CertPath /p $password /fd SHA256 /tr $TimestampServer /td SHA256 "installer\ObeliskCLI-0.1.0-x64.msi"

if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Failed to sign MSI!" -ForegroundColor Red
    exit 1
}
Write-Host "✓ MSI signed" -ForegroundColor Green

# Verify signatures
Write-Host ""
Write-Host "Verifying signatures..." -ForegroundColor Yellow

& $signtool verify /pa "bin\obelisk.exe"
& $signtool verify /pa "installer\ObeliskCLI-0.1.0-x64.msi"

Write-Host ""
Write-Host "✅ All files signed successfully!" -ForegroundColor Green
```

**Usage:**

```powershell
$password = Read-Host -AsSecureString "Certificate password"
.\installer\sign-release.ps1 -CertPath "certificate.pfx" -CertPassword $password
```

---

## 🔄 CI/CD Integration

### GitHub Actions

```yaml
name: Build and Sign Release

on:
    push:
        tags:
            - "v*"

jobs:
    build-and-sign:
        runs-on: windows-latest
        steps:
            - uses: actions/checkout@v3

            - name: Setup Go
              uses: actions/setup-go@v4
              with:
                  go-version: "1.21"

            - name: Build executable
              run: go build -o bin/obelisk.exe

            - name: Download certificate
              run: |
                  echo "${{ secrets.CERT_BASE64 }}" | base64 -d > certificate.pfx

            - name: Sign executable
              run: |
                  $password = ConvertTo-SecureString -String "${{ secrets.CERT_PASSWORD }}" -Force -AsPlainText
                  .\installer\sign-release.ps1 -CertPath certificate.pfx -CertPassword $password

            - name: Upload signed artifacts
              uses: actions/upload-artifact@v3
              with:
                  name: signed-release
                  path: |
                      bin/obelisk.exe
                      installer/*.msi
```

**Setup:**

1. Convert certificate to base64: `[Convert]::ToBase64String([IO.File]::ReadAllBytes("certificate.pfx"))`
2. Add to GitHub Secrets:
    - `CERT_BASE64` - Base64 encoded certificate
    - `CERT_PASSWORD` - Certificate password

---

## 🛡️ Security Best Practices

### Certificate Storage

**✅ DO:**

- Store certificate in secure location
- Use strong password
- Limit access to signing machine
- Use hardware security module (HSM) for production
- Rotate certificates before expiration

**❌ DON'T:**

- Commit certificate to Git
- Share certificate password
- Use weak passwords
- Store in cloud without encryption
- Use expired certificates

### Timestamp Servers

Always use timestamp servers:

```powershell
# Recommended timestamp servers
http://timestamp.digicert.com
http://timestamp.sectigo.com
http://timestamp.globalsign.com
```

**Why?** Timestamps allow signatures to remain valid after certificate expires.

---

## 🧪 Testing Signed Binaries

### Verify Signature

```powershell
# Check signature details
Get-AuthenticodeSignature "bin\obelisk.exe" | Format-List

# Should show:
# Status: Valid
# SignerCertificate: CN=OneDev PH
# TimeStamperCertificate: CN=DigiCert
```

### Test on Clean System

1. Copy signed files to clean Windows VM
2. Double-click MSI installer
3. Should NOT show SmartScreen warning
4. Should show "Verified publisher: OneDev PH"

---

## 💰 Cost Comparison

| Provider    | Price/Year | Validation   | Trust Level    |
| ----------- | ---------- | ------------ | -------------- |
| DigiCert    | $474       | Organization | Highest        |
| Sectigo     | $199       | Organization | High           |
| SSL.com     | $199       | Organization | High           |
| GlobalSign  | $249       | Organization | High           |
| Self-Signed | Free       | None         | ⚠️ Not Trusted |

**Recommendation:** Start with Sectigo or SSL.com for best value.

---

## 📝 Certificate Renewal

Certificates expire after 1-3 years:

1. **60 days before expiration:**
    - Order renewal from CA
    - Update business information if changed

2. **30 days before expiration:**
    - Receive new certificate
    - Test signing with new cert
    - Update CI/CD secrets

3. **Before expiration:**
    - Sign new releases with new cert
    - Old signatures remain valid (if timestamped)

---

## 🆘 Troubleshooting

### "Certificate not found"

```powershell
# List all certificates
Get-ChildItem Cert:\CurrentUser\My -CodeSigningCert

# Import if missing
Import-PfxCertificate -FilePath "certificate.pfx" -CertStoreLocation Cert:\CurrentUser\My
```

### "Timestamp server unavailable"

```powershell
# Try different timestamp server
-TimestampServer "http://timestamp.sectigo.com"
```

### "Invalid certificate"

- Check certificate hasn't expired
- Verify it's a code signing certificate
- Ensure proper key usage extensions

### SmartScreen still shows warning

- Certificate needs to build reputation (takes time)
- Ensure certificate is from trusted CA
- Self-signed certificates always show warnings

---

## 📚 Resources

- [Microsoft Code Signing Guide](https://docs.microsoft.com/en-us/windows/win32/seccrypto/cryptography-tools)
- [DigiCert Code Signing](https://www.digicert.com/signing/code-signing-certificates)
- [Windows SDK Download](https://developer.microsoft.com/en-us/windows/downloads/windows-sdk/)
- [SignTool Documentation](https://docs.microsoft.com/en-us/windows/win32/seccrypto/signtool)

---

**Note:** Code signing is optional but highly recommended for professional distribution. Start without it and add later as your project grows!
