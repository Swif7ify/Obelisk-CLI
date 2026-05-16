# Obelisk CLI Extension Build Script
# Builds and packages the VS Code extension

param(
    [switch]$Install = $false,
    [switch]$Publish = $false,
    [string]$Registry = "vscode"  # "vscode" or "openvsx"
)

$ErrorActionPreference = "Stop"

Write-Host "Obelisk CLI Extension Builder" -ForegroundColor Cyan
Write-Host "==============================" -ForegroundColor Cyan
Write-Host ""

# Ensure we're in the extension directory
$extensionDir = Split-Path -Parent $MyInvocation.MyCommand.Path
Push-Location $extensionDir

try {
    # Check if Node.js is installed
    $nodeVersion = node --version 2>$null
    if (-not $nodeVersion) {
        Write-Host "ERROR: Node.js not found!" -ForegroundColor Red
        Write-Host "Please install Node.js from: https://nodejs.org/" -ForegroundColor Yellow
        exit 1
    }
    Write-Host "Node.js version: $nodeVersion" -ForegroundColor Green

    # Install dependencies
    Write-Host ""
    Write-Host "Installing dependencies..." -ForegroundColor Yellow
    npm install
    if ($LASTEXITCODE -ne 0) {
        Write-Host "ERROR: npm install failed!" -ForegroundColor Red
        exit 1
    }
    Write-Host "Dependencies installed" -ForegroundColor Green

    # Compile TypeScript
    Write-Host ""
    Write-Host "Compiling TypeScript..." -ForegroundColor Yellow
    npm run compile
    if ($LASTEXITCODE -ne 0) {
        Write-Host "ERROR: TypeScript compilation failed!" -ForegroundColor Red
        exit 1
    }
    Write-Host "TypeScript compiled successfully" -ForegroundColor Green

    # Package extension
    Write-Host ""
    Write-Host "Packaging extension..." -ForegroundColor Yellow
    
    # Check if vsce is installed
    $vsceInstalled = Get-Command vsce -ErrorAction SilentlyContinue
    if (-not $vsceInstalled) {
        Write-Host "Installing @vscode/vsce..." -ForegroundColor Yellow
        npm install -g @vscode/vsce
    }

    npm run package
    if ($LASTEXITCODE -ne 0) {
        Write-Host "ERROR: Extension packaging failed!" -ForegroundColor Red
        exit 1
    }

    # Find the generated .vsix file
    $vsixFile = Get-ChildItem -Filter "*.vsix" | Select-Object -First 1
    if (-not $vsixFile) {
        Write-Host "ERROR: No .vsix file found!" -ForegroundColor Red
        exit 1
    }

    $vsixSize = [math]::Round($vsixFile.Length / 1KB, 2)
    Write-Host "Extension packaged successfully!" -ForegroundColor Green
    Write-Host "  File: $($vsixFile.Name)" -ForegroundColor Cyan
    Write-Host "  Size: $vsixSize KB" -ForegroundColor Gray

    # Install locally if requested
    if ($Install) {
        Write-Host ""
        Write-Host "Installing extension locally..." -ForegroundColor Yellow
        code --install-extension $vsixFile.FullName --force
        if ($LASTEXITCODE -eq 0) {
            Write-Host "Extension installed successfully!" -ForegroundColor Green
            Write-Host "Restart VS Code to activate the extension." -ForegroundColor Yellow
        } else {
            Write-Host "WARNING: Installation may have failed. Check VS Code." -ForegroundColor Yellow
        }
    }

    # Publish if requested
    if ($Publish) {
        Write-Host ""
        Write-Host "Publishing extension..." -ForegroundColor Yellow
        
        if ($Registry -eq "vscode") {
            Write-Host "Publishing to VS Code Marketplace..." -ForegroundColor Cyan
            vsce publish
            if ($LASTEXITCODE -eq 0) {
                Write-Host "Published to VS Code Marketplace!" -ForegroundColor Green
            } else {
                Write-Host "ERROR: Publishing failed!" -ForegroundColor Red
                Write-Host "Make sure you're logged in: vsce login <publisher>" -ForegroundColor Yellow
                exit 1
            }
        } elseif ($Registry -eq "openvsx") {
            Write-Host "Publishing to Open VSX Registry..." -ForegroundColor Cyan
            
            # Check if ovsx is installed
            $ovsxInstalled = Get-Command ovsx -ErrorAction SilentlyContinue
            if (-not $ovsxInstalled) {
                Write-Host "Installing ovsx..." -ForegroundColor Yellow
                npm install -g ovsx
            }

            $token = $env:OVSX_TOKEN
            if (-not $token) {
                Write-Host "ERROR: OVSX_TOKEN environment variable not set!" -ForegroundColor Red
                Write-Host "Set it with: `$env:OVSX_TOKEN = 'your-token'" -ForegroundColor Yellow
                exit 1
            }

            ovsx publish $vsixFile.FullName -p $token
            if ($LASTEXITCODE -eq 0) {
                Write-Host "Published to Open VSX Registry!" -ForegroundColor Green
            } else {
                Write-Host "ERROR: Publishing failed!" -ForegroundColor Red
                exit 1
            }
        } else {
            Write-Host "ERROR: Invalid registry '$Registry'" -ForegroundColor Red
            Write-Host "Use 'vscode' or 'openvsx'" -ForegroundColor Yellow
            exit 1
        }
    }

    # Summary
    Write-Host ""
    Write-Host "Build complete!" -ForegroundColor Green
    Write-Host ""
    Write-Host "Next steps:" -ForegroundColor Yellow
    if (-not $Install) {
        Write-Host "  1. Install locally:  .\build-extension.ps1 -Install" -ForegroundColor White
    }
    if (-not $Publish) {
        Write-Host "  2. Publish to VS Code Marketplace:  .\build-extension.ps1 -Publish -Registry vscode" -ForegroundColor White
        Write-Host "  3. Publish to Open VSX (IBM Bob):   .\build-extension.ps1 -Publish -Registry openvsx" -ForegroundColor White
    }
    Write-Host "  4. Share .vsix file:  $($vsixFile.Name)" -ForegroundColor White
    Write-Host ""

} finally {
    Pop-Location
}

# Made with Bob
