# Obelisk CLI Cleanup Script
# Called by the MSI installer during uninstall when user opts to remove all data.
# This script removes:
#   1. The ~/.obelisk configuration directory (config.json)
#   2. The OS keyring entry for the Gemini API key

param(
    [switch]$RemoveUserData = $false
)

$ErrorActionPreference = "SilentlyContinue"

if (-not $RemoveUserData) {
    exit 0
}

# --- Remove configuration directory ---
$configDir = Join-Path $env:USERPROFILE ".obelisk"
if (Test-Path $configDir) {
    Remove-Item -Path $configDir -Recurse -Force
}

# --- Remove OS keyring entry for Gemini API key ---
# Windows Credential Manager stores keyring entries as Generic Credentials.
# The go-keyring library stores them with the service name as the target.
try {
    # Use cmdkey to remove the stored credential
    cmdkey /delete:obelisk-cli 2>$null
} catch {
    # Silently ignore if not found
}

exit 0
