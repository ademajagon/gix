$ErrorActionPreference = 'Stop'

$version   = '0.3.0'
$repo      = 'ademajagon/gix'
$arch      = if ([System.Environment]::Is64BitOperatingSystem -and $env:PROCESSOR_ARCHITECTURE -eq 'ARM64') { 'arm64' } else { 'amd64' }
$assetName = "gix-windows-$arch.exe"
$url       = "https://github.com/$repo/releases/download/v$version/$assetName"

$checksums = @{
    'amd64' = 'REPLACE_WITH_SHA256_AMD64'
    'arm64' = 'REPLACE_WITH_SHA256_ARM64'
}

$packageArgs = @{
    packageName    = 'gix'
    unzipLocation  = $null
    fileType       = 'exe'
    url64bit       = $url
    softwareName   = 'gix*'
    checksum64     = $checksums[$arch]
    checksumType64 = 'sha256'
}

$toolsDir   = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$destFile   = Join-Path $toolsDir 'gix.exe'

Get-ChocolateyWebFile `
  -PackageName  $packageArgs.packageName `
  -FileFullPath $destFile `
  -Url64bit     $packageArgs.url64bit `
  -Checksum64   $packageArgs.checksum64 `
  -ChecksumType $packageArgs.checksumType64

Unblock-File -Path $destFile

Write-Host ""
Write-Host "gix $version installed. Run: gix --help" -ForegroundColor Green