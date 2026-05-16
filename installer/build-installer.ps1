# Obelisk CLI MSI Installer Build Script
# Requires WiX Toolset 3.11+ to be installed
# Download from: https://wixtoolset.org/releases/

param(
    [string]$Version = "0.1.0",
    [switch]$SkipBuild = $false
)

$ErrorActionPreference = "Stop"

Write-Host "Obelisk CLI Installer Builder" -ForegroundColor Cyan
Write-Host "=================================" -ForegroundColor Cyan
Write-Host ""

# Ensure we are running from the project root
$projectRoot = Split-Path -Parent (Split-Path -Parent $MyInvocation.MyCommand.Path)
Push-Location $projectRoot

try {
    # -- Check if WiX is installed ------------------------------------------------
    $wixPath = "${env:WIX}bin"
    if (-not (Test-Path $wixPath)) {
        # Fallback: try finding candle.exe on PATH
        $candleOnPath = Get-Command candle.exe -ErrorAction SilentlyContinue
        if ($candleOnPath) {
            $wixPath = Split-Path $candleOnPath.Source
        } else {
            Write-Host "[ERROR] WiX Toolset not found!" -ForegroundColor Red
            Write-Host ""
            Write-Host "Please install WiX Toolset 3.11 or later:" -ForegroundColor Yellow
            Write-Host "  Option A: choco install wixtoolset" -ForegroundColor White
            Write-Host "  Option B: Download from https://wixtoolset.org/releases/" -ForegroundColor White
            Write-Host ""
            Write-Host "After installing, restart PowerShell and run this script again." -ForegroundColor Yellow
            exit 1
        }
    }

    Write-Host "[OK] WiX Toolset found at: $wixPath" -ForegroundColor Green

    # -- Build the Go executable ---------------------------------------------------
    $ldflagsPkg = "github.com/Swif7ify/Obelisk-CLI/cmd"
    $commitHash = git rev-parse --short HEAD 2>$null
    if (-not $commitHash) {
        $commitHash = "unknown"
    }
    $buildDate = Get-Date -Format "yyyy-MM-ddTHH:mm:ssZ"
    $ldflags = "-s -w -X $ldflagsPkg.Version=$Version -X $ldflagsPkg.Commit=$commitHash -X $ldflagsPkg.BuildDate=$buildDate"

    if (-not $SkipBuild) {
        Write-Host ""
        Write-Host "Building Obelisk executable..." -ForegroundColor Yellow

        if (-not (Test-Path "main.go")) {
            Write-Host "[ERROR] main.go not found. Please run this script from the project root." -ForegroundColor Red
            exit 1
        }

        # Clean previous build
        if (Test-Path "bin") {
            Remove-Item -Path "bin" -Recurse -Force
        }
        New-Item -ItemType Directory -Path "bin" -Force | Out-Null

        # Build for Windows x64
        $env:GOOS = "windows"
        $env:GOARCH = "amd64"

        go build -o "bin\obelisk.exe" -ldflags $ldflags .

        if ($LASTEXITCODE -ne 0) {
            Write-Host "[ERROR] Build failed!" -ForegroundColor Red
            exit 1
        }

        $exeSize = [math]::Round((Get-Item "bin\obelisk.exe").Length / 1MB, 2)
        Write-Host "[OK] Build successful: bin\obelisk.exe - $exeSize MB" -ForegroundColor Green
    } else {
        Write-Host ""
        Write-Host "[WARN] Skipping build (using existing bin\obelisk.exe)" -ForegroundColor Yellow

        if (-not (Test-Path "bin\obelisk.exe")) {
            Write-Host "[ERROR] bin\obelisk.exe not found! Run without -SkipBuild first." -ForegroundColor Red
            exit 1
        }
    }

    # -- Verify required installer assets ------------------------------------------
    Write-Host ""
    Write-Host "Checking installer assets..." -ForegroundColor Yellow

    $requiredFiles = @(
        "installer\obelisk.wxs",
        "installer\License.rtf",
        "README.md",
        "LICENSE",
        "CHANGELOG.md"
    )

    $missingFiles = @()
    foreach ($file in $requiredFiles) {
        if (-not (Test-Path $file)) {
            $missingFiles += $file
        }
    }

    if ($missingFiles.Count -gt 0) {
        Write-Host "[ERROR] Missing required files:" -ForegroundColor Red
        foreach ($file in $missingFiles) {
            Write-Host "  - $file" -ForegroundColor Red
        }
        exit 1
    }

    Write-Host "[OK] All required files present" -ForegroundColor Green

    # -- Compile WiX source --------------------------------------------------------
    Write-Host ""
    Write-Host "Compiling WiX source..." -ForegroundColor Yellow

    $candleExe = Join-Path $wixPath "candle.exe"
    $lightExe  = Join-Path $wixPath "light.exe"
    $binDir    = Join-Path $projectRoot "bin"

    # Run candle (compiler) - pass variables for portable paths
    & $candleExe -nologo `
        -arch x64 `
        "-dBinDir=$binDir" `
        "-dProjectDir=$projectRoot" `
        -out "installer\obelisk.wixobj" `
        "installer\obelisk.wxs"

    if ($LASTEXITCODE -ne 0) {
        Write-Host "[ERROR] WiX compilation failed!" -ForegroundColor Red
        exit 1
    }

    Write-Host "[OK] WiX compilation successful" -ForegroundColor Green

    # -- Link MSI installer --------------------------------------------------------
    Write-Host ""
    Write-Host "Linking MSI installer..." -ForegroundColor Yellow

    $msiOutput = "installer\ObeliskCLI-$Version-x64.msi"

    & $lightExe -nologo `
        -ext WixUIExtension `
        -cultures:en-us `
        -out $msiOutput `
        "installer\obelisk.wixobj"

    if ($LASTEXITCODE -ne 0) {
        Write-Host "[ERROR] MSI linking failed!" -ForegroundColor Red
        exit 1
    }

    # -- Clean up intermediate files -----------------------------------------------
    Remove-Item "installer\obelisk.wixobj" -Force -ErrorAction SilentlyContinue
    Remove-Item "installer\ObeliskCLI-$Version-x64.wixpdb" -Force -ErrorAction SilentlyContinue

    # -- Done ----------------------------------------------------------------------
    $msiSize = [math]::Round((Get-Item $msiOutput).Length / 1MB, 2)
    Write-Host ""
    Write-Host "[SUCCESS] Installer built successfully!" -ForegroundColor Green
    Write-Host ""
    Write-Host "  Output:    $msiOutput" -ForegroundColor Cyan
    Write-Host "  File size: $msiSize MB" -ForegroundColor Gray
    Write-Host ""
    Write-Host "Next steps:" -ForegroundColor Yellow
    Write-Host "  1. Test the installer:  .\$msiOutput" -ForegroundColor White
    Write-Host "  2. Upload to GitHub Releases" -ForegroundColor White
    Write-Host "  3. Users download and double-click to install" -ForegroundColor White
}
finally {
    Pop-Location
}
