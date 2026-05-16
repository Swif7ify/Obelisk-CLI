# Obelisk CLI Uninstaller for Windows
# Run this script with: powershell -ExecutionPolicy Bypass -File uninstall.ps1

$ErrorActionPreference = "Stop"

Write-Host "🏛️  Obelisk CLI Uninstaller" -ForegroundColor Cyan
Write-Host ""

# Check both possible install locations
$systemPath = "C:\Program Files\Obelisk"
$userPath = "$env:LOCALAPPDATA\Obelisk"

$installPath = $null
if (Test-Path $systemPath) {
    $installPath = $systemPath
} elseif (Test-Path $userPath) {
    $installPath = $userPath
}

if ($null -eq $installPath) {
    Write-Host "✓ Obelisk is not installed" -ForegroundColor Green
    exit 0
}

Write-Host "Found installation at: $installPath" -ForegroundColor Yellow

# Remove from PATH
$currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($currentPath -like "*$installPath*") {
    $newPath = ($currentPath -split ';' | Where-Object { $_ -ne $installPath }) -join ';'
    [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
    Write-Host "✓ Removed from PATH" -ForegroundColor Green
}

# Remove directory
Remove-Item -Path $installPath -Recurse -Force
Write-Host "✓ Removed installation directory" -ForegroundColor Green

Write-Host ""
Write-Host "✅ Uninstallation complete!" -ForegroundColor Green
Write-Host "   Please restart your terminal for PATH changes to take effect" -ForegroundColor Yellow

# Made with Bob
