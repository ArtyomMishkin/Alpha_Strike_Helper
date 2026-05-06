$ErrorActionPreference = "Stop"

function Import-DotEnv {
    param(
        [Parameter(Mandatory = $true)][string]$Path
    )
    if (-not (Test-Path -LiteralPath $Path)) { return }

    Get-Content -LiteralPath $Path | ForEach-Object {
        $line = $_.Trim()
        if (-not $line) { return }
        if ($line.StartsWith("#")) { return }
        if ($line.StartsWith(";")) { return }

        $idx = $line.IndexOf("=")
        if ($idx -lt 1) { return }

        $key = $line.Substring(0, $idx).Trim()
        $val = $line.Substring($idx + 1).Trim()

        if (-not $key) { return }

        if (($val.StartsWith('"') -and $val.EndsWith('"')) -or ($val.StartsWith("'") -and $val.EndsWith("'"))) {
            if ($val.Length -ge 2) { $val = $val.Substring(1, $val.Length - 2) }
        }

        if ([string]::IsNullOrEmpty([Environment]::GetEnvironmentVariable($key, "Process"))) {
            Set-Item -Path "Env:$key" -Value $val
        }
    }
}

$repoRoot = Resolve-Path (Join-Path $PSScriptRoot "..")
Set-Location $repoRoot

Write-Host "Importing MUL data (BattleMech + Vehicle + Infantry + Aerospace)..." -ForegroundColor Cyan
Write-Host "Repo: $repoRoot" -ForegroundColor DarkGray

$dotEnvPath = Join-Path $repoRoot ".env"
Import-DotEnv -Path $dotEnvPath

# Default envs for local import (can be overridden before script call).
if (-not $env:DB_HOST) { $env:DB_HOST = "localhost" }
if (-not $env:DB_PORT) { $env:DB_PORT = "5432" }
if (-not $env:DB_USER) { $env:DB_USER = "postgres" }
if (-not $env:DB_PASSWORD) { $env:DB_PASSWORD = "change-me" }
if (-not $env:DB_NAME) { $env:DB_NAME = "alpha_strike" }
if (-not $env:JWT_SECRET) { $env:JWT_SECRET = "change-me-local" }

# 18 BattleMech, 19 Combat Vehicle, 17 Aerospace, 21 Infantry
go run ./cmd/masterunitlist_sync --http-timeout=180s --unit-type-id=18 --replace=true --include-faction-eras=true
go run ./cmd/masterunitlist_sync --http-timeout=180s --unit-type-id=19 --replace=false --include-faction-eras=true
go run ./cmd/masterunitlist_sync --http-timeout=180s --unit-type-id=17 --replace=false --include-faction-eras=true
go run ./cmd/masterunitlist_sync --http-timeout=180s --unit-type-id=21 --replace=false --include-faction-eras=true

Write-Host "Import finished." -ForegroundColor Green
