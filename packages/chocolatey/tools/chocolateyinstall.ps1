# Chocolatey install script for Obelisk CLI

$ErrorActionPreference = 'Stop'

$packageName = 'obelisk-cli'
$toolsDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$url64 = 'https://github.com/Swif7ify/Obelisk-CLI/releases/download/v0.1.0/obelisk-windows-amd64.exe'

$packageArgs = @{
  packageName   = $packageName
  unzipLocation = $toolsDir
  fileType      = 'exe'
  url64bit      = $url64
  
  softwareName  = 'Obelisk CLI*'
  
  checksum64    = ''  # Will be filled during package creation
  checksumType64= 'sha256'
  
  silentArgs    = ""
  validExitCodes= @(0)
}

# Download and install
Install-ChocolateyPackage @packageArgs

# Rename to obelisk.exe
$exePath = Join-Path $toolsDir 'obelisk-windows-amd64.exe'
$targetPath = Join-Path $toolsDir 'obelisk.exe'

if (Test-Path $exePath) {
  Rename-Item -Path $exePath -NewName 'obelisk.exe' -Force
}

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
