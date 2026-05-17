# Chocolatey install script for Obelisk CLI

$ErrorActionPreference = 'Stop'

$packageName = 'obelisk-cli'
$toolsDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$url64 = 'https://github.com/Swif7ify/Obelisk-CLI/releases/download/v1.1.0/obelisk.exe'
$exePath = Join-Path $toolsDir 'obelisk.exe'

# Download the portable executable directly into the tools directory
# Chocolatey will automatically create a shim for any .exe placed in this folder
Get-ChocolateyWebFile -PackageName $packageName `
                      -FileFullPath $exePath `
                      -Url64bit $url64 `
                      -Url $url64

Write-Host "Obelisk CLI installed successfully!" -ForegroundColor Green
Write-Host ""
Write-Host "Usage:" -ForegroundColor Cyan
Write-Host "  obelisk              # Launch interactive TUI" -ForegroundColor White
Write-Host "  obelisk check        # Run health check" -ForegroundColor White
Write-Host "  obelisk scan         # Headless scan" -ForegroundColor White
Write-Host "  obelisk --help       # Show all commands" -ForegroundColor White
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Cyan
Write-Host "  1. Set your Gemini API key: obelisk config set api-key YOUR_KEY" -ForegroundColor White
Write-Host "  2. Run: obelisk" -ForegroundColor White

# Made with Bob
