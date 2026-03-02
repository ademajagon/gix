$ErrorActionPreference = 'Stop'

$toolsDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$binary   = Join-Path $toolsDir 'gix.exe'

if (Test-Path $binary) {
    Remove-Item $binary -Force
    Write-Host "gix uninstalled." -ForegroundColor Yellow
}