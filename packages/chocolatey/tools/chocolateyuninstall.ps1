# Chocolatey uninstall script for Obelisk CLI

$ErrorActionPreference = 'Stop'

$packageName = 'obelisk-cli'
$toolsDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"

# Remove the executable
$exePath = Join-Path $toolsDir 'obelisk.exe'
if (Test-Path $exePath) {
  Remove-Item $exePath -Force
  Write-Host "Removed obelisk.exe" -ForegroundColor Green
}

Write-Host "Obelisk CLI uninstalled successfully!" -ForegroundColor Green
Write-Host ""
Write-Host "Note: Configuration files in ~/.obelisk/ were preserved." -ForegroundColor Yellow
Write-Host "To remove them manually: Remove-Item -Recurse -Force ~/.obelisk" -ForegroundColor Gray

# Made with Bob
