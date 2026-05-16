# Obelisk CLI Uninstaller for Windows
# Run this script with: powershell -ExecutionPolicy Bypass -File uninstall.ps1

$ErrorActionPreference = "Stop"

Write-Host "Obelisk CLI Uninstaller" -ForegroundColor Cyan
Write-Host ""

# Check both possible install locations
$systemPath = "C:\Program Files\Obelisk"
$userPath = "$env:LOCALAPPDATA\Obelisk"

$installPath = $null
$envScope = "User"
if (Test-Path $systemPath) {
    $installPath = $systemPath
    $envScope = "Machine"
} elseif (Test-Path $userPath) {
    $installPath = $userPath
    $envScope = "User"
}

if ($null -eq $installPath) {
    Write-Host "[OK] Obelisk is not installed" -ForegroundColor Green
    exit 0
}

Write-Host "Found installation at: $installPath" -ForegroundColor Yellow

# Remove from PATH
$currentPath = [Environment]::GetEnvironmentVariable("Path", $envScope)
if ($currentPath -like "*$installPath*") {
    $newPath = ($currentPath -split ';' | Where-Object { $_ -ne $installPath }) -join ';'
    [Environment]::SetEnvironmentVariable("Path", $newPath, $envScope)
    Write-Host "[OK] Removed from PATH" -ForegroundColor Green
}

# Remove directory
Remove-Item -Path $installPath -Recurse -Force
Write-Host "[OK] Removed installation directory" -ForegroundColor Green

# --- User Data Cleanup ---
Write-Host ""
$configDir = Join-Path $env:USERPROFILE ".obelisk"
$hasUserData = Test-Path $configDir

if ($hasUserData) {
    Write-Host "User data found at: $configDir" -ForegroundColor Yellow
    Write-Host "This includes your configuration and saved preferences." -ForegroundColor Gray
    Write-Host ""
    Write-Host "  [1] Keep user data   (recommended if you plan to reinstall)" -ForegroundColor White
    Write-Host "  [2] Delete everything (removes config, API keys, all user data)" -ForegroundColor White
    Write-Host ""

    $choice = Read-Host "Choose an option [1/2] (default: 1)"

    if ($choice -eq "2") {
        # Remove config directory
        Remove-Item -Path $configDir -Recurse -Force -ErrorAction SilentlyContinue
        Write-Host "[OK] Removed user data directory: $configDir" -ForegroundColor Green

        # Remove OS keyring entry for Gemini API key
        try {
            cmdkey /delete:obelisk-cli 2>$null | Out-Null
            Write-Host "[OK] Removed stored API key from Windows Credential Manager" -ForegroundColor Green
        } catch {
            # Silently ignore if not found
        }

        Write-Host "[OK] All user data has been removed" -ForegroundColor Green
    } else {
        Write-Host "[OK] User data retained at: $configDir" -ForegroundColor Green
    }
} else {
    Write-Host "[OK] No user data found" -ForegroundColor Gray
}

Write-Host ""
Write-Host "[SUCCESS] Uninstallation complete!" -ForegroundColor Green
Write-Host "  Please restart your terminal for PATH changes to take effect" -ForegroundColor Yellow
