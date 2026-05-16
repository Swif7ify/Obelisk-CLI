# Obelisk CLI Installer for Windows
# Run this script with: powershell -ExecutionPolicy Bypass -File install.ps1

$ErrorActionPreference = "Stop"

Write-Host "🏛️  Obelisk CLI Installer" -ForegroundColor Cyan
Write-Host ""

# Check if running as administrator
$isAdmin = ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)

# Determine install location
if ($isAdmin) {
    $installPath = "C:\Program Files\Obelisk"
    Write-Host "Installing to: $installPath (System-wide)" -ForegroundColor Green
} else {
    $installPath = "$env:LOCALAPPDATA\Obelisk"
    Write-Host "Installing to: $installPath (User-only)" -ForegroundColor Yellow
    Write-Host "Tip: Run as Administrator for system-wide installation" -ForegroundColor Gray
}

Write-Host ""

# Create install directory
if (!(Test-Path $installPath)) {
    New-Item -ItemType Directory -Path $installPath -Force | Out-Null
    Write-Host "✓ Created directory: $installPath" -ForegroundColor Green
}

# Copy executable
$exePath = ".\bin\obelisk.exe"
if (!(Test-Path $exePath)) {
    Write-Host "✗ Error: obelisk.exe not found at $exePath" -ForegroundColor Red
    Write-Host "  Please run 'go build -o bin/obelisk.exe' first" -ForegroundColor Yellow
    exit 1
}

Copy-Item $exePath "$installPath\obelisk.exe" -Force
Write-Host "✓ Copied obelisk.exe to $installPath" -ForegroundColor Green

# Add to PATH
$currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($currentPath -notlike "*$installPath*") {
    [Environment]::SetEnvironmentVariable(
        "Path",
        "$currentPath;$installPath",
        "User"
    )
    Write-Host "✓ Added to PATH" -ForegroundColor Green
    Write-Host ""
    Write-Host "⚠️  Please restart your terminal for PATH changes to take effect" -ForegroundColor Yellow
} else {
    Write-Host "✓ Already in PATH" -ForegroundColor Green
}

Write-Host ""
Write-Host "✅ Installation complete!" -ForegroundColor Green
Write-Host ""
Write-Host "Usage:" -ForegroundColor Cyan
Write-Host "  obelisk              # Launch interactive TUI" -ForegroundColor White
Write-Host "  obelisk check        # Run health check" -ForegroundColor White
Write-Host "  obelisk scan         # Headless scan" -ForegroundColor White
Write-Host "  obelisk config       # Manage configuration" -ForegroundColor White
Write-Host "  obelisk --help       # Show all commands" -ForegroundColor White
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Cyan
Write-Host "  1. Restart your terminal" -ForegroundColor White
Write-Host "  2. Run: obelisk config set api-key YOUR_KEY" -ForegroundColor White
Write-Host "  3. Run: obelisk" -ForegroundColor White

# Made with Bob
