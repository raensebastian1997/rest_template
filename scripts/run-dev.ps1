$ErrorActionPreference = "Stop"

$projectRoot = Split-Path -Parent $PSScriptRoot
$pointerPath = Join-Path $projectRoot ".air-tmp\current.txt"

if ([string]::IsNullOrWhiteSpace($env:JWT_SECRET)) {
    $env:JWT_SECRET = "development-only-jwt-secret-change-me"
}

if (-not (Test-Path -LiteralPath $pointerPath)) {
    throw "Development binary pointer was not created: $pointerPath"
}

$binaryPath = [System.IO.File]::ReadAllText($pointerPath).Trim()
if (-not (Test-Path -LiteralPath $binaryPath)) {
    throw "Development binary was not found: $binaryPath"
}

& $binaryPath
exit $LASTEXITCODE
