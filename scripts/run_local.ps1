$ErrorActionPreference = "Stop"

function Import-DotEnv {
    param(
        [Parameter(Mandatory = $true)][string]$Path
    )
    if (-not (Test-Path -LiteralPath $Path)) { return }

    Get-Content -LiteralPath $Path | ForEach-Object {
        $line = $_.Trim().TrimStart([char]0xFEFF)
        if (-not $line) { return }
        if ($line.StartsWith("#")) { return }
        if ($line.StartsWith(";")) { return }

        $idx = $line.IndexOf("=")
        if ($idx -lt 1) { return }

        $key = $line.Substring(0, $idx).Trim()
        $val = $line.Substring($idx + 1).Trim()

        if (-not $key) { return }

        # Strip optional quotes
        if (($val.StartsWith('"') -and $val.EndsWith('"')) -or ($val.StartsWith("'") -and $val.EndsWith("'"))) {
            if ($val.Length -ge 2) { $val = $val.Substring(1, $val.Length - 2) }
        }

        # Do not override variables already set in the shell/session
        if ([string]::IsNullOrEmpty([Environment]::GetEnvironmentVariable($key, "Process"))) {
            Set-Item -Path "Env:$key" -Value $val
        }
    }
}

$repoRoot = Resolve-Path (Join-Path $PSScriptRoot "..")
Set-Location $repoRoot

Write-Host "Starting Alpha Strike Helper (local mode)..." -ForegroundColor Cyan
Write-Host "Repo: $repoRoot" -ForegroundColor DarkGray

$dotEnvPath = Join-Path $repoRoot ".env"
Import-DotEnv -Path $dotEnvPath

# Default envs for local launch (can be overridden before script call).
if (-not $env:DB_HOST) { $env:DB_HOST = "localhost" }
if (-not $env:DB_PORT) { $env:DB_PORT = "5432" }
if (-not $env:DB_USER) { $env:DB_USER = "postgres" }
if (-not $env:DB_PASSWORD) { $env:DB_PASSWORD = "change-me" }
if (-not $env:DB_NAME) { $env:DB_NAME = "alpha_strike" }
if (-not $env:SERVER_PORT) { $env:SERVER_PORT = "8080" }
if (-not $env:JWT_SECRET) { $env:JWT_SECRET = "change-me-local" }

Write-Host "DB: $($env:DB_HOST):$($env:DB_PORT) / $($env:DB_NAME)" -ForegroundColor Yellow
Write-Host "Server port: $($env:SERVER_PORT)" -ForegroundColor Yellow

go run ./cmd/server
