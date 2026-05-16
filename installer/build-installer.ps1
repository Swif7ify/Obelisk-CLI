# Obelisk CLI MSI Installer Build Script
# Requires WiX Toolset 3.11+ to be installed
# Download from: https://wixtoolset.org/releases/

param(
    [string]$Version = "0.1.0",
    [switch]$SkipBuild = $false
)

$ErrorActionPreference = "Stop"

Write-Host "🏛️  Obelisk CLI Installer Builder" -ForegroundColor Cyan
Write-Host "=================================" -ForegroundColor Cyan
Write-Host ""

# Check if WiX is installed
$wixPath = "${env:WIX}bin"
if (-not (Test-Path $wixPath)) {
    Write-Host "❌ WiX Toolset not found!" -ForegroundColor Red
    Write-Host ""
    Write-Host "Please install WiX Toolset 3.11 or later:" -ForegroundColor Yellow
    Write-Host "  1. Download from: https://wixtoolset.org/releases/" -ForegroundColor White
    Write-Host "  2. Install WiX Toolset" -ForegroundColor White
    Write-Host "  3. Restart PowerShell" -ForegroundColor White
    Write-Host "  4. Run this script again" -ForegroundColor White
    exit 1
}

Write-Host "✓ WiX Toolset found at: $wixPath" -ForegroundColor Green

# Build the Go executable first
if (-not $SkipBuild) {
    Write-Host ""
    Write-Host "Building Obelisk executable..." -ForegroundColor Yellow
    
    if (-not (Test-Path "main.go")) {
        Write-Host "❌ main.go not found. Please run this script from the project root." -ForegroundColor Red
        exit 1
    }
    
    # Clean previous build
    if (Test-Path "bin") {
        Remove-Item -Path "bin" -Recurse -Force
    }
    New-Item -ItemType Directory -Path "bin" -Force | Out-Null
    
    # Build for Windows
    $env:GOOS = "windows"
    $env:GOARCH = "amd64"
    
    go build -o "bin/obelisk.exe" -ldflags "-s -w -X main.Version=$Version" .
    
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ Build failed!" -ForegroundColor Red
        exit 1
    }
    
    Write-Host "✓ Build successful: bin/obelisk.exe" -ForegroundColor Green
} else {
    Write-Host ""
    Write-Host "⚠️  Skipping build (using existing bin/obelisk.exe)" -ForegroundColor Yellow
    
    if (-not (Test-Path "bin/obelisk.exe")) {
        Write-Host "❌ bin/obelisk.exe not found!" -ForegroundColor Red
        exit 1
    }
}

# Create installer directory if it doesn't exist
if (-not (Test-Path "installer")) {
    New-Item -ItemType Directory -Path "installer" -Force | Out-Null
}

# Check for required installer assets
Write-Host ""
Write-Host "Checking installer assets..." -ForegroundColor Yellow

$requiredFiles = @(
    "installer/obelisk.wxs",
    "installer/License.rtf"
)

$missingFiles = @()
foreach ($file in $requiredFiles) {
    if (-not (Test-Path $file)) {
        $missingFiles += $file
    }
}

if ($missingFiles.Count -gt 0) {
    Write-Host "❌ Missing required files:" -ForegroundColor Red
    foreach ($file in $missingFiles) {
        Write-Host "  - $file" -ForegroundColor Red
    }
    exit 1
}

Write-Host "✓ All required files present" -ForegroundColor Green

# Optional assets (will use defaults if missing)
if (-not (Test-Path "installer/icon.ico")) {
    Write-Host "⚠️  installer/icon.ico not found (using default)" -ForegroundColor Yellow
}
if (-not (Test-Path "installer/banner.bmp")) {
    Write-Host "⚠️  installer/banner.bmp not found (using default)" -ForegroundColor Yellow
}
if (-not (Test-Path "installer/dialog.bmp")) {
    Write-Host "⚠️  installer/dialog.bmp not found (using default)" -ForegroundColor Yellow
}

# Compile WiX source
Write-Host ""
Write-Host "Compiling WiX source..." -ForegroundColor Yellow

$candleExe = Join-Path $wixPath "candle.exe"
$lightExe = Join-Path $wixPath "light.exe"

# Run candle (compiler)
& $candleExe -nologo -out "installer/obelisk.wixobj" "installer/obelisk.wxs"

if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ WiX compilation failed!" -ForegroundColor Red
    exit 1
}

Write-Host "✓ WiX compilation successful" -ForegroundColor Green

# Run light (linker)
Write-Host ""
Write-Host "Linking MSI installer..." -ForegroundColor Yellow

& $lightExe -nologo `
    -ext WixUIExtension `
    -cultures:en-us `
    -out "installer/ObeliskCLI-$Version-x64.msi" `
    "installer/obelisk.wixobj"

if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ MSI linking failed!" -ForegroundColor Red
    exit 1
}

# Clean up intermediate files
Remove-Item "installer/obelisk.wixobj" -Force -ErrorAction SilentlyContinue

Write-Host ""
Write-Host "✅ Installer built successfully!" -ForegroundColor Green
Write-Host ""
Write-Host "Output: installer/ObeliskCLI-$Version-x64.msi" -ForegroundColor Cyan
Write-Host ""
Write-Host "File size: $((Get-Item "installer/ObeliskCLI-$Version-x64.msi").Length / 1MB) MB" -ForegroundColor Gray
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Yellow
Write-Host "  1. Test the installer: .\installer\ObeliskCLI-$Version-x64.msi" -ForegroundColor White
Write-Host "  2. Distribute via GitHub Releases" -ForegroundColor White
Write-Host "  3. Users can double-click to install" -ForegroundColor White

# Made with Bob
