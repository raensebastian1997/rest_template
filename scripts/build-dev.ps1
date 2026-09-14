$ErrorActionPreference = "Stop"

$projectRoot = Split-Path -Parent $PSScriptRoot
$outputDirectory = Join-Path $projectRoot ".air-tmp"
$binaryName = "api-dev-{0}.exe" -f [DateTimeOffset]::UtcNow.ToUnixTimeMilliseconds()
$binaryPath = Join-Path $outputDirectory $binaryName
$pointerPath = Join-Path $outputDirectory "current.txt"

[System.IO.Directory]::CreateDirectory($outputDirectory) | Out-Null

$buildArguments = @(
    "build"
    "-tags=nomsgpack"
    "-buildmode=exe"
    "-buildvcs=false"
    "-ldflags=-s -w -buildid="
    "-o"
    $binaryPath
    "./cmd/api"
)

Push-Location $projectRoot
try {
    & go @buildArguments
    if ($LASTEXITCODE -ne 0) {
        exit $LASTEXITCODE
    }
} finally {
    Pop-Location
}

[System.IO.File]::WriteAllText(
    $pointerPath,
    $binaryPath,
    [System.Text.UTF8Encoding]::new($false)
)

Get-ChildItem -LiteralPath $outputDirectory -Filter "api-dev-*.exe" -File |
    Where-Object { $_.FullName -ne $binaryPath } |
    ForEach-Object {
        try {
            Remove-Item -LiteralPath $_.FullName -Force -ErrorAction Stop
        } catch {
            # Windows may briefly retain the previous executable. It will be
            # removed during a later build or when Air exits.
        }
    }
